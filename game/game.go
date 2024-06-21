package game

import (
	"errors"
	"fmt"
	"strings"
	"tictactoe/util"
)

var WinPositionsCache = map[string][][][2]int{}

type Step struct {
	Player
	X int
	Y int
}

func (s Step) String() string {
	return fmt.Sprintf("%c(%d,%d)", s.Player, s.X, s.Y)
}

type Game struct {
	PlayerTurn   Player
	PlayerWon    Player
	BoardHistory []Step
	Board        [][]Player
	BoardWidth   int
	BoardHeight  int
	WinLength    int
	WinPositions [][][2]int
}

func NewGame(w, h, l int) (*Game, error) {
	g := &Game{
		PlayerTurn:   PlayerX,
		PlayerWon:    PlayerNone,
		WinLength:    l,
		BoardWidth:   w,
		BoardHeight:  h,
		BoardHistory: make([]Step, 0),
		Board:        make([][]Player, w),
	}

	for i := range g.Board {
		g.Board[i] = make([]Player, h)

		for j := range g.Board[i] {
			g.Board[i][j] = PlayerNone
		}
	}

	err := g.refreshWinPositions()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) Copy() *Game {
	newGame := &Game{}

	newGame.PlayerTurn = g.PlayerTurn
	newGame.PlayerWon = g.PlayerWon
	newGame.WinLength = g.WinLength
	newGame.BoardWidth = g.BoardWidth
	newGame.BoardHeight = g.BoardHeight

	newGame.BoardHistory = make([]Step, len(g.BoardHistory))

	for i, step := range g.BoardHistory {
		newGame.BoardHistory[i] = step
	}

	newGame.Board = make([][]Player, len(g.Board))

	for x := range g.Board {
		newGame.Board[x] = make([]Player, len(g.Board[x]))

		for y := range g.Board[x] {
			newGame.Board[x][y] = g.Board[x][y]
		}
	}

	newGame.WinPositions = make([][][2]int, len(g.WinPositions))

	for i, positions := range g.WinPositions {
		newGame.WinPositions[i] = make([][2]int, len(positions))

		for j, position := range positions {
			newGame.WinPositions[i][j] = position
		}
	}

	return newGame
}

func (g *Game) refreshWinPositions() error {
	if g.BoardWidth < g.WinLength || g.BoardHeight < g.WinLength {
		return errors.New(fmt.Sprintf("Board %dx%d is too small for winning length %d", g.BoardWidth, g.BoardHeight, g.WinLength))
	}

	cacheKey := util.GetMapKey(g.BoardWidth, g.BoardHeight, g.WinLength)

	if WinPositionsCache[cacheKey] != nil {
		g.WinPositions = WinPositionsCache[cacheKey]
		return nil
	}

	// Horizontal
	for xOffset := 0; xOffset <= g.BoardWidth-g.WinLength; xOffset++ {
		for x := 0; x < g.WinLength; x++ {
			var column [][2]int

			for y := 0; y < g.WinLength; y++ {
				column = append(column, [2]int{x + xOffset, y})
			}

			g.WinPositions = append(g.WinPositions, column)
		}
	}

	// Vertical
	for yOffset := 0; yOffset <= g.BoardHeight-g.WinLength; yOffset++ {
		for y := 0; y < g.WinLength; y++ {
			var row [][2]int

			for x := 0; x < g.WinLength; x++ {
				row = append(row, [2]int{y + yOffset, x})
			}

			g.WinPositions = append(g.WinPositions, row)
		}
	}

	// Diagonal
	for xOffset := 0; xOffset <= g.BoardWidth-g.WinLength; xOffset++ {
		for yOffset := 0; yOffset <= g.BoardHeight-g.WinLength; yOffset++ {
			var diagonal1 [][2]int
			var diagonal2 [][2]int

			for i := 0; i < g.WinLength; i++ {
				diagonal1 = append(diagonal1, [2]int{xOffset + i, yOffset + i})
				diagonal2 = append(diagonal2, [2]int{xOffset + g.WinLength - 1 - i, yOffset + i})
			}

			g.WinPositions = append(g.WinPositions, diagonal1)
			g.WinPositions = append(g.WinPositions, diagonal2)
		}
	}

	WinPositionsCache[cacheKey] = g.WinPositions

	return nil
}

func (g *Game) PrintBoard() {
	for x := range g.Board {
		for y := range g.Board[x] {
			fmt.Print(g.Board[x][y], " ")
		}
		fmt.Println()
	}
}

func (g *Game) MakeMove(x, y int) {
	g.BoardHistory = append(g.BoardHistory, Step{g.PlayerTurn, x, y})
	g.Board[x][y] = g.PlayerTurn
	g.PlayerTurn = g.PlayerTurn.Opponent()
}

func (g *Game) CheckWin() {
	g.PlayerWon = PlayerNone

	for _, positions := range g.WinPositions {
		var x, y = positions[0][0], positions[0][1]
		var player = g.Board[x][y]
		var count int

		if player == PlayerNone {
			continue
		}

		for _, position := range positions {
			x, y := position[0], position[1]

			if player != g.Board[x][y] {
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

func (g *Game) String() string {
	str := ""
	str += fmt.Sprintf("b{%d,%d,%d};", g.BoardWidth, g.BoardHeight, g.WinLength)
	str += fmt.Sprintf("p{%s};", g.BoardHistory[0].Player.String())

	for _, step := range g.BoardHistory {
		str += fmt.Sprintf("m{%d,%d};", step.X, step.Y)
	}

	str += fmt.Sprintf("r{%s};", g.PlayerWon.String())

	return str
}

func FromString(s string) (*Game, error) {
	var w, h, l int
	_, err := fmt.Sscanf(s, "b{%d,%d,%d};", &w, &h, &l)
	if err != nil {
		return nil, errors.New("invalid game string: " + s)
	}

	g, err := NewGame(w, h, l)
	if err != nil {
		return nil, err
	}

	for _, p := range strings.Split(s, ";") {
		cmd := p[0]

		switch cmd {
		case 'b':
			// ignore
		case 'p':
			var playerChar string
			_, err := fmt.Sscanf(p, "p{%c}", &playerChar)
			if err != nil {
				return nil, errors.New("invalid game string: " + s)
			}

			player, err := StringToPlayer(playerChar)
			if err != nil {
				return nil, errors.New("invalid game string: " + s)
			}

			g.PlayerTurn = player
		case 'm':
			var x, y int
			_, err := fmt.Sscanf(p, "m{%d,%d}", &x, &y)
			if err != nil {
				return nil, errors.New("invalid game string: " + s)
			}

			g.MakeMove(x, y)
		case 'r':
			var result string
			_, err := fmt.Sscanf(p, "r{%c}", &result)
			if err != nil {
				return nil, errors.New("invalid game string: " + s)
			}

			player, err := StringToPlayer(result)
			if err != nil {
				return nil, errors.New("invalid game string: " + s)
			}

			g.PlayerWon = player
		}
	}

	return g, nil
}
