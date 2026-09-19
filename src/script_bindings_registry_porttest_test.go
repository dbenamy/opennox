//go:build porttest

package opennox

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestScriptBindingsStringRegistry(t *testing.T) {
	newWorldCollisionOwner(t)
	region := serverConfigOwnBytes(t, 0x973F18, 26660, 4104)
	for i := range region {
		region[i] = 0xa5
	}
	words := unsafe.Slice(memmap.PtrUint32(0x973F18, 26664), 1024)
	clear(words)
	counter, restore := legacy.PortTestScriptBindingRegistryOwner()
	t.Cleanup(restore)
	*counter = 0
	type row struct {
		Call, Index int
		Count       uint32
		Text        string
	}
	var rows []row
	for i := 0; i < 1027; i++ {
		text := []string{"", "café Ω", "𐐷", "before\x00after"}[i%4] + fmt.Sprint(i)
		p, free := alloc.CString16(text)
		t.Cleanup(free)
		before := append([]uint32(nil), words...)
		got := legacy.PortTestScriptBindingIntern(unsafe.Pointer(p))
		want := i
		if want > 1023 {
			want = 1023
		}
		count := uint32(i + 1)
		if count > 1024 {
			count = 1024
		}
		if got != want || *counter != count {
			t.Fatal("registry index/capacity", i, got, *counter)
		}
		for j := range words {
			w := before[j]
			if i < 1024 && i == j {
				w = uint32(uintptr(unsafe.Pointer(p)))
			}
			if words[j] != w {
				t.Fatal("registry changed wrong slot", i, j)
			}
		}
		stored := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(words[got]))))
		if i < 1024 && stored != strings.SplitN(text, "\x00", 2)[0] {
			t.Fatal("registry UTF16 content")
		}
		rows = append(rows, row{i, got, *counter, stored})
	}
	before := append([]uint32(nil), words...)
	noxServer.noxScript.Reset()
	if *counter != 0 {
		t.Fatal("script reset did not reset string counter")
	}
	for i := range words {
		if words[i] != before[i] {
			t.Fatal("C reset unexpectedly cleared borrowed string slot")
		}
	}
	p, free := alloc.CString16("after reset")
	defer free()
	if got := legacy.PortTestScriptBindingIntern(unsafe.Pointer(p)); got != 0 || *counter != 1 || words[0] != uint32(uintptr(unsafe.Pointer(p))) {
		t.Fatal("registry reset/reuse")
	}
	for _, b := range append(append([]byte(nil), region[:4]...), region[4100:]...) {
		if b != 0xa5 {
			t.Fatal("registry crossed table boundary")
		}
	}
	spellbookCapture(t, "script-bindings-string-registry", rows, "db24e916dad9dd1ff8a0430df084ba77f9933977e1b36c8af8bc58d550223d73")
}
