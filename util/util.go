package util

import (
	"fmt"
	"os"
	"os/exec"
)

func Factorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	return uint64(n) * Factorial(n-1)
}

func GetMapKey(w, h, l int) string {
	return fmt.Sprintf("%dx%d_%d", w, h, l)
}

func ClearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
