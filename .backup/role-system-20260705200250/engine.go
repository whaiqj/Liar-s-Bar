package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Card string

const (
	Ace   Card = "A"
	King  Card = "K"
	Queen Card = "Q"
	Jack  Card = "J"
)

var CardOrder = []Card{Ace, King, Queen, Jack}

type GamePhase string

const (
	PhaseWaiting    GamePhase = "WAITING"
	PhaseMatched    GamePhase = "MATCHED"
	PhasePlaying    GamePhase = "PLAYING"
	PhaseChallenge  GamePhase = "CHALLENGE"
	PhasePunishment GamePhase = "PUNISHMENT"
	PhaseRoundEnd   GamePhase = "ROUND_END"
	PhaseGameOver   GamePhase = "GAME_OVER"
)

type ActionType string

const (
	ActPlayCard   ActionType = "PLAY_CARD"
	ActChallenge  ActionType = "CHALLENGE"
	ActPass       ActionType = "PASS"
	ActChat       ActionType = "CHAT"
	ActPunishment ActionType = "PUNISHMENT"
	ActEliminated ActionType = "ELIMINATED"
	ActGameOver   ActionType = "GAME_OVER"
)

type RoundRecord struct {
	PlayerID   uint
	CardIDs    []int
	Claim      Card
	Truthful   bool
	Cards      []Card
	Challenged bool
}

type Player struct {
	ID               uint   `json:"id"`
	Nickname         string `json:"nickname"`
	IsAI             bool   `json:"is_ai"`
	IsOnline         bool   `json:"is_online"`
	IsAlive          bool   `json:"is_alive"`
	IsReady          bool   `json:"is_ready"`
	AITakeover       bool   `json:"ai_takeover"`
	SeatIndex        int    `json:"seat_index"`
	Hand             []Card `json:"-"`
	HandCount        int    `json:"hand_count"`
	PunishmentCount  int    `json:"punishment_count"`
	PlayCount        int    `json:"play_count"`
	LieCount         int    `json:"lie_count"`
	ChallengeCount   int    `json:"challenge_count"`
	ChallengeSuccess int    `json:"challenge_success"`
}

type GameState struct {
	Phase         GamePhase     `json:"phase"`
	CurrentPlayer int           `json:"current_player"`
	CurrentRound  int           `json:"current_round"`
	CurrentTurn   int           `json:"current_turn"`
	TargetCard    Card          `json:"target_card"`
	Players       []*Player     `json:"players"`
	LastPlay      *RoundRecord  `json:"last_play"`
	RoundHistory  []RoundRecord `json:"-"`
	Deck          []Card        `json:"-"`
	DiscardPile   []Card        `json:"-"`
	RoundCounter  int           `json:"-"`
	WinnerID      *uint         `json:"winner_id"`
	AliveCount    int           `json:"alive_count"`
}

