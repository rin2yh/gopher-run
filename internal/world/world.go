package world

import "math/rand"

type Segment struct {
	X      int
	Width  int
	IsHole bool
}

type World struct {
	Segments []Segment
}

func New() *World {
	w := &World{
		Segments: []Segment{
			{X: 0, Width: 400, IsHole: false},
		},
	}
	return w
}

func (w *World) Fill(cameraX, screenWidth int) {
	rightX := 0
	if n := len(w.Segments); n > 0 {
		last := w.Segments[n-1]
		rightX = last.X + last.Width
	}
	for rightX < cameraX+screenWidth+400 {
		seg := Segment{X: rightX}
		if rand.Intn(3) == 0 {
			seg.IsHole = true
			seg.Width = 40 + rand.Intn(21)
		} else {
			seg.Width = 200 + rand.Intn(201)
		}
		w.Segments = append(w.Segments, seg)
		rightX += seg.Width
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

func (w *World) IsGroundAt(worldX int) bool {
	for _, s := range w.Segments {
		if !s.IsHole && worldX >= s.X && worldX < s.X+s.Width {
			return true
		}
	}
	return false
}
