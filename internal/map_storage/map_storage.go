package map_storage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"tictactoe/internal/game"
	"tictactoe/internal/util"
)

func Write(g *game.Game) error {
	path := getChunkFilePath(g)
	var file *os.File
	var err error

	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		file, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
	} else {
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
	}

	if _, err = file.WriteString(g.String() + "\n"); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %v", err)
	}

	return nil
}

func GetChunkFiles(g *game.Game) ([]string, error) {
	return filepath.Glob(
		filepath.Join(
			getChunksDir(),
			util.GetMapKey(g.Size, g.WinLength),
			"*",
		),
	)
}

func IsMapExist(s, l int) bool {
	_, err := os.Stat(getChunksDir() + "/" + util.GetMapKey(s, l))
	return !os.IsNotExist(err)
}

func getChunkFilePath(g *game.Game) string {
	return filepath.Join(
		getChunksDir(),
		util.GetMapKey(g.Size, g.WinLength),
		string(g.Board[0:len(g.Board)-6]),
	)
}

func getChunksDir() string {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get current file path")
	}

	return filepath.Join(filepath.Dir(currentFilePath), "../../../tic-tac-toe-maps")
}
