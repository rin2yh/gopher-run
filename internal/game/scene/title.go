package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"gopher-run/internal/entity/player"
	"gopher-run/internal/input"
	"gopher-run/internal/world"
)

type TitleScene struct {
	assets *Assets
	input  *input.Handler
	world  *world.World
}

func NewTitleScene(assets *Assets, h *input.Handler) *TitleScene {
	return &TitleScene{assets: assets, input: h, world: world.NewFlat(ScreenWidth + 400)}
}

func (s *TitleScene) Update() Scene {
	if s.input.IsJustPressed() {
		return NewPlayingScene(s.assets, s.input)
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	s.world.Draw(screen, world.DrawParams{
		CameraX:     0,
		ScreenWidth: ScreenWidth,
		GroundY:     player.GroundY,
		TileSize:    TileSize,
		FillHeight:  float64(ScreenHeight - player.GroundY - TileSize),
	}, s.assets.GrassTile, s.assets.Dirt)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.ScreenX), float64(player.GroundY-player.Height))
	screen.DrawImage(s.assets.Gopher, op)

	ebitenutil.DebugPrintAt(screen, "GOPHER RUN", ScreenWidth/2-40, ScreenHeight/2-30)
	ebitenutil.DebugPrintAt(screen, "Press SPACE / Click to start", ScreenWidth/2-100, ScreenHeight/2-10)
}
