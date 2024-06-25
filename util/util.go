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
func ParseMapKey(mapKey string) (int, int, int, error) {
	var w, h, l int
	if _, err := fmt.Sscanf(mapKey, "%dx%d_%d", &w, &h, &l); err != nil {
		return 0, 0, 0, err
	}
	return w, h, l, nil
}

func ClearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// Equal2DSlices compares two 2D slices for equality
func Equal2DSlices(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
