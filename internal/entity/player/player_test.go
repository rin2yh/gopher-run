package player

import (
	"testing"
)

func TestReset(t *testing.T) {
	p := &Player{}
	p.Reset()

	wantY16 := (GroundY - Height) * 16
	if p.y16 != wantY16 {
		t.Errorf("y16 = %d, want %d", p.y16, wantY16)
	}
	if p.vy16 != 0 {
		t.Errorf("vy16 = %d, want 0", p.vy16)
	}
	if p.jumpFrames != 0 {
		t.Errorf("jumpFrames = %d, want 0", p.jumpFrames)
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
