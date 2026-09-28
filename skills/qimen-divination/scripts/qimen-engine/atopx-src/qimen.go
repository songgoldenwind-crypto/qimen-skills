package qimen

import (
	"iter"
	"time"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/hexagram"
	"github.com/atopx/qimen/internal/compute"
	"github.com/atopx/qimen/internal/tables"
	"github.com/atopx/qimen/palace"
	"github.com/atopx/qimen/pattern"
	"github.com/atopx/qimen/plate"
	"github.com/atopx/qimen/shensha"
	"github.com/atopx/qimen/terrain"
)

// Chart is the result of one qimen 起局. Once built, all fields are
// read-only and safe to share across goroutines.
type Chart struct {
	method enum.Method
	style  enum.Style
	juRule enum.JuRule

	// Pillars + solar context
	solar    almanac.SolarTime
	pillars  almanac.Pillars
	lead     almanac.Cycle // duty pillar (主柱) per the method
	term     almanac.Term
	juTerm   almanac.Term // 用局节气 (== term under 拆补)
	yinYang  almanac.YinYang
	ju       uint8
	yuan     enum.Yuan
	xunShou  almanac.Stem
	zhiFu    DutyStar
	zhiShi   DutyDoor
	kongWang [2]almanac.Branch
	lunarDay almanac.LunarDay // cached at build for O(1) LunarDay() access

	// Palaces (1..9; index = palace number - 1).
	// Stored by value to avoid 9 heap allocations per chart.
	palaces [9]palace.Palace
}

// config holds optional construction parameters set via Option.
type config struct {
	method enum.Method
	style  enum.Style
	juRule enum.JuRule
}

// Option configures a Chart constructor.
type Option func(*config)

// WithMethod selects the 起局 method: 时家 (default), 日家, 月家 or
// 年家. The method picks the duty pillar (时/日/月/年柱) and the 局
// source — 时家 and 日家 share the per-day 节气三元 局; 月/年家 use
// the year/month calendars and are always 阴遁. See [enum.Method].
func WithMethod(m enum.Method) Option {
	return func(c *config) { c.method = m }
}

// WithStyle selects the chart style: 转盘 (default, LuoShu ring
// rotation with the center stem 寄坤) or 飞盘 (palace-number flying
// with the center as a regular stop and 天禽 on its real palace).
func WithStyle(s enum.Style) Option {
	return func(c *config) { c.style = s }
}

// WithJuRule selects how the 时家 / 日家 用局节气 is resolved: 置闰
// (default) keeps 符头 aligned with the solstices by intercalating
// 芒种 / 大雪, 拆补 keys the 局 to the astronomical term in effect.
// Ignored by the 月/年家 methods. See [enum.JuRule].
func WithJuRule(r enum.JuRule) Option {
	return func(c *config) { c.juRule = r }
}

// DutyStar is the 值符 entry: the 九星 currently acting as 值符, its
// originating palace, and the (real, possibly center) palace it landed
// on.
type DutyStar struct {
	Star           enum.Star
	OriginalPalace uint8
	Palace         uint8
}

// DutyDoor is the 值使 entry: the 八门 currently acting as 值使, its
// originating palace, and the (real, possibly center) palace it landed
// on.
type DutyDoor struct {
	Door           enum.Door
	OriginalPalace uint8
	Palace         uint8
}

// defaultConfig returns the default construction options:
// 时家 / 转盘 / 置闰.
func defaultConfig() config {
	return config{
		method: enum.MethodTime,
		style:  enum.StyleRotate,
		juRule: enum.JuRuleZhiRun,
	}
}

// New builds a chart for the current instant (UTC+8 wall clock) with
// default options (时家 / 转盘 / 置闰).
func New() *Chart {
	return build(almanac.Now(), defaultConfig())
}

// From builds a chart from a [almanac.SolarTime]. Every Method / Style
// / JuRule combination is implemented, so construction is total — no
// error to handle.
func From(t almanac.SolarTime, opts ...Option) *Chart {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return build(t, cfg)
}

// FromTime builds a chart from a standard library [time.Time], using the
// time's own wall clock — its Location is respected as-is, never
// converted. Classical qimen casts charts in the local civil time of the
// event; convert explicitly (t.In(loc)) before calling when the value's
// zone is not the intended one.
func FromTime(t time.Time, opts ...Option) (*Chart, error) {
	st, err := almanac.SolarTimeFromTime(t)
	if err != nil {
		return nil, err // already wraps ErrInvalidTime
	}
	return From(st, opts...), nil
}

