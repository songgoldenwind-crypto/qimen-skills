package almanac

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidTime indicates an out-of-range or otherwise invalid SolarTime input.
var ErrInvalidTime = errors.New("almanac: invalid solar time")

// locCST is the canonical China Standard Time zone (UTC+8) used by all
// entry points that must interpret an absolute instant as a wall clock.
var locCST = time.FixedZone("CST", 8*3600)

// SolarTime is an immutable solar (Gregorian) instant at second precision.
// All fields are exported for ergonomic field access; the type is treated as
// a value (copied freely, never mutated in place).
type SolarTime struct {
	Year   int16
	Month  uint8
	Day    uint8
	Hour   uint8
	Minute uint8
	Second uint8
}

// SolarTimeOf builds a SolarTime, validating the calendar components.
func SolarTimeOf(year, month, day, hour, minute, second int) (SolarTime, error) {
	return newSolarTime(year, month, day, hour, minute, second)
}

// SolarTimeFromTime converts a Go time.Time (any location) to a SolarTime
// in that time's wall clock.
func SolarTimeFromTime(t time.Time) (SolarTime, error) {
	y, mo, d := t.Date()
	h, mi, s := t.Clock()
	return newSolarTime(y, int(mo), d, h, mi, s)
}

// SolarTimeFromUnix interprets a Unix-seconds timestamp as a UTC+8
// (China Standard Time) wall clock and returns the corresponding
// SolarTime. For other zones convert explicitly:
// SolarTimeFromTime(time.Unix(unix, 0).In(loc)).
func SolarTimeFromUnix(unix int64) (SolarTime, error) {
	return SolarTimeFromTime(time.Unix(unix, 0).In(locCST))
}

// Now returns SolarTime equivalent to time.Now() in UTC+8.
func Now() SolarTime {
	t, _ := SolarTimeFromTime(time.Now().In(locCST))
	return t
}

// newSolarTime constructs a SolarTime after validating each component.
func newSolarTime(year, month, day, hour, minute, second int) (SolarTime, error) {
	if year < -9999 || year > 9999 {
		return SolarTime{}, fmt.Errorf("%w: year %d out of [-9999..9999]", ErrInvalidTime, year)
	}
	if month < 1 || month > 12 {
		return SolarTime{}, fmt.Errorf("%w: month %d not in 1..12", ErrInvalidTime, month)
	}
	maxDay := daysInMonth(year, month)
	if day < 1 || day > maxDay {
		return SolarTime{}, fmt.Errorf("%w: day %d not in 1..%d for %04d-%02d", ErrInvalidTime, day, maxDay, year, month)
	}
	if hour < 0 || hour > 23 {
		return SolarTime{}, fmt.Errorf("%w: hour %d not in 0..23", ErrInvalidTime, hour)
	}
	if minute < 0 || minute > 59 {
		return SolarTime{}, fmt.Errorf("%w: minute %d not in 0..59", ErrInvalidTime, minute)
	}
	if second < 0 || second > 59 {
		return SolarTime{}, fmt.Errorf("%w: second %d not in 0..59", ErrInvalidTime, second)
	}
	return SolarTime{
		Year:   int16(year),
		Month:  uint8(month),
		Day:    uint8(day),
		Hour:   uint8(hour),
		Minute: uint8(minute),
		Second: uint8(second),
	}, nil
}

// daysInMonth applies Gregorian leap-year rules.
func daysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	}
	return 0
}

// JulianDay returns the Julian Day number corresponding to this instant.
func (s SolarTime) JulianDay() float64 {
	return julianDayFromYmdhms(int(s.Year), int(s.Month), int(s.Day), int(s.Hour), int(s.Minute), int(s.Second))
}

// ToTime converts to a Go time.Time using the provided location.
// If loc is nil, time.UTC is used.
func (s SolarTime) ToTime(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.UTC
	}
	return time.Date(int(s.Year), time.Month(s.Month), int(s.Day),
		int(s.Hour), int(s.Minute), int(s.Second), 0, loc)
}

// Before reports whether s < other in solar time.
func (s SolarTime) Before(other SolarTime) bool {
	switch {
	case s.Year != other.Year:
		return s.Year < other.Year
	case s.Month != other.Month:
		return s.Month < other.Month
	case s.Day != other.Day:
		return s.Day < other.Day
	case s.Hour != other.Hour:
		return s.Hour < other.Hour
	case s.Minute != other.Minute:
		return s.Minute < other.Minute
	}
	return s.Second < other.Second
}

// After reports whether s > other in solar time.
func (s SolarTime) After(other SolarTime) bool { return other.Before(s) }

// Equal reports whether two SolarTime values represent the same instant.
func (s SolarTime) Equal(other SolarTime) bool { return s == other }

// AddSeconds returns a new SolarTime shifted by n seconds.
//
// Pure integer arithmetic over the y/m/d/h/m/s components — avoids the
// allocation overhead of round-tripping through time.Time. Handles
// arbitrary positive or negative n.
func (s SolarTime) AddSeconds(n int) SolarTime {
	if n == 0 {
		return s
	}
	sec := int(s.Second) + n
	minute := int(s.Minute) + sec/60
	sec %= 60
	if sec < 0 {
		sec += 60
		minute--
	}
	hour := int(s.Hour) + minute/60
	minute %= 60
	if minute < 0 {
		minute += 60
		hour--
	}
	dayShift := hour / 24
	hour %= 24
	if hour < 0 {
		hour += 24
		dayShift--
	}
	year, month, day := int(s.Year), int(s.Month), int(s.Day)+dayShift
	// Normalize day → month/year via daysInMonth.
	for day < 1 {
		month--
		if month < 1 {
			month = 12
			year--
		}
		day += daysInMonth(year, month)
	}
	for {
		dim := daysInMonth(year, month)
		if day <= dim {
			break
		}
		day -= dim
		month++
		if month > 12 {
			month = 1
			year++
		}
	}
	return SolarTime{
		Year: int16(year), Month: uint8(month), Day: uint8(day),
		Hour: uint8(hour), Minute: uint8(minute), Second: uint8(sec),
	}
}

// String returns "YYYY-MM-DD HH:MM:SS".
func (s SolarTime) String() string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		s.Year, s.Month, s.Day, s.Hour, s.Minute, s.Second)
}
