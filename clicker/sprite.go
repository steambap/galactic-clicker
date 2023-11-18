package clicker

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/steambap/galactic-clicker/assets"
)

// Sprite represents an image.
type Sprite struct {
	image *ebiten.Image
	x     int
	y     int
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) In(x, y int) bool {
	_, _, _, a := s.image.At(x-s.x, y-s.y).RGBA()
	return a > 0
}

// Draw draws the sprite.
func (s *Sprite) Draw(screen *ebiten.Image, dx, dy int, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x+dx), float64(s.y+dy))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(s.image, op)
}

type ShipBuyButton struct {
	Sprite
	shipID int
}

func (s *ShipBuyButton) Draw(screen *ebiten.Image, g *Game) {
	s.Sprite.Draw(screen, 0, 0, 1)
	leftOffset := IMG_HEIGHT + 5
	topOffset := 24
	text.Draw(screen, shipDataTable[s.shipID].Name, assets.Font24, s.x+leftOffset, s.y+topOffset, color.White)
	text.Draw(screen, "$"+formatBigFloat(g.shipCost[s.shipID]), assets.Font24, s.x+leftOffset, s.y+topOffset*2, color.White)
	text.Draw(screen, formatBigFloat(g.shipDPS[s.shipID])+" /s", assets.Font24, s.x+leftOffset, s.y+topOffset*3, color.White)
}
