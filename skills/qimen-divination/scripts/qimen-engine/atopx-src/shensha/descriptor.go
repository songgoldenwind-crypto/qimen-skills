package shensha

import "github.com/atopx/qimen/auspice"

type descriptor struct {
	name    string
	summary string
	auspice auspice.Auspice
}

var registry = [...]descriptor{
	YiMa:     {"驿马", "主迁移、出行、变动", auspice.Neutral},
	TaoHua:   {"桃花", "主异性缘、魅力、风流韵事", auspice.Neutral},
	HuaGai:   {"华盖", "主孤高、艺术、宗教、清高", auspice.Neutral},
	TianYi:   {"天乙贵人", "八字神煞之首,主逢凶化吉、贵人相助", auspice.GreatAuspicious},
	TianDe:   {"天德贵人", "上天之德庇佑,主化灾解难", auspice.GreatAuspicious},
	YueDe:    {"月德贵人", "月令之德庇佑,主品行端正", auspice.GreatAuspicious},
	GuoYin:   {"国印贵人", "主掌权印信、仕途亨通", auspice.Auspicious},
	WenChang: {"文昌", "主学业、文采、考试、功名", auspice.Auspicious},
	LuShen:   {"禄神", "代表俸禄、衣食、财源", auspice.Auspicious},
	YangRen:  {"羊刃", "刚猛之气的极端表现,主血光、刑伤、果敢", auspice.Inauspicious},
}

// Name returns the Chinese name.
func (k Kind) Name() string { return registry[k].name }

// Summary returns a one-line description.
func (k Kind) Summary() string { return registry[k].summary }

// Auspice returns the literature-defined auspice level.
func (k Kind) Auspice() auspice.Auspice { return registry[k].auspice }

// String implements fmt.Stringer.
func (k Kind) String() string { return k.Name() }
