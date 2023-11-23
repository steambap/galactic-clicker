package assets

import (
	"bytes"
	_ "embed"
	"image"
	"log"
)

//go:embed images/battleship.png
var BattleshipBytes []byte

//go:embed images/carrier.png
var CarrierBytes []byte

//go:embed images/corvette.png
var CorvetteBytes []byte

//go:embed images/cruiser.png
var CruiserBytes []byte

//go:embed images/destroyer.png
var DestroyerBytes []byte

//go:embed images/dreadnought.png
var DreadnoughtBytes []byte

//go:embed images/escort.png
var EscortBytes []byte

//go:embed images/fighter.png
var FighterBytes []byte

//go:embed images/frigate.png
var FrigateBytes []byte

//go:embed images/gunship.png
var GunshipBytes []byte

//go:embed images/scout.png
var ScoutBytes []byte

//go:embed images/warship.png
var WarshipBytes []byte

var ShipImageList []image.Image

func init() {
	battleship, _, err := image.Decode(bytes.NewReader(BattleshipBytes))
	if err != nil {
		log.Fatal(err)
	}
	carrier, _, err := image.Decode(bytes.NewReader(CarrierBytes))
	if err != nil {
		log.Fatal(err)
	}
	corvette, _, err := image.Decode(bytes.NewReader(CorvetteBytes))
	if err != nil {
		log.Fatal(err)
	}
	cruiser, _, err := image.Decode(bytes.NewReader(CruiserBytes))
	if err != nil {
		log.Fatal(err)
	}
	destroyer, _, err := image.Decode(bytes.NewReader(DestroyerBytes))
	if err != nil {
		log.Fatal(err)
	}
	dreadnought, _, err := image.Decode(bytes.NewReader(DreadnoughtBytes))
	if err != nil {
		log.Fatal(err)
	}
	escort, _, err := image.Decode(bytes.NewReader(EscortBytes))
	if err != nil {
		log.Fatal(err)
	}
	fighter, _, err := image.Decode(bytes.NewReader(FighterBytes))
	if err != nil {
		log.Fatal(err)
	}
	frigate, _, err := image.Decode(bytes.NewReader(FrigateBytes))
	if err != nil {
		log.Fatal(err)
	}
	gunship, _, err := image.Decode(bytes.NewReader(GunshipBytes))
	if err != nil {
		log.Fatal(err)
	}
	scout, _, err := image.Decode(bytes.NewReader(ScoutBytes))
	if err != nil {
		log.Fatal(err)
	}
	warship, _, err := image.Decode(bytes.NewReader(WarshipBytes))
	if err != nil {
		log.Fatal(err)
	}

	ShipImageList = []image.Image{
		scout, fighter, corvette, escort, gunship, frigate,
		cruiser, warship, destroyer, battleship, carrier, dreadnought,
	}
}
