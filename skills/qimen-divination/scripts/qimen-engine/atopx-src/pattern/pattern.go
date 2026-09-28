package pattern

import (
	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/auspice"
	"github.com/atopx/qimen/enum"
)

// Pattern is a single 格局 instance. Kind is the discriminator that
// determines which side fields are meaningful:
//
//	Kind=FanYin    → OriginalPalace meaningful, Stem/Door/Branch ignored
//	Kind=RuMu      → Stem meaningful
//	Kind=MenPo     → Door meaningful
//	Kind=KongWang  → Branch meaningful
//	Kind=anything else → only Palace + Kind are meaningful
//
// Zero values in unused fields are never read by callers that respect
// Kind; previous (Stem/Door/Branch)Set flags were redundant with Kind
// and have been removed.
type Pattern struct {
	Kind           Kind
	Palace         uint8
	OriginalPalace uint8          // FanYin only
	Stem           almanac.Stem   // RuMu only
	Door           enum.Door      // MenPo only
	Branch         almanac.Branch // KongWang only
}

// Name delegates to Kind.
func (p Pattern) Name() string { return p.Kind.Name() }

// Summary delegates to Kind.
func (p Pattern) Summary() string { return p.Kind.Summary() }

// Auspice delegates to Kind.
func (p Pattern) Auspice() auspice.Auspice { return p.Kind.Auspice() }

// String returns the Chinese name.
func (p Pattern) String() string { return p.Name() }
