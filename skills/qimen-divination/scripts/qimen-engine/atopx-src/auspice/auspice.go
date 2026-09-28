// Package auspice defines the 5-level traditional 吉凶 classification
// shared by all qimen domain entities (patterns, shensha, hexagrams,
// long-sheng-12 positions).
//
// Values are fixed by classical literature, not computed dynamically.
package auspice

// Auspice 吉凶等级 (5 级).
//
// This is the **literature-defined** auspice of an entity (e.g. 青龙返首
// is 大吉 by tradition). It is not a numeric score; callers can combine
// multiple Auspice values themselves if they need a composite metric.
type Auspice uint8

const (
	// GreatAuspicious 大吉 — extremely favorable.
	GreatAuspicious Auspice = iota
	// Auspicious 吉 — generally favorable.
	Auspicious
	// Neutral 中和 — stable, no strong direction.
	Neutral
	// Inauspicious 凶 — generally unfavorable.
	Inauspicious
	// GreatInauspicious 大凶 — extremely unfavorable.
	GreatInauspicious
)

var names = [5]string{"大吉", "吉", "中和", "凶", "大凶"}

// Name returns the Chinese label.
func (a Auspice) Name() string { return names[a] }

// String implements fmt.Stringer.
func (a Auspice) String() string { return a.Name() }

// IsAuspicious reports 大吉 or 吉.
func (a Auspice) IsAuspicious() bool { return a == GreatAuspicious || a == Auspicious }

// IsInauspicious reports 大凶 or 凶.
func (a Auspice) IsInauspicious() bool { return a == GreatInauspicious || a == Inauspicious }

// IsExtreme reports 大吉 or 大凶.
func (a Auspice) IsExtreme() bool { return a == GreatAuspicious || a == GreatInauspicious }

// Auspicable is the shared interface implemented by every qimen domain
// entity: name + one-line summary + literature-defined auspice level.
type Auspicable interface {
	Name() string
	Summary() string
	Auspice() Auspice
}
