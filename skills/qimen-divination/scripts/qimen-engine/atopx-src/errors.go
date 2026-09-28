package qimen

import "github.com/atopx/qimen/almanac"

// ErrInvalidTime is returned by FromTime / FromTimestamp for invalid
// solar-time inputs. It is the same sentinel as [almanac.ErrInvalidTime]
// (the almanac layer is where time validation happens), so errors.Is
// matches against either name.
//
// Chart construction itself is total: every Method / Style / JuRule is
// implemented, so [From] and [New] have no error path.
var ErrInvalidTime = almanac.ErrInvalidTime
