//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/spf13/viper"
)

type optionsOwner struct {
	*entryOwner
	root     *gui.Window
	controls map[int]*gui.Window
}

func newOptionsOwner(t *testing.T) *optionsOwner {
	o := &optionsOwner{entryOwner: newEntryOwner(t), controls: make(map[int]*gui.Window)}
	o.c.Inp = input.New(o.c.Log, &entrySeat{}, false, 0)
	o.root = o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 640, 480, nil)
	for _, id := range []int{311, 312, 313, 314, 316, 318, 331, 332, 333, 334, 341, 351, 352, 353, 361, 362, 363, 371, 999} {
		w := o.c.GUI.NewWindowRaw(o.root, 8, 0, 0, 10, 10, nil)
		w.SetID(uint(id))
		o.controls[id] = w
	}
	oldGamma, oldDirty, oldRO, oldCfg := getGamma(), configDirty, configReadOnly, viper.Get(configVideoGamma)
	oldCut, oldFull, oldScaled, oldRes := nox_video_cutSize, g_fullscreen_cfg, g_scaled_cfg, guiOptionsRes
	t.Cleanup(func() {
		nox_video_cutSize = oldCut
		g_fullscreen_cfg = oldFull
		g_scaled_cfg = oldScaled
		guiOptionsRes = oldRes
	})
	configReadOnly = true
	t.Cleanup(func() {
		setGamma(oldGamma)
		configDirty = oldDirty
		configReadOnly = oldRO
		viper.Set(configVideoGamma, oldCfg)
	})
	return o
}
