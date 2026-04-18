package scene

import "testing"

func TestScorerInitialValue(t *testing.T) {
	s := NewScorer(100)
	if got := s.Value(); got != 0 {
		t.Fatalf("初期スコアは 0 のはず: got=%d", got)
	}
}

func TestScorerAddDistanceAccumulatesDelta(t *testing.T) {
	s := NewScorer(0)
	s.AddDistance(120, false, false)
	want := int(120.0 * baseRatePerPixel)
	if got := s.Value(); got != want {
		t.Fatalf("距離差分加算: want=%d got=%d", want, got)
	}
	s.AddDistance(240, false, false)
	want = int(240.0 * baseRatePerPixel)
	if got := s.Value(); got != want {
		t.Fatalf("lastCameraX 更新後の加算: want=%d got=%d", want, got)
	}
}

func TestScorerAirborneBonus(t *testing.T) {
	s := NewScorer(0)
	s.AddDistance(120, true, false)
	want := 120.0 * baseRatePerPixel * (1 + airborneBonusRate)
	if got := s.points; got != want {
		t.Fatalf("airborne 倍率: want=%v got=%v", want, got)
	}
}

func TestScorerDiggingBonus(t *testing.T) {
	s := NewScorer(0)
	s.AddDistance(120, false, true)
	want := 120.0 * baseRatePerPixel * (1 + diggingBonusRate)
	if got := s.points; got != want {
		t.Fatalf("digging 倍率: want=%v got=%v", want, got)
	}
}

func TestScorerAirborneAndDiggingCombined(t *testing.T) {
	s := NewScorer(0)
	s.AddDistance(120, true, true)
	want := 120.0 * baseRatePerPixel * (1 + airborneBonusRate) * (1 + diggingBonusRate)
	if got := s.points; got != want {
		t.Fatalf("airborne+digging 重畳: want=%v got=%v", want, got)
	}
}

func TestScorerNoticeEagleDodged(t *testing.T) {
	s := NewScorer(0)
	s.NoticeEagleDodged()
	if got := s.Value(); got != eagleDodgeBonus {
		t.Fatalf("Eagle 回避加点: want=%d got=%d", eagleDodgeBonus, got)
	}
}

func TestScorerNoticeHoleCleared(t *testing.T) {
	s := NewScorer(0)
	s.NoticeHoleCleared()
	if got := s.Value(); got != holeClearBonus {
		t.Fatalf("穴越え加点: want=%d got=%d", holeClearBonus, got)
	}
}
