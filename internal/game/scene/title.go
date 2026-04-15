package scene

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"gopher-run/internal/entity/player"
	"gopher-run/internal/input"
)

type TitleScene struct {
	assets *Assets
	input  *input.Handler
}

func NewTitleScene(assets *Assets, h *input.Handler) *TitleScene {
	return &TitleScene{assets: assets, input: h}
}

func (s *TitleScene) Update() Scene {
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

func (s *TitleScene) Draw(screen *ebiten.Image) {
	s.drawGround(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.ScreenX), float64(player.GroundY-player.Height))
	screen.DrawImage(s.assets.Gopher, op)

	ebitenutil.DebugPrintAt(screen, "GOPHER RUN", ScreenWidth/2-40, ScreenHeight/2-30)
	ebitenutil.DebugPrintAt(screen, "Press SPACE / Click to start", ScreenWidth/2-100, ScreenHeight/2-10)
}
