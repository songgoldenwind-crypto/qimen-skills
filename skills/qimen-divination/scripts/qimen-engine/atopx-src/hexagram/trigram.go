// Package hexagram implements the 八卦 (Trigram) and 六十四卦 (Hexagram)
// types used by qimen's 门-宫 演卦.
package hexagram

// Trigram 八卦.
//
//	0乾 ☰ 天 / 1兑 ☱ 泽 / 2离 ☲ 火 / 3震 ☳ 雷
//	4巽 ☴ 风 / 5坎 ☵ 水 / 6艮 ☶ 山 / 7坤 ☷ 地
type Trigram uint8

const (
	Qian Trigram = iota // 乾
	Dui                 // 兑
	Li                  // 离
	Zhen                // 震
	Xun                 // 巽
	Kan                 // 坎
	Gen                 // 艮
	Kun                 // 坤
)

var trigramNames = [8]string{"乾", "兑", "离", "震", "巽", "坎", "艮", "坤"}
var trigramSymbols = [8]string{"☰", "☱", "☲", "☳", "☴", "☵", "☶", "☷"}
var trigramElems = [8]string{"天", "泽", "火", "雷", "风", "水", "山", "地"}

// palaceTrigram maps palace 1..9 → Trigram. Index 0 and 5 (center) hold
// the zero value Qian — callers MUST gate by `palace != 5` first, since
// center has no trigram.
//
//	1坎 / 2坤 / 3震 / 4巽 / 6乾 / 7兑 / 8艮 / 9离.
var palaceTrigram = [10]Trigram{0, Kan, Kun, Zhen, Xun, 0, Qian, Dui, Gen, Li}

// TrigramOfPalace returns the canonical 八卦 for a non-center palace.
// Precondition: palace ∈ [1, 9] AND palace != 5 (center has no trigram).
func TrigramOfPalace(palace uint8) Trigram { return palaceTrigram[palace] }

// Name returns the Chinese trigram name.
func (t Trigram) Name() string { return trigramNames[t] }

// Symbol returns the Unicode trigram symbol (☰..☷).
func (t Trigram) Symbol() string { return trigramSymbols[t] }

// ElementName returns the natural-element name (天/泽/火/雷/风/水/山/地).
func (t Trigram) ElementName() string { return trigramElems[t] }

// String implements fmt.Stringer.
func (t Trigram) String() string { return t.Name() }
