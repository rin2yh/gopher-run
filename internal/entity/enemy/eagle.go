package enemy

import "github.com/hajimehoshi/ebiten/v2"

const (
	eagleW = 55.0
	eagleH = 50.0

	EagleSpawnX = 900.0
	EagleSpeedX = 5.0

	eagleSpawnY      = 150.0
	eagleSpeedY      = 1.5
	eagleDiveFrames  = 60
	eagleCycleFrames = 120
)

type eagle struct {
	x      float64
	y      float64
	frames int
	drawOp ebiten.DrawImageOptions
}

func NewEagle() *eagle {
	return NewEagleAt(EagleSpawnX)
}

func NewEagleAt(x float64) *eagle {
	return &eagle{x: x, y: eagleSpawnY}
}

func (e *eagle) X() float64 {
	return e.x
}

func (e *eagle) Move() {
	e.x -= EagleSpeedX
	e.frames++
	if e.frames < eagleDiveFrames {
		e.y += eagleSpeedY
	} else {
		e.y -= eagleSpeedY
		if e.frames > eagleCycleFrames {
			e.frames = 0
		}
	}
}

func (e *eagle) Hit(px, py, pw, ph float64) bool {
	return px < e.x+eagleW && px+pw > e.x &&
		py < e.y+eagleH && py+ph > e.y
}

func (e *eagle) Draw(screen *ebiten.Image, img *ebiten.Image) {
	e.drawOp.GeoM.Reset()
	e.drawOp.GeoM.Translate(e.x, e.y)
	screen.DrawImage(img, &e.drawOp)
}
