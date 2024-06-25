package map_builder

import (
	"fmt"
	"sync/atomic"
	"tictactoe/game"
	"tictactoe/util"
	"time"
)

type Stats struct {
	gamesCountEstimated uint64
	gamesCountElapsed   uint64
	gamesInProgress     uint64
	games               StatsGames
}

type StatsGames struct {
	played uint64
	won    uint64
	lose   uint64
	draw   uint64
}

func NewStats() *Stats {
	s := &Stats{}
	s.Reset()
	return s
}

func (s *Stats) Reset() {
	s.gamesCountEstimated = 0
	s.gamesCountElapsed = 0
	s.gamesInProgress = 0
	s.games.played = 0
	s.games.won = 0
	s.games.lose = 0
	s.games.draw = 0
}

func (s *Stats) GameStarted() {
	atomic.AddUint64(&s.gamesInProgress, 1)
}

func (s *Stats) GameFinished() {
	atomic.AddUint64(&s.gamesInProgress, ^uint64(0))
}

func (s *Stats) BuildStarted(g *game.Game) {
	atomic.AddUint64(&s.gamesCountEstimated, util.Factorial(g.BoardWidth*g.BoardHeight))
	atomic.AddUint64(&s.gamesCountElapsed, util.Factorial(g.BoardWidth*g.BoardHeight))
}

func (s *Stats) GamePlayed(g *game.Game) {
	c := util.Factorial((g.BoardWidth * g.BoardHeight) - g.StepsCount)
	atomic.AddUint64(&s.gamesCountElapsed, ^(c - 1))
	atomic.AddUint64(&s.games.played, 1)

	switch g.PlayerWon {
	case game.PlayerX:
		atomic.AddUint64(&s.games.won, 1)
	case game.PlayerO:
		atomic.AddUint64(&s.games.lose, 1)
	case game.PlayerNone:
		atomic.AddUint64(&s.games.draw, 1)
	}
}

func (s *Stats) GetPercent() float64 {
	return 100 - float64(atomic.LoadUint64(&s.gamesCountElapsed))/float64(atomic.LoadUint64(&s.gamesCountEstimated))*100
}

func (s *Stats) Print(start time.Time) {
	p := s.GetPercent()
	fmt.Println("time elapsed:", time.Since(start))
	fmt.Println("time estimated:", time.Duration(float64(time.Since(start))/p*100))
	fmt.Println("games in progress:", atomic.LoadUint64(&s.gamesInProgress))
	fmt.Println("games estimated:", atomic.LoadUint64(&s.gamesCountEstimated))
	fmt.Println("games played:", atomic.LoadUint64(&s.games.played))
	fmt.Println("games %:", p)
	fmt.Println("games won:", atomic.LoadUint64(&s.games.won))
	fmt.Println("games lose:", atomic.LoadUint64(&s.games.lose))
	fmt.Println("games draw:", atomic.LoadUint64(&s.games.draw))
}
