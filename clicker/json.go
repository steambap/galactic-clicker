//go:build windows || darwin || linux || android

package clicker

import (
	"encoding/json"
	"log"

	"github.com/steambap/galactic-clicker/storage"
)

func (g *Game) Load() {
	bytes, err := storage.LoadBytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &g.state)
	if err != nil {
		log.Fatalf("Save data corrupted: %v", err)
	}
}

func (g *Game) Save() {
	bytes, err := json.Marshal(&g.state)
	if err != nil {
		log.Fatal(err)
	}

	storage.SaveBytes(bytes)
}
