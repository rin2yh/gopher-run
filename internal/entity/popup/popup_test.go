package popup

import "testing"

func TestNewHoleClear(t *testing.T) {
	p := NewHoleClear(10, 20, nil)
	if p.text != "+5" {
		t.Errorf("text = %q, want %q", p.text, "+5")
	}
	if p.col != holeColor {
		t.Errorf("col = %v, want %v", p.col, holeColor)
	}
	if p.life <= 0 {
		t.Errorf("life = %d, want > 0", p.life)
	}
	if p.maxLife != p.life {
		t.Errorf("maxLife = %d, want %d", p.maxLife, p.life)
	}
	if p.x != 10 || p.y != 20 {
		t.Errorf("pos = (%v,%v), want (10,20)", p.x, p.y)
	}
}

func TestNewEagleDodge(t *testing.T) {
	p := NewEagleDodge(30, 40, nil)
	if p.text != "+10!" {
		t.Errorf("text = %q, want %q", p.text, "+10!")
	}
	if p.col != eagleColor {
		t.Errorf("col = %v, want %v", p.col, eagleColor)
	}
	if p.life <= 0 {
		t.Errorf("life = %d, want > 0", p.life)
	}
}

func TestSpawnRespectsMax(t *testing.T) {
	var ps []Popup
	for range MaxPopups + 3 {
		ps = Spawn(ps, NewHoleClear(0, 0, nil))
	}
	if len(ps) != MaxPopups {
		t.Errorf("len = %d, want %d", len(ps), MaxPopups)
	}
}

func TestUpdateDecrementsLifeAndMovesY(t *testing.T) {
	ps := []Popup{NewHoleClear(0, 100, nil)}
	initLife := ps[0].life
	initY := ps[0].y
	expectedVY := ps[0].vy

	ps = Update(ps)
	if len(ps) != 1 {
		t.Fatalf("len after update = %d, want 1", len(ps))
	}
	if ps[0].life != initLife-1 {
		t.Errorf("life = %d, want %d", ps[0].life, initLife-1)
	}
	if ps[0].y != initY+expectedVY {
		t.Errorf("y = %v, want %v", ps[0].y, initY+expectedVY)
	}
}

func TestUpdateRemovesDeadPopups(t *testing.T) {
	p := NewHoleClear(0, 0, nil)
	p.life = 1
	ps := []Popup{p}
	ps = Update(ps)
	if len(ps) != 0 {
		t.Errorf("len = %d, want 0", len(ps))
	}
}