// FromTimestamp builds a chart from a Unix-seconds timestamp, interpreted
// as a UTC+8 (China Standard Time) wall clock. For another zone use
// FromTime(time.Unix(unix, 0).In(loc)).
func FromTimestamp(unix int64, opts ...Option) (*Chart, error) {
	st, err := almanac.SolarTimeFromUnix(unix)
	if err != nil {
		return nil, err // already wraps ErrInvalidTime
	}
	return From(st, opts...), nil
}

// build is the shared chart-construction routine, total over its
// closed-domain inputs.
func build(t almanac.SolarTime, cfg config) *Chart {
	term := t.Term()
	pillars := almanac.PillarsAt(t, term)

	// The method picks the duty pillar (主柱) and the 局 source. 时家
	// and 日家 share the 节气三元 source (the 局 is a per-day fact;
	// 置闰 may shift the 用局节气 to keep 符头 aligned with the
	// solstices); 月/年家 use their year/month calendars and are always 阴遁.
	juTerm := term
	var lead almanac.Cycle
	var yinYang almanac.YinYang
	var ju uint8
	var yuan enum.Yuan
	switch cfg.method {
	case enum.MethodMonth:
		lead = pillars.Month
		yinYang = almanac.Yin
		ju, yuan = compute.MonthJu(pillars.Year.Branch(), pillars.Month.Branch())
	case enum.MethodYear:
		lead = pillars.Year
		yinYang = almanac.Yin
		ju, yuan = compute.YearJu(term)
	default: // enum.MethodTime / enum.MethodDay
		lead = pillars.Hour
		if cfg.method == enum.MethodDay {
			lead = pillars.Day
		}
		if cfg.juRule == enum.JuRuleZhiRun {
			juTerm = compute.ZhiRunTerm(almanac.DayNumber(t), term)
		}
		yinYang = juTerm.YinYang()
		yuan = compute.Yuan(pillars.Day)
		ju = compute.Ju(juTerm, yuan)
	}

	earth := plate.BuildEarth(yinYang, ju)
	xunShou := compute.XunShou(lead)
	leadStem := lead.Stem()

	// 值符原宫 / 落宫: the earth plate always contains the 旬首 stem
	// (one of 戊..癸), possibly in the center palace. Duty records keep
	// the real palaces (which may be 5).
	zhiFuOriginalPalace := uint8(2)
	if p, ok := plate.FindStem(&earth, xunShou); ok {
		zhiFuOriginalPalace = p
	}
	zhiFuPalace := plate.FindHourStem(&earth, leadStem, zhiFuOriginalPalace)

	// Lay the rotating plates. 天盘 follows the 值符, 暗干 follows the
	// 值使 (same shift as the door plate). The 值使 marches from the
	// REAL origin palace while its door is the projected home's door.
	var heaven, hidden plate.StemPlate
	var stars plate.StarPlate
	var doors plate.DoorPlate
	var gods plate.GodPlate
	var zhiShiPalace uint8
	if cfg.style == enum.StyleFly {
		// 飞盘: palace-number flying through the full 1..9 sequence;
		// the center is a regular stop, nothing is projected.
		delta := plate.FlyDelta(zhiFuOriginalPalace, zhiFuPalace)
		heaven = plate.FlyStems(&earth, delta)
		stars = plate.FlyStars(delta)
		zhiShiPalace = plate.MoveBy(zhiFuOriginalPalace, plate.MarchSteps(lead), yinYang)
		doorDelta := plate.FlyDelta(projectCenter(zhiFuOriginalPalace), zhiShiPalace)
		doors = plate.FlyDoors(doorDelta)
		gods = plate.FlyGods(yinYang, zhiFuPalace)
		hidden = plate.FlyStems(&earth, doorDelta)
	} else {
		// 转盘: rigid LuoShu ring rotation; ring positions use the
		// 寄坤 projection (5 → 2).
		origEff := projectCenter(zhiFuOriginalPalace)
		landEff := projectCenter(zhiFuPalace)
		heaven = plate.RotateStems(&earth, origEff, landEff)
		stars = plate.BuildStar(origEff, landEff)
		doors, zhiShiPalace = plate.BuildDoor(yinYang, origEff, zhiFuOriginalPalace, lead)
		gods = plate.BuildGod(yinYang, landEff)
		hidden = plate.RotateStems(&earth, origEff, projectCenter(zhiShiPalace))
	}

	// 值符 is the home star of the real origin palace (center → 天禽);
	// 值使 is the home door of its projection.
	zhiFuStar := enum.StarOfPalace(zhiFuOriginalPalace)
	zhiShiDoor := enum.DoorOfPalace(projectCenter(zhiFuOriginalPalace))

	kongWang := lead.EmptyBranches()

	// Assemble 9 palaces in place (no heap allocations).
	var palaces [9]palace.Palace
	for i := range 9 {
		n := uint8(i + 1)
		earthStem, _ := earth.Get(n)
		heavenStem := earthStem
		if v, ok := heaven.Get(n); ok {
			heavenStem = v
		}
		hiddenStem := earthStem
		if v, ok := hidden.Get(n); ok {
			hiddenStem = v
		}
		p := &palaces[i]
		p.Number = n
		p.Name = tables.PalaceNames[n]
		p.Direction = almanac.DirectionOfPalace(n)
		p.Branches = compute.BranchesForPalace(n)
		p.EarthStem = earthStem
		p.HeavenStem = heavenStem
		p.HiddenStem = hiddenStem
		// Presence depends on the style: 转盘 leaves the center empty,
		// 飞盘 flies through it — the plate state is the truth.
		p.Star, p.StarSet = stars.Get(n)
		p.Door, p.DoorSet = doors.Get(n)
		p.God, p.GodSet = gods.Get(n)
	}

	// Derived attributes (十神 / 长生 / 六十四卦) — only for non-center
	// palaces (the center has no branches); the hexagram also needs a
	// door for its upper trigram.
	dayStem := pillars.Day.Stem()
	for i := range palaces {
		p := &palaces[i]
		if p.IsCenter() {
			continue
		}
		p.TenStar = dayStem.TenStarOf(p.EarthStem)
		p.Terrain = terrain.Of(dayStem.TerrainOf(p.Branches[0]))
		if p.DoorSet {
			// Hexagram: upper = door's home trigram, lower = palace trigram.
			lower := hexagram.TrigramOfPalace(p.Number)
			upper := hexagram.TrigramOfPalace(p.Door.HomePalace())
			p.Hexagram = hexagram.Of(upper, lower)
			p.HexagramSet = true
		}
	}

	// Pattern detection (build input view, then bucket results per palace).
	pIn := pattern.DetectInput{
		ZhiFuOriginalPalace: zhiFuOriginalPalace,
		ZhiFuPalace:         zhiFuPalace,
		KongWang:            kongWang,
	}
	for i := range palaces {
		p := &palaces[i]
		pIn.EarthStems[i] = p.EarthStem
		pIn.HeavenStems[i] = p.HeavenStem
		pIn.Doors[i] = p.Door
		pIn.DoorsSet[i] = p.DoorSet
		pIn.Gods[i] = p.God
		pIn.GodsSet[i] = p.GodSet
		pIn.Branches[i] = p.Branches
	}
	pats := pattern.AppendAll(make([]pattern.Pattern, 0, 12), &pIn)
	for i, ps := range bucketByPalace(pats, patternPalace) {
		palaces[i].Patterns = ps
	}

	// 神煞 detection.
	ssIn := shensha.DetectInput{
		YearStem:    pillars.Year.Stem(),
		MonthBranch: pillars.Month.Branch(),
		DayStem:     dayStem,
		DayBranch:   pillars.Day.Branch(),
	}
	for i := range 9 {
		ssIn.EarthStems[i] = palaces[i].EarthStem
	}
	// The 10 神煞 kinds yield at most 12 instances (each 地盘 stem is
	// unique, so stem-anchored kinds land in at most one palace each).
	sss := shensha.AppendAll(make([]shensha.ShenSha, 0, 12), &ssIn)
	for i, ss := range bucketByPalace(sss, shenshaPalace) {
		palaces[i].ShenSha = ss
	}

	return &Chart{
		method:   cfg.method,
		style:    cfg.style,
		juRule:   cfg.juRule,
		solar:    t,
		pillars:  pillars,
		lead:     lead,
		term:     term,
		juTerm:   juTerm,
		yinYang:  yinYang,
		ju:       ju,
		yuan:     yuan,
		xunShou:  xunShou,
		zhiFu:    DutyStar{Star: zhiFuStar, OriginalPalace: zhiFuOriginalPalace, Palace: zhiFuPalace},
		zhiShi:   DutyDoor{Door: zhiShiDoor, OriginalPalace: zhiFuOriginalPalace, Palace: zhiShiPalace},
		kongWang: kongWang,
		lunarDay: t.LunarDay(),
		palaces:  palaces,
	}
}

