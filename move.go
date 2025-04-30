package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func canMove(newX, newY float32) bool {
	centerX := newX + playerDest.Width/2
	centerY := newY + playerDest.Height/2

	tileX := int(centerX / 16)
	tileY := int(centerY / 16)

	if tileX < 0 || tileX >= mapW || tileY < 0 || tileY >= mapH {
		showMessages("You can't go beyond the map!", 0.5)
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
			showMessages("You can't walk in water!", 0.5)
		case "f":
			showMessages("The fence blocks your path.", 0.5)
		case "h", "l":
			showMessages("You can't walk through this.", 0.5)
		case "d":
			showMessages("The door is closed.", 0.5)
		case "s":
			showMessages("Welcome to the shop Press 'B'.", 0.75)
		case "b":
			showMessages("Welcome to the bakery Press 'B'.", 0.75)
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
