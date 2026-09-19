//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"math/rand"
	"testing"
)

func TestSessionServerFilterPredicates(t *testing.T) {
	version := legacy.PortTestSessionFilterVersion()
	rng := rand.New(rand.NewSource(0x4899c0))
	modes := []uint32{0, 1, 2, 3, 0xffffffff}
	pings := []uint32{0, 1, 100, 9998, 9999, 10000, 0x7fffffff, 0x80000000, 0xffffffff}
	nonzero := []uint32{0, 1, 2, 0x80000000}
	for i := 0; i < 8192; i++ {
		mode := modes[i%len(modes)]
		status := byte(rng.Uint32())
		video := byte(i)
		ping := pings[rng.Intn(len(pings))]
		limit := pings[rng.Intn(len(pings))]
		recordVersion := version
		if i&1 != 0 {
			recordVersion ^= 1 << uint(i%32)
		}
		var filter [11]uint32
		filter[0] = nonzero[rng.Intn(len(nonzero))]
		filter[1] = nonzero[rng.Intn(len(nonzero))]
		filter[2] = nonzero[rng.Intn(len(nonzero))]
		filter[3] = uint32(i % 130)
		filter[4] = limit
		filter[10] = nonzero[rng.Intn(len(nonzero))]
		var record [169]byte
		binary.LittleEndian.PutUint32(record[48:], recordVersion)
		binary.LittleEndian.PutUint32(record[96:], ping)
		record[100] = status
		record[102] = video
		want := true
		switch mode {
		case 1:
			want = status&0x30 == 0 && recordVersion == version
		case 2:
			want = (filter[0] == 0 || ping <= limit || ping == 9999) &&
				(filter[1] == 0 || status&0x10 == 0) && (filter[2] == 0 || status&0x20 == 0) &&
				(video&0x80 == 0 || uint32(video&0x7f) >= filter[3]) &&
				(filter[10] == 0 || recordVersion == version)
		}
		got, intact := legacy.PortTestSessionFilter(mode, filter, record)
		expected := 0
		if want {
			expected = 1
		}
		if !intact {
			t.Fatalf("case%d modified input record", i)
		}
		if got != expected {
			t.Fatalf("case%d mode%d status%02x video%02x ping%d version%x filters%v: got%d want%d", i, mode, status, video, ping, recordVersion, filter, got, expected)
		}
	}
}
