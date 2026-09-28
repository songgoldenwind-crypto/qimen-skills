package enum

// Star 九星 (天蓬/天芮/天冲/天辅/天禽/天心/天柱/天任/天英).
//
// StarQinRui is the special "禽芮" marker used when 天禽 and 天芮 both
// fall in palace 2 (since 天禽 has no native palace, by tradition it is
// merged into 天芮's slot).
type Star uint8

const (
	StarTianPeng  Star = iota // 天蓬
	StarTianRui               // 天芮
	StarTianChong             // 天冲
	StarTianFu                // 天辅
	StarTianQin               // 天禽
	StarTianXin               // 天心
	StarTianZhu               // 天柱
	StarTianRen               // 天任
	StarTianYing              // 天英
	StarQinRui                // 禽芮 (merged 天禽/天芮 marker)
)

var starNames = [10]string{"天蓬", "天芮", "天冲", "天辅", "天禽", "天心", "天柱", "天任", "天英", "禽芮"}

// Name returns the Chinese label.
func (s Star) Name() string { return starNames[s] }

// String implements fmt.Stringer.
func (s Star) String() string { return s.Name() }

// starHomePalace maps Star index → home palace (1..9).
//
//	天蓬→1, 天芮/禽芮→2, 天冲→3, 天辅→4, 天禽→5,
//	天心→6, 天柱→7, 天任→8, 天英→9.
var starHomePalace = [10]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 2}

// HomePalace returns the canonical 本位宫 (1..9).
func (s Star) HomePalace() uint8 { return starHomePalace[s] }

// palaceStar maps palace 1..9 → Star. Index 0 reserved.
var palaceStar = [10]Star{
	0,
	StarTianPeng, StarTianRui, StarTianChong, StarTianFu, StarTianQin,
	StarTianXin, StarTianZhu, StarTianRen, StarTianYing,
}

// StarOfPalace returns the canonical 本位星 for a palace.
// Precondition: palace ∈ [1, 9].
func StarOfPalace(palace uint8) Star { return palaceStar[palace] }
