package clicker

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/steambap/galactic-clicker/assets"
)

type ShipData struct {
	Name     string
	Img      *ebiten.Image
	BaseCost float64
	Exponent float64
	DPS      float64
}

var shipDataTable []ShipData = []ShipData{
	{"Scout", nil, 4, 7, 1},
	{"Fighter", nil, 60, 12, 8},
	{"Corvette", nil, 480, 11.5, 48},
	{"Escort", nil, 3840, 11, 288},
	{"Gunship", nil, 30720, 10.5, 1728},
	{"Frigate", nil, 245760, 10, 10368},
	{"Cruiser", nil, 1966080, 9.5, 62208},
	{"Warship", nil, 15728640, 9, 373248},
	{"Destroyer", nil, 125829120, 8.5, 2239488},
	{"Battleship", nil, 1006632960, 8, 13436928},
	{"Carrier", nil, 8053063680, 7.5, 80621568},
	{"Dreadnought", nil, 64424509440, 7, 483729408},
}

var shipLevelCount [32]int = [32]int{
	25, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 1100, 1200, 1300,
	1400, 1500, 1600, 1700, 1800, 1900, 2000, 2100, 2200, 2300, 2400, 2500,
	2600, 2700, 2800, 2900, 3000,
}

var shipLevelMulti [32]float64 = [32]float64{
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 5, 2, 2, 2, 2, 3, 2, 2, 2, 2, 3, 2, 2, 2,
	2, 3, 2, 2, 2, 2, 4,
}

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(assets.ShipBytes))
	if err != nil {
		log.Fatal(err)
	}
	shipImage := ebiten.NewImageFromImage(img)

	for i := range shipDataTable {
		shipDataTable[i].Img = shipImage
	}
}
