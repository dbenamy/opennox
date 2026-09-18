//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func prefabScriptsResetFile(t *testing.T, f *binfile.File, data []byte) {
	t.Helper()
	if err := f.Truncate(0); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
}

func TestPrefabScriptsInstructions(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	gl, restore := legacy.PortTestPrefabScriptsGlobals()
	defer restore()
	raw, files := prefabScriptsFiles(t, nil, nil)
	type row struct {
		Name     string
		Return   uint64
		Position int64
		Output   []byte
		Globals  [3]uint32
	}
	var rows []row
	for op := uint32(0); op <= 74; op++ {
		values := []uint32{0}
		if op <= 6 || op >= 19 && op <= 21 || op == 70 {
			values = []uint32{0, 3, 4, 127, 128, 255, 256, 0x7fffffff, 0x80000000, 0xffffffff}
		}
		if op == 5 {
			values = []uint32{0, 0x80000000, 0x3f800000, 0xbf800000, 1, 0x00800000, 0x7f7fffff, 0x7f800000, 0xff800000, 0x7fc01234}
		}
		if op == 69 {
			values = make([]uint32, 212)
			for i := range values {
				values[i] = uint32(i)
			}
		}
		flags := []uint32{0}
		if op <= 3 {
			flags = []uint32{0, 1, 0xffffffff}
		}
		configs := [][3]uint32{{9, 257, 13}}
		if op <= 3 || op == 6 || op == 70 {
			configs = append(configs, [3]uint32{0xffffffff, 0xfffffffd, 0xfffffffe})
		}
		for _, cfg := range configs {
			for _, flag := range flags {
				for _, value := range values {
					for _, remap := range []uint32{0, 1, 0xffffffff} {
						name := fmt.Sprintf("op%d/value%x/flag%x/remap%x/cfg%x", op, value, flag, remap, cfg)
						prefix := []uint32{4, 101, 4, 202, 4, 303}
						code := append([]uint32(nil), prefix...)
						code = append(code, op)
						switch {
						case op <= 3:
							code = append(code, flag, value)
						case op <= 6 || op >= 19 && op <= 21 || op == 69 || op == 70:
							code = append(code, value)
						}
						if op != 72 {
							code = append(code, 72)
						}
						if op < 74 {
							if _, err := asm.Decode(code); err != nil {
								t.Fatalf("%s: fixture decode: %v", name, err)
							}
						}
						want := append([]uint32(nil), code...)
						if remap != 0 {
							switch {
							case op <= 2 && flag != 0 && int32(value) >= 4:
								want[8] += cfg[0] - 4
							case op == 3 && flag != 0:
								want[8] += cfg[1]
							case op == 6:
								want[7] += cfg[1]
							case op == 70:
								want[7] += cfg[2] - 2
							case op == 69:
								switch value {
								case 9, 10, 46, 47, 126, 190:
									want[5] += cfg[2] - 2
								}
								if value == 126 {
									want[3] += cfg[2] - 2
								}
							}
						}
						result := uint64(1)
						position := int64(len(code) * 4)
						if op == 74 {
							result = 74
							position = 28
							want = prefix
						}
						input := append(prefabScriptsWords(code...), 0xde, 0xad, 0xbe, 0xef)
						prefabScriptsResetFile(t, files[0], input)
						prefabScriptsResetFile(t, files[1], nil)
						for i := range gl {
							*gl[i] = cfg[i]
						}
						ret := legacy.PortTestPrefabScriptsCall(4, raw[0], raw[1], nil, remap)
						pos, err := files[0].Seek(0, io.SeekCurrent)
						if err != nil {
							t.Fatal(err)
						}
						out := prefabScriptsOutput(t, files[1])
						after := [3]uint32{*gl[0], *gl[1], *gl[2]}
						if ret != result || pos != position || !bytes.Equal(out, prefabScriptsWords(want...)) || after != cfg {
							t.Fatalf("%s return=%x position=%d output=%x expected=%x globals=%x", name, ret, pos, out, prefabScriptsWords(want...), after)
						}
						rows = append(rows, row{name, ret, pos, out, after})
					}
				}
			}
		}
	}
	spellbookCapture(t, "prefab-scripts-instructions", rows, "4fea75fdbe7823de766a93eb683ab1372f3af885ef9b7153012511bf693db4db")
}

type prefabScriptFunction struct {
	Name         string
	Vars, Code   []uint32
	Return, Args uint32
}

