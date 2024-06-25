package game

import (
	"fmt"
	"tictactoe/util"
)

type Game struct {
	PlayerTurn  Player
	PlayerWon   Player
	StepsCount  int
	Board       []Player
	BoardWidth  int
	BoardHeight int
	WinLength   int
}

func NewGame(w, h, l int) (*Game, error) {
	g := &Game{
		PlayerTurn:  PlayerX,
		PlayerWon:   PlayerNone,
		StepsCount:  0,
		Board:       make([]Player, w*h),
		BoardWidth:  w,
		BoardHeight: h,
		WinLength:   l,
	}

	for i := 0; i < w*h; i++ {
		g.Board[i] = PlayerNone
	}

	return g, nil
}

func (g *Game) Copy() *Game {
	newGame := &Game{}

	newGame.PlayerTurn = g.PlayerTurn
	newGame.PlayerWon = g.PlayerWon
	newGame.StepsCount = g.StepsCount
	newGame.BoardWidth = g.BoardWidth
	newGame.BoardHeight = g.BoardHeight
	newGame.WinLength = g.WinLength

	newGame.Board = make([]Player, len(g.Board))

	for i, player := range g.Board {
		newGame.Board[i] = player
	}

	return newGame
}

func (g *Game) GetMapKey() string {
	return util.GetMapKey(g.BoardWidth, g.BoardHeight, g.WinLength)
}

func (g *Game) MakeMoveByIndex(i int) {
	g.Board[i] = g.PlayerTurn
	g.PlayerTurn = g.PlayerTurn.Opponent()
	g.StepsCount++
	g.CheckWin()
}

func (g *Game) MakeMoveByCoordinates(x, y int) {
	g.MakeMoveByIndex(x + y*g.BoardWidth)
}

func (g *Game) CheckWin() {
	var w, h, l = g.BoardWidth, g.BoardHeight, g.WinLength

	for _, positions := range GetWinPositions(w, h, l) {
		var player = g.Board[positions[0]]
		var count int

		if player == PlayerNone {
			continue
		}

		for _, i := range positions {
			if player != g.Board[i] {
				break
			}

			count++
		}

		if count == g.WinLength {
			g.PlayerWon = player
			break
		}
	}
}

func (g *Game) IsFulfilled() bool {
	return g.StepsCount == g.BoardWidth*g.BoardHeight
}

func (g *Game) IsOver() bool {
	return g.PlayerWon != PlayerNone || g.IsFulfilled()
}

func (g *Game) String() string {
	return fmt.Sprintf("%c %s", g.PlayerWon, g.Board)
}

func FromString(str string) (*Game, error) {
	g := &Game{}

	_, err := fmt.Sscanf(str, "%c %s", &g.PlayerWon, &g.Board)
	if err != nil {
		return nil, err
	}

	return g, nil
}
