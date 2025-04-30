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
	//TODO: add more
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

	addedSuccessfully := false
	remainingQty := qtyToAdd

	if isStackable {
		for i := range playerInv.Slots {
			slot := &playerInv.Slots[i]
			if slot.ItemID == id && slot.ItemQuantity < maxStack {
				canAdd := maxStack - slot.ItemQuantity
				addAmount := 0
				if remainingQty <= canAdd {
					addAmount = remainingQty
				} else {
					addAmount = canAdd
				}

				slot.ItemQuantity += addAmount
				remainingQty -= addAmount
				addedSuccessfully = true
				fmt.Printf("Added %d to existing stack of %s (Slot %d). New Qty: %d. Remaining to Add: %d\n", addAmount, name, i, slot.ItemQuantity, remainingQty)

				if remainingQty <= 0 {
					return true
				}
			}
		}
	}

	if remainingQty > 0 {
		for i := range playerInv.Slots {
			slot := &playerInv.Slots[i]
			if slot.ItemID == 0 {
				addAmount := 0
				if isStackable {
					if remainingQty <= maxStack {
						addAmount = remainingQty
					} else {
						addAmount = maxStack
					}
				} else {
					if remainingQty > 0 {
						addAmount = 1
					} else {
						addAmount = 0
					}
				}

				if addAmount > 0 {
					slot.ItemID = id
					slot.ItemName = name
					slot.ItemReusable = reusable
					slot.ItemQuantity = addAmount
					remainingQty -= addAmount
					addedSuccessfully = true
					fmt.Printf("Added %d %s to new slot %d. Remaining to Add: %d\n", addAmount, name, i, remainingQty)

					if remainingQty <= 0 {
						return true
					}

					if !isStackable {
						return true
					}

				} else {
					// This case shouldn't be reached if remainingQty > 0, but safety break
					break
				}

			}
		}
	}

	if !addedSuccessfully {
		showMessages("Inventory full! Cannot add "+name+".", 1.0)
		return false
	}

	if remainingQty > 0 {
		showMessages(fmt.Sprintf("Inventory full! Could not add %d %s.", remainingQty, name), 1.0)
		return true
	}

	return true
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
