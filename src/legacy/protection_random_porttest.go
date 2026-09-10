//go:build porttest

package legacy

import (
	"encoding/binary"
	"math"

	"github.com/opennox/opennox/v1/internal/protection"
)

type PortTestRandomOp struct {
	Mode int
	A, B uint32
}
type PortTestRandomSnapshot struct {
	Random protection.Random
	Result uint64
}

func portTestRandomState(r protection.Random) (state [40]byte, spanMaxMin [3]uint32) {
	for i, v := range r.State {
		binary.LittleEndian.PutUint64(state[i*8:], math.Float64bits(v))
	}
	return state, [3]uint32{r.Span, r.Max, r.Min}
}

func PortTestProtectionRandom(seed uint32, ops []PortTestRandomOp) []PortTestRandomSnapshot {
	old := protectionRandom
	defer func() { protectionRandom = old }()
	snapshot := func(result uint64) PortTestRandomSnapshot {
		return PortTestRandomSnapshot{Random: protectionRandom, Result: result}
	}
	protectionRandom.Seed(seed)
	out := []PortTestRandomSnapshot{snapshot(0)}
	for _, op := range ops {
		var result uint64
		switch op.Mode {
		case 0:
			result = math.Float64bits(protectionRandom.Next())
		case 1:
			result = uint64(protectionRandom.Range(op.A, op.B))
		case 2:
			result = uint64(protectionRandom.Draw())
		case 3:
			protectionRandom.Seed(op.A)
		default:
			panic("invalid random fixture mode")
		}
		out = append(out, snapshot(result))
	}
	return out
}