func prefabScriptsEncode(strings []string, funcs []prefabScriptFunction) []byte {
	var b bytes.Buffer
	word := func(v uint32) { b.Write(prefabScriptsWords(v)) }
	str := func(s string) { word(uint32(len(s))); b.WriteString(s) }
	b.WriteString("SCRIPT03STRG")
	word(uint32(len(strings)))
	for _, s := range strings {
		str(s)
	}
	b.WriteString("CODE")
	word(uint32(len(funcs)))
	for _, f := range funcs {
		b.WriteString("FUNC")
		str(f.Name)
		word(f.Return)
		word(f.Args)
		b.WriteString("SYMB")
		word(uint32(len(f.Vars)))
		word(0)
		b.Write(prefabScriptsWords(f.Vars...))
		b.WriteString("DATA")
		word(uint32(len(f.Code) * 4))
		b.Write(prefabScriptsWords(f.Code...))
	}
	b.WriteString("DONE")
	return b.Bytes()
}
func TestPrefabScriptsCompleteMerge(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	gl, restore := legacy.PortTestPrefabScriptsGlobals()
	defer restore()
	raw, files := prefabScriptsFiles(t, nil, nil, nil)
	type row struct {
		Name      string
		Return    uint64
		Positions [3]int64
		Output    []byte
		Globals   [3]uint32
	}
	var rows []row
	for _, extraA := range []int{0, 1, 4} {
		for _, extraB := range []int{0, 1, 3} {
			first := prefabScriptFunction{Name: "GLOBAL", Code: []uint32{72}}
			a := []prefabScriptFunction{first, {Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1, 2, 3}, Code: []uint32{4, 65536, 72}}}
			b := []prefabScriptFunction{first, {Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1, 4}, Code: []uint32{6, 0, 72}}}
			for i := 0; i < extraA; i++ {
				a = append(a, prefabScriptFunction{Name: fmt.Sprintf("A%d", i), Vars: []uint32{1, 2}, Code: []uint32{4, 0x12345678, 72}})
			}
			for i := 0; i < extraB; i++ {
				b = append(b, prefabScriptFunction{Name: fmt.Sprintf("B%d", i), Vars: []uint32{1}, Code: []uint32{6, 1, 70, 2, 72}})
			}
			ab := prefabScriptsEncode([]string{"first", "second"}, a)
			bb := prefabScriptsEncode([]string{"third", "fourth", "fifth"}, b)
			for _, data := range [][]byte{ab, bb} {
				if _, err := asm.ReadScript(bytes.NewReader(data)); err != nil {
					t.Fatal("invalid input fixture", err)
				}
			}
			prefabScriptsResetFile(t, files[0], ab)
			prefabScriptsResetFile(t, files[1], bb)
			prefabScriptsResetFile(t, files[2], nil)
			for _, p := range gl {
				*p = 0xa5a5a5a5
			}
			ret := legacy.PortTestPrefabScriptsCall(9, raw[0], raw[1], raw[2], 0)
			var pos [3]int64
			for i := range files {
				var err error
				pos[i], err = files[i].Seek(0, io.SeekCurrent)
				if err != nil {
					t.Fatal(err)
				}
			}
			out := prefabScriptsOutput(t, files[2])
			sc, err := asm.ReadScript(bytes.NewReader(out))
			if err != nil {
				t.Fatalf("merged script is unreadable: %v", err)
			}
			if ret != 4 || pos[0] != int64(len(ab)) || pos[1] != int64(len(bb)) || pos[2] != int64(len(out)) {
				t.Fatal("merge return/file positions", ret, pos)
			}
			if !reflect.DeepEqual(sc.Strings, []string{"first", "second", "third", "fourth", "fifth"}) || len(sc.Funcs) != len(a)+len(b)-2 {
				t.Fatal("merged table counts")
			}
			if len(sc.Funcs[1].Vars) != 7 || !reflect.DeepEqual(sc.Funcs[1].Code, []uint32{4, 65536, 6, 2, 72}) {
				t.Fatalf("global merge %+v", sc.Funcs[1])
			}
			for i := 0; i < extraA; i++ {
				if !reflect.DeepEqual(sc.Funcs[2+i].Code, a[2+i].Code) {
					t.Fatal("first-input code changed")
				}
			}
			for i := 0; i < extraB; i++ {
				f := sc.Funcs[2+extraA+i]
				if f.Name != b[2+i].Name || !reflect.DeepEqual(f.Code, []uint32{6, 3, 70, uint32(2 + extraA), 72}) {
					t.Fatalf("second-input remap %+v", f)
				}
			}
			counters := [3]uint32{*gl[0], *gl[1], *gl[2]}
			if counters != ([3]uint32{6, 2, uint32(len(a))}) {
				t.Fatal("merge counters", counters)
			}
			rows = append(rows, row{fmt.Sprintf("extra%d/%d", extraA, extraB), ret, pos, out, counters})
		}
	}
	spellbookCapture(t, "prefab-scripts-merge", rows, "5a16ebfaf33c3731ef58412af035f8d79559c4de46cb6fe5b3f70afddfd22f89")
}
