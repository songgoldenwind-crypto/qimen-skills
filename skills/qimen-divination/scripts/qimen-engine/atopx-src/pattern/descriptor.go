package pattern

import "github.com/atopx/qimen/auspice"

type descriptor struct {
	name    string
	summary string
	auspice auspice.Auspice
}

// registry indexes every recognized Kind to its name/summary/auspice tuple.
var registry = [...]descriptor{
	FanYin:          {"反吟", "奇门凶格。星门反吟,反复无常,事情易变。", auspice.Inauspicious},
	FuYin:           {"伏吟", "奇门凶格。星门伏吟,事情停滞,宜守不宜进。", auspice.Inauspicious},
	RuMu:            {"入墓", "奇门凶格。天盘干落入墓库之宫,艰难阻塞,难以发展。", auspice.Inauspicious},
	KongWang:        {"落空亡", "奇门凶格。落入空亡宫,心愿落空,难以实现。", auspice.Inauspicious},
	MenPo:           {"门迫", "奇门凶格。门克宫位,门被迫害,谋事不成,阻碍重重。", auspice.Inauspicious},
	YiQiDeShi:       {"乙奇得使", "奇门吉格。乙奇临开门,利于谋划,贵人相助,诸事吉利。", auspice.Auspicious},
	BingQiDeShi:     {"丙奇得使", "奇门吉格。丙奇临休门,光明正大,官司必胜,声名可得。", auspice.Auspicious},
	DingQiDeShi:     {"丁奇得使", "奇门吉格。丁奇临生门,才思敏捷,求财必得,生意兴隆。", auspice.Auspicious},
	TianDun:         {"天遁", "奇门吉格。丙丁同临生门,天助之,万事亨通,大吉大利。", auspice.Auspicious},
	DiDun:           {"地遁", "奇门吉格。乙奇临开门加九地,隐匿藏形,避凶趋吉。", auspice.Auspicious},
	RenDun:          {"人遁", "奇门吉格。丁奇临休门加太阴,人和之象,贵人暗助。", auspice.Auspicious},
	ShenDun:         {"神遁", "奇门吉格。丙奇临生门加九天,神助之象,心想事成。", auspice.Auspicious},
	GuiDun:          {"鬼遁", "奇门吉格。丁奇临杜门加九地,神秘莫测,暗中成事。", auspice.Auspicious},
	FengDun:         {"风遁", "奇门吉格。乙奇临杜门在巽宫,运筹帷幄,避开祸端。", auspice.Auspicious},
	YunDun:          {"云遁", "奇门吉格。乙奇临开门在乾宫,腾云驾雾,步步高升。", auspice.Auspicious},
	LongDun:         {"龙遁", "奇门吉格。乙奇临休门在坎宫,龙入大海,鸿图大展。", auspice.Auspicious},
	HuDun:           {"虎遁", "奇门吉格。乙奇临开门在兑宫,猛虎添翼,势不可挡。", auspice.Auspicious},
	QingLongFanShou: {"青龙返首", "奇门大吉格。天盘戊临地盘丙,大吉大利,名利双收。", auspice.GreatAuspicious},
	FeiNiaoDieXue:   {"飞鸟跌穴", "奇门大吉格。天盘丙临地盘戊,诸事顺遂,不求自得。", auspice.GreatAuspicious},
	DaGe:            {"大格", "奇门凶格。庚临癸上,谋事难成,处处受制,大凶。", auspice.GreatInauspicious},
	XiaoGe:          {"小格", "奇门凶格。庚临壬上,小有阻碍,谋事迟缓。", auspice.Inauspicious},
	XingGe:          {"刑格", "奇门凶格。庚临己上,官司牢狱,纷争不断。", auspice.Inauspicious},
	BoGe:            {"悖格", "奇门凶格。庚金克制三奇,悖逆阻碍,诸事不顺,主凶。", auspice.GreatInauspicious},
	TianWangSiZhang: {"天网四张", "奇门凶格。癸水入火域,身陷天网,行动招祸,主凶。", auspice.GreatInauspicious},
}

// Name returns the Chinese name.
func (k Kind) Name() string { return registry[k].name }

// Summary returns the one-line classical description.
func (k Kind) Summary() string { return registry[k].summary }

// Auspice returns the literature-defined auspice level.
func (k Kind) Auspice() auspice.Auspice { return registry[k].auspice }

// String implements fmt.Stringer.
func (k Kind) String() string { return k.Name() }
