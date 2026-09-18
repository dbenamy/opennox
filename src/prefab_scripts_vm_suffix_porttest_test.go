//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"image"
	"testing"

	"github.com/opennox/opennox/v1/server"
)

func TestPrefabScriptsVMSuffixCoordinates(t *testing.T) {
	var rows []struct {
		Name     string
		Suffix   string
		Position image.Point
	}
	for _, xy := range [][2]int{{0, 0}, {-46, 92}, {123, -456}, {2147483647, -2147483648}} {
		name := fmt.Sprintf("OnEnter%%7%%%d%%%d", xy[0], xy[1])
		funcs := []prefabScriptFunction{{Name: "GLOBAL", Code: []uint32{72}}, {Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1}, Code: []uint32{72}}, {Name: name, Code: []uint32{72}}}
		var vm server.NoxScriptVM
		if err := vm.ReadScript(bytes.NewReader(prefabScriptsEncode(nil, funcs))); err != nil {
			t.Fatal(err)
		}
		fn := vm.Funcs()[2]
		want := image.Pt(xy[0], xy[1])
		if fn.NamePref != "%7" || fn.PosOff != want {
			t.Errorf("%s VM suffix=%q offset=%v want %v", name, fn.NamePref, fn.PosOff, want)
		}
		rows = append(rows, struct {
			Name     string
			Suffix   string
			Position image.Point
		}{name, fn.NamePref, fn.PosOff})
	}
	spellbookCapture(t, "prefab-scripts-vm-suffix", rows, "3cd5b4e1a22206e268fa54e58e1d861c6fd8e29f64187ee715582bb5e448f70f")
}
