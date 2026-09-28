package plate

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/internal/tables"
)

// FindStem searches a StemPlate for a particular stem across all nine
// palaces. Returns (palace, true) when found.
func FindStem(p *StemPlate, stem almanac.Stem) (uint8, bool) {
	for palace := uint8(1); palace <= 9; palace++ {
		if v, ok := p.Get(palace); ok && v == stem {
			return palace, true
		}
	}
	return 0, false
}

// FindHourStem locates the palace of the current 时辰 heavenly stem on
// the earth plate. 甲 is hidden and maps to the 值符原宫. The center
// palace (5) is returned as-is — callers project it to 2 (坤) where a
// plate rotation needs a ring position.
func FindHourStem(p *StemPlate, stem almanac.Stem, zhiFuOriginalPalace uint8) uint8 {
	if stem == almanac.Jia {
		return zhiFuOriginalPalace
	}
	if palace, ok := FindStem(p, stem); ok {
		return palace
	}
	return zhiFuOriginalPalace
}

// IsStemInTomb reports whether a stem lands in its 入墓 palace.
func IsStemInTomb(stem almanac.Stem, palace uint8) bool {
	return tables.StemTombPalace[stem.Index()] == palace
}
