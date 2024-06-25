package map_builder

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"tictactoe/game"
	"tictactoe/util"
)

func Write(g *game.Game) error {
	path := GetChunkFilePath(g)
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

func GetChunkFilePath(g *game.Game) string {
	path := ""

	for i := 0; i < len(g.Board)-6; i += 3 {
		path += string(g.Board[i:i+3]) + "/"
	}

	return filepath.Join(
		GetChunksDir(),
		util.GetMapKey(g.BoardWidth, g.BoardHeight, g.WinLength),
		path,
	)
}

func GetChunksDir() string {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get current file path")
	}

	return filepath.Join(filepath.Dir(currentFilePath), "../maps")
}
