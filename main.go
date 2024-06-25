package main

import "tictactoe/server"

func main() {
	server.NewServer().Start(4000)
}
