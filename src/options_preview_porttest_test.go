//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestOptionsDialogPreview(t *testing.T) {
	type record struct {
		Start, Enabled, Busy, Missing, Step int
		Cursor                              uint32
		File                                string
	}
	var rows []record
	for start := 0; start < 3; start++ {
		for _, enabled := range []int{0, 1, 2} {
			for busy := 0; busy < 2; busy++ {
				for missing := 0; missing < 2; missing++ {
					t.Run(fmt.Sprintf("start%d/enabled%d/busy%d/missing%d", start, enabled, busy, missing), func(t *testing.T) {
						o := newOptionsAudioOwner(t, false)
						*o.words[122848] = uint32(enabled)
						if missing != 0 {
							var entries []strman.Entry
							for _, name := range []string{"OptionsPreviewA", "OptionsPreviewB", "OptionsPreviewC"} {
								entries = append(entries, strman.Entry{ID: strman.ID("Options.c:" + name), Vals: []strman.Variant{{Str: name}}})
							}
							configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
							t.Cleanup(restore)
							configure(0)
						}
						cursor := (*uint32)(unsafe.Pointer(&memmap.BlobByAddr(0x5D4594).Data[1309744]))
						*cursor = uint32(start)
						wantFile := ""
						if busy != 0 {
							legacy.Dialogs.PlayFile("already.wav", 100)
							if enabled != 0 {
								wantFile = "already.wav"
							}
						}
						wantCursor := uint32(start)
						for step := 0; step < 5; step++ {
							if step == 2 {
								legacy.Dialogs.Sub_44D8F0()
								wantFile = ""
							}
							if enabled == 0 || wantFile == "" {
								if enabled != 0 && missing == 0 {
									wantFile = []string{"OptionsPreviewA.wav", "OptionsPreviewB.wav", "OptionsPreviewC.wav"}[wantCursor]
								}
								wantCursor = (wantCursor + 1) % 3
							}
							legacy.PortTestOptionsAction(6, 0)
							if *cursor != wantCursor || legacy.Dialogs.FileToRead() != wantFile {
								t.Fatalf("step%d cursor=%d file=%q want=%d/%q", step, *cursor, legacy.Dialogs.FileToRead(), wantCursor, wantFile)
							}
							rows = append(rows, record{start, enabled, busy, missing, step, *cursor, legacy.Dialogs.FileToRead()})
						}
					})
				}
			}
		}
	}
	spellbookCapture(t, "options-preview", rows, "9714869c059c73ca48f4b525303cf07bf83193401a9e01a5c1c1b73d05aa96f5")
}
