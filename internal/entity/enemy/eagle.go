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

type Eagle struct {
	x      float64
	y      float64
	frames int
	img    *ebiten.Image
	drawOp ebiten.DrawImageOptions
}

func NewEagle(img *ebiten.Image) *Eagle {
	return NewEagleAt(EagleSpawnX, img)
}

func NewEagleAt(x float64, img *ebiten.Image) *Eagle {
	return &Eagle{x: x, y: eagleSpawnY, img: img}
}

func (e *Eagle) X() float64 {
	return e.x
}

func (e *Eagle) Move() {
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

func (e *Eagle) Hit(px, py, pw, ph float64, _ bool) bool {
	return aabbOverlap(px, py, pw, ph, e.x, e.y, eagleW, eagleH)
}

func (e *Eagle) Draw(screen *ebiten.Image) {
	e.drawOp.GeoM.Reset()
	e.drawOp.GeoM.Translate(e.x, e.y)
	screen.DrawImage(e.img, &e.drawOp)
}
