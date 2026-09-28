package enum

// JuRule 定局规则 — how the 用局节气 (the term whose 局 table row is
// used) is resolved from the calendar.
//
// Both rules share the same 三元 derivation (day-pillar 符头 grid);
// they differ only in which solar term supplies the 局:
//
//   - 拆补 uses the astronomical term in effect at the instant. When a
//     term changes mid-旬, the running 三元 is "split" across the old
//     term and "patched" onto the new one.
//   - 置闰 keeps 符头 (upper-yuan leader days) aligned with the
//     solstices: a leader shortly before a term start adopts the coming
//     term early (超神), and once the lead reaches nine days counted
//     inclusively an extra 三元 of 芒种 / 大雪 is repeated (置闰),
//     after which terms run behind their leaders (接气).
type JuRule uint8

const (
	// JuRuleZhiRun 置闰法 (default).
	JuRuleZhiRun JuRule = iota
	// JuRuleChaiBu 拆补法.
	JuRuleChaiBu
)

var juRuleNames = [2]string{"置闰", "拆补"}

// Name returns the Chinese label.
func (r JuRule) Name() string { return juRuleNames[r] }

// String implements fmt.Stringer.
func (r JuRule) String() string { return r.Name() }
