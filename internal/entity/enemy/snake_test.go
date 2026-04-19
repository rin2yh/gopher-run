package enemy

import "testing"

type fakeGround struct {
	hasGround bool
}

func (f fakeGround) IsGroundAt(_ int) bool { return f.hasGround }

var groundAlways = fakeGround{hasGround: true}
var groundNever = fakeGround{hasGround: false}

func TestNewSnake(t *testing.T) {
	s := NewSnake(nil)
	if s.x != SnakeSpawnX {
		t.Errorf("x = %v, want %v", s.x, SnakeSpawnX)
	}
	if s.y != snakeInitialY {
		t.Errorf("y = %v, want %v", s.y, snakeInitialY)
	}
}

func TestNewSnakeAt(t *testing.T) {
	s := NewSnakeAt(500, nil)
	if s.x != 500 {
		t.Errorf("x = %v, want 500", s.x)
	}
	if s.y != snakeInitialY {
		t.Errorf("y = %v, want %v", s.y, snakeInitialY)
	}
}

func TestSnakeX(t *testing.T) {
	s := NewSnakeAt(500, nil)
	if s.X() != 500 {
		t.Errorf("X() = %v, want 500", s.X())
	}
}

func TestSnakeY(t *testing.T) {
	s := NewSnakeAt(500, nil)
	if s.Y() != snakeInitialY {
		t.Errorf("Y() = %v, want %v", s.Y(), snakeInitialY)
	}
}

func TestSnakeMove_HorizontalSpeed(t *testing.T) {
	s := NewSnakeAt(200, nil)
	s.Move(groundAlways, 0)
	if s.x != 200-SnakeSpeedX {
		t.Errorf("x after 1 move = %v, want %v", s.x, 200-SnakeSpeedX)
	}
}

func TestSnakeMove_OverGround_NoFall(t *testing.T) {
	s := NewSnakeAt(200, nil)
	for range 10 {
		s.Move(groundAlways, 0)
	}
	if s.y != snakeInitialY {
		t.Errorf("y should remain %v on ground, got %v", snakeInitialY, s.y)
	}
	if s.vy != 0 {
		t.Errorf("vy should remain 0 on ground, got %v", s.vy)
	}
}

func TestSnakeMove_OverHole_Falls(t *testing.T) {
	s := NewSnakeAt(200, nil)
	prevY := s.y
	for range 5 {
		s.Move(groundNever, 0)
		if s.y <= prevY {
			t.Errorf("y should increase over hole, got %v <= %v", s.y, prevY)
		}
		prevY = s.y
	}
}

func TestSnakeMove_FallVYCapped(t *testing.T) {
	s := NewSnakeAt(200, nil)
	for range 1000 {
		s.Move(groundNever, 0)
	}
	if s.vy > snakeMaxFallVY {
		t.Errorf("vy = %v, want <= %v", s.vy, snakeMaxFallVY)
	}
}

func TestSnakeHit_NotDigging(t *testing.T) {
	s := &Snake{x: 100, y: snakeInitialY}
	// even with perfect geometry overlap, !digging → no hit
	if s.Hit(100, snakeInitialY, 60, 75, false) {
		t.Error("expected no hit when player is not digging")
	}
}

func TestSnakeHit_DiggingOverlap(t *testing.T) {
	s := &Snake{x: 100, y: snakeInitialY}
	if !s.Hit(100, snakeInitialY, 60, 75, true) {
		t.Error("expected hit when player digging and overlaps snake")
	}
}

func TestSnakeHit_PlayerLeft(t *testing.T) {
	s := &Snake{x: 100, y: snakeInitialY}
	// player right edge exactly at snake left edge: px+pw == s.x → no hit
	if s.Hit(40, snakeInitialY, 60, 75, true) {
		t.Error("expected no hit when player right edge meets snake left edge")
	}
}

func TestSnakeHit_PlayerRight(t *testing.T) {
	s := &Snake{x: 100, y: snakeInitialY}
	if s.Hit(100+snakeW, snakeInitialY, 60, 75, true) {
		t.Error("expected no hit when player left edge meets snake right edge")
	}
}

func TestSnakeHit_PlayerOnGround(t *testing.T) {
	s := &Snake{x: 100, y: snakeInitialY}
	// player on ground: top y = 245, bottom = 320. Snake y = 330-372. No overlap.
	if s.Hit(100, 245, 60, 75, true) {
		t.Error("expected no hit when player on ground y range")
	}
}
