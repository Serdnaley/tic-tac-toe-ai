package predictor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"tictactoe/game"
	"tictactoe/util"
)

type ResultStorage struct {
	Width       int
	Height      int
	WinLength   int
	ChunksCount int
	Chunks      map[int]*os.File
}

func NewResultStorage(w, h, l int) (*ResultStorage, error) {
	count := int((util.Factorial(w) + util.Factorial(h) + util.Factorial(l)) / 3)

	st := &ResultStorage{
		Width:       w,
		Height:      h,
		WinLength:   l,
		ChunksCount: count,
		Chunks:      make(map[int]*os.File, count),
	}

	dir := st.GetChunksDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	return st, nil
}

func (st *ResultStorage) Create() error {
	for i := 0; i < st.ChunksCount; i++ {
		file, err := os.Create(st.GetChunkFilePath(i))
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		st.Chunks[i] = file
	}

	return nil
}

func (st *ResultStorage) Close() error {
	for _, file := range st.Chunks {
		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close file: %v", err)
		}
	}

	return nil
}

func (st *ResultStorage) Write(g *game.Game) {
	if _, ok := st.Chunks[st.GetGameChunk(g)]; !ok {
		panic(fmt.Errorf("chunk is not found: %d", st.GetGameChunk(g)))
	}

	str := fmt.Sprintf("p{%s};", g.BoardHistory[0].Player.String())
	for _, step := range g.BoardHistory {
		str += fmt.Sprintf("m{%d,%d};", step.X, step.Y)
	}
	str += fmt.Sprintf("r{%s};", g.PlayerWon.String())
	str += "\n"

	file := st.Chunks[st.GetGameChunk(g)]
	if _, err := file.WriteString(str); err != nil {
		panic(err)
	}
}

func (st *ResultStorage) GetGameChunk(g *game.Game) int {
	r := util.Murmur3Hash32([]byte(g.String()), 0)

	return int(r % uint32(st.ChunksCount))
}

func (st *ResultStorage) GetChunksDir() string {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get current file path")
	}

	return filepath.Join(filepath.Dir(currentFilePath), "maps", util.GetMapKey(st.Width, st.Height, st.WinLength))
}

func (st *ResultStorage) GetChunkFileName(n int) string {
	return fmt.Sprintf("chunk_%d.txt", n)
}

func (st *ResultStorage) GetChunkFilePath(n int) string {

	return filepath.Join(st.GetChunksDir(), st.GetChunkFileName(n))
}
