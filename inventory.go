package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	InventorySize = 5
	pickupRadius  = 20.0
)

var (
	itemSprites map[int]ItemSprite
)

type InventorySlot struct {
	ItemID       int
	ItemName     string
	ItemReusable bool
	ItemQuantity int
}

type Inventory struct {
	Slots       [InventorySize]InventorySlot
	Open        bool
	Cursor      int
	HoveredSlot int
}

type ItemData struct {
	ID             int
	Name           string
	Category       string
	IsPlantable    bool
	IsEdible       bool
	GrowsIntoID    int
	GrowthStages   int
	TimePerStage   float32
	HarvestYield   int
	RequiresTilled bool
	MaxStack       int
	SeedTexture    rl.Texture2D
	PlantTexture   rl.Texture2D
}

type ItemSprite struct {
	Texture rl.Texture2D
	Src     rl.Rectangle
}

var Items = map[int]ItemData{
	1: {ID: 1, Name: "Wheat Seeds", Category: "seed", IsPlantable: true, GrowsIntoID: 10, GrowthStages: 4, TimePerStage: 15.0, HarvestYield: 1, RequiresTilled: true, MaxStack: 99},
	2: {ID: 2, Name: "Strawberry Seeds", Category: "seed", IsPlantable: true, GrowsIntoID: 11, GrowthStages: 4, TimePerStage: 20.0, HarvestYield: 2, RequiresTilled: true, MaxStack: 99},

	3: {ID: 3, Name: "Watering Can", Category: "tool"},
	4: {ID: 4, Name: "Hoe", Category: "tool"},
	5: {ID: 5, Name: "Axe", Category: "tool"},

	10: {ID: 10, Name: "Wheat", Category: "crop", MaxStack: 99},
	11: {ID: 11, Name: "Strawberry", Category: "crop", IsEdible: true, MaxStack: 99},

	20: {ID: 20, Name: "Milk", Category: "animal", IsEdible: true, MaxStack: 99},
	21: {ID: 21, Name: "Butter", Category: "animal", IsEdible: true, MaxStack: 99},
	22: {ID: 22, Name: "Egg", Category: "animal", IsEdible: true, MaxStack: 99},

	30: {ID: 30, Name: "Bread", Category: "baked", IsEdible: true, MaxStack: 99},
	31: {ID: 31, Name: "Strawberry Tart", Category: "baked", IsEdible: true, MaxStack: 99},
	32: {ID: 32, Name: "Strawberry Milk Cake", Category: "baked", IsEdible: true, MaxStack: 99},
	33: {ID: 33, Name: "Burnt Pie", Category: "baked", IsEdible: true, MaxStack: 99},
	34: {ID: 34, Name: "Experimental Jam", Category: "baked", IsEdible: true, MaxStack: 99},

	40: {ID: 40, Name: "Eldermint Leaves", Category: "quest", MaxStack: 99},
	41: {ID: 41, Name: "Cow Flute", Category: "quest", MaxStack: 99},
}

func drawInventory() {
	rl.DrawRectangle(50, 50, screenWidth-100, 100, rl.NewColor(0, 0, 0, 180))
	slotWidth := float32((screenWidth - 100) / InventorySize)
	y := float32(60)
	for i := range InventorySize {
		x := 50 + slotWidth*float32(i)
		box := rl.NewRectangle(x, y, slotWidth-10, 80)
		color := rl.DarkGray

		if i == playerInv.Cursor {
			color = rl.Gray
		}
		if i == playerInv.HoveredSlot {
			color = rl.LightGray
		}
		rl.DrawRectangleRec(box, color)

		slot := &playerInv.Slots[i]
		sprite, ok := itemSprites[slot.ItemID]
		if ok {
			rl.DrawTextureRec(sprite.Texture, sprite.Src, rl.NewVector2(x+10, y+10), rl.White)
		} else {
			name := Items[slot.ItemID].Name
			if name != "" {
				textX := x + (box.Width-float32(rl.MeasureText(name, 16)))/2
				rl.DrawText(name, int32(textX), int32(y+box.Height/2-8), 16, rl.White)
			} else {
				rl.DrawText("-", int32(x+box.Width/2-4), int32(y+box.Height/2-8), 16, rl.White)
			}
		}
	}
}

