package qimen

import (
	"fmt"
	"testing"

	"github.com/atopx/qimen/enum"
)

// goldenChart is one authoritative chart transcribed from reference
// paipan software output, covering 时/日/月/年家, 转盘/飞盘 and both
// 定局规则. Palaces are encoded palace 1..9 as six runes: 天盘干 地盘干
// 暗干 星 门 神 (star/door/god are single-rune abbreviations; '-' marks
// an empty slot).
type goldenChart struct {
	when     string // "YYYY-MM-DD HH:MM"
	note     string
	method   enum.Method // zero value = MethodTime
	style    enum.Style  // zero value = StyleRotate
	juRule   enum.JuRule // zero value = JuRuleZhiRun
	pillars  string      // "年 月 日 时"
	yinYang  string
	ju       uint8
	xun      string // 旬 name, e.g. "甲寅"
	xunStem  string // 旬首遁干
	kongWang string // hour-pillar 空亡 pair
	zhiFu    string // duty star name
	zfOrig   uint8
	zfLand   uint8
	zhiShi   string // duty door name
	zsLand   uint8
	palaces  [9]string
}

var goldenCharts = []goldenChart{
	{
		when: "2026-02-18 22:59", note: "立春下元 (伏吟)",
		pillars: "丙午 庚寅 癸亥 癸亥", yinYang: "阳", ju: 2,
		xun: "甲寅", xunStem: "癸", kongWang: "子丑",
		zhiFu: "天柱", zfOrig: 7, zfLand: 7, zhiShi: "惊", zsLand: 7,
		palaces: [9]string{
			"乙乙乙蓬休阴", "戊戊戊芮死天", "己己己冲伤虎",
			"庚庚庚辅杜玄", "辛辛辛---", "壬壬壬心开蛇",
			"癸癸癸柱惊符", "丁丁丁任生合", "丙丙丙英景地",
		},
	},
	{
		when: "2026-02-18 23:00", note: "晚子时换日柱: 雨水上元",
		pillars: "丙午 庚寅 甲子 甲子", yinYang: "阳", ju: 9,
		xun: "甲子", xunStem: "戊", kongWang: "戌亥",
		zhiFu: "天英", zfOrig: 9, zfLand: 9, zhiShi: "景", zsLand: 9,
		palaces: [9]string{
			"己己己蓬休虎", "庚庚庚芮死蛇", "辛辛辛冲伤地",
			"壬壬壬辅杜天", "癸癸癸---", "丁丁丁心开合",
			"丙丙丙柱惊阴", "乙乙乙任生玄", "戊戊戊英景符",
		},
	},
	{
		when: "2026-01-14 18:45", note: "小寒下元",
		pillars: "乙巳 己丑 戊子 辛酉", yinYang: "阳", ju: 8,
		xun: "甲寅", xunStem: "癸", kongWang: "子丑",
		zhiFu: "天辅", zfOrig: 4, zfLand: 2, zhiShi: "杜", zsLand: 2,
		palaces: [9]string{
			"乙庚乙柱惊合", "癸辛癸辅杜符", "庚壬庚蓬休玄",
			"戊癸戊任生地", "丁丁丁---", "辛丙辛芮死阴",
			"己乙己英景蛇", "丙戊丙心开虎", "壬己壬冲伤天",
		},
	},
	{
		when: "2026-01-05 13:00", note: "超神: 节气未交已用小寒上元; 时干落中五",
		pillars: "乙巳 戊子 己卯 辛未", yinYang: "阳", ju: 2,
		xun: "甲子", xunStem: "戊", kongWang: "戌亥",
		zhiFu: "天芮", zfOrig: 2, zfLand: 5, zhiShi: "死", zsLand: 9,
		palaces: [9]string{
			"乙乙丁蓬生合", "戊戊癸芮惊符", "己己庚冲杜玄",
			"庚庚丙辅景地", "辛辛辛---", "壬壬乙心休阴",
			"癸癸壬柱开蛇", "丁丁己任伤虎", "丙丙戊英死天",
		},
	},
	{
		when: "2026-03-02 18:30", note: "旬首遁干落中五: 值符天禽, 值使从五宫起步",
		pillars: "丙午 庚寅 乙亥 乙酉", yinYang: "阳", ju: 3,
		xun: "甲申", xunStem: "庚", kongWang: "午未",
		zhiFu: "天禽", zfOrig: 5, zfLand: 2, zhiShi: "死", zsLand: 6,
		palaces: [9]string{
			"丙丙壬蓬惊合", "乙乙己芮杜符", "戊戊丙冲休玄",
			"己己癸辅生地", "庚庚庚---", "辛辛乙心死阴",
			"壬壬丁柱景蛇", "癸癸辛任开虎", "丁丁戊英伤天",
		},
	},
	{
		when: "2026-10-31 12:02", note: "霜降下元",
		pillars: "丙午 戊戌 戊寅 戊午", yinYang: "阴", ju: 2,
		xun: "甲寅", xunStem: "癸", kongWang: "子丑",
		zhiFu: "天心", zfOrig: 6, zfLand: 2, zhiShi: "开", zsLand: 2,
		palaces: [9]string{
			"乙己乙冲伤玄", "癸戊癸心开符", "庚乙庚英景合",
			"戊丙戊芮死阴", "丁丁丁---", "辛癸辛任生地",
			"己壬己蓬休天", "丙辛丙辅杜虎", "壬庚壬柱惊蛇",
		},
	},
	{
		when: "2024-12-06 12:00", note: "大雪下元",
		pillars: "甲辰 乙亥 甲辰 庚午", yinYang: "阴", ju: 1,
		xun: "甲子", xunStem: "戊", kongWang: "戌亥",
		zhiFu: "天蓬", zfOrig: 1, zfLand: 8, zhiShi: "休", zsLand: 4,
		palaces: [9]string{
			"壬戊乙心死蛇", "己乙丙英伤虎", "庚丙壬任开天",
			"丙丁戊冲休地", "癸癸癸---", "辛壬己柱景阴",
			"乙辛丁芮杜合", "戊庚辛蓬惊符", "丁己庚辅生玄",
		},
	},
	{
		when: "2026-02-03 23:30", note: "晚子时 + 正授: 立春上元第1天",
		pillars: "乙巳 己丑 己酉 甲子", yinYang: "阳", ju: 8,
		xun: "甲子", xunStem: "戊", kongWang: "戌亥",
		zhiFu: "天任", zfOrig: 8, zfLand: 8, zhiShi: "生", zsLand: 8,
		palaces: [9]string{
			"庚庚庚蓬休天", "辛辛辛芮死虎", "壬壬壬冲伤蛇",
			"癸癸癸辅杜阴", "丁丁丁---", "丙丙丙心开地",
			"乙乙乙柱惊玄", "戊戊戊任生符", "己己己英景合",
		},
	},
	{
		when: "2024-12-31 12:00", note: "冬至中元第1天",
		pillars: "甲辰 丙子 己巳 庚午", yinYang: "阳", ju: 7,
		xun: "甲子", xunStem: "戊", kongWang: "戌亥",
		zhiFu: "天柱", zfOrig: 7, zfLand: 9, zhiShi: "惊", zsLand: 4,
		palaces: [9]string{
			"癸辛丁冲杜虎", "乙壬辛心休蛇", "庚癸壬英死地",
			"壬丁戊芮惊天", "丙丙丙---", "己乙癸任伤合",
			"辛戊己蓬生阴", "丁己庚辅景玄", "戊庚乙柱开符",
		},
	},
	{
		when: "2024-12-21 12:00", note: "闰大雪下元第1天 (冬至日, 接气)",
		pillars: "甲辰 丙子 己未 庚午", yinYang: "阴", ju: 1,
		xun: "甲子", xunStem: "戊", kongWang: "戌亥",
		zhiFu: "天蓬", zfOrig: 1, zfLand: 8, zhiShi: "休", zsLand: 4,
		palaces: [9]string{
			"壬戊乙心死蛇", "己乙丙英伤虎", "庚丙壬任开天",
			"丙丁戊冲休地", "癸癸癸---", "辛壬己柱景阴",
			"乙辛丁芮杜合", "戊庚辛蓬惊符", "丁己庚辅生玄",
		},
	},
	{
		when: "2024-12-28 12:00", note: "冬至上元第3天 (接气后)",
		pillars: "甲辰 丙子 丙寅 甲午", yinYang: "阳", ju: 1,
		xun: "甲午", xunStem: "辛", kongWang: "辰巳",
		zhiFu: "天辅", zfOrig: 4, zfLand: 4, zhiShi: "杜", zsLand: 4,
		palaces: [9]string{
			"戊戊戊蓬休玄", "己己己芮死阴", "庚庚庚冲伤天",
			"辛辛辛辅杜符", "壬壬壬---", "癸癸癸心开虎",
			"丁丁丁柱惊合", "丙丙丙任生地", "乙乙乙英景蛇",
		},
	},
	{
		when: "2033-06-14 12:00", note: "闰芒种上元第3天 (lead 9, 含首尾计数)",
		pillars: "癸丑 戊午 丙申 甲午", yinYang: "阳", ju: 6,
		xun: "甲午", xunStem: "辛", kongWang: "辰巳",
		zhiFu: "天英", zfOrig: 9, zfLand: 9, zhiShi: "景", zsLand: 9,
		palaces: [9]string{
			"壬壬壬蓬休虎", "癸癸癸芮死蛇", "丁丁丁冲伤地",
			"丙丙丙辅杜天", "乙乙乙---", "戊戊戊心开合",
			"己己己柱惊阴", "庚庚庚任生玄", "辛辛辛英景符",
		},
	},
	{
		when: "2027-12-18 12:00", note: "闰大雪中元第3天",
		pillars: "丁未 壬子 辛未 甲午", yinYang: "阴", ju: 7,
		xun: "甲午", xunStem: "辛", kongWang: "辰巳",
		zhiFu: "天辅", zfOrig: 4, zfLand: 4, zhiShi: "杜", zsLand: 4,
		palaces: [9]string{
			"丁丁丁蓬休合", "癸癸癸芮死地", "壬壬壬冲伤蛇",
			"辛辛辛辅杜符", "庚庚庚---", "己己己心开虎",
			"戊戊戊柱惊玄", "乙乙乙任生阴", "丙丙丙英景天",
		},
	},
	{
		when: "1995-12-15 12:00", note: "闰大雪上元第2天 (lead 9, 含首尾计数)",
		pillars: "乙亥 戊子 庚辰 壬午", yinYang: "阴", ju: 4,
		xun: "甲戌", xunStem: "己", kongWang: "申酉",
		zhiFu: "天冲", zfOrig: 3, zfLand: 9, zhiShi: "伤", zsLand: 4,
		palaces: [9]string{
			"丁辛丙柱开虎", "戊庚壬辅景天", "辛己癸蓬生阴",
			"癸戊己任伤蛇", "乙乙乙---", "庚丙丁芮惊玄",
			"壬丁庚英死地", "丙癸辛心休合", "己壬戊冲杜符",
		},
	},
	{
		when: "2024-12-10 04:00", note: "大雪下元第5天 (阴遁, 旬首遁干落中五, 甲寅时伏吟: 值符值使皆落中五)",
		pillars: "甲辰 丙子 戊申 甲寅", yinYang: "阴", ju: 1,
		xun: "甲寅", xunStem: "癸", kongWang: "子丑",
		zhiFu: "天禽", zfOrig: 5, zfLand: 5, zhiShi: "死", zsLand: 5,
		palaces: [9]string{
			"戊戊戊蓬休玄", "乙乙乙芮死符", "丙丙丙冲伤合",
			"丁丁丁辅杜阴", "癸癸癸---", "壬壬壬心开地",
			"辛辛辛柱惊天", "庚庚庚任生虎", "己己己英景蛇",
		},
	},
	{
		when: "2024-12-25 12:00", note: "闰大雪下元第5天 (接气末日, 次日甲子换冬至上元)",
		pillars: "甲辰 丙子 癸亥 戊午", yinYang: "阴", ju: 1,
		xun: "甲寅", xunStem: "癸", kongWang: "子丑",
		zhiFu: "天禽", zfOrig: 5, zfLand: 1, zhiShi: "死", zsLand: 1,
		palaces: [9]string{
			"乙戊乙芮死符", "丙乙丙冲伤合", "壬丙壬心开地",
			"戊丁戊蓬休玄", "癸癸癸---", "己壬己英景蛇",
			"丁辛丁辅杜阴", "辛庚辛柱惊天", "庚己庚任生虎",
		},
	},
	{
		when: "2026-01-14 01:03", note: "时家飞盘置闰: 九星数飞, 天禽落实宫, 九神含太常, 值使走满九步回本位",
		style:   enum.StyleFly,
		pillars: "乙巳 己丑 戊子 癸丑", yinYang: "阳", ju: 8,
		xun: "甲辰", xunStem: "壬", kongWang: "寅卯",
		zhiFu: "天冲", zfOrig: 3, zfLand: 4, zhiShi: "伤", zsLand: 3,
		palaces: [9]string{
			"己庚庚英休玄", "庚辛辛蓬死地", "辛壬壬芮伤天",
			"壬癸癸冲杜符", "癸丁丁辅-蛇", "丁丙丙禽开阴",
			"丙乙乙心惊合", "乙戊戊柱生虎", "戊己己任景常",
		},
	},
	{
		when: "2024-12-26 12:00", note: "日家转盘 (拆补): 冬至上元阳一, 甲子日全盘伏吟",
		method: enum.MethodDay, juRule: enum.JuRuleChaiBu,
		pillars: "甲辰 丙子 甲子 庚午", yinYang: "阳", ju: 1,
		xun: "甲子", xunStem: "戊", kongWang: "戌亥",
		zhiFu: "天蓬", zfOrig: 1, zfLand: 1, zhiShi: "休", zsLand: 1,
		palaces: [9]string{
			"戊戊戊蓬休符", "己己己芮死玄", "庚庚庚冲伤阴",
			"辛辛辛辅杜合", "壬壬壬---", "癸癸癸心开天",
			"丁丁丁柱惊地", "丙丙丙任生蛇", "乙乙乙英景虎",
		},
	},
	{
		when: "2024-12-25 12:00", note: "日家转盘 (拆补): 冬至下元阳四 — 日家与时家共用节气三元",
		method: enum.MethodDay, juRule: enum.JuRuleChaiBu,
		pillars: "甲辰 丙子 癸亥 戊午", yinYang: "阳", ju: 4,
		xun: "甲寅", xunStem: "癸", kongWang: "子丑",
		zhiFu: "天英", zfOrig: 9, zfLand: 9, zhiShi: "景", zsLand: 9,
		palaces: [9]string{
			"丁丁丁蓬休虎", "丙丙丙芮死蛇", "乙乙乙冲伤地",
			"戊戊戊辅杜天", "己己己---", "庚庚庚心开合",
			"辛辛辛柱惊阴", "壬壬壬任生玄", "癸癸癸英景符",
		},
	},
	{
		when: "2024-12-25 12:00", note: "月家转盘: 甲辰年子月 阴遁四局 (寅月起五逆行), 主柱月柱",
		method:  enum.MethodMonth,
		pillars: "甲辰 丙子 癸亥 戊午", yinYang: "阴", ju: 4,
		xun: "甲戌", xunStem: "己", kongWang: "申酉",
		zhiFu: "天冲", zfOrig: 3, zfLand: 6, zhiShi: "伤", zsLand: 1,
		palaces: [9]string{
			"戊辛己辅伤天", "辛庚丙蓬开阴", "庚己壬芮景玄",
			"丁戊庚柱死虎", "乙乙乙---", "己丙癸冲生符",
			"癸丁辛任休蛇", "壬癸戊英杜地", "丙壬丁心惊合",
		},
	},
	{
		when: "2024-12-25 12:00", note: "年家转盘: 甲辰年 阴遁七局 (下元定局), 主柱年柱, 全盘伏吟",
		method:  enum.MethodYear,
		pillars: "甲辰 丙子 癸亥 戊午", yinYang: "阴", ju: 7,
		xun: "甲辰", xunStem: "壬", kongWang: "寅卯",
		zhiFu: "天冲", zfOrig: 3, zfLand: 3, zhiShi: "伤", zsLand: 3,
		palaces: [9]string{
			"丁丁丁蓬休阴", "癸癸癸芮死玄", "壬壬壬冲伤符",
			"辛辛辛辅杜天", "庚庚庚---", "己己己心开合",
			"戊戊戊柱惊虎", "乙乙乙任生蛇", "丙丙丙英景地",
		},
	},
}

