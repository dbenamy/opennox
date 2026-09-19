//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"unsafe"
)

func TestResourceDefinitionsRetainedABI(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	ud, freeUD := alloc.New(server.MonsterUpdateData{})
	defer freeUD()
	set, freeSet := alloc.New(uint32(99))
	defer freeSet()
	u.UpdateData = unsafe.Pointer(ud)
	if legacy.PortTestResourceMonsterSoundABI(nil) != nil {
		t.Fatal("nil ABI return")
	}
	for _, flags := range []uint32{0, 1, 2, 3, 4, 0x102, 0xffffffff} {
		for _, present := range []bool{false, true} {
			u.ObjClass = object.Class(flags)
			ud.SoundSet122 = nil
			if present {
				ud.SoundSet122 = unsafe.Pointer(set)
			}
			want := unsafe.Pointer(nil)
			if flags&2 != 0 && present {
				want = unsafe.Pointer(set)
			}
			if legacy.PortTestResourceMonsterSoundABI(u) != want {
				t.Fatal("retained C ABI", flags, present)
			}
		}
	}
}
func TestResourceDefinitionsFailedReaderCleanup(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified Linux target uses procfs for descriptor ownership")
	}
	resourceSoundOwner(t)
	resourceSoundFile(t, "A UNKNOWN invalid END ")
	path, err := filepath.Abs("soundset.bin")
	if err != nil {
		t.Fatal(err)
	}
	count := func() int {
		entries, err := os.ReadDir("/proc/self/fd")
		if err != nil {
			t.Fatal(err)
		}
		n := 0
		for _, e := range entries {
			dst, _ := os.Readlink("/proc/self/fd/" + e.Name())
			if dst == path {
				n++
			}
		}
		return n
	}
	before := count()
	for i := 0; i < 32; i++ {
		if legacy.Nox_xxx_parseSoundSetBin_424170("soundset.bin") != 0 {
			t.Fatal("unknown sound field admitted")
		}
		state := legacy.PortTestResourceSounds()
		if len(state) != 1 || state[0].Name != "A" {
			t.Fatal("partial soundset changed")
		}
		legacy.PortTestResourceSoundClear()
	}
	if got := count(); got != before {
		t.Fatalf("failed loads leaked readers: before%d after%d", before, got)
	}
}
