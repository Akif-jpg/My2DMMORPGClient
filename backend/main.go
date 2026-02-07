package main

import (
	gameapp "github.com/Akif-jpg/MyHobieMMORPGGame/app"
)

func main() {
	app := gameapp.New()
	app.Init(gameapp.DEVELOPER)
}
