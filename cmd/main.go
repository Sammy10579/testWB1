package main

import (
	"log"

	"testWB/app"
)

func main() {
	a := app.New()
	go func() {
		if err := a.RunSet(":8081"); err != nil {
			log.Print(err)
		}
	}()

	if err := a.RunGet(":8080"); err != nil {
		log.Print(err)
	}
}
