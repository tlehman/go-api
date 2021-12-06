/* This is where the rules of Go are defined and enforced
 */
package main

import (
  "time"
	"errors"
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

// Game stores both player1 and player2's ids. Empty cells 
// are represented as 0, and players 1 and 2's moves are 
// represented as 1 and 2, respectively
func (g *Game) playerNumber(m Move) uint8 {
  if g.player1_id == m.player_id {
    return uint8(1)
  } else if g.player2_id == m.player_id {
    return uint8(2)
  } else {
    // the play function will prevent it from being set on the board
    return uint8(3)
  }
}

func (g *Game) NewMove(x uint8, y uint8) Move {
  var player_id string
  if(len(g.moves) == 0) {
    player_id = g.player1_id
  } else {
    lastMove := g.moves[len(g.moves)-1]
    if lastMove.player_id == g.player1_id {
      player_id = g.player2_id
    } else {
      player_id = g.player1_id
    }
  }

  return Move{x: x, y: y, t: time.Now(), player_id: player_id}
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
func state(g Game) [][]uint8 {
	s := create2dSlice(g.size, g.size)
  // Walk through moves, compute new state
  for _, move := range g.moves {
    s[move.y-1][move.x-1] = g.playerNumber(move)
  }
	return s
}

func create2dSlice(w, h uint8) [][]uint8 {
    a := make([]uint8, w*h)
    s := make([][]uint8, h)
    lo, hi := uint8(0), w
    for i := range s {
        s[i] = a[lo:hi:hi]
        lo, hi = hi, hi+w
    }
    return s
}