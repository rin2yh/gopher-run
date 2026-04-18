package scene

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"gopher-run/internal/entity/player"
	"gopher-run/internal/input"
)

const (
	titleFontSize      = 68.0
	subtitleFontSize   = 20.0
	controlFontSize    = 15.0
	titleYRatio        = 0.35
	subtitleOffset     = 100.0
	controlsOffset     = 140.0
	controlsLineHeight = 25.0
	blinkCycle         = 60
	blinkVisible       = 40
	shadowOffset       = 3.0
)

var titleColor = color.RGBA{0xFF, 0xE0, 0x00, 0xFF}
var controlColor = color.RGBA{0xCC, 0xCC, 0xCC, 0xFF}

type TitleScene struct {
	assets       *Assets
	input        *input.Handler
	tick         int
	titleFace    *ebitentext.GoTextFace
	subtitleFace *ebitentext.GoTextFace
	controlFace  *ebitentext.GoTextFace
}

func NewTitleScene(assets *Assets, h *input.Handler) *TitleScene {
	return &TitleScene{
		assets:       assets,
		input:        h,
		titleFace:    &ebitentext.GoTextFace{Source: assets.FontSource, Size: titleFontSize},
		subtitleFace: &ebitentext.GoTextFace{Source: assets.FontSource, Size: subtitleFontSize},
		controlFace:  &ebitentext.GoTextFace{Source: assets.FontSource, Size: controlFontSize},
	}
}

func (s *TitleScene) Update() Scene {
	s.tick++
	if s.input.IsJustPressed() {
		return NewPlayingScene(s.assets, s.input)
	}
	return nil
}

func (s *TitleScene) drawGround(screen *ebiten.Image) {
	fillH := float64(ScreenHeight - player.GroundY - TileSize)
	op := &ebiten.DrawImageOptions{}
	for worldTileX := 0; worldTileX < ScreenWidth+400; worldTileX += TileSize {
		w := min(TileSize, ScreenWidth+400-worldTileX)

		op.GeoM.Reset()
		op.GeoM.Translate(float64(worldTileX), float64(player.GroundY))
		screen.DrawImage(s.assets.GrassTile.SubImage(image.Rect(0, 0, w, TileSize)).(*ebiten.Image), op)

		op.GeoM.Reset()
		op.GeoM.Scale(float64(w), fillH)
		op.GeoM.Translate(float64(worldTileX), float64(player.GroundY+TileSize))
		screen.DrawImage(s.assets.Dirt, op)
	}
}

func (s *TitleScene) drawCenteredText(
	screen *ebiten.Image,
	str string,
	face *ebitentext.GoTextFace,
	x, y float64,
	clr color.Color,
	shadow bool,
) {
	if shadow {
		shadowOpts := &ebitentext.DrawOptions{}
		shadowOpts.PrimaryAlign = ebitentext.AlignCenter
		shadowOpts.ColorScale.ScaleWithColor(color.RGBA{0, 0, 0, 200})
		shadowOpts.GeoM.Translate(x+shadowOffset, y+shadowOffset)
		ebitentext.Draw(screen, str, face, shadowOpts)
	}
	opts := &ebitentext.DrawOptions{}
	opts.PrimaryAlign = ebitentext.AlignCenter
	opts.ColorScale.ScaleWithColor(clr)
	opts.GeoM.Translate(x, y)
	ebitentext.Draw(screen, str, face, opts)
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	s.drawGround(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.ScreenX), float64(player.GroundY-player.Height))
	screen.DrawImage(s.assets.Gopher, op)

	titleY := float64(ScreenHeight) * titleYRatio
	cx := float64(ScreenWidth) / 2

	s.drawCenteredText(screen, "Gopher Run", s.titleFace, cx, titleY, titleColor, true)

	if s.tick%blinkCycle < blinkVisible {
		s.drawCenteredText(screen, "Press SPACE / Click to start", s.subtitleFace, cx, titleY+subtitleOffset, color.White, false)
	}

	s.drawCenteredText(screen, "SPACE / Tap : Jump", s.controlFace, cx, titleY+controlsOffset, controlColor, false)
	s.drawCenteredText(screen, "\u2193 / 2-finger : Burrow", s.controlFace, cx, titleY+controlsOffset+controlsLineHeight, controlColor, false)
}