func patternPalace(p pattern.Pattern) uint8 { return p.Palace }
func shenshaPalace(s shensha.ShenSha) uint8 { return s.Palace }

// projectCenter maps the center palace to its 寄宫 (坤 2) for ring
// positions; non-center palaces pass through.
func projectCenter(palace uint8) uint8 {
	if palace == 5 {
		return 2
	}
	return palace
}

// bucketByPalace regroups detection results (each tagged with a palace
// number 1..9) into 9 per-palace sub-slices of the SAME backing array:
// items are stably sorted by palace in place (they arrive nearly sorted,
// so insertion sort is effectively linear), then sliced per palace with
// the capacity clamped so appending to one palace's slice cannot clobber
// a neighbour's. Zero additional allocations.
func bucketByPalace[T any](items []T, palaceOf func(T) uint8) [9][]T {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && palaceOf(items[j]) < palaceOf(items[j-1]); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
	var out [9][]T
	for start := 0; start < len(items); {
		p := palaceOf(items[start])
		end := start + 1
		for end < len(items) && palaceOf(items[end]) == p {
			end++
		}
		out[p-1] = items[start:end:end]
		start = end
	}
	return out
}

// ===================== context accessors =====================

// SolarTime returns the solar instant used to build the chart.
func (c *Chart) SolarTime() almanac.SolarTime { return c.solar }

