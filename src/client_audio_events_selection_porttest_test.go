//go:build porttest

package opennox

import (
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

type audioEventRandomOwner struct {
	legacy.Server
	s *server.Server
}

func (o *audioEventRandomOwner) S() *server.Server { return o.s }

func TestClientAudioEventsSampleSelection(t *testing.T) {
	s := new(server.Server)
	s.Rand.Logic = prand.New(17)
	s.Rand.Other = prand.New(31)
	old := legacy.GetServer
	legacy.GetServer = func() legacy.Server { return &audioEventRandomOwner{s: s} }
	defer func() { legacy.GetServer = old }()
	event, free := alloc.Make([]uint32{}, 144)
	defer free()
	meta, free := alloc.Make([]uint32{}, 50)
	defer free()
	entries, free := alloc.Make([]uint32{}, 21*32)
	defer free()
	ep, mp := audioEventPointer(unsafe.Pointer(&event[0])), audioEventPointer(unsafe.Pointer(&meta[0]))
	base := audioEventPointer(unsafe.Pointer(&entries[0]))
	var rows []map[string]any
	for _, n := range []int{0, 1, 2, 8, 32} {
		for _, flags := range []uint32{0, 1, 2, 3} {
			for _, limit := range []uint32{0, 1, 2, 3} {
				clear(event)
				clear(meta)
				event[9] = uint32(mp)
				event[42] = uint32(n)
				event[43] = 12345
				meta[1], meta[15] = flags, limit
				for i := 0; i < 32; i++ {
					event[10+i] = uint32(base) + uint32(84*i)
				}
				s.Rand.Other.Reset(31)
				rng := prand.New(31)
				s.Rand.Logic.Reset(17)
				next := 0
				remaining := n
				iteration := uint32(0)
				last := int32(12345)
				var bag []int
				reset := func() {
					bag = make([]int, n)
					for i := range bag {
						bag[i] = i
					}
					remaining = n
					next = 0
					if n != 0 {
						last = -1
					}
				}
				reset()
				var chosen []int32
				for step := 0; step < 2*n+3; step++ {
					name := "sub_451CF0"
					if step == 0 {
						name = "sub_451CA0"
					}
					if step > 0 && remaining == 0 && flags&1 != 0 {
						if limit != 0 {
							iteration++
						}
						if limit == 0 || iteration < limit {
							reset()
						}
					}
					want := uint64(0)
					selected := int32(-1)
					if remaining > 0 {
						if flags&2 != 0 {
							i := rng.Int(0, remaining-1)
							selected = int32(bag[i])
							bag = append(bag[:i], bag[i+1:]...)
						} else {
							selected = int32(next)
							next++
						}
						last = selected
						remaining--
						want = base + uint64(selected)*84 + 24
					}
					got := legacy.PortTestAudioEventCall(name, ep)
					if uint32(got) != uint32(want) || event[108] != uint32(remaining) || int32(event[43]) != last || event[109] != iteration || s.Rand.Other.Index() != rng.Index() || s.Rand.Logic.Index() != 17 {
						t.Fatal("sample cycle/loop/RNG", n, flags, limit, step, got, want, event[108], remaining, event[109], iteration)
					}
					chosen = append(chosen, selected)
				}
				rows = append(rows, map[string]any{"count": n, "flags": flags, "limit": limit, "selected": chosen, "iterations": iteration, "rng": rng.Index()})
			}
		}
	}
	spellbookCapture(t, "client-audio-events-sample-selection", rows, "12f95ad035cb3236b0df15137aacdb246d6c8001e387e415250c20ef6825e689")
}
