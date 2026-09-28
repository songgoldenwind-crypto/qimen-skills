package compute

import "github.com/atopx/qimen/almanac"

// zhiRunThreshold is the maximum tolerated 超神 lead, as the day-grid
// interval between the upper-yuan 符头 and the solstice day. Classical
// texts count inclusively ("符头超前九日则置闰", the leader day itself
// being day one), so the traditional bound of nine days corresponds to
// an interval of 8: a solstice 0..7 days after its leader is adopted by
// it (超神/正授); at 8..14 days the preceding 芒种 / 大雪 repeats one
// extra 三元 (置闰) and the solstice runs 接气.
const zhiRunThreshold = 7

// ZhiRunTerm resolves the 用局节气 (the term that supplies the 局) for
// a day under the 置闰 rule. dayNum is the almanac.DayNumber of the
// instant; realTerm is the astronomical term in effect (s.Term()).
//
// The rule keeps the 15-day 符头 grid (upper-yuan leader days 甲子 /
// 己卯 / 甲午 / 己酉) aligned with the solstices:
//
//   - Each solstice is "adopted" by a leader: the one at most
//     zhiRunThreshold days before it (超神 / 正授), or — when the lead
//     interval reaches 8..14 days — the next leader after an extra
//     repeated 三元 of the preceding 芒种 / 大雪 (置闰), leaving the
//     solstice 接气.
//   - From the adopting leader, every leader advances the working term
//     by one, so each term occupies exactly one 15-day 三元 (30 days
//     for the intercalated 芒种 / 大雪).
//
// Both half-year segments tile exactly: solstices are ~182.62 days
// apart, i.e. 12 leaders (180 days, no intercalation) or 13 (195 days,
// one intercalation), so the anchor chosen below is locally decidable
// with no historical state.
func ZhiRunTerm(dayNum int, realTerm almanac.Term) almanac.Term {
	f := dayNum - floorMod(dayNum, 15) // upper-yuan leader (符头) of this day

	// Candidate anchor solstices, nearest first: the upcoming one (the
	// leader may already adopt it — 超神), the one opening the current
	// half-year, and one more back (covers 接气 days that still belong
	// to the previous segment's tail or its intercalated 三元).
	var candidates [3]almanac.Term
	if realTerm.Index() < 12 {
		candidates[0] = almanac.TermOf(realTerm.Year(), 12)
		candidates[1] = almanac.TermOf(realTerm.Year(), 0)
	} else {
		candidates[0] = almanac.TermOf(realTerm.Year()+1, 0)
		candidates[1] = almanac.TermOf(realTerm.Year(), 12)
	}
	candidates[2] = candidates[1].Next(-12)

	for _, z := range candidates {
		g := adoptingLeader(z)
		if g > f {
			continue
		}
		// Leaders count terms forward from the anchor; the 13th leader
		// (k == 12) exists only in an intercalated segment and repeats
		// the segment's closing 芒种 / 大雪.
		k := (f - g) / 15
		if k > 11 {
			k = 11
		}
		return z.Next(k)
	}
	return realTerm // unreachable: candidates[2] always precedes f
}

// adoptingLeader returns the day number of the 符头 that adopts the
// given solstice: the leader 0..7 days before it (超神 / 正授), or the
// following leader when the lead reaches 8..14 days (the preceding
// 芒种 / 大雪 was intercalated and the solstice runs 接气).
func adoptingLeader(z almanac.Term) int {
	zd := z.DayNumber()
	lead := floorMod(zd, 15)
	g := zd - lead
	if lead > zhiRunThreshold {
		return g + 15
	}
	return g
}
