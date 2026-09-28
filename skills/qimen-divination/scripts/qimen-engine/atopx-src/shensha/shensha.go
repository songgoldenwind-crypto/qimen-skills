package shensha

import (
	"fmt"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/auspice"
)

// Target identifies the 神煞 anchor entity (stem or branch).
type Target struct {
	Stem   almanac.Stem
	Branch almanac.Branch
	IsStem bool // true → Stem populated; false → Branch populated
}

// String returns the Chinese name of the underlying stem/branch.
func (t Target) String() string {
	if t.IsStem {
		return t.Stem.Name()
	}
	return t.Branch.Name()
}

// ShenSha is a single 神煞 instance landing in one palace.
type ShenSha struct {
	Kind   Kind
	Target Target
	Palace uint8
}

// Name delegates to Kind.
func (s ShenSha) Name() string { return s.Kind.Name() }

// Summary delegates to Kind.
func (s ShenSha) Summary() string { return s.Kind.Summary() }

// Auspice delegates to Kind.
func (s ShenSha) Auspice() auspice.Auspice { return s.Kind.Auspice() }

// String formats as "<name>(<target>→<palace>宫)".
func (s ShenSha) String() string {
	return fmt.Sprintf("%s(%s→%d宫)", s.Kind.Name(), s.Target.String(), s.Palace)
}
