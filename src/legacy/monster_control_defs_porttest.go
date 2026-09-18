//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestMonsterDefsOwner() func() {
	old := monsterDefinitions
	monsterDefinitions = nil
	cb := portTestCallbackTablesEnvironment()
	files.Lock()
	oldFiles := files.byHandle
	files.byHandle = make(map[unsafe.Pointer]*binfile.File)
	files.Unlock()
	return func() {
		monsterDefinitionFree()
		monsterDefinitions = old
		cb()
		files.Lock()
		defer files.Unlock()
		for _, f := range files.byHandle {
			_ = f.Close()
		}
		files.byHandle = oldFiles
	}
}
func PortTestMonsterDefs(op string, id int) int {
	switch op {
	case "load":
		return monsterDefinitionLoad()
	case "free":
		return int(monsterDefinitionFree())
	case "bind":
		return monsterDefinitionBind()
	case "lookup":
		p := unsafe.Pointer(monsterDefinitionByType(uint32(id)))
		if p == nil {
			return 0
		}
		for i, d := 1, (*server.MonsterDef)(monsterDefinitions); d != nil; i, d = i+1, d.Next244 {
			if unsafe.Pointer(d) == p {
				return i
			}
		}
		panic("definition lookup outside owned list")
	default:
		panic(op)
	}
}

// Copy every definition byte, identifying callback slots and list links rather
// than retaining process-specific pointers. No parser result is reconstructed.
func PortTestMonsterDefsSnapshot() [][62]uint32 {
	var out [][62]uint32
	seen := map[*server.MonsterDef]bool{}
	for d := (*server.MonsterDef)(monsterDefinitions); d != nil; d = d.Next244 {
		if seen[d] {
			panic("definition list cycle")
		}
		seen[d] = true
		words := *(*[62]uint32)(unsafe.Pointer(d))
		for _, off := range []int{57, 58, 59} {
			if words[off] == 0 {
				continue
			}
			found := false
			for _, base := range []uintptr{287096, 287192, 287280} {
				for j := uintptr(0); *memmap.PtrPtr(0x587000, base+8*j) != nil; j++ {
					if words[off] == *memmap.PtrUint32(0x587000, base+8*j+4) {
						words[off] = uint32(base + 8*j)
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				panic("unknown definition callback")
			}
		}
		words[61] = 0
		if d.Next244 != nil {
			words[61] = uint32(len(out) + 2)
		}
		out = append(out, words)
	}
	return out
}

type PortTestMonsterToken struct {
	Return int
	Buffer []byte
	Intact bool
}

func PortTestMonsterTokens(f *binfile.Binfile) []PortTestMonsterToken {
	h := NewFileHandle(f.File)
	defer nox_fs_close(h)
	b, free := alloc.Make([]byte{}, 272)
	defer free()
	var out []PortTestMonsterToken
	for i := 0; i < 1024; i++ {
		for j := range b {
			b[j] = 0xa5
		}
		rv := bool2int(monsterDefinitionToken(f, b[8:264]))
		r := PortTestMonsterToken{Return: rv, Buffer: bytes.Clone(b[8:264]), Intact: true}
		for _, v := range append(bytes.Clone(b[:8]), b[264:]...) {
			r.Intact = r.Intact && v == 0xa5
		}
		out = append(out, r)
		if rv == 0 {
			return out
		}
	}
	panic("token fixture did not reach EOF")
}
