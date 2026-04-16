package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

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
	assets := &scene.Assets{Gopher: gopherImg, Dirt: dirtImg, GrassTile: grassTileImg, Eagle: eagleImg}
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
