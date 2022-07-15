package main

import (
	"go-IM/src/server"
)

func main() {
	server := server.NewServer("127.0.0.1", 8888)
	server.Start()
}
