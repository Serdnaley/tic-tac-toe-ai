package map_builder

import (
	"fmt"
	"sync"
	"tictactoe/internal/game"
	"tictactoe/internal/map_storage"
	"time"
)

type Task struct {
	wg   *sync.WaitGroup
	game *game.Game
	move int
}

type MapBuilder struct {
	buildWinMapChan chan *game.Game
	todoChan        chan Task
	doneChan        chan Task
	stats           *Stats
}

func NewMapBuilder() *MapBuilder {
	mb := &MapBuilder{
		buildWinMapChan: make(chan *game.Game, 1),
		todoChan:        make(chan Task, 1_000_000_000),
		doneChan:        make(chan Task),
		stats:           NewStats(),
	}

	for i := 0; i < 2000; i++ {
		go mb.todoWorker()
	}

	go mb.doneWorker()
	go mb.buildWinMapWorker()

	return mb
}

func (mb *MapBuilder) BuildWinMap(g *game.Game) error {
	map_storage.SaveProgress(g, 0)

	if g.IsOver() {
		return fmt.Errorf("game is already over")
	}

	mb.buildWinMapChan <- g

	return nil
}

type GetMapStatusResponse struct {
	Status   string `json:"status"`
	Progress uint8  `json:"progress"`
}

func (mb *MapBuilder) buildWinMapWorker() {
	for g := range mb.buildWinMapChan {
		start := time.Now()
		fmt.Println("building win map started", g)

		mb.stats.Reset()
		mb.stats.BuildStarted(g)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		mb.doneChan <- Task{wg: wg, game: g, move: -1}

		stopped := false

		go func() {
			for {
				if stopped {
					break
				}

				mb.stats.Print(start)
				map_storage.SaveProgress(g, uint8(mb.stats.GetPercent()))
				time.Sleep(time.Second)
			}
		}()

		wg.Wait()
		stopped = true

		mb.stats.Print(start)
		fmt.Println("building win map finished in", time.Since(start), g)

		map_storage.RemoveDuplicates(g)
		map_storage.SaveProgress(g, 100)
	}
}

func (mb *MapBuilder) todoWorker() {
	for task := range mb.todoChan {
		mb.stats.GameStarted()

		g := task.game.Copy()

		g.MakeMoveByIndex(task.move)

		if g.IsOver() {
			mb.saveResult(g)
		} else {
			task.wg.Add(1)
			mb.doneChan <- Task{wg: task.wg, game: g, move: 0}
		}

		task.wg.Add(1)
		mb.doneChan <- task

		mb.stats.GameFinished()
		task.wg.Done()
	}
}

func (mb *MapBuilder) doneWorker() {
	for task := range mb.doneChan {
		for i := task.move + 1; i < len(task.game.Board); i++ {
			if task.game.Board[i] == game.PlayerNone {
				task.wg.Add(1)
				mb.todoChan <- Task{wg: task.wg, game: task.game, move: i}
				break
			}
		}

		task.wg.Done()
	}
}

func (mb *MapBuilder) saveResult(g *game.Game) {
	if err := map_storage.Write(g); err != nil {
		panic(err)
	}
	mb.stats.GamePlayed(g)
}
