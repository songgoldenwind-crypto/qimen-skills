package almanac

// Direction labels the 9-grid (洛书) position by compass orientation.
//
//	0=北 (palace 1) 1=东北 (8) 2=东 (3) 3=东南 (4)
//	4=中宫 (5) 5=西北 (6) 6=西 (7) 7=西南 (2) 8=南 (9)
//
// Index 4 is the central palace; the 8 cardinal/intercardinal indices
// correspond to the eight surrounding 八卦 positions.
type Direction uint8

var directionNames = [9]string{"北", "东北", "东", "东南", "中", "西北", "西", "西南", "南"}

// palaceDirection maps palace 1..9 → Direction. Index 0 reserved.
//
//	1坎→北 2坤→西南 3震→东 4巽→东南 5中→中
//	6乾→西北 7兑→西 8艮→东北 9离→南
var palaceDirection = [10]Direction{0, 0, 7, 2, 3, 4, 5, 6, 1, 8}

// DirectionOfPalace maps 1..9 palace numbers to the Direction value.
// Precondition: palace ∈ [1, 9].
func DirectionOfPalace(palace uint8) Direction { return palaceDirection[palace] }

// Name returns the Chinese label.
func (d Direction) Name() string { return directionNames[d] }

// String implements fmt.Stringer.
func (d Direction) String() string { return d.Name() }
