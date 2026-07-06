package game

import "testing"

func newTestState() *GameState {
	players := []*Player{
		{ID: 1, Nickname: "p1", SeatIndex: 0, IsAlive: true, Hand: []Card{Ace, King}, HandCount: 2},
		{ID: 2, Nickname: "p2", SeatIndex: 1, IsAlive: true, Hand: []Card{Queen, Jack}, HandCount: 2},
	}
	return &GameState{
		Phase:         PhasePlaying,
		CurrentPlayer: 0,
		CurrentRound:  1,
		TargetCard:    Ace,
		Players:       players,
		AliveCount:    len(players),
	}
}

func TestPlayCardPreservesLastPlayAndTruthfulness(t *testing.T) {
	gs := newTestState()

	if err := gs.PlayCard(1, []int{0}, Ace); err != nil {
		t.Fatalf("PlayCard returned error: %v", err)
	}

	if gs.LastPlay == nil {
		t.Fatal("LastPlay was cleared after playing a card")
	}
	if !gs.LastPlay.Truthful {
		t.Fatal("truthful play was marked as a lie")
	}
	if len(gs.LastPlay.Cards) != 1 || gs.LastPlay.Cards[0] != Ace {
		t.Fatalf("unexpected challenged cards: %#v", gs.LastPlay.Cards)
	}
	if got := gs.GetPlayer(1).HandCount; got != 1 {
		t.Fatalf("player hand count = %d, want 1", got)
	}
	if gs.CurrentPlayer != 1 {
		t.Fatalf("current player index = %d, want 1", gs.CurrentPlayer)
	}
}

func TestPlayCardRejectsDuplicateCardIndices(t *testing.T) {
	gs := newTestState()

	if err := gs.PlayCard(1, []int{0, 0}, Ace); err == nil {
		t.Fatal("PlayCard accepted duplicate card indices")
	}
}

func TestChallengeRequiresCurrentPlayerAndClearsLastPlay(t *testing.T) {
	gs := newTestState()
	if err := gs.PlayCard(1, []int{1}, Ace); err != nil {
		t.Fatalf("PlayCard returned error: %v", err)
	}

	if _, err := gs.Challenge(1); err == nil {
		t.Fatal("non-current player was allowed to challenge")
	}

	result, err := gs.Challenge(2)
	if err != nil {
		t.Fatalf("Challenge returned error: %v", err)
	}
	if !result.Success || result.LoserID != 1 {
		t.Fatalf("unexpected challenge result: %#v", result)
	}
	if gs.LastPlay != nil {
		t.Fatal("LastPlay was not cleared after challenge resolution")
	}
}

func TestPassAdvancesTurnWithoutLastPlay(t *testing.T) {
	gs := newTestState()

	if err := gs.Pass(1); err != nil {
		t.Fatalf("Pass returned error: %v", err)
	}
	if gs.CurrentPlayer != 1 {
		t.Fatalf("current player index = %d, want 1", gs.CurrentPlayer)
	}
	if gs.CurrentTurn != 1 {
		t.Fatalf("current turn = %d, want 1", gs.CurrentTurn)
	}
}
