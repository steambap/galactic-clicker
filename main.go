package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func formatBigFloat(input float64) string {
	digit := math.Log(input) * math.Log10E
	if digit < 6 {
		return fmt.Sprintf("%.f", input)
	}

	return fmt.Sprintf("%.2e", input)
}

type Game struct {
	CurMoney   float64
	TotalMoney float64
	LastUpdate time.Time
}

func (g *Game) Update() error {
	newTime := time.Now()
	dt := newTime.Sub(g.LastUpdate).Milliseconds()
	newMoney := 1 * float64(dt) / 1000
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		newMoney += 1
	}
	g.CurMoney += newMoney
	g.TotalMoney += newMoney
	g.LastUpdate = newTime
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, formatBigFloat(g.CurMoney))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{
		LastUpdate: time.Now(),
	}); err != nil {
		log.Fatal(err)
	}
}