var goldenStarNames = map[rune]string{
	'蓬': "天蓬", '任': "天任", '冲': "天冲", '辅': "天辅",
	'英': "天英", '芮': "禽芮", '柱': "天柱", '心': "天心",
}

// goldenFlyStarNames differs in the 禽/芮 pair: fly charts place 天禽
// on its real palace, so 芮 is plain 天芮 (no merge marker).
var goldenFlyStarNames = map[rune]string{
	'蓬': "天蓬", '任': "天任", '冲': "天冲", '辅': "天辅",
	'英': "天英", '芮': "天芮", '柱': "天柱", '心': "天心",
	'禽': "天禽",
}

var goldenGodNames = map[rune]string{
	'符': "值符", '蛇': "腾蛇", '阴': "太阴", '合': "六合",
	'虎': "白虎", '玄': "玄武", '地': "九地", '天': "九天",
	'常': "太常", // fly style: ninth god
}

// TestGoldenZhiRunCharts replays authoritative 置闰 charts and checks
// every layer of the layout against the reference output.
func TestGoldenZhiRunCharts(t *testing.T) {
	for _, g := range goldenCharts {
		t.Run(g.when, func(t *testing.T) {
			var y, mo, d, h, mi int
			if _, err := fmt.Sscanf(g.when, "%d-%d-%d %d:%d", &y, &mo, &d, &h, &mi); err != nil {
				t.Fatalf("parse %q: %v", g.when, err)
			}
			c := From(solarTime(t, y, mo, d, h, mi, 0),
				WithMethod(g.method), WithStyle(g.style), WithJuRule(g.juRule))

			pillars := fmt.Sprintf("%s %s %s %s", c.Year(), c.Month(), c.Day(), c.Hour())
			if pillars != g.pillars {
				t.Errorf("pillars: got %s, want %s", pillars, g.pillars)
			}
			if c.YinYang().Name() != g.yinYang || c.Ju() != g.ju {
				t.Errorf("ju: got %s遁%d局, want %s遁%d局 (%s)",
					c.YinYang().Name(), c.Ju(), g.yinYang, g.ju, g.note)
			}
			if got := c.Lead().Ten().Name(); got != g.xun {
				t.Errorf("xun: got %s, want %s", got, g.xun)
			}
			if got := c.XunShou().Name(); got != g.xunStem {
				t.Errorf("xunshou: got %s, want %s", got, g.xunStem)
			}
			kw := c.KongWang()
			if got := kw[0].Name() + kw[1].Name(); got != g.kongWang {
				t.Errorf("kongwang: got %s, want %s", got, g.kongWang)
			}

			zf, zs := c.ZhiFu(), c.ZhiShi()
			if zf.Star.Name() != g.zhiFu || zf.OriginalPalace != g.zfOrig || zf.Palace != g.zfLand {
				t.Errorf("zhifu: got %s %d→%d, want %s %d→%d",
					zf.Star.Name(), zf.OriginalPalace, zf.Palace, g.zhiFu, g.zfOrig, g.zfLand)
			}
			if zs.Door.Name() != g.zhiShi || zs.Palace != g.zsLand {
				t.Errorf("zhishi: got %s →%d, want %s →%d",
					zs.Door.Name(), zs.Palace, g.zhiShi, g.zsLand)
			}

			for n := uint8(1); n <= 9; n++ {
				p := c.Palace(n)
				spec := []rune(g.palaces[n-1])
				if len(spec) != 6 {
					t.Fatalf("palace %d spec %q: want 6 runes", n, g.palaces[n-1])
				}
				if got, want := p.HeavenStem.Name(), string(spec[0]); got != want {
					t.Errorf("palace %d heaven: got %s, want %s", n, got, want)
				}
				if got, want := p.EarthStem.Name(), string(spec[1]); got != want {
					t.Errorf("palace %d earth: got %s, want %s", n, got, want)
				}
				if got, want := p.HiddenStem.Name(), string(spec[2]); got != want {
					t.Errorf("palace %d hidden: got %s, want %s", n, got, want)
				}
				starTable := goldenStarNames
				if g.style == enum.StyleFly {
					starTable = goldenFlyStarNames
				}
				assertSlot(t, n, "star", spec[3], p.StarSet, func() string { return p.Star.Name() }, starTable)
				assertSlot(t, n, "door", spec[4], p.DoorSet, func() string { return p.Door.Name() }, nil)
				assertSlot(t, n, "god", spec[5], p.GodSet, func() string { return p.God.Name() }, goldenGodNames)
			}
		})
	}
}

// assertSlot checks one palace slot against its golden rune: '-' means
// the slot must be unset; otherwise the slot must be set and its name
// must match (via the rune table when provided, verbatim otherwise).
func assertSlot(t *testing.T, n uint8, slot string, want rune, set bool, name func() string, table map[rune]string) {
	t.Helper()
	if want == '-' {
		if set {
			t.Errorf("palace %d %s: got %s, want empty", n, slot, name())
		}
		return
	}
	if !set {
		t.Errorf("palace %d %s: empty, want %c", n, slot, want)
		return
	}
	wantName := string(want)
	if table != nil {
		wantName = table[want]
	}
	if got := name(); got != wantName {
		t.Errorf("palace %d %s: got %s, want %s", n, slot, got, wantName)
	}
}
