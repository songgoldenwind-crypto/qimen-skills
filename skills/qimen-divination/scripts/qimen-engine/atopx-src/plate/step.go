package plate

import "github.com/atopx/qimen/almanac"

// StepPalace advances one step in 阴/阳 遁 order through the full 1..9
// ring (the center palace is a regular stop — 地盘 lays a stem there).
func StepPalace(palace uint8, yy almanac.YinYang) uint8 {
	if yy == almanac.Yang {
		if palace == 9 {
			return 1
		}
		return palace + 1
	}
	if palace == 1 {
		return 9
	}
	return palace - 1
}

// MoveBy steps n times in 阴/阳 遁 order (forward for 阳, backward for
// 阴, both through the full 1..9 ring) and returns the real destination
// palace — including the center (5). Callers needing a ring position
// project the center to its 寄宫 (坤 2) themselves.
// Precondition: palace ∈ [1, 9], steps ≥ 0.
func MoveBy(palace uint8, steps int, yy almanac.YinYang) uint8 {
	var idx int
	if yy == almanac.Yang {
		idx = (int(palace) - 1 + steps) % 9
	} else {
		idx = (int(palace)-1-steps)%9 + 9
	}
	return uint8(idx%9 + 1)
}

// AreOpposite reports whether two palaces are diametrically opposite.
//
//	1↔9, 2↔8, 3↔7, 4↔6.
func AreOpposite(a, b uint8) bool {
	switch {
	case a == 1 && b == 9, a == 9 && b == 1:
		return true
	case a == 2 && b == 8, a == 8 && b == 2:
		return true
	case a == 3 && b == 7, a == 7 && b == 3:
		return true
	case a == 4 && b == 6, a == 6 && b == 4:
		return true
	}
	return false
}