func NewDeck() []Card {
	deck := make([]Card, 0, 24)
	cards := []Card{Ace, King, Queen, Jack}
	for _, c := range cards {
		for i := 0; i < 6; i++ {
			deck = append(deck, c)
		}
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func DealCards(deck []Card, numPlayers int) ([][]Card, []Card) {
	hands := make([][]Card, numPlayers)
	for i := 0; i < numPlayers; i++ {
		hands[i] = make([]Card, 6)
		for j := 0; j < 6; j++ {
			hands[i][j] = deck[i*6+j]
		}
	}
	remaining := deck[numPlayers*6:]
	return hands, remaining
}

func (gs *GameState) InitGame(players []*Player) {
	gs.Players = players
	gs.Phase = PhasePlaying
	gs.CurrentRound = 1
	gs.CurrentTurn = 0
	gs.TargetCard = Ace
	gs.AliveCount = len(players)
	gs.RoundHistory = make([]RoundRecord, 0)

	deck := NewDeck()
	hands, remaining := DealCards(deck, len(players))
	for i, p := range players {
		p.Hand = hands[i]
		p.HandCount = len(hands[i])
		p.IsAlive = true
		p.IsReady = false
		p.PunishmentCount = 0
	}
	gs.Deck = remaining
	gs.CurrentPlayer = rand.Intn(len(players))
}

func (gs *GameState) GetPlayer(playerID uint) *Player {
	for _, p := range gs.Players {
		if p.ID == playerID {
			return p
		}
	}
	return nil
}

func (gs *GameState) GetCurrentPlayer() *Player {
	if gs.CurrentPlayer < len(gs.Players) {
		return gs.Players[gs.CurrentPlayer]
	}
	return nil
}

func (gs *GameState) NextPlayer() {
	start := (gs.CurrentPlayer + 1) % len(gs.Players)
	for i := 0; i < len(gs.Players); i++ {
		idx := (start + i) % len(gs.Players)
		if gs.Players[idx].IsAlive {
			gs.CurrentPlayer = idx
			gs.CurrentTurn++
			return
		}
	}
}

func (gs *GameState) GetPreviousPlayer() *Player {
	if gs.LastPlay == nil {
		return nil
	}
	for _, p := range gs.Players {
		if p.ID == gs.LastPlay.PlayerID {
			return p
		}
	}
	return nil
}

func (gs *GameState) PlayCard(playerID uint, cardIndices []int, claim Card) error {
	player := gs.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found")
	}
	if !player.IsAlive {
		return fmt.Errorf("player is eliminated")
	}
	currentPlayer := gs.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return fmt.Errorf("not your turn")
	}
	if len(cardIndices) < 1 || len(cardIndices) > 3 {
		return fmt.Errorf("must play 1-3 cards")
	}
	if gs.Phase != PhasePlaying {
		return fmt.Errorf("invalid phase")
	}
	if claim != gs.TargetCard {
		return fmt.Errorf("must claim current target card")
	}

	seen := make(map[int]bool, len(cardIndices))
	for _, idx := range cardIndices {
		if idx < 0 || idx >= len(player.Hand) {
			return fmt.Errorf("invalid card index")
		}
		if seen[idx] {
			return fmt.Errorf("duplicate card index")
		}
		seen[idx] = true
	}

	selected := make([]Card, 0, len(cardIndices))
	newHand := make([]Card, 0)
	skipSet := make(map[int]bool)
	for _, idx := range cardIndices {
		skipSet[idx] = true
	}
	for i, c := range player.Hand {
		if skipSet[i] {
			selected = append(selected, c)
		} else {
			newHand = append(newHand, c)
		}
	}

	truthful := true
	for _, c := range selected {
		if c != gs.TargetCard {
			truthful = false
			break
		}
	}

	player.Hand = newHand
	player.HandCount = len(newHand)
	player.PlayCount++
	if !truthful {
		player.LieCount++
	}
	gs.DiscardPile = append(gs.DiscardPile, selected...)

	gs.LastPlay = &RoundRecord{
		PlayerID: playerID,
		CardIDs:  cardIndices,
		Claim:    claim,
		Truthful: truthful,
		Cards:    selected,
	}

	gs.NextPlayer()
	gs.checkEmptyHands()
	gs.RecordAction(playerID, ActPlayCard, map[string]interface{}{
		"count":    len(selected),
		"claim":    claim,
		"truthful": truthful,
	})
	return nil
}

func (gs *GameState) Challenge(challengerID uint) (*ChallengeResult, error) {
	challenger := gs.GetPlayer(challengerID)
	if challenger == nil || !challenger.IsAlive {
		return nil, fmt.Errorf("invalid challenger")
	}
	currentPlayer := gs.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != challengerID {
		return nil, fmt.Errorf("not your turn")
	}
	if gs.LastPlay == nil {
		return nil, fmt.Errorf("no cards to challenge")
	}
	prevPlayer := gs.GetPreviousPlayer()
	if prevPlayer == nil || prevPlayer.ID == challengerID {
		return nil, fmt.Errorf("cannot challenge yourself")
	}

	challenger.ChallengeCount++

	result := &ChallengeResult{
		Success: !gs.LastPlay.Truthful,
		LiarID:  0,
		LoserID: 0,
	}

	var loser *Player
	if !gs.LastPlay.Truthful {
		result.LiarID = prevPlayer.ID
		result.LoserID = prevPlayer.ID
		loser = prevPlayer
		challenger.ChallengeSuccess++
	} else {
		result.LoserID = challengerID
		loser = challenger
	}

	result.Truthful = gs.LastPlay.Truthful
	result.ChallengedCards = gs.LastPlay.Cards
	result.ChallengerID = challengerID
	gs.LastPlay.Challenged = true

	gs.Phase = PhasePunishment
	gs.RecordAction(challengerID, ActChallenge, map[string]interface{}{
		"target":  prevPlayer.ID,
		"success": result.Success,
	})

	_ = loser
	gs.punishPlayer(loser)
	return result, nil
}

type ChallengeResult struct {
	Success         bool   `json:"success"`
	Truthful        bool   `json:"truthful"`
	LiarID          uint   `json:"liar_id"`
	LoserID         uint   `json:"loser_id"`
	ChallengerID    uint   `json:"challenger_id"`
	ChallengedCards []Card `json:"challenged_cards"`
}

