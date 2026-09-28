package almanac

import "math"

// jdJiaZiDayAnchor is the noon Julian Day of 2000-01-07 UT — a 甲子 day
// (j2000 + 6). Day-pillar cycles are counted from this anchor.
const jdJiaZiDayAnchor = j2000 + 6

// Pillars holds the four sexagenary pillars (year/month/day/hour)
// derived from a solar instant using term-direct rules:
//
//   - Year:  立春 (term index 3) marks the year switch. (Y - 4) mod 60
//     gives 甲子 numbering from year 4 AD.
//   - Month: Each 节 (odd term) starts a new month. 五虎遁年起月 gives
//     the stem of 寅月 from the year stem.
//   - Day:   2000-01-07 noon UT = 甲子日 (JD 2451551). Hour 23 rolls
//     forward by 1 day.
//   - Hour:  五鼠遁日起时. Two-hour 时辰 from 23-01=子.
type Pillars struct {
	Year, Month, Day, Hour Cycle
}

// PillarsOf computes the four pillars for a solar instant.
func PillarsOf(s SolarTime) Pillars { return PillarsAt(s, TermOfSolarTime(s)) }

// PillarsAt computes the four pillars for a solar instant whose current
// solar term is already known. Precondition: term == s.Term(). Callers
// that already resolved the term (e.g. chart construction) avoid a
// second term lookup this way.
func PillarsAt(s SolarTime, term Term) Pillars {
	yearCycle := yearPillar(term)
	dayCycle := dayPillar(s)
	return Pillars{
		Year:  yearCycle,
		Month: monthPillar(term, yearCycle),
		Day:   dayCycle,
		Hour:  hourPillarFromDay(s, dayCycle),
	}
}

// Pillars returns the four pillars for s. Convenience wrapper.
func (s SolarTime) Pillars() Pillars { return PillarsOf(s) }

// Term returns the current solar term for s.
func (s SolarTime) Term() Term { return TermOfSolarTime(s) }

// yearPillar uses 立春 (term index 3) as the year boundary, derived
// directly from the current term: terms 冬至/小寒/大寒 (index < 3) of
// cycle year Y lie before 立春, so the pillar year is Y-1; from 立春
// onward it is Y. Pillar cycle is (year - 4) mod 60 (甲子 = AD 4).
func yearPillar(term Term) Cycle {
	y := term.Year()
	if term.Index() < 3 {
		y--
	}
	return CycleOf(y - 4)
}

// monthPillar uses 节气-based month boundaries (节 = odd-index terms).
// Month sequence starts at 寅 (after 立春): each pair of adjacent terms
// (one 节 + one 气) is one month, so seq = ((termIdx - 3 + 24) % 24) / 2.
// 五虎遁: month stem of 寅月 = (year_stem * 2 + 2) mod 10, hence the
// sixty-cycle index at seq 0 = ((year_stem mod 5) * 12 + 2) mod 60.
func monthPillar(term Term, year Cycle) Cycle {
	seq := ((term.Index() - 3 + 24) % 24) / 2
	yearStem := int(year.Stem())
	idx0 := ((yearStem%5)*12 + 2) % 60
	return CycleOf(idx0 + seq)
}

// dayCycleAtNoon returns the day-pillar cycle for the calendar day
// containing the given date (anchored at jdJiaZiDayAnchor).
func dayCycleAtNoon(year, month, day int) Cycle {
	noonJD := julianDayFromYmdhms(year, month, day, 12, 0, 0)
	return CycleOf(int(math.Floor(noonJD+0.5)) - jdJiaZiDayAnchor)
}

// DayNumber returns the continuous day ordinal aligned with the day
// pillar: day 0 is 2000-01-07 (a 甲子 day), and hour 23 rolls to the
// next day (晚子时 convention), so CycleOf(DayNumber(s)) always equals
// the day pillar of s. Useful for day-grid arithmetic such as 符头
// (15-day leader) alignment.
func DayNumber(s SolarTime) int {
	noonJD := julianDayFromYmdhms(int(s.Year), int(s.Month), int(s.Day), 12, 0, 0)
	n := int(math.Floor(noonJD+0.5)) - jdJiaZiDayAnchor
	if s.Hour == 23 {
		n++
	}
	return n
}

// dayPillar derives the day pillar; hour 23 rolls to the next day
// (晚子时 convention: the day switches at 23:00).
func dayPillar(s SolarTime) Cycle { return CycleOf(DayNumber(s)) }

// hourPillarFromDay derives the hour pillar from the day pillar.
// At day stem ds, 子时 (hour branch 0) has stem (ds*2) mod 10, so the
// hour-pillar sixty-cycle index at 子时 is ((ds mod 5) * 12) mod 60.
// Each subsequent 时辰 adds 1 to the cycle index.
func hourPillarFromDay(s SolarTime, day Cycle) Cycle {
	branchIdx := hourBranchIndex(int(s.Hour))
	dayStem := int(day.Stem())
	idx0 := ((dayStem % 5) * 12) % 60
	return CycleOf(idx0 + branchIdx)
}

// hourBranchIndex maps 0..23 hours to 0..11 branches.
// 23, 0 → 0 (子); 1, 2 → 1 (丑); ...; 21, 22 → 11 (亥).
func hourBranchIndex(h int) int {
	return ((h + 1) / 2) % 12
}
