//go:build porttest

package opennox

import (
	"encoding/json"
	"fmt"
	"image"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInventoryDisplayLocalizationOrder(t *testing.T) {
	if os.Getenv("OPENNOX_DISPLAY_RANDOM_CHILD") != "1" {
		exe, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(exe, "-test.run=^TestClientInventoryDisplayLocalizationOrder$", "-test.count=1", "-test.timeout=60s")
		cmd.Env = append(os.Environ(), "OPENNOX_DISPLAY_RANDOM_CHILD=1", "GODEBUG="+os.Getenv("GODEBUG")+",randautoseed=0,randseednop=0")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("isolated localization contract: %v\n%s", err, output)
		}
		return
	}
	// Keep process-global localization randomness inside this child, preserving
	// the parent suite's random state. The original C orders label, then item name.
	o := newInventoryDisplayOwner(t)
	prefixes := []string{"Inspect-A", "Inspect-B"}
	unknown := []string{"Unknown-A %S", "Unknown-B %S", "Unknown-C %S"}
	entries := []strman.Entry{}
	add := func(id string, values []string) {
		e := strman.Entry{ID: strman.ID(id)}
		for _, v := range values {
			e.Vals = append(e.Vals, strman.Variant{Str: v})
		}
		entries = append(entries, e)
	}
	add("guiinv.c:IdentifyItem", prefixes)
	add("ToolTip.c:NoArmsInfo", unknown)
	add("guiinv.c:IdentifyWeight", []string{"Weight %d"})
	add("guiinv.c:IdentifyDurabilityIndestructable", []string{"Indestructible"})
	add("guiinv.c:IdentifySpecialAttributes", []string{"Special attributes"})
	add("guiinv.c:IdentifyUnknown", []string{"Unknown"})
	data, err := json.Marshal(struct {
		Lang    int            `json:"lang"`
		Entries []strman.Entry `json:"entries"`
	}{0, entries})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "strings.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := o.c.srv.Strings().ReadJSON(path); err != nil {
		t.Fatal(err)
	}
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	for seed := int64(0); seed < 32; seed++ {
		o.reset(t)
		o.inventoryRects()
		label, _, _ := o.identifyWindows(t)
		o.c.Mouse = image.Pt(16, 22)
		dr := o.item(t, "RedApple", 123)
		dr.ObjClass = 0x10000000
		dr.ObjSubClass = 0
		*o.displayWords["dword_5d4594_1063116"] = uint32(uintptr(dr.C()))
		expected := rand.New(rand.NewSource(seed))
		want := fmt.Sprintf("%s %s", prefixes[expected.Intn(len(prefixes))], strings.ReplaceAll(unknown[expected.Intn(len(unknown))], "%S", "RedApple"))
		next := expected.Int63()
		rand.Seed(seed)
		o.call(t, int(seed), 3, txptr(unsafe.Pointer(&pos[0])), 0, 0)
		got := alloc.GoString16((*gui.StaticTextData)(label.WidgetData).Text)
		if got != want {
			t.Fatalf("seed%d label %q want %q (label lookup must precede item lookup)", seed, got, want)
		}
		if got := rand.Int63(); got != next {
			t.Fatalf("seed%d localization consumed unexpected randomness: %d want%d", seed, got, next)
		}
	}
}