func (gs *GameState) Pass(playerID uint) error {
	currentPlayer := gs.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.ID != playerID {
		return fmt.Errorf("not your turn")
	}
	if gs.Phase != PhasePlaying {
		return fmt.Errorf("invalid phase")
	}
	gs.NextPlayer()
	return nil
}

func (gs *GameState) punishPlayer(player *Player) {
	player.PunishmentCount++
	bulletSlots := player.PunishmentCount
	if bulletSlots > 6 {
		bulletSlots = 6
	}

	hitChamber := rand.Intn(6) + 1
	survived := hitChamber > bulletSlots

	gs.RecordAction(player.ID, ActPunishment, map[string]interface{}{
		"bullet_count": bulletSlots,
		"chamber":      hitChamber,
		"survived":     survived,
	})

	gs.LastPlay = nil
	if !survived {
		player.IsAlive = false
		gs.AliveCount--
		gs.RecordAction(player.ID, ActEliminated, nil)

		if gs.AliveCount <= 1 {
			gs.endGame()
			return
		}
	}

	gs.Phase = PhasePlaying
	gs.NextPlayer()
	gs.checkEmptyHands()
}

func (gs *GameState) checkEmptyHands() {
	allEmpty := true
	for _, p := range gs.Players {
		if p.IsAlive && len(p.Hand) > 0 {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		gs.startNewRound()
	}
}

func (gs *GameState) startNewRound() {
	gs.RoundCounter++
	gs.CurrentRound++
	gs.TargetCard = CardOrder[gs.RoundCounter%len(CardOrder)]
	gs.LastPlay = nil

	deck := NewDeck()
	alivePlayers := make([]*Player, 0)
	for _, p := range gs.Players {
		if p.IsAlive {
			alivePlayers = append(alivePlayers, p)
		}
	}

	hands, _ := DealCards(deck, len(alivePlayers))
	for i, p := range alivePlayers {
		p.Hand = hands[i]
		p.HandCount = len(hands[i])
	}
	gs.Deck = nil
	gs.DiscardPile = nil
}

func (gs *GameState) endGame() {
	gs.Phase = PhaseGameOver
	for _, p := range gs.Players {
		if p.IsAlive {
			winnerID := p.ID
			gs.WinnerID = &winnerID
			break
		}
	}
}

func (gs *GameState) RecordAction(playerID uint, actionType ActionType, data interface{}) {
}

func (gs *GameState) ToPublic(playerID uint) map[string]interface{} {
	players := make([]map[string]interface{}, len(gs.Players))
	for i, p := range gs.Players {
		showHand := false
		if p.ID == playerID {
			showHand = true
		}
		pm := map[string]interface{}{
			"id":               p.ID,
			"nickname":         p.Nickname,
			"is_ai":            p.IsAI,
			"is_online":        p.IsOnline,
			"is_alive":         p.IsAlive,
			"is_ready":         p.IsReady,
			"seat_index":       p.SeatIndex,
			"hand_count":       p.HandCount,
			"punishment_count": p.PunishmentCount,
		}
		if showHand {
			pm["hand"] = p.Hand
		}
		players[i] = pm
	}

	result := map[string]interface{}{
		"phase":          gs.Phase,
		"current_player": gs.CurrentPlayer,
		"current_round":  gs.CurrentRound,
		"current_turn":   gs.CurrentTurn,
		"target_card":    gs.TargetCard,
		"players":        players,
		"alive_count":    gs.AliveCount,
	}

	if gs.LastPlay != nil {
		result["last_play"] = map[string]interface{}{
			"player_id": gs.LastPlay.PlayerID,
			"count":     len(gs.LastPlay.CardIDs),
			"claim":     gs.LastPlay.Claim,
		}
	}

	if gs.WinnerID != nil {
		result["winner_id"] = *gs.WinnerID
	}

	return result
}

func (gs *GameState) GetLegalActions(playerID uint) []string {
	player := gs.GetPlayer(playerID)
	if player == nil || !player.IsAlive || gs.Phase != PhasePlaying {
		return nil
	}

	actions := make([]string, 0)
	if gs.GetCurrentPlayer().ID == playerID {
		actions = append(actions, "PLAY_CARD")
	}

	if gs.LastPlay != nil && gs.LastPlay.PlayerID != playerID {
		actions = append(actions, "CHALLENGE")
	}

	actions = append(actions, "PASS", "CHAT")
	return actions
}

func (gs *GameState) MarshalJSON() ([]byte, error) {
	type Alias GameState
	return json.Marshal(&struct{ *Alias }{Alias: (*Alias)(gs)})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
