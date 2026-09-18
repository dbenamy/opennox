//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestSessionEntryStateWords(t *testing.T) {
	mapState := serverConfigOwnBytes(t, 0x5D4594, 1523072, 4)
	loadState := serverConfigOwnBytes(t, 0x5D4594, 1563072, 4)
	scavenger := serverConfigOwnBytes(t, 0x5D4594, 1548508, 4)
	maximum := serverConfigOwnBytes(t, 0x5D4594, 1548528, 4)
	type row struct {
		Value                                           uint32
		MapSet, MapGet, LoadSet, Maximum                int32
		MapWord, LoadWord, ResetScavenger, ResetMaximum uint32
	}
	var rows []row
	for _, v := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff, 0x13579bdf} {
		binary.LittleEndian.PutUint32(scavenger, v)
		binary.LittleEndian.PutUint32(maximum, v)
		r := row{Value: v}
		r.MapSet = legacy.PortTestSessionEntryScalar("map-state-set", int32(v))
		r.MapGet = legacy.PortTestSessionEntryScalar("map-state-get", 0)
		r.LoadSet = legacy.PortTestSessionEntryScalar("load-state-set", int32(v))
		r.Maximum = legacy.PortTestSessionEntryScalar("scavenger-max", 0)
		legacy.PortTestSessionEntryScalar("scavenger-reset", 0)
		legacy.PortTestSessionEntryScalar("scavenger-max-reset", 0)
		r.MapWord = binary.LittleEndian.Uint32(mapState)
		r.LoadWord = binary.LittleEndian.Uint32(loadState)
		r.ResetScavenger = binary.LittleEndian.Uint32(scavenger)
		r.ResetMaximum = binary.LittleEndian.Uint32(maximum)
		if r.MapSet != int32(v) || r.MapGet != int32(v) || r.LoadSet != int32(v) || r.Maximum != int32(v) || r.MapWord != v || r.LoadWord != v || r.ResetScavenger != 0 || r.ResetMaximum != 0 {
			t.Fatalf("state round trip %+v", r)
		}
		rows = append(rows, r)
	}
	spellbookCapture(t, "session-entry-state-words", rows, "872947a3050a424edfab47ac3f0a3b7dabedc78b40c1ee63a60c430ccfad132a")
}
