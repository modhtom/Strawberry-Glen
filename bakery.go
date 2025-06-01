package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	bakeryOpen   = false
	bakeryCursor = 0
	bakeryMode   = 0 // 0 = Buy, 1 = Sell

	bakingMode          = false
	selectedIngredients []int
	recipeCursor        = 0
)

type bakeryItem struct {
	ID        int
	PriceBuy  int
	PriceSell int
	Stock     int
}

var bakeryItems = []bakeryItem{
	{ID: 30, PriceBuy: 25, PriceSell: 12, Stock: 15}, // Bread (needs 1 Wheat)
	{ID: 31, PriceBuy: 40, PriceSell: 18, Stock: 15}, // Strawberry Tart

	{ID: 32, PriceBuy: 85, PriceSell: 40, Stock: 5}, // Strawberry Milk Cake
	{ID: 33, PriceBuy: 10, PriceSell: 5, Stock: 5},  // Burnt Pie (punishment)
	{ID: 34, PriceBuy: 45, PriceSell: 17, Stock: 5}, // Experimental Jam

	{ID: 40, PriceBuy: 200, PriceSell: 80, Stock: 1},  // Eldermint
	{ID: 41, PriceBuy: 300, PriceSell: 150, Stock: 1}, // Cow Flute
}

type Recipe struct {
	Ingredients    []int
	OutputID       int
	OutputQuantity int
}

var bakeryRecipes = []Recipe{
	{Ingredients: []int{10, 11}, OutputID: 31, OutputQuantity: 1},     // Wheat + Strawberry = Tart
	{Ingredients: []int{10, 20, 11}, OutputID: 32, OutputQuantity: 1}, // Wheat + Milk + Strawberry = Cake
	{Ingredients: []int{10}, OutputID: 30, OutputQuantity: 1},         // Wheat = Bread
	{Ingredients: []int{10, 11, 34}, OutputID: 33, OutputQuantity: 1}, // Failed recipe = Burnt Pie
}

func drawBakery() {
	w, h := float32(400), float32(300)
	x, y := float32(screenWidth)/2-w/2, float32(screenHeight)/2-h/2
	rl.DrawRectangleRec(rl.NewRectangle(x, y, w, h), rl.NewColor(0, 0, 0, 200))

	if bakingMode {
		drawBakingInterface(x, y, w, h)
	} else {
		rl.DrawText(fmt.Sprintf("bakery — Gold: %d", playerGold), int32(x+20), int32(y+20), 20, rl.White)
		modeText := []string{"[Buy]", "[Sell]"}[bakeryMode]
		rl.DrawText(modeText, int32(x+300), int32(y+20), 20, rl.Yellow)
		for i, it := range bakeryItems {
			ty := int32(int(y) + 60 + 24*i)
			name := Items[it.ID].Name
			price := it.PriceBuy
			if bakeryMode == 1 {
				price = it.PriceSell
			}
			stock := it.Stock
			if bakeryMode == 1 {
				stock = 0
				for _, slot := range playerInv.Slots {
					if slot.ItemID == it.ID {
						stock += slot.ItemQuantity
					}
				}
			}
			color := rl.White
			if i == bakeryCursor {
				color = rl.Yellow
			}
			rl.DrawText(fmt.Sprintf("%-15s %3dG  Stock:%2d", name, price, stock),
				int32(x+20), ty, 20, color)
		}
	}
	rl.DrawText("Press TAB to switch modes", int32(x+20), int32(y+h-60), 16, rl.LightGray)
	rl.DrawText("left arrow/right arrow Switch Buy/Sell.", int32(x+20), int32(y+h-40), 16, rl.LightGray)
	rl.DrawText("Use 'B' to close the bakery.", int32(x+20), int32(y+h-20), 16, rl.LightGray)
}

