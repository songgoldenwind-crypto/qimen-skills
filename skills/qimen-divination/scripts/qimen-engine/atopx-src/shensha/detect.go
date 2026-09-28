package shensha

import (
	"iter"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/internal/tables"
)

// DetectInput is the chart view used by Detect.
//
// EarthStems is the 地盘 三奇六仪 layout (palace - 1 indexed) used to find
// stem-anchored 神煞 (天德/月德 when keyed by stem).
type DetectInput struct {
	YearStem    almanac.Stem
	MonthBranch almanac.Branch
	DayStem     almanac.Stem
	DayBranch   almanac.Branch

	EarthStems [9]almanac.Stem
}

// AppendAll appends every 神煞 instance landing in the chart to dst and
// returns the extended slice. This is the allocation-friendly core;
// Detect wraps it as a sequence.
func AppendAll(dst []ShenSha, in *DetectInput) []ShenSha {
	dayBranchIdx := in.DayBranch.Index()
	dayStemIdx := in.DayStem.Index()
	monthBranchIdx := in.MonthBranch.Index()
	yearStemIdx := in.YearStem.Index()

	// Branch-anchored shensha from day branch.
	dst = appendBranch(dst, YiMa, int(yiMaTable[dayBranchIdx]))
	dst = appendBranch(dst, TaoHua, int(taoHuaTable[dayBranchIdx]))
	dst = appendBranch(dst, HuaGai, int(huaGaiTable[dayBranchIdx]))

	// Stem-anchored shensha from day stem.
	pair := tianYiTable[dayStemIdx]
	dst = appendBranch(dst, TianYi, int(pair[0]))
	dst = appendBranch(dst, TianYi, int(pair[1]))
	dst = appendBranch(dst, WenChang, int(wenChangTable[dayStemIdx]))
	dst = appendBranch(dst, LuShen, int(luShenTable[dayStemIdx]))
	if idx := yangRenTable[dayStemIdx]; idx >= 0 {
		dst = appendBranch(dst, YangRen, idx)
	}

	// Month-anchored 天德/月德.
	entry := tianDeTable[monthBranchIdx]
	if entry.kind == 0 {
		dst = appendStem(dst, TianDe, int(entry.value), &in.EarthStems)
	} else {
		dst = appendBranch(dst, TianDe, int(entry.value))
	}
	dst = appendStem(dst, YueDe, int(yueDeTable[monthBranchIdx]), &in.EarthStems)

	// Year-anchored 国印贵人.
	return appendBranch(dst, GuoYin, int(guoYinTable[yearStemIdx]))
}

// Detect yields every 神煞 instance landing in the chart.
func Detect(in DetectInput) iter.Seq[ShenSha] {
	return func(yield func(ShenSha) bool) {
		for _, s := range AppendAll(nil, &in) {
			if !yield(s) {
				return
			}
		}
	}
}

// appendBranch appends a single branch-anchored shensha.
// Precondition: branchIdx ∈ [0, 11] (always true for table-sourced values).
func appendBranch(dst []ShenSha, kind Kind, branchIdx int) []ShenSha {
	return append(dst, ShenSha{
		Kind:   kind,
		Target: Target{Branch: almanac.BranchOf(branchIdx)},
		Palace: tables.BranchToPalace[branchIdx],
	})
}

// appendStem appends one shensha instance per palace where the target
// stem appears in the EarthStems layout.
// Precondition: stemIdx ∈ [0, 9] (always true for table-sourced values).
func appendStem(dst []ShenSha, kind Kind, stemIdx int, earth *[9]almanac.Stem) []ShenSha {
	stem := almanac.StemOf(stemIdx)
	for palace := uint8(1); palace <= 9; palace++ {
		if earth[palace-1] == stem {
			dst = append(dst, ShenSha{
				Kind:   kind,
				Target: Target{Stem: stem, IsStem: true},
				Palace: palace,
			})
		}
	}
	return dst
}
