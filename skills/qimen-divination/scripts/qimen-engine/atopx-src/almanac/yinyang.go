package almanac

// YinYang denotes the yin/yang polarity of a stem, year, or other
// entity. The numeric values follow the classical binary reading:
// 阴 = 0 (broken line), 阳 = 1 (solid line).
type YinYang uint8

const (
	// Yin 阴
	Yin YinYang = iota
	// Yang 阳
	Yang
)

var yinYangNames = [2]string{"阴", "阳"}

// Name returns the Chinese character "阴" or "阳".
func (y YinYang) Name() string { return yinYangNames[y] }

// String implements fmt.Stringer.
func (y YinYang) String() string { return y.Name() }
