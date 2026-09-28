package almanac

// Cycle is one of the 60 sexagenary cycles (六十甲子), 0..59.
//
//	0=甲子, 1=乙丑, 2=丙寅, ..., 59=癸亥
type Cycle uint8

var cycleNames = [60]string{
	"甲子", "乙丑", "丙寅", "丁卯", "戊辰", "己巳", "庚午", "辛未", "壬申", "癸酉",
	"甲戌", "乙亥", "丙子", "丁丑", "戊寅", "己卯", "庚辰", "辛巳", "壬午", "癸未",
	"甲申", "乙酉", "丙戌", "丁亥", "戊子", "己丑", "庚寅", "辛卯", "壬辰", "癸巳",
	"甲午", "乙未", "丙申", "丁酉", "戊戌", "己亥", "庚子", "辛丑", "壬寅", "癸卯",
	"甲辰", "乙巳", "丙午", "丁未", "戊申", "己酉", "庚戌", "辛亥", "壬子", "癸丑",
	"甲寅", "乙卯", "丙辰", "丁巳", "戊午", "己未", "庚申", "辛酉", "壬戌", "癸亥",
}

// CycleOf wraps an integer into Cycle, normalizing into 0..59.
func CycleOf(i int) Cycle {
	i = ((i % 60) + 60) % 60
	return Cycle(i)
}

// Index returns the 0..59 ordinal.
func (c Cycle) Index() int { return int(c) }

// Name returns the two-character Chinese name.
func (c Cycle) Name() string { return cycleNames[c] }

// String implements fmt.Stringer.
func (c Cycle) String() string { return c.Name() }

// Stem returns the stem component.
func (c Cycle) Stem() Stem { return Stem(int(c) % 10) }

// Branch returns the branch component.
func (c Cycle) Branch() Branch { return Branch(int(c) % 12) }

// Ten returns the 旬 (decade group): 0=甲子, 1=甲戌, 2=甲申, 3=甲午, 4=甲辰, 5=甲寅.
func (c Cycle) Ten() Ten {
	return Ten((int(c.Stem()) - int(c.Branch()) + 12) / 2 % 6)
}

// Next returns the cycle advanced by n steps (negative goes backward).
func (c Cycle) Next(n int) Cycle { return CycleOf(int(c) + n) }

// EmptyBranches returns the two 空亡 branches for the 旬 containing c.
// In every decade the stems pair with 10 branches, leaving 2 "extra"
// branches: those are the 空亡 pair.
func (c Cycle) EmptyBranches() [2]Branch {
	first := BranchOf(10 + int(c.Branch()) - int(c.Stem()))
	return [2]Branch{first, first.Next(1)}
}
