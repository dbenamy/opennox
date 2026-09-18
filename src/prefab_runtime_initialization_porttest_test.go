//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestPrefabRuntimeInitialization(t *testing.T) {
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	sizes := map[string]int{"path": 2048, "alternate": 2048, "metadata": 2048 * 76}
	if legacy.PortTestPrefabCall(7, [6]uint32{}) != 0 {
		t.Fatal("empty initialization return")
	}
	defer func() {
		for key := range sizes {
			legacy.PortTestPrefabReleasePayload(unsafe.Pointer(uintptr(*words[key])))
			*words[key] = 0
		}
	}()
	pointers := map[string]uint32{}
	for key, size := range sizes {
		ptr := *words[key]
		if ptr == 0 {
			t.Fatal("missing allocated buffer", key)
		}
		pointers[key] = ptr
		for _, v := range unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size) {
			if v != 0 {
				t.Fatal("nonzero new buffer", key)
			}
		}
	}
	*words["count"] = 123
	if legacy.PortTestPrefabCall(7, [6]uint32{}) != 0 || *words["count"] != 0 {
		t.Fatal("empty repeat")
	}
	for key, p := range pointers {
		if *words[key] != p {
			t.Fatal("reallocated existing buffer", key)
		}
	}
}
