package almanac

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

// ErrUnsupportedTerm signals an out-of-range solar term index or name.
var ErrUnsupportedTerm = errors.New("almanac: unsupported solar term")

// termNames lists the 24 solar terms in the canonical 寿星 order, starting
// from 冬至 (winter solstice). Index 0..23 maps directly to a Term.Index.
var termNames = [24]string{
	"冬至", "小寒", "大寒", "立春", "雨水", "惊蛰",
	"春分", "清明", "谷雨", "立夏", "小满", "芒种",
	"夏至", "小暑", "大暑", "立秋", "处暑", "白露",
	"秋分", "寒露", "霜降", "立冬", "小雪", "大雪",
}

// Term identifies one of the 24 solar terms keyed by the 寿星 year
// convention: the cycle of year Y starts at its 冬至 (index 0), which
// physically occurs near the end of calendar year Y-1. So TermOf(2025, 0)
// is the winter solstice of 2024-12-21.
// Term is an immutable value type — methods never mutate the receiver.
type Term struct {
	year             int     // 寿星 cycle year (cycle starts at 冬至 of calendar year-1)
	index            uint8   // 0..23 (冬至=0, 大雪=23)
	cursoryJulianDay float64 // calendar-rounded JD (noon precision)
	julianDay        float64 // second-precision JD (cached on construction)
}

// termCache memoizes Term values across calls. Each Term entry already
// holds the expensive julianDay precomputation, so a hit is purely a
// read-locked map lookup — no VSOP / qiAccurate work and no key boxing.
var (
	termCacheMu sync.RWMutex
	termCache   = map[int64]Term{}
)

// TermOf returns the (year, index) solar term. index is normalized into
// 0..23 with floor semantics, so TermOf(y, i) depends only on y*24 + i
// (negative indices and BCE years wrap consistently).
func TermOf(year, index int) Term {
	n := int64(year)*24 + int64(index)
	y := int(n / 24)
	idx := int(n % 24)
	if idx < 0 {
		y--
		idx += 24
	}
	key := int64(y)*24 + int64(idx)

	termCacheMu.RLock()
	t, ok := termCache[key]
	termCacheMu.RUnlock()
	if ok {
		return t
	}

	cjd := initTermJD(y, idx)
	t = Term{
		year:             y,
		index:            uint8(idx),
		cursoryJulianDay: cjd,
		julianDay:        qiAccurate2(cjd) + j2000,
	}
	termCacheMu.Lock()
	termCache[key] = t
	termCacheMu.Unlock()
	return t
}

// TermOfName returns the term in the given year with the given Chinese name.
// Returns ErrUnsupportedTerm when name is unknown.
func TermOfName(year int, name string) (Term, error) {
	for i, n := range termNames {
		if n == name {
			return TermOf(year, i), nil
		}
	}
	return Term{}, fmt.Errorf("%w: %q", ErrUnsupportedTerm, name)
}

// initTermJD computes a coarse Julian Day for the term's start (≈noon)
// based on the calcQi table search.
func initTermJD(year, offset int) float64 {
	jd := math.Floor(float64(year-2000)*365.2422 + 180)
	// 355 ≈ JD offset of 2000-12-21 winter solstice (J2000 = 2451545)
	w := math.Floor((jd-355+183)/365.2422)*365.2422 + 355
	if calcQi(w) > jd {
		w -= 365.2422
	}
	return calcQi(w + 15.2184*float64(offset))
}

// Year returns the 寿星 cycle year of this term. The cycle starts at
// 冬至 (index 0), which falls near the end of calendar year Year()-1;
// terms with index ≥ 3 (立春 onward) lie inside calendar year Year().
func (t Term) Year() int { return t.year }

// Index returns the term's index in canonical order (0=冬至 .. 23=大雪).
func (t Term) Index() int { return int(t.index) }

// Name returns the Chinese name of the term.
func (t Term) Name() string { return termNames[t.index] }

// IsJie reports whether the term is a "节" (odd index — 小寒, 立春, 惊蛰, ...).
func (t Term) IsJie() bool { return t.index%2 == 1 }

// IsQi reports whether the term is a "气" (even index — 冬至, 大寒, 雨水, ...).
func (t Term) IsQi() bool { return t.index%2 == 0 }

// JulianDay returns the second-precision Julian Day of the term's start
// (cached at construction).
func (t Term) JulianDay() float64 { return t.julianDay }

// SolarTime returns the solar instant when this term begins.
func (t Term) SolarTime() SolarTime {
	return solarFromJulianDay(t.julianDay)
}

// DayNumber returns the day ordinal (see [DayNumber]) of the term's
// start, computed directly from the cached Julian Day. Equivalent to
// DayNumber(t.SolarTime()) — including the second rounding that decides
// the 23:00 day roll — without the calendar round-trip.
func (t Term) DayNumber() int {
	d := math.Floor(t.julianDay + 0.5)
	n := int(d) - jdJiaZiDayAnchor
	if math.Round((t.julianDay+0.5-d)*86400) >= 23*3600 {
		n++
	}
	return n
}

// Next returns the term n steps ahead (negative = backward).
func (t Term) Next(n int) Term { return TermOf(t.year, int(t.index)+n) }

// YinYang returns the 阴/阳 遁 polarity of the half-year this term
// belongs to: 冬至..芒种 (index 0..11) → 阳遁, 夏至..大雪 (12..23) → 阴遁.
func (t Term) YinYang() YinYang {
	if t.index < 12 {
		return Yang
	}
	return Yin
}

// String implements fmt.Stringer.
func (t Term) String() string { return t.Name() }

// jdEpsilon is half a second in Julian-Day units. SolarTime ↔ JD
// round-trips lose sub-second precision because SolarTime is integer-
// second; this tolerance lets boundary instants (instants that were
// themselves derived from a Term's SolarTime) still match their term.
const jdEpsilon = 0.5 / 86400

// meanTermDays is the mean interval between adjacent solar terms
// (tropical year / 24), used only to seed the O(1) index estimate.
const meanTermDays = 365.2422 / 24

// TermOfSolarTime returns the solar term whose start ≤ s, in the canonical
// order. This is the "current" term for the supplied instant.
//
// Strategy: TermOf(Y, 0) by convention sits in calendar year (Y-1), so the
// 冬至 of calendar year Y is TermOf(Y+1, 0). Probe that first; otherwise
// s lies inside cycle year Y — estimate the index from the mean term
// interval and correct by at most a step or two. O(1) cache lookups.
func TermOfSolarTime(s SolarTime) Term {
	jd := s.JulianDay()
	year := int(s.Year)
	// 冬至 of current calendar year — last possible term ≤ s.
	if t := TermOf(year+1, 0); t.julianDay <= jd+jdEpsilon {
		return t
	}
	winter := TermOf(year, 0)
	est := int((jd - winter.julianDay) / meanTermDays)
	if est < 0 {
		est = 0
	} else if est > 23 {
		est = 23
	}
	t := TermOf(year, est)
	for t.julianDay > jd+jdEpsilon {
		est--
		t = TermOf(year, est) // est < 0 wraps into the previous cycle year
	}
	for {
		next := TermOf(year, est+1)
		if next.julianDay > jd+jdEpsilon {
			return t
		}
		est++
		t = next
	}
}
