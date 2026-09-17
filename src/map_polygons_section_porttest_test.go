//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

// Independent wire builder: old vertex IDs deliberately differ from loaded slots.
func mapPolygonWire(version, count int, editor bool, sequential bool) []byte {
	var b bytes.Buffer
	put := func(v any) {
		if err := binary.Write(&b, binary.LittleEndian, v); err != nil {
			panic(err)
		}
	}
	put(uint16(version))
	put(uint32(3 * count))
	for i := 0; i < 3*count; i++ {
		id := uint32(50 + i*7)
		if sequential {
			id = uint32(i + 1)
		}
		put(id)
		x := float32(100 + 30*(i/3))
		y := float32(200)
		switch i % 3 {
		case 0:
			x += 0.75
			y += 0.25
		case 1:
			x += 10.5
			y += 0.75
		case 2:
			x += 3.25
			y += 10.5
		}
		put(math.Float32bits(x))
		put(math.Float32bits(y))
		if version < 3 {
			put(uint32(0x98765432))
		}
	}
	put(uint32(count))
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("Region%d", i)
		put(uint8(len(name)))
		b.WriteString(name)
		rgb := []byte{byte(17 + i), byte(89 + i), byte(211 - i)}
		if version < 3 {
			put(uint32(0x11223344))
			for _, c := range rgb {
				put(uint32(c) | 0x98760000)
			}
		} else {
			b.Write(rgb)
		}
		put(uint8(5 + i))
		put(uint16(3))
		for j := 0; j < 3; j++ {
			id := uint32(50 + (3*i+j)*7)
			if sequential {
				id = uint32(3*i + j + 1)
			}
			put(id)
		}
		if version >= 2 {
			for j := 0; j < 2; j++ {
				put(uint16(1))
				name := ""
				if editor {
					name = fmt.Sprintf("handler%d-%d", i, j)
				}
				put(uint32(len(name)))
				b.WriteString(name)
				put(uint32(2 + j))
			}
		}
		if version >= 4 {
			put(uint32(0x12340000 + i))
		}
	}
	return b.Bytes()
}

func TestMapPolygonsSectionVersions(t *testing.T) {
	old := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(old) })
	words, restore := legacy.PortTestQuestRuntimeGlobals()
	t.Cleanup(restore)
	path := filepath.Join(t.TempDir(), "polygons.bin")
	type row struct {
		Name              string
		Records, Metadata [][]byte
		Vertices          []byte
		Stage, Previous   uint32
		Wire              []byte
	}
	var rows []row
	for _, version := range []int{1, 2, 3, 4} {
		for _, count := range []int{0, 2} {
			for _, editor := range []bool{false, true} {
				name := fmt.Sprintf("v%d/count%d/editor%t", version, count, editor)
				t.Run(name, func(t *testing.T) {
					o := newMapPolygonsOwner(t)
					if editor {
						noxflags.PortTestGameFlags(0x200000)
					}
					*words["previousStage"] = 73
					*words["1556132"] = 91
					raw := mapPolygonWire(version, count, editor, false)
					if err := os.WriteFile(path, raw, 0600); err != nil {
						t.Fatal(err)
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					ret := legacy.PortTestMapPolygonSection(0)
					pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
					if err != nil {
						t.Fatal(err)
					}
					cryptfile.Close()
					if ret != 1 || pos != int64(len(raw)) {
						t.Fatalf("read return %d position %d want %d", ret, pos, len(raw))
					}
					wantSecrets := uint32(0)
					if version == 4 && count == 2 {
						wantSecrets = 1
					}
					if *words["previousStage"] != wantSecrets || *words["1556132"] != 73 {
						t.Fatal("secret count state", *words["previousStage"], *words["1556132"])
					}
					r := row{Name: name, Stage: *words["previousStage"], Previous: *words["1556132"], Vertices: bytes.Clone(o.vertices[16 : 16*(3*count+1)])}
					for i := 1; i <= count; i++ {
						p := o.polygon(i)
						if binary.LittleEndian.Uint32(p[80:]) != uint32(i) || binary.LittleEndian.Uint32(p[84:]) != 1 || binary.LittleEndian.Uint16(p[128:]) != 3 || p[130] != byte(i+4) {
							t.Fatal("loaded region identity/count/level")
						}
						for j, want := range []int32{int32(100 + 30*(i-1)), 200, int32(110 + 30*(i-1)), 210} {
							if int32(binary.LittleEndian.Uint32(p[88+4*j:])) != want {
								t.Fatal("integer bounds", i, j)
							}
						}
						ids := unsafe.Slice(*(**uint32)(unsafe.Pointer(&p[108])), 3)
						for j, id := range ids {
							if id != uint32(3*(i-1)+j+1) {
								t.Fatal("vertex ID remapping", ids)
							}
						}
						copyP := bytes.Clone(p)
						clear(copyP[:4])
						clear(copyP[108:112])
						r.Records = append(r.Records, copyP)
						if meta := *(*unsafe.Pointer)(unsafe.Pointer(&p[0])); meta != nil {
							r.Metadata = append(r.Metadata, bytes.Clone(unsafe.Slice((*byte)(meta), 256)))
						}
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
						t.Fatal(err)
					}
					if legacy.PortTestMapPolygonSection(0) != 1 {
						t.Fatal("write")
					}
					if err := cryptfile.Close(); err != nil {
						t.Fatal(err)
					}
					r.Wire, err = os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					// Version 4 round trips must exactly match the independent modern wire
					// builder after remapping. Earlier versions retain default flags/callbacks.
					if version == 4 && !bytes.Equal(r.Wire, mapPolygonWire(4, count, editor, true)) {
						t.Fatalf("round trip wire\ngot %x\nwant %x", r.Wire, mapPolygonWire(4, count, editor, true))
					}
					rows = append(rows, r)
				})
			}
		}
	}
	t.Run("bypass-and-version", func(t *testing.T) {
		for _, version := range []uint16{5, 32767} {
			raw := []byte{byte(version), byte(version >> 8), 0x91, 0x92}
			if err := os.WriteFile(path, raw, 0600); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			if legacy.PortTestMapPolygonSection(1) != 1 {
				t.Fatal("bypass")
			}
			pos, _ := cryptfile.Global().File.Seek(0, io.SeekCurrent)
			if pos != 0 {
				t.Fatal("bypass consumed bytes")
			}
			if legacy.PortTestMapPolygonSection(0) != 0 {
				t.Fatal("unsupported version")
			}
			pos, _ = cryptfile.Global().File.Seek(0, io.SeekCurrent)
			if pos != 2 {
				t.Fatal("version position", pos)
			}
			cryptfile.Close()
		}
	})
	spellbookCapture(t, "map-polygons-section-versions", rows, "bace02065cb0367da8838352caccf569aafccb988a08f301efafba0574f7bd32")
}
