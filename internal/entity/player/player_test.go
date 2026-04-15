package player

import (
	"testing"
)

type mockInput struct {
	justPressed bool
	held        bool
	digging     bool
}

func (m *mockInput) IsJustPressed() bool { return m.justPressed }
func (m *mockInput) IsHeld() bool        { return m.held }
func (m *mockInput) IsDigging() bool     { return m.digging }

type mockWorld struct {
	groundXSet map[int]bool
}

func (m *mockWorld) IsGroundAt(worldX int) bool { return m.groundXSet[worldX] }

func solidMock() *mockWorld {
	ground := map[int]bool{}
	for x := 0; x < 600; x++ {
		ground[x] = true
	}
	return &mockWorld{groundXSet: ground}
}

func holeMock() *mockWorld {
	ground := map[int]bool{}
	for x := 0; x < 80; x++ {
		ground[x] = true
	}
	for x := 180; x < 600; x++ {
		ground[x] = true
	}
	return &mockWorld{groundXSet: ground}
}

func TestReset(t *testing.T) {
	p := &Player{state: StateFalling}
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
	if p.state != StateOnGround {
		t.Errorf("state = %v, want StateOnGround", p.state)
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

func TestUpdate_FallsIntoHole(t *testing.T) {
	w := holeMock() // cameraX=0 → ScreenX=80 is in hole
	h := &mockInput{}
	p := &Player{state: StateJumping, y16: groundY16 + 1}
	p.Update(w, 0, h)
	if p.state != StateFalling {
		t.Errorf("state = %v, want StateFalling", p.state)
	}
}

func TestUpdate_LandingSkippedAfterFalling(t *testing.T) {
	w := solidMock()
	h := &mockInput{}
	p := &Player{state: StateFalling, y16: groundY16 + 10, vy16: 8}
	p.Update(w, 0, h)
	if p.state != StateFalling {
		t.Errorf("state = %v, want StateFalling", p.state)
	}
	// 物理演算は継続する（地面にスナップされない）
	if p.y16 == groundY16 {
		t.Error("y16 was snapped to groundY16, want falling to continue")
	}
}

func TestUpdate_DiggingFixesY(t *testing.T) {
	w := solidMock()
	h := &mockInput{digging: true}
	p := &Player{state: StateOnGround, y16: groundY16}
	p.Update(w, 0, h)
	if p.state != StateDigging {
		t.Errorf("state = %v, want StateDigging", p.state)
	}
	if p.y16 != diggingY16 {
		t.Errorf("y16 = %d, want %d (diggingY16)", p.y16, diggingY16)
	}
	if p.vy16 != 0 {
		t.Errorf("vy16 = %d, want 0", p.vy16)
	}
}

func TestUpdate_StopsDiggingOnRelease(t *testing.T) {
	w := solidMock()
	h := &mockInput{digging: false}
	p := &Player{state: StateDigging, y16: diggingY16}
	p.Update(w, 0, h)
	if p.state != StateOnGround {
		t.Errorf("state = %v, want StateOnGround", p.state)
	}
	if p.y16 != groundY16 {
		t.Errorf("y16 = %d, want %d (groundY16)", p.y16, groundY16)
	}
}

func TestUpdate_DiggingOverHoleCausesFall(t *testing.T) {
	w := holeMock() // cameraX=0 → ScreenX=80 is in hole
	h := &mockInput{digging: true}
	p := &Player{state: StateOnGround, y16: groundY16}
	p.Update(w, 0, h)
	if p.state != StateFalling {
		t.Errorf("state = %v, want StateFalling", p.state)
	}
}

func TestUpdate_Jump(t *testing.T) {
	w := solidMock()
	h := &mockInput{justPressed: true, held: true}
	p := &Player{state: StateOnGround, y16: groundY16}
	p.Update(w, 0, h)
	if p.vy16 >= 0 {
		t.Errorf("vy16 = %d, want negative (jumping)", p.vy16)
	}
}

func TestUpdate_LandAfterJump(t *testing.T) {
	w := solidMock()
	h := &mockInput{}
	// 落下中（vy16 > 0）で地面付近
	p := &Player{state: StateJumping, y16: groundY16 - 10, vy16: 20}
	p.Update(w, 0, h)
	if p.state != StateOnGround {
		t.Errorf("state = %v, want StateOnGround", p.state)
	}
	if p.y16 != groundY16 {
		t.Errorf("y16 = %d, want %d", p.y16, groundY16)
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
