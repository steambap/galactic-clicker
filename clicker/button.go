package clicker

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/steambap/galactic-clicker/assets"
	"github.com/steambap/galactic-clicker/draw"
	"golang.org/x/image/font"
)

var grayScale colorm.ColorM
var defaultScale = colorm.ColorM{}
var disableColor color.Color = color.RGBA{64, 64, 64, 127}

type Button interface {
	In(x, y int) bool
	Draw(screen *ebiten.Image, g *Game)
	Update(g *Game)
	OnPressed(g *Game) bool
}

type BaseButton struct {
	x      float64
	y      float64
	width  float64
	height float64
}

func (b *BaseButton) In(x, y int) bool {
	x1, y1 := float64(x), float64(y)

	return x1 >= b.x && x1 <= (b.x+b.width) &&
		y1 >= b.y && y1 <= (b.y+b.height)
}

func (b *BaseButton) Draw(screen *ebiten.Image) {
	draw.StrokeRoundedRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height))
}

func (b *BaseButton) Update(g *Game) {}

func (b *BaseButton) OnPressed(g *Game) bool {
	return false
}

type EventButton struct {
	BaseButton
	eventID int
}

func (b *EventButton) isVisible(g *Game) bool {
	if b.eventID < 0 || b.eventID >= len(eventDataTable) {
		return false
	}
	event := eventDataTable[b.eventID]

	return event.Stage <= g.curStage
}

func (b *EventButton) Draw(screen *ebiten.Image, g *Game) {
	if !b.isVisible(g) {
		return
	}

	event := eventDataTable[b.eventID]
	if g.state.CurMoney < event.Cost {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), disableColor, true)
	}
	b.BaseButton.Draw(screen)

	// Title
	text.Draw(screen, event.Title, assets.Font20, int(b.x)+2, int(b.y)+18, color.White)
	// Cost
	costLabel := fmt.Sprintf("$%s", formatBigFloat(event.Cost))
	text.Draw(screen, costLabel, assets.Font14, int(b.x)+2, int(b.y)+36, color.White)
	// Hint
	hintLabel := ""
	if event.Type1 < len(shipDataTable) {
		// ship upgrade
		hintLabel = fmt.Sprintf("%s x%v", shipDataTable[event.Type1].Name, int(event.Type2))
	} else if event.Type1 == 12 {
		// production multiplier
		hintLabel = fmt.Sprintf("Production x%v", int(event.Type2))
	} else {
		hintLabel = fmt.Sprintf("PPI +%.1f%%", event.Type2)
	}
	text.Draw(screen, hintLabel, assets.Font14, int(b.x)+2, int(b.y)+54, color.White)

	// Feature
	featLabel := ""
	if event.Feature != "" {
		featLabel = "!"
	}
	text.Draw(screen, featLabel, assets.Font24, int(b.x)+EVENT_BUTTON_WIDTH-10, int(b.y)+EVENT_BUTTON_HEIGHT-14, color.White)
}

func (b *EventButton) OnPressed(g *Game) bool {
	event := eventDataTable[b.eventID]
	if g.state.CurMoney > event.Cost {
		g.state.CurMoney -= event.Cost
		g.state.EventPurchased[b.eventID] = true
		stage := event.Stage
		g.numEvents[stage] += 1
		if event.Feature != "" {
			g.renderFeature[event.Feature] = true
		}
		for i := 0; i < len(shipDataTable); i++ {
			g.calculateDPS(i)
			g.shipCost[i] = g.calculateCost(i)
		}
		if g.finalEventInStage[g.curStage] == b.eventID && g.curStage != GALAXY {
			g.curStage += 1
		}
		g.refreshEventButtons()
		// save game

		return true
	} else {
		g.state.CurMoney += event.Cost
	}

	return false
}

func newEventButton(x, y int, eventID int) *EventButton {
	return &EventButton{
		BaseButton: BaseButton{
			x:      float64(x),
			y:      float64(y),
			width:  EVENT_BUTTON_WIDTH - 2,
			height: EVENT_BUTTON_HEIGHT,
		},
		eventID: eventID,
	}
}

