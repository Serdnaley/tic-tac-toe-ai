package game

import (
	"tictactoe/internal/util"
)

var WinPositionsCache = map[string][][]int{}

func GetWinPositions(s, l int) [][]int {
	cacheKey := util.GetMapKey(s, l)

	if WinPositionsCache[cacheKey] != nil {
		return WinPositionsCache[cacheKey]
	}

	res := make([][]int, 0, s*s)

	// Vertical
	for xOffset := 0; xOffset <= s-l; xOffset++ {
		for x := 0; x < l; x++ {
			var column []int

			for y := 0; y < l; y++ {
				column = append(column, (y*s)+(x+xOffset))
			}

			res = append(res, column)
		}
	}

	// Horizontal
	for yOffset := 0; yOffset <= s-l; yOffset++ {
		for y := 0; y < l; y++ {
			var row []int

			for x := 0; x < l; x++ {
				row = append(row, (y+yOffset)*s+x)
			}

			res = append(res, row)
		}
	}

	// Diagonal
	for xOffset := 0; xOffset <= s-l; xOffset++ {
		for yOffset := 0; yOffset <= s-l; yOffset++ {
			var diagonal1 []int
			var diagonal2 []int

			for i := 0; i < l; i++ {
				diagonal1 = append(diagonal1, xOffset+i+(yOffset+i)*s)
				diagonal2 = append(diagonal2, xOffset+l-1-i+(yOffset+i)*s)
			}

			res = append(res, diagonal1)
			res = append(res, diagonal2)
		}
	}

	WinPositionsCache[cacheKey] = res

	return res
}
