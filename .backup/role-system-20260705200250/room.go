package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"liars-bar/internal/game"
)

type GameEvent struct {
	Type     string          `json:"type"`
	PlayerID uint            `json:"player_id,omitempty"`
	AIPlayer bool            `json:"ai_player,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}

type GameRoom struct {
	ID          uint
	Name        string
	Players     map[uint]*game.Player
	State       *game.GameState
	Events      chan GameEvent
	Hub         *Hub
	mu          sync.RWMutex
	aiService   *AIProxy
	turnTimer   *time.Timer
	turnTimeout time.Duration
	closed      bool

	// game recording
	GameRecordID  uint
	RoundNo       int
	TurnNo        int
	statsRecorded int32 // accessed atomically; ensures OnGameOver fires once
}

func NewGameRoom(id uint, name string, hub *Hub) *GameRoom {
	room := &GameRoom{
		ID:          id,
		Name:        name,
		Players:     make(map[uint]*game.Player),
		State:       &game.GameState{Phase: game.PhaseWaiting},
		Events:      make(chan GameEvent, 256),
		Hub:         hub,
		turnTimeout: 30 * time.Second,
	}
	go room.eventLoop()
	return room
}

func (r *GameRoom) eventLoop() {
	for evt := range r.Events {
		r.processEvent(evt)
	}
}

func (r *GameRoom) HandleEvent(evt GameEvent) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.closed {
		return
	}
	select {
	case r.Events <- evt:
	default:
		log.Printf("Room %d event channel full, dropping %s", r.ID, evt.Type)
	}
}

func (r *GameRoom) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	r.closed = true
	close(r.Events)
}

func (r *GameRoom) CanJoin(userID uint) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, exists := r.Players[userID]; exists {
		return true
	}
	if r.State.Phase == game.PhasePlaying || r.State.Phase == game.PhaseGameOver {
		return false
	}
	return len(r.Players) < 4
}

func (r *GameRoom) processEvent(evt GameEvent) {
	switch evt.Type {
	case "PLAYER_JOIN":
		r.handlePlayerJoin(evt)
	case "PLAYER_LEAVE":
		r.handlePlayerLeave(evt)
	case "PLAYER_READY":
		r.handlePlayerReady(evt)
	case "START_GAME":
		r.handleStartGame(evt)
	case "PLAY_CARD":
		r.handlePlayCard(evt)
	case "CHALLENGE":
		r.handleChallenge(evt)
	case "PASS":
		r.handlePass(evt)
	case "CHAT":
		r.handleChat(evt)
	case "AI_ACTION":
		r.handleAIAction(evt)
	case "RECONNECT":
		r.handleReconnect(evt)
	case "GAME_OVER":
		r.handleGameOver(evt)
	}
}

func (r *GameRoom) handlePlayerJoin(evt GameEvent) {
	if r.State.Phase == game.PhasePlaying || r.State.Phase == game.PhaseGameOver {
		r.sendError(evt.PlayerID, "game already started")
		return
	}

	nickname := fmt.Sprintf("Player %d", evt.PlayerID)
	if client := r.Hub.GetClient(evt.PlayerID); client != nil {
		client.RoomID = r.ID
		if client.Nickname != "" {
			nickname = client.Nickname
		}
	}

	if player, ok := r.Players[evt.PlayerID]; ok {
		player.IsOnline = true
		player.IsAI = false
		player.AITakeover = false
		player.Nickname = nickname
	} else {
		if len(r.Players) >= 4 {
			r.sendError(evt.PlayerID, "room is full")
			return
		}
		r.Players[evt.PlayerID] = &game.Player{
			ID:        evt.PlayerID,
			Nickname:  nickname,
			SeatIndex: r.nextSeatIndex(),
			IsAlive:   true,
			IsOnline:  true,
			IsReady:   false,
		}
	}

	r.broadcast(Message{
		Type: "PLAYER_JOINED",
		Payload: map[string]interface{}{
			"player_id": evt.PlayerID,
			"nickname":  nickname,
		},
	})
	r.broadcastRoomState()
}

func (r *GameRoom) handlePlayerLeave(evt GameEvent) {
	player, ok := r.Players[evt.PlayerID]
	if !ok {
		return
	}

	if r.State.Phase == game.PhasePlaying {
		// A player leaving mid-game ends the game immediately. The remaining
		// alive players (if any) are declared winners — but typically this
		// means the game is abandoned. We broadcast a game-over signal so the
		// frontend can show "玩家退出，游戏结束" and the room is destroyed once
		// everyone has left.
		player.IsOnline = false
		delete(r.Players, evt.PlayerID)

		// Pick a winner among remaining alive players (nil if none).
		var winnerID *uint
		for _, p := range r.Players {
			if p.IsAlive && !p.IsAI {
				id := p.ID
				winnerID = &id
				break
			}
		}
		r.State.Phase = game.PhaseGameOver
		r.State.WinnerID = winnerID

		r.broadcast(Message{
			Type: "PLAYER_LEFT",
			Payload: map[string]interface{}{
				"player_id":  evt.PlayerID,
				"nickname":   player.Nickname,
				"game_over":  true,
				"reason":     "玩家退出，游戏结束",
				"winner_id":  winnerID,
			},
		})

		// Record stats once.
		if atomic.CompareAndSwapInt32(&r.statsRecorded, 0, 1) {
			wid := uint(0)
			if r.State.WinnerID != nil {
				wid = *r.State.WinnerID
			}
			if r.Hub.OnGameOver != nil {
				r.Hub.OnGameOver(r.ID, wid, r.State.Players)
			}
		}

		// If no human players remain, destroy immediately.
		hasHuman := false
		for _, p := range r.Players {
			if !p.IsAI {
				hasHuman = true
				break
			}
		}
		if !hasHuman {
			r.Hub.DestroyRoom(r.ID)
		}
		return
	}

	// Pre-game leave: just remove the player.
	delete(r.Players, evt.PlayerID)
	r.broadcast(Message{
		Type: "PLAYER_LEFT",
		Payload: map[string]interface{}{
			"player_id": evt.PlayerID,
		},
	})
	if len(r.Players) == 0 {
		r.Hub.DestroyRoom(r.ID)
		return
	}
	r.broadcastRoomState()
}

func (r *GameRoom) handlePlayerReady(evt GameEvent) {
	if r.State.Phase != game.PhaseWaiting && r.State.Phase != game.PhaseMatched {
		return
	}
	player, ok := r.Players[evt.PlayerID]
	if !ok {
		return
	}
	player.IsReady = true
	r.broadcast(Message{
		Type: "PLAYER_READY",
		Payload: map[string]interface{}{
			"player_id": evt.PlayerID,
		},
	})
	r.broadcastRoomState()

	if len(r.Players) == 4 && r.allPlayersReady() {
		r.startGame()
	}
}

func (r *GameRoom) handleStartGame(evt GameEvent) {
	if len(r.Players) != 4 {
		return
	}
	r.startGame()
}

func (r *GameRoom) startGame() {
	players := make([]*game.Player, 0, len(r.Players))
	for _, p := range r.Players {
		players = append(players, p)
	}
	// Sort by SeatIndex so gs.Players[i].SeatIndex == i. The frontend treats
	// GameState.CurrentPlayer (a slice index) as a seat index, so the slice
	// must be in seat order — iterating r.Players (a map) gives random order.
	sort.Slice(players, func(i, j int) bool {
		return players[i].SeatIndex < players[j].SeatIndex
	})
	r.State.InitGame(players)
	r.State.Phase = game.PhasePlaying

	r.broadcast(Message{
		Type:    "GAME_STARTED",
		Payload: r.State.ToPublic(0),
	})

	for _, p := range r.Players {
		if client := r.Hub.GetClient(p.ID); client != nil && !p.IsAI {
			client.SendMessage(Message{
				Type:    "GAME_STATE",
				Payload: r.State.ToPublic(p.ID),
			})
		}
	}

	r.processAITurns()
}

func (r *GameRoom) nextSeatIndex() int {
	used := make(map[int]bool, len(r.Players))
	for _, p := range r.Players {
		used[p.SeatIndex] = true
	}
	for i := 0; i < 4; i++ {
		if !used[i] {
			return i
		}
	}
	return len(r.Players)
}

func (r *GameRoom) allPlayersReady() bool {
	for _, p := range r.Players {
		if !p.IsReady {
			return false
		}
	}
	return true
}

func (r *GameRoom) readyCount() int {
	count := 0
	for _, p := range r.Players {
		if p.IsReady {
			count++
		}
	}
	return count
}

func (r *GameRoom) roomStatePayload() map[string]interface{} {
	players := make([]*game.Player, 0, len(r.Players))
	for _, p := range r.Players {
		players = append(players, p)
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].SeatIndex < players[j].SeatIndex
	})

	publicPlayers := make([]map[string]interface{}, 0, len(players))
	for _, p := range players {
		publicPlayers = append(publicPlayers, map[string]interface{}{
			"id":         p.ID,
			"nickname":   p.Nickname,
			"is_ai":      p.IsAI,
			"is_online":  p.IsOnline,
			"is_ready":   p.IsReady,
			"seat_index": p.SeatIndex,
		})
	}

	return map[string]interface{}{
		"id":           r.ID,
		"name":         r.Name,
		"phase":        r.State.Phase,
		"players":      publicPlayers,
		"player_count": len(players),
		"max_players":  4,
		"ready_count":  r.readyCount(),
		"can_join":     r.State.Phase != game.PhasePlaying && r.State.Phase != game.PhaseGameOver && len(players) < 4,
	}
}

func (r *GameRoom) broadcastRoomState() {
	r.broadcast(Message{Type: "ROOM_STATE", Payload: r.roomStatePayload()})
}

func (r *GameRoom) sendError(playerID uint, msg string) {
	if client := r.Hub.GetClient(playerID); client != nil {
		client.SendMessage(Message{Type: "ERROR", Payload: map[string]interface{}{"msg": msg}})
	}
}

// finalizeGameIfOver broadcasts GAME_OVER and records stats exactly once when
// the phase has transitioned to PhaseGameOver. Returns true if the game is over.
func (r *GameRoom) finalizeGameIfOver() bool {
	if r.State.Phase != game.PhaseGameOver {
		return false
	}
	winnerID := uint(0)
	if r.State.WinnerID != nil {
		winnerID = *r.State.WinnerID
	}
	r.broadcast(Message{
		Type: "GAME_OVER",
		Payload: map[string]interface{}{
			"winner_id": winnerID,
		},
	})
	if r.State.WinnerID != nil && atomic.CompareAndSwapInt32(&r.statsRecorded, 0, 1) {
		if r.Hub.OnGameOver != nil {
			r.Hub.OnGameOver(r.ID, *r.State.WinnerID, r.State.Players)
		}
	}
	return true
}

func (r *GameRoom) handlePlayCard(evt GameEvent) {
	if r.State.Phase != game.PhasePlaying {
		return
	}

	payload := struct {
		CardIDs []int  `json:"card_ids"`
		Claim   string `json:"claim"`
	}{}
	if err := json.Unmarshal(evt.Payload, &payload); err != nil {
		return
	}

	err := r.State.PlayCard(evt.PlayerID, payload.CardIDs, game.Card(payload.Claim))
	if err != nil {
		r.Hub.GetClient(evt.PlayerID).SendMessage(Message{
			Type:    "ERROR",
			Payload: map[string]interface{}{"msg": err.Error()},
		})
		return
	}

	r.broadcastGameState()
	if r.finalizeGameIfOver() {
		return
	}
	r.processAITurns()
}

func (r *GameRoom) handleChallenge(evt GameEvent) {
	if r.State.Phase != game.PhasePlaying {
		return
	}

	payload := struct {
		TargetPlayerID uint `json:"target_player_id"`
	}{}
	json.Unmarshal(evt.Payload, &payload)

	result, err := r.State.Challenge(evt.PlayerID)
	if err != nil {
		r.Hub.GetClient(evt.PlayerID).SendMessage(Message{
			Type:    "ERROR",
			Payload: map[string]interface{}{"msg": err.Error()},
		})
		return
	}

	r.broadcast(Message{
		Type:    "CHALLENGE_RESULT",
		Payload: result,
	})

	loser := r.State.GetPlayer(result.LoserID)
	if loser != nil {
		r.broadcast(Message{
			Type: "RUSSIAN_ROULETTE",
			Payload: map[string]interface{}{
				"player_id":    loser.ID,
				"bullet_count": loser.PunishmentCount,
				"survived":     loser.IsAlive,
			},
		})

		if !loser.IsAlive {
			r.broadcast(Message{
				Type: "PLAYER_ELIMINATED",
				Payload: map[string]interface{}{
					"player_id": loser.ID,
				},
			})
		}
	}

	if r.finalizeGameIfOver() {
		return
	}

	r.broadcastGameState()
	r.processAITurns()
}

func (r *GameRoom) handlePass(evt GameEvent) {
	if err := r.State.Pass(evt.PlayerID); err != nil {
		r.Hub.GetClient(evt.PlayerID).SendMessage(Message{
			Type:    "ERROR",
			Payload: map[string]interface{}{"msg": err.Error()},
		})
		return
	}
	r.broadcastGameState()
	r.processAITurns()
}

func (r *GameRoom) handleChat(evt GameEvent) {
	payload := struct {
		Content string `json:"content"`
	}{}
	json.Unmarshal(evt.Payload, &payload)

	if payload.Content == "" || len(payload.Content) > 500 {
		return
	}

	nickname := fmt.Sprintf("玩家%d", evt.PlayerID)
	if p := r.Players[evt.PlayerID]; p != nil && p.Nickname != "" {
		nickname = p.Nickname
	} else if client := r.Hub.GetClient(evt.PlayerID); client != nil && client.Nickname != "" {
		nickname = client.Nickname
	}

	r.broadcast(Message{
		Type: "CHAT",
		Payload: map[string]interface{}{
			"sender_id":   evt.PlayerID,
			"sender_name": nickname,
			"content":     payload.Content,
			"is_ai":       evt.AIPlayer,
		},
	})
}

func (r *GameRoom) handleAIAction(evt GameEvent) {
}

func (r *GameRoom) handleReconnect(evt GameEvent) {
	player := r.State.GetPlayer(evt.PlayerID)
	if player != nil {
		player.IsOnline = true
		player.IsAI = false
		player.AITakeover = false
		if client := r.Hub.GetClient(evt.PlayerID); client != nil {
			client.SendMessage(Message{
				Type:    "GAME_STATE",
				Payload: r.State.ToPublic(evt.PlayerID),
			})
		}
	}
}

func (r *GameRoom) handleGameOver(evt GameEvent) {
	r.State.Phase = game.PhaseGameOver
	r.broadcast(Message{
		Type: "GAME_OVER",
		Payload: map[string]interface{}{
			"winner_id": *r.State.WinnerID,
		},
	})
}

func (r *GameRoom) broadcast(msg Message) {
	r.Hub.BroadcastToRoom(r.ID, msg)
}

func (r *GameRoom) broadcastGameState() {
	for _, p := range r.Players {
		if !p.IsAI || !p.AITakeover {
			if client := r.Hub.GetClient(p.ID); client != nil {
				client.SendMessage(Message{
					Type:    "GAME_STATE",
					Payload: r.State.ToPublic(p.ID),
				})
			}
		}
	}
}

func (r *GameRoom) processAITurns() {
	if r.State.Phase != game.PhasePlaying {
		return
	}

	currentPlayer := r.State.GetCurrentPlayer()
	if currentPlayer == nil || !currentPlayer.IsAI {
		return
	}

	go func() {
		time.Sleep(1 * time.Second) // simulate thinking
		r.mu.Lock()
		defer r.mu.Unlock()
		r.executeAITurn()
	}()
}

func (r *GameRoom) executeAITurn() {
	currentPlayer := r.State.GetCurrentPlayer()
	if currentPlayer == nil || !currentPlayer.IsAI {
		return
	}

	// Simple AI strategy
	actions := r.State.GetLegalActions(currentPlayer.ID)

	if r.aiService != nil {
		action := r.aiService.GetAction(r.State, currentPlayer.ID, actions)
		r.executeAction(action, currentPlayer.ID)
		return
	}

	// Fallback: rule-based AI
	r.simpleAITurn(currentPlayer, actions)
}

func (r *GameRoom) simpleAITurn(player *game.Player, actions []string) {
	hasChallenge := false
	hasPlay := false
	for _, a := range actions {
		if a == "CHALLENGE" {
			hasChallenge = true
		}
		if a == "PLAY_CARD" {
			hasPlay = true
		}
	}

	if hasChallenge && r.State.LastPlay != nil && r.shouldChallenge(player) {
		r.executeChallenge(player.ID, r.State.LastPlay.PlayerID)
	} else if hasPlay {
		r.executeAIPlay(player)
	} else {
		r.executePass(player.ID)
	}
}

func (r *GameRoom) shouldChallenge(player *game.Player) bool {
	if r.State.LastPlay == nil {
		return false
	}

	// Count how many of the target card the player has
	targetCount := 0
	for _, c := range player.Hand {
		if c == r.State.TargetCard {
			targetCount++
		}
	}

	// If the last player claimed many cards and you have some of that kind,
	// it's more likely they're lying
	if len(r.State.LastPlay.CardIDs) >= 2 && targetCount >= 2 {
		return true
	}

	return false
}

func (r *GameRoom) executeAIPlay(player *game.Player) {
	// Choose cards to play
	playCount := 1
	if len(player.Hand) >= 3 {
		playCount = 2
	}

	indices := make([]int, playCount)
	for i := 0; i < playCount && i < len(player.Hand); i++ {
		indices[i] = i
	}

	err := r.State.PlayCard(player.ID, indices, r.State.TargetCard)
	if err == nil {
		r.broadcastGameState()
		if !r.finalizeGameIfOver() {
			r.processAITurns()
		}
	}
}

func (r *GameRoom) executeChallenge(challengerID, targetID uint) {
	result, err := r.State.Challenge(challengerID)
	if err != nil {
		return
	}

	r.broadcast(Message{
		Type:    "CHALLENGE_RESULT",
		Payload: result,
	})

	loser := r.State.GetPlayer(result.LoserID)
	if loser != nil {
		r.broadcast(Message{
			Type: "RUSSIAN_ROULETTE",
			Payload: map[string]interface{}{
				"player_id":    loser.ID,
				"bullet_count": loser.PunishmentCount,
				"survived":     loser.IsAlive,
			},
		})
		if !loser.IsAlive {
			r.broadcast(Message{
				Type: "PLAYER_ELIMINATED",
				Payload: map[string]interface{}{
					"player_id": loser.ID,
				},
			})
		}
	}

	r.broadcastGameState()
	if !r.finalizeGameIfOver() {
		r.processAITurns()
	}
}

func (r *GameRoom) executePass(playerID uint) {
	r.State.Pass(playerID)
	r.broadcastGameState()
	if !r.finalizeGameIfOver() {
		r.processAITurns()
	}
}

func (r *GameRoom) executeAction(action AIAction, playerID uint) {
	switch action.Type {
	case "PLAY_CARD":
		r.State.PlayCard(playerID, action.CardIDs, r.State.TargetCard)
	case "CHALLENGE":
		r.State.Challenge(playerID)
	case "PASS":
		r.State.Pass(playerID)
	}
	r.broadcastGameState()
	if !r.finalizeGameIfOver() {
		r.processAITurns()
	}
}

type AIAction struct {
	Type    string `json:"type"`
	CardIDs []int  `json:"card_ids,omitempty"`
	Message string `json:"message,omitempty"`
}

type AIProxy struct {
	serviceURL string
}

func NewAIProxy(url string) *AIProxy {
	return &AIProxy{serviceURL: url}
}

func (ai *AIProxy) GetAction(state *game.GameState, playerID uint, actions []string) AIAction {
	return AIAction{Type: "PASS"}
}
