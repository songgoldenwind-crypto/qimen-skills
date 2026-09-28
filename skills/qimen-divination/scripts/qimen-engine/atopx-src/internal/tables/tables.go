// Package tables exposes the static tables used by qimen layout / detection.
//
// All tables are package-level vars so they can be tree-shaken into a
// single read-only segment by the linker. Callers in other internal
// packages reference these directly without indirection.
package tables

import "github.com/atopx/qimen/enum"

// SanQiLiuYi 三奇六仪 fill order: 戊 己 庚 辛 壬 癸 丁 丙 乙
// (heavenly-stem indices 4, 5, 6, 7, 8, 9, 3, 2, 1).
var SanQiLiuYi = [9]uint8{4, 5, 6, 7, 8, 9, 3, 2, 1}

// LuoShuOrder 9-palace layout in 洛书 order (excluding center 5):
// 坎 艮 震 巽 离 坤 兑 乾.
var LuoShuOrder = [8]uint8{1, 8, 3, 4, 9, 2, 7, 6}

// LuoShuIndex reverse lookup: palace 1..9 → position in LuoShuOrder.
// Center (palace 5) maps to 0 as a sentinel; callers must filter.
var LuoShuIndex = [10]uint8{0, 0, 5, 2, 3, 0, 7, 6, 1, 4}

// GodsOrder 九神 in canonical sequence: 值符 腾蛇 太阴 六合 白虎 玄武 九地 九天.
var GodsOrder = [8]enum.God{
	enum.GodZhiFu,
	enum.GodTengShe,
	enum.GodTaiYin,
	enum.GodLiuHe,
	enum.GodBaiHu,
	enum.GodXuanWu,
	enum.GodJiuDi,
	enum.GodJiuTian,
}

// Grid 3×3 nine-palace display [巽 离 坤 / 震 中 兑 / 艮 坎 乾].
var Grid = [3][3]uint8{{4, 9, 2}, {3, 5, 7}, {8, 1, 6}}

// PalaceNames index 1..9 → 中文 palace name. Index 0 reserved.
var PalaceNames = [10]string{"", "坎", "坤", "震", "巽", "中", "乾", "兑", "艮", "离"}

// TermJu solar-term index (0..23) → [上元局, 中元局, 下元局].
var TermJu = [24][3]uint8{
	{1, 7, 4}, // 0  冬至
	{2, 8, 5}, // 1  小寒
	{3, 9, 6}, // 2  大寒
	{8, 5, 2}, // 3  立春
	{9, 6, 3}, // 4  雨水
	{1, 7, 4}, // 5  惊蛰
	{3, 9, 6}, // 6  春分
	{4, 1, 7}, // 7  清明
	{5, 2, 8}, // 8  谷雨
	{4, 1, 7}, // 9  立夏
	{5, 2, 8}, // 10 小满
	{6, 3, 9}, // 11 芒种
	{9, 3, 6}, // 12 夏至
	{8, 2, 5}, // 13 小暑
	{7, 1, 4}, // 14 大暑
	{2, 5, 8}, // 15 立秋
	{1, 4, 7}, // 16 处暑
	{9, 3, 6}, // 17 白露
	{7, 1, 4}, // 18 秋分
	{6, 9, 3}, // 19 寒露
	{5, 8, 2}, // 20 霜降
	{6, 9, 3}, // 21 立冬
	{5, 8, 2}, // 22 小雪
	{4, 7, 1}, // 23 大雪
}

// TenXunShou 旬首 stem indices: 甲子→戊, 甲戌→己, 甲申→庚, 甲午→辛, 甲辰→壬, 甲寅→癸.
var TenXunShou = [6]uint8{4, 5, 6, 7, 8, 9}

// BranchToPalace 地支 (0..11) → 落宫.
//
//	子1, 丑寅8, 卯3, 辰巳4, 午9, 未申2, 酉7, 戌亥6.
var BranchToPalace = [12]uint8{1, 8, 8, 3, 4, 4, 9, 2, 2, 7, 6, 6}

// StemTombPalace 天干 (0..9) → 入墓宫位.
//
//	甲癸→2, 乙丙戊→6, 丁己庚→8, 辛壬→4.
var StemTombPalace = [10]uint8{2, 6, 6, 8, 6, 8, 8, 4, 4, 2}

// PalaceBranches indexes palace 1..9 → 地支 indices held by that palace.
// Index 0 and 5 (center) are nil. Slices share package-level backing
// arrays — callers MUST NOT mutate the returned slice.
//
//	1坎→子, 2坤→未申, 3震→卯, 4巽→辰巳,
//	5中→(空), 6乾→戌亥, 7兑→酉, 8艮→丑寅, 9离→午.
var PalaceBranches = [10][]uint8{
	1: {0},
	2: {7, 8},
	3: {3},
	4: {4, 5},
	6: {10, 11},
	7: {9},
	8: {1, 2},
	9: {6},
}
