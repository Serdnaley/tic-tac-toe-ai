package main

import (
	"fmt"
	"tictactoe/game"
	"tictactoe/map_builder"
	"time"
)

func main() {
	mb := map_builder.NewMapBuilder()

	start := time.Now()

	g, err := game.NewGame(3, 3, 3)
	if err != nil {
		panic(err)
	}

	if err := mb.BuildWinMap(g); err != nil {
		panic(err)
	}
	fmt.Println("Time elapsed:", time.Since(start))
}
