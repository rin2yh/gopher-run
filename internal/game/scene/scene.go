package scene

import "github.com/hajimehoshi/ebiten/v2"

const (
	ScreenWidth  = 800
	ScreenHeight = 400
	TileSize     = 32
)

type Scene interface {
	Update() Scene
	Draw(screen *ebiten.Image)
}

type Assets struct {
	Gopher    *ebiten.Image
	Dirt      *ebiten.Image
	GrassTile *ebiten.Image
}