func drawBakingInterface(x, y, w, h float32) {
	rl.DrawText("BAKING - Select Recipe", int32(x+20), int32(y+20), 20, rl.White)
	for i, recipe := range bakeryRecipes {
		color := rl.White
		if i == recipeCursor {
			color = rl.Yellow
		}

		var ingredients string
		for _, ing := range recipe.Ingredients {
			ingredients += Items[ing].Name + " + "
		}
		ingredients = ingredients[:len(ingredients)-3]

		rl.DrawText(fmt.Sprintf("%s ➜ %s", ingredients, Items[recipe.OutputID].Name),
			int32(x+20), int32(int(y)+60+24*i), 20, color)
	}

	startY := y + h - 140
	rl.DrawText("Selected Ingredients:", int32(x+20), int32(startY), 20, rl.White)
	for i, itemID := range selectedIngredients {
		rl.DrawText(fmt.Sprintf("%d. %s", i+1, Items[itemID].Name),
			int32(x+20), int32(int(startY)+30+24*i), 20, rl.LightGray)
	}

	rl.DrawText(fmt.Sprintf("Press 1-%v to add ingredients from inventory slots", InventorySize),
		int32(x+20), int32(y+h-80), 16, rl.LightGray)

	rl.DrawText(fmt.Sprintf("Press ALT + 1-%v to remove ingredients from inventory slots", InventorySize),
		int32(x+20), int32(y+h-100), 16, rl.LightGray)
}

func tryBake() bool {
	if len(selectedIngredients) == 0 {
		showMessages("Select ingredients first!", 1.0)
		return false
	}
	selectedRecipe := bakeryRecipes[recipeCursor]

	if !compareIngredients(selectedRecipe.Ingredients, selectedIngredients) {
		showMessages("Ingredients don't match recipe!", 1.0)
		return false
	}

	if !hasInventorySpace(selectedRecipe.OutputID) {
		showMessages("No space for baked goods!", 1.0)
		return false
	}

	for _, ingID := range selectedRecipe.Ingredients {
		if slot, ok := findInventorySlot(ingID); ok {
			slot.ItemQuantity--
			if slot.ItemQuantity <= 0 {
				slot.ItemID = 0
				slot.ItemName = ""
			}
		}
	}

	addToInventory(selectedRecipe.OutputID, Items[selectedRecipe.OutputID].Name, false, selectedRecipe.OutputQuantity)
	showMessages(fmt.Sprintf("Baked %s!", Items[selectedRecipe.OutputID].Name), 2.0)
	selectedIngredients = []int{}
	return true
}

func compareIngredients(recipe, selected []int) bool {
	if len(recipe) != len(selected) {
		return false
	}
	counts := make(map[int]int)
	for _, id := range selected {
		counts[id]++
	}

	for _, id := range recipe {
		if counts[id] <= 0 {
			return false
		}
		counts[id]--
	}
	return true
}

func hasInventorySpace(itemID int) bool {
	itemData := Items[itemID]
	if slot, ok := findInventorySlot(itemID); ok && itemData.MaxStack > slot.ItemQuantity {
		return true
	}
	for _, slot := range playerInv.Slots {
		if slot.ItemID == 0 {
			return true
		}
	}
	return false
}

func handleBakingInput() {
	if rl.IsKeyPressed(rl.KeyUp) {
		recipeCursor = (recipeCursor - 1 + len(bakeryRecipes)) % len(bakeryRecipes)
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		recipeCursor = (recipeCursor + 1) % len(bakeryRecipes)
	}
	if rl.IsKeyPressed(rl.KeyEnter) {
		tryBake()
	}

	// Add ingredients with number keys 1-InventorySize
	for i := range InventorySize {
		if rl.IsKeyPressed(int32(rl.KeyOne) + int32(i)) {
			if i < len(playerInv.Slots) {
				itemID := playerInv.Slots[i].ItemID
				if itemID != 0 && Items[itemID].Category != "tool" {
					selectedIngredients = append(selectedIngredients, itemID)
				}
			}
		}
	}

	// Remove ingredients with ALT + number keys 1-InventorySize
	for i := 0; i < len(selectedIngredients); i++ {
		keyIndex := int32(rl.KeyOne) + int32(i)
		if rl.IsKeyDown(rl.KeyLeftAlt) && rl.IsKeyPressed(keyIndex) {
			itemID := selectedIngredients[i]
			itemData := Items[itemID]
			addToInventory(itemID, itemData.Name, false, 1)
			selectedIngredients = append(selectedIngredients[:i], selectedIngredients[i+1:]...)
			break
		}
	}

}
