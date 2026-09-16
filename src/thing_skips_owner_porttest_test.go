//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"os"
	"testing"
	"unsafe"
)

type thingSkipResult struct {
	Op, Offset, Return, Consumed int
	Scratch                      [32]byte
}

func thingSkipCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	b, e := json.Marshal(rows)
	if e != nil {
		t.Fatal(e)
	}
	if p := os.Getenv("OPENNOX_THING_SKIPS_CAPTURE"); p != "" {
		if e := os.WriteFile(p+"-"+label+".json", b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s %s", label, got)
	if got != want {
		t.Fatalf("%s got%s want%s", label, got, want)
	}
}
func thingSkipWord(b []byte, v uint32) []byte { return binary.LittleEndian.AppendUint32(b, v) }
func thingSkipFill(b []byte, n int, seed byte) []byte {
	for i := 0; i < n; i++ {
		b = append(b, seed+byte(i%17))
	}
	return b
}
func thingSkipName(b []byte, n int, seed byte) []byte {
	b = append(b, byte(n))
	return thingSkipFill(b, n, seed)
}
func thingSkipRef(b []byte, inline bool, n int) []byte {
	if !inline {
		return thingSkipWord(b, 0x12345678)
	}
	b = thingSkipWord(b, 0xffffffff)
	b = append(b, 0xf3)
	return thingSkipName(b, n, 0x31)
}
func thingSkipInvoke(t *testing.T, op, offset int, payload []byte, wantRet, wantConsumed int, wantScratch []byte) thingSkipResult {
	t.Helper()
	input := bytes.Repeat([]byte{0xa7}, offset)
	input = append(input, payload...)
	input = append(input, bytes.Repeat([]byte{0xc3}, 16)...)
	raw, _ := alloc.CloneSlice(input)
	f := binfile.NewMemFile(unsafe.Pointer(&raw[0]), len(raw))
	defer f.Free()
	f.Skip(offset)
	scratch, free := alloc.Make([]byte{}, 256*1024)
	defer free()
	for i := range scratch {
		scratch[i] = 0xa5
	}
	ret := legacy.PortTestThingSkip(op, f, scratch)
	consumed := len(raw) - len(f.Data()) - offset
	if !bytes.Equal(raw, input) {
		t.Fatal("input mutated")
	}
	if ret != wantRet || consumed != wantConsumed {
		t.Fatalf("op%d offset%d ret%d/%d cursor%d/%d", op, offset, ret, wantRet, consumed, wantConsumed)
	}
	want := bytes.Repeat([]byte{0xa5}, len(scratch))
	copy(want, wantScratch)
	if !bytes.Equal(scratch, want) {
		t.Fatal("unexpected scratch changes")
	}
	return thingSkipResult{op, offset, ret, consumed, sha256.Sum256(scratch)}
}
