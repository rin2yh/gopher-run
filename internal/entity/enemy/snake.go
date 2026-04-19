package enemy

import "github.com/hajimehoshi/ebiten/v2"

const (
	snakeW = 102.0
	snakeH = 42.0

	SnakeSpawnX = 1200.0
	SnakeSpeedX = 4.0

	snakeY = 330.0
)

type Snake struct {
	x      float64
	img    *ebiten.Image
	drawOp ebiten.DrawImageOptions
}

func NewSnake(img *ebiten.Image) *Snake {
	return NewSnakeAt(SnakeSpawnX, img)
}

func NewSnakeAt(x float64, img *ebiten.Image) *Snake {
	return &Snake{x: x, img: img}
}

func (s *Snake) X() float64 {
	return s.x
}

func (s *Snake) Move() {
	s.x -= SnakeSpeedX
}

func (s *Snake) Hit(px, py, pw, ph float64, digging bool) bool {
	if !digging {
		return false
	}
	return aabbOverlap(px, py, pw, ph, s.x, snakeY, snakeW, snakeH)
}

func (s *Snake) Draw(screen *ebiten.Image) {
	s.drawOp.GeoM.Reset()
	s.drawOp.GeoM.Translate(s.x, snakeY)
	screen.DrawImage(s.img, &s.drawOp)
}
