package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Crop struct {
	ID            int
	CropTypeID    int
	GrowthStage   int
	MaxGrowth     int
	GrowthTimer   float32
	TimePerStage  float32
	PosX, PosY    int
	PlantedAt     time.Time
	NeedsWatering bool
	IsWatered     bool
	WaterTimer    float32
}

func (c *Crop) IsHarvestable() bool {
	return c.GrowthStage >= c.MaxGrowth
}

func (c *Crop) GetSpriteRect() rl.Rectangle {
	offsetX := float32(16 * c.GrowthStage)
	return rl.NewRectangle(offsetX, 0, 16, 16)
}

type Object struct {
	Type       string
	PosX, PosY int
	Metadata   map[string]interface{}
}

func getFacingTilePos() (int, int) {
	tileX := int((playerDest.X + playerDest.Width/2) / tileDest.Width)
	tileY := int((playerDest.Y + playerDest.Height/2) / tileDest.Height)

	switch playerDir {
	case 0: // Down
		tileY++
	case 1: // Up
		tileY--
	case 2: // Left
		tileX--
	case 3: // Right
		tileX++
	}
	return tileX, tileY
}

func isTileValid(tileX, tileY int) bool {
	return tileX >= 0 && tileX < mapW && tileY >= 0 && tileY < mapH
}

func getTileIndex(tileX, tileY int) int {
	if !isTileValid(tileX, tileY) {
		return -1
	}
	return tileY*mapW + tileX
}

func tryTillSoil() {
	tileX, tileY := getFacingTilePos()
	index := getTileIndex(tileX, tileY)

	if index == -1 {
		showMessages("Can't till outside the map!", 0.5)
		return
	}

	if srcMap[index] == "t" {
		showMessages("Already tilled!", 0.5)
		return
	}
	if srcMap[index] != "g" {
		showMessages("Cannot till this ground!", 0.5)
		return
	}

	if world.GetCropAt(tileX, tileY) != nil {
		showMessages("Something is planted here!", 0.5)
		return
	}

	srcMap[index] = "t"
	showMessages("Tilled the soil!", 0.5)
}

func tryPlantSeed(seedData ItemData) bool {
	tileX, tileY := getFacingTilePos()
	index := getTileIndex(tileX, tileY)

	if index == -1 {
		showMessages("Can't plant outside the map!", 0.5)
		return false
	}

	if seedData.RequiresTilled && srcMap[index] != "t" {
		showMessages("Needs tilled soil!", 0.5)
		return false
	}

	if srcMap[index] == "w" || srcMap[index] == "h" {
		showMessages("Cannot plant here!", 0.5)
		return false
	}

	if world.GetCropAt(tileX, tileY) != nil {
		showMessages("Something is already growing here!", 0.5)
		return false
	}

	newCrop := &Crop{
		ID:            seedData.ID,
		CropTypeID:    seedData.GrowsIntoID,
		GrowthStage:   0,
		MaxGrowth:     seedData.GrowthStages - 1,
		GrowthTimer:   0,
		TimePerStage:  seedData.TimePerStage,
		PosX:          tileX,
		PosY:          tileY,
		PlantedAt:     time.Now(),
		NeedsWatering: true,
		IsWatered:     false,
		WaterTimer:    0,
	}
	world.AddCrop(newCrop)

	return true
}

func tryWaterCrop() {
	tileX, tileY := getFacingTilePos()
	crop := world.GetCropAt(tileX, tileY)

	if crop == nil {
		showMessages("Nothing to water here.", 0.5)
		return
	}

	if crop.IsWatered {
		showMessages("Already watered!", 0.5)
		return
	}

	if !crop.NeedsWatering {
		showMessages("This doesn't need watering.", 0.5)
		return
	}

	crop.IsWatered = true
	crop.WaterTimer = 18 * 60
	showMessages("Watered the "+Items[crop.ID].Name+"!", 0.5)
}

func tryInteract() {
	tileX, tileY := getFacingTilePos()

	crop := world.GetCropAt(tileX, tileY)
	if crop != nil && crop.IsHarvestable() {
		harvestCrop(crop)
		return
	}

	entities := world.FindEntitiesAt(Vector2{X: tileX, Y: tileY}, 0)
	for _, e := range entities {
		switch e.Type {
		case "Door":
			showMessages("It's a door.", 0.5)

		case "NPC":
			showMessages("Talking to NPC...", 0.5)

		default:

			index := getTileIndex(tileX, tileY)
			if index != -1 {
				tileType := srcMap[index]
				switch tileType {
				case "s":
					showMessages("This is the shop counter. Press 'B' to shop.", 1.0)
				default:
					showMessages("Nothing interesting here.", 0.5)
				}
			}
		}
		return
	}
	showMessages("Nothing to interact with.", 0.5)
}

func harvestCrop(crop *Crop) {
	cropData, ok := Items[crop.ID]
	if !ok {
		fmt.Printf("Error: Cannot find item data for seed ID %d during harvest\n", crop.ID)
		return
	}

	harvestID := cropData.GrowsIntoID
	harvestAmount := cropData.HarvestYield

	harvestItemData, ok := Items[harvestID]
	if !ok {
		fmt.Printf("Error: Cannot find item data for harvest ID %d\n", harvestID)
		return
	}
	successfullyHarvestedCount := 0
	couldNotAddCount := 0

	for i := 0; i < harvestAmount; i++ {
		if addToInventory(harvestID, harvestItemData.Name, harvestItemData.IsEdible, 1) {
			successfullyHarvestedCount++
		} else {
			couldNotAddCount++
			break
		}
	}

	if successfullyHarvestedCount > 0 {
		if couldNotAddCount > 0 {
			showMessages(fmt.Sprintf("Harvested %d %s. Inventory full, %d left behind.", successfullyHarvestedCount, harvestItemData.Name, couldNotAddCount), 1.5)
		} else {
			showMessages(fmt.Sprintf("Harvested %d %s!", successfullyHarvestedCount, harvestItemData.Name), 1.0)
		}
		world.RemoveCrop(crop)

		index := getTileIndex(crop.PosX, crop.PosY)
		if index != -1 {
			srcMap[index] = "t"
		}
	} else {
		showMessages("Inventory full! Cannot harvest "+harvestItemData.Name+".", 1.0)
	}
}
