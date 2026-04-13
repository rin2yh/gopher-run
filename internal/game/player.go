package game

import (
	"gopher-run/internal/input"
	"gopher-run/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	groundY       = 320
	gopherScreenX = 80
	gopherHeight  = 75
	minJumpVY     = -80
	maxJumpFrames = 15
)

type Player struct {
	y16        int
	vy16       int
	jumpFrames int
}

func (p *Player) Reset() {
	p.y16 = (groundY - gopherHeight) * 16
	p.vy16 = 0
	p.jumpFrames = 0
}

// 後ろ足（左端）基準: 前半分が穴に掛かっていても接地でき、着地も自然に見える
func (p *Player) isOverGround(w *world.World, cameraX int) bool {
	return w.IsGroundAt(gopherScreenX + cameraX)
}

func (p *Player) Update(w *world.World, cameraX int) {
	if input.IsJustPressed() {
		onGround := p.y16 >= (groundY-gopherHeight)*16
		if onGround && p.isOverGround(w, cameraX) {
			p.vy16 = minJumpVY
			p.jumpFrames = 1
		}
	}

	if p.jumpFrames > 0 && p.jumpFrames <= maxJumpFrames && input.IsHeld() {
		p.vy16 -= 4
		p.jumpFrames++
	} else {
		p.jumpFrames = 0
	}

	p.vy16 += 4
	if p.vy16 > 128 {
		p.vy16 = 128
	}
	p.y16 += p.vy16

	if p.isOverGround(w, cameraX) {
		groundY16 := (groundY - gopherHeight) * 16
		if p.y16 >= groundY16 {
			p.y16 = groundY16
			p.vy16 = 0
		}
	}
}

func (p *Player) IsFallen(screenHeight int) bool {
	return p.y16/16 > screenHeight
}

func (p *Player) ScreenY() int {
	return p.y16 / 16
}

func (p *Player) Draw(screen *ebiten.Image, img *ebiten.Image, isTitle bool) {
	gy := p.ScreenY()
	if isTitle {
		gy = groundY - gopherHeight
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gopherScreenX), float64(gy))
	screen.DrawImage(img, op)
}
