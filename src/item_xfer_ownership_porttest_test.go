//go:build porttest

package opennox

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func TestItemXferGeneratorFailedChildOwnership(t *testing.T) {
	s := newItemXferOwner(t)
	typ := s.Types.ByID("PortInvisibleLight")
	// Match the earlier inventory rejection contract: real callback, no child
	// side buffers, rejection before payload access, actual pool ownership.
	typ.Weight = 255
	parent := newItemXferObject(t, s, "MonsterGenerator")
	before := s.Objs.Alive
	p := objectXferCurrentRecord(63)
	p.u8(3)
	p.Write(make([]byte, 3))
	p.u8(0)
	p.u8(0)
	p.u32(0)
	for i := 0; i < 4; i++ {
		p.u16(1)
		p.u32(0)
		p.u32(0)
	}
	p.u8(1)
	p.u8(1)
	name := "PortInvisibleLight"
	p.u8(byte(len(name)))
	p.WriteString(name)
	p.u16(0)
	p.u32(2)
	p.u16(61)
	path := filepath.Join(t.TempDir(), "rejected-generator-child.bin")
	if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
		t.Fatal(err)
	}
	defer cryptfile.Close()
	if parent.CallXfer(nil) == nil {
		t.Fatal("rejected child accepted")
	}
	if *(*unsafe.Pointer)(parent.UpdateData) != nil || s.Objs.Alive != before {
		t.Fatalf("failed generator child retained: live=%d want=%d", s.Objs.Alive, before)
	}
}
