package enemy

import (
	"github.com/hajimehoshi/ebiten/v2"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"gopher-run/internal/entity/popup"
)

type GroundChecker interface {
	IsGroundAt(worldX int) bool
}

type DodgeObserver interface {
	NoticeEagleDodged()
	SpawnPopup(p popup.Popup)
}

type DodgeContext struct {
	PopupX    float64
	PopupY    float64
	FaceLarge *ebitentext.GoTextFace
}

type RespawnService interface {
	SafeEagleSpawnX(fromX float64) float64
	SafeSnakeSpawnX(fromX float64) float64
	EagleImage() *ebiten.Image
	SnakeImage() *ebiten.Image
}

type Enemy interface {
	Move(w GroundChecker, cameraX int)
	Hit(px, py, pw, ph float64, digging bool) bool
	Draw(screen *ebiten.Image)
	X() float64
	IsOffScreen(screenHeight int) bool
	OnDodged(obs DodgeObserver, ctx DodgeContext)
	Respawn(svc RespawnService) Enemy
}

func aabbOverlap(px, py, pw, ph, ex, ey, ew, eh float64) bool {
	return px < ex+ew && px+pw > ex && py < ey+eh && py+ph > ey
}
