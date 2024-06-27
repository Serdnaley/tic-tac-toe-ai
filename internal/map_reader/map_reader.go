package map_reader

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"tictactoe/internal/game"
	"tictactoe/internal/map_storage"
	"tictactoe/internal/util"
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

			if util.CompareGamePattern(task.game.String()[2:], gameStr) {
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

		if !util.CompareGamePattern(pattern, fileName) {
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

func (mr *MapReader) GetNextMove(g *game.Game) (int, int, error) {
	wg := &sync.WaitGroup{}
	results := map[int]Result{}

	for i, p := range g.Board {
		if p != game.PlayerNone {
			continue
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			gCopy := g.Copy()
			gCopy.Board[i] = game.PlayerX

			res, err := mr.GetGameStats(gCopy)
			if err != nil {
				fmt.Println("failed to get game stats", err)
				return
			}

			results[i] = res
		}(i)
	}

	wg.Wait()

	haveResults := false
	var bestMove int
	var bestResult Result

	for i, res := range results {
		haveResults = true
		if (bestResult == Result{}) {
			bestMove = i
			bestResult = res
			continue
		}
		if res.Win > bestResult.Win {
			bestMove = i
			bestResult = res
			continue
		}
		if res.Win == bestResult.Win && res.Lose < bestResult.Lose {
			bestMove = i
			bestResult = res
			continue
		}
		if res.Win == bestResult.Win && res.Lose == bestResult.Lose && res.Draw > bestResult.Draw {
			bestMove = i
			bestResult = res
			continue
		}
	}

	if !haveResults {
		return 0, 0, errors.New("no results")
	}

	x := bestMove % g.Size
	y := bestMove / g.Size

	return x, y, nil
}
