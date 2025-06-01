package main

import rl "github.com/gen2brain/raylib-go/raylib"

func drawMainMenu() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(147, 211, 196, 255))

	// Draw title
	title := "Strawberry Glen"
	titleWidth := rl.MeasureText(title, 40)
	rl.DrawText(title, int32(screenWidth/2)-titleWidth/2, 100, 40, rl.White)

	// Draw menu items
	for i, item := range menuItems {
		y := 200 + i*60
		color := rl.White
		if i == menuCursor {
			color = rl.Gold
		}

		// Gray out Continue if no save
		if i == 1 && !saveFileExists {
			color = rl.Gray
		}

		textWidth := rl.MeasureText(item, 30)
		rl.DrawText(item, int32(screenWidth/2)-textWidth/2, int32(y), 30, color)
	}

	// Draw version info
	rl.DrawText("v0.1", 10, int32(screenHeight)-30, 20, rl.White)
	rl.EndDrawing()
}

func drawCredits() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(0, 0, 0, 128)) // Semi-transparent background

	// Credits window
	width := 600
	height := 400
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2

	rl.DrawRectangle(int32(x), int32(y), int32(width), int32(height), rl.Black)
	rl.DrawRectangleLines(int32(x), int32(y), int32(width), int32(height), rl.White)

	// Credits content
	credits := []string{
		"Strawberry Glen",
		"Created by MODHTOM",
		" ",
		"Programming: MODHTOM",
		"Art: ",
		"Music: ",
		" ",
		"Press ESC to return",
	}

	startY := y + 30
	for i, line := range credits {
		rl.DrawText(line,
			int32(x+30),
			int32(startY+i*30),
			20,
			rl.White)
	}

	rl.EndDrawing()
}

func handleCredits() {
	if rl.IsKeyPressed(rl.KeyEscape) {
		gameState = StateMainMenu
	}
}
