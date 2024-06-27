package game

import (
	"fmt"
)

type Player byte

const (
	PlayerNone Player = '_'
	PlayerX           = 'X'
	PlayerO           = 'O'
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
