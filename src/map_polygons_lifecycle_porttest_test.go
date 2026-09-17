//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"os/exec"
	"testing"
	"unsafe"
)

func TestMapPolygonsResetAndStorage(t *testing.T) {
	o := newMapPolygonsOwner(t)
	if *o.words["polygons"] != 1 || *o.words["vertices"] != 1 || *o.words["remap"] != 0 || binary.LittleEndian.Uint32(o.control[4:]) != 1 {
		t.Fatal("reset counters")
	}
	for i := 1; i < 256; i++ {
		p := o.polygon(i)
		if binary.LittleEndian.Uint32(p[116:]) != 0xffffffff || binary.LittleEndian.Uint32(p[124:]) != 0xffffffff || binary.LittleEndian.Uint16(p[128:]) != 0 {
			t.Fatal("reset record", i)
		}
	}
	var rows []any
	for _, index := range []int{1, 7, 1023} {
		p := legacy.PortTestMapPolygonVertex(float32(index)+0.25, -float32(index)-0.5, index, 1000+index)
		if p != unsafe.Pointer(&o.vertex(index)[0]) || o.vertexID(p) != index || binary.LittleEndian.Uint32(o.vertex(index)[12:]) != 1 {
			t.Fatal("vertex fields", index)
		}
		rows = append(rows, []any{index, append([]byte(nil), o.vertex(index)...), *o.words["vertices"]})
	}
	var found []int
	for p := legacy.PortTestMapPolygonPointer("vertex-first", nil, 0); p != nil; p = legacy.PortTestMapPolygonPointer("vertex-next", p, 0) {
		found = append(found, o.vertexID(p))
	}
	if len(found) != 3 || found[0] != 1 || found[1] != 7 || found[2] != 1023 {
		t.Fatal("sparse iteration", found)
	}
	if got := legacy.PortTestMapPolygonScalar("free-vertex", nil); got != 2 {
		t.Fatal("free slot", got)
	}
	for _, index := range []int{0, 1, 7, 1023} {
		if legacy.PortTestMapPolygonPointer("vertex-get", nil, index) != unsafe.Pointer(&o.vertex(index)[0]) {
			t.Fatal("vertex address", index)
		}
	}
	if legacy.PortTestMapPolygonPointer("remap-clear", nil, 0) != nil || *o.words["remap"] != 0 {
		t.Fatal("remap clear")
	}
	result := legacy.PortTestMapPolygonPointer("reset-vertices", nil, 0)
	if result != unsafe.Add(unsafe.Pointer(&o.polygons[0]), 12) || *o.words["vertices"] != 1 {
		t.Fatal("reset vertex result")
	}
	rows = append(rows, found, append([]byte(nil), o.vertex(7)...))
	spellbookCapture(t, "map-polygons-reset-storage", rows, "1d97ea3a13709c53c99cb70bc3ff7e27c72087097735d070d6fbfea3f6efac7b")
}
func TestMapPolygonsRecordDefaults(t *testing.T) {
	o := newMapPolygonsOwner(t)
	var rows []any
	for i := 1; i <= 4; i++ {
		p := legacy.PortTestMapPolygonPointer("new", nil, 0)
		data := o.polygon(i)
		if p != o.pointer(i) || o.polygonID(p) != i || binary.LittleEndian.Uint32(data[84:]) != 1 || data[104] != 117 || data[105] != 38 || data[106] != 219 || data[130] != 9 {
			t.Fatal("record defaults", i, data)
		}
		// Omit allocation addresses from the record; this non-editor path has no metadata.
		rows = append(rows, append([]byte(nil), data...))
	}
	binary.LittleEndian.PutUint32(o.polygon(2)[84:], 0)
	if legacy.PortTestMapPolygonScalar("free-polygon", nil) != 2 {
		t.Fatal("reuse inactive slot")
	}
	var found []int
	for p := legacy.PortTestMapPolygonPointer("first", nil, 0); p != nil; p = legacy.PortTestMapPolygonPointer("next", p, 0) {
		found = append(found, o.polygonID(p))
	}
	if len(found) != 3 || found[0] != 1 || found[1] != 3 || found[2] != 4 {
		t.Fatal("record iteration", found)
	}
	for _, id := range []int{0, 1, 2, 3, 4, 5, 255, -559023410} {
		got := legacy.PortTestMapPolygonPointer("get", nil, id)
		if id == -559023410 {
			if got != nil {
				t.Fatal("sentinel lookup")
			}
		} else if got != o.pointer(id) {
			t.Fatal("raw record address", id)
		}
		rows = append(rows, []int{id, o.polygonID(got)})
	}
	p := o.polygon(3)
	copy(p[4:], "Changed\x00")
	p[104] = 13
	p[105] = 17
	p[106] = 19
	p[130] = 21
	binary.LittleEndian.PutUint32(p[132:], 0x12345678)
	binary.LittleEndian.PutUint32(p[136:], 0xffffffff)
	if legacy.PortTestMapPolygonScalar("defaults-write", o.pointer(3)) != 0 || binary.LittleEndian.Uint32(p[132:]) != 0 || binary.LittleEndian.Uint32(p[136:]) != 0 {
		t.Fatal("default update")
	}
	if legacy.PortTestMapPolygonScalar("defaults-read", o.pointer(1)) != 0 || o.polygon(1)[104] != 13 || o.polygon(1)[130] != 21 {
		t.Fatal("default copy")
	}
	rows = append(rows, found, append([]byte(nil), o.defaults...), append([]byte(nil), o.polygon(1)...), append([]byte(nil), p...))
	spellbookCapture(t, "map-polygons-record-defaults", rows, "529907abf39663b16a3d01a69300471ec33c2a081b9422c40562f07f22ef97e9")
}
func TestMapPolygonsNilActors(t *testing.T) {
	if os.Getenv("OPENNOX_POLYGON_NIL_CHILD") == "" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestMapPolygonsNilActors$", "-test.count=1")
		cmd.Env = append(os.Environ(), "OPENNOX_POLYGON_NIL_CHILD=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("nil actor child: %v\n%s", err, out)
		}
		return
	}
	legacy.PortTestMapPolygonActor("player", nil)
	legacy.PortTestMapPolygonActor("monster", nil)
}

