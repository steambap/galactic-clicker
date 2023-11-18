package clicker

import (
	"time"
)

type GameState struct {
	Version    string    `json:"version"`
	CurMoney   float64   `json:"curMoney"`
	TotalMoney float64   `json:"totalMoney"`
	LastUpdate time.Time `json:"lastUpdate"`
	ShipCounts []int     `json:"shipCOunts"`
	PPI        float64   `json:"ppi"`
}
