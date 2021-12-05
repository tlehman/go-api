package main

import (
	"testing"
	"time"
)

<<<<<<< HEAD
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
	s := state(g)
	if s != [3][3]uint8{{0,0,0},
	                    {0,0,0},
										  {0,0,0}} {
    t.Fatalf("%v is supposed to be all zeros", s)
	}
=======
func TestNewGame(t *testing.T) {
	
>>>>>>> origin/main
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
	if(err == nil) {
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
	var err error

	// Alice and Bob are playing a game
	a := User{name: "Alice", id: "1"}
	b := User{name: "Bob", id: "2"}

	_, err = NewGame(3, a, b, 13)
	if err != nil {
		t.Fatal(err)
	}


	/*
	C4
	K10
	D4
	C3
	D10
	D3
	D9
	E4
	E10
	D5
	F10
	C5
	F11
	B4 (then C4,C5 are captured)	
	 */
	
}