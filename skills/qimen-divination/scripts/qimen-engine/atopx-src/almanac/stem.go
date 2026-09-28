package almanac

// Stem represents one of the 10 heavenly stems (天干).
// Values 0..9 correspond to 甲, 乙, 丙, 丁, 戊, 己, 庚, 辛, 壬, 癸.
type Stem uint8

// Canonical stem indices.
const (
	Jia Stem = iota
	Yi
	Bing
	Ding
	Wu
	Ji
	Geng
	Xin
	Ren
	Gui
)

var stemNames = [10]string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// Element index ordering: 0=木, 1=火, 2=土, 3=金, 4=水.
//
// 甲乙→木, 丙丁→火, 戊己→土, 庚辛→金, 壬癸→水.
var stemElement = [10]uint8{0, 0, 1, 1, 2, 2, 3, 3, 4, 4}

// terrainBase encodes 长生 location for each stem indexed by 12-branch step.
// Pattern: 阳干顺行, 阴干逆行 (mirrored for 乙丁己辛癸).
//
// Indices: base[stem] is the 12-position of "长生" for that stem in 寅..丑 cycle.
//
// 甲长生在亥(11)→ shifted relative to branch 寅(2).
// Direct table: stem → branch index of 长生 (0..11 = 子..亥).
var stemTerrainStart = [10]uint8{
	11, // 甲 长生在亥
	6,  // 乙 长生在午
	2,  // 丙 长生在寅
	9,  // 丁 长生在酉
	2,  // 戊 长生在寅
	9,  // 己 长生在酉
	5,  // 庚 长生在巳
	0,  // 辛 长生在子
	8,  // 壬 长生在申
	3,  // 癸 长生在卯
}

// StemOf wraps an integer into Stem, normalizing into 0..9.
func StemOf(i int) Stem {
	i = ((i % 10) + 10) % 10
	return Stem(i)
}

// Index returns the 0..9 ordinal.
func (s Stem) Index() int { return int(s) }

// Name returns the Chinese character.
func (s Stem) Name() string { return stemNames[s] }

// String implements fmt.Stringer.
func (s Stem) String() string { return s.Name() }

// YinYang of the stem: even index = 阳, odd index = 阴.
func (s Stem) YinYang() YinYang {
	if s&1 == 0 {
		return Yang
	}
	return Yin
}

// TenStarOf returns the 十神 relationship from self s to target.
//
// Categories (element relation between self and target, in pairs of 2):
//
//	0,1 比肩/劫财 — same element
//	2,3 食神/伤官 — self generates target (我生他)
//	4,5 偏财/正财 — self controls target (我克他)
//	6,7 七杀/正官 — target controls self (克我)
//	8,9 偏印/正印 — target generates self (生我)
//
// Within each pair, same-parity (both 阳 or both 阴) picks the even index
// (比肩/食神/偏财/七杀/偏印); mixed parity picks the odd index.
func (s Stem) TenStarOf(target Stem) TenStar {
	se := int(stemElement[s])
	te := int(stemElement[target])
	var cat int
	switch {
	case se == te:
		cat = 0
	case (se+1)%5 == te:
		cat = 1 // self generates target (木→火→土→金→水→木)
	case (se+2)%5 == te:
		cat = 2 // self controls target (木→土→水→火→金→木)
	case (te+2)%5 == se:
		cat = 3 // target controls self
	case (te+1)%5 == se:
		cat = 4 // target generates self
	}
	sub := 0
	if (int(s) & 1) != (int(target) & 1) {
		sub = 1
	}
	return TenStar(cat*2 + sub)
}

// TerrainOf returns the 长生十二宫 position for stem s with branch b.
//
//	0=长生 1=沐浴 2=冠带 3=临官 4=帝旺 5=衰 6=病 7=死 8=墓 9=绝 10=胎 11=养
//
// Yang stems advance forward through the 12 branches; yin stems retrograde.
func (s Stem) TerrainOf(b Branch) Terrain {
	start := int(stemTerrainStart[s])
	branchIdx := int(b)
	var off int
	if s.YinYang() == Yang {
		off = ((branchIdx - start) + 12) % 12
	} else {
		off = ((start - branchIdx) + 12) % 12
	}
	return Terrain(off)
}
