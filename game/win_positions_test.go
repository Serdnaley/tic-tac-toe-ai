package game

import (
	"testing"
	"tictactoe/util"
)

func TestGetWinPositions(t *testing.T) {
	// Test vertical win positions
	verticalWinPositions := GetWinPositions(3, 3, 3)
	expectedVerticalWinPositions := [][]int{
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	if !util.Equal2DSlices(verticalWinPositions, expectedVerticalWinPositions) {
		t.Fatalf("Expected vertical win positions to be %v, got %v", expectedVerticalWinPositions, verticalWinPositions)
	}

	// Test horizontal win positions
	horizontalWinPositions := GetWinPositions(3, 3, 3)
	expectedHorizontalWinPositions := [][]int{
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	if !util.Equal2DSlices(horizontalWinPositions, expectedHorizontalWinPositions) {
		t.Fatalf("Expected horizontal win positions to be %v, got %v", expectedHorizontalWinPositions, horizontalWinPositions)
	}

	// Test diagonal win positions
	diagonalWinPositions := GetWinPositions(3, 3, 3)
	expectedDiagonalWinPositions := [][]int{
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	if !util.Equal2DSlices(diagonalWinPositions, expectedDiagonalWinPositions) {
		t.Fatalf("Expected diagonal win positions to be %v, got %v", expectedDiagonalWinPositions, diagonalWinPositions)
	}

	// Test cache functionality
	cacheKey := util.GetMapKey(3, 3, 3)
	if WinPositionsCache[cacheKey] == nil {
		t.Fatalf("Expected cache to contain key %s", cacheKey)
	}
}
