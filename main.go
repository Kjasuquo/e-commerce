package main

import (
	"fmt"
	"github.com/decadevs/shoparena/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	fmt.Println("Starting... Oja")

	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server in main", err)
		return
	}

}
