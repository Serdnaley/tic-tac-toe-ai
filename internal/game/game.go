package game

import (
	"fmt"
	"tictactoe/internal/util"
)

type Game struct {
	PlayerTurn Player
	PlayerWon  Player
	StepsCount int
	Board      []Player
	Size       int
	WinLength  int
}

func NewGame(s, l int) (*Game, error) {
	g := &Game{
		PlayerTurn: PlayerX,
		PlayerWon:  PlayerNone,
		StepsCount: 0,
		Board:      make([]Player, s*s),
		Size:       s,
		WinLength:  l,
	}

	for i := 0; i < s*s; i++ {
		g.Board[i] = PlayerNone
	}

	return g, nil
}

func (g *Game) Copy() *Game {
	newGame := &Game{}

	newGame.PlayerTurn = g.PlayerTurn
	newGame.PlayerWon = g.PlayerWon
	newGame.StepsCount = g.StepsCount
	newGame.Size = g.Size
	newGame.WinLength = g.WinLength

	newGame.Board = make([]Player, len(g.Board))

	for i, player := range g.Board {
		newGame.Board[i] = player
	}

	return newGame
}

func (g *Game) GetMapKey() string {
	return util.GetMapKey(g.Size, g.WinLength)
}

func (g *Game) MakeMoveByIndex(i int) {
	g.Board[i] = g.PlayerTurn
	g.PlayerTurn = g.PlayerTurn.Opponent()
	g.StepsCount++
	g.CheckWin()
}

func (g *Game) MakeMoveByCoordinates(x, y int) {
	g.MakeMoveByIndex(x + y*g.Size)
}

func (g *Game) CheckWin() {
	var s, l = g.Size, g.WinLength

	for _, positions := range GetWinPositions(s, l) {
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

func (g *Game) ScaleBoard(s, l int) error {
	newGame, err := NewGame(s, l)
	if err != nil {
		return err
	}

	offset := (s - g.Size) / 2

	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			newX := x + offset
			newY := y + offset
			newGame.Board[newX+newY*s] = g.Board[x+y*g.Size]
		}
	}

	newGame.PlayerTurn = g.PlayerTurn
	newGame.StepsCount = g.StepsCount
	newGame.CheckWin()

	*g = *newGame

	return nil
}

func (g *Game) IsFulfilled() bool {
	return g.StepsCount == g.Size*g.Size
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
			g.Size = s
			break
		}
	}

	if g.Size <= 4 {
		g.WinLength = g.Size
	} else {
		g.WinLength = g.Size - 1
	}

	return g, nil
}
