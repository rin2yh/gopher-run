package scene

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"gopher-run/internal/entity/player"
	"gopher-run/internal/input"
	"gopher-run/internal/world"
)

type PlayingScene struct {
	assets  *Assets
	input   *input.Handler
	score   int
	cameraX int
	player  player.Player
	world   *world.World
}

func NewPlayingScene(assets *Assets, h *input.Handler) *PlayingScene {
	s := &PlayingScene{assets: assets, input: h}
	s.player.Reset()
	s.world = world.New(player.Width, player.ScreenX)
	s.world.Fill(s.cameraX, ScreenWidth)
	return s
}

func (s *PlayingScene) Update() Scene {
	s.score++
	s.cameraX += 2

	s.player.Update(s.world, s.cameraX, s.input)
	if s.player.IsFallen(ScreenHeight) {
		return NewGameOverScene(s.assets, s.input, s.score, s.world, s.player, s.cameraX)
	}

	s.world.Prune(s.cameraX)
	s.world.Fill(s.cameraX, ScreenWidth)
	return nil
}

func (s *PlayingScene) Draw(screen *ebiten.Image) {
	s.world.Draw(screen, world.DrawParams{
		CameraX:     s.cameraX,
		ScreenWidth: ScreenWidth,
		GroundY:     player.GroundY,
		TileSize:    TileSize,
		FillHeight:  float64(ScreenHeight - player.GroundY - TileSize),
	}, s.assets.GrassTile, s.assets.Dirt)

	s.player.Draw(screen, s.assets.Gopher)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", s.score/60), 10, 10)
}