// LunarDay returns the lunar date of the chart's solar instant
// (cached at build for O(1) access).
func (c *Chart) LunarDay() almanac.LunarDay { return c.lunarDay }

// Method returns the chart's 起局法门 (currently always MethodTime).
func (c *Chart) Method() enum.Method { return c.method }

// Style returns the chart's 盘式 (currently always StyleRotate).
func (c *Chart) Style() enum.Style { return c.style }

// JuRule returns the chart's 定局规则 (置闰 by default).
func (c *Chart) JuRule() enum.JuRule { return c.juRule }

// Lead returns the duty pillar (主柱) the chart is keyed to: 时柱 for
// 时家, 日柱 for 日家, 月柱 for 月家, 年柱 for 年家. The 旬首, 值符 /
// 值使 movement and 空亡 all derive from this pillar.
func (c *Chart) Lead() almanac.Cycle { return c.lead }

// Year returns the 年柱 sixty cycle.
func (c *Chart) Year() almanac.Cycle { return c.pillars.Year }

// Month returns the 月柱 sixty cycle.
func (c *Chart) Month() almanac.Cycle { return c.pillars.Month }

// Day returns the 日柱 sixty cycle.
func (c *Chart) Day() almanac.Cycle { return c.pillars.Day }

// Hour returns the 时柱 sixty cycle.
func (c *Chart) Hour() almanac.Cycle { return c.pillars.Hour }

// Term returns the astronomical solar term in effect at the chart's
// instant. The 局 may be keyed to a different term under 置闰 — see
// [Chart.JuTerm].
func (c *Chart) Term() almanac.Term { return c.term }

// JuTerm returns the 用局节气 — the term whose row of the 局 table the
// chart uses. Equal to [Chart.Term] under 拆补; under 置闰 it may lead
// the astronomical term by one (超神) or trail it by one (接气 / during
// an intercalated 芒种 / 大雪).
func (c *Chart) JuTerm() almanac.Term { return c.juTerm }

// YinYang returns the 阴/阳 遁 (derived from the 用局节气).
func (c *Chart) YinYang() almanac.YinYang { return c.yinYang }

// Ju returns the local 局 number (1..9).
func (c *Chart) Ju() uint8 { return c.ju }

// Yuan returns the 三元 segment.
func (c *Chart) Yuan() enum.Yuan { return c.yuan }

