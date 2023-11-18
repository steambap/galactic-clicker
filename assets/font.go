package assets

import (
	_ "embed"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/Kenney_High.ttf
var kenneyFont []byte

var Font16 font.Face
var Font24 font.Face

func init() {
	ttf, err := opentype.Parse(kenneyFont)
	if err != nil {
		log.Fatal("fail to parse font")
	}

	const dpi = 72
	Font16, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("fail to construct 16px font")
	}
	Font24, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("fail to construct 24px font")
	}
}
