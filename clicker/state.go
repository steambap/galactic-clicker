package clicker

import (
	"time"
)

type GameState struct {
	Version        string    `json:"version"`
	LastUpdate     time.Time `json:"lastUpdate"`
	CurMoney       float64   `json:"curMoney"`
	TotalMoney     float64   `json:"totalMoney"`
	PlayTime       int       `json:"playTime"`
	PPI            float64   `json:"ppi"`
	ShipCounts     []int     `json:"shipCOunts"`
	EventPurchased []int     `json:"eventPurchased"`
}
