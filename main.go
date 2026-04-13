package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"gopher-run/internal/game"
)

const (
	screenWidth  = 800
	screenHeight = 400
)

var (
	gopherImage    *ebiten.Image
	dirtImage      *ebiten.Image
	grassTileImage *ebiten.Image
)

func init() {
	img, err := png.Decode(bytes.NewReader(gopherPng))
	if err != nil {
		log.Fatal(err)
	}
	gopherImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(tilesPng))
	if err != nil {
		log.Fatal(err)
	}
	tilesImg := ebiten.NewImageFromImage(img)
	grassTileImage = tilesImg.SubImage(image.Rect(0, 0, game.TileSize, game.TileSize)).(*ebiten.Image)

	dr, dg, db, _ := img.At(game.TileSize/2, game.TileSize-4).RGBA()
	dirtImage = ebiten.NewImage(1, 1)
	dirtImage.Fill(color.RGBA{uint8(dr >> 8), uint8(dg >> 8), uint8(db >> 8), 0xFF})
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Gopher Run")
	g := game.New(gopherImage, dirtImage, grassTileImage)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
