package main

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type World struct {
	MapWidth          int
	MapHeight         int
	Tiles             []string
	Entities          []Entity
	Crops             []*Crop
	Objects           []*GameObject
	LastUpdate        time.Time
	ShopCounterPos    Vector2
	bakeryCounterPos  Vector2
	Animals           []*Animal
	ChickenHousePos   Vector2
	LastEggCollection time.Time
	EggCooldown       time.Duration
	EggsAvailable     int
}

type Animal struct {
	Type       string
	Position   Vector2
	Texture    rl.Texture2D
	LastMilked time.Time
	HasMilk    bool
}

type Entity struct {
	Type     string
	Position Vector2
	Metadata map[string]interface{}
}

type GameObject struct {
	Type     string
	Position Vector2
	Active   bool
}

type Vector2 struct {
	X, Y int
}

func (w *World) Init(mapWidth, mapHeight int, tiles []string) {
	w.MapWidth = mapWidth
	w.MapHeight = mapHeight
	w.Tiles = tiles

	w.LastUpdate = time.Now()
}

func (w *World) AddCrop(crop *Crop) {
	w.Crops = append(w.Crops, crop)
	w.Entities = append(w.Entities, Entity{
		Type:     "crop",
		Position: Vector2{crop.PosX, crop.PosY},
		Metadata: map[string]interface{}{
			"crop": crop,
		},
	})
}

func (w *World) AddObject(obj *GameObject) {
	w.Objects = append(w.Objects, obj)
	w.Entities = append(w.Entities, Entity{
		Type:     obj.Type,
		Position: obj.Position,
	})
}

func (w *World) InitShopCounter() {
	w.ShopCounterPos = Vector2{X: 15, Y: 4}
	w.AddObject(&GameObject{
		Type:     "ShopCounter",
		Position: w.ShopCounterPos,
		Active:   true,
	})
}

func (w *World) InitCows() {
	w.Animals = []*Animal{
		{
			Type:       "cow",
			Position:   Vector2{X: 5, Y: 3},
			Texture:    cowSprite,
			LastMilked: time.Now().Add(-1 * time.Hour),
		},
	}
}

func (w *World) InitbakeryCounter() {
	w.bakeryCounterPos = Vector2{X: 23, Y: 7}
	w.AddObject(&GameObject{
		Type:     "BakeryCounter",
		Position: w.bakeryCounterPos,
		Active:   true,
	})
}

func (w *World) InitChickenHouse() {
	w.ChickenHousePos = Vector2{X: 23, Y: 4}
	w.AddObject(&GameObject{
		Type:     "ChickenHouse",
		Position: w.ChickenHousePos,
		Active:   true,
	})
	w.EggsAvailable = 3
	w.EggCooldown = 2 * time.Minute
	w.LastEggCollection = time.Now().Add(-w.EggCooldown)
}

func (w *World) Update() {
	now := time.Now()
	delta := now.Sub(w.LastUpdate).Seconds()
	w.LastUpdate = now

	for _, crop := range w.Crops {
		if crop.GrowthStage < 3 {
			crop.GrowthTimer += float32(delta)
			if crop.GrowthTimer >= crop.TimePerStage*60 {
				crop.GrowthStage++
				crop.GrowthTimer = 0
			}
		}
	}
}

func (w *World) GetTileType(x, y int) string {
	if x < 0 || x >= w.MapWidth || y < 0 || y >= w.MapHeight {
		return ""
	}
	return w.Tiles[y*w.MapWidth+x]
}

func (w *World) FindEntitiesAt(pos Vector2, radius int) []Entity {
	var found []Entity
	for _, e := range w.Entities {
		dx := math.Abs(float64(e.Position.X - pos.X))
		dy := math.Abs(float64(e.Position.Y - pos.Y))
		if dx <= float64(radius) && dy <= float64(radius) {
			found = append(found, e)
		}
	}
	return found
}

func (w *World) GetCropAt(tileX, tileY int) *Crop {
	for _, crop := range w.Crops {
		if crop.PosX == tileX && crop.PosY == tileY {
			return crop
		}
	}
	return nil
}

func (w *World) RemoveCrop(cropToRemove *Crop) {
	found := false
	newCrops := w.Crops[:0]
	for _, c := range w.Crops {
		if c != cropToRemove {
			newCrops = append(newCrops, c)
		} else {
			found = true
		}
	}
	w.Crops = newCrops

	newEntities := w.Entities[:0]
	for _, e := range w.Entities {
		isCropToRemove := false
		if e.Type == "crop" {
			if cropRef, ok := e.Metadata["crop"].(*Crop); ok {
				if cropRef == cropToRemove {
					isCropToRemove = true
				}
			}
		}
		if !isCropToRemove {
			newEntities = append(newEntities, e)
		}
	}
	w.Entities = newEntities

	if found {
		fmt.Printf("Removed crop ID %d at %d,%d. Remaining crops: %d\n", cropToRemove.ID, cropToRemove.PosX, cropToRemove.PosY, len(w.Crops))
	} else {
		fmt.Printf("Warning: Tried to remove crop at %d,%d but it wasn't found.\n", cropToRemove.PosX, cropToRemove.PosY)
	}
}
