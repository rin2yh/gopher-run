package enemy

import "testing"

func TestNewSnake(t *testing.T) {
	s := NewSnake(nil)
	if s.x != SnakeSpawnX {
		t.Errorf("x = %v, want %v", s.x, SnakeSpawnX)
	}
}

func TestNewSnakeAt(t *testing.T) {
	s := NewSnakeAt(500, nil)
	if s.x != 500 {
		t.Errorf("x = %v, want 500", s.x)
	}
}

func TestSnakeX(t *testing.T) {
	s := NewSnakeAt(500, nil)
	if s.X() != 500 {
		t.Errorf("X() = %v, want 500", s.X())
	}
}

func TestSnakeMove_HorizontalSpeed(t *testing.T) {
	s := NewSnakeAt(200, nil)
	s.Move()
	if s.x != 200-SnakeSpeedX {
		t.Errorf("x after 1 move = %v, want %v", s.x, 200-SnakeSpeedX)
	}
}

func TestSnakeHit_NotDigging(t *testing.T) {
	s := &Snake{x: 100}
	// even with perfect geometry overlap, !digging → no hit
	if s.Hit(100, snakeY, 60, 75, false) {
		t.Error("expected no hit when player is not digging")
	}
}

func TestSnakeHit_DiggingOverlap(t *testing.T) {
	s := &Snake{x: 100}
	if !s.Hit(100, snakeY, 60, 75, true) {
		t.Error("expected hit when player digging and overlaps snake")
	}
}

func TestSnakeHit_PlayerLeft(t *testing.T) {
	s := &Snake{x: 100}
	// player right edge exactly at snake left edge: px+pw == s.x → no hit
	if s.Hit(40, snakeY, 60, 75, true) {
		t.Error("expected no hit when player right edge meets snake left edge")
	}
}

func TestSnakeHit_PlayerRight(t *testing.T) {
	s := &Snake{x: 100}
	if s.Hit(100+snakeW, snakeY, 60, 75, true) {
		t.Error("expected no hit when player left edge meets snake right edge")
	}
}

func TestSnakeHit_PlayerOnGround(t *testing.T) {
	s := &Snake{x: 100}
	// player on ground: top y = 245, bottom = 320. Snake y = 330-372. No overlap.
	if s.Hit(100, 245, 60, 75, true) {
		t.Error("expected no hit when player on ground y range")
	}
}
