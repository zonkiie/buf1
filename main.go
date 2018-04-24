package main

import (
	"log"

	"buf1/actions"
)

func main() {
	app := actions.App()
	log.Fatal(app.Serve())
}
