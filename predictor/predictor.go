package predictor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"tictactoe/game"
)

type Predictor struct {
	winMap map[string]map[string]game.Player
}

func NewPredictor() *Predictor {
	pm := &Predictor{}
	pm.winMap = make(map[string]map[string]game.Player)
	return pm
}

func (pm *Predictor) play(g *game.Game) {
	if g.PlayerWon != game.PlayerNone || len(g.BoardHistory) == g.BoardWidth*g.BoardHeight {
		pm.putGameResult(g)
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

			pm.play(newGame)
		}
	}
}

func (pm *Predictor) putGameResult(g *game.Game) {
	boardKey := fmt.Sprintf("%d-%d-%d", g.BoardWidth, g.BoardHeight, g.WinLength)

	historyKey := fmt.Sprintf("p{%s};", g.BoardHistory[0].Player.String())
	for _, step := range g.BoardHistory {
		historyKey += fmt.Sprintf("m{%d,%d};", step.X, step.Y)
	}

	if pm.winMap[boardKey] == nil {
		pm.winMap[boardKey] = make(map[string]game.Player)
	}
	pm.winMap[boardKey][historyKey] = g.PlayerWon
}

func (pm *Predictor) BuildWinMap(w, h, l int) error {
	g, err := game.NewGame(w, h, l)
	if err != nil {
		return err
	}

	pm.play(g)
	if err = pm.saveWinMap(w, h, l); err != nil {
		return err
	}

	return nil
}

func (pm *Predictor) saveWinMap(w, h, l int) error {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	key := fmt.Sprintf("%d-%d-%d", w, h, l)
	name := fmt.Sprintf("%s.txt", key)
	dir := filepath.Join(filepath.Dir(currentFilePath), "maps")
	path := filepath.Join(dir, name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("failed to close file:", err)
		}
	}()

	for k, v := range pm.winMap[key] {
		str := fmt.Sprintf("%sr{%s};\n", k, v.String())
		if _, err = file.WriteString(str); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	fmt.Println("predictor saved to", path)

	return nil
}
