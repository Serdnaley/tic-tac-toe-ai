package map_builder

import (
	"fmt"
	"sync"
	"tictactoe/game"
	"tictactoe/util"
	"time"
)

type Task struct {
	game *game.Game
	move int
}

type MapBuilder struct {
	stats *Stats
}

func NewMapBuilder() *MapBuilder {
	mb := &MapBuilder{
		stats: NewStats(),
	}

	return mb
}

func (mb *MapBuilder) BuildWinMap(g *game.Game) error {
	mb.stats.BuildStarted(g)

	wg := &sync.WaitGroup{}
	todoChan := make(chan Task, 1_000_000_000)
	doneChan := make(chan Task)

	for i := 0; i < 2000; i++ {
		go mb.todoWorker(wg, todoChan, doneChan)
	}

	go mb.doneWorker(wg, todoChan, doneChan)

	wg.Add(1)
	doneChan <- Task{game: g, move: -1}

	go func() {
		start := time.Now()
		i := 0
		for {
			i++

			util.ClearConsole()
			fmt.Println("loop:", i)

			mb.stats.Print(start)

			time.Sleep(time.Second / 2)
		}
	}()

	wg.Wait()
	close(todoChan)
	close(doneChan)

	fmt.Println()
	fmt.Println("done")

	return nil
}

func (mb *MapBuilder) todoWorker(
	wg *sync.WaitGroup,
	todoChan chan Task,
	doneChan chan Task,
) {
	for task := range todoChan {
		mb.stats.GameStarted()

		g := task.game.Copy()

		g.MakeMoveByIndex(task.move)

		if g.IsOver() {
			mb.saveResult(g)
		} else {
			wg.Add(1)
			doneChan <- Task{game: g, move: 0}
		}

		wg.Add(1)
		doneChan <- task

		mb.stats.GameFinished()
		wg.Done()
	}
}

func (mb *MapBuilder) doneWorker(
	wg *sync.WaitGroup,
	todoChan chan Task,
	doneChan chan Task,
) {
	for task := range doneChan {
		for i := task.move + 1; i < len(task.game.Board); i++ {
			if task.game.Board[i] == game.PlayerNone {
				wg.Add(1)
				todoChan <- Task{game: task.game, move: i}
				break
			}
		}

		wg.Done()
	}
}

func (mb *MapBuilder) saveResult(g *game.Game) {
	if err := Write(g); err != nil {
		panic(err)
	}
	mb.stats.GamePlayed(g)
}
