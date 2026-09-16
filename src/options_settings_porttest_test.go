//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/opennox/libs/cfg"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/timer"
	"github.com/spf13/viper"
)

func TestOptionsSettingsRoundTrip(t *testing.T) {
	type record struct {
		Menu                    bool
		Mask, SensitivitySlider int
		Volumes                 [3]int
		Sensitivity, Gamma      uint32
	}
	var rows []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprint(menu), func(t *testing.T) {
			o := newOptionsAudioOwner(t, menu)
			oldVolumes := [3]int{configGetVolume(VolumeFX), configGetVolume(VolumeDialog), configGetVolume(VolumeMusic)}
			t.Cleanup(func() {
				for i, v := range oldVolumes {
					configSetVolume(v, VolumeControl(i))
				}
			})
			dir := t.TempDir()
			for mask := 0; mask < 8; mask++ {
				for _, slider := range []int{0, 50, 100} {
					var volumes [3]int
					for ch, off := range []int{126996, 122848, 93156} {
						*o.words[off] = 1
						button := o.buttons[ch][1]
						button.DrawData().Field0 = 4
						*o.words[o.buttonOffset(ch)] = uint32(uintptr(button.C()))
						*o.timers[ch] = timer.Timer{Current: 8192 << 16, Target: 8192 << 16}
						if mask&(1<<ch) != 0 {
							volumes[ch] = 1234 + ch*456
						}
					}
					for ch, v := range volumes {
						legacy.PortTestOptionsEvent(menu, o.root, 16393, o.controls[351+ch], v)
						o.timers[ch].Update()
					}
					legacy.PortTestOptionsEvent(menu, o.root, 16393, o.controls[318], slider)
					legacy.PortTestOptionsEvent(menu, o.root, 16393, o.controls[316], slider)
					sensitivity, gamma := o.c.GetSensitivity(), getGamma()
					var all, selected cfg.Section
					writeConfigLegacyMain(&all)
					for ch, key := range []string{"FXVolume", "DialogVolume", "MusicVolume"} {
						v, ok := all.Get(key)
						if !ok || v != strconv.Itoa(volumes[ch]) {
							t.Fatalf("saved %s=%q want=%d", key, v, volumes[ch])
						}
						selected.Set(key, v)
					}
					v, ok := all.Get("InputSensitivity")
					if !ok {
						t.Fatal("missing saved sensitivity")
					}
					selected.Set("InputSensitivity", v)
					path := filepath.Join(dir, "options.cfg")
					f, err := os.Create(path)
					if err != nil {
						t.Fatal(err)
					}
					err = (&cfg.File{Sections: []cfg.Section{selected}}).WriteTo(f)
					closeErr := f.Close()
					if err != nil {
						t.Fatal(err)
					}
					if closeErr != nil {
						t.Fatal(closeErr)
					}
					o.c.SetSensitivity(7)
					for ch := 0; ch < 3; ch++ {
						configSetVolume(9, VolumeControl(ch))
					}
					f, err = os.Open(path)
					if err != nil {
						t.Fatal(err)
					}
					err = parseLegacyConfig(f, false)
					closeErr = f.Close()
					if err != nil {
						t.Fatal(err)
					}
					if closeErr != nil {
						t.Fatal(closeErr)
					}
					for ch, want := range volumes {
						if got := configGetVolume(VolumeControl(ch)); got != want {
							t.Fatalf("loaded channel%d=%d want=%d", ch, got, want)
						}
					}
					if math.Float32bits(o.c.GetSensitivity()) != math.Float32bits(sensitivity) {
						t.Fatal("sensitivity did not round-trip exactly")
					}
					path = filepath.Join(dir, "modern.json")
					if err = viper.WriteConfigAs(path); err != nil {
						t.Fatal(err)
					}
					loaded := viper.New()
					loaded.SetConfigFile(path)
					if err = loaded.ReadInConfig(); err != nil {
						t.Fatal(err)
					}
					if math.Float32bits(float32(loaded.GetFloat64(configVideoGamma))) != math.Float32bits(gamma) {
						t.Fatal("modern gamma did not round-trip exactly")
					}
					rows = append(rows, record{menu, mask, slider, volumes, math.Float32bits(sensitivity), math.Float32bits(gamma)})
				}
			}
		})
	}
	spellbookCapture(t, "options-settings", rows, "f670776d0b563031d9a3b29e29ede8768a58c119ef35a9ab58e1214648e18617")
}
