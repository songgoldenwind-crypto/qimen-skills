// Package terrain wraps almanac.Terrain to add qimen-specific
// Auspice classification.
package terrain

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/auspice"
)

// Terrain is the qimen-side wrapper around almanac.Terrain with auspice.
type Terrain struct {
	Inner almanac.Terrain
}

// Of wraps an almanac.Terrain.
func Of(t almanac.Terrain) Terrain { return Terrain{Inner: t} }

// Index returns the 0..11 ordinal.
func (t Terrain) Index() int { return t.Inner.Index() }

// Name returns the Chinese name.
func (t Terrain) Name() string { return t.Inner.Name() }

// Summary returns a one-line classical description.
func (t Terrain) Summary() string { return summaries[t.Index()] }

// Auspice returns the literature-defined auspice level.
//
//   - 长生 / 冠带 / 临官 / 养 → 吉
//   - 帝旺 → 大吉
//   - 沐浴 / 衰 / 胎 → 中和
//   - 病 / 墓 → 凶
//   - 死 / 绝 → 大凶
func (t Terrain) Auspice() auspice.Auspice { return auspices[t.Index()] }

// String implements fmt.Stringer.
func (t Terrain) String() string { return t.Name() }

var auspices = [12]auspice.Auspice{
	auspice.Auspicious,        // 0  长生
	auspice.Neutral,           // 1  沐浴
	auspice.Auspicious,        // 2  冠带
	auspice.Auspicious,        // 3  临官
	auspice.GreatAuspicious,   // 4  帝旺
	auspice.Neutral,           // 5  衰
	auspice.Inauspicious,      // 6  病
	auspice.GreatInauspicious, // 7  死
	auspice.Inauspicious,      // 8  墓
	auspice.GreatInauspicious, // 9  绝
	auspice.Neutral,           // 10 胎
	auspice.Auspicious,        // 11 养
}

var summaries = [12]string{
	"长生十二宫之首。万物初生之态,如人之出生,充满生机与希望。",
	"长生十二宫第二位。又名\"桃花\",如人初生后沐浴净身,象征不稳定与桃花。",
	"长生十二宫第三位。如人成年加冠,开始走向社会,代表成长、自信。",
	"长生十二宫第四位。又名\"禄\",如人入仕做官,代表壮年、权力。",
	"长生十二宫第五位。五行力量巅峰,如帝王之旺,极盛之态。盛极必衰。",
	"长生十二宫第六位。帝旺之后力量开始衰退,如人步入中年体力下降。",
	"长生十二宫第七位。力量进一步衰弱,如人生病,需休养调理。",
	"长生十二宫第八位。万物气绝,如人之死亡,代表终结与静止。",
	"长生十二宫第九位。又名\"库\",如人入土归葬,代表收藏、积蓄、隐藏。",
	"长生十二宫第十位。万物气息断绝,旧形全消,但蕴含新生的转机。",
	"长生十二宫第十一位。新生命在母体中孕育,代表新计划、新希望的萌芽。",
	"长生十二宫第十二位。胎儿在母体中成长发育,即将出世,代表蓄养、培育。",
}
