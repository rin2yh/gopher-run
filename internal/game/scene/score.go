package scene

const (
	baseRatePerPixel  = 1.0 / 12.0
	airborneBonusRate = 0.5
	diggingBonusRate  = 0.3
	eagleDodgeBonus   = 10
	holeClearBonus    = 5
)

type Scorer struct {
	points      float64
	lastCameraX int
}

func NewScorer(initialCameraX int) *Scorer {
	return &Scorer{lastCameraX: initialCameraX}
}

func (s *Scorer) AddDistance(cameraX int, airborne, digging bool) {
	delta := float64(cameraX - s.lastCameraX)
	rate := baseRatePerPixel
	if airborne {
		rate *= 1 + airborneBonusRate
	}
	if digging {
		rate *= 1 + diggingBonusRate
	}
	s.points += delta * rate
	s.lastCameraX = cameraX
}

func (s *Scorer) NoticeHoleCleared() {
	s.points += holeClearBonus
}

func (s *Scorer) NoticeEagleDodged() {
	s.points += eagleDodgeBonus
}

func (s *Scorer) Value() int {
	return int(s.points)
}
