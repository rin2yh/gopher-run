package enemy

import "github.com/hajimehoshi/ebiten/v2"

type Enemy interface {
	Move()
	Hit(px, py, pw, ph float64, digging bool) bool
	Draw(screen *ebiten.Image)
	X() float64
}

func aabbOverlap(px, py, pw, ph, ex, ey, ew, eh float64) bool {
	return px < ex+ew && px+pw > ex && py < ey+eh && py+ph > ey
}
