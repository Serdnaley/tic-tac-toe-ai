package map_storage

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"tictactoe/internal/game"
	"time"
)

func RemoveDuplicates(g *game.Game) {
	start := time.Now()
	files, err := GetChunkFiles(g)
	if err != nil {
		return
	}

	fmt.Println("removing duplicates from", len(files), "files")

	wg := &sync.WaitGroup{}
	wg.Add(len(files))

	todoChan := make(chan string)
	doneChan := make(chan struct{})

	for i := 0; i < 10; i++ {
		go todoWorker(todoChan, doneChan)
	}

	go doneWorker(doneChan, wg)

	for _, file := range files {
		todoChan <- file
	}

	wg.Wait()
	close(todoChan)
	close(doneChan)

	fmt.Println("done removing duplicates in", time.Since(start))
}

func todoWorker(todoChan chan string, doneChan chan struct{}) {
	for file := range todoChan {
		if err := removeDuplicates(file); err != nil {
			fmt.Println("error removing duplicates", err)
		}
		doneChan <- struct{}{}
	}
}

func doneWorker(doneChan chan struct{}, wg *sync.WaitGroup) {
	for range doneChan {
		wg.Done()
	}
}

func removeDuplicates(filePath string) error {
	inFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	seenLines := make(map[string]bool)

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		if _, seen := seenLines[line]; !seen {
			seenLines[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error removing input file: %w", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for line := range seenLines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing output file: %w", err)
	}

	return nil
}
