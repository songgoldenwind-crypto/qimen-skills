// Command qimen-json exposes a small, source-labelled JSON boundary for atopx/qimen.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/atopx/qimen"
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

type palaceJSON struct {
	Name       string  `json:"name"`
	SkyStem    string  `json:"sky_stem"`
	EarthStem  string  `json:"earth_stem"`
	HiddenStem string  `json:"hidden_stem"`
	StemPair   string  `json:"stem_pair"`
	Door       *string `json:"door"`
	Star       *string `json:"star"`
	Spirit     *string `json:"spirit"`
	LeadVoid   bool    `json:"lead_void"`
	HourVoid   *bool   `json:"hour_void"`
	DayVoid    *bool   `json:"day_void"`
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}

func pointer(value string) *string { return &value }

func displayPalace(real uint8, style string) uint8 {
	if style == "rotate" && real == 5 {
		return 2
	}
	return real
}

func main() {
	input := flag.String("datetime", "", "local civil time: YYYY-MM-DDTHH:MM:SS")
	family := flag.String("family", "time", "time, day, month or year")
	style := flag.String("style", "rotate", "rotate or fly")
	rule := flag.String("method", "zhirun", "zhirun or chaibu (time/day only)")
	flag.Parse()

	when, err := time.Parse("2006-01-02T15:04:05", *input)
	if err != nil {
		fail("invalid --datetime: %v", err)
	}
	methodMap := map[string]enum.Method{
		"time": enum.MethodTime, "day": enum.MethodDay,
		"month": enum.MethodMonth, "year": enum.MethodYear,
	}
	selectedMethod, ok := methodMap[*family]
	if !ok {
		fail("unsupported family %q", *family)
	}
	styleMap := map[string]enum.Style{"rotate": enum.StyleRotate, "fly": enum.StyleFly}
	selectedStyle, ok := styleMap[*style]
	if !ok {
		fail("unsupported style %q", *style)
	}
	ruleMap := map[string]enum.JuRule{"zhirun": enum.JuRuleZhiRun, "chaibu": enum.JuRuleChaiBu}
	selectedRule, ok := ruleMap[*rule]
	if !ok {
		fail("unsupported method %q", *rule)
	}
	if (*family == "month" || *family == "year") && *rule != "zhirun" {
		fail("month/year charts do not use a zhirun/chaibu choice")
	}

	solar, err := almanac.SolarTimeOf(when.Year(), int(when.Month()), when.Day(), when.Hour(), when.Minute(), when.Second())
	if err != nil {
		fail("invalid local time: %v", err)
	}
	chart := qimen.From(solar, qimen.WithMethod(selectedMethod), qimen.WithStyle(selectedStyle), qimen.WithJuRule(selectedRule))
	void := chart.KongWang()
	voidNames := []string{void[0].Name(), void[1].Name()}
	palaces := make(map[string]palaceJSON, 9)
	for n := uint8(1); n <= 9; n++ {
		p := chart.Palace(n)
		item := palaceJSON{
			Name: p.Name, SkyStem: p.HeavenStem.Name(), EarthStem: p.EarthStem.Name(),
			HiddenStem: p.HiddenStem.Name(),
			StemPair:   p.HeavenStem.Name() + "+" + p.EarthStem.Name(),
			LeadVoid:   false,
		}
		if p.DoorSet {
			item.Door = pointer(p.Door.Name())
		}
		if p.StarSet {
			item.Star = pointer(p.Star.Name())
		}
		if p.GodSet {
			item.Spirit = pointer(p.God.Name())
		}
		for _, branch := range p.Branches {
			if branch == void[0] || branch == void[1] {
				item.LeadVoid = true
			}
		}
		if *family == "time" {
			item.HourVoid = &item.LeadVoid
		}
		if *family == "day" {
			item.DayVoid = &item.LeadVoid
		}
		palaces[fmt.Sprint(n)] = item
	}
	result := map[string]any{
		"raw": map[string]any{
			"干支":  chart.Year().Name() + "年" + chart.Month().Name() + "月" + chart.Day().Name() + "日" + chart.Hour().Name() + "时",
			"阴阳遁": chart.YinYang().Name(), "局数": chart.Ju(), "三元": chart.Yuan().Name(),
			"节气": chart.Term().Name(), "用局节气": chart.JuTerm().Name(),
			"旬首": chart.XunShou().Name(), "旬空": voidNames,
			"值符": map[string]any{"星": chart.ZhiFu().Star.Name(), "原宫": chart.ZhiFu().OriginalPalace, "落宫": chart.ZhiFu().Palace, "盘面宫": displayPalace(chart.ZhiFu().Palace, *style)},
			"值使": map[string]any{"门": chart.ZhiShi().Door.Name(), "原宫": chart.ZhiShi().OriginalPalace, "落宫": chart.ZhiShi().Palace, "盘面宫": displayPalace(chart.ZhiShi().Palace, *style)},
		},
		"palaces": palaces,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(result); err != nil {
		fail("encode chart: %v", err)
	}
}
