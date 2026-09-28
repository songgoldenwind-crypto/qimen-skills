package enum

// God 九神 (值符/腾蛇/太阴/六合/白虎/玄武/九地/九天).
type God uint8

const (
	GodZhiFu    God = iota // 值符
	GodTengShe             // 腾蛇
	GodTaiYin              // 太阴
	GodLiuHe               // 六合
	GodBaiHu               // 白虎
	GodXuanWu              // 玄武
	GodJiuDi               // 九地
	GodJiuTian             // 九天
	GodTaiChang            // 太常 (ninth god, fly-style charts only)
)

var godNames = [9]string{"值符", "腾蛇", "太阴", "六合", "白虎", "玄武", "九地", "九天", "太常"}

// Name returns the Chinese label.
func (g God) Name() string { return godNames[g] }

// String implements fmt.Stringer.
func (g God) String() string { return g.Name() }
