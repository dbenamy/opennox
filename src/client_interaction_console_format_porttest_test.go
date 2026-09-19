//go:build porttest

package opennox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionConsoleFormatBoundaries(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	format, freeFormat := alloc.CString16("%s|%S|%08d")
	defer freeFormat()
	narrow, freeNarrow := alloc.CString("b\xff")
	defer freeNarrow()
	type row struct {
		Length int
		Number int32
		Result int
		Text   string
	}
	var rows []row
	for _, length := range []int{0, 1, 490, 499, 500, 501, 510, 511, 512, 513, 4096} {
		for _, number := range []int32{0, 1, 12345, 0x7fffffff} {
			text := strings.Repeat("w", length)
			wide, freeWide := alloc.CString16(text)
			o.console = nil
			ret := legacy.PortTestClientInteractionConsole(format, wide, narrow, number)
			freeWide()
			want := text + "|bÿ|" + fmt.Sprintf("%08d", number)
			if ret != 1 || len(o.console) != 1 || !strings.HasSuffix(o.console[0], want) {
				t.Fatal("console formatter boundary", length, number, ret, len(o.console))
			}
			rows = append(rows, row{length, number, ret, o.console[0]})
		}
	}
	interactionCapture(t, "console-format", rows)
}
