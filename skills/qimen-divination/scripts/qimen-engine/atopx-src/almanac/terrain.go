package almanac

// Terrain (长生十二宫) marks one of the 12 life-cycle phases of a stem.
//
//	0=长生 1=沐浴 2=冠带 3=临官 4=帝旺 5=衰 6=病 7=死 8=墓 9=绝 10=胎 11=养
type Terrain uint8

var terrainNames = [12]string{"长生", "沐浴", "冠带", "临官", "帝旺", "衰", "病", "死", "墓", "绝", "胎", "养"}

// Index returns the 0..11 ordinal.
func (t Terrain) Index() int { return int(t) }

// Name returns the Chinese name.
func (t Terrain) Name() string { return terrainNames[t] }

// String implements fmt.Stringer.
func (t Terrain) String() string { return t.Name() }
