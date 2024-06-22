package predictor

import (
	"fmt"
	"sync"
	"sync/atomic"
	"tictactoe/game"
	"tictactoe/util"
	"time"
)

type Predictor struct {
	ResultStorage

	CountEstimated  *uint64
	CountElapsed    *uint64
	CountInProgress *uint64
	CountPlayed     *uint64
	CountWon        *uint64
	CountLose       *uint64
	CountDraw       *uint64
}

func NewPredictor() *Predictor {
	p := &Predictor{}

	p.CountEstimated = new(uint64)
	p.CountElapsed = new(uint64)
	p.CountInProgress = new(uint64)
	p.CountPlayed = new(uint64)
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

	count := uint16(0)
	wg2 := &sync.WaitGroup{}
	for x := range g.Board {
		for y := range g.Board[x] {
			if g.Board[x][y] != game.PlayerNone {
				continue
			}

			atomic.AddUint64(pm.CountInProgress, 1)
			wg.Add(1)
			wg2.Add(1)
			count++
			go func(x, y int) {
				newGame := g.Copy()
				newGame.MakeMove(x, y)

				if len(newGame.BoardHistory) > g.WinLength*2-1 {
					newGame.CheckWin()
				}

				pm.play(newGame, st, wg)
				atomic.AddUint64(pm.CountInProgress, ^uint64(0))
				wg.Done()
				wg2.Done()
				count--
			}(x, y)

			if count >= uint16(len(g.BoardHistory)/3+1) {
				wg2.Wait()
			}
		}
	}
}

func (pm *Predictor) putGameResult(g *game.Game, st *ResultStorage) {
	c := util.Factorial((g.BoardWidth * g.BoardHeight) - len(g.BoardHistory))
	atomic.AddUint64(pm.CountElapsed, ^(c - 1))
	atomic.AddUint64(pm.CountPlayed, 1)

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

	start := time.Now()
	wg := &sync.WaitGroup{}

	atomic.AddUint64(pm.CountEstimated, util.Factorial(w*h))
	atomic.AddUint64(pm.CountElapsed, util.Factorial(w*h))
	atomic.AddUint64(pm.CountInProgress, 1)
	wg.Add(1)
	go func() {
		pm.play(g, st, wg)
		atomic.AddUint64(pm.CountInProgress, ^uint64(0))
		wg.Done()
	}()

	go func() {
		i := 0
		for {
			i++
			p := 100 - float64(atomic.LoadUint64(pm.CountElapsed))/float64(atomic.LoadUint64(pm.CountEstimated))*100

			util.ClearConsole()
			fmt.Println("loop:", i)
			fmt.Println("time elapsed:", time.Since(start))
			fmt.Println("time estimated:", time.Duration(float64(time.Since(start))/p*100))
			fmt.Println("games in progress:", atomic.LoadUint64(pm.CountInProgress))
			fmt.Println("games estimated:", atomic.LoadUint64(pm.CountEstimated))
			fmt.Println("games played:", atomic.LoadUint64(pm.CountPlayed))
			fmt.Println("games %:", p)
			fmt.Println("games won:", atomic.LoadUint64(pm.CountWon))
			fmt.Println("games lose:", atomic.LoadUint64(pm.CountLose))
			fmt.Println("games draw:", atomic.LoadUint64(pm.CountDraw))

			time.Sleep(time.Second / 2)
		}
	}()

	wg.Wait()

	fmt.Println()
	fmt.Println("done")

	if err := st.Close(); err != nil {
		return err
	}

	return nil
}
