package predictor

import (
	"fmt"
	"sync"
	"sync/atomic"
	"tictactoe/game"
	"time"
)

type Predictor struct {
	ResultStorage

	Count     *uint64
	CountWon  *uint64
	CountLose *uint64
	CountDraw *uint64
}

func NewPredictor() *Predictor {
	p := &Predictor{}

	p.Count = new(uint64)
	p.CountWon = new(uint64)
	p.CountLose = new(uint64)
	p.CountDraw = new(uint64)

	return p
}

func (pm *Predictor) play(g *game.Game, st *ResultStorage, wg *sync.WaitGroup) {
	if g.PlayerWon != game.PlayerNone || len(g.BoardHistory) == g.BoardWidth*g.BoardHeight {
		pm.putGameResult(g, st)
		return
	}

	for x := range g.Board {
		for y := range g.Board[x] {
			if g.Board[x][y] != game.PlayerNone {
				continue
			}

			newGame := g.Copy()
			newGame.MakeMove(x, y)

			if len(newGame.BoardHistory) > g.WinLength*2-1 {
				newGame.CheckWin()
			}

			wg.Add(1)
			go func() {
				pm.play(newGame, st, wg)
				wg.Done()
			}()
		}
	}
}

func (pm *Predictor) putGameResult(g *game.Game, st *ResultStorage) {
	atomic.AddUint64(pm.Count, 1)
	switch g.PlayerWon {
	case game.PlayerX:
		atomic.AddUint64(pm.CountWon, 1)
	case game.PlayerO:
		atomic.AddUint64(pm.CountLose, 1)
	case game.PlayerNone:
		atomic.AddUint64(pm.CountDraw, 1)
	}

	st.Write(g)
}

func (pm *Predictor) BuildWinMap(w, h, l int) error {
	st, err := NewResultStorage(w, h, l)
	if err != nil {
		return err
	}

	if err := st.Create(); err != nil {
		return err
	}

	g, err := game.NewGame(w, h, l)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		pm.play(g, st, wg)
		wg.Done()
	}()

	go func() {
		i := 0
		for {
			i++
			fmt.Println()
			fmt.Println("loop:", i)
			fmt.Println("games played:", atomic.LoadUint64(pm.Count))
			fmt.Println("games won:", atomic.LoadUint64(pm.CountWon))
			fmt.Println("games lose:", atomic.LoadUint64(pm.CountLose))
			fmt.Println("games draw:", atomic.LoadUint64(pm.CountDraw))
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	fmt.Println()
	fmt.Println("games played:", pm.Count)
	fmt.Println("games won:", pm.CountWon)
	fmt.Println("games lose:", pm.CountLose)
	fmt.Println("games draw:", pm.CountDraw)
	fmt.Println()
	fmt.Println("done")

	if err := st.Close(); err != nil {
		return err
	}

	return nil
}
