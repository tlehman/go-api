/* This is where the rules of Go are defined and enforced
 */
package main

import (
	"errors"
	"fmt"
	"time"
)

// A Go Game has two players, and proceeds as a sequence of moves
type Game struct {
	id                    uint64
	size                  uint8  // if size == 9, then the board is 9x9
	player1_id            string // User id of player 1
	player2_id            string // User id of player 2
	moves                 []Move
	components            []Component // Connected components
	state                 [][]uint8
	pointToComponentIndex map[Point]int // used to quickly look up component index of point
}

type Point struct {
	x uint8
	y uint8
}

func (p *Point) north() Point {
	return Point{x: p.x, y: p.y - 1}
}
func (p *Point) south() Point {
	return Point{x: p.x, y: p.y + 1}
}
func (p *Point) east() Point {
	return Point{x: p.x + 1, y: p.y}
}
func (p *Point) west() Point {
	return Point{x: p.x - 1, y: p.y}
}

type Component struct {
	positions    []Point
	playerNumber uint8 // 0=empty, 1=player1, 2=player2
	liberties    uint
}

func NewGame(id uint64, player1 User, player2 User, size uint8) (Game, error) {
	if size == uint8(0) {
		return Game{}, errors.New("zero board size invalid")
	}
	s := create2dSlice(size, size)
	g := Game{
		id:                    id,
		player1_id:            player1.id,
		player2_id:            player2.id,
		state:                 s,
		size:                  size,
		components:            make([]Component, 0),
		pointToComponentIndex: make(map[Point]int),
	}
	return g, nil
}

// During a Game, a player makes a Move at a point in time and space
type Move struct {
	x         uint8
	y         uint8
	t         time.Time
	player_id string // User id of player who made this move
}

func (m *Move) toPoint() Point {
	return Point{x: m.x, y: m.y}
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
	if len(g.moves) == 0 {
		player_id = g.player1_id
	} else {
		lastMove := g.moves[len(g.moves)-1]
		if lastMove.player_id == g.player1_id {
			player_id = g.player2_id
		} else {
			player_id = g.player1_id
		}
	}
	// translate 1..size to 0..size-1
	return Move{x: x-1, y: y-1, t: time.Now(), player_id: player_id}
}

// A Move is only valid if it follows the rules of Go
func (g *Game) play(nextMove Move) error {
	// (x,y) must be in the [1..size]x[1..size] integer plane
	var lastMove Move
	for _, m := range g.moves {
		lastMove = m
	}
	if lastMove.player_id == nextMove.player_id {
		return errors.New("last player must be different than current player")
	}
	if g.player1_id != nextMove.player_id && g.player2_id != nextMove.player_id {
		return errors.New("current player must be in the game")
	}
	var playerNumber = uint8(2)
	if g.player1_id == nextMove.player_id {
		playerNumber = uint8(1)
	}

	// Find which component nextMove would belong to.
	// There are three possibilities, nextMove either:
	// 1. adds a new component
	// 2. expands an existing component
	// 3. merges 2-4 existing components
	// After the component update, then we can implement the capture logic
	p := nextMove.toPoint()
	// Find all player's neighbors (north, east, south and west)
	neighborCount := g.countNeighbors(p, playerNumber)

	if neighborCount == 0 {
		// 1. Add new component
		g.state[p.y][p.x] = playerNumber
		g.components = append(g.components,
			Component{
				positions:    []Point{p},
				liberties:    4, // number of open spaces around it
				playerNumber: playerNumber,
			})
		// Store component index for point p
		g.pointToComponentIndex[p] = len(g.components)-1
		// Commit move 
		g.moves = append(g.moves, nextMove)
	} else if neighborCount == 1 {
		// 2. Expand existing component
		// TODO: find neighbor and join that component
		var neighbor Point
		if g.state[p.north().y][p.north().x] == playerNumber {
			neighbor = p.north()
			fmt.Printf(" north neighbor = %v\n", neighbor)
		} else if g.state[p.south().y][p.south().x] == playerNumber {
			neighbor = p.south()
		} else if g.state[p.east().y][p.east().x] == playerNumber {
			neighbor = p.east()
		} else {
			neighbor = p.west()
		}
		// add p to neighbor's component
		var compIdx int = g.pointToComponentIndex[neighbor]
		g.pointToComponentIndex[p] = compIdx
		g.components[compIdx].positions = append(g.components[compIdx].positions, p)
		// Update state
		g.state[p.y][p.x] = playerNumber
		// Commit move 
		g.moves = append(g.moves, nextMove)
	}

	return nil
}

func slicefmt(sl [][]uint8) string {
	var s string
	for i := 0; i < len(sl); i++ {
		s += fmt.Sprintf("%v\n", sl[i])
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

// Count all neighbors belonging to the same player
func (g *Game) countNeighbors(p Point, playerNumber uint8) int {
	var count int = 0

	// Look west
	if p.x > 0 {
		wp := p.west()
		if g.state[wp.y][wp.x] == playerNumber {
			count++
		}
	}
	// Look east
	if p.x < g.size-1 {
		ep := p.east()
		if g.state[ep.y][ep.x] == playerNumber {
			count++
		}
	}
	// Look north
	if p.y > 0 {
		np := p.north()
		if g.state[np.y][np.x] == playerNumber {
			count++
		}
	}
	// Look south
	if p.y < g.size-1 {
		sp := p.south()
		if g.state[sp.y][sp.x] == playerNumber {
			count++
		}
	}
	return count
}