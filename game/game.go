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

func (g *Game) ScaleBoard(w, h, l int) error {
	newGame, err := NewGame(w, h, l)
	if err != nil {
		return err
	}

	offsetX := (w - g.BoardWidth) / 2
	offsetY := (h - g.BoardHeight) / 2

	for x := 0; x < g.BoardWidth; x++ {
		for y := 0; y < g.BoardHeight; y++ {
			newX := x + offsetX
			newY := y + offsetY
			newGame.Board[newX+newY*w] = g.Board[x+y*g.BoardWidth]
		}
	}

	newGame.PlayerTurn = g.PlayerTurn
	newGame.StepsCount = g.StepsCount
	newGame.CheckWin()

	*g = *newGame

	return nil
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

	countX, countO := 0, 0
	for _, p := range g.Board {
		switch p {
		case PlayerX:
			countX++
		case PlayerO:
			countO++
		}
	}

	if countX < countO {
		g.PlayerTurn = PlayerX
	} else {
		g.PlayerTurn = PlayerO
	}

	g.StepsCount = countX + countO

	for s := range MapSizes {
		if s*s == len(g.Board) {
			g.BoardWidth = s
			g.BoardHeight = s
			break
		}
	}

	if g.BoardWidth == 3 {
		g.WinLength = 3
	} else {
		g.WinLength = g.BoardWidth - 1
	}

	return g, nil
}
