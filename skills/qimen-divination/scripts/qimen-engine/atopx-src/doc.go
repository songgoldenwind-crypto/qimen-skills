// Package qimen builds 奇门遁甲 (Qimen Dunjia) charts from a solar
// instant. The package is fully self-contained: no external dependencies.
//
// # Quick start
//
//	// Build a chart from a solar instant
//	chart := qimen.New()                        // current time, default options
//	chart, _ := qimen.FromTimestamp(1735632000) // Unix seconds (UTC+8 wall clock)
//	t, _ := time.Parse("2006-01-02 15:04", "2026-01-14 18:45")
//	chart, _ := qimen.FromTime(t)               // standard library time.Time
//	st, _ := almanac.SolarTimeOf(2026, 1, 14, 18, 45, 0)
//	chart := qimen.From(st)                     // first-class SolarTime, total
//
//	// Iterate the 9 palaces
//	for n, p := range chart.Palaces() {
//	    fmt.Printf("%d宫 %s — 地盘:%s 天盘:%s\n",
//	        n, p.Name, p.EarthStem.Name(), p.HeavenStem.Name())
//	}
//
//	// Stream patterns + shensha
//	for p := range chart.Patterns() {
//	    fmt.Printf("格局 %s [%s]\n", p.Name(), p.Auspice().Name())
//	}
//
// # Time zones
//
// Charts are cast in the local civil (wall clock) time of the event.
// [FromTime] uses the wall clock of the supplied time.Time as-is;
// [FromTimestamp] and [New] interpret the instant in UTC+8 (China
// Standard Time). For another zone, convert before calling:
// FromTime(time.Unix(unix, 0).In(loc)).
//
// # Options
//
//   - [WithMethod] selects the 起局法门: 时家 (default), 日家, 月家,
//     年家. The method picks the duty pillar (主柱) and the 局 source —
//     时家 and 日家 share the per-day 节气三元 局; 月/年家 use the
//     year/month calendars and are always 阴遁.
//   - [WithStyle] selects the 盘式: 转盘 (default, LuoShu ring rotation
//     with the center 寄坤) or 飞盘 (palace-number flying with the
//     center as a regular stop, 天禽 on its real palace and the NINE
//     gods — 太常 included — covering every palace).
//   - [WithJuRule] selects the 时家 / 日家 定局规则: 置闰 (default) or
//     拆补.
//
// Every option combination is implemented, so [From] and [New] are
// total — only [FromTime] / [FromTimestamp] can fail, with
// [ErrInvalidTime], on out-of-range calendar input.
//
// # Conventions
//
// The layout follows mainstream 时家转盘 practice; points where schools
// diverge are fixed as follows:
//
//   - 三元 is derived from the day pillar 符头 grid (index mod 15). The
//     时家 局 is keyed to the solstice-aligned 置闰 schedule by default
//     (超神 / 接气, intercalating 芒种 / 大雪 when the leader's advance
//     reaches nine days counted inclusively) — validated against
//     authoritative reference charts (see qimen_golden_test.go); 拆补
//     keys it to the astronomical term instead. [Chart.JuTerm] exposes
//     the term actually used.
//   - The day pillar switches at 23:00 (晚子时 counts as the next day).
//   - 暗干 follows the 值使: the earth stems shift together with the
//     door plate (门下藏干).
//   - Duty records report real palaces (the center included); the 寄坤
//     projection exists only inside rotate-style plate arithmetic.
package qimen
