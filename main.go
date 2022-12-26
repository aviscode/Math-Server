package main

import (
	"MathServer/server"
)

func main() {
	server := server.InitServer(":8989")
	if err := server.ListenAndServe(); err != nil {
		return
	}
}
