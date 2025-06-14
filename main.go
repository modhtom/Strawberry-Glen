package main

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	StateMainMenu = iota
	StatePlaying
	StateCredits
)

var (
	screenHeight  = 480
	screenWidth   = 1000
	windowResized = false
	prevScreenW   = screenWidth
	prevScreenH   = screenHeight

	running = true

	gameState      = StateMainMenu
	saveFileExists = false
	menuCursor     = 0
	menuItems      = []string{"Play", "Continue", "Credits", "Exit"}

	tex      rl.Texture2D
	bkgColor = rl.NewColor(147, 211, 196, 255)

	ShopCounterPos Vector2

	wheatGrowthSprite      rl.Texture2D
	strawberryGrowthSprite rl.Texture2D

	cowSprite          rl.Texture2D
	chickenHouseSprite rl.Texture2D

	grassSprite  rl.Texture2D
	houseSprite  rl.Texture2D
	shopSprite   rl.Texture2D
	bakerySprite rl.Texture2D
	hillSprite   rl.Texture2D
	waterSprite  rl.Texture2D
	fenceSprite  rl.Texture2D
	tilledSprite rl.Texture2D
	doorSprite   rl.Texture2D
	wasOnMud     bool
	tileDest     rl.Rectangle
	tileSrc      rl.Rectangle
	tileMap      []int
	srcMap       []string
	mapW, mapH   int

	playerSprite                                  rl.Texture2D
	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerSpeed                                   float32 = 1.5
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int
	playerInv                                     Inventory

	frameCount int

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D

	paused = false

	showMessage  bool
	messageText  string
	messageTimer float32

	showKeyBindings bool
	gameSave        bool

	timeOfDay          float32
	dayDuration        float32 = 24
	transitionDuration float32 = 0.01
	numberOfDays       int64   = 1

	world World
)

func addGold(amount int) {
	playerGold += amount
	progression.TotalGoldEarned += amount
}

