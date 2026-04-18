package player

import (
	"math"

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
	diggingY16    = (GroundY + 10) * 16

	tiltDegrees  = 8.0
	tiltAngle    = tiltDegrees * math.Pi / 180.0
	bobAmplitude = 3.0
	bobSpeed     = 0.35
	diggingAngle = math.Pi / 2

	cx = float64(Width) / 2
	cy = float64(Height) / 2
)

type GroundChecker interface {
	IsGroundAt(worldX int) bool
}

type InputReader interface {
	IsJustPressed() bool
	IsHeld() bool
	IsDigging() bool
}

type PlayerState int

const (
	StateOnGround PlayerState = iota
	StateJumping
	StateDigging
	StateFalling
)

type Player struct {
	y16        int
	vy16       int
	jumpFrames int
	state      PlayerState
	overGround bool
	bobFrame   int
	drawOp     ebiten.DrawImageOptions
}

func (p *Player) Reset() {
	p.y16 = groundY16
	p.vy16 = 0
	p.jumpFrames = 0
	p.state = StateOnGround
	p.overGround = true
	p.bobFrame = 0
}

func (p *Player) isOverGround(w GroundChecker, cameraX int) bool {
	return w.IsGroundAt(ScreenX+cameraX) || w.IsGroundAt(ScreenX+Width-1+cameraX)
}

func (p *Player) processDigging(h InputReader, overGround bool) {
	if p.state == StateDigging {
		if !h.IsDigging() {
			p.state = StateOnGround
			return
		}
		if !overGround {
			p.state = StateFalling
			return
		}
		p.y16 = diggingY16
		p.vy16 = 0
		return
	}

	if h.IsDigging() && p.state == StateOnGround {
		if overGround {
			p.state = StateDigging
			p.y16 = diggingY16
			p.vy16 = 0
		} else {
			p.state = StateFalling
		}
	}
}

func (p *Player) processJump(h InputReader, overGround bool) {
	if h.IsJustPressed() {
		onGround := p.y16 >= groundY16
		if onGround && overGround {
			p.vy16 = minJumpVY
			p.jumpFrames = 1
			p.state = StateJumping
		}
	}

	if p.jumpFrames > 0 && p.jumpFrames <= maxJumpFrames && h.IsHeld() {
		p.vy16 -= jumpBoost
		p.jumpFrames++
	} else {
		p.jumpFrames = 0
	}
}

func (p *Player) applyPhysics() {
	p.vy16 += gravity
	if p.vy16 > maxFallVY {
		p.vy16 = maxFallVY
	}
	p.y16 += p.vy16
}

func (p *Player) updateGroundState(overGround bool) {
	if !overGround && p.y16 > groundY16 {
		p.state = StateFalling
		return
	}

	if overGround && p.y16 >= groundY16 {
		p.y16 = groundY16
		p.vy16 = 0
		p.state = StateOnGround
		p.bobFrame++
	} else {
		p.bobFrame = 0
		p.state = StateJumping
	}
}

func (p *Player) Update(w GroundChecker, cameraX int, h InputReader) {
	overGround := p.isOverGround(w, cameraX)
	p.overGround = overGround
	p.processDigging(h, overGround)
	if p.state == StateDigging {
		return
	}
	if p.state != StateFalling {
		p.processJump(h, overGround)
	}
	p.applyPhysics()
	if p.state != StateFalling {
		p.updateGroundState(overGround)
	}
}

func (p *Player) IsDigging() bool {
	return p.state == StateDigging
}

func (p *Player) IsAirborne() bool {
	return p.state == StateJumping || p.state == StateFalling
}

func (p *Player) IsOverGround() bool {
	return p.overGround
}

func (p *Player) IsFallen(screenHeight int) bool {
	return p.y16/16 > screenHeight
}

func (p *Player) ScreenY() int {
	return p.y16 / 16
}

func (p *Player) Draw(screen *ebiten.Image, img *ebiten.Image) {
	angle := tiltAngle
	bobOffset := 0.0

	switch p.state {
	case StateDigging:
		angle = diggingAngle
	case StateOnGround:
		bobOffset = bobAmplitude * math.Sin(bobSpeed*float64(p.bobFrame))
	}

	p.drawOp.GeoM.Reset()
	p.drawOp.GeoM.Translate(-cx, -cy)
	p.drawOp.GeoM.Rotate(angle)
	p.drawOp.GeoM.Translate(cx, cy)
	p.drawOp.GeoM.Translate(float64(ScreenX), float64(p.ScreenY())+bobOffset)
	screen.DrawImage(img, &p.drawOp)
}
