package clicker

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SHIP_CONTROL_TOP_OFFSET = 70
	SHIP_BOTTOM_PADDING     = 6
	IMG_HEIGHT              = 75
	SHIP_COLUMN_WIDTH       = 410
	BASE_PPI_BONUS          = 2
	BG_LINE_WIDTH           = 2
	GAME_WIDTH              = 1024
	GAME_HEIGHT             = 576
	STATUS_BAR_HEIGHT       = 64
	BOTTOM_BAR_HEIGHT       = 64
)

var GRID_COLOR = color.RGBA{R: 0, G: 0, B: 64, A: 255}

func formatBigFloat(input float64) string {
	digit := math.Log(input) * math.Log10E
	if digit < 6 {
		return fmt.Sprintf("%.f", input)
	}

	return fmt.Sprintf("%.2e", input)
}

type Game struct {
	state          GameState
	buyAmount      int
	shipCost       []float64
	shipControls   []*ShipBuyButton
	shipLevel      []int
	shipLevelMulti []float64
	shipDPS        []float64
}

func (g *Game) Update() error {
	newTime := time.Now()
	dt := newTime.Sub(g.state.LastUpdate).Milliseconds()
	totalShipDPS := 0.0
	for i := 0; i < len(shipDataTable); i++ {
		totalShipDPS += g.shipDPS[i]
	}
	newMoney := totalShipDPS * float64(dt) / 1000
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if captured := g.onPressed(); !captured {
			newMoney += math.Max(totalShipDPS*0.25, 1)
		}
	}
	g.state.CurMoney += newMoney
	g.state.TotalMoney += newMoney
	g.state.LastUpdate = newTime
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, formatBigFloat(g.state.CurMoney))

	g.drawBackground(screen)
	for _, c := range g.shipControls {
		c.Draw(screen, g)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 576
}

func (g *Game) onPressed() (captured bool) {
	x, y := ebiten.CursorPosition()

	for i, c := range g.shipControls {
		if c.In(x, y) {
			g.buyShip(i)
			return true
		}
	}

	return
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	vector.StrokeLine(screen,
		SHIP_COLUMN_WIDTH, 0, SHIP_COLUMN_WIDTH, GAME_HEIGHT, BG_LINE_WIDTH,
		color.White, false)
	vector.StrokeLine(screen,
		0, STATUS_BAR_HEIGHT, GAME_WIDTH, STATUS_BAR_HEIGHT, BG_LINE_WIDTH,
		color.White, false)
	vector.StrokeLine(screen,
		SHIP_COLUMN_WIDTH, GAME_HEIGHT-BOTTOM_BAR_HEIGHT, GAME_WIDTH, GAME_HEIGHT-BOTTOM_BAR_HEIGHT, BG_LINE_WIDTH,
		color.White, false)
}

func (g *Game) buyShip(shipID int) {
	cost := g.calculateCost(shipID)
	if g.state.CurMoney >= cost {
		g.state.CurMoney -= cost
		g.state.ShipCounts[shipID] += g.buyAmount
		g.calculateDPS(shipID)
		g.shipCost[shipID] = g.calculateCost(shipID)
		// save game
	}
}

func (g *Game) calculateCost(shipID int) float64 {
	var sum float64 = 0
	for i := 0; i < g.buyAmount; i++ {
		exp := 1 + float64(shipDataTable[shipID].Exponent)/100
		pow := math.Pow(exp, float64(g.state.ShipCounts[shipID]+i))
		perShipCost := pow * shipDataTable[shipID].BaseCost
		sum += perShipCost
	}

	return sum
}

func (g *Game) calculateDPS(shipID int) {
	g.shipLevelMulti[shipID] = 1
	lv := 0
	for (lv < len(shipLevelCount)) && (shipLevelCount[lv] <= g.state.ShipCounts[shipID]) {
		g.shipLevelMulti[shipID] *= shipLevelMulti[lv]
		lv += 1
	}
	g.shipLevel[shipID] = lv

	shipDPS := float64(g.state.ShipCounts[shipID]) * shipDataTable[shipID].DPS
	// Event bonus
	ppiModifier := 1 + (g.state.PPI * g.getPPIBonus() / 100)
	shipDPS = shipDPS * ppiModifier * g.shipLevelMulti[shipID]
	g.shipDPS[shipID] = shipDPS
}

func (g *Game) getPPIBonus() float64 {
	return float64(BASE_PPI_BONUS)
}

func NewGame() *Game {
	g := &Game{
		state: GameState{
			CurMoney:   4,
			LastUpdate: time.Now(),
			ShipCounts: []int{},
		},
		buyAmount:      1,
		shipCost:       []float64{},
		shipControls:   []*ShipBuyButton{},
		shipLevel:      make([]int, 12),
		shipLevelMulti: make([]float64, 12),
		shipDPS:        make([]float64, 12),
	}

	leftOffset := 5
	topOffset := SHIP_CONTROL_TOP_OFFSET
	for i := 0; i < len(shipDataTable); i++ {
		g.state.ShipCounts = append(g.state.ShipCounts, 0)
		g.shipCost = append(g.shipCost, g.calculateCost(i))
		g.calculateDPS(i)
		sb := &ShipBuyButton{
			Sprite: Sprite{
				image: shipDataTable[i].Img,
				x:     leftOffset,
				y:     topOffset,
			},
			shipID: i,
		}

		g.shipControls = append(g.shipControls, sb)

		if i == 5 {
			leftOffset += SHIP_COLUMN_WIDTH / 2
			topOffset = SHIP_CONTROL_TOP_OFFSET
		} else {
			topOffset += IMG_HEIGHT + SHIP_BOTTOM_PADDING
		}
	}

	return g
}
