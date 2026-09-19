//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"os"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionDollAssets(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	o.reset(t)
	blob, err := os.ReadFile("common/memmap/nox/blobdata/blob_587000.dat")
	if err != nil {
		t.Fatal(err)
	}
	type group struct {
		Offset, Count, Dest int
		Names               []string
	}
	groups := []group{{180960, 2, 16, nil}, {180968, 2, 24, nil}, {180976, 52, 32, nil}, {181184, 54, 256, nil}}
	outputs := serverConfigOwnBytes(t, 0x973A20, 16, 456)
	cloak := serverConfigOwnBytes(t, 0x5D4594, 1319052, 4)
	ids := map[string]int{}
	for gi := range groups {
		g := &groups[gi]
		table := serverConfigOwnBytes(t, 0x587000, uintptr(g.Offset), g.Count*4)
		for i := 0; i < g.Count; i++ {
			address := binary.LittleEndian.Uint32(blob[g.Offset+4*i:])
			name := ""
			if address != 0 {
				off := int(address - 0x587000)
				end := bytes.IndexByte(blob[off:], 0)
				if end < 0 {
					t.Fatal("doll asset name terminator")
				}
				name = string(blob[off : off+end])
				p, free := alloc.CString(name)
				defer free()
				binary.LittleEndian.PutUint32(table[4*i:], uint32(uintptr(unsafe.Pointer(p))))
			}
			g.Names = append(g.Names, name)
			if name != "" {
				ids[name] = 1 + len(ids)%len(o.images)
			}
		}
	}
	ids["MaleMedievalCloakTop"] = 7
	var called []string
	oldLoad := legacy.Nox_xxx_gLoadImg
	defer func() { legacy.Nox_xxx_gLoadImg = oldLoad }()
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		called = append(called, name)
		if name == "" {
			return nil
		}
		id, ok := ids[name]
		if !ok {
			t.Fatal("unexpected doll asset", name)
		}
		return o.images[id-1]
	}
	if interactionCall("sub_4BFAD0") != 1 {
		t.Fatal("doll asset initialization")
	}
	var want []string
	for gender := 0; gender < 2; gender++ {
		want = append(want, groups[0].Names[gender], groups[1].Names[gender])
		want = append(want, groups[2].Names[26*gender:26*(gender+1)]...)
		want = append(want, groups[3].Names[27*gender:27*(gender+1)]...)
	}
	want = append(want, "MaleMedievalCloakTop")
	if !slices.Equal(called, want) || len(called) != 111 {
		t.Fatal("doll asset load order", len(called))
	}
	for _, g := range groups {
		for i, name := range g.Names {
			value := binary.LittleEndian.Uint32(outputs[g.Dest-16+4*i:])
			var expected uint32
			if name != "" {
				expected = uint32(uintptr(o.images[ids[name]-1].C()))
			}
			if value != expected {
				t.Fatal("doll asset table", g.Offset, i, name)
			}
		}
	}
	if binary.LittleEndian.Uint32(cloak) != uint32(uintptr(o.images[6].C())) {
		t.Fatal("cloak overlay owner")
	}
	interactionCapture(t, "doll-assets", called)
}
