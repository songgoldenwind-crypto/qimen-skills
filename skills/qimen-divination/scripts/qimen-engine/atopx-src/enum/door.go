package enum

// Door 八门 (休/生/伤/杜/景/死/惊/开).
type Door uint8

const (
	DoorRest  Door = iota // 休门
	DoorLife              // 生门
	DoorHurt              // 伤门
	DoorBlock             // 杜门
	DoorView              // 景门
	DoorDeath             // 死门
	DoorFear              // 惊门
	DoorOpen              // 开门
)

var doorNames = [8]string{"休", "生", "伤", "杜", "景", "死", "惊", "开"}

// Name returns the Chinese label.
func (d Door) Name() string { return doorNames[d] }

// String implements fmt.Stringer.
func (d Door) String() string { return d.Name() }

// doorHomePalace maps Door index → home palace.
//
//	休→1, 生→8, 伤→3, 杜→4, 景→9, 死→2, 惊→7, 开→6.
var doorHomePalace = [8]uint8{1, 8, 3, 4, 9, 2, 7, 6}

// HomePalace returns the canonical 本位宫 (1..9).
func (d Door) HomePalace() uint8 { return doorHomePalace[d] }

// palaceDoor maps palace 1..9 → Door. Index 0 and 5 (center) hold the
// zero value DoorRest — callers MUST gate by `palace != 5` first, since
// center has no door.
var palaceDoor = [10]Door{
	0,
	DoorRest, DoorDeath, DoorHurt, DoorBlock, 0,
	DoorOpen, DoorFear, DoorLife, DoorView,
}

// DoorOfPalace returns the canonical 本位门 for a palace.
// Precondition: palace ∈ [1, 9] AND palace != 5 (center has no door).
func DoorOfPalace(palace uint8) Door { return palaceDoor[palace] }
