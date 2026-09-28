package almanac

import (
	"fmt"
	"math"
)

// LunarDay identifies a day within a lunar month.
//
// The zero value (Day == 0) is the "not found" sentinel returned by
// LunarDayOf for instants outside the computable range; Name() and
// String() return "" for it.
type LunarDay struct {
	Month LunarMonth
	Day   uint8 // 1..30; 0 = invalid/unresolved
}

// LunarDayOf maps a solar instant to its lunar date.
//
// Resolution: solar instant → noon-JD-of-its-day → 朔 ≤ noon-JD → containing
// lunar month → day offset. Returns the zero LunarDay when the instant
// cannot be located (astronomical tables out of range).
func LunarDayOf(s SolarTime) LunarDay {
	// Solar day at noon
	jd := julianDayFromYmdhms(int(s.Year), int(s.Month), int(s.Day), 12, 0, 0)
	noonJD := math.Floor(jd + 0.5) // integer JD of noon

	// Scan candidate lunar years (current and adjacent) for the month containing noonJD.
	for offset := -1; offset <= 1; offset++ {
		year := int(s.Year) + offset
		info := lunarYearCache(year)
		for i, fd := range info.firstDays {
			next := int64(math.MaxInt32)
			if i+1 < len(info.firstDays) {
				next = info.firstDays[i+1]
			} else {
				ny := lunarYearCache(year + 1)
				if len(ny.firstDays) > 0 {
					next = ny.firstDays[0]
				}
			}
			if int64(noonJD) >= fd && int64(noonJD) < next {
				dayNum := uint8(int64(noonJD) - fd + 1)
				return LunarDay{
					Month: LunarMonth{
						Year:  LunarYear{Year: year, LeapMonth: info.leapMonth},
						Month: info.monthNumbers[i],
						Leap:  info.leapMask[i],
					},
					Day: dayNum,
				}
			}
		}
	}
	return LunarDay{}
}

// LunarDay returns the lunar date corresponding to the solar time.
func (s SolarTime) LunarDay() LunarDay { return LunarDayOf(s) }

// Year returns the parent lunar year.
func (d LunarDay) Year() LunarYear { return d.Month.Year }

// SolarDay returns the solar date of this lunar day (00:00:00).
func (d LunarDay) SolarDay() SolarTime {
	first := d.Month.FirstSolarDay()
	if first == (SolarTime{}) {
		return SolarTime{}
	}
	t := first.ToTime(nil).AddDate(0, 0, int(d.Day)-1)
	out, _ := SolarTimeFromTime(t)
	return out
}

// Cycle returns the day-pillar sixty cycle for the lunar day's solar
// equivalent at noon (matches the day pillar in PillarsOf for hour<23).
func (d LunarDay) Cycle() Cycle {
	s := d.SolarDay()
	if s == (SolarTime{}) {
		return CycleOf(0)
	}
	return dayCycleAtNoon(int(s.Year), int(s.Month), int(s.Day))
}

// Name returns the day's Chinese name (e.g. "初一", "十五", "廿三"),
// or "" for the zero (unresolved) LunarDay.
func (d LunarDay) Name() string {
	if d.Day == 0 {
		return ""
	}
	return lunarDayNames[d.Day-1]
}

// String returns "<year><month><day>", e.g. "农历甲辰年正月初一",
// or "" for the zero (unresolved) LunarDay.
func (d LunarDay) String() string {
	if d.Day == 0 {
		return ""
	}
	return fmt.Sprintf("%s%s%s", d.Month.Year.Name(), d.Month.Name(), d.Name())
}

var lunarDayNames = [30]string{
	"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十",
	"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
	"廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十",
}
