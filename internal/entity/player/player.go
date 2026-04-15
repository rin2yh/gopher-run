package player

import (
	"math"

	"gopher-run/internal/input"
	"gopher-run/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GroundY = 320
	ScreenX = 80
	Width   = 60
	Height  = 75
)

const (
	minJumpVY     = -80
	jumpBoost     = 4
	maxJumpFrames = 15
	gravity       = 4
	maxFallVY     = 128
	groundY16     = (GroundY - Height) * 16

	tiltDegrees  = 8.0
	bobAmplitude = 3.0
	bobSpeed     = 0.35

	cx = float64(Width) / 2
	cy = float64(Height) / 2
)

type Player struct {
	y16        int
	vy16       int
	jumpFrames int
	isFalling  bool
	isOnGround bool
	bobFrame   int
	drawOp     ebiten.DrawImageOptions
}

func (p *Player) Reset() {
	p.y16 = groundY16
	p.vy16 = 0
	p.jumpFrames = 0
	p.isFalling = false
	p.isOnGround = true
	p.bobFrame = 0
}

func (p *Player) isOverGround(w *world.World, cameraX int) bool {
	return w.IsGroundAt(ScreenX+cameraX) || w.IsGroundAt(ScreenX+Width-1+cameraX)
}

func (p *Player) Update(w *world.World, cameraX int, h *input.Handler) {
	overGround := p.isOverGround(w, cameraX)

	if h.IsJustPressed() {
		onGround := p.y16 >= groundY16
		if onGround && overGround {
			p.vy16 = minJumpVY
			p.jumpFrames = 1
		}
	}

	if p.jumpFrames > 0 && p.jumpFrames <= maxJumpFrames && h.IsHeld() {
		p.vy16 -= jumpBoost
		p.jumpFrames++
	} else {
		p.jumpFrames = 0
	}

	p.vy16 += gravity
	if p.vy16 > maxFallVY {
		p.vy16 = maxFallVY
	}
	p.y16 += p.vy16

	fallingIntoHole := !overGround && p.y16 > groundY16
	if fallingIntoHole {
		p.isFalling = true
	}

	canLand := !p.isFalling && overGround
	if canLand && p.y16 >= groundY16 {
		p.y16 = groundY16
		p.vy16 = 0
		p.isOnGround = true
		p.bobFrame++
	} else {
		p.bobFrame = 0
		p.isOnGround = false
	}
}

func (p *Player) IsFallen(screenHeight int) bool {
	return p.y16/16 > screenHeight
}

func (p *Player) ScreenY() int {
	return p.y16 / 16
}

func (p *Player) Draw(screen *ebiten.Image, img *ebiten.Image) {
	const tiltAngle = tiltDegrees * math.Pi / 180.0

	bobOffset := 0.0
	if p.isOnGround {
		bobOffset = bobAmplitude * math.Sin(bobSpeed*float64(p.bobFrame))
	}

	p.drawOp.GeoM.Reset()
	p.drawOp.GeoM.Translate(-cx, -cy)
	p.drawOp.GeoM.Rotate(tiltAngle)
	p.drawOp.GeoM.Translate(cx, cy)
	p.drawOp.GeoM.Translate(float64(ScreenX), float64(p.ScreenY())+bobOffset)
	screen.DrawImage(img, &p.drawOp)
}
