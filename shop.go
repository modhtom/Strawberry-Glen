package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	playerGold                  = 100 // starting money
	shopOpen                    = false
	shopCursor                  = 0
	shopMode                    = 0 // 0 = Buy, 1 = Sell
	shopRestockTimer    float32 = 0
	shopRestockInterval float32 = 300.0 // 5 minutes
	initialShopStock            = make([]int, len(shopItems))
)

type ShopItem struct {
	ID        int
	PriceBuy  int
	PriceSell int
	Stock     int
}

var shopItems = []ShopItem{
	{ID: 1, PriceBuy: 5, PriceSell: 2, Stock: 10},   // Wheat Seeds
	{ID: 10, PriceBuy: 20, PriceSell: 8, Stock: 5},  // Wheat
	{ID: 11, PriceBuy: 25, PriceSell: 10, Stock: 5}, // Strawberry
	// add more
}

func findInventorySlot(id int) (*InventorySlot, bool) {
	for i := range playerInv.Slots {
		if playerInv.Slots[i].ItemID == id {
			return &playerInv.Slots[i], true
		}
	}
	return nil, false
}

func addToInventory(id int, name string, reusable bool, qty int) bool {
	for i := range playerInv.Slots {
		if playerInv.Slots[i].ItemID == 0 {
			playerInv.Slots[i] = InventorySlot{
				ItemID:       id,
				ItemName:     name,
				ItemReusable: reusable,
				ItemQuantity: qty,
			}
			return true
		}
	}
	showMessages("There is no room left in your Inventory.")
	return false
}

func drawShop() {
	w, h := float32(400), float32(300)
	x, y := float32(screenWidth)/2-w/2, float32(screenHeight)/2-h/2
	rl.DrawRectangleRec(rl.NewRectangle(x, y, w, h), rl.NewColor(0, 0, 0, 200))

	rl.DrawText(fmt.Sprintf("SHOP — Gold: %d", playerGold), int32(x+20), int32(y+20), 20, rl.White)
	modeText := []string{"[Buy]", "[Sell]"}[shopMode]
	rl.DrawText(modeText, int32(x+300), int32(y+20), 20, rl.Yellow)
	for i, it := range shopItems {
		ty := int32(int(y) + 60 + 24*i)
		name := Items[it.ID].Name
		price := it.PriceBuy
		if shopMode == 1 {
			price = it.PriceSell
		}
		stock := it.Stock
		if shopMode == 1 {
			if slot, ok := findInventorySlot(it.ID); ok {
				stock = slot.ItemQuantity
			} else {
				stock = 0
			}
		}
		color := rl.White
		if i == shopCursor {
			color = rl.Yellow
		}
		rl.DrawText(fmt.Sprintf("%-15s %3dG  Stock:%2d", name, price, stock),
			int32(x+20), ty, 20, color)
	}
	rl.DrawText("left arrow/right arrow Switch Buy/Sell.", int32(x+20), int32(y+h-40), 16, rl.LightGray)
	rl.DrawText("Use 'B' to close the shop.", int32(x+20), int32(y+h-20), 16, rl.LightGray)
}
