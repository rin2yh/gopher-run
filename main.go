package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"gopher-run/internal/input"
	"gopher-run/internal/world"
)

const (
	screenWidth  = 800
	screenHeight = 400
	groundY      = 320
	tileSize     = 32

	gopherScreenX = 80
	gopherHeight  = 75

	minJumpVY     = -80
	maxJumpFrames = 15
)

type Mode int

const (
	ModeTitle   Mode = iota
	ModePlaying
	ModeGameOver
)

type Game struct {
	mode       Mode
	score      int
	cameraX    int
	gopherY16  int
	gopherVY16 int
	jumpFrames int
	world      *world.World
}

var (
	gopherImage    *ebiten.Image
	tilesImage     *ebiten.Image
	dirtImage      *ebiten.Image
	grassTileImage *ebiten.Image
)

func floorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

func floorMod(x, y int) int {
	return x - floorDiv(x, y)*y
}

func init() {
	img, err := png.Decode(bytes.NewReader(gopherPng))
	if err != nil {
		log.Fatal(err)
	}
	gopherImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(tilesPng))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
	grassTileImage = tilesImage.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image)

	dr, dg, db, _ := img.At(tileSize/2, tileSize-4).RGBA()
	dirtImage = ebiten.NewImage(1, 1)
	dirtImage.Fill(color.RGBA{uint8(dr >> 8), uint8(dg >> 8), uint8(db >> 8), 0xFF})
}

func (g *Game) reset() {
	g.score = 0
	g.cameraX = 0
	g.gopherY16 = (groundY - gopherHeight) * 16
	g.gopherVY16 = 0
	g.jumpFrames = 0
	g.world = world.New()
	g.world.Fill(g.cameraX, screenWidth)
}

// 後ろ足（左端）基準: 前半分が穴に掛かっていても接地でき、着地も自然に見える
func (g *Game) isOverGround() bool {
	return g.world.IsGroundAt(gopherScreenX + g.cameraX)
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

		if input.IsJustPressed() {
			onGround := g.gopherY16 >= (groundY-gopherHeight)*16
			if onGround && g.isOverGround() {
				g.gopherVY16 = minJumpVY
				g.jumpFrames = 1
			}
		}

		if g.jumpFrames > 0 && g.jumpFrames <= maxJumpFrames && input.IsHeld() {
			g.gopherVY16 -= 4
			g.jumpFrames++
		} else {
			g.jumpFrames = 0
		}

		g.gopherVY16 += 4
		if g.gopherVY16 > 128 {
			g.gopherVY16 = 128
		}
		g.gopherY16 += g.gopherVY16

		if g.isOverGround() {
			groundY16 := (groundY - gopherHeight) * 16
			if g.gopherY16 >= groundY16 {
				g.gopherY16 = groundY16
				g.gopherVY16 = 0
			}
		}

		if g.gopherY16/16 > screenHeight {
			g.mode = ModeGameOver
		}

		g.world.Prune(g.cameraX)
		g.world.Fill(g.cameraX, screenWidth)

	case ModeGameOver:
		if input.IsJustPressed() {
			g.reset()
			g.mode = ModeTitle
		}
	}
	return nil
}

func (g *Game) drawGround(screen *ebiten.Image) {
	const fillH = float64(screenHeight - groundY - tileSize)

	op := &ebiten.DrawImageOptions{}
	for _, s := range g.world.Segments {
		if s.IsHole {
			continue
		}
		if s.X+s.Width-g.cameraX <= 0 || s.X-g.cameraX >= screenWidth {
			continue
		}

		// ワールド座標でタイル位置を計算し、セグメント境界でクリップ
		firstWorldTileX := floorDiv(s.X, tileSize) * tileSize
		lastWorldTileX := floorDiv(s.X+s.Width-1, tileSize) * tileSize

		for worldTileX := firstWorldTileX; worldTileX <= lastWorldTileX; worldTileX += tileSize {
			screenTileX := worldTileX - g.cameraX
			if screenTileX+tileSize <= 0 || screenTileX >= screenWidth {
				continue
			}

			// セグメント境界でクリップ（ワールド座標）
			clippedLeft := worldTileX
			if clippedLeft < s.X {
				clippedLeft = s.X
			}
			clippedRight := worldTileX + tileSize
			if clippedRight > s.X+s.Width {
				clippedRight = s.X + s.Width
			}
			srcStartX := clippedLeft - worldTileX
			srcEndX := clippedRight - worldTileX
			if srcEndX <= srcStartX {
				continue
			}
			drawX := float64(clippedLeft - g.cameraX)
			tileW := clippedRight - clippedLeft

			op.GeoM.Reset()
			op.GeoM.Translate(drawX, float64(groundY))
			screen.DrawImage(grassTileImage.SubImage(image.Rect(srcStartX, 0, srcEndX, tileSize)).(*ebiten.Image), op)

			op.GeoM.Reset()
			op.GeoM.Scale(float64(tileW), fillH)
			op.GeoM.Translate(drawX, float64(groundY+tileSize))
			screen.DrawImage(dirtImage, op)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	g.drawGround(screen)

	gy := g.gopherY16 / 16
	if g.mode == ModeTitle {
		gy = groundY - gopherHeight
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gopherScreenX), float64(gy))
	screen.DrawImage(gopherImage, op)

	switch g.mode {
	case ModeTitle:
		ebitenutil.DebugPrintAt(screen, "GOPHER RUN", screenWidth/2-40, screenHeight/2-30)
		ebitenutil.DebugPrintAt(screen, "Press SPACE / Click to start", screenWidth/2-100, screenHeight/2-10)
	case ModePlaying:
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score/60), 10, 10)
	case ModeGameOver:
		ebitenutil.DebugPrintAt(screen, "GAME OVER", screenWidth/2-35, screenHeight/2-30)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score/60), screenWidth/2-30, screenHeight/2-10)
		ebitenutil.DebugPrintAt(screen, "Press SPACE / Click to restart", screenWidth/2-105, screenHeight/2+10)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Gopher Run")
	g := &Game{}
	g.reset()
	g.mode = ModeTitle
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
