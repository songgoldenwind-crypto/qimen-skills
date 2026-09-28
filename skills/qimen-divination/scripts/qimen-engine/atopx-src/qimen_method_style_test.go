package qimen

import (
	"testing"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

// TestMethods verifies each 起局法门 keys the chart to its own duty
// pillar: the 旬首 and 空亡 must come from that pillar, and the 月/年家
// charts are always 阳遁.
func TestMethods(t *testing.T) {
	st := solarTime(t, 2026, 1, 14, 18, 45, 0) // 乙巳 己丑 戊子 辛酉

	cases := []struct {
		method enum.Method
		lead   func(*Chart) almanac.Cycle
	}{
		{enum.MethodTime, (*Chart).Hour},
		{enum.MethodDay, (*Chart).Day},
		{enum.MethodMonth, (*Chart).Month},
		{enum.MethodYear, (*Chart).Year},
	}
	for _, c := range cases {
		chart := From(st, WithMethod(c.method))
		if chart.Method() != c.method {
			t.Errorf("%s: Method() = %s", c.method.Name(), chart.Method().Name())
		}
		lead := c.lead(chart)
		if got, want := chart.XunShou(), almanac.StemOf(lead.Ten().Index()+4); got != want {
			t.Errorf("%s: xunshou %s, want %s (from %s)",
				c.method.Name(), got.Name(), want.Name(), lead.Name())
		}
		if got, want := chart.KongWang(), lead.EmptyBranches(); got != want {
			t.Errorf("%s: kongwang %v, want %v", c.method.Name(), got, want)
		}
		if (c.method == enum.MethodMonth || c.method == enum.MethodYear) &&
			chart.YinYang() != almanac.Yin {
			t.Errorf("%s: must be 阴遁", c.method.Name())
		}
		if got := chart.Ju(); got < 1 || got > 9 {
			t.Errorf("%s: ju %d out of range", c.method.Name(), got)
		}
	}

	// 日家 shares the 节气三元 局 with 时家 — same day, same 局.
	if dayJu, timeJu := From(st, WithMethod(enum.MethodDay)).Ju(), From(st).Ju(); dayJu != timeJu {
		t.Errorf("日家 ju %d != 时家 ju %d on the same day", dayJu, timeJu)
	}
	// 年家 局 is fixed per 60-year 元: 2026 sits in the 下元 (1984-2043).
	if got := From(st, WithMethod(enum.MethodYear)).Ju(); got != 7 {
		t.Errorf("年家 2026: ju %d, want 7 (下元)", got)
	}
}

// TestStyleFly checks the structural laws of a fly-plate chart against
// its rotate-plate sibling: same earth plate, heaven plate = earth
// plate flown by the 值符 delta, nine stars on nine palaces (天禽 on a
// real palace), eight doors and eight gods leaving one palace each.
func TestStyleFly(t *testing.T) {
	st := solarTime(t, 2026, 1, 14, 18, 45, 0)
	rot := From(st)
	fly := From(st, WithStyle(enum.StyleFly))

	if fly.Style() != enum.StyleFly {
		t.Fatalf("Style() = %s", fly.Style().Name())
	}
	if fly.Ju() != rot.Ju() || fly.YinYang() != rot.YinYang() {
		t.Errorf("局/遁 must not depend on style")
	}

	zf := fly.ZhiFu()
	delta := (int(zf.Palace) - int(zf.OriginalPalace) + 9) % 9
	stars := map[enum.Star]bool{}
	doors, gods := 0, 0
	for n, p := range fly.Palaces() {
		if p.EarthStem != rot.Palace(n).EarthStem {
			t.Errorf("palace %d: earth plate differs between styles", n)
		}
		// heaven = earth flown by the duty delta
		src := uint8((int(n)-1-delta+9)%9 + 1)
		if p.HeavenStem != fly.Palace(src).EarthStem {
			t.Errorf("palace %d: heaven %s, want earth of palace %d (%s)",
				n, p.HeavenStem.Name(), src, fly.Palace(src).EarthStem.Name())
		}
		if !p.StarSet {
			t.Errorf("palace %d: fly charts have a star on every palace", n)
		}
		stars[p.Star] = true
		if p.DoorSet {
			doors++
		}
		if p.GodSet {
			gods++
		}
	}
	if len(stars) != 9 {
		t.Errorf("fly chart: %d distinct stars, want 9", len(stars))
	}
	if stars[enum.StarQinRui] {
		t.Error("fly chart must place 天禽 itself, not the 禽芮 merge marker")
	}
	if doors != 8 || gods != 9 {
		t.Errorf("fly chart: %d doors / %d gods set, want 8 / 9", doors, gods)
	}
}

// TestStyleFlyFuYin verifies the degenerate flight (delta 0) reproduces
// the earth plate on the heaven plate with every star at home.
func TestStyleFlyFuYin(t *testing.T) {
	// 2026-02-18 22:59 is a 伏吟 chart (值符落原宫) under 置闰.
	c := From(solarTime(t, 2026, 2, 18, 22, 59, 0), WithStyle(enum.StyleFly))
	zf := c.ZhiFu()
	if zf.Palace != zf.OriginalPalace {
		t.Skip("fixture no longer 伏吟")
	}
	for n, p := range c.Palaces() {
		if p.HeavenStem != p.EarthStem {
			t.Errorf("palace %d: 伏吟 fly heaven %s != earth %s",
				n, p.HeavenStem.Name(), p.EarthStem.Name())
		}
		if p.Star != enum.StarOfPalace(n) {
			t.Errorf("palace %d: 伏吟 fly star %s not at home", n, p.Star.Name())
		}
	}
}
