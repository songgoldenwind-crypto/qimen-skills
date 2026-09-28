package almanac

import (
	"iter"
	"math"
	"sync"
)

// LunarYear is one Chinese lunar year, identified by the Gregorian year
// that contains its 春节 (1st day of 正月).
type LunarYear struct {
	Year      int   // Gregorian year containing 春节
	LeapMonth uint8 // 1..12 if year has a leap month, else 0
}

// LunarYearOf returns metadata for the lunar year corresponding to a
// given Gregorian year.
func LunarYearOf(year int) LunarYear {
	cache := lunarYearCache(year)
	return LunarYear{Year: year, LeapMonth: cache.leapMonth}
}

// Name returns the localized name, e.g. "甲辰年".
func (y LunarYear) Name() string {
	return CycleOf(y.Year - 4).Name()
}

// String implements fmt.Stringer.
func (y LunarYear) String() string { return y.Name() }

// Cycle returns the sexagenary cycle of the lunar year.
func (y LunarYear) Cycle() Cycle { return CycleOf(y.Year - 4) }

// MonthCount returns 12 (normal) or 13 (leap year).
func (y LunarYear) MonthCount() int {
	if y.LeapMonth == 0 {
		return 12
	}
	return 13
}

// Months iterates lunar months in calendar order. Pushes a leap month
// immediately after the month whose number matches y.LeapMonth.
func (y LunarYear) Months() iter.Seq[LunarMonth] {
	return func(yield func(LunarMonth) bool) {
		for m := uint8(1); m <= 12; m++ {
			if !yield(LunarMonth{Year: y, Month: m}) {
				return
			}
			if m == y.LeapMonth {
				if !yield(LunarMonth{Year: y, Month: m, Leap: true}) {
					return
				}
			}
		}
	}
}

// --- internal cache: leap-month + 12-month first-day Julian Days ---

type lunarYearInfo struct {
	leapMonth uint8
	// firstDays[i] = JD offset (from j2000, noon) of 1st day of month
	// at index i (0..11 or 0..12 with leap inserted after leapMonth-1).
	firstDays []int64
	// monthNumbers parallel to firstDays, holds month number (1..12).
	monthNumbers []uint8
	// leapMask[i]=true if month index i is the leap month.
	leapMask []bool
}

var (
	lunarYearCacheMu sync.RWMutex
	lunarYearMap     = map[int]*lunarYearInfo{}
)

func lunarYearCache(year int) *lunarYearInfo {
	lunarYearCacheMu.RLock()
	info, ok := lunarYearMap[year]
	lunarYearCacheMu.RUnlock()
	if ok {
		return info
	}
	lunarYearCacheMu.Lock()
	defer lunarYearCacheMu.Unlock()
	if info, ok = lunarYearMap[year]; ok {
		return info
	}
	info = computeLunarYearInfo(year)
	lunarYearMap[year] = info
	return info
}

// windowMonth is one lunar month inside a solstice-to-solstice window:
// its first-day JD (noon-floored, absolute), month number and leap flag.
type windowMonth struct {
	firstJD int64
	num     uint8
	leap    bool
}

// solsticeWindow resolves the lunar months of one 岁实 window
// [冬至(cycle year Y), 冬至(cycle year Y+1)) per GB/T 33661-2017:
//
//  1. The month containing a 冬至 is always 十一月, so the window opens
//     at 十一月 and closes right before the next 十一月.
//  2. The window spans 12 朔望月 (no leap) or 13 (leap). In a leap
//     window, the first month without a 中气 (even-indexed term) is the
//     leap month and takes the previous month's number.
//  3. Month numbers run 11, 12, 1, 2, ... 10 (leap variants interleaved).
//
// The trailing sentinel 朔 (start of the next window's 十一月) is NOT
// included; callers needing the last month's length peek at the next
// window.
func solsticeWindow(year int) []windowMonth {
	dz1 := math.Floor(TermOf(year, 0).JulianDay() + 0.5)
	dz2 := math.Floor(TermOf(year+1, 0).JulianDay() + 0.5)

	// calcShuo(p) returns the new moon within ±15 days of probe p, so a
	// p+35 probe advances exactly one lunation. Seed strictly before dz1,
	// then walk forward to the 朔 opening the month that contains dz1.
	seed := calcShuo(dz1 - j2000 - 30)
	for {
		next := calcShuo(seed + 35)
		if next+j2000 > dz1 {
			break
		}
		seed = next
	}

	// 朔 sequence: month containing dz1 through the month containing dz2
	// (the latter is the next window's 十一月 and acts as the sentinel).
	shuos := []float64{seed}
	for len(shuos) < 16 {
		next := calcShuo(shuos[len(shuos)-1] + 35)
		if next+j2000 > dz2 {
			break
		}
		shuos = append(shuos, next)
	}
	monthCount := len(shuos) - 1 // 12 or 13
	if monthCount < 12 {
		return nil // astronomical tables out of range
	}

	// Leap detection: first month in the window without a 中气.
	// The window's 中气 are exactly TermOf(year, 0/2/../22) plus dz2.
	leapIdx := -1
	if monthCount == 13 {
		zq := make([]float64, 0, 13)
		for i := 0; i <= 22; i += 2 {
			zq = append(zq, math.Floor(TermOf(year, i).JulianDay()+0.5))
		}
		zq = append(zq, dz2)
		for i := 0; i < monthCount; i++ {
			start := math.Floor(shuos[i] + j2000 + 0.5)
			end := math.Floor(shuos[i+1] + j2000 + 0.5)
			hasZQ := false
			for _, z := range zq {
				if z >= start && z < end {
					hasZQ = true
					break
				}
			}
			if !hasZQ {
				leapIdx = i
				break
			}
		}
	}

	months := make([]windowMonth, monthCount)
	num := uint8(10) // pre-increments to 11 for the first month
	for i := 0; i < monthCount; i++ {
		leap := i == leapIdx
		if !leap {
			if num == 12 {
				num = 1
			} else {
				num++
			}
		}
		months[i] = windowMonth{
			firstJD: int64(math.Floor(shuos[i] + j2000 + 0.5)),
			num:     num,
			leap:    leap,
		}
	}
	return months
}

// computeLunarYearInfo assembles lunar year `year` (正月..十二月 plus a
// possible leap month) from the two solstice windows it straddles:
//
//   - window A = [冬至(year-1), 冬至(year)) supplies 正月..十月 of `year`
//     (months numbered 1..10 there belong to this year by construction);
//   - window B = [冬至(year), 冬至(year+1)) supplies 十一月 and 十二月
//     of `year` (months numbered 11..12 there, including 闰十一月 /
//     闰十二月 — the placement window A alone can never see).
func computeLunarYearInfo(year int) *lunarYearInfo {
	a := solsticeWindow(year)
	b := solsticeWindow(year + 1)
	if a == nil || b == nil {
		return &lunarYearInfo{}
	}

	months := make([]windowMonth, 0, 13)
	for _, m := range a {
		if m.num >= 1 && m.num <= 10 {
			months = append(months, m)
		}
	}
	for _, m := range b {
		if m.num == 11 || m.num == 12 {
			months = append(months, m)
		}
	}

	info := &lunarYearInfo{
		firstDays:    make([]int64, len(months)),
		monthNumbers: make([]uint8, len(months)),
		leapMask:     make([]bool, len(months)),
	}
	for i, m := range months {
		info.firstDays[i] = m.firstJD
		info.monthNumbers[i] = m.num
		info.leapMask[i] = m.leap
		if m.leap {
			info.leapMonth = m.num
		}
	}
	return info
}