func input() {
	if !paused && !shopOpen && !bakeryOpen {
		for key := rl.KeyOne; key <= rl.KeyNine; key++ {
			if rl.IsKeyPressed(int32(key)) {
				slotIndex := int(key - rl.KeyOne)
				if slotIndex < len(playerInv.Slots) {
					useInventoryItem(slotIndex)
					playerInv.Cursor = slotIndex
				}
			}
		}
	}

	if rl.IsKeyPressed(rl.KeyEscape) && !shopOpen && !bakeryOpen {
		if gameState == StatePlaying {
			paused = !paused
		} else {
			gameState = StateMainMenu
		}
	}

	if paused {
		return
	}
	if rl.IsKeyPressed(rl.KeyF11) {
		rl.ToggleFullscreen()
	}

	if rl.IsWindowResized() {
		screenWidth = int(rl.GetScreenWidth())
		screenHeight = int(rl.GetScreenHeight())
		windowResized = true
	}

	if playerInv.Open {
		mp := rl.GetMousePosition()
		slotWidth := float32((screenWidth - 100) / InventorySize)
		baseX, baseY := float32(50), float32(60)

		playerInv.HoveredSlot = -1

		for i := range InventorySize {
			x := baseX + slotWidth*float32(i)
			r := rl.NewRectangle(x, baseY, slotWidth-10, 80)
			if rl.CheckCollisionPointRec(mp, r) {
				playerInv.HoveredSlot = i
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) || rl.IsKeyPressed(rl.KeyEnter) {
					useInventoryItem(i)
				}
				break
			}
		}
		if rl.IsKeyPressed(rl.KeyI) {
			playerInv.Open = !playerInv.Open
			playerInv.HoveredSlot = -1
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			playerInv.Cursor--
			if playerInv.Cursor < 0 {
				playerInv.Cursor = InventorySize - 1
			}
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			playerInv.Cursor = (playerInv.Cursor + 1) % InventorySize
		}
		return
	}

	if rl.IsKeyPressed(rl.KeyH) {
		tryTillSoil()
	}

	if rl.IsKeyPressed(rl.KeyB) {
		if isPlayerNearShop() {
			shopOpen = !shopOpen
			shopCursor = 0
			shopMode = 0
		} else if isPlayerNearBakery() {
			bakeryOpen = !bakeryOpen
			bakeryCursor = 0
			bakeryMode = 0
		} else if isPlayerNearChickenHouse() {
			if time.Since(world.LastEggCollection) < world.EggCooldown {
				remaining := world.EggCooldown - time.Since(world.LastEggCollection)
				showMessages(fmt.Sprintf("Chickens resting! Come back in %0.1f minutes.",
					remaining.Minutes()), 2.0)
				return
			}

			if world.EggsAvailable <= 0 {
				showMessages("No eggs available - check back later!", 2.0)
				return
			}

			if added := addToInventory(22, "Egg", false, 1); added {
				world.EggsAvailable--
				showMessages(fmt.Sprintf("Collected egg! %d remaining.", world.EggsAvailable), 1.5)
				world.LastEggCollection = time.Now()
			} else {
				showMessages("Inventory full - make space for eggs!", 1.5)
			}
		} else {
			showMessages("You need to be at a counter!", 0.4)
		}
	}

	if shopOpen {
		if rl.IsKeyPressed(rl.KeyUp) {
			shopCursor = (shopCursor - 1 + len(shopItems)) % len(shopItems)
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			shopCursor = (shopCursor + 1) % len(shopItems)
		}
		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight) {
			shopMode ^= 1
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			item := &shopItems[shopCursor]
			slot, found := findInventorySlot(item.ID)
			switch shopMode {
			case 0: // BUY
				if playerGold >= item.PriceBuy && item.Stock > 0 {
					playerGold -= item.PriceBuy
					item.Stock--
					if found {
						slot.ItemQuantity++
					} else {
						if ok := addToInventory(item.ID, Items[item.ID].Name, false, 1); !ok {
							playerGold += item.PriceBuy
							item.Stock++
						}
					}
				}
			case 1: // SELL
				if found && slot.ItemQuantity > 0 {
					addGold(item.PriceSell)
					playerGold += item.PriceSell
					item.Stock++
					slot.ItemQuantity--
					if slot.ItemQuantity == 0 {
						slot.ItemID = 0
						slot.ItemName = ""
					}
					progression.ItemsSold++
				}
			}
		}
		return
	}

	if bakeryOpen {
		if rl.IsKeyPressed(rl.KeyTab) {
			bakingMode = !bakingMode
			selectedIngredients = []int{}
			recipeCursor = 0
		}
		if bakingMode {
			handleBakingInput()
			return
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			bakeryCursor = (bakeryCursor - 1 + len(bakeryItems)) % len(bakeryItems)
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			bakeryCursor = (bakeryCursor + 1) % len(bakeryItems)
		}
		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight) {
			bakeryMode ^= 1
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			item := &bakeryItems[bakeryCursor]
			slot, found := findInventorySlot(item.ID)
			switch bakeryMode {
			case 0: // BUY
				if playerGold >= item.PriceBuy && item.Stock > 0 {
					playerGold -= item.PriceBuy
					item.Stock--
					if found {
						slot.ItemQuantity++
					} else {
						if ok := addToInventory(item.ID, Items[item.ID].Name, false, 1); !ok {
							playerGold += item.PriceBuy
							item.Stock++
						}
					}
				}
			case 1: // SELL
				if found && slot.ItemQuantity > 0 {
					playerGold += item.PriceSell
					item.Stock++
					slot.ItemQuantity--
					if slot.ItemQuantity == 0 {
						slot.ItemID = 0
						slot.ItemName = ""
					}
				}
			}
		}
		return
	}

	world.Update()

	if rl.IsKeyPressed(rl.KeyE) {
		tryInteract()
	}

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
	}

	if rl.IsKeyPressed(rl.KeyI) {
		playerInv.Open = !playerInv.Open
		playerInv.HoveredSlot = -1
	}

	if rl.IsKeyPressed(rl.KeyZ) {
		cam.Zoom += 0.5
	}
	if rl.IsKeyPressed(rl.KeyX) {
		if cam.Zoom == 0 {
			cam.Zoom = 1.0
		} else {
			cam.Zoom -= 0.5
		}
	}
}

