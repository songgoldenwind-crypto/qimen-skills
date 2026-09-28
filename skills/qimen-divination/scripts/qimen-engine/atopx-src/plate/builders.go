package plate

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/internal/tables"
)

// BuildEarth populates 地盘 (三奇六仪) from the local 局 and 阴/阳 遁.
func BuildEarth(yy almanac.YinYang, ju uint8) StemPlate {
	var p StemPlate
	palace := ju
	for _, stemIdx := range tables.SanQiLiuYi {
		p.Set(palace, almanac.StemOf(int(stemIdx)))
		palace = StepPalace(palace, yy)
	}
	return p
}

// RotateStems rotates the 地盘 stems rigidly around the LuoShu ring by
// the arc from palace `from` to palace `to` (the center stem stays in
// place). Rotation is direction-free: 阴/阳 遁 only decide where the
// duty lands, not how the ring turns. This single primitive lays both
// rotated stem plates of a chart:
//
//   - 天盘: the stems follow the 值符 (from 值符原宫 to 时干落宫);
//   - 暗干: the stems follow the 值使 (from 值使门本位宫 to 值使落宫).
//
// Precondition: from, to ∈ [1, 9] \ {5} (center is projected to 坤 2
// by the caller).
func RotateStems(earth *StemPlate, from, to uint8) StemPlate {
	var p StemPlate
	shift := (int(tables.LuoShuIndex[to]) - int(tables.LuoShuIndex[from]) + 8) % 8
	for i, palace := range tables.LuoShuOrder {
		if stem, ok := earth.Get(palace); ok {
			p.Set(tables.LuoShuOrder[(i+shift)%8], stem)
		}
	}
	if center, ok := earth.Get(5); ok {
		p.Set(5, center)
	}
	return p
}

// BuildStar populates 九星盘. 禽 falls into 2 with the merged 禽芮 marker.
func BuildStar(zhiFuOriginalPalace, hourPalace uint8) StarPlate {
	var p StarPlate
	zhiFuIdx := int(tables.LuoShuIndex[zhiFuOriginalPalace])
	hourIdx := int(tables.LuoShuIndex[hourPalace])
	steps := (hourIdx + 8 - zhiFuIdx) % 8
	for i, originalPalace := range tables.LuoShuOrder {
		targetPalace := tables.LuoShuOrder[(i+steps)%8]
		// 禽 has no native palace; falls into palace 2's slot as 禽芮.
		var star enum.Star
		if originalPalace == 2 {
			star = enum.StarQinRui
		} else {
			star = enum.StarOfPalace(originalPalace)
		}
		p.Set(targetPalace, star)
	}
	return p
}

// MarchSteps returns how many palaces the 值使 marches: the duty
// pillar's branch distance from the first branch of its 旬.
func MarchSteps(lead almanac.Cycle) int {
	return (lead.Branch().Index() + 12 - lead.Ten().FirstBranch().Index()) % 12
}

// BuildDoor populates the rotate-style 八门盘 and returns the (real,
// possibly center) 值使落宫 alongside it.
//
// The 值使 door marches from the real 值符原宫 (marchFrom, which may be
// the center palace 5) by the number of duty-pillar steps since the
// 旬首; the eight doors then rotate rigidly by the arc from the door's
// home palace (zhiShiHome, the 寄坤 projection of the origin) to the
// landing palace (likewise projected when it is the center).
func BuildDoor(yy almanac.YinYang, zhiShiHome, marchFrom uint8, lead almanac.Cycle) (DoorPlate, uint8) {
	var p DoorPlate
	zhiShiPalace := MoveBy(marchFrom, MarchSteps(lead), yy)
	landEff := zhiShiPalace
	if landEff == 5 {
		landEff = 2
	}
	shift := (int(tables.LuoShuIndex[landEff]) - int(tables.LuoShuIndex[zhiShiHome]) + 8) % 8
	for i, originalPalace := range tables.LuoShuOrder {
		// LuoShuOrder excludes center palace 5; every entry has a door.
		p.Set(tables.LuoShuOrder[(i+shift)%8], enum.DoorOfPalace(originalPalace))
	}
	return p, zhiShiPalace
}

