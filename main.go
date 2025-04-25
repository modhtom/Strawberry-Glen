package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenHeight = 480
	screenWidth  = 1000
)

var (
	running = true

	tex      rl.Texture2D
	bkgColor = rl.NewColor(147, 211, 196, 255)

	grassSprite  rl.Texture2D
	houseSprite  rl.Texture2D
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

	showMessage     bool
	messageText     string
	messageTimer    float32
	messageDuration float32 = 2.0

	showKeyBindings bool
	gameSave        bool
)

func drawScene() {
	for i := range tileMap {
		if tileMap[i] == 0 {
			continue
		}
		tileDest.X = tileDest.Width * float32(i%mapW)
		tileDest.Y = tileDest.Height * float32(i/mapW)
		switch srcMap[i] {
		case "g":
			tex = grassSprite
		case "l":
			tex = hillSprite
		case "h":
			tex = houseSprite
		case "w":
			tex = waterSprite
		case "t":
			tex = tilledSprite
		case "f":
			tex = fenceSprite
		case "d":
			tex = doorSprite
		default:
			tex = grassSprite
		}
		if srcMap[i] == "h" || srcMap[i] == "d" || srcMap[i] == "f" {
			tileSrc.X = 10
			tileSrc.Y = 25
			rl.DrawTexturePro(grassSprite, tileSrc, tileDest,
				rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}

		tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
		tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))

		rl.DrawTexturePro(tex, tileSrc, tileDest,
			rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
	}
	rl.DrawTexturePro(
		playerSprite,
		playerSrc,
		rl.NewRectangle(playerDest.X, playerDest.Y, 64, 64),
		rl.NewVector2(32, 32),
		0,
		rl.White,
	)
}

func drawPauseMenu() {
	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.NewColor(0, 0, 0, 180))

	box := rl.NewRectangle(screenWidth/2-120, screenHeight/2-100, 400, 220)
	rl.DrawRectangleRec(box, rl.NewColor(80, 10, 10, 230))
	rl.DrawRectangleLines(int32(box.X), int32(box.Y), int32(box.Width), int32(box.Height), rl.White)
	rl.DrawText("PAUSED", int32(box.X+60), int32(box.Y+10), 30, rl.White)

	buttons := []struct {
		rect   rl.Rectangle
		text   string
		action func()
	}{
		{rl.NewRectangle(box.X+20, box.Y+60, 200, 40), "Return to Game", func() { paused = false }},
		{rl.NewRectangle(box.X+20, box.Y+110, 200, 40), "Toggle Sound", func() { musicPaused = !musicPaused }},
		{rl.NewRectangle(box.X+20, box.Y+160, 200, 40), "Save", func() { gameSave = !gameSave }},
		{rl.NewRectangle(box.X+20, box.Y+210, 200, 40), "keybindings", func() { showKeyBindings = !showKeyBindings }},
		{rl.NewRectangle(box.X+20, box.Y+260, 200, 40), "Main Menu", func() { running = false }},
	}

	mousePos := rl.GetMousePosition()
	for _, btn := range buttons {
		color := rl.DarkGray
		if rl.CheckCollisionPointRec(mousePos, btn.rect) {
			color = rl.Gray
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				btn.action()
			}
		}
		rl.DrawRectangleRec(btn.rect, color)
		rl.DrawText(btn.text, int32(btn.rect.X+10), int32(btn.rect.Y+10), 20, rl.White)
	}
}

