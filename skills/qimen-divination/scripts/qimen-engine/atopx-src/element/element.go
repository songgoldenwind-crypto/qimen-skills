// Package element defines the 五行 type and 生克 relationship enum
// (used by both palace classifications and door/star × palace analyses).
package element

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/auspice"
	"github.com/atopx/qimen/enum"
)

// Element 五行.
type Element uint8

const (
	// Wood 木
	Wood Element = iota
	// Fire 火
	Fire
	// Earth 土
	Earth
	// Metal 金
	Metal
	// Water 水
	Water
)

var names = [5]string{"木", "火", "土", "金", "水"}

// Name returns the Chinese label.
func (e Element) Name() string { return names[e] }

// String implements fmt.Stringer.
func (e Element) String() string { return e.Name() }

// stemElements maps stem index (0..9 = 甲..癸) → Element.
//
//	甲乙→木, 丙丁→火, 戊己→土, 庚辛→金, 壬癸→水.
var stemElements = [10]Element{Wood, Wood, Fire, Fire, Earth, Earth, Metal, Metal, Water, Water}

// OfStem returns the element of a heavenly stem.
func OfStem(s almanac.Stem) Element { return stemElements[s.Index()] }

// branchElements maps branch index (0..11 = 子..亥) → Element.
//
//	子=水 丑=土 寅卯=木 辰=土 巳午=火 未=土 申酉=金 戌=土 亥=水.
var branchElements = [12]Element{Water, Earth, Wood, Wood, Earth, Fire, Fire, Earth, Metal, Metal, Earth, Water}

// OfBranch returns the element of an earthly branch.
func OfBranch(b almanac.Branch) Element { return branchElements[b.Index()] }

// palaceElements maps palace 1..9 → Element. Index 0 reserved.
//
//	1坎=水 2坤/5中/8艮=土 3震/4巽=木 6乾/7兑=金 9离=火.
var palaceElements = [10]Element{Earth, Water, Earth, Wood, Wood, Earth, Metal, Metal, Earth, Fire}

// FromPalace returns the element of a 九宫 palace.
// Precondition: palace ∈ [1, 9].
func FromPalace(palace uint8) Element { return palaceElements[palace] }

// starElements maps enum.Star (0..9) → Element.
//
//	天蓬=水 天芮/天禽/天任/禽芮=土 天冲/天辅=木
//	天心/天柱=金 天英=火.
var starElements = [10]Element{Water, Earth, Wood, Wood, Earth, Metal, Metal, Earth, Fire, Earth}

// doorElements maps enum.Door (0..7) → Element.
//
//	休=水 生/死=土 伤/杜=木 景=火 惊/开=金.
var doorElements = [8]Element{Water, Earth, Wood, Wood, Fire, Earth, Metal, Metal}

// OfStar returns the 五行 of a 九星.
func OfStar(s enum.Star) Element { return starElements[s] }

// OfDoor returns the 五行 of a 八门.
func OfDoor(d enum.Door) Element { return doorElements[d] }

// Relation classifies the 五行 generative/restraining relationship
// from the subject (self) to a counterpart (other).
type Relation uint8

const (
	// Same 比和 (same element)
	Same Relation = iota
	// Generates 生出 (self → other; self gives, other receives)
	Generates
	// Generated 受生 (other → self; self receives)
	Generated
	// Restrains 克出 (self → other; self attacks)
	Restrains
	// Restrained 受克 (other → self; self is attacked)
	Restrained
)

var relationNames = [5]string{"比和", "生出", "受生", "克出", "受克"}

// Name returns the Chinese label.
func (r Relation) Name() string { return relationNames[r] }

// String implements fmt.Stringer.
func (r Relation) String() string { return r.Name() }

// relationAuspice maps Relation → auspice (from subject's perspective).
//
//   - 比和 / 克出 → 中和
//   - 受生 → 吉
//   - 生出 / 受克 → 凶
var relationAuspice = [5]auspice.Auspice{
	Same:       auspice.Neutral,
	Generates:  auspice.Inauspicious,
	Generated:  auspice.Auspicious,
	Restrains:  auspice.Neutral,
	Restrained: auspice.Inauspicious,
}

// AuspiceAsSelf returns the auspice from the subject's perspective.
func (r Relation) AuspiceAsSelf() auspice.Auspice { return relationAuspice[r] }

// relationTable[e][other] is precomputed e.RelationTo(other) for all 5×5
// element pairs. Built at init from the 生 / 克 cycles.
//
// Element ordering (iota): Wood=0, Fire=1, Earth=2, Metal=3, Water=4.
var relationTable [5][5]Relation

func init() {
	// gen[e] = what e generates (相生: 木→火→土→金→水→木).
	//   gen[Wood]=Fire, gen[Fire]=Earth, gen[Earth]=Metal, gen[Metal]=Water, gen[Water]=Wood.
	gen := [5]Element{Fire, Earth, Metal, Water, Wood}
	// res[e] = what e restrains (相克: 木→土→水→火→金→木).
	//   res[Wood]=Earth, res[Fire]=Metal, res[Earth]=Water, res[Metal]=Wood, res[Water]=Fire.
	res := [5]Element{Earth, Metal, Water, Wood, Fire}
	for e := Element(0); e < 5; e++ {
		for o := Element(0); o < 5; o++ {
			switch {
			case e == o:
				relationTable[e][o] = Same
			case gen[e] == o:
				relationTable[e][o] = Generates
			case gen[o] == e:
				relationTable[e][o] = Generated
			case res[e] == o:
				relationTable[e][o] = Restrains
			default:
				relationTable[e][o] = Restrained
			}
		}
	}
}

// RelationTo classifies e's relation to other.
// O(1) table lookup over the 25-entry 五行 generation/restraint matrix.
func (e Element) RelationTo(other Element) Relation { return relationTable[e][other] }
