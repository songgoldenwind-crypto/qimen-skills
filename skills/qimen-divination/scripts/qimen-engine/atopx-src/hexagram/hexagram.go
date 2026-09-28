package hexagram

import (
	"fmt"

	"github.com/atopx/qimen/auspice"
)

// Hexagram 六十四卦, formed from (upper, lower) trigram pairs.
type Hexagram struct {
	upper Trigram
	lower Trigram
	index uint8
}

// Of builds the hexagram from an (upper, lower) trigram pair.
func Of(upper, lower Trigram) Hexagram {
	return Hexagram{upper: upper, lower: lower, index: indexTable[upper][lower]}
}

// Upper returns the upper trigram.
func (h Hexagram) Upper() Trigram { return h.upper }

// Lower returns the lower trigram.
func (h Hexagram) Lower() Trigram { return h.lower }

// Index returns the 周易序卦传 index (0..63).
func (h Hexagram) Index() uint8 { return h.index }

// Symbol returns the Unicode hexagram character (䷀..䷿).
func (h Hexagram) Symbol() string { return dataTable[h.index].symbol }

// Name returns the Chinese name.
func (h Hexagram) Name() string { return dataTable[h.index].name }

// Summary returns the one-line classical description.
func (h Hexagram) Summary() string { return dataTable[h.index].summary }

// Auspice returns the literature-defined auspice level.
func (h Hexagram) Auspice() auspice.Auspice { return auspiceTable[h.index] }

// String returns symbol + name.
func (h Hexagram) String() string { return fmt.Sprintf("%s %s", h.Symbol(), h.Name()) }
