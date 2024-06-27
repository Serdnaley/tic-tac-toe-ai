package main

import "tictactoe/internal/server"

func main() {
	server.NewServer().Start(4000)
}