func isPlayerNearChickenHouse() bool {
	playerTileX := int((playerDest.X + playerDest.Width/2) / tileDest.Width)
	playerTileY := int((playerDest.Y + playerDest.Height/2) / tileDest.Height)

	return math.Abs(float64(playerTileX-world.ChickenHousePos.X)) <= 1 &&
		math.Abs(float64(playerTileY-world.ChickenHousePos.Y)) <= 1
}

func isPlayerNearCow() bool {
	playerTileX := int((playerDest.X + playerDest.Width/2) / tileDest.Width)
	playerTileY := int((playerDest.Y + playerDest.Height/2) / tileDest.Height)

	for _, animal := range world.Animals {
		if animal.Type == "cow" {
			if math.Abs(float64(playerTileX-animal.Position.X)) <= 1 &&
				math.Abs(float64(playerTileY-animal.Position.Y)) <= 1 {
				return true
			}
		}
	}
	return false
}

func tryMilkCow() {
	if !isPlayerNearCow() {
		showMessages("No cows nearby!", 0.5)
		return
	}

	for _, cow := range world.Animals {
		if cow.Type == "cow" && time.Since(cow.LastMilked) >= 2*time.Minute {
			if added := addToInventory(20, "Milk", false, 1); added {
				cow.LastMilked = time.Now()
				showMessages("Milked cow! Got fresh milk.", 0.3)
				return
			} else {
				showMessages("Inventory full - make space!", 0.5)
				return
			}
		} else if cow.Type == "cow" {
			remaining := 2*time.Minute - time.Since(cow.LastMilked)
			showMessages(fmt.Sprintf("Cow needs rest! Come back in %0.1f minutes.",
				remaining.Minutes()), 1.5)
			return
		}
	}
}

func isPlayerNearShop() bool {
	playerTileX := int((playerDest.X + playerDest.Width/2) / tileDest.Width)
	playerTileY := int((playerDest.Y + playerDest.Height/2) / tileDest.Height)

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			checkX := playerTileX + dx
			checkY := playerTileY + dy
			if checkX == world.ShopCounterPos.X && checkY == world.ShopCounterPos.Y {
				return true
			}
			if world.ShopCounterPos.X > 0 && checkX == world.ShopCounterPos.X+1 && checkY == world.ShopCounterPos.Y {
				return true
			}
		}
	}
	return false
}

func isPlayerNearBakery() bool {
	playerTileX := int((playerDest.X + playerDest.Width/2) / tileDest.Width)
	playerTileY := int((playerDest.Y + playerDest.Height/2) / tileDest.Height)

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			checkX := playerTileX + dx
			checkY := playerTileY + dy
			if checkX == world.bakeryCounterPos.X && checkY == world.bakeryCounterPos.Y {
				return true
			}
			if world.bakeryCounterPos.X > 0 && checkX == world.bakeryCounterPos.X+1 && checkY == world.bakeryCounterPos.Y {
				return true
			}
		}
	}
	return false
}

