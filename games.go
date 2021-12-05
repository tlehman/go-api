/* This is where the rules of Go are defined and enforced
 */
package main

import (
	"errors"
	"time"
)

// A Go Game has two players, and proceeds as a sequence of moves
type Game struct {
	id         uint64
	player1_id string // User id of player 1
	player2_id string // User id of player 2
	moves      []Move
	size       uint8 // if size == 9, then the board is 9x9
}

func NewGame(id uint64, player1 User, player2 User, size uint8) (Game, error) {
	if size == uint8(0) {
		return Game{}, errors.New("zero board size invalid")
	}
	g := Game{id: id, player1_id: player1.id, player2_id: player2.id, size: size}
	return g, nil
}

// During a Game, a player makes a Move at a point in time and space
type Move struct {
	x         uint8
	y         uint8
	t         time.Time
	player_id string // User id of player who made this move
}

// A Move is only valid if it follows the rules of Go
func play(g Game, nextMove Move) (Game, error) {
	// (x,y) must be in the [1..size]x[1..size] integer plane

	var lastMove Move
	for _, m := range g.moves {
		lastMove = m
	}
	if lastMove.player_id == nextMove.player_id {
		return g, errors.New("last player must be different than current player")
	}
	if g.player1_id != nextMove.player_id && g.player2_id != nextMove.player_id {
		return g, errors.New("current player must be in the game")
	}

	g.moves = append(g.moves, nextMove)

	return g, nil
}

// The state of the game
func state(g Game) [3][3]uint8 {
	return [3][3]uint8 {{0,0,0},
											{0,0,0},
											{0,0,0}}
}