package main

import (
	"testing"
	"time"
)

func equal(s [][]uint8, t [][]uint8) bool {
	if len(s) != len(t) || cap(s) != cap(t) {
		return false
	}
	for i := 0; i < len(s); i += 1 {
		for j := 0; j < len(s[0]); j += 1 {
			if s[i][j] != t[i][j] {
				return false
			}
		}
	}
	return true
}

func TestNewGameState(t *testing.T) {
	var g Game
	var err error

	// Alice and Bob are playing a game
	a := User{name: "Alice", id: "1"}
	b := User{name: "Bob", id: "2"}

	// Before they start, the board state should equal this
	g, err = NewGame(1, a, b, 3)
	if err != nil {
		t.Fatal(err)
	}
	actual := state(g)
	expected := [][]uint8{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	if !equal(expected, actual) {
		t.Fatalf("%v is supposed to be all zeros", actual)
	}
}

func TestGameStateNoCapture(t *testing.T) {
	var g Game
	var err error

	// Alice and Bob are playing a game
	a := User{name: "Alice", id: "1"}
	b := User{name: "Bob", id: "2"}

	// Before they start, the board state should equal this
	g, err = NewGame(1, a, b, 3)
	if err != nil {
		t.Fatal(err)
	}

	g, err = play(g, g.NewMove(1, 1)) // Player1 plays (1,1)
	g, err = play(g, g.NewMove(2, 1)) // Player2 plays (2,1)

	actual := state(g)
	expected := [][]uint8{
		{1, 2, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	if !equal(expected, actual) {
		t.Fatalf("%v is supposed to be %v", actual, expected)
	}
}

// Verify that only players attached to a game can play
func TestPlayUserRules(t *testing.T) {
	var g1, g2 Game
	var err error

	// Alice and Bob are playing a game
	a := User{name: "Alice", id: "1"}
	b := User{name: "Bob", id: "2"}
	k := User{name: "Karen", id: "3"}

	// Games cannot be 0x0, they must be at least 9x9
	_, err = NewGame(1, a, b, 0)
	if err == nil {
		t.Fatalf("size must be at least 1")
	}

	g1, err = NewGame(1, a, b, 9)

	// Alice should be able to play (1,1) on an empty board
	m1 := Move{x: 1, y: 1, t: time.Now(), player_id: a.id}
	g1, err = play(g1, m1)
	if len(g1.moves) != 1 {
		t.Fatalf("Move #%v should be playable in game %d", m1, g1.id)
	}

	// Bob should be able to play
	m2 := Move{x: 5, y: 3, t: time.Now(), player_id: b.id}
	g1, err = play(g1, m2)
	if len(g1.moves) != 2 {
		t.Fatalf("%v should be playable in %v", m2, g1)
	}

	// But Karen is not part of this game,
	// so she should not be able to play in this game
	m3 := Move{x: 8, y: 2, t: time.Now(), player_id: k.id}
	g1, err = play(g1, m3)
	if len(g1.moves) != 2 {
		t.Fatalf("%v should not be playable in %v", m3, g1)
	}

	// Karen starts a game with Bob, and she can play
	// in that game
	g2, _ = NewGame(2, k, b, 9)
	g2, err = play(g2, m3)
	if len(g2.moves) != 1 {
		t.Fatalf("%v should be playable in %v", m3, g2)
	}
}

// Verify that the capture rules work
// Reference: https://www.britgo.org/files/rules/GoQuickRef.pdf
func TestPlayCaptureRules(t *testing.T) {
	var g Game
	var err error

	// Alice and Bob are playing a game
	a := User{name: "Alice", id: "1"}
	b := User{name: "Bob", id: "2"}

	g, err = NewGame(3, a, b, 13)
	if err != nil {
		t.Fatal(err)
	}

	g, err = play(g, g.NewMove(3, 4))
	g, err = play(g, g.NewMove(11,10))
	g, err = play(g, g.NewMove( 4, 4))
	g, err = play(g, g.NewMove( 3, 3))
	g, err = play(g, g.NewMove(4,10))
	g, err = play(g, g.NewMove(4, 3))
	g, err = play(g, g.NewMove(4, 9))
	g, err = play(g, g.NewMove(5, 4))
	g, err = play(g, g.NewMove(5,10))
	g, err = play(g, g.NewMove(4, 5))
	g, err = play(g, g.NewMove(6,10))
	g, err = play(g, g.NewMove(3, 5))
	g, err = play(g, g.NewMove(6,11))
	//g, err = play(g, g.NewMove(2, 4)) // should capture C4,C5
	
	actual := state(g)
	expected := [][]uint8{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 2, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	if !equal(expected, actual) {
		t.Fatalf("actual = \n%s\n\nexpected = \n%s", slicefmt(actual), slicefmt(expected))
	}
	/*
		C4.   3, 4
		K10. 11,10
		D4.   4, 4
		C3.   3, 3
		D10.  4,10
		D3.   4, 3
		D9.   4, 9
		E4.   5, 4
		E10.  5,10
		D5.   4, 5
		F10.  6,10
		C5.   3, 5
		F11.  6,11
		B4 (then C4,C5 are captured)
	*/

}
