package enum

// Style 奇门盘式 (turn-plate vs fly-plate vs sizhu).
type Style uint8

const (
	// StyleRotate 转盘 (三元盘 / rotate-plate, default, currently the only
	// fully implemented style).
	StyleRotate Style = iota
	// StyleFly 飞盘 (reserved for future implementation).
	StyleFly
)

var styleNames = [2]string{"转盘", "飞盘"}

// Name returns the Chinese label.
func (s Style) Name() string { return styleNames[s] }

// String implements fmt.Stringer.
func (s Style) String() string { return s.Name() }
