package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"gopher-run/internal/game"
)

var (
	gopherImage    *ebiten.Image
	dirtImage      *ebiten.Image
	grassTileImage *ebiten.Image
	eagleImage     *ebiten.Image
	snakeImage     *ebiten.Image
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

	img, err = png.Decode(bytes.NewReader(eaglePng))
	if err != nil {
		log.Fatal(err)
	}
	eagleImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(snakePng))
	if err != nil {
		log.Fatal(err)
	}
	snakeImage = composeSnake(ebiten.NewImageFromImage(img))
}

func composeSnake(sheet *ebiten.Image) *ebiten.Image {
	const (
		tile = 42
		step = 30
	)
	head := sheet.SubImage(image.Rect(0, tile*2, tile, tile*3)).(*ebiten.Image)
	body := sheet.SubImage(image.Rect(tile*2, tile*2, tile*3, tile*3)).(*ebiten.Image)
	tail := sheet.SubImage(image.Rect(tile, tile*2, tile*2, tile*3)).(*ebiten.Image)

	dst := ebiten.NewImage(tile+step*2, tile)

	rotated := func(piece *ebiten.Image, x, angle float64) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-tile/2.0, -tile/2.0)
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(tile/2.0, tile/2.0)
		op.GeoM.Translate(x, 0)
		dst.DrawImage(piece, op)
	}

	rotated(body, step, math.Pi/2)
	rotated(tail, step*2, -math.Pi/2)

	headOp := &ebiten.DrawImageOptions{}
	headOp.GeoM.Scale(-1, 1)
	headOp.GeoM.Translate(tile, 0)
	dst.DrawImage(head, headOp)

	return dst
}

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Gopher Run")
	g := game.New(gopherImage, dirtImage, grassTileImage, eagleImage, snakeImage)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
