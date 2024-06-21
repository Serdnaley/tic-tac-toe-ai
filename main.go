package main

import "tictactoe/predictor"

func main() {
	pm := predictor.NewPredictor()

	err := pm.BuildWinMap(3, 3, 3)
	if err != nil {
		panic(err)
	}
}
