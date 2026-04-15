package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) IsJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (h *Handler) IsHeld() bool {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return len(ebiten.AppendTouchIDs(nil)) > 0
}

func (h *Handler) IsDigging() bool {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		return true
	}
	return len(ebiten.AppendTouchIDs(nil)) >= 2
}
