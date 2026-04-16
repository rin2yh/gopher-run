package scene

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"gopher-run/internal/entity/enemy"
	"gopher-run/internal/entity/particle"
	"gopher-run/internal/entity/player"
	"gopher-run/internal/input"
	"gopher-run/internal/world"
)

type PlayingScene struct {
	assets      *Assets
	input       *input.Handler
	score       int
	cameraX     int
	player      player.Player
	world       *world.World
	particles   []particle.Particle
	particleImg *ebiten.Image
	enemies     []enemy.Enemy
}

const cameraSpeedPerFrame = 2

func NewPlayingScene(assets *Assets, h *input.Handler) *PlayingScene {
	s := &PlayingScene{assets: assets, input: h}
	s.player.Reset()
	s.world = world.New(player.Width, player.ScreenX)
	s.world.Fill(s.cameraX, ScreenWidth)
	s.particles = make([]particle.Particle, 0, particle.MaxParticles)
	s.particleImg = particle.NewImage()
	const eagleSpacing = 1200.0
	s.enemies = []enemy.Enemy{
		enemy.NewEagleAt(s.safeEagleSpawnX(enemy.EagleSpawnX)),
		enemy.NewEagleAt(s.safeEagleSpawnX(enemy.EagleSpawnX + eagleSpacing)),
	}
	return s
}

// safeEagleSpawnX は fromX を起点に、Eagle がプレイヤー位置へ到達するタイミングと
// 穴が重ならないよう spawn X を調整して返す。
func (s *PlayingScene) safeEagleSpawnX(fromX float64) float64 {
	pScreen := float64(player.ScreenX)

	frames := (fromX - pScreen) / enemy.EagleSpeedX
	targetWorldX := float64(s.cameraX) + frames*cameraSpeedPerFrame + pScreen
	targetWorldX = s.world.ShiftPastHole(targetWorldX, player.Width)

	spawnX := (targetWorldX-float64(s.cameraX)-pScreen)*(enemy.EagleSpeedX/cameraSpeedPerFrame) + pScreen
	if spawnX < fromX {
		spawnX = fromX
	}
	return spawnX
}

func (s *PlayingScene) Update() Scene {
	s.score++
	s.cameraX += cameraSpeedPerFrame

	s.player.Update(s.world, s.cameraX, s.input)
	if s.player.IsFallen(ScreenHeight) {
		return NewGameOverScene(s.assets, s.input, s.score, s.world, s.player, s.cameraX)
	}

	if s.player.IsDigging() {
		s.particles = particle.SpawnDirt(s.particles, player.ScreenX, player.GroundY, player.Width)
	}
	s.particles = particle.Update(s.particles, cameraSpeedPerFrame)

	for i, e := range s.enemies {
		e.Move()
		if e.Hit(float64(player.ScreenX), float64(s.player.ScreenY()), player.Width, player.Height) {
			return NewGameOverScene(s.assets, s.input, s.score, s.world, s.player, s.cameraX)
		}
		if e.X() < 0 {
			s.enemies[i] = enemy.NewEagleAt(s.safeEagleSpawnX(enemy.EagleSpawnX))
		}
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

	particle.Draw(screen, s.particles, s.particleImg)
	for _, e := range s.enemies {
		e.Draw(screen, s.assets.Eagle)
	}
	s.player.Draw(screen, s.assets.Gopher)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", s.score/60), 10, 10)
}
