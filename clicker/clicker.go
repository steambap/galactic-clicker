package clicker

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/steambap/galactic-clicker/assets"
)

const (
	SHIP_CONTROL_TOP_OFFSET = 96
	SHIP_BOTTOM_PADDING     = 6
	IMG_SIZE                = 72
	SHIP_COLUMN_WIDTH       = 410
	BASE_PPI_BONUS          = 2
	BG_LINE_WIDTH           = 2
	GAME_WIDTH              = 1024
	GAME_HEIGHT             = 576
	STATUS_BAR_HEIGHT       = 64
	BOTTOM_BAR_HEIGHT       = 64
	NUM_EVENT_BUTTONS       = 5
	EVENT_BUTTON_WIDTH      = 136
	EVENT_BUTTON_HEIGHT     = 72
)

func formatBigFloat(input float64) string {
	digit := math.Log(input) * math.Log10E
	if digit < 6 {
		return fmt.Sprintf("%.f", input)
	}

	return fmt.Sprintf("%.2e", input)
}

type Game struct {
	state             GameState
	buyAmount         int
	shipCost          []float64
	shipLevel         []int
	shipLevelMulti    []float64
	shipDPS           []float64
	newPpi            float64
	finalEventInStage map[int]int
	numEvents         map[int]int
	eventButtons      []Button
	curStage          int
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
	g.newPpi = math.Floor(math.Sqrt(g.state.TotalMoney / 1.0e12))
	g.state.LastUpdate = newTime

	for _, b := range g.eventButtons {
		b.Update(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)
	g.drawTitle(screen)
	for _, b := range g.eventButtons {
		b.Draw(screen, g)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GAME_WIDTH, GAME_HEIGHT
}

func (g *Game) onPressed() (captured bool) {
	x, y := ebiten.CursorPosition()

	for _, b := range g.eventButtons {
		if b.In(x, y) {
			b.OnPressed(g)
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

func (g *Game) drawTitle(screen *ebiten.Image) {
	title := "Galactic Clicker"
	textCenter := getTextCenter(assets.Font36, title)
	text.Draw(screen, title, assets.Font36, SHIP_COLUMN_WIDTH/2-int(textCenter), 36, color.White)

	money := fmt.Sprintf("$ %s", formatBigFloat(g.state.CurMoney))
	textCenter = getTextCenter(assets.Font36, money)
	tx := (GAME_WIDTH+SHIP_COLUMN_WIDTH)/2 - int(textCenter)
	text.Draw(screen, money, assets.Font36, tx, 36, color.White)

	// render month,year
}

func (g *Game) buyShip(shipID int) bool {
	cost := g.calculateCost(shipID)
	if g.state.CurMoney >= cost {
		g.state.CurMoney -= cost
		g.state.ShipCounts[shipID] += g.buyAmount
		g.calculateDPS(shipID)
		g.shipCost[shipID] = g.calculateCost(shipID)
		// save game
		return true
	}

	return false
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

func (g *Game) initAmountButton() {
	ls := []int{1, 10, 100}

	w := 90
	for i, a := range ls {
		b := newBuyAmountButton(float64(w*i+5), 72, a)
		g.eventButtons = append(g.eventButtons, b)
	}
}

func (g *Game) initEventButtons() {
	leftOffset := GAME_WIDTH - EVENT_BUTTON_WIDTH - 4
	topOffset := STATUS_BAR_HEIGHT + 26

	for i := 0; i < NUM_EVENT_BUTTONS; i++ {
		button := newEventButton(leftOffset, topOffset, i)
		g.eventButtons = append(g.eventButtons, button)
		topOffset += EVENT_BUTTON_HEIGHT + 8
	}
}

func NewGame() *Game {
	g := &Game{
		state: GameState{
			CurMoney:   4,
			LastUpdate: time.Now(),
			ShipCounts: []int{},
		},
		buyAmount:         1,
		shipCost:          []float64{},
		shipLevel:         make([]int, 12),
		shipLevelMulti:    make([]float64, 12),
		shipDPS:           make([]float64, 12),
		finalEventInStage: map[int]int{},
		numEvents:         map[int]int{},
		eventButtons:      []Button{},
	}

	leftOffset := 5
	topOffset := SHIP_CONTROL_TOP_OFFSET
	for i := 0; i < len(shipDataTable); i++ {
		g.state.ShipCounts = append(g.state.ShipCounts, 0)
		g.shipCost = append(g.shipCost, g.calculateCost(i))
		g.calculateDPS(i)
		sb := newShipButton(float64(leftOffset), float64(topOffset), i, shipDataTable[i].Img)

		g.eventButtons = append(g.eventButtons, sb)

		if i == 5 {
			leftOffset += SHIP_COLUMN_WIDTH / 2
			topOffset = SHIP_CONTROL_TOP_OFFSET
		} else {
			topOffset += IMG_SIZE + SHIP_BOTTOM_PADDING
		}
	}

	prevStage := EARTH
	for i := 0; i < len(eventDataTable); i++ {
		stage := eventDataTable[i].Stage
		if _, ok := g.numEvents[stage]; !ok {
			g.numEvents[stage] = 0
		}
		g.numEvents[stage] += 1
		if stage != prevStage {
			g.finalEventInStage[prevStage] = i - 1
			prevStage = stage
		}
	}
	g.finalEventInStage[GALAXY] = len(eventDataTable) - 1
	g.initAmountButton()
	g.initEventButtons()

	return g
}

func init() {
	for i := range shipDataTable {
		shipDataTable[i].Img = ebiten.NewImageFromImage(assets.ShipImageList[i])
	}

	grayScale = colorm.ColorM{}
	grayScale.ChangeHSV(0, 0, 1)
}
