//go:build porttest

package opennox

import (
	"bytes"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

// Full-sized actions are covered by the frozen C captures. Every incomplete
// prefix must be rejected before accessing gameplay owners or changing input.
func TestGameMessageServerIncompleteActions(t *testing.T) {
	for _, action := range []struct {
		kind, sub byte
		size      int
	}{
		{64, 0, 7}, {114, 0, 7}, {115, 0, 3}, {116, 0, 3}, {117, 0, 3},
		{118, 0, 3}, {120, 0, 4}, {121, 0, 22}, {123, 0, 3}, {165, 0, 10},
		{224, 0, 3}, {226, 0, 4}, {241, 0, 3},
		{238, 0, 52}, {238, 1, 52}, {238, 2, 52}, {238, 3, 52}, {238, 4, 2}, {238, 5, 2},
		{201, 14, 2}, {201, 15, 4}, {201, 16, 4}, {201, 17, 2}, {201, 18, 2},
		{201, 21, 4}, {201, 22, 4}, {201, 23, 5}, {201, 24, 4}, {201, 25, 5},
		{201, 26, 4}, {201, 28, 4}, {201, 30, 4}, {240, 3, 2}, {240, 27, 2},
	} {
		for size := 0; size < action.size; size++ {
			data := make([]byte, size)
			if size > 0 {
				data[0] = action.kind
			}
			if size > 1 {
				data[1] = action.sub
			}
			before := bytes.Clone(data)
			if n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(0, data, nil, nil, nil); n != -1 || !bytes.Equal(data, before) {
				t.Fatalf("incomplete action kind=%d sub=%d size=%d: result=%d", action.kind, action.sub, size, n)
			}
		}
	}
}
