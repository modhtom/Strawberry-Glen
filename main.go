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
	flowerSprite rl.Texture2D
	tilledSprite rl.Texture2D

	tileDest   rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	playerSprite                                  rl.Texture2D
	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerSpeed                                   float32 = 3
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int

	frameCount int

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D

	paused = false

	showMessage     bool
	messageText     string
	messageTimer    float32
	messageDuration float32 = 2.0
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
			tex = flowerSprite
		default:
			tex = grassSprite
		}
		if srcMap[i] == "h" || srcMap[i] == "f" {
			tileSrc.X = 0
			tileSrc.Y = 0
			rl.DrawTexturePro(grassSprite, tileSrc, tileDest,
				rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}

		tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
		tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))

		rl.DrawTexturePro(tex, tileSrc, tileDest,
			rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
	}
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest,
		rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func drawPauseMenu() {
	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.NewColor(0, 0, 0, 180))

	box := rl.NewRectangle(screenWidth/2-120, screenHeight/2-100, 240, 220)
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
		{rl.NewRectangle(box.X+20, box.Y+160, 200, 40), "Main Menu", func() { running = false }},
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

func input() {
	if rl.IsKeyPressed(rl.KeyEscape) {
		paused = !paused
	}

	if paused {
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

func update() {
	running = !rl.WindowShouldClose()

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
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

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
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

	if paused {
		drawPauseMenu()
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
		fmt.Printf("expected %d tiles + %d src glyphs, got %d tokens",
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

	messageTimer += 0.015
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
	showMessages("Whoa there, speedster! Use arrow keys or WASD to move.\nZoom in with 'Z', zoom out with 'X'.\nPress ESC to pause and catch your breath.")

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	houseSprite = rl.LoadTexture("assets/Tilesets/Wooden_House_Walls_Tilset.png")
	hillSprite = rl.LoadTexture("assets/Tilesets/Hills.png")
	tilledSprite = rl.LoadTexture("assets/Tilesets/Tilled_Dirt.png")
	waterSprite = rl.LoadTexture("assets/Tilesets/Water.png")
	flowerSprite = rl.LoadTexture("assets/Tilesets/Fences.png")

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSprite = rl.LoadTexture("assets/Characters/BasicCharakterSpritesheet.png")
	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 100, 100)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/audio/amb.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)),
			float32(playerDest.Y-(playerDest.Height/2))), 0.0, 1.0)

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
