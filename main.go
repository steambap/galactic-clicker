package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/steambap/galactic-clicker/clicker"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(clicker.GAME_WIDTH, clicker.GAME_HEIGHT)
	ebiten.SetWindowTitle("Galactic Clicker!")
	g := clicker.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
