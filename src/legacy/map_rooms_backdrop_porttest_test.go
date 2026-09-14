//go:build porttest

package legacy

import (
	"testing"
	"unsafe"

	"github.com/opennox/libs/platform"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestMapOrchestrationBackdropWeights(t *testing.T) {
	old := platform.Get()
	platform.Set(platform.New())
	defer platform.Set(old)
	r := mapRoomNew(13, 13)
	defer mapRoomFree(r)
	list, freeList := alloc.New([8]uint32{})
	defer freeList()
	a, freeA := alloc.New([56]uint32{})
	defer freeA()
	b, freeB := alloc.New([56]uint32{})
	defer freeB()
	list[0] = mapRoomRaw(unsafe.Pointer(a))
	a[55] = mapRoomRaw(unsafe.Pointer(b))
	a[18] = 1000
	b[18] = 1000
	a[20] = 99
	b[20] = 99
	// A phase-specific entry must not contribute to an unflagged backdrop's draw.
	a[16] = 1
	for seed := uint32(0); seed < 64; seed++ {
		mapRoomSeed(seed)
		if got := mapRoomSelectDecoration(r, unsafe.Pointer(list)); got != unsafe.Pointer(b) {
			t.Fatalf("unrestricted selection seed %d", seed)
		}
	}
	a[16] = 0
	selected := map[unsafe.Pointer]bool{}
	for seed := uint32(0); seed < 64; seed++ {
		mapRoomSeed(seed)
		got := mapRoomSelectDecoration(r, unsafe.Pointer(list))
		if got != unsafe.Pointer(a) && got != unsafe.Pointer(b) {
			t.Fatal("weighted selection outside eligible list")
		}
		selected[got] = true
	}
	if len(selected) != 2 {
		t.Fatal("weighted selection did not reach both eligible entries")
	}
	a[16] = 1
	b[16] = 2
	if mapRoomSelectDecoration(r, unsafe.Pointer(list)) != nil {
		t.Fatal("ineligible list selected a backdrop")
	}
	list[0] = 0
	if mapRoomSelectDecoration(r, unsafe.Pointer(list)) != nil {
		t.Fatal("empty list selected a backdrop")
	}
}
