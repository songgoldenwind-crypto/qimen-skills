// Package pattern detects 奇门格局 (二十四格局) given an unpacked
// chart view.
package pattern

// Kind is one of the 24 recognized 奇门格局.
type Kind uint8

const (
	FanYin          Kind = iota // 反吟 — 值符落入与原宫对冲的宫位
	FuYin                       // 伏吟 — 值符落入原宫
	RuMu                        // 入墓 — 天盘天干落自己的墓库宫
	KongWang                    // 落空亡 — 旬空亡支所在宫位
	MenPo                       // 门迫 — 门被宫位五行所克
	YiQiDeShi                   // 乙奇得使 — 乙奇临开门
	BingQiDeShi                 // 丙奇得使 — 丙奇临休门
	DingQiDeShi                 // 丁奇得使 — 丁奇临生门
	TianDun                     // 天遁 — 天盘丙、地盘丁、生门
	DiDun                       // 地遁 — 乙奇、开门、九地
	RenDun                      // 人遁 — 丁奇、休门、太阴
	ShenDun                     // 神遁 — 丙奇、生门、九天
	GuiDun                      // 鬼遁 — 丁奇、杜门、九地
	FengDun                     // 风遁 — 乙奇、开/杜门、巽 4 宫
	YunDun                      // 云遁 — 乙奇、开门、乾 6 宫
	LongDun                     // 龙遁 — 乙奇、休门、坎 1 宫
	HuDun                       // 虎遁 — 乙奇、开门、兑 7 宫
	QingLongFanShou             // 青龙返首 — 天盘戊、地盘丙
	FeiNiaoDieXue               // 飞鸟跌穴 — 天盘丙、地盘戊
	DaGe                        // 大格 — 天盘庚、地盘癸
	XiaoGe                      // 小格 — 天盘庚、地盘壬
	XingGe                      // 刑格 — 天盘庚、地盘己
	BoGe                        // 悖格 — 丙/庚 互克
	TianWangSiZhang             // 天网四张 — 双癸
)