// XunShou returns the 旬首 stem.
func (c *Chart) XunShou() almanac.Stem { return c.xunShou }

// ZhiFu returns the 值符 duty record.
func (c *Chart) ZhiFu() DutyStar { return c.zhiFu }

// ZhiShi returns the 值使 duty record.
func (c *Chart) ZhiShi() DutyDoor { return c.zhiShi }

// KongWang returns the 旬空亡 branch pair.
func (c *Chart) KongWang() [2]almanac.Branch { return c.kongWang }

// ===================== palace access =====================

// Palace returns the palace at number n (1..9), or nil when out of range.
// Returned pointer aliases the chart's internal storage — do not mutate.
func (c *Chart) Palace(n uint8) *palace.Palace {
	if n < 1 || n > 9 {
		return nil
	}
	return &c.palaces[n-1]
}

// Palaces streams the 9 palaces in canonical (1..9) order.
//
// Yields (n, palace) pairs where n is the palace number (1-indexed).
// The yielded pointer aliases the chart's internal storage.
func (c *Chart) Palaces() iter.Seq2[uint8, *palace.Palace] {
	return func(yield func(uint8, *palace.Palace) bool) {
		for i := range 9 {
			if !yield(uint8(i+1), &c.palaces[i]) {
				return
			}
		}
	}
}

// Grid returns the 3×3 display grid [巽离坤; 震中兑; 艮坎乾].
func (c *Chart) Grid() [3][3]*palace.Palace {
	var out [3][3]*palace.Palace
	for row := range 3 {
		for col := range 3 {
			out[row][col] = &c.palaces[tables.Grid[row][col]-1]
		}
	}
	return out
}

// ===================== user-stem queries =====================

// StemPalace returns the palace where a heavenly stem appears in the 地盘.
//
//   - 甲 (idx 0) has no visible position; falls back to the 值符原宫.
//   - Other stems are looked up; center palace stems are mapped to 2 (坤).
func (c *Chart) StemPalace(stem almanac.Stem) uint8 {
	if stem == almanac.Jia {
		return c.zhiFu.OriginalPalace
	}
	for i := range c.palaces {
		p := &c.palaces[i]
		if p.EarthStem == stem {
			if p.Number == 5 {
				return 2
			}
			return p.Number
		}
	}
	return c.zhiFu.OriginalPalace
}

// SelfPalace returns the palace of the user (subject), keyed by 日干.
func (c *Chart) SelfPalace() uint8 { return c.StemPalace(c.pillars.Day.Stem()) }

// OpponentPalace returns the palace of the counterpart / topic, which
// classical qimen places at 值符落宫.
func (c *Chart) OpponentPalace() uint8 { return c.zhiFu.Palace }

// ===================== aggregated streams =====================

// Patterns streams every 格局 across the 9 palaces in palace order.
func (c *Chart) Patterns() iter.Seq[pattern.Pattern] {
	return func(yield func(pattern.Pattern) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			for _, pat := range p.Patterns {
				if !yield(pat) {
					return
				}
			}
		}
	}
}

// ShenSha streams every 神煞 across the 9 palaces in palace order.
func (c *Chart) ShenSha() iter.Seq[shensha.ShenSha] {
	return func(yield func(shensha.ShenSha) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			for _, ss := range p.ShenSha {
				if !yield(ss) {
					return
				}
			}
		}
	}
}

// EarthStems streams the 地盘 stems with palace numbers.
func (c *Chart) EarthStems() iter.Seq2[uint8, almanac.Stem] {
	return func(yield func(uint8, almanac.Stem) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			if !yield(p.Number, p.EarthStem) {
				return
			}
		}
	}
}

// HeavenStems streams the 天盘 stems with palace numbers.
func (c *Chart) HeavenStems() iter.Seq2[uint8, almanac.Stem] {
	return func(yield func(uint8, almanac.Stem) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			if !yield(p.Number, p.HeavenStem) {
				return
			}
		}
	}
}

// HiddenStems streams the 暗干 stems with palace numbers.
func (c *Chart) HiddenStems() iter.Seq2[uint8, almanac.Stem] {
	return func(yield func(uint8, almanac.Stem) bool) {
		for i := range c.palaces {
			p := &c.palaces[i]
			if !yield(p.Number, p.HiddenStem) {
				return
			}
		}
	}
}
