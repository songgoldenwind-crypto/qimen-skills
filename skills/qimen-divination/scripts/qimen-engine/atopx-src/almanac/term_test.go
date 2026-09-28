package almanac

import (
	"fmt"
	"testing"
)

// TestJulianDayRoundTrip verifies SolarTime ↔ JulianDay conversion is
// lossless across a representative range.
func TestJulianDayRoundTrip(t *testing.T) {
	cases := []struct {
		y, m, d, h, mi, s int
	}{
		{1900, 1, 1, 0, 0, 0},
		{1900, 6, 15, 12, 30, 45},
		{2000, 1, 1, 12, 0, 0},
		{2025, 5, 5, 5, 5, 0},
		{2050, 12, 31, 23, 59, 59},
	}
	for _, c := range cases {
		st, err := SolarTimeOf(c.y, c.m, c.d, c.h, c.mi, c.s)
		if err != nil {
			t.Fatalf("SolarTimeOf(%v): %v", c, err)
		}
		jd := st.JulianDay()
		back := solarFromJulianDay(jd)
		if back != st {
			t.Errorf("round-trip mismatch: %v → JD %.9f → %v", st, jd, back)
		}
	}
}

// TestKnownTerms anchors a handful of well-known 节气 to their reference
// times. Values verified against published 中国科学院紫金山天文台 ephemerides.
func TestKnownTerms(t *testing.T) {
	cases := []struct {
		year, idx int
		want      string // SolarTime formatted "YYYY-MM-DD HH:MM:SS"
	}{
		{2025, 0, "2024-12-21 17:20:34"},  // 冬至 of 2024
		{2025, 3, "2025-02-03 22:10:13"},  // 立春 of 2025
		{2025, 12, "2025-06-21 10:42:00"}, // 夏至 of 2025
	}
	for _, c := range cases {
		got := TermOf(c.year, c.idx).SolarTime().String()
		// Allow ±60s tolerance (different ephemerides differ at sub-minute level)
		if !nearTime(t, got, c.want, 60) {
			t.Errorf("TermOf(%d, %d): got %s, want %s (±60s)", c.year, c.idx, got, c.want)
		}
	}
}

// TestTermBoundary covers term lookup at the exact start instant.
func TestTermBoundary(t *testing.T) {
	winter := TermOf(2026, 0)
	want := winter.SolarTime()
	got := TermOfSolarTime(want)
	if got.Index() != winter.Index() || got.Year() != winter.Year() {
		t.Errorf("at winter solstice start: got %v/%d, want %v/%d",
			got.Name(), got.Year(), winter.Name(), winter.Year())
	}
	prev := TermOfSolarTime(want.AddSeconds(-1))
	if prev.Index() == winter.Index() && prev.Year() == winter.Year() {
		t.Errorf("1s before should be prev term, got same")
	}
}

// TestTermIndexWrap verifies (year, index) normalization uses floor
// division so negative inputs (including BCE years) wrap consistently.
// The invariant: TermOf(y, i) depends only on y*24 + i.
func TestTermIndexWrap(t *testing.T) {
	cases := []struct {
		year, index         int
		wantYear, wantIndex int
	}{
		{2025, -1, 2024, 23},
		{2025, 24, 2026, 0},
		{2025, 25, 2026, 1},
		{-100, -1, -101, 23},
		{-100, 24, -99, 0},
		{-1, -25, -3, 23},
	}
	for _, c := range cases {
		got := TermOf(c.year, c.index)
		if got.Year() != c.wantYear || got.Index() != c.wantIndex {
			t.Errorf("TermOf(%d, %d): got (%d, %d), want (%d, %d)",
				c.year, c.index, got.Year(), got.Index(), c.wantYear, c.wantIndex)
		}
	}
	// Next must agree with direct construction across the wrap.
	for _, year := range []int{2025, -100} {
		prev := TermOf(year, 0).Next(-1)
		want := TermOf(year, -1)
		if prev.Year() != want.Year() || prev.Index() != want.Index() {
			t.Errorf("TermOf(%d,0).Next(-1) = (%d,%d), want (%d,%d)",
				year, prev.Year(), prev.Index(), want.Year(), want.Index())
		}
	}
}

// nearTime parses two "YYYY-MM-DD HH:MM:SS" strings and reports whether
// they are within tolSeconds of each other.
func nearTime(t *testing.T, a, b string, tolSeconds int) bool {
	t.Helper()
	var ay, am, ad, ah, ami, as int
	var by, bm, bd, bh, bmi, bs int
	if _, err := parseStamp(a, &ay, &am, &ad, &ah, &ami, &as); err != nil {
		t.Fatalf("parse %q: %v", a, err)
	}
	if _, err := parseStamp(b, &by, &bm, &bd, &bh, &bmi, &bs); err != nil {
		t.Fatalf("parse %q: %v", b, err)
	}
	st1, _ := SolarTimeOf(ay, am, ad, ah, ami, as)
	st2, _ := SolarTimeOf(by, bm, bd, bh, bmi, bs)
	diff := int(st1.JulianDay()*86400 - st2.JulianDay()*86400)
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolSeconds
}

func parseStamp(s string, y, m, d, h, mi, sec *int) (int, error) {
	return fmt.Sscanf(s, "%d-%d-%d %d:%d:%d", y, m, d, h, mi, sec)
}
