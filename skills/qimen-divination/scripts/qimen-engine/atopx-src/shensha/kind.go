// Package shensha detects 神煞 (10 types) given a chart view.
package shensha

// Kind is one of the recognized 神煞.
type Kind uint8

const (
	YiMa     Kind = iota // 驿马 — 主迁移、出行、变动
	TaoHua               // 桃花 — 主异性缘、魅力 (又名 咸池)
	HuaGai               // 华盖 — 主孤高、艺术、宗教、清高
	TianYi               // 天乙贵人 — 主逢凶化吉、贵人相助
	TianDe               // 天德贵人 — 上天之德庇佑,主化灾解难
	YueDe                // 月德贵人 — 月令之德庇佑,主品行端正
	GuoYin               // 国印贵人 — 主掌权印信、仕途亨通
	WenChang             // 文昌 — 主学业、文采、考试、功名
	LuShen               // 禄神 — 代表俸禄、衣食、财源
	YangRen              // 羊刃 — 仅阳干, 主血光、刑伤、果敢
)
