package enemy

import "github.com/hajimehoshi/ebiten/v2"

const (
	snakeW = 102.0
	snakeH = 42.0

	SnakeSpawnX = 1200.0
	SnakeSpeedX = 4.0

	snakeInitialY  = 330.0
	snakeGravity   = 4.0
	snakeMaxFallVY = 30.0
)

type Snake struct {
	x      float64
	y      float64
	vy     float64
	img    *ebiten.Image
	drawOp ebiten.DrawImageOptions
}

func NewSnake(img *ebiten.Image) *Snake {
	return NewSnakeAt(SnakeSpawnX, img)
}

func NewSnakeAt(x float64, img *ebiten.Image) *Snake {
	return &Snake{x: x, y: snakeInitialY, img: img}
}

func (s *Snake) X() float64 {
	return s.x
}

func (s *Snake) Y() float64 {
	return s.y
}

func (s *Snake) Move(w GroundChecker, cameraX int) {
	s.x -= SnakeSpeedX
	leftWorldX := int(s.x) + cameraX
	rightWorldX := int(s.x+snakeW-1) + cameraX
	overGround := w.IsGroundAt(leftWorldX) || w.IsGroundAt(rightWorldX)
	if overGround {
		s.vy = 0
		if s.y > snakeInitialY {
			s.y = snakeInitialY
		}
		return
	}
	s.vy += snakeGravity
	if s.vy > snakeMaxFallVY {
		s.vy = snakeMaxFallVY
	}
	s.y += s.vy
}

func (s *Snake) Hit(px, py, pw, ph float64, digging bool) bool {
	if !digging {
		return false
	}
	return aabbOverlap(px, py, pw, ph, s.x, s.y, snakeW, snakeH)
}

func (s *Snake) IsOffScreen(screenHeight int) bool {
	return s.x < 0 || s.y > float64(screenHeight)
}

func (s *Snake) OnDodged(_ DodgeObserver, _ DodgeContext) {}

func (s *Snake) Respawn(svc RespawnService) Enemy {
	return NewSnakeAt(svc.SafeSnakeSpawnX(SnakeSpawnX), svc.SnakeImage())
}

func (s *Snake) Draw(screen *ebiten.Image) {
	s.drawOp.GeoM.Reset()
	s.drawOp.GeoM.Translate(s.x, s.y)
	screen.DrawImage(s.img, &s.drawOp)
}
