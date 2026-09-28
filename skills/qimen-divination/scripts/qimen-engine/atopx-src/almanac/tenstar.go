package almanac

// TenStar (十神) classifies the relationship between two heavenly stems
// from the perspective of the day master (日主).
//
//	0=比肩 1=劫财 2=食神 3=伤官 4=偏财 5=正财 6=七杀 7=正官 8=偏印 9=正印
type TenStar uint8

var tenStarNames = [10]string{"比肩", "劫财", "食神", "伤官", "偏财", "正财", "七杀", "正官", "偏印", "正印"}

// Index returns the 0..9 ordinal.
func (t TenStar) Index() int { return int(t) }

// Name returns the Chinese name.
func (t TenStar) Name() string { return tenStarNames[t] }

// String implements fmt.Stringer.
func (t TenStar) String() string { return t.Name() }
