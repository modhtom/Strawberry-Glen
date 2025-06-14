package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var messageDuration float32

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
		case "s":
			tex = shopSprite
		case "b":
			tex = bakerySprite
		case "c":
			tex = chickenHouseSprite
		case "o":
			tex = cowSprite
		default:
			tex = grassSprite
		}
		if srcMap[i] == "h" || srcMap[i] == "d" || srcMap[i] == "f" || srcMap[i] == "s" || srcMap[i] == "b" || srcMap[i] == "c" || srcMap[i] == "o" {
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

	for _, crop := range world.Crops {
		if crop != nil {
			var cropTexture rl.Texture2D
			switch crop.CropTypeID {
			case 11: // Strawberry
				cropTexture = strawberryGrowthSprite
			default: // Default to wheat (10)
				cropTexture = wheatGrowthSprite
			}

			cropSrcRect := crop.GetSpriteRect()

			cropDestRect := rl.NewRectangle(
				float32(crop.PosX-1)*tileDest.Width,
				float32(crop.PosY-1)*tileDest.Height,
				tileDest.Width,
				tileDest.Height,
			)
			rl.DrawTexturePro(cropTexture, cropSrcRect, cropDestRect, rl.NewVector2(0, 0), 0, rl.White)
		}
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
	rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.NewColor(0, 0, 0, 180))

	box := rl.NewRectangle(
		float32(screenWidth/2)-120,
		float32(screenHeight/2)-100,
		400,
		220,
	)
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
		{rl.NewRectangle(box.X+20, box.Y+260, 200, 40), "Main Menu", func() { gameState = StateMainMenu }},
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
	rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.NewColor(0, 0, 0, 180))

	box := rl.NewRectangle(
		float32(screenWidth/2)-200,
		float32(screenHeight/2)-150,
		400,
		400,
	)
	rl.DrawRectangleRec(box, rl.NewColor(80, 10, 10, 230))
	rl.DrawRectangleLines(int32(box.X), int32(box.Y), int32(box.Width), int32(box.Height), rl.White)
	rl.DrawText("KEY BINDINGS", int32(box.X+140), int32(box.Y+20), 20, rl.White)

	startY := int32(box.Y + 60)
	rl.DrawText("Movement:", int32(box.X+20), startY, 20, rl.Gold)
	rl.DrawText("WASD/Arrow Keys - Move", int32(box.X+40), startY+30, 20, rl.White)
	rl.DrawText("E - Interact", int32(box.X+40), startY+60, 20, rl.White)
	rl.DrawText("H - Tile Soil", int32(box.X+40), startY+80, 20, rl.White)

	rl.DrawText("Inventory/Shop:", int32(box.X+20), startY+100, 20, rl.Gold)
	rl.DrawText("I - Open Inventory", int32(box.X+40), startY+130, 20, rl.White)
	rl.DrawText("B - Open Shop/Bakery", int32(box.X+40), startY+160, 20, rl.White)
	rl.DrawText(fmt.Sprintf("1-%v - Quick Select Slots", InventorySize), int32(box.X+40), startY+190, 20, rl.White)

	rl.DrawText("Baking:", int32(box.X+20), startY+230, 20, rl.Gold)
	rl.DrawText("TAB - Switch Buy/Bake", int32(box.X+40), startY+260, 20, rl.White)
	rl.DrawText("ALT+1-6 - Remove Ingredients", int32(box.X+40), startY+290, 20, rl.White)
	rl.DrawText("ENTER - Confirm Action", int32(box.X+40), startY+320, 20, rl.White)

	rl.DrawText("System:", int32(box.X+20), startY+360, 20, rl.Gold)
	rl.DrawText("ESC - Pause/Menu", int32(box.X+40), startY+390, 20, rl.White)
	rl.DrawText("F11 - Fullscreen", int32(box.X+40), startY+420, 20, rl.White)
}

func showMessages(text string, aMessageDuration float32) {
	showMessage = true
	messageText = text
	messageTimer = 0
	messageDuration = aMessageDuration
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

	bgX := int32(screenWidth/2) - maxWidth/2 - 10
	bgY := int32(screenHeight) - totalHeight - 20
	rl.DrawRectangle(bgX, bgY, maxWidth+20, totalHeight+10, rl.NewColor(0, 0, 0, 200))

	for i, line := range lines {
		textWidth := rl.MeasureText(line, fontSize)
		rl.DrawText(line,
			int32(screenWidth/2)-textWidth/2,
			bgY+5+int32(i)*lineHeight,
			fontSize,
			rl.White)
	}
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
