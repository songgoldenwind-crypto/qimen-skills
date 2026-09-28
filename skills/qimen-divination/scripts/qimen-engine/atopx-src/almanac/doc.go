// Package almanac implements Chinese calendrical computations
// (solar terms, new moons, sexagenary cycles, lunar dates)
// used by qimen divination.
//
// Astronomical algorithms are based on the public-domain
// 寿星万年历 (Shouxing Wannianli):
//
//   - VSOP87 truncated series for ecliptic longitude of the Sun and Moon
//   - IAU 1980 nutation series (truncated to 10 leading periodic terms)
//   - A piecewise ΔT polynomial approximation calibrated for ~1820..2050
//
// The package is fully self-contained — it has no external dependencies
// and is intended to be embedded into other libraries that need accurate
// historical and contemporary Chinese calendar values without pulling
// in a general-purpose astronomical library.
//
// # Stability
//
// All exported types and functions are immutable value types
// (Stem / Branch / Cycle / SolarTime / Term / LunarDay / ...).
// Methods never mutate the receiver and never return shared mutable
// pointers — callers can safely share values across goroutines.
package almanac
