package game

import (
	"fmt"
)

type Player int8

const (
	PlayerNone Player = iota
	PlayerX
	PlayerO
)

func (p Player) Opponent() Player {
	switch p {
	case PlayerNone:
		return PlayerNone
	case PlayerX:
		return PlayerO
	case PlayerO:
		return PlayerX
	default:
		panic("invalid player: " + fmt.Sprint(p))
	}
}

func (p Player) String() string {
	switch p {
	case PlayerNone:
		return "-"
	case PlayerX:
		return "X"
	case PlayerO:
		return "O"
	default:
		panic("invalid player: " + fmt.Sprint(p))
	}
}

func StringToPlayer(c string) (Player, error) {
	switch c {
	case "-":
		return PlayerNone, nil
	case "X":
		return PlayerX, nil
	case "O":
		return PlayerO, nil
	default:
		return PlayerNone, fmt.Errorf("invalid player: %s", c)
	}
}
