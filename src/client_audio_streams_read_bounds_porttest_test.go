//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsReadBounds(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	cat, free := alloc.Make([]uint32{}, 72)
	defer free()
	entry, free := alloc.Make([]uint32{}, 9)
	defer free()
	dst, free := alloc.Make([]byte{}, 32)
	defer free()
	cp := audioStreamPointer(unsafe.Pointer(&cat[0]))
	cat[0] = audioStreamPointer(unsafe.Pointer(&entry[0]))
	entry[5] = 4
	files, _ := prefabScriptsFiles(t, []byte("abcd"))
	cat[67] = audioStreamPointer(files[0])
	dp := audioStreamPointer(unsafe.Pointer(&dst[0]))
	call := legacy.PortTestAudioStreamCall
	defer call("sub_486E00", cp)
	var rows []map[string]any
	for _, remaining := range []uint32{0, 1, 3, 12, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, request := range []uint32{0xffffffff, 0, 1, 4, 12, 0x7fffffff, 0x80000000} {
			if call("sub_486B60", cp, 0) != 1 {
				t.Fatal("open")
			}
			cat[71] = remaining
			for i := range dst {
				dst[i] = 0xcc
			}
			want := int32(request)
			if int32(remaining) < want {
				want = int32(remaining)
			}
			if want < 0 {
				want = 0
			}
			if want > int32(len(dst)) {
				continue
			}
			if want > 4 {
				want = 4
			}
			got := call("sub_486DB0", cp, dp, request)
			if got != uint32(want) || cat[71] != remaining-got || !bytes.Equal(dst[:got], []byte("abcd")[:got]) {
				t.Fatal("signed bounds/short read", remaining, request, got, want)
			}
			for _, v := range dst[got:] {
				if v != 0xcc {
					t.Fatal("short read guard")
				}
			}
			rows = append(rows, map[string]any{"remaining": remaining, "requested": request, "read": got, "left": cat[71]})
		}
	}
	spellbookCapture(t, "client-audio-streams-read-bounds", rows, "e6823b5cc70db53e5de59d293944ffd24ec0b0c998fea6deb49b289cd6b38a6b")
}
