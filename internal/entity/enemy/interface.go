package enemy

import "github.com/hajimehoshi/ebiten/v2"

type Enemy interface {
	Move()
	Hit(px, py, pw, ph float64) bool
	Draw(screen *ebiten.Image, img *ebiten.Image)
	X() float64
}
