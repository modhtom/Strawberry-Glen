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
	Unlocked  bool
}

var shopItems = []ShopItem{
	{ID: 1, PriceBuy: 5, PriceSell: 2, Stock: 10, Unlocked: true},
	{ID: 2, PriceBuy: 8, PriceSell: 3, Stock: 5, Unlocked: false},
	{ID: 10, PriceBuy: 20, PriceSell: 8, Stock: 5, Unlocked: true},
	{ID: 40, PriceBuy: 100, PriceSell: 40, Stock: 3, Unlocked: false},
}

func getAvailableShopItems() []*ShopItem {
	var available []*ShopItem
	for i := range shopItems {
		if shopItems[i].Unlocked || progression.UnlockedSeeds[shopItems[i].ID] {
			available = append(available, &shopItems[i])
		}
	}
	return available
}

func findInventorySlot(id int) (*InventorySlot, bool) {
	for i := range playerInv.Slots {
		if playerInv.Slots[i].ItemID == id {
			return &playerInv.Slots[i], true
		}
	}
	return nil, false
}
func addToInventory(id int, name string, reusable bool, qtyToAdd int) bool {
	itemData, ok := Items[id]
	if !ok {
		fmt.Printf("Warning: Trying to add unknown item ID %d\n", id)
		return false
	}

	maxStack := itemData.MaxStack
	isStackable := !reusable && maxStack > 1
	remainingQty := qtyToAdd
	added := false
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	if isStackable {
		for i := range playerInv.Slots {
			if remainingQty <= 0 {
				break
			}

			slot := &playerInv.Slots[i]
			if slot.ItemID == id && slot.ItemQuantity < maxStack {
				availableSpace := maxStack - slot.ItemQuantity
				toAdd := min(remainingQty, availableSpace)

				slot.ItemQuantity += toAdd
				remainingQty -= toAdd
				added = true
			}
		}
	}

	if remainingQty > 0 {
		if isStackable {
			for i := range playerInv.Slots {
				if remainingQty <= 0 {
					break
				}

				slot := &playerInv.Slots[i]
				if slot.ItemID == 0 {
					toAdd := min(remainingQty, maxStack)

					slot.ItemID = id
					slot.ItemName = name
					slot.ItemReusable = reusable
					slot.ItemQuantity = toAdd
					remainingQty -= toAdd
					added = true
				}
			}
		} else {
			for i := 0; i < remainingQty; i++ {
				foundSlot := false

				for j := range playerInv.Slots {
					slot := &playerInv.Slots[j]
					if slot.ItemID == 0 {
						slot.ItemID = id
						slot.ItemName = name
						slot.ItemReusable = reusable
						slot.ItemQuantity = 1
						foundSlot = true
						added = true
						break
					}
				}

				if !foundSlot {
					remainingQty = i
					break
				}
			}
		}
	}

	if remainingQty > 0 {
		showMessages(fmt.Sprintf("Inventory full! Couldn't add %d %s", remainingQty, name), 1.0)
	}

	return added || qtyToAdd > remainingQty
}

func drawShop() {
	w, h := float32(400), float32(300)
	x, y := float32(screenWidth)/2-w/2, float32(screenHeight)/2-h/2
	rl.DrawRectangleRec(rl.NewRectangle(x, y, w, h), rl.NewColor(0, 0, 0, 200))

	rl.DrawText(fmt.Sprintf("SHOP — Gold: %d", playerGold), int32(x+20), int32(y+20), 20, rl.White)
	modeText := []string{"[Buy]", "[Sell]"}[shopMode]
	rl.DrawText(modeText, int32(x+300), int32(y+20), 20, rl.Yellow)

	availableItems := getAvailableShopItems()
	for i, it := range availableItems {
		ty := int32(int(y) + 60 + 24*i)
		name := Items[it.ID].Name
		price := it.PriceBuy
		if shopMode == 1 {
			price = it.PriceSell
		}
		color := rl.White
		if i == shopCursor {
			color = rl.Yellow
		}
		rl.DrawText(fmt.Sprintf("%-15s %3dG  Stock:%2d", name, price, it.Stock),
			int32(x+20), ty, 20, color)
	}

	rl.DrawText("left arrow/right arrow Switch Buy/Sell.", int32(x+20), int32(y+h-40), 16, rl.LightGray)
	rl.DrawText("Use 'B' to close the shop.", int32(x+20), int32(y+h-20), 16, rl.LightGray)
}
