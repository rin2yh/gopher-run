package game

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gobold"

	"gopher-run/internal/game/scene"
	"gopher-run/internal/input"
)

const (
	ScreenWidth  = scene.ScreenWidth
	ScreenHeight = scene.ScreenHeight
	TileSize     = scene.TileSize
)

type Game struct {
	mode   scene.Scene
	input  *input.Handler
	assets *scene.Assets
}

func New(gopherImg, dirtImg, grassTileImg, eagleImg *ebiten.Image) *Game {
	src, err := ebitentext.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		log.Fatal("failed to load gobold font:", err)
	}
	assets := &scene.Assets{Gopher: gopherImg, Dirt: dirtImg, GrassTile: grassTileImg, Eagle: eagleImg, FontSource: src}
	h := input.NewHandler()
	return &Game{mode: scene.NewTitleScene(assets, h), input: h, assets: assets}
}

func (g *Game) Update() error {
	if next := g.mode.Update(); next != nil {
		g.mode = next
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	g.mode.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
