package map_reader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"tictactoe/game"
	"tictactoe/map_storage"
)

type Task struct {
	path       string
	game       *game.Game
	resultChan chan Result
}

type Result struct {
	Win  uint64 `json:"win"`
	Lose uint64 `json:"lose"`
	Draw uint64 `json:"draw"`
}

type MapReader struct {
	checkFileChan chan Task
}

func NewMapReader() *MapReader {
	mr := &MapReader{
		checkFileChan: make(chan Task, 1000),
	}

	for i := 0; i < 10; i++ {
		go mr.checkFileWorker()
	}

	return mr
}

func (mr *MapReader) checkFileWorker() {
	for task := range mr.checkFileChan {
		res := Result{}

		file, err := os.OpenFile(task.path, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Println("failed to check the file", task.path, err)
			task.resultChan <- Result{}
			continue
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()

			var winner game.Player
			var gameStr string
			if _, err := fmt.Sscanf(line, "%c %s", &winner, &gameStr); err != nil {
				fmt.Println("failed to parse line", line, err)
				continue
			}

			if compareGamePattern(task.game.String()[2:], gameStr) {
				switch winner {
				case game.PlayerX:
					res.Win++
				case game.PlayerO:
					res.Lose++
				case game.PlayerNone:
					res.Draw++
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("failed to scan file:", err)
		}

		if err := file.Close(); err != nil {
			fmt.Println("failed to close file:", err)
		}

		task.resultChan <- res
	}
}

func (mr *MapReader) GetGameStats(g *game.Game) (Result, error) {
	gameStr := g.String()
	pattern := gameStr[2 : len(gameStr)-6]

	paths, err := map_storage.GetChunkFiles(g)
	if err != nil {
		return Result{}, err
	}

	var pathsFiltered []string

	for _, path := range paths {
		_, fileName := filepath.Split(path)

		if !compareGamePattern(pattern, fileName) {
			continue
		}

		pathsFiltered = append(pathsFiltered, path)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(pathsFiltered))
	resultChan := make(chan Result)
	res := Result{}

	go func() {
		for r := range resultChan {
			res.Win += r.Win
			res.Lose += r.Lose
			res.Draw += r.Draw
			wg.Done()
		}
	}()

	for _, path := range pathsFiltered {
		mr.checkFileChan <- Task{
			path:       path,
			game:       g,
			resultChan: resultChan,
		}
	}

	wg.Wait()
	close(resultChan)

	fmt.Println("processed", len(pathsFiltered), "files for game", g)

	return res, nil
}

func compareGamePattern(pattern, target string) bool {
	if len(pattern) != len(target) {
		return false
	}

	for i, p := range pattern {
		if p == int32(game.PlayerNone) {
			continue
		}

		if target[i] != byte(p) {
			return false
		}
	}

	return true
}
