//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

// Complete records are compared with frozen C results. These independent wire
// contracts require every incomplete prefix to fail before gameplay access.
func TestGameMessageClientSessionIncompleteRecords(t *testing.T) {
	records := []struct {
		op, sub byte
		size    int
	}{
		{166, 0, 4}, {167, 0, 4}, {171, 0, 5}, {174, 0, 3}, {175, 0, 20}, {176, 0, 49}, {177, 0, 60},
		{178, 0, 4}, {179, 0, 4}, {180, 0, 4}, {181, 0, 14}, {189, 0, 2}, {195, 0, 12},
		{196, 0, 18}, {196, 1, 10}, {196, 2, 6}, {196, 3, 10}, {196, 4, 46}, {196, 5, 6}, {196, 6, 6}, {196, 7, 2}, {196, 8, 10}, {196, 9, 2}, {196, 12, 5},
		{197, 0, 1}, {198, 0, 1},
		{201, 1, 2}, {201, 2, 2}, {201, 3, 3}, {201, 4, 15}, {201, 5, 4}, {201, 6, 14}, {201, 7, 2}, {201, 8, 18}, {201, 9, 4}, {201, 12, 52}, {201, 13, 86}, {201, 27, 4}, {201, 29, 8}, {201, 31, 8},
		{202, 0, 3}, {203, 0, 1}, {204, 0, 4}, {205, 0, 3}, {206, 0, 3}, {207, 0, 3}, {209, 0, 3}, {210, 0, 7}, {211, 0, 13},
		{213, 1, 68}, {213, 2, 68}, {213, 3, 68}, {214, 0, 3}, {215, 0, 5}, {216, 0, 6}, {217, 0, 4}, {218, 0, 4}, {219, 0, 5}, {220, 0, 3},
		{223, 0, 6}, {224, 0, 4}, {225, 0, 3}, {226, 0, 4}, {229, 0, 3}, {230, 0, 3}, {231, 0, 3}, {232, 0, 3}, {233, 0, 9}, {234, 0, 5}, {235, 0, 2}, {237, 0, 2}, {238, 6, 3}, {238, 7, 2},
		{240, 0, 2}, {240, 1, 4}, {240, 2, 14}, {240, 4, 5}, {240, 5, 4}, {240, 6, 4}, {240, 7, 4}, {240, 8, 4}, {240, 9, 4}, {240, 10, 4}, {240, 11, 16}, {240, 12, 90}, {240, 13, 69}, {240, 14, 69}, {240, 15, 4}, {240, 16, 12}, {240, 17, 4}, {240, 18, 4}, {240, 19, 4}, {240, 20, 2}, {240, 21, 8}, {240, 22, 5}, {240, 23, 5}, {240, 24, 3}, {240, 25, 7}, {240, 26, 6}, {240, 28, 3}, {240, 29, 4}, {240, 30, 5}, {240, 31, 5}, {240, 32, 5}, {240, 33, 52},
	}
	check := func(op byte, record []byte) {
		t.Helper()
		for size := 0; size < len(record); size++ {
			data := bytes.Clone(record[:size])
			before := bytes.Clone(data)
			if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(op), data); n != -1 || !bytes.Equal(data, before) {
				t.Fatalf("op=%d size=%d/%d returned %d", op, size, len(record), n)
			}
		}
	}
	for _, r := range records {
		data := make([]byte, r.size)
		data[0] = r.op
		if r.size > 1 {
			data[1] = r.sub
		}
		check(r.op, data)
	}
	for count := 1; count <= 255; count++ {
		team := make([]byte, 18+count*2)
		team[0] = 196
		team[15] = byte(count)
		check(196, team)
		seq := make([]byte, 4+count)
		seq[0] = 204
		seq[3] = byte(count)
		check(204, seq)
	}
}