func update() {
	running = !rl.WindowShouldClose()

	if windowResized {
		cam.Offset = rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2))
		cam.Zoom = rl.Clamp(cam.Zoom, 0.5, 3.0)
		windowResized = false
		prevScreenW = screenWidth
		prevScreenH = screenHeight
	}

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	var dx, dy float32 = 0, 0

	if playerMoving {
		if playerUp {
			dy -= playerSpeed
		}
		if playerDown {
			dy += playerSpeed
		}
		if playerLeft {
			dx -= playerSpeed
		}
		if playerRight {
			dx += playerSpeed
		}

		proposedX := playerDest.X + dx
		proposedY := playerDest.Y + dy

		if canMove(proposedX, proposedY) {
			playerDest.X = proposedX
			playerDest.Y = proposedY
		}

		if frameCount%8 == 1 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		playerFrame++
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}
	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	currentTileTypes := getCurrentTileTypes(playerDest.X, playerDest.Y)
	onMud := false

	for _, t := range currentTileTypes {
		if t == "t" {
			onMud = true
			break
		}
	}
	if onMud && !wasOnMud {
		showMessages("You slip in the mud!", 0.4)
	}
	wasOnMud = onMud

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	shopRestockTimer += rl.GetFrameTime()
	if shopRestockTimer >= shopRestockInterval {
		shopRestockTimer -= shopRestockInterval
		for i := range shopItems {
			shopItems[i].Stock = initialShopStock[i]
		}
		showMessages("Shop has been restocked!", 1.0)
	}

	timeScale := rl.GetFrameTime()
	timeOfDay += transitionDuration * 60 * timeScale
	if timeOfDay >= dayDuration {
		timeOfDay -= dayDuration
		numberOfDays++
	}

	if timeOfDay >= dayDuration {
		progression.DaysPlayed++
	}

	applyUpgrades()
	checkProgression()

	cam.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)),
		float32(playerDest.Y-(playerDest.Height/2)))

	playerMoving = false
	playerUp, playerDown, playerLeft, playerRight = false, false, false, false
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()

	drawMessage()

	//hour := int(timeOfDay * 24 / dayDuration)
	//minute := int(math.Mod(float64(timeOfDay*24*60/dayDuration), 60))
	//ampm := "AM"
	//displayHour := hour
	//if hour == 0 {
	//displayHour = 12
	//} else if hour == 12 {
	//ampm = "PM"
	//} else if hour > 12 {
	//displayHour = hour - 12
	//ampm = "PM"
	//}
	// rl.DrawText(fmt.Sprintf("Day: %d", numberOfDays),  10, 10, 20, rl.White)
	//rl.DrawText(fmt.Sprintf("Time: %02d:%02d %s", displayHour, minute, ampm), 10, 35, 20, rl.White)

	//if hour >= 18 || hour <= 6 {
	//	alpha := uint8(math.Abs(float64(hour-18)) * 10)
	//	rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.NewColor(0, 0, 0, alpha))
	//}

	TextX := screenWidth - 150
	if screenWidth < 800 {
		TextX = screenWidth - 130
	}
	rl.DrawText(fmt.Sprintf("Day: %d", numberOfDays), int32(TextX), 10, 20, rl.White)
	rl.DrawText(fmt.Sprintf("Gold: %d G", playerGold), int32(TextX), 30, 20, rl.Gold)

	if shopOpen {
		drawShop()
	} else if bakeryOpen {
		drawBakery()
	} else {
		if playerInv.Open {
			drawInventory()
		}
		if paused {
			drawPauseMenu()
		}
		if showKeyBindings {
			drawKeybindings()
		}
	}
	rl.EndDrawing()
}

