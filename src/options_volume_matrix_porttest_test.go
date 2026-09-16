//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/timer"
)

func TestOptionsVolumeMatrix(t *testing.T) {
	type record struct {
		Menu                                                            bool
		Event, Channel, Enabled, Ready, Callback, Active, Value, Return int
		State                                                           [3]uint32
		Checked                                                         bool
		Timer                                                           timer.Timer
		Sounds                                                          [][2]int
		Dialog                                                          string
		PreviewIndex                                                    uint32
	}
	var rows []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprint(menu), func(t *testing.T) {
			o := newOptionsAudioOwner(t, menu)
			activeWindow := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 1, 1, nil)
			for _, event := range []int{16393, 16396} {
				for ch := 0; ch < 3; ch++ {
					for _, enabled := range []int{0, 1, 2} {
						for ready := 0; ready < 2; ready++ {
							for callback := 0; callback < 2; callback++ {
								for active := 0; active < 2; active++ {
									for _, value := range []int{0, 1, 8192, 16384, -1} {
										cursor := (*uint32)(unsafe.Pointer(&memmap.BlobByAddr(0x5D4594).Data[1309744]))
										*cursor = 0
										legacy.Dialogs.Sub_44D8F0()
										o.sounds = nil
										o.c.GUI.WinYYY = nil
										*o.words[831092] = uint32(ready)
										*o.words[816376] = uint32(ready)
										for n, off := range []int{126996, 122848, 93156} {
											*o.words[off] = uint32(enabled)
											button := o.buttons[n][callback]
											button.DrawData().Field0 = 0
											if enabled == 1 {
												button.DrawData().Field0 = 4
											}
											*o.words[o.buttonOffset(n)] = uint32(uintptr(button.C()))
											*o.timers[n] = timer.Timer{Flags: 4, Current: 0x12340000, Target: 0x43210000}
										}
										if active != 0 {
											activeWindow.SetParent(o.controls[351+ch])
											o.c.GUI.WinYYY = activeWindow
										}
										ret := legacy.PortTestOptionsEvent(menu, o.root, event, o.controls[351+ch], value)
										handled := event == 16393 || ch < 2
										suppress := event == 16393 && ch == 0 && active != 0
										toggle := handled && !suppress && callback != 0 && ((value != 0 && enabled == 0) || (value == 0 && enabled == 1))
										wantState := uint32(enabled)
										if toggle {
											if enabled == 1 {
												wantState = 0
											} else if ch == 0 || ready != 0 {
												wantState = 1
											}
										}
										var state [3]uint32
										for n, off := range []int{126996, 122848, 93156} {
											state[n] = *o.words[off]
											want := uint32(enabled)
											if n == ch {
												want = wantState
											}
											if state[n] != want {
												t.Fatalf("state menu=%v event=%d ch=%d enabled=%d ready=%d callback=%d active=%d value=%d: channel%d=%d want=%d", menu, event, ch, enabled, ready, callback, active, value, n, state[n], want)
											}
										}
										checked := o.buttons[ch][callback].DrawData().Field0&4 != 0
										wantChecked := enabled == 1
										if toggle {
											wantChecked = !wantChecked
										}
										if checked != wantChecked {
											t.Fatal("checkbox selection differs")
										}
										wantTimer := timer.Timer{Flags: 4, Current: 0x12340000, Target: 0x43210000}
										if handled {
											wantTimer.Flags |= 1
											wantTimer.Target = uint32(value) << 16
										}
										if toggle && ch == 2 && enabled != 1 {
											wantTimer.Target = wantTimer.Current
										}
										if *o.timers[ch] != wantTimer {
											t.Fatalf("timer event=%d ch=%d value=%d got=%+v want=%+v", event, ch, value, *o.timers[ch], wantTimer)
										}
										var sounds [][2]int
										if toggle {
											sounds = append(sounds, [2]int{921, 100})
										}
										if handled && !suppress && ch == 0 && value != 0 {
											sounds = append(sounds, [2]int{768, 100})
										}
										if !reflect.DeepEqual(o.sounds, sounds) {
											t.Fatalf("preview sounds got=%v want=%v", o.sounds, sounds)
										}
										dialog := ""
										previewIndex := uint32(0)
										if handled && ch == 1 && value != 0 {
											previewIndex = 1
											if wantState != 0 {
												dialog = "OptionsPreviewA.wav"
											}
										}
										if legacy.Dialogs.FileToRead() != dialog || *cursor != previewIndex {
											t.Fatalf("dialog=%q index=%d want=%q/%d", legacy.Dialogs.FileToRead(), *cursor, dialog, previewIndex)
										}
										if ret != 0 {
											t.Fatalf("unexpected return %d", ret)
										}
										rows = append(rows, record{menu, event, ch, enabled, ready, callback, active, value, ret, state, checked, *o.timers[ch], append([][2]int(nil), o.sounds...), legacy.Dialogs.FileToRead(), *cursor})
									}
								}
							}
						}
					}
				}
			}
		})
	}
	spellbookCapture(t, "options-volume", rows, "7e720133b8f493b07fd671d6c3636d7bedaeb057efc3c770e29279f427839d8e")
}
