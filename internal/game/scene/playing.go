package scene

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"gopher-run/internal/entity/enemy"
	"gopher-run/internal/entity/particle"
	"gopher-run/internal/entity/player"
	"gopher-run/internal/entity/popup"
	"gopher-run/internal/input"
	"gopher-run/internal/world"
)

type PlayingScene struct {
	assets         *Assets
	input          *input.Handler
	scorer         *Scorer
	cameraX        int
	player         player.Player
	world          *world.World
	particles      []particle.Particle
	particleImg    *ebiten.Image
	enemies        []enemy.Enemy
	wasOverHole    bool
	popups         []popup.Popup
	popupFaceSmall *ebitentext.GoTextFace
	popupFaceLarge *ebitentext.GoTextFace
}

const (
	cameraSpeedPerFrame = 2
	popupSpawnX         = float64(player.ScreenX + player.Width/2)

	// Eagle 回避（潜る）と Snake 回避（潜らない）は排他操作。プレイヤーが両方へ反応できるよう、
	// 到達フレームの最小間隔を設ける。
	separationFrames = 80.0
)

func NewPlayingScene(assets *Assets, h *input.Handler) *PlayingScene {
	s := &PlayingScene{assets: assets, input: h}
	s.player.Reset()
	s.world = world.New(player.Width, player.ScreenX)
	s.world.Fill(s.cameraX, ScreenWidth)
	s.scorer = NewScorer(s.cameraX)
	s.particles = make([]particle.Particle, 0, particle.MaxParticles)
	s.particleImg = particle.NewImage()
	s.popupFaceSmall = &ebitentext.GoTextFace{Source: assets.FontSource, Size: popup.SmallFontSize}
	s.popupFaceLarge = &ebitentext.GoTextFace{Source: assets.FontSource, Size: popup.LargeFontSize}
	const eagleSpacing = 1200.0
	s.enemies = append(s.enemies, enemy.NewEagleAt(s.SafeEagleSpawnX(enemy.EagleSpawnX), assets.Eagle))
	s.enemies = append(s.enemies, enemy.NewEagleAt(s.SafeEagleSpawnX(enemy.EagleSpawnX+eagleSpacing), assets.Eagle))
	s.enemies = append(s.enemies, enemy.NewSnakeAt(s.SafeSnakeSpawnX(enemy.SnakeSpawnX), assets.Snake))
	return s
}

func shiftPastArrivals(arrival float64, forbidden []float64, separation float64) float64 {
	// 1 回のパスで arrival が複数の f を跨ぐ可能性があるため、最大 len(forbidden)+1 回の再評価が必要。
	for range len(forbidden) + 1 {
		conflict := false
		for _, f := range forbidden {
			if f <= 0 {
				continue
			}
			if math.Abs(arrival-f) < separation {
				arrival = f + separation
				conflict = true
			}
		}
		if !conflict {
			break
		}
	}
	return arrival
}

func arrivalsOf[T enemy.Enemy](enemies []enemy.Enemy, speed float64) []float64 {
	pScreen := float64(player.ScreenX)
	var arrivals []float64
	for _, e := range enemies {
		if _, ok := e.(T); ok {
			arrivals = append(arrivals, (e.X()-pScreen)/speed)
		}
	}
	return arrivals
}

func (s *PlayingScene) SafeSnakeSpawnX(fromX float64) float64 {
	pScreen := float64(player.ScreenX)
	arrival := (fromX - pScreen) / enemy.SnakeSpeedX
	arrival = shiftPastArrivals(arrival, arrivalsOf[*enemy.Eagle](s.enemies, enemy.EagleSpeedX), separationFrames)
	spawnX := arrival*enemy.SnakeSpeedX + pScreen
	if spawnX < fromX {
		spawnX = fromX
	}
	return spawnX
}

// Eagle はプレイヤー位置で穴と衝突すると回避不可のため、Snake 到達シフトを挟んだ前後で穴を避ける。
func (s *PlayingScene) SafeEagleSpawnX(fromX float64) float64 {
	pScreen := float64(player.ScreenX)
	spawnX := fromX

	shiftPastHole := func(x float64) float64 {
		frames := (x - pScreen) / enemy.EagleSpeedX
		worldX := float64(s.cameraX) + frames*cameraSpeedPerFrame + pScreen
		worldX = s.world.ShiftPastHole(worldX, player.Width)
		return (worldX-float64(s.cameraX)-pScreen)*(enemy.EagleSpeedX/cameraSpeedPerFrame) + pScreen
	}

	spawnX = shiftPastHole(spawnX)
	arrival := (spawnX - pScreen) / enemy.EagleSpeedX
	arrival = shiftPastArrivals(arrival, arrivalsOf[*enemy.Snake](s.enemies, enemy.SnakeSpeedX), separationFrames)
	spawnX = arrival*enemy.EagleSpeedX + pScreen
	spawnX = shiftPastHole(spawnX)

	if spawnX < fromX {
		spawnX = fromX
	}
	return spawnX
}

func (s *PlayingScene) NoticeEagleDodged() {
	s.scorer.NoticeEagleDodged()
}

func (s *PlayingScene) SpawnPopup(p popup.Popup) {
	s.popups = popup.Spawn(s.popups, p)
}

func (s *PlayingScene) EagleImage() *ebiten.Image { return s.assets.Eagle }
func (s *PlayingScene) SnakeImage() *ebiten.Image { return s.assets.Snake }

func (s *PlayingScene) Update() Scene {
	s.cameraX += cameraSpeedPerFrame

	s.player.Update(s.world, s.cameraX, s.input)
	if s.player.IsFallen(ScreenHeight) {
		return NewGameOverScene(s.assets, s.input, s.scorer.Value(), s.world, s.player, s.cameraX)
	}

	airborne := s.player.IsAirborne()
	digging := s.player.IsDigging()
	s.scorer.AddDistance(s.cameraX, airborne, digging)

	overGroundNow := s.player.IsOverGround()
	if s.wasOverHole && overGroundNow && !airborne {
		s.scorer.NoticeHoleCleared()
		s.popups = popup.Spawn(s.popups, popup.NewHoleClear(popupSpawnX, float64(s.player.ScreenY()), s.popupFaceSmall))
	}
	s.wasOverHole = !overGroundNow

	if digging {
		s.particles = particle.SpawnDirt(s.particles, player.ScreenX, player.GroundY, player.Width)
	}
	s.particles = particle.Update(s.particles, cameraSpeedPerFrame)
	s.popups = popup.Update(s.popups)

	for i, e := range s.enemies {
		e.Move(s.world, s.cameraX)
		if e.Hit(float64(player.ScreenX), float64(s.player.ScreenY()), player.Width, player.Height, digging) {
			return NewGameOverScene(s.assets, s.input, s.scorer.Value(), s.world, s.player, s.cameraX)
		}
		if e.IsOffScreen(ScreenHeight) {
			e.OnDodged(s, enemy.DodgeContext{
				PopupX:    popupSpawnX,
				PopupY:    float64(s.player.ScreenY()),
				FaceLarge: s.popupFaceLarge,
			})
			s.enemies[i] = e.Respawn(s)
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
		e.Draw(screen)
	}
	s.player.Draw(screen, s.assets.Gopher)
	popup.Draw(screen, s.popups)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", s.scorer.Value()), 10, 10)
}
