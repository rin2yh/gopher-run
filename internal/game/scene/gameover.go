package scene

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"gopher-run/internal/entity/player"
	"gopher-run/internal/input"
	"gopher-run/internal/world"
)

type GameOverScene struct {
	assets  *Assets
	input   *input.Handler
	score   int
	world   *world.World
	player  player.Player
	cameraX int
}

func NewGameOverScene(assets *Assets, h *input.Handler, score int, w *world.World, p player.Player, cameraX int) *GameOverScene {
	return &GameOverScene{
		assets:  assets,
		input:   h,
		score:   score,
		world:   w,
		player:  p,
		cameraX: cameraX,
	}
}

func (s *GameOverScene) Update() Scene {
	if s.input.IsJustPressed() {
		return NewTitleScene(s.assets, s.input)
	}
	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	s.world.Draw(screen, world.DrawParams{
		CameraX:     s.cameraX,
		ScreenWidth: ScreenWidth,
		GroundY:     player.GroundY,
		TileSize:    TileSize,
		FillHeight:  float64(ScreenHeight - player.GroundY - TileSize),
	}, s.assets.GrassTile, s.assets.Dirt)

	s.player.Draw(screen, s.assets.Gopher)
	ebitenutil.DebugPrintAt(screen, "GAME OVER", ScreenWidth/2-35, ScreenHeight/2-30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", s.score), ScreenWidth/2-30, ScreenHeight/2-10)
	ebitenutil.DebugPrintAt(screen, "Press SPACE / Click to restart", ScreenWidth/2-105, ScreenHeight/2+10)
}
