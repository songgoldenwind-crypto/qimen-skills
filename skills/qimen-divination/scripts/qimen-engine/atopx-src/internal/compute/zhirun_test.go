package compute

import (
	"testing"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/enum"
)

// termSeq linearizes a term into a continuous ordinal for comparisons.
func termSeq(t almanac.Term) int { return t.Year()*24 + t.Index() }

// TestZhiRunInvariants scans ~30 years day by day and checks the
// structural properties any correct 置闰 schedule must satisfy:
//
//  1. the working term never moves backward and advances by at most
//     one term per day;
//  2. it never drifts more than one term from the astronomical term
//     (超神 leads by at most one, 接气/闰 trails by at most one);
//  3. every term occupies exactly 15 days — except an intercalated
//     芒种 or 大雪, which occupies 30;
//  4. intercalations happen roughly every 2.9 years (the 符头 grid of
//     15 days loses ~5.2 days per year against the tropical year);
//  5. a 符头 day (day number ≡ 0 mod 15) is always an upper-yuan day.
func TestZhiRunInvariants(t *testing.T) {
	st, err := almanac.SolarTimeOf(2010, 1, 1, 12, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	const days = 30 * 366

	prevSeq := 0
	termDays := map[int]int{}
	minSeq, maxSeq := 1<<62, 0
	for i := 0; i < days; i++ {
		dn := almanac.DayNumber(st)
		real := almanac.TermOfSolarTime(st)
		ju := ZhiRunTerm(dn, real)
		seq := termSeq(ju)

		if i > 0 && (seq < prevSeq || seq > prevSeq+1) {
			t.Fatalf("%s: working term jumped %d → %d", st, prevSeq, seq)
		}
		if d := termSeq(real) - seq; d < -1 || d > 1 {
			t.Fatalf("%s: working term %s drifts %d terms from real %s",
				st, ju.Name(), d, real.Name())
		}
		if floorMod(dn, 15) == 0 && Yuan(almanac.CycleOf(dn)) != enum.YuanUpper {
			t.Fatalf("%s: leader day is not upper yuan", st)
		}

		termDays[seq]++
		if seq < minSeq {
			minSeq = seq
		}
		if seq > maxSeq {
			maxSeq = seq
		}
		prevSeq = seq
		st = st.AddSeconds(86400)
	}

	leaps := 0
	for seq, n := range termDays {
		if seq == minSeq || seq == maxSeq {
			continue // clipped by the scan window
		}
		switch n {
		case 15:
		case 30:
			idx := ((seq % 24) + 24) % 24
			if idx != 11 && idx != 23 {
				t.Errorf("term seq %d (index %d) intercalated; only 芒种/大雪 may repeat", seq, idx)
			}
			leaps++
		default:
			t.Errorf("term seq %d occupies %d days; want 15 or 30", seq, n)
		}
	}
	if leaps < 8 || leaps > 13 {
		t.Errorf("got %d intercalations in ~30 years, want ≈10 (one per ≈2.9 years)", leaps)
	}
}

// TestZhiRunSolsticeAdoption checks the anchor rule on every solstice in
// the scan range: with a lead interval ≤ 7 the solstice's working term
// starts on the leader day itself (超神/正授); at ≥ 8 (nine days counted
// inclusively) the preceding 芒种/大雪 is still working on the solstice
// day (接气 after 置闰) and the next leader opens the solstice term.
func TestZhiRunSolsticeAdoption(t *testing.T) {
	for year := 2010; year <= 2040; year++ {
		for _, idx := range []int{0, 12} {
			z := almanac.TermOf(year, idx)
			zd := z.DayNumber()
			lead := floorMod(zd, 15)

			leader := zd - lead
			if lead > zhiRunThreshold {
				leader += 15
			}
			at := func(dn int) almanac.Term {
				return ZhiRunTerm(dn, almanac.TermOfSolarTime(noonOf(t, dn)))
			}
			if got := at(leader); termSeq(got) != termSeq(z) {
				t.Errorf("%d/%s: adopting leader day works %s, want %s",
					year, z.Name(), got.Name(), z.Name())
			}
			if lead > zhiRunThreshold {
				if got := at(zd); termSeq(got) != termSeq(z)-1 {
					t.Errorf("%d/%s 接气: solstice day works %s, want previous term",
						year, z.Name(), got.Name())
				}
			}
		}
	}
}

// noonOf converts a day number back to its noon SolarTime via the
// almanac day-pillar anchor (2000-01-07 = day 0).
func noonOf(t *testing.T, dn int) almanac.SolarTime {
	t.Helper()
	base, err := almanac.SolarTimeOf(2000, 1, 7, 12, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	return base.AddSeconds(dn * 86400)
}
