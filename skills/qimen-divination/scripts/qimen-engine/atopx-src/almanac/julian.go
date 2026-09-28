package almanac

import "math"

// julianDayFromYmdhms returns the Julian Day for a Gregorian/Julian
// calendar instant. Switches to the Gregorian system at 1582-10-15.
//
// Implementation follows Meeus' formula (Astronomical Algorithms, ch. 7).
func julianDayFromYmdhms(year, month, day, hour, minute, second int) float64 {
	d := float64(day) + ((float64(second)/60+float64(minute))/60+float64(hour))/24
	n := 0
	g := year*372+month*31+int(d) >= 588829
	if month <= 2 {
		month += 12
		year--
	}
	if g {
		n = int(float64(year) * 0.01)
		n = 2 - n + int(float64(n)*0.25)
	}
	return float64(int(365.25*float64(year+4716))) +
		float64(int(30.6001*float64(month+1))) + d + float64(n) - 1524.5
}

// solarFromJulianDay converts a Julian Day to a SolarTime (year, month, day, hour, minute, second).
//
// Inverse of julianDayFromYmdhms. Handles the 1582 calendar switch.
func solarFromJulianDay(jd float64) SolarTime {
	d := int(jd + 0.5)
	f := jd + 0.5 - float64(d)

	if d >= 2299161 {
		c := int((float64(d) - 1867216.25) / 36524.25)
		d += 1 + c - int(float64(c)*0.25)
	}
	d += 1524
	y := int((float64(d) - 122.1) / 365.25)
	d -= int(365.25 * float64(y))
	m := int(float64(d) / 30.601)
	d -= int(30.601 * float64(m))
	if m > 13 {
		m -= 12
	} else {
		y -= 1
	}
	m -= 1
	y -= 4715

	f *= 24
	hour := int(f)
	f -= float64(hour)
	f *= 60
	minute := int(f)
	f -= float64(minute)
	f *= 60
	second := int(math.Round(f))
	if second >= 60 {
		// rollover (≥1 day)
		t, _ := newSolarTime(y, m, d, hour, minute, second-60)
		return t.AddSeconds(60)
	}
	t, _ := newSolarTime(y, m, d, hour, minute, second)
	return t
}
