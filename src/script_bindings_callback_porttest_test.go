//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestScriptBindingsCallbackRead(t *testing.T) {
	o := newWorldCollisionOwner(t)
	funcs := []prefabScriptFunction{{Name: "GLOBAL", Code: []uint32{72}}, {Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1}, Code: []uint32{72}}, {Name: "OnEnter", Code: []uint32{72}}}
	if err := o.s.NoxScriptVM.ReadScript(bytes.NewReader(prefabScriptsEncode(nil, funcs))); err != nil {
		t.Fatal(err)
	}
	record, free := alloc.New([2]uint32{})
	t.Cleanup(free)
	buffer, free := alloc.Make([]byte{}, 1028)
	t.Cleanup(free)
	path := filepath.Join(t.TempDir(), "callback.bin")
	type row struct {
		Version  uint16
		Flags    uint32
		Name     string
		Declared uint32
		Result   int
		Consumed int64
		Words    [2]uint32
		Buffer   []byte
	}
	var rows []row
	for _, version := range []uint16{0, 1, 2, 0x7fff, 0x8000, 0xffff} {
		for _, flags := range []uint32{0, 0x200000, 0x400000, 0x600000} {
			for _, name := range []string{"", "OnEnter", "Missing", "OnEnter\x00ignored", strings.Repeat("x", 1023), strings.Repeat("x", 1024)} {
				for _, oversize := range []bool{false, true} {
					declared := uint32(len(name))
					if oversize {
						declared = 0xffffffff
					}
					t.Run(fmt.Sprintf("v%x/flags%x/len%d/oversize%v", version, flags, len(name), oversize), func(t *testing.T) {
						t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
						*record = [2]uint32{0x13579bdf, 0x2468ace0}
						for i := range buffer {
							buffer[i] = 0xa5
						}
						var stream bytes.Buffer
						binary.Write(&stream, binary.LittleEndian, version)
						binary.Write(&stream, binary.LittleEndian, declared)
						stream.WriteString(name)
						binary.Write(&stream, binary.LittleEndian, uint32(0x89abcdef))
						stream.WriteString("TAIL")
						if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
							t.Fatal(err)
						}
						if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
							t.Fatal(err)
						}
						defer cryptfile.Close()
						got := legacy.PortTestScriptBindingCallback(unsafe.Pointer(record), unsafe.Pointer(&buffer[0]))
						pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
						if err != nil {
							t.Fatal(err)
						}
						wantResult, wantPos := 0, int64(2)
						wantWords := [2]uint32{0x13579bdf, 0x2468ace0}
						wantBuffer := bytes.Repeat([]byte{0xa5}, len(buffer))
						if int16(version) <= 1 {
							wantPos = 6
							if declared < 1024 {
								wantResult, wantPos = 1, int64(10+len(name))
								wantWords[0] = 0x89abcdef
								if len(name) != 0 {
									cname := strings.SplitN(name, "\x00", 2)[0]
									if flags != 0 {
										copy(wantBuffer, cname)
										wantBuffer[len(cname)] = 0
									} else if cname == "OnEnter" {
										wantWords[1] = 2
									} else {
										wantWords[1] = 0xffffffff
									}
								}
							}
						}
						if got != wantResult || pos != wantPos || *record != wantWords || !bytes.Equal(buffer, wantBuffer) {
							t.Fatalf("callback read result/offset/record/buffer mismatch: result %d/%d offset %d/%d words %x/%x", got, wantResult, pos, wantPos, *record, wantWords)
						}
						rows = append(rows, row{version, flags, name, declared, got, pos, *record, bytes.Clone(buffer)})
					})
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-callback-read", rows, "bd7fbd8df288c3e279580cb388b179dba9314a86ae214f89dd63e0df09c65a46")
}

func TestScriptBindingsCallbackWrite(t *testing.T) {
	o := newWorldCollisionOwner(t)
	funcs := []prefabScriptFunction{{Name: "GLOBAL", Code: []uint32{72}}, {Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1}, Code: []uint32{72}}, {Name: "OnEnter", Code: []uint32{72}}, {Name: "", Code: []uint32{72}}}
	if err := o.s.NoxScriptVM.ReadScript(bytes.NewReader(prefabScriptsEncode(nil, funcs))); err != nil {
		t.Fatal(err)
	}
	record, free := alloc.New([2]uint32{})
	t.Cleanup(free)
	path := filepath.Join(t.TempDir(), "callback.bin")
	var rows [][]byte
	for _, flags := range []uint32{0, 0x200000, 0x400000, 0x600000} {
		for _, index := range []uint32{0xffffffff, 2, 3} {
			for _, name := range []string{"", "OnEnter", "OnEnter\x00ignored", strings.Repeat("x", 1023), strings.Repeat("y", 1024), strings.Repeat("z", 1025)} {
				for _, nilName := range []bool{false, true} {
					t.Run(fmt.Sprintf("flags%x/index%x/len%d/nil%v", flags, index, len(name), nilName), func(t *testing.T) {
						t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
						*record = [2]uint32{0x89abcdef, index}
						text, free := alloc.CString(name)
						defer free()
						p := unsafe.Pointer(text)
						if nilName {
							p = nil
						}
						if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
							t.Fatal(err)
						}
						got := legacy.PortTestScriptBindingCallback(unsafe.Pointer(record), p)
						cryptfile.Close()
						data, err := os.ReadFile(path)
						if err != nil {
							t.Fatal(err)
						}
						wantName := ""
						if flags != 0 {
							if !nilName {
								wantName = strings.SplitN(name, "\x00", 2)[0]
							}
						} else if index == 2 {
							wantName = "OnEnter"
						}
						var want bytes.Buffer
						binary.Write(&want, binary.LittleEndian, uint16(1))
						binary.Write(&want, binary.LittleEndian, uint32(len(wantName)))
						want.WriteString(wantName)
						binary.Write(&want, binary.LittleEndian, uint32(0x89abcdef))
						if got != 1 || !bytes.Equal(data, want.Bytes()) || *record != [2]uint32{0x89abcdef, index} {
							t.Fatal("callback write result/bytes/record mismatch")
						}
						if !bytes.Equal(unsafe.Slice(text, len(name)+1), append([]byte(name), 0)) {
							t.Fatal("callback write changed name")
						}
						rows = append(rows, data)
					})
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-callback-write", rows, "dad06a5baf9d7ce4362fcddccfe1a6f6edd30d71f38e324c1a8c211afc98123a")
}
