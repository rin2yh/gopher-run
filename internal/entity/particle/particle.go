package particle

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const MaxParticles = 60

var dirtColors = []color.RGBA{
	{R: 0x8B, G: 0x65, B: 0x13, A: 0xFF},
	{R: 0xA0, G: 0x78, B: 0x30, A: 0xFF},
	{R: 0x6B, G: 0x8C, B: 0x21, A: 0xFF},
}

type Particle struct {
	x, y    float64
	vx, vy  float64
	life    int
	maxLife int
	size    float64
	col     color.RGBA
	op      ebiten.DrawImageOptions
}

func NewImage() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
	return img
}

// SpawnDirt adds dirt particles around (originX, originY) with the given horizontal spread.
func SpawnDirt(ps []Particle, originX, originY float64, spread int) []Particle {
	if len(ps) >= MaxParticles {
		return ps
	}
	count := 2 + rand.Intn(2)
	for i := 0; i < count && len(ps) < MaxParticles; i++ {
		x := originX + float64(rand.Intn(spread+20)) - 10
		y := originY + float64(rand.Intn(10)) - 5
		vx := -3.0 + rand.Float64()*6.0
		vy := -5.0 + rand.Float64()*3.5
		life := 20 + rand.Intn(11)
		size := 2.0 + rand.Float64()*3.0
		col := dirtColors[rand.Intn(len(dirtColors))]
		ps = append(ps, Particle{x: x, y: y, vx: vx, vy: vy, life: life, maxLife: life, size: size, col: col})
	}
	return ps
}

func Update(ps []Particle, cameraShift float64) []Particle {
	n := 0
	for i := range ps {
		ps[i].vy += 0.3
		ps[i].x += ps[i].vx - cameraShift
		ps[i].y += ps[i].vy
		ps[i].life--
		if ps[i].life > 0 {
			if n != i {
				ps[n] = ps[i]
			}
			n++
		}
	}
	return ps[:n]
}

func Draw(screen *ebiten.Image, ps []Particle, img *ebiten.Image) {
	for i := range ps {
		p := &ps[i]
		alpha := float32(p.life) / float32(p.maxLife)
		r := float32(p.col.R) / 255 * alpha
		g := float32(p.col.G) / 255 * alpha
		b := float32(p.col.B) / 255 * alpha
		p.op.GeoM.Reset()
		p.op.ColorScale.Reset()
		p.op.GeoM.Scale(p.size, p.size)
		p.op.GeoM.Translate(p.x, p.y)
		p.op.ColorScale.Scale(r, g, b, alpha)
		screen.DrawImage(img, &p.op)
	}
}
