package qimen

import (
	"errors"
	"testing"
	"time"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

// solarTime is a test helper: fail on err.
func solarTime(t *testing.T, y, mo, d, h, mi, s int) almanac.SolarTime {
	t.Helper()
	st, err := almanac.SolarTimeOf(y, mo, d, h, mi, s)
	if err != nil {
		t.Fatalf("SolarTimeOf: %v", err)
	}
	return st
}

// TestEntryPointsAgree verifies the chart constructors all produce the
// same chart for an equivalent instant.
func TestEntryPointsAgree(t *testing.T) {
	cst := time.FixedZone("CST", 8*3600)
	tt := time.Date(2026, 1, 14, 18, 45, 0, 0, cst)
	st, _ := almanac.SolarTimeFromTime(tt)

	c1 := From(st)
	c2, _ := FromTime(tt)
	c3, _ := FromTimestamp(tt.Unix())

	for _, pair := range [][2]*Chart{{c1, c2}, {c1, c3}} {
		a, b := pair[0], pair[1]
		if a.Year() != b.Year() ||
			a.Month() != b.Month() ||
			a.Day() != b.Day() ||
			a.Hour() != b.Hour() ||
			a.Ju() != b.Ju() {
			t.Errorf("entry points disagree: %+v vs %+v",
				a.SolarTime(), b.SolarTime())
		}
	}
}

// TestSelfPalaceJiaDay covers the 甲 → 值符原宫 fallback path.
// 2025-05-05 is a 甲戌 day, so the day stem is 甲 by construction.
func TestSelfPalaceJiaDay(t *testing.T) {
	c := From(solarTime(t, 2025, 5, 5, 5, 5, 0))
	if got := c.Day().Stem(); got != almanac.Jia {
		t.Fatalf("fixture day stem: got %s, want 甲", got.Name())
	}
	if got, want := c.SelfPalace(), c.ZhiFu().OriginalPalace; got != want {
		t.Errorf("self palace on jia day: got %d, want %d", got, want)
	}
}

// TestJuRule covers the 拆补 option against the 置闰 default on a known
// 超神 instant: on 2024-01-01 the leader day precedes 小寒 (2024-01-06),
// so 置闰 (default) adopts 小寒's row early while 拆补 stays on the
// astronomical 冬至.
func TestJuRule(t *testing.T) {
	st := solarTime(t, 2024, 1, 1, 12, 0, 0)

	zr := From(st)
	if zr.JuRule() != enum.JuRuleZhiRun {
		t.Errorf("default JuRule: got %s, want 置闰", zr.JuRule().Name())
	}
	if zr.Term().Name() != "冬至" {
		t.Errorf("置闰 must keep the astronomical Term(): got %s", zr.Term().Name())
	}
	if zr.JuTerm().Name() != "小寒" || zr.Ju() != 2 {
		t.Errorf("置闰 超神: got %s %d局, want 小寒 2局", zr.JuTerm().Name(), zr.Ju())
	}

	cb := From(st, WithJuRule(enum.JuRuleChaiBu))
	if cb.JuRule() != enum.JuRuleChaiBu {
		t.Errorf("JuRule: got %s, want 拆补", cb.JuRule().Name())
	}
	if cb.JuTerm() != cb.Term() {
		t.Errorf("拆补 JuTerm: got %s, want Term() %s", cb.JuTerm().Name(), cb.Term().Name())
	}
	if cb.Term().Name() != "冬至" || cb.Ju() != 1 {
		t.Errorf("拆补: got %s %d局, want 冬至 1局", cb.Term().Name(), cb.Ju())
	}
	// Pillars are calendar facts — identical under both rules.
	if zr.Month() != cb.Month() || zr.Day() != cb.Day() {
		t.Error("pillars must not depend on JuRule")
	}
}

// TestInvalidTimeErrorChain verifies invalid-time errors from the chart
// entry points match the sentinel under both its qimen and almanac
// names (ErrInvalidTime aliases almanac.ErrInvalidTime).
func TestInvalidTimeErrorChain(t *testing.T) {
	_, err := FromTime(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC))
	if err == nil {
		t.Fatal("expected error for year 10000")
	}
	if !errors.Is(err, ErrInvalidTime) {
		t.Errorf("err %v: want wrapping qimen.ErrInvalidTime", err)
	}
	if !errors.Is(err, almanac.ErrInvalidTime) {
		t.Errorf("err %v: want wrapping almanac.ErrInvalidTime", err)
	}
}

// TestPalaceBounds covers nil-returning out-of-range palace lookups.
func TestPalaceBounds(t *testing.T) {
	c := From(solarTime(t, 2026, 1, 14, 18, 45, 0))
	if c.Palace(0) != nil {
		t.Error("palace 0 should be nil")
	}
	if c.Palace(10) != nil {
		t.Error("palace 10 should be nil")
	}
}

// TestPalacesIter verifies the iter.Seq2 yields all 9 palaces in order.
func TestPalacesIter(t *testing.T) {
	c := From(solarTime(t, 2026, 1, 14, 18, 45, 0))
	count := 0
	for n, p := range c.Palaces() {
		count++
		if p.Number != n {
			t.Errorf("yield mismatch: n=%d, palace=%d", n, p.Number)
		}
	}
	if count != 9 {
		t.Errorf("expected 9 palaces, got %d", count)
	}
}
