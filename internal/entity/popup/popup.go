package popup

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

const MaxPopups = 8

const (
	popupLife     = 50
	popupVY       = -1.0
	shadowOffset  = 2.0
	SmallFontSize = 16.0
	LargeFontSize = 22.0
)

var (
	holeColor  = color.RGBA{R: 0x6B, G: 0xD0, B: 0x50, A: 0xFF}
	eagleColor = color.RGBA{R: 0xFF, G: 0xE0, B: 0x00, A: 0xFF}
)

type Popup struct {
	x, y    float64
	vy      float64
	life    int
	maxLife int
	text    string
	col     color.RGBA
	face    *ebitentext.GoTextFace
}

func NewHoleClear(x, y float64, face *ebitentext.GoTextFace) Popup {
	return Popup{
		x: x, y: y, vy: popupVY,
		life: popupLife, maxLife: popupLife,
		text: "+5", col: holeColor, face: face,
	}
}

func NewEagleDodge(x, y float64, face *ebitentext.GoTextFace) Popup {
	return Popup{
		x: x, y: y, vy: popupVY,
		life: popupLife, maxLife: popupLife,
		text: "+10!", col: eagleColor, face: face,
	}
}

func Spawn(ps []Popup, p Popup) []Popup {
	if len(ps) >= MaxPopups {
		return ps
	}
	return append(ps, p)
}

func Update(ps []Popup) []Popup {
	n := 0
	for i := range ps {
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

func Draw(screen *ebiten.Image, ps []Popup) {
	for i := range ps {
		p := &ps[i]
		alpha := float32(p.life) / float32(p.maxLife)

		shadow := &ebitentext.DrawOptions{}
		shadow.PrimaryAlign = ebitentext.AlignCenter
		shadow.ColorScale.ScaleWithColor(color.RGBA{0, 0, 0, 200})
		shadow.ColorScale.ScaleAlpha(alpha)
		shadow.GeoM.Translate(p.x+shadowOffset, p.y+shadowOffset)
		ebitentext.Draw(screen, p.text, p.face, shadow)

		opts := &ebitentext.DrawOptions{}
		opts.PrimaryAlign = ebitentext.AlignCenter
		opts.ColorScale.ScaleWithColor(p.col)
		opts.ColorScale.ScaleAlpha(alpha)
		opts.GeoM.Translate(p.x, p.y)
		ebitentext.Draw(screen, p.text, p.face, opts)
	}
}
