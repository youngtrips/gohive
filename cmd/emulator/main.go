package main

import (
	"gohive/emulator"
)

func main() {
	app := emulator.New()
	app.Run()
	app.Close()
}
