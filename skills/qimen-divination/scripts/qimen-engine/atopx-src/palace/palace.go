// Package palace defines the per-cell Palace value that aggregates
// every derived attribute (stems, star, door, god, ten-star, terrain,
// hexagram, patterns, shensha) for one of the 9 grid positions.
package palace

import (
	"fmt"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/auspice"
	"github.com/atopx/qimen/element"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/hexagram"
	"github.com/atopx/qimen/pattern"
	"github.com/atopx/qimen/shensha"
	"github.com/atopx/qimen/terrain"
)

// Palace holds the data for one of the 9 grid cells of a qimen chart.
//
// Star / Door / God presence depends on the chart style: in 转盘 the
// center palace holds none of them; in 飞盘 the center participates in
// the flight, nine stars cover nine palaces and the eight doors / gods
// leave exactly one palace (not necessarily the center) empty. The
// *Set flags are the source of truth for presence — the value fields
// hold zero values when unset.
//
// The derived attributes TenStar / Terrain are populated for every
// non-center palace (the center has no branches to derive from), and
// Hexagram additionally requires a door (HexagramSet).
type Palace struct {
	Number     uint8
	Name       string
	Direction  almanac.Direction
	Branches   []almanac.Branch
	EarthStem  almanac.Stem // 地盘干 (三奇六仪)
	HeavenStem almanac.Stem
	HiddenStem almanac.Stem

	Star    enum.Star
	StarSet bool
	Door    enum.Door
	DoorSet bool
	God     enum.God
	GodSet  bool

	// Populated iff Number != 5.
	TenStar almanac.TenStar
	Terrain terrain.Terrain
	// Populated iff Number != 5 and DoorSet.
	Hexagram    hexagram.Hexagram
	HexagramSet bool

	Patterns []pattern.Pattern
	ShenSha  []shensha.ShenSha
}

// IsCenter reports whether this is the center palace (Number == 5).
func (p *Palace) IsCenter() bool { return p.Number == 5 }

// Element returns the 五行 of this palace.
func (p *Palace) Element() element.Element { return element.FromPalace(p.Number) }

// DoorPalaceRelation computes how the 八门 of this palace relates to
// its 宫位 五行. Returns the zero Relation (empty Description) when the
// palace holds no door.
func (p *Palace) DoorPalaceRelation() Relation {
	if !p.DoorSet {
		return Relation{}
	}
	return relationFromSubject(p.Door.Name(), element.OfDoor(p.Door), p.Element())
}

// StarPalaceRelation computes how the 九星 of this palace relates to
// its 宫位 五行. Returns the zero Relation (empty Description) when the
// palace holds no star.
func (p *Palace) StarPalaceRelation() Relation {
	if !p.StarSet {
		return Relation{}
	}
	return relationFromSubject(p.Star.Name(), element.OfStar(p.Star), p.Element())
}

// Relation is the result of a subject (door/star) vs. palace 五行 comparison.
type Relation struct {
	Description string
	Element     element.Relation
	Auspice     auspice.Auspice
}

// String formats as "<desc>[<auspice>]".
func (r Relation) String() string {
	return fmt.Sprintf("%s[%s]", r.Description, r.Auspice.Name())
}

// relationDescPrefix / Suffix wrap the subject name with the relation phrase.
// Indexed by element.Relation.
var (
	relationDescPrefix = [5]string{
		element.Same:       "",
		element.Generates:  "",
		element.Generated:  "宫生",
		element.Restrains:  "",
		element.Restrained: "宫克",
	}
	relationDescSuffix = [5]string{
		element.Same:       "与宫比和",
		element.Generates:  "生宫",
		element.Generated:  "",
		element.Restrains:  "克宫",
		element.Restrained: "",
	}
)

func relationFromSubject(subjectName string, subjEl, palaceEl element.Element) Relation {
	rel := subjEl.RelationTo(palaceEl)
	return Relation{
		Description: relationDescPrefix[rel] + subjectName + relationDescSuffix[rel],
		Element:     rel,
		Auspice:     rel.AuspiceAsSelf(),
	}
}