func init() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Strawberry Glen")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
	showMessages("Whoa there\nPress ESC to pause catch your breath,\nknow more about the game keybindings, and Save your progress.", 2.0)

	plants := rl.LoadTexture("assets/Objects/Basic_Plants.png")
	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	houseSprite = rl.LoadTexture("assets/Tilesets/Wooden_House_Walls_Tilset.png")
	shopSprite = rl.LoadTexture("assets/Objects/Basic_Furniture.png")
	bakerySprite = rl.LoadTexture("assets/Objects/Basic_Furniture.png")
	hillSprite = rl.LoadTexture("assets/Tilesets/Hills.png")
	tilledSprite = rl.LoadTexture("assets/Tilesets/Tilled_Dirt.png")
	waterSprite = rl.LoadTexture("assets/Tilesets/Water.png")
	fenceSprite = rl.LoadTexture("assets/Tilesets/Fences.png")
	doorSprite = rl.LoadTexture("assets/Tilesets/Doors.png")

	wheatGrowthSprite = plants
	strawberryGrowthSprite = plants

	loadItemSprites()

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	cowSprite = rl.LoadTexture("assets/Characters/FreeCowSprites.png")
	chickenHouseSprite = rl.LoadTexture("assets/Objects/Free_Chicken_House.png")

	playerSprite = rl.LoadTexture("assets/Characters/BasicCharakterSpritesheet.png")
	playerSrc = rl.NewRectangle(0, 0, 48, 48)

	playerDest = rl.NewRectangle(200, 200, 32, 32)
	playerInv = Inventory{
		Open:        false,
		Cursor:      0,
		HoveredSlot: -1,
		Slots:       make([]InventorySlot, InventorySize),
	}
	playerInv.Slots[0] = InventorySlot{ItemID: 3, ItemName: Items[3].Name, ItemReusable: true}
	playerInv.Slots[1] = InventorySlot{ItemID: 4, ItemName: Items[4].Name, ItemReusable: true}
	playerInv.Slots[2] = InventorySlot{ItemID: 5, ItemName: Items[5].Name, ItemReusable: true}
	playerInv.Slots[3] = InventorySlot{ItemID: 6, ItemName: Items[6].Name, ItemReusable: true}

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/audio/amb.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)), // Offset
		rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), // Target
			float32(playerDest.Y-(playerDest.Height/2))),
		0.0, // Rotation
		2.0, // Zoom
	)

	gameSave = false

	showKeyBindings = false

	for i, it := range shopItems {
		initialShopStock[i] = it.Stock
	}

	loadMap("assets/one.map")
	world.InitShopCounter()
	world.InitbakeryCounter()
	world.InitChickenHouse()
	world.InitCows()
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)

	rl.UnloadTexture(wheatGrowthSprite)
	rl.UnloadTexture(strawberryGrowthSprite)
	rl.CloseAudioDevice()
	defer rl.CloseWindow()
}

func selectMenuItem(index int) {
	switch index {
	case 0: // Play
		gameState = StatePlaying
	case 1: // Continue
		if saveFileExists {
			gameState = StatePlaying
		}
	case 2: // Credits
		gameState = StateCredits
	case 3: // Exit
		running = false
	}
}

func handleMainMenu() {
	mp := rl.GetMousePosition()
	menuItemHeight := int32(40)
	baseY := int32(screenHeight/2 - 50)
	newCursor := -1

	for i := range menuItems {
		itemY := baseY + int32(i)*menuItemHeight
		textWidth := rl.MeasureText(menuItems[i], 32)
		bounds := rl.NewRectangle(
			float32(screenWidth/2)-float32(textWidth/2),
			float32(itemY),
			float32(textWidth),
			float32(menuItemHeight),
		)

		if rl.CheckCollisionPointRec(mp, bounds) {
			newCursor = i
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				selectMenuItem(i)
				return
			}
		}
	}

	if newCursor != -1 {
		menuCursor = newCursor
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		menuCursor = (menuCursor + 1) % len(menuItems)
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		menuCursor = (menuCursor - 1 + len(menuItems)) % len(menuItems)
	}

	if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeySpace) {
		switch menuCursor {
		case 0: // Play
			gameState = StatePlaying
			// Add game initialization here
		case 1: // Continue
			if saveFileExists {
				gameState = StatePlaying
				//TODO: Add save loading here
			}
		case 2: // Credits
			gameState = StateCredits
		case 3: // Exit
			running = false
		}
	}
}

func main() {
	for running {
		switch gameState {
		case StateMainMenu:
			handleMainMenu()
			drawMainMenu()
		case StateCredits:
			handleCredits()
			drawCredits()
		default:
			input()
			update()
			render()
		}
	}
	quit()
}
