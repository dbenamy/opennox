//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestConsoleCommandsGUI(t *testing.T) {
	for _, flags := range []uint32{0, 8, 8192, 8200} {
		t.Run(fmt.Sprint(flags), func(t *testing.T) {
			o := newConsoleCommandOwner(t)
			t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
			before := *o.optionWords["root"]
			if !o.call(t, "menu options", false) {
				t.Fatal("options return")
			}
			want := before
			if flags == 8192 {
				want = 0
			}
			if *o.optionWords["root"] != want {
				t.Fatal("options mode gate", *o.optionWords["root"], want)
			}
		})
	}
	t.Run("video", func(t *testing.T) {
		o := newConsoleCommandOwner(t)
		words, restore := legacy.PortTestOptionsWords()
		t.Cleanup(restore)
		root := o.c.GUI.NewWindowRaw(nil, 24, 0, 0, 160, 120, nil)
		sub := o.c.GUI.NewWindowRaw(nil, 24, 0, 0, 120, 80, nil)
		t.Cleanup(func() { root.Destroy(); sub.Destroy() })
		*words[1309820] = uint32(uintptr(root.C()))
		*words[1309824] = uint32(uintptr(sub.C()))
		oldPause := legacy.Sub_413A00
		paused := -1
		legacy.Sub_413A00 = func(v int) { paused = v }
		t.Cleanup(func() { legacy.Sub_413A00 = oldPause })
		if !o.call(t, "menu vidopt", false) || root.GetFlags().IsHidden() || sub.GetFlags().IsHidden() || paused != 1 || o.c.r.TabWidth() != 15 {
			t.Fatal("video options state")
		}
	})
	t.Run("motd", func(t *testing.T) {
		o := newConsoleCommandOwner(t)
		oldEngine := noxflags.GetEngine()
		noxflags.SetEngine(noxflags.EngineNoRendering)
		t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
		state := serverConfigOwnBytes(t, 0x5D4594, 826068, 4)
		binary.LittleEndian.PutUint32(state, 123)
		if !o.call(t, "show motd", false) || binary.LittleEndian.Uint32(state) != 0 {
			t.Fatal("MOTD dispatch")
		}
	})
}
