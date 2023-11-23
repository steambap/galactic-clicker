package assets

import (
	_ "embed"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/RobotoCondensed-Medium.ttf
var roboto []byte

var Font14 font.Face
var Font20 font.Face
var Font24 font.Face
var Font36 font.Face

func init() {
	ttf, err := opentype.Parse(roboto)
	if err != nil {
		log.Fatal("fail to parse font")
	}

	const dpi = 72
	Font14, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("fail to construct 14px font")
	}
	Font20, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("fail to construct 20px font")
	}
	Font24, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("fail to construct 24px font")
	}
	Font36, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    36,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("fail to construct 36px font")
	}
}
