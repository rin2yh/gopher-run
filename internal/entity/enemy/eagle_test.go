package enemy

import "testing"

func TestNewEagle(t *testing.T) {
	e := NewEagle()
	if e.x != EagleSpawnX {
		t.Errorf("x = %v, want %v", e.x, EagleSpawnX)
	}
	if e.y != eagleSpawnY {
		t.Errorf("y = %v, want %v", e.y, eagleSpawnY)
	}
}

func TestNewEagleAt(t *testing.T) {
	e := NewEagleAt(1000)
	if e.x != 1000 {
		t.Errorf("x = %v, want 1000", e.x)
	}
	if e.y != eagleSpawnY {
		t.Errorf("y = %v, want %v", e.y, eagleSpawnY)
	}
}

func TestEagleX(t *testing.T) {
	e := NewEagleAt(500)
	if e.X() != 500 {
		t.Errorf("X() = %v, want 500", e.X())
	}
}

func TestEagleMove_HorizontalSpeed(t *testing.T) {
	e := NewEagleAt(200)
	e.Move()
	if e.x != 200-EagleSpeedX {
		t.Errorf("x after 1 move = %v, want %v", e.x, 200-EagleSpeedX)
	}
}

func TestEagleMove_DivePhase(t *testing.T) {
	e := NewEagleAt(200)
	// frames 1..eagleDiveFrames-1: y increases
	for range eagleDiveFrames - 1 {
		prev := e.y
		e.Move()
		if e.y <= prev {
			t.Errorf("y should increase during dive, got %v <= %v", e.y, prev)
		}
	}
}

func TestEagleMove_AscentPhase(t *testing.T) {
	e := NewEagleAt(200)
	// advance to ascent phase
	for range eagleDiveFrames {
		e.Move()
	}
	yAtTransition := e.y
	e.Move() // first ascent frame
	if e.y >= yAtTransition {
		t.Errorf("y should decrease in ascent phase, got %v >= %v", e.y, yAtTransition)
	}
}

func TestEagleMove_CycleReset(t *testing.T) {
	e := NewEagleAt(200)
	// eagleCycleFrames+1 moves triggers frames reset to 0
	for range eagleCycleFrames + 1 {
		e.Move()
	}
	if e.frames != 0 {
		t.Errorf("frames after full cycle = %v, want 0", e.frames)
	}
}

func TestEagleHit_Overlap(t *testing.T) {
	e := &eagle{x: 100, y: 100}
	if !e.Hit(100, 100, 60, 75) {
		t.Error("expected hit when player overlaps eagle")
	}
}

func TestEagleHit_PlayerLeft(t *testing.T) {
	e := &eagle{x: 100, y: 100}
	// player right edge exactly at eagle left edge: px+pw == e.x → no hit
	if e.Hit(40, 100, 60, 75) {
		t.Error("expected no hit when player right edge meets eagle left edge")
	}
}

func TestEagleHit_PlayerRight(t *testing.T) {
	e := &eagle{x: 100, y: 100}
	// player left edge exactly at eagle right edge: px == e.x+eagleW → no hit
	if e.Hit(100+eagleW, 100, 60, 75) {
		t.Error("expected no hit when player left edge meets eagle right edge")
	}
}

func TestEagleHit_PlayerAbove(t *testing.T) {
	e := &eagle{x: 100, y: 100}
	// player bottom edge exactly at eagle top edge: py+ph == e.y → no hit
	if e.Hit(100, 50, 60, 50) {
		t.Error("expected no hit when player bottom edge meets eagle top edge")
	}
}

func TestEagleHit_PlayerBelow(t *testing.T) {
	e := &eagle{x: 100, y: 100}
	// player top edge exactly at eagle bottom edge: py == e.y+eagleH → no hit
	if e.Hit(100, 100+eagleH, 60, 75) {
		t.Error("expected no hit when player top edge meets eagle bottom edge")
	}
}