func drawKeybindings() {
	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.NewColor(0, 0, 0, 180))

	box := rl.NewRectangle(screenWidth/2-120, screenHeight/2-100, 250, 300)
	rl.DrawRectangleRec(box, rl.NewColor(80, 10, 10, 230))
	rl.DrawRectangleLines(int32(box.X), int32(box.Y), int32(box.Width), int32(box.Height), rl.White)
	rl.DrawText("keybindings", int32(box.X+60), int32(box.Y+10), 30, rl.White)
	rl.DrawText("Use arrow keys or WASD to move.", int32(box.X+20), int32(box.Y+50), 30, rl.White)
	rl.DrawText("Zoom in with 'Z', zoom out with 'X'.", int32(box.X+20), int32(box.Y+90), 30, rl.White)
	rl.DrawText("Use 'I' to open your Inventory.", int32(box.X+20), int32(box.Y+130), 30, rl.White)
	rl.DrawText("Use 'B' to open Shop.", int32(box.X+20), int32(box.Y+170), 30, rl.White)
	buttons := []struct {
		rect   rl.Rectangle
		text   string
		action func()
	}{
		{rl.NewRectangle(box.X+100, box.Y+200, 200, 40), "Close Window", func() { showKeyBindings = false }}}

	mousePos := rl.GetMousePosition()
	for _, btn := range buttons {
		color := rl.DarkGray
		if rl.CheckCollisionPointRec(mousePos, btn.rect) {
			color = rl.Gray
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				btn.action()
			}
		}
		rl.DrawRectangleRec(btn.rect, color)
		rl.DrawText(btn.text, int32(btn.rect.X+10), int32(btn.rect.Y+10), 20, rl.White)
	}
}