func TestMapPolygonsBorrowedReset(t *testing.T) {
	o := newMapPolygonsOwner(t)
	noxflags.PortTestGameFlags(0x200000)
	o.construct(t, [][2]float32{{100, 200}, {110, 200}, {105, 210}})
	p := o.polygon(1)
	meta := *(*unsafe.Pointer)(unsafe.Pointer(&p[0]))
	ids := *(*unsafe.Pointer)(unsafe.Pointer(&p[108]))
	metadata := unsafe.Slice((*byte)(meta), 256)
	vertices := unsafe.Slice((*byte)(ids), 12)
	metadata[0] = 0x91
	metadata[255] = 0xe7
	wantMeta, wantVertices := bytes.Clone(metadata), bytes.Clone(vertices)
	// Return detached allocations to the fixture's real C cleanup afterward.
	defer func() {
		*(*unsafe.Pointer)(unsafe.Pointer(&p[0])) = meta
		*(*unsafe.Pointer)(unsafe.Pointer(&p[108])) = ids
	}()
	binary.LittleEndian.PutUint32(o.control[4:], 0)
	legacy.PortTestMapPolygonPointer("reset-records", nil, 0)
	if *o.words["polygons"] != 1 || binary.LittleEndian.Uint32(p[0:]) != 0 || binary.LittleEndian.Uint32(p[108:]) != 0 || binary.LittleEndian.Uint16(p[128:]) != 0 {
		t.Fatal("borrowed reset record")
	}
	if !bytes.Equal(metadata, wantMeta) || !bytes.Equal(vertices, wantVertices) {
		t.Fatal("borrowed allocation mutated")
	}
}
