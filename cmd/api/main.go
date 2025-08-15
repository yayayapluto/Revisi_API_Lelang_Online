package main

import (
	server2 "github.com/API_Lelang_Online_Go/internal/server"
)

func main() {
	server := server2.New()
	server.RegisterFiberServer()
	server.Listen(":8080")
}
