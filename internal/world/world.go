package world

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxHoleWidth        = 120 // 最大ジャンプ距離（約136px）から安全マージンを引いた上限
	maxGroundWidth      = 250
	holeWidthMargin     = 20  // Phase3（両足が同時に穴に入るウィンドウ）を21フレーム確保し、偶数cameraXで必ず落下判定が入るよう設定
)

type Segment struct {
	X      int
	Width  int
	IsHole bool
}

type World struct {
	Segments       []Segment
	minHoleWidth   int
	minGroundWidth int
}

func New(playerWidth, playerScreenX int) *World {
	return &World{
		Segments:       []Segment{{X: 0, Width: 400, IsHole: false}},
		minHoleWidth:   playerWidth + holeWidthMargin,
		minGroundWidth: playerWidth + playerScreenX,
	}
}

func (w *World) Fill(cameraX, screenWidth int) {
	rightX := 0
	prevIsHole := false

	if n := len(w.Segments); n > 0 {
		last := w.Segments[n-1]
		rightX = last.X + last.Width
		prevIsHole = last.IsHole
	}

	for rightX < cameraX+screenWidth+400 {
		seg := Segment{X: rightX}

		if prevIsHole || rand.Intn(3) != 0 {
			seg.IsHole = false
			seg.Width = w.minGroundWidth + rand.Intn(maxGroundWidth-w.minGroundWidth+1)
		} else {
			seg.IsHole = true
			seg.Width = w.minHoleWidth + rand.Intn(maxHoleWidth-w.minHoleWidth+1)
		}

		w.Segments = append(w.Segments, seg)
		rightX += seg.Width
		prevIsHole = seg.IsHole
	}
}

func (w *World) Prune(cameraX int) {
	cutoff := cameraX - 200
	i := 0
	for i < len(w.Segments) && w.Segments[i].X+w.Segments[i].Width < cutoff {
		i++
	}
	w.Segments = w.Segments[i:]
}

// ShiftPastHole は worldX から width の範囲が穴と重なる場合、穴の末端まで worldX をずらして返す。
func (w *World) ShiftPastHole(worldX, width float64) float64 {
	for _, s := range w.Segments {
		if !s.IsHole {
			continue
		}
		if worldX < float64(s.X+s.Width) && worldX+width > float64(s.X) {
			worldX = float64(s.X + s.Width)
		}
	}
	return worldX
}

func (w *World) IsGroundAt(worldX int) bool {
	for _, s := range w.Segments {
		if !s.IsHole && worldX >= s.X && worldX < s.X+s.Width {
			return true
		}
	}
	return false
}

func floorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

type DrawParams struct {
	CameraX     int
	ScreenWidth int
	GroundY     int
	TileSize    int
	FillHeight  float64
}

func (w *World) Draw(screen *ebiten.Image, p DrawParams, grassTileImage, dirtImage *ebiten.Image) {
	drawSegments(w.Segments, screen, p, grassTileImage, dirtImage)
}

func drawSegments(segments []Segment, screen *ebiten.Image, p DrawParams, grassTileImage, dirtImage *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	for _, s := range segments {
		if s.IsHole {
			continue
		}
		if s.X+s.Width-p.CameraX <= 0 || s.X-p.CameraX >= p.ScreenWidth {
			continue
		}

		firstWorldTileX := floorDiv(s.X, p.TileSize) * p.TileSize
		lastWorldTileX := floorDiv(s.X+s.Width-1, p.TileSize) * p.TileSize

		for worldTileX := firstWorldTileX; worldTileX <= lastWorldTileX; worldTileX += p.TileSize {
			screenTileX := worldTileX - p.CameraX
			if screenTileX+p.TileSize <= 0 || screenTileX >= p.ScreenWidth {
				continue
			}

			clippedLeft := worldTileX
			if clippedLeft < s.X {
				clippedLeft = s.X
			}
			clippedRight := worldTileX + p.TileSize
			if clippedRight > s.X+s.Width {
				clippedRight = s.X + s.Width
			}
			srcStartX := clippedLeft - worldTileX
			srcEndX := clippedRight - worldTileX
			if srcEndX <= srcStartX {
				continue
			}
			drawX := float64(clippedLeft - p.CameraX)
			tileW := clippedRight - clippedLeft

			op.GeoM.Reset()
			op.GeoM.Translate(drawX, float64(p.GroundY))
			screen.DrawImage(grassTileImage.SubImage(image.Rect(srcStartX, 0, srcEndX, p.TileSize)).(*ebiten.Image), op)

			op.GeoM.Reset()
			op.GeoM.Scale(float64(tileW), p.FillHeight)
			op.GeoM.Translate(drawX, float64(p.GroundY+p.TileSize))
			screen.DrawImage(dirtImage, op)
		}
	}
}
