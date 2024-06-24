package game

import (
	"tictactoe/util"
)

var WinPositionsCache = map[string][][]int{}

func GetWinPositions(w, h, l int) [][]int {
	cacheKey := util.GetMapKey(w, h, l)

	if WinPositionsCache[cacheKey] != nil {
		return WinPositionsCache[cacheKey]
	}

	res := make([][]int, 0, w*h)

	// Vertical
	for xOffset := 0; xOffset <= w-l; xOffset++ {
		for x := 0; x < l; x++ {
			var column []int

			for y := 0; y < l; y++ {
				column = append(column, (y*w)+(x+xOffset))
			}

			res = append(res, column)
		}
	}

	// Horizontal
	for yOffset := 0; yOffset <= h-l; yOffset++ {
		for y := 0; y < l; y++ {
			var row []int

			for x := 0; x < l; x++ {
				row = append(row, (y+yOffset)*w+x)
			}

			res = append(res, row)
		}
	}

	// Diagonal
	for xOffset := 0; xOffset <= w-l; xOffset++ {
		for yOffset := 0; yOffset <= h-l; yOffset++ {
			var diagonal1 []int
			var diagonal2 []int

			for i := 0; i < l; i++ {
				diagonal1 = append(diagonal1, xOffset+i+(yOffset+i)*w)
				diagonal2 = append(diagonal2, xOffset+l-1-i+(yOffset+i)*w)
			}

			res = append(res, diagonal1)
			res = append(res, diagonal2)
		}
	}

	WinPositionsCache[cacheKey] = res

	return res
}
