package compute

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

// The 月/年家 charts derive their 局 per the fixed year/month rule set
// (validated against authoritative reference charts; both are 阴遁):
//
//   - 年家: 180-year cycle of three 60-year 元 anchored at the 上元
//     甲子 year 1864; every year of an 元 shares its fixed 局 —
//     上元 1, 中元 4, 下元 7.
//   - 月家: the 寅月 (first month) 局 is fixed by the year-branch
//     triad — 子午卯酉 → 8, 辰戌丑未 → 5, 寅申巳亥 → 2 — and retreats
//     by one per month.
//
// The 日家 chart shares the 时家 局 source (节气三元) and needs no
// dedicated derivation here.

// floorMod returns x mod m with the sign of m (m > 0).
func floorMod(x, m int) int { return ((x % m) + m) % m }

// yearJuAnchor is the 上元甲子 year of the current 180-year 三元 cycle.
const yearJuAnchor = 1864

// YearJu returns the 年家 局 (阴遁; fixed per 60-year 元: 上1 中4 下7)
// and 元 for the pillar year containing the given term (立春 boundary,
// per the year pillar).
func YearJu(term almanac.Term) (uint8, enum.Yuan) {
	y := term.Year()
	if term.Index() < 3 {
		y--
	}
	n := y - yearJuAnchor
	el := floorMod((n-floorMod(n, 60))/60, 3)
	return [3]uint8{1, 4, 7}[el], enum.Yuan(el)
}

// MonthJu returns the 月家 局 (1..9, 阴遁) and 元 from the year branch
// and the month pillar's branch (寅 = first month).
func MonthJu(yearBranch, monthBranch almanac.Branch) (uint8, enum.Yuan) {
	// Triads share branch index mod 3: 子午卯酉 → 0, 辰戌丑未 → 1,
	// 寅申巳亥 → 2; their 寅月 leading numbers are 8 / 5 / 2 and the
	// classical 元 assignment is 上 / 下 / 中 respectively.
	triad := yearBranch.Index() % 3
	leading := [3]int{8, 5, 2}[triad]
	yuan := [3]enum.Yuan{enum.YuanUpper, enum.YuanLower, enum.YuanMiddle}[triad]
	monthSeq := floorMod(monthBranch.Index()-int(almanac.Yin_), 12)
	return uint8(floorMod(leading-1-monthSeq, 9) + 1), yuan
}
