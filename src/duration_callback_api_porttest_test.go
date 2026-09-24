//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestDurationNativeCallbackDispatch(t *testing.T) {
	key := unsafe.Pointer(new(byte))
	record, free := alloc.New(server.DurSpell{})
	defer free()
	var value int32
	var received *server.DurSpell
	calls := 0
	server.RegisterDurSpellCallback(key, func(p *server.DurSpell) int32 {
		calls++
		received = p
		if p != nil {
			p.Field80 = 0x90000000 | uint32(calls)
		}
		return value
	})
	for _, word := range []int32{0, 1, 255, 256, 0x7fffffff, -0x80000000, -1} {
		value = word
		for _, p := range []*server.DurSpell{nil, record} {
			for _, discard := range []bool{false, true} {
				calls = 0
				received = nil
				record.Field80 = 0
				runtime.GC()
				if discard {
					server.CallDurSpellDiscard(key, p)
				} else if got := server.CallDurSpellResult(key, p); got != word {
					t.Fatal("native duration result", got, word)
				}
				if calls != 1 || received != p {
					t.Fatal("native duration forwarding", calls, received, p)
				}
				want := uint32(0)
				if p != nil {
					want = 0x90000001
				}
				if record.Field80 != want {
					t.Fatal("native duration state", record.Field80, want)
				}
			}
		}
	}
}