type BuyAmountButton struct {
	BaseButton
	amount int
	// font width / 2
	dx float64
}

func (b *BuyAmountButton) OnPressed(g *Game) bool {
	if b.amount == g.buyAmount {
		return false
	}
	g.buyAmount = b.amount
	for i := 0; i < len(shipDataTable); i++ {
		g.shipCost[i] = g.calculateCost(i)
	}

	return true
}

func (b *BuyAmountButton) Draw(screen *ebiten.Image, g *Game) {
	if b.amount == g.buyAmount {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), disableColor, true)
	}
	b.BaseButton.Draw(screen)
	text.Draw(screen, fmt.Sprint(b.amount), assets.Font20, int(b.x+b.dx), int(b.y+b.height)-2, color.White)
}

func newBuyAmountButton(x, y float64, amount int) *BuyAmountButton {
	w, h := measure(assets.Font20, fmt.Sprint(amount))

	return &BuyAmountButton{
		BaseButton: BaseButton{
			x: x, y: y,
			width:  IMG_SIZE,
			height: h + 4,
		},
		amount: amount,
		dx:     IMG_SIZE/2 - w/2,
	}
}

type ShipButton struct {
	BaseButton
	shipImage *ebiten.Image
	shipID    int
}

func (b *ShipButton) Draw(screen *ebiten.Image, g *Game) {
	isEnabled := g.state.CurMoney >= g.shipCost[b.shipID]
	if !isEnabled {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), disableColor, true)
	}
	b.BaseButton.Draw(screen)
	// Sprite
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	if isEnabled {
		colorm.DrawImage(screen, b.shipImage, defaultScale, op)
	} else {
		colorm.DrawImage(screen, b.shipImage, grayScale, op)
	}
	// level
	lvLabel := fmt.Sprintf("Lv.%v", g.shipLevel[b.shipID])
	advanceLV := font.MeasureString(assets.Font14, lvLabel)
	widthLV := float64(advanceLV >> 6)
	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(widthLV+8), 14, color.White, false)
	text.Draw(screen, lvLabel, assets.Font14, int(b.x+4), int(b.y+12), color.Black)
	leftOffset := IMG_SIZE + 5
	topOffset := 24
	// name
	text.Draw(screen, shipDataTable[b.shipID].Name, assets.Font20, int(b.x)+leftOffset, int(b.y)+topOffset-4, color.White)
	// cost
	text.Draw(screen, "$"+formatBigFloat(g.shipCost[b.shipID]), assets.Font20, int(b.x)+leftOffset, int(b.y)+topOffset*2-4, color.White)
	// dps per ship
	text.Draw(screen, formatBigFloat(g.shipDPS[b.shipID])+" /s", assets.Font20, int(b.x)+leftOffset, int(b.y)+topOffset*3-4, color.White)
	// ship count
	countLabel := fmt.Sprintf("%v", g.state.ShipCounts[b.shipID])
	if countLabel == "0" {
		return
	}
	advanceCount := font.MeasureString(assets.Font14, countLabel)
	widthCount := float32(advanceCount >> 6)
	vector.DrawFilledRect(screen, float32(b.x)+IMG_SIZE-widthCount-8, float32(b.y+IMG_SIZE-14), float32(widthCount+8), 14, color.White, false)
	text.Draw(screen, countLabel, assets.Font14, int(b.x)+IMG_SIZE-int(widthCount)-4, int(b.y)+IMG_SIZE-2, color.Black)
}

func (b *ShipButton) OnPressed(g *Game) bool {
	return g.buyShip(b.shipID)
}

func newShipButton(x, y float64, id int, img *ebiten.Image) *ShipButton {
	return &ShipButton{
		BaseButton: BaseButton{
			x:      x,
			y:      y,
			width:  IMG_SIZE,
			height: IMG_SIZE,
		},
		shipID:    id,
		shipImage: img,
	}
}
