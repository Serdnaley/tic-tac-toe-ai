package main

import (
	"fmt"
	"tictactoe/predictor"
	"time"
)

func main() {
	pm := predictor.NewPredictor()

	start := time.Now()
	err := pm.BuildWinMap(3, 3, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println("Time elapsed:", time.Since(start))
}
