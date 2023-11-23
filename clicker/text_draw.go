package clicker

import (
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func getTextCenter(f font.Face, s string) float64 {
	advance := font.MeasureString(f, s)

	return fixed26_6ToFloat64(advance) / 2
}

func fixed26_6ToFloat64(x fixed.Int26_6) float64 {
	return float64(x>>6) + float64(x&((1<<6)-1))/float64(1<<6)
}

func measure(f font.Face, s string) (w, h float64) {
	bound := text.BoundString(f, s)

	w = float64(bound.Dx())
	h = float64(bound.Dy())
	return
}
