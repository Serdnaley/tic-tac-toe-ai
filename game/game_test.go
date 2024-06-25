package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game, err := NewGame(3, 3, 3)
	if err != nil {
		t.Fatalf("Failed to create a new game: %v", err)
	}

	if game.BoardWidth != 3 {
		t.Fatalf("Expected board width to be 3, got %d", game.BoardWidth)
	}

	if game.BoardHeight != 3 {
		t.Fatalf("Expected board height to be 3, got %d", game.BoardHeight)
	}

	if game.WinLength != 3 {
		t.Fatalf("Expected win length to be 3, got %d", game.WinLength)
	}
}

func TestCopy(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	copy := game.Copy()

	if copy.BoardWidth != game.BoardWidth || copy.BoardHeight != game.BoardHeight || copy.WinLength != game.WinLength {
		t.Fatalf("Copy did not preserve game dimensions")
	}

	if copy.PlayerTurn != game.PlayerTurn || copy.PlayerWon != game.PlayerWon || copy.StepsCount != game.StepsCount {
		t.Fatalf("Copy did not preserve game state")
	}

	for i := range game.Board {
		if copy.Board[i] != game.Board[i] {
			t.Fatalf("Copy did not preserve board state")
		}
	}
}

func TestGetMapKey(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	key := game.GetMapKey()

	expectedKey := "3x3_3"
	if key != expectedKey {
		t.Fatalf("Expected map key to be %s, got %s", expectedKey, key)
	}
}

func TestMakeMoveByIndex(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	game.MakeMoveByIndex(0)

	if game.Board[0] != PlayerX {
		t.Fatalf("Expected PlayerX at index 0, got %v", game.Board[0])
	}

	if game.PlayerTurn != PlayerO {
		t.Fatalf("Expected PlayerO's turn, got %v", game.PlayerTurn)
	}
}

func TestMakeMoveByCoordinates(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	game.MakeMoveByCoordinates(1, 1)

	if game.Board[4] != PlayerX {
		t.Fatalf("Expected PlayerX at coordinates (1,1), got %v", game.Board[4])
	}

	if game.PlayerTurn != PlayerO {
		t.Fatalf("Expected PlayerO's turn, got %v", game.PlayerTurn)
	}
}

func TestCheckWin(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	game.MakeMoveByIndex(0)
	game.MakeMoveByIndex(1)
	game.MakeMoveByIndex(3)
	game.MakeMoveByIndex(4)
	game.MakeMoveByIndex(6)

	if game.PlayerWon != PlayerX {
		t.Fatalf("Expected PlayerX to win, got %v", game.PlayerWon)
	}
}

func TestScaleBoard(t *testing.T) {
	game, _ := NewGame(3, 3, 3)

	for i := 0; i < 9; i++ {
		game.MakeMoveByIndex(i)
	}

	err := game.ScaleBoard(5, 5, 4)
	if err != nil {
		t.Fatalf("Failed to scale the board: %v", err)
	}

	if game.BoardWidth != 5 {
		t.Fatalf("Expected board width to be 5, got %d", game.BoardWidth)
	}

	if game.BoardHeight != 5 {
		t.Fatalf("Expected board height to be 5, got %d", game.BoardHeight)
	}

	if game.WinLength != 4 {
		t.Fatalf("Expected win length to be 4, got %d", game.WinLength)
	}

	expectedBoard := []Player{
		PlayerNone, PlayerNone, PlayerNone, PlayerNone, PlayerNone,
		PlayerNone, PlayerX, PlayerO, PlayerX, PlayerNone,
		PlayerNone, PlayerO, PlayerX, PlayerO, PlayerNone,
		PlayerNone, PlayerX, PlayerO, PlayerX, PlayerNone,
		PlayerNone, PlayerNone, PlayerNone, PlayerNone, PlayerNone,
	}

	for i := range game.Board {
		if game.Board[i] != expectedBoard[i] {
			t.Fatalf("Expected board to be %v, got %v", expectedBoard, game.Board)
		}
	}
}

func TestIsFulfilled(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	for i := 0; i < 9; i++ {
		game.MakeMoveByIndex(i)
	}

	if !game.IsFulfilled() {
		t.Fatalf("Expected game to be fulfilled")
	}
}

func TestIsOver(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	game.MakeMoveByIndex(0)
	game.MakeMoveByIndex(1)
	game.MakeMoveByIndex(3)
	game.MakeMoveByIndex(4)
	game.MakeMoveByIndex(6)

	if !game.IsOver() {
		t.Fatalf("Expected game to be over")
	}
}

func TestString(t *testing.T) {
	game, _ := NewGame(3, 3, 3)
	game.MakeMoveByIndex(0)
	game.MakeMoveByIndex(1)
	game.MakeMoveByIndex(3)
	game.MakeMoveByIndex(4)
	game.MakeMoveByIndex(6)

	expectedString := "X XO_XO_X__"
	if game.String() != expectedString {
		t.Fatalf("Expected game string to be %s, got %s", expectedString, game.String())
	}
}

func TestFromString(t *testing.T) {
	game, err := FromString("X XO_XO_X__")
	if err != nil {
		t.Fatalf("Failed to create a new game from string: %v", err)
	}

	if game.BoardWidth != 3 {
		t.Fatalf("Expected board width to be 3, got %d", game.BoardWidth)
	}

	if game.BoardHeight != 3 {
		t.Fatalf("Expected board height to be 3, got %d", game.BoardHeight)
	}

	if game.WinLength != 3 {
		t.Fatalf("Expected win length to be 3, got %d", game.WinLength)
	}

	if game.PlayerWon != PlayerX {
		t.Fatalf("Expected PlayerX to have won, got %v", game.PlayerWon)
	}
}
