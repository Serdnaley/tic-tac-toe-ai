package map_storage

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"tictactoe/internal/game"
	"tictactoe/internal/util"
)

func GetProgress(g *game.Game) (uint8, bool) {
	path, err := getRelevantProgressFile(g)
	if err != nil {
		log.Fatalf("error getting relevant progress file: %s", err)
	}

	if path == "" {
		return 0, false
	}

	file, err := os.Open(path)
	if err != nil {
		return 0, false
	}

	defer file.Close()

	var str string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str = scanner.Text()
		break
	}

	p, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("error converting string to int: %s", err)
	}

	if p < 0 {
		return 0, true
	}

	if p > 100 {
		return 100, true
	}

	return uint8(p), true
}

func SaveProgress(g *game.Game, p uint8) {
	path := getProgressPath(g)

	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("error removing progress file: %s", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		log.Fatalf("error creating progress dir: %s", err)
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("error creating progress file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(int(p)))
	if err != nil {
		log.Fatalf("error writing progress to file: %s", err)
	}
}

func getRelevantProgressFile(g *game.Game) (string, error) {
	files, err := filepath.Glob(getChunksDir() + "/progress/*")
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if util.CompareGamePattern(filepath.Base(f)[2:], g.String()[2:]) {
			return f, nil
		}
	}

	return "", nil
}

func getProgressPath(g *game.Game) string {
	return filepath.Join(
		getChunksDir(),
		"progress",
		g.String(),
	)
}
