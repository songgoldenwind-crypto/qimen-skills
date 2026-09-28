package almanac

// Ten represents one of the six 旬 (decade groups) in the sexagenary cycle.
// Each 旬 starts with 甲 paired with a different branch:
//
//	0=甲子旬, 1=甲戌旬, 2=甲申旬, 3=甲午旬, 4=甲辰旬, 5=甲寅旬
type Ten uint8

var tenNames = [6]string{"甲子", "甲戌", "甲申", "甲午", "甲辰", "甲寅"}

// Index returns the 0..5 ordinal.
func (t Ten) Index() int { return int(t) }

// FirstBranch returns the branch paired with 甲 at the start of this 旬:
// 甲子→子, 甲戌→戌, 甲申→申, 甲午→午, 甲辰→辰, 甲寅→寅 (each 旬
// starts two branches earlier than the previous one).
func (t Ten) FirstBranch() Branch { return BranchOf(-2 * int(t)) }

// Name returns the 旬 starting pair (e.g. "甲子").
func (t Ten) Name() string { return tenNames[t] }

// String implements fmt.Stringer.
func (t Ten) String() string { return t.Name() }
