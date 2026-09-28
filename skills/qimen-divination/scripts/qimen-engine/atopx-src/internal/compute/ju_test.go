package compute

import (
	"testing"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

// termIn returns the in-effect term for a noon instant.
func termIn(t *testing.T, y, mo, d int) almanac.Term {
	t.Helper()
	st, err := almanac.SolarTimeOf(y, mo, d, 12, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	return almanac.TermOfSolarTime(st)
}

// TestYearJu anchors the 年家局: fixed per 60-year 元 (上1 中4 下7,
// anchored at the 上元甲子 year 1864), switching years at 立春.
func TestYearJu(t *testing.T) {
	cases := []struct {
		y, m, d int
		ju      uint8
		yuan    enum.Yuan
	}{
		{1864, 6, 1, 1, enum.YuanUpper},
		{1923, 6, 1, 1, enum.YuanUpper},  // last year of 上元
		{1924, 6, 1, 4, enum.YuanMiddle}, // 中元甲子
		{1983, 6, 1, 4, enum.YuanMiddle},
		{1984, 6, 1, 7, enum.YuanLower}, // 下元甲子
		{2024, 12, 25, 7, enum.YuanLower},
		{2043, 6, 1, 7, enum.YuanLower},
		{2044, 6, 1, 1, enum.YuanUpper},  // next 上元甲子
		{2024, 1, 10, 7, enum.YuanLower}, // before 立春 → still 癸卯/下元
	}
	for _, c := range cases {
		ju, yuan := YearJu(termIn(t, c.y, c.m, c.d))
		if ju != c.ju || yuan != c.yuan {
			t.Errorf("YearJu(%d-%02d-%02d): got %d局/%s, want %d局/%s",
				c.y, c.m, c.d, ju, yuan.Name(), c.ju, c.yuan.Name())
		}
	}
}

// TestMonthJu anchors the 月家局 mnemonic: 子午卯酉 → 寅月八局,
// 辰戌丑未 → 五局, 寅申巳亥 → 二局, retreating by one per month.
func TestMonthJu(t *testing.T) {
	cases := []struct {
		yearBranch, monthBranch almanac.Branch
		ju                      uint8
		yuan                    enum.Yuan
	}{
		{almanac.Zi, almanac.Yin_, 8, enum.YuanUpper}, // 子年寅月 八局
		{almanac.Zi, almanac.Mao, 7, enum.YuanUpper},  // 子年卯月 七局
		{almanac.Chen, almanac.Yin_, 5, enum.YuanLower},
		{almanac.Chen, almanac.Zi, 4, enum.YuanLower}, // 甲辰年子月 四局 (golden)
		{almanac.Yin_, almanac.Yin_, 2, enum.YuanMiddle},
		{almanac.Wu_, almanac.Chou, 6, enum.YuanUpper}, // 午年丑月 (12th) 8-11→6
	}
	for _, c := range cases {
		ju, yuan := MonthJu(c.yearBranch, c.monthBranch)
		if ju != c.ju || yuan != c.yuan {
			t.Errorf("MonthJu(%s年%s月): got %d局/%s, want %d局/%s",
				c.yearBranch.Name(), c.monthBranch.Name(), ju, yuan.Name(), c.ju, c.yuan.Name())
		}
	}
}
