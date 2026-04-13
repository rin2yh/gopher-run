package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"gopher-run/internal/input"
	"gopher-run/internal/world"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 400
	TileSize     = 32
)

type Mode int

const (
	ModeTitle   Mode = iota
	ModePlaying
	ModeGameOver
)

type Game struct {
	mode           Mode
	score          int
	cameraX        int
	player         Player
	world          *world.World
	gopherImage    *ebiten.Image
	dirtImage      *ebiten.Image
	grassTileImage *ebiten.Image
}

func New(gopherImg, dirtImg, grassTileImg *ebiten.Image) *Game {
	g := &Game{
		gopherImage:    gopherImg,
		dirtImage:      dirtImg,
		grassTileImage: grassTileImg,
	}
	g.reset()
	g.mode = ModeTitle
	return g
}

func (g *Game) reset() {
	g.score = 0
	g.cameraX = 0
	g.player.Reset()
	g.world = world.New()
	g.world.Fill(g.cameraX, ScreenWidth)
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if input.IsJustPressed() {
			g.mode = ModePlaying
		}

	case ModePlaying:
		g.score++
		g.cameraX += 2

		g.player.Update(g.world, g.cameraX)
		if g.player.IsFallen(ScreenHeight) {
			g.mode = ModeGameOver
		}

		g.world.Prune(g.cameraX)
		g.world.Fill(g.cameraX, ScreenWidth)

	case ModeGameOver:
		if input.IsJustPressed() {
			g.reset()
			g.mode = ModeTitle
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	g.world.Draw(screen, world.DrawParams{
		CameraX:     g.cameraX,
		ScreenWidth: ScreenWidth,
		GroundY:     groundY,
		TileSize:    TileSize,
		FillHeight:  float64(ScreenHeight - groundY - TileSize),
	}, g.grassTileImage, g.dirtImage)

	g.player.Draw(screen, g.gopherImage, g.mode == ModeTitle)

	switch g.mode {
	case ModeTitle:
		ebitenutil.DebugPrintAt(screen, "GOPHER RUN", ScreenWidth/2-40, ScreenHeight/2-30)
		ebitenutil.DebugPrintAt(screen, "Press SPACE / Click to start", ScreenWidth/2-100, ScreenHeight/2-10)
	case ModePlaying:
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score/60), 10, 10)
	case ModeGameOver:
		ebitenutil.DebugPrintAt(screen, "GAME OVER", ScreenWidth/2-35, ScreenHeight/2-30)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score/60), ScreenWidth/2-30, ScreenHeight/2-10)
		ebitenutil.DebugPrintAt(screen, "Press SPACE / Click to restart", ScreenWidth/2-105, ScreenHeight/2+10)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
