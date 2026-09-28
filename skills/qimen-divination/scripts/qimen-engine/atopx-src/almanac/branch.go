package almanac

// Branch represents one of the 12 earthly branches (地支).
// Values 0..11 correspond to 子, 丑, 寅, 卯, 辰, 巳, 午, 未, 申, 酉, 戌, 亥.
type Branch uint8

// Canonical branch indices.
const (
	Zi Branch = iota
	Chou
	Yin_ // 寅 (named Yin_ to avoid clash with YinYang.Yin)
	Mao
	Chen
	Si
	Wu_ // 午 (clash with stem 戊 Wu)
	Wei
	Shen
	You
	Xu
	Hai
)

var branchNames = [12]string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

var branchAnimals = [12]string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}

// BranchOf wraps an integer into Branch, normalizing into 0..11.
func BranchOf(i int) Branch {
	i = ((i % 12) + 12) % 12
	return Branch(i)
}

// Index returns the 0..11 ordinal.
func (b Branch) Index() int { return int(b) }

// Name returns the Chinese character.
func (b Branch) Name() string { return branchNames[b] }

// String implements fmt.Stringer.
func (b Branch) String() string { return b.Name() }

// Next returns the branch advanced by n steps (negative goes backward).
func (b Branch) Next(n int) Branch { return BranchOf(int(b) + n) }

// Animal returns the Chinese zodiac animal name.
func (b Branch) Animal() string { return branchAnimals[b] }
