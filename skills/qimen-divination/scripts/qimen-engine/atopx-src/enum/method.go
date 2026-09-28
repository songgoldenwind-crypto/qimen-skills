// Package enum holds the qimen-specific domain enumerations
// (Method / Style / Yuan / Star / Door / God).
//
// All values are exported with type-qualifying prefixes (`MethodTime`,
// `StarTianPeng`, ...) so they read naturally as `enum.StarTianPeng`
// without the package name doubling up.
package enum

// Method 奇门起局法门. The method picks the duty pillar (主柱) that the
// 旬首, 值符 / 值使 movement and 空亡 derive from, plus the 局 source:
// 时家 and 日家 share the 节气三元 table (the 局 is a per-day fact);
// 月家 and 年家 use the year/month calendars and are always 阴遁.
type Method uint8

const (
	// MethodTime 时家 (default): duty pillar = 时柱, 局 from 节气三元.
	MethodTime Method = iota
	// MethodDay 日家: duty pillar = 日柱, same per-day 局 as 时家.
	MethodDay
	// MethodMonth 月家: duty pillar = 月柱, 阴遁; the 寅月 局 is 8/5/2
	// by year-branch triad (子午卯酉/辰戌丑未/寅申巳亥), retreating
	// monthly.
	MethodMonth
	// MethodYear 年家: duty pillar = 年柱, 阴遁; the 局 is fixed per
	// 60-year 元 (上1 中4 下7, anchored at the 上元甲子 year 1864).
	MethodYear
)

var methodNames = [4]string{"时家", "日家", "月家", "年家"}

// Name returns the Chinese label.
func (m Method) Name() string { return methodNames[m] }

// String implements fmt.Stringer.
func (m Method) String() string { return m.Name() }
