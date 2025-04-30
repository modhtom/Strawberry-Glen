package main

import (
	"fmt"
	"os"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	editorScreenWidth  = 832 // 26 tiles * 32px
	editorScreenHeight = 608 // 16 tiles * 32px + room for palette
	tileSize           = 32
)

var (
	tiles = map[int]string{
		1: "g", // grass
		2: "h", // house
		3: "l", // hill
		4: "w", // water
		5: "t", // tilled
		6: "f", // fence
		7: "d", // door
	}
	tileTextures = make(map[int]rl.Texture2D)
	tileColors   = map[int]rl.Color{
		1: rl.Green,
		2: rl.Brown,
		3: rl.DarkGreen,
		4: rl.Blue,
		5: rl.DarkBrown,
		6: rl.Gray,
		7: rl.Orange,
	}

	selectedTile = 1

	editorMapW, editorMapH = 26, 16
	tileIDs                []int
	srcIDs                 []string
)

func initMap(w, h int) {
	tileIDs = make([]int, w*h)
	srcIDs = make([]string, w*h)
	for i := range tileIDs {
		tileIDs[i] = 0
		srcIDs[i] = "g"
	}
}

func drawGrid() {
	for y := 0; y < editorMapH; y++ {
		for x := 0; x < editorMapW; x++ {
			idx := y*editorMapW + x
			tile := tileIDs[idx]
			color := rl.LightGray
			if tile != 0 {
				color = tileColors[tile]
			}
			rl.DrawRectangle(int32(x*tileSize), int32(y*tileSize)+tileSize, tileSize, tileSize, color)
			rl.DrawRectangleLines(int32(x*tileSize), int32(y*tileSize+tileSize), tileSize, tileSize, rl.Black)
		}
	}
}

func drawPalette() {
	for i := 1; i <= len(tiles); i++ {
		color := tileColors[i]
		x := (i - 1) * tileSize
		y := 0
		if i == selectedTile {
			rl.DrawRectangle(int32(x), int32(y), tileSize, tileSize, rl.White)
			rl.DrawRectangle(int32(x)+4, int32(y)+4, tileSize-8, tileSize-8, color)
		} else {
			rl.DrawRectangle(int32(x), int32(y), tileSize, tileSize, color)
		}
		rl.DrawRectangleLines(int32(x), int32(y), tileSize, tileSize, rl.Black)
	}
}

func handleInput() {
	mouse := rl.GetMousePosition()
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		x := int(mouse.X) / tileSize
		y := int(mouse.Y-tileSize) / tileSize
		if x >= 0 && x < editorMapW && y >= 0 && y < editorMapH && mouse.Y > float32(tileSize) {
			idx := y*editorMapW + x
			tileIDs[idx] = selectedTile
			srcIDs[idx] = tiles[selectedTile]
		}
	}

	for i := 1; i <= len(tiles); i++ {
		if rl.IsKeyPressed(int32(rl.KeyOne) + int32(i-1)) {
			selectedTile = i
		}
	}

	if rl.IsKeyPressed(rl.KeyS) {
		saveMap("assets/new.map")
	}
}

func saveMap(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to save map:", err)
		return
	}
	defer file.Close()

	total := editorMapW * editorMapH
	tokens := []string{
		fmt.Sprintf("%03d", editorMapW),
		fmt.Sprintf("%03d", editorMapH),
	}

	for i := 0; i < total; i++ {
		tokens = append(tokens, fmt.Sprintf("%02d", tileIDs[i]))
	}
	for i := 0; i < total; i++ {
		tokens = append(tokens, srcIDs[i])
	}

	output := strings.Join(tokens, " ")
	file.WriteString(output)
	fmt.Println("Map saved to", filename)
}

func runEditor() {
	rl.InitWindow(editorScreenWidth, editorScreenHeight, "Map Editor")
	rl.SetTargetFPS(60)

	initMap(editorMapW, editorMapH)

	for !rl.WindowShouldClose() {
		handleInput()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		drawPalette()
		drawGrid()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
