package draw

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	strokeWidth  = 2
	cornerRadius = 4
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	// This is hacky, but WritePixels is better than Fill in term of automatic texture packing.
	whiteImage.WritePixels(pix)
}

func drawVerticesForUtil(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color, antialias bool) {
	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	op.AntiAlias = antialias
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}

func StrokeRoundedRect(dst *ebiten.Image, x, y, width, height float32) {
	offset := float32(math.Min(cornerRadius, math.Min(float64(x/2), float64(y/2))))
	var path vector.Path
	path.MoveTo(x+offset, y)
	path.LineTo(x+width-offset, y)
	path.Arc(
		x+width-offset,
		y+offset,
		offset,
		(math.Pi*3)/2,
		0,
		vector.Clockwise,
	)
	path.LineTo(x+width, y+height-offset)
	path.Arc(
		x+width-offset,
		y+height-offset,
		offset,
		0,
		math.Pi/2,
		vector.Clockwise,
	)
	path.LineTo(x+offset, y+height)
	path.Arc(
		x+offset,
		y+height-offset,
		offset,
		math.Pi/2,
		math.Pi,
		vector.Clockwise,
	)
	path.LineTo(x, y+offset)
	path.Arc(x+offset, y+offset, offset, math.Pi, (math.Pi*3)/2, vector.Clockwise)
	path.Close()

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	strokeOp.MiterLimit = 10
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

	drawVerticesForUtil(dst, vs, is, color.White, true)
}
