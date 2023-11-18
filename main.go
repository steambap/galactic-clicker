package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/steambap/galactic-clicker/clicker"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(1024, 576)
	ebiten.SetWindowTitle("Galactic Clicker!")
	g := clicker.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
