package map_storage

import (
	"os"
	"path/filepath"
	"testing"
	"tictactoe/game"
)

func TestWrite(t *testing.T) {
	// Create a new game
	g, err := game.NewGame(3, 3, 3)
	if err != nil {
		t.Fatalf("Failed to create a new game: %v", err)
	}

	// Write the game state to a file
	err = Write(g)
	if err != nil {
		t.Fatalf("Failed to write game state to file: %v", err)
	}

	// Check if the file was created
	path := getChunkFilePath(g)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Expected file to be created at %s", path)
	}

	// Clean up
	os.RemoveAll(filepath.Dir(path))
}

func TestGetChunkFiles(t *testing.T) {
	// Create a new game
	g, err := game.NewGame(3, 3, 3)
	if err != nil {
		t.Fatalf("Failed to create a new game: %v", err)
	}

	// Write the game state to a file
	err = Write(g)
	if err != nil {
		t.Fatalf("Failed to write game state to file: %v", err)
	}

	// Get chunk files
	files, err := GetChunkFiles(g)
	if err != nil {
		t.Fatalf("Failed to get chunk files: %v", err)
	}

	// Check if the file was retrieved
	if len(files) == 0 {
		t.Fatalf("Expected to retrieve chunk files")
	}

	// Clean up
	os.RemoveAll(filepath.Dir(files[0]))
}

func TestIsMapExist(t *testing.T) {
	// Create a new game
	g, err := game.NewGame(3, 3, 3)
	if err != nil {
		t.Fatalf("Failed to create a new game: %v", err)
	}

	// Write the game state to a file
	err = Write(g)
	if err != nil {
		t.Fatalf("Failed to write game state to file: %v", err)
	}

	// Check if the map exists
	exists := IsMapExist(3, 3, 3)
	if !exists {
		t.Fatalf("Expected map to exist")
	}

	// Clean up
	path := getChunkFilePath(g)
	os.RemoveAll(filepath.Dir(path))
}