func input() {
	if rl.IsKeyPressed(rl.KeyEscape) {
		paused = !paused
	}

	if paused {
		return
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

	if rl.IsKeyPressed(rl.KeyB) {
		shopOpen = !shopOpen
		shopCursor = 0
		shopMode = 0
	}
	if shopOpen {
		// navigate shop
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

func canMove(newX, newY float32) bool {
	centerX := newX + playerDest.Width/2
	centerY := newY + playerDest.Height/2

	tileX := int(centerX / 16)
	tileY := int(centerY / 16)

	if tileX < 0 || tileX >= mapW || tileY < 0 || tileY >= mapH {
		showMessages("You can't go beyond the map!")
		return false
	}

	index := tileY*mapW + tileX
	if index < 0 || index >= len(srcMap) {
		return false
	}

	tileType := srcMap[index]
	if isImpassable(tileType) {
		switch tileType {
		case "w":
			showMessages("You can't walk in water!")
		case "f":
			showMessages("The fence blocks your path.")
		case "h", "l":
			showMessages("You can't walk through this.")
		case "d":
			showMessages("The door is closed.")
		}
		return false
	}

	return true
}

func isImpassable(tileType string) bool {
	switch tileType {
	case "w", "f", "h", "d", "l":
		return true
	default:
		return false
	}
}

func getCurrentTileTypes(x, y float32) []string {
	playerRect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  playerDest.Width,
		Height: playerDest.Height,
	}

	tileSize := float32(16)
	startTileX := int(playerRect.X / tileSize)
	endTileX := int((playerRect.X + playerRect.Width) / tileSize)
	startTileY := int(playerRect.Y / tileSize)
	endTileY := int((playerRect.Y + playerRect.Height) / tileSize)

	var tiles []string
	for y := startTileY; y <= endTileY; y++ {
		for x := startTileX; x <= endTileX; x++ {
			index := y*mapW + x
			if index >= 0 && index < len(srcMap) {
				tiles = append(tiles, srcMap[index])
			}
		}
	}
	return tiles
}

func update() {
	running = !rl.WindowShouldClose()

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
		showMessages("You slip in the mud!")
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
		showMessages("Shop has been restocked!")
	}

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

	if shopOpen {
		drawShop()
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

func loadMap(mapFile string) {
	data, err := os.ReadFile(mapFile)
	if err != nil {
		fmt.Printf("failed to read map file: %v", err)
	}

	tokens := strings.Fields(string(data))
	if len(tokens) < 2 {
		fmt.Printf("map file too short: %d tokens", len(tokens))
	}

	var errW, errH error
	mapW, errW = strconv.Atoi(tokens[0])
	mapH, errH = strconv.Atoi(tokens[1])
	if errW != nil || errH != nil {
		fmt.Printf("invalid dimensions in map file: %v, %v", errW, errH)
	}

	total := mapW * mapH
	if len(tokens) < 2+2*total {
		fmt.Printf("expected %d tiles %d src glyphs, got %d tokens",
			total, total, len(tokens)-2)
	}
	tileMap = make([]int, total)
	srcMap = make([]string, total)
	for i := range total {
		tileMap[i], _ = strconv.Atoi(tokens[2+i])
	}
	for i := range total {
		srcMap[i] = tokens[2+total+i]
	}
}

func showMessages(text string) {
	showMessage = true
	messageText = text
	messageTimer = 0
}
func drawMessage() {
	if !showMessage {
		return
	}

	messageTimer += 0.005
	if messageTimer > messageDuration {
		showMessage = false
		return
	}

	lines := strings.Split(messageText, "\n")
	var fontSize int32 = 20
	lineHeight := fontSize + 5
	totalHeight := int32(len(lines)) * lineHeight

	var maxWidth int32 = 0
	for _, line := range lines {
		width := rl.MeasureText(line, fontSize)
		if width > maxWidth {
			maxWidth = width
		}
	}

	bgX := screenWidth/2 - maxWidth/2 - 10
	bgY := screenHeight - totalHeight - 20
	rl.DrawRectangle(bgX, bgY, maxWidth+20, totalHeight+10, rl.NewColor(0, 0, 0, 200))

	for i, line := range lines {
		rl.DrawText(line, screenWidth/2-rl.MeasureText(line, fontSize)/2, bgY+5+int32(i)*lineHeight, fontSize, rl.White)
	}
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "Strawberry Glen")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
	showMessages("Whoa there\nPress ESC to pause catch your breath,\nknow more about the game keybindings, and Save your progress.")

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	houseSprite = rl.LoadTexture("assets/Tilesets/Wooden_House_Walls_Tilset.png")
	hillSprite = rl.LoadTexture("assets/Tilesets/Hills.png")
	tilledSprite = rl.LoadTexture("assets/Tilesets/Tilled_Dirt.png")
	waterSprite = rl.LoadTexture("assets/Tilesets/Water.png")
	fenceSprite = rl.LoadTexture("assets/Tilesets/Fences.png")
	doorSprite = rl.LoadTexture("assets/Tilesets/Doors.png")
	loadItemSprites()

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSprite = rl.LoadTexture("assets/Characters/BasicCharakterSpritesheet.png")
	playerSrc = rl.NewRectangle(0, 0, 48, 48)

	playerDest = rl.NewRectangle(200, 200, 32, 32)
	playerInv = Inventory{Open: false, Cursor: 0, HoveredSlot: -1}
	playerInv.Slots[0] = InventorySlot{ItemID: 3, ItemName: Items[3].Name, ItemReusable: true}                     // Watering Can
	playerInv.Slots[1] = InventorySlot{ItemID: 4, ItemName: Items[4].Name, ItemReusable: true}                     // Hoe
	playerInv.Slots[2] = InventorySlot{ItemID: 5, ItemName: Items[5].Name, ItemReusable: true}                     // Axe
	playerInv.Slots[3] = InventorySlot{ItemID: 1, ItemName: Items[1].Name, ItemReusable: false, ItemQuantity: 5}   // Wheat Seeds
	playerInv.Slots[4] = InventorySlot{ItemID: 22, ItemName: Items[22].Name, ItemReusable: false, ItemQuantity: 6} // Egg

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/audio/amb.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)),
			float32(playerDest.Y-(playerDest.Height/2))), 0.0, 2.0)

	gameSave = false

	showKeyBindings = false

	for i, it := range shopItems {
		initialShopStock[i] = it.Stock
	}

	loadMap("assets/one.map")
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	defer rl.CloseWindow()
}

func main() {
	for running {
		input()
		update()
		render()
	}
	quit()
}
