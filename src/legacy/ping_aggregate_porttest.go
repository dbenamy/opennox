//go:build porttest

package legacy

import (
	"bytes"
	"fmt"

	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/server"
)

type PortTestPingAggregateCall struct {
	Kind    string // min or average
	Wrapper bool   // exercise the public Go wrapper rather than the private helper
}

type PortTestPingTrace struct {
	Player int
	Value  int
}

type PortTestPingAggregateResult struct {
	Value            uint32
	Trace            []PortTestPingTrace
	PlayersUnchanged bool
}

// PortTestPingAggregates invokes either aggregate independently, with each
// Sub_554240 callback consuming its next per-player value. Exhausted sequences
// return zero, which makes repeated-read behavior explicit in test inputs.
func PortTestPingAggregates(slots []int, timings map[int][]int, calls []PortTestPingAggregateCall) ([]PortTestPingAggregateResult, error) {
	core, raw, freePlayers := server.PortTestPingServer(slots)
	defer freePlayers()
	before := append([]byte(nil), raw...)
	oldGet, oldTiming := GetServer, Sub_554240
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	defer func() { GetServer, Sub_554240 = oldGet, oldTiming }()
	pos := make(map[int]int)
	var trace []PortTestPingTrace
	Sub_554240 = func(ind ntype.PlayerInd) int {
		id := int(ind)
		i := pos[id]
		pos[id] = i + 1
		vals := timings[id]
		v := 0
		if i < len(vals) {
			v = vals[i]
		}
		trace = append(trace, PortTestPingTrace{Player: id, Value: v})
		return v
	}
	out := make([]PortTestPingAggregateResult, 0, len(calls))
	for _, call := range calls {
		start := len(trace)
		var value uint32
		switch call.Kind {
		case "min":
			if call.Wrapper {
				value = Sub_554290()
			} else {
				value = pingMinimum()
			}
		case "average":
			if call.Wrapper {
				value = Sub_554300()
			} else {
				value = pingAverage()
			}
		default:
			return nil, fmt.Errorf("unknown aggregate %q", call.Kind)
		}
		out = append(out, PortTestPingAggregateResult{
			Value:            value,
			Trace:            append([]PortTestPingTrace(nil), trace[start:]...),
			PlayersUnchanged: bytes.Equal(raw, before),
		})
	}
	return out, nil
}
