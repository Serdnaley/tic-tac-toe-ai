package main

import (
	"fmt"
	"tictactoe/predictor"
	"time"
)

func main() {
	pm := predictor.NewPredictor()

	start := time.Now()
	if err := pm.BuildWinMap(4, 4, 2); err != nil {
		panic(err)
	}
	fmt.Println("Time elapsed:", time.Since(start))
}