func useInventoryItem(idx int) {
	slot := &playerInv.Slots[idx]
	if slot.ItemID == 0 {
		showMessages("Nothing to use here!", 0.5)
		return
	}

	itemData, ok := Items[slot.ItemID]
	if !ok {
		showMessages("Unknown item!", 0.5)
		return
	}

	if itemData.Category == "tool" {
		switch itemData.ID {
		case 4: // Hoe
			tryTillSoil()
		case 3: // Watering Can
			tryWaterCrop()
		default:
			showMessages("Used "+slot.ItemName, 0.5)
		}
		return
	}

	if itemData.IsPlantable {
		planted := tryPlantSeed(itemData)
		if planted {
			showMessages("Planted "+slot.ItemName, 0.5)
			slot.ItemQuantity--
			if slot.ItemQuantity == 0 {
				slot.ItemID = 0
				slot.ItemName = ""
				slot.ItemReusable = false
			}
		} else {
			showMessages("Can't Planted "+slot.ItemName, 0.5)
		}
		return
	}

	showMessages("Used "+slot.ItemName, 0.5)
	if !slot.ItemReusable {
		slot.ItemQuantity--
		if slot.ItemQuantity == 0 {
			slot.ItemID = 0
			slot.ItemName = ""
			slot.ItemReusable = false
		} else {
			showMessages(fmt.Sprintf("%s left: %d", slot.ItemName, slot.ItemQuantity), 0.5)
		}
	}

}

func loadItemSprites() {
	plants := rl.LoadTexture("assets/Objects/Basic_Plants.png")
	tools := rl.LoadTexture("assets/Objects/Basic_tools_and_meterials.png")
	milk := rl.LoadTexture("assets/Objects/Simple_Milk_and_grass_item.png")

	itemSprites = map[int]ItemSprite{
		1:  {plants, rl.NewRectangle(0, 0, 16, 16)},  // Wheat Seeds
		2:  {plants, rl.NewRectangle(7, 0, 16, 16)},  // Strawberry Seeds
		3:  {tools, rl.NewRectangle(0, 0, 16, 16)},   // Watering Can
		4:  {tools, rl.NewRectangle(13, 0, 16, 16)},  // Hoe
		5:  {tools, rl.NewRectangle(37, 0, 16, 16)},  // Axe
		10: {plants, rl.NewRectangle(6, 0, 16, 16)},  // Wheat
		11: {plants, rl.NewRectangle(26, 0, 16, 16)}, // Strawberry
		20: {milk, rl.NewRectangle(3, 0, 16, 16)},    // Milk
		21: {rl.LoadTexture("assets/Objects/Butter.png"), rl.NewRectangle(0, 0, 16, 16)},
		22: {rl.LoadTexture("assets/Objects/Egg_item.png"), rl.NewRectangle(0, 0, 16, 16)},
		30: {rl.LoadTexture("assets/Objects/Bread.png"), rl.NewRectangle(0, 0, 16, 16)},
		31: {rl.LoadTexture("assets/Objects/Strawberry_Tart.png"), rl.NewRectangle(0, 0, 16, 16)},
		32: {rl.LoadTexture("assets/Objects/Strawberry_Milk_Cake.png"), rl.NewRectangle(0, 0, 16, 16)},
		33: {rl.LoadTexture("assets/Objects/Burnt_Pie.png"), rl.NewRectangle(0, 0, 16, 16)},
		34: {rl.LoadTexture("assets/Objects/Experimental_Jam.png"), rl.NewRectangle(0, 0, 16, 16)},
		40: {rl.LoadTexture("assets/Objects/Eldermint_Leaves.png"), rl.NewRectangle(0, 0, 16, 16)},
		41: {rl.LoadTexture("assets/Objects/Cow_Flute.png"), rl.NewRectangle(0, 0, 16, 16)},
	}
}
