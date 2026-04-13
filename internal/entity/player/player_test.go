package player

import (
	"testing"

	"gopher-run/internal/world"
)

func TestReset(t *testing.T) {
	p := &Player{isFalling: true}
	p.Reset()

	if p.y16 != groundY16 {
		t.Errorf("y16 = %d, want %d", p.y16, groundY16)
	}
	if p.vy16 != 0 {
		t.Errorf("vy16 = %d, want 0", p.vy16)
	}
	if p.jumpFrames != 0 {
		t.Errorf("jumpFrames = %d, want 0", p.jumpFrames)
	}
	if p.isFalling {
		t.Error("isFalling = true, want false")
	}
}

func TestScreenY(t *testing.T) {
	cases := []struct {
		name string
		y16  int
		want int
	}{
		{"y16がリセット後の初期値のとき画面Y座標245を返す", 3920, 245},
		{"y16が0のとき0を返す", 0, 0},
		{"y16が16で割り切れないとき切り捨てる", 32, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := &Player{y16: c.y16}
			got := p.ScreenY()
			if got != c.want {
				t.Errorf("ScreenY() = %d, want %d", got, c.want)
			}
		})
	}
}

func holeWorld() *world.World {
	return &world.World{
		Segments: []world.Segment{
			{X: 0, Width: 80, IsHole: false},
			{X: 80, Width: 100, IsHole: true},
			{X: 180, Width: 400, IsHole: false},
		},
	}
}

func TestInHole_SetWhenBelowGroundOverHole(t *testing.T) {
	w := holeWorld()

	p := &Player{y16: groundY16 + 1}

	if p.isOverGround(w, 0) {
		t.Fatal("player must be over hole")
	}

	fallingIntoHole := !p.isOverGround(w, 0) && p.y16 > groundY16
	if fallingIntoHole {
		p.isFalling = true
	}

	if !p.isFalling {
		t.Error("isFalling = false, want true")
	}
}

func TestInHole_LandingSkippedAfterFallingInHole(t *testing.T) {
	w := holeWorld()

	p := &Player{y16: groundY16 + 10, isFalling: true}

	if !p.isOverGround(w, 100) {
		t.Fatal("player must be over ground")
	}

	canLand := !p.isFalling && p.isOverGround(w, 100)
	if canLand && p.y16 >= groundY16 {
		p.y16 = groundY16
		p.vy16 = 0
	}

	if p.y16 != groundY16+10 {
		t.Errorf("y16 = %d, want %d", p.y16, groundY16+10)
	}
}

func TestIsFallen(t *testing.T) {
	cases := []struct {
		name         string
		y16          int
		screenHeight int
		want         bool
	}{
		{"y16が画面高さを超えているとき落下と判定する", 401 * 16, 400, true},
		{"y16が画面高さちょうどのとき落下と判定しない", 400 * 16, 400, false},
		{"y16が0のとき落下と判定しない", 0, 400, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := &Player{y16: c.y16}
			got := p.IsFallen(c.screenHeight)
			if got != c.want {
				t.Errorf("IsFallen(%d) = %v, want %v", c.screenHeight, got, c.want)
			}
		})
	}
}
