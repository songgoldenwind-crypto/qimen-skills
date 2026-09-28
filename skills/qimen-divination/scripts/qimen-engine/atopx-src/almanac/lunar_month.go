package almanac

import (
	"fmt"
	"iter"
)

// LunarMonth identifies a lunar month within a LunarYear.
type LunarMonth struct {
	Year  LunarYear
	Month uint8 // 1..12
	Leap  bool  // true if this is the leap variant
}

// Name returns the Chinese name (e.g. "正月" or "闰二月").
func (m LunarMonth) Name() string {
	prefix := ""
	if m.Leap {
		prefix = "闰"
	}
	return prefix + lunarMonthDigits[m.Month-1]
}

// String implements fmt.Stringer.
func (m LunarMonth) String() string {
	return fmt.Sprintf("%s%s", m.Year.Name(), m.Name())
}

// DayCount returns the number of days in the month (29 or 30).
func (m LunarMonth) DayCount() int {
	info := lunarYearCache(m.Year.Year)
	for i, num := range info.monthNumbers {
		if num == m.Month && info.leapMask[i] == m.Leap {
			next := i + 1
			if next < len(info.firstDays) {
				return int(info.firstDays[next] - info.firstDays[i])
			}
			// Last month of year — peek into next year's first month.
			nextYear := lunarYearCache(m.Year.Year + 1)
			if len(nextYear.firstDays) > 0 {
				return int(nextYear.firstDays[0] - info.firstDays[i])
			}
			return 29
		}
	}
	return 29
}

// FirstSolarDay returns the solar date of the 1st day of this lunar month.
func (m LunarMonth) FirstSolarDay() SolarTime {
	info := lunarYearCache(m.Year.Year)
	for i, num := range info.monthNumbers {
		if num == m.Month && info.leapMask[i] == m.Leap {
			return solarFromJulianDay(float64(info.firstDays[i]))
		}
	}
	return SolarTime{}
}

// SixtyCycle returns the lunar month sixty-cycle.
// Uses tyme4go-compatible formula: (Year*12 + Month - 47) mod 60.
func (m LunarMonth) SixtyCycle() Cycle {
	idx := m.Year.Year*12 + int(m.Month) - 47
	return CycleOf(idx)
}

// Days iterates lunar days from day 1 to DayCount.
func (m LunarMonth) Days() iter.Seq[LunarDay] {
	return func(yield func(LunarDay) bool) {
		dc := m.DayCount()
		for d := uint8(1); d <= uint8(dc); d++ {
			if !yield(LunarDay{Month: m, Day: d}) {
				return
			}
		}
	}
}

var lunarMonthDigits = [12]string{"正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"}