// BuildGod populates 九神盘: the gods follow the LuoShu ring from the
// 值符落宫, forward in 阳遁 and backward in 阴遁.
func BuildGod(yy almanac.YinYang, zhiFuPalace uint8) GodPlate {
	var p GodPlate
	startPalace := zhiFuPalace
	if startPalace == 5 {
		startPalace = 2
	}
	start := int(tables.LuoShuIndex[startPalace])
	for i, god := range tables.GodsOrder {
		var pos int
		if yy == almanac.Yang {
			pos = (start + i) % 8
		} else {
			pos = (start - i + 8) % 8
		}
		p.Set(tables.LuoShuOrder[pos], god)
	}
	return p
}

// ===================== StyleFly (飞盘) builders =====================
//
// The fly-plate style replaces the LuoShu ring rotation with palace-
// number flying: every element shifts by the same delta through the
// full 1..9 sequence, the center palace participating as a regular
// stop. 天禽 occupies its real landing palace (no 寄宫 merge), so a
// fly chart has nine stars on nine palaces, while the eight doors and
// eight gods leave exactly one palace empty.

// flyTo returns the palace reached from `palace` after `delta` steps of
// palace-number flying (delta ∈ [0, 8]).
func flyTo(palace uint8, delta int) uint8 {
	return uint8((int(palace)-1+delta)%9 + 1)
}

// FlyDelta returns the flying shift that carries palace `from` onto
// palace `to` through the 1..9 sequence.
func FlyDelta(from, to uint8) int {
	return floorMod9(int(to) - int(from))
}

func floorMod9(x int) int { return ((x % 9) + 9) % 9 }

// FlyStems lays a stem plate by flying every 地盘 stem `delta` palaces
// forward (used for the fly-style 天盘 and 暗干).
func FlyStems(earth *StemPlate, delta int) StemPlate {
	var p StemPlate
	for palace := uint8(1); palace <= 9; palace++ {
		if stem, ok := earth.Get(palace); ok {
			p.Set(flyTo(palace, delta), stem)
		}
	}
	return p
}

// FlyStars lays the nine-star plate by flying each star from its home
// palace; 天禽 lands on a real palace instead of merging into 芮.
func FlyStars(delta int) StarPlate {
	var p StarPlate
	for home := uint8(1); home <= 9; home++ {
		p.Set(flyTo(home, delta), enum.StarOfPalace(home))
	}
	return p
}

// FlyDoors lays the eight-door plate by flying each door from its home
// palace; the palace receiving the (door-less) center slot stays empty.
func FlyDoors(delta int) DoorPlate {
	var p DoorPlate
	for _, home := range tables.LuoShuOrder {
		p.Set(flyTo(home, delta), enum.DoorOfPalace(home))
	}
	return p
}

// flyGodsOrder is the fly-style NINE-god sequence: the rotate-style
// eight gods with 太常 inserted after 白虎, filling all nine palaces.
var flyGodsOrder = [9]enum.God{
	enum.GodZhiFu,
	enum.GodTengShe,
	enum.GodTaiYin,
	enum.GodLiuHe,
	enum.GodBaiHu,
	enum.GodTaiChang,
	enum.GodXuanWu,
	enum.GodJiuDi,
	enum.GodJiuTian,
}

// FlyGods lays the nine gods flying from the 值符落宫 — forward through
// the palace numbers in 阳遁, backward in 阴遁 — covering every palace.
func FlyGods(yy almanac.YinYang, zhiFuPalace uint8) GodPlate {
	var p GodPlate
	for i, god := range flyGodsOrder {
		var palace int
		if yy == almanac.Yang {
			palace = floorMod9(int(zhiFuPalace) - 1 + i)
		} else {
			palace = floorMod9(int(zhiFuPalace) - 1 - i)
		}
		p.Set(uint8(palace+1), god)
	}
	return p
}
