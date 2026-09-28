package almanac

import "testing"

// TestPillarsBoundary covers cross-day, cross-month (节气), and 23 点换日 edges.
//
// Reference values verified against 寿星天文历 + 三命通会 公式.
func TestPillarsBoundary(t *testing.T) {
	cases := []struct {
		desc      string
		y, m, d   int
		h, mi, s  int
		wantYear  string
		wantMonth string
		wantDay   string
		wantHour  string
	}{
		// 23 点换日: 2025-05-05 22:59 → 旧日; 23:00 → 次日 子时
		// 2025-05-05 14:57 立夏后, 月柱进 辛巳; 年柱 乙巳 (立春后)
		{"2025-05-05 22:59 boundary", 2025, 5, 5, 22, 59, 0, "乙巳", "辛巳", "甲戌", "乙亥"},
		{"2025-05-05 23:00 boundary", 2025, 5, 5, 23, 0, 0, "乙巳", "辛巳", "乙亥", "丙子"},
		// 立春 2025: 2025-02-03 22:10:13. 立春前 年柱 甲辰, 月柱 丁丑
		{"立春前 2025-02-03 22:00", 2025, 2, 3, 22, 0, 0, "甲辰", "丁丑", "癸卯", "癸亥"},
		// 立春后 年柱 乙巳, 月柱 戊寅
		{"立春后 2025-02-04 00:00", 2025, 2, 4, 0, 0, 0, "乙巳", "戊寅", "甲辰", "甲子"},
	}
	for _, c := range cases {
		st, err := SolarTimeOf(c.y, c.m, c.d, c.h, c.mi, c.s)
		if err != nil {
			t.Fatal(err)
		}
		p := st.Pillars()
		check := func(field, got, want string) {
			if got != want {
				t.Errorf("%s %s: got %q, want %q", c.desc, field, got, want)
			}
		}
		check("year", p.Year.Name(), c.wantYear)
		check("month", p.Month.Name(), c.wantMonth)
		check("day", p.Day.Name(), c.wantDay)
		check("hour", p.Hour.Name(), c.wantHour)
	}
}

// TestStemTenStarReflexive covers the 比肩 baseline: any stem paired with
// itself produces 比肩.
func TestStemTenStarReflexive(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := Stem(i)
		got := s.TenStarOf(s).Name()
		if got != "比肩" {
			t.Errorf("Stem(%d).TenStarOf(self) = %q, want 比肩", i, got)
		}
	}
}

// TestStemTerrainOf samples representative 长生十二宫 positions verified
// against the canonical 长生 table.
func TestStemTerrainOf(t *testing.T) {
	cases := []struct {
		stem   Stem
		branch Branch
		want   string
	}{
		{Jia, Hai, "长生"},   // 甲长生在亥
		{Jia, Wu_, "死"},    // 甲死于午
		{Yi, Wu_, "长生"},    // 乙长生在午
		{Bing, Yin_, "长生"}, // 丙长生在寅
		{Gui, Mao, "长生"},   // 癸长生在卯
	}
	for _, c := range cases {
		got := c.stem.TerrainOf(c.branch).Name()
		if got != c.want {
			t.Errorf("Stem(%s).TerrainOf(Branch(%s)) = %q, want %q",
				c.stem.Name(), c.branch.Name(), got, c.want)
		}
	}
}
