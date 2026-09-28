package pattern

import (
	"iter"

	"github.com/atopx/qimen/almanac"
	"github.com/atopx/qimen/element"
	"github.com/atopx/qimen/enum"
	"github.com/atopx/qimen/plate"
)

// DetectInput is the unpacked view of a chart used by detection. Using
// flat arrays here avoids a circular dependency between pattern and
// palace packages.
//
// Arrays are indexed by (palace - 1). The center palace is excluded
// from per-palace rules in both styles: its rotate-style stems are a
// constructional copy (not a rotation result), and classical pattern
// reading stays on the eight outer palaces.
type DetectInput struct {
	ZhiFuOriginalPalace uint8
	ZhiFuPalace         uint8

	EarthStems  [9]almanac.Stem
	HeavenStems [9]almanac.Stem

	Doors    [9]enum.Door
	DoorsSet [9]bool

	Gods    [9]enum.God
	GodsSet [9]bool

	Branches [9][]almanac.Branch

	KongWang [2]almanac.Branch
}

// AppendAll appends every 格局 found in the chart view to dst and
// returns the extended slice. This is the allocation-friendly core;
// Detect wraps it as a sequence.
func AppendAll(dst []Pattern, in *DetectInput) []Pattern {
	// Global patterns
	if in.ZhiFuOriginalPalace == in.ZhiFuPalace {
		dst = append(dst, Pattern{Kind: FuYin, Palace: in.ZhiFuPalace})
	} else if plate.AreOpposite(in.ZhiFuOriginalPalace, in.ZhiFuPalace) {
		dst = append(dst, Pattern{
			Kind:           FanYin,
			Palace:         in.ZhiFuPalace,
			OriginalPalace: in.ZhiFuOriginalPalace,
		})
	}

	// Per-palace patterns (eight outer palaces; see DetectInput).
	for i := range 9 {
		if n := uint8(i + 1); n != 5 {
			dst = appendPalace(dst, n, in)
		}
	}
	return dst
}

// Detect yields every 格局 found in the chart view. Use slices.Collect
// to materialize if needed.
func Detect(in DetectInput) iter.Seq[Pattern] {
	return func(yield func(Pattern) bool) {
		for _, p := range AppendAll(nil, &in) {
			if !yield(p) {
				return
			}
		}
	}
}

func appendPalace(dst []Pattern, n uint8, in *DetectInput) []Pattern {
	idx := n - 1
	heaven := in.HeavenStems[idx]
	earth := in.EarthStems[idx]
	door := in.Doors[idx]
	doorOK := in.DoorsSet[idx]
	god := in.Gods[idx]
	godOK := in.GodsSet[idx]

	// 入墓
	if plate.IsStemInTomb(heaven, n) {
		dst = append(dst, Pattern{Kind: RuMu, Palace: n, Stem: heaven})
	}

	// 落空亡 — both kongwang branches may match this palace independently
	for _, kb := range in.KongWang {
		for _, pb := range in.Branches[idx] {
			if pb == kb {
				dst = append(dst, Pattern{Kind: KongWang, Palace: n, Branch: kb})
				break
			}
		}
	}

	// 门迫: door element 受克 by palace element (palace.element clobbers door.element)
	if doorOK && element.OfDoor(door).RelationTo(element.FromPalace(n)) == element.Restrained {
		dst = append(dst, Pattern{Kind: MenPo, Palace: n, Door: door})
	}

	// 三奇得使
	if doorOK {
		if heaven == almanac.Yi && door == enum.DoorOpen {
			dst = append(dst, Pattern{Kind: YiQiDeShi, Palace: n})
		}
		if heaven == almanac.Bing && door == enum.DoorRest {
			dst = append(dst, Pattern{Kind: BingQiDeShi, Palace: n})
		}
		if heaven == almanac.Ding && door == enum.DoorLife {
			dst = append(dst, Pattern{Kind: DingQiDeShi, Palace: n})
		}
	}

	// 八遁
	if doorOK {
		if heaven == almanac.Bing && earth == almanac.Ding && door == enum.DoorLife {
			dst = append(dst, Pattern{Kind: TianDun, Palace: n})
		}
		if heaven == almanac.Yi && door == enum.DoorOpen && godOK && god == enum.GodJiuDi {
			dst = append(dst, Pattern{Kind: DiDun, Palace: n})
		}
		if heaven == almanac.Ding && door == enum.DoorRest && godOK && god == enum.GodTaiYin {
			dst = append(dst, Pattern{Kind: RenDun, Palace: n})
		}
		if heaven == almanac.Bing && door == enum.DoorLife && godOK && god == enum.GodJiuTian {
			dst = append(dst, Pattern{Kind: ShenDun, Palace: n})
		}
		if heaven == almanac.Ding && door == enum.DoorBlock && godOK && god == enum.GodJiuDi {
			dst = append(dst, Pattern{Kind: GuiDun, Palace: n})
		}
		if heaven == almanac.Yi && (door == enum.DoorOpen || door == enum.DoorBlock) && n == 4 {
			dst = append(dst, Pattern{Kind: FengDun, Palace: n})
		}
		if heaven == almanac.Yi && door == enum.DoorOpen && n == 6 {
			dst = append(dst, Pattern{Kind: YunDun, Palace: n})
		}
		if heaven == almanac.Yi && door == enum.DoorRest && n == 1 {
			dst = append(dst, Pattern{Kind: LongDun, Palace: n})
		}
		if heaven == almanac.Yi && door == enum.DoorOpen && n == 7 {
			dst = append(dst, Pattern{Kind: HuDun, Palace: n})
		}
	}

	// 青龙返首 / 飞鸟跌穴
	if heaven == almanac.Wu && earth == almanac.Bing {
		dst = append(dst, Pattern{Kind: QingLongFanShou, Palace: n})
	}
	if heaven == almanac.Bing && earth == almanac.Wu {
		dst = append(dst, Pattern{Kind: FeiNiaoDieXue, Palace: n})
	}

	// 凶格
	if heaven == almanac.Geng && earth == almanac.Gui {
		dst = append(dst, Pattern{Kind: DaGe, Palace: n})
	}
	if heaven == almanac.Geng && earth == almanac.Ren {
		dst = append(dst, Pattern{Kind: XiaoGe, Palace: n})
	}
	if heaven == almanac.Geng && earth == almanac.Ji {
		dst = append(dst, Pattern{Kind: XingGe, Palace: n})
	}

	// 悖格
	if (heaven == almanac.Bing && earth == almanac.Geng) || (heaven == almanac.Geng && earth == almanac.Bing) {
		dst = append(dst, Pattern{Kind: BoGe, Palace: n})
	}

	// 天网四张
	if heaven == almanac.Gui && earth == almanac.Gui {
		dst = append(dst, Pattern{Kind: TianWangSiZhang, Palace: n})
	}

	return dst
}
