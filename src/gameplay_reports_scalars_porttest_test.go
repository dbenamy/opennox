//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

// Populate only after auditing original C results. Every capture is mandatory.
var gameplayReportHashes = map[string]string{
	"player-fields": "cbc6f037e28e7028ca799b9cde1bbda942436391b7935842bf4fbd6ca056e2af",
	"height-fields": "f4f047bca1fc0bbd5c12ec335be8aac84ce4d5881af49d0bb9f93d66d736d9b7",
	"object-fields": "4eda9691ce8fbaac7d1a32dc18098c626061b5c292abea6ced08e8343b5db3eb",
	"health-fields": "5e4468131bd9468f3fbb6e1a8cb83cc4240f8cbc24b3acb4c920c0e19e003892",
	"scalar-fields": "1c1bd2819428de0849b808f19dbba0d23b55c3d85ef0521236bac101d9bdefb7",
	"health-delta":  "e27f68756fb497f89266eb6c913a3ca3459af364eed62158def73e777ea39ef5",
}

func gameplayReportsCapture(t *testing.T, label string, out []legacy.PortTestRoamResult) {
	t.Helper()
	for i, r := range out {
		if !r.Intact || !r.Callbacks.Intact || !r.Spells.Intact || !r.Combat.Intact || !r.MonsterState.Intact || !r.Callbacks.Shop.Intact {
			t.Fatalf("%s case%d guards", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_GAMEPLAY_REPORTS_CAPTURE"); prefix != "" {
		if err = os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d cases %s", label, len(out), hash)
	want, ok := gameplayReportHashes[label]
	if !ok {
		t.Fatalf("missing audited C capture: %s", label)
	}
	if hash != want {
		t.Fatalf("%s hash %s want %s", label, hash, want)
	}
}
func TestGameplayReportsScalarFields(t *testing.T) {
	values := []uint32{0, 1, 127, 128, 255, 256, 32767, 32768, 65535, 65536, 0x7fffffff, 0x80000000, 0xffffffff}
	var cases []legacy.PortTestRoamSpec
	type expected struct {
		op        int
		recipient byte
		data      []byte
		ordered   byte
		priority  uint32
	}
	var checks []expected
	for _, op := range []int{3, 40, 44, 56, 57, 58, 59, 62, 64, 65, 67, 72} {
		for _, recipient := range []uint32{1, 7, 31} {
			for i, v := range values {
				w := values[len(values)-1-i]
				x := v ^ 0xa5
				y := v*17 + 3
				s := gameplayReportsBase(op, reportValue(recipient), reportValue(v), reportValue(w), reportValue(x), reportValue(y))
				var data []byte
				ordered, priority := byte(1), uint32(1)
				switch op {
				case 3:
					data = []byte{237, byte(v)}
				case 40:
					data = make([]byte, 13)
					data[0] = 211
					binary.LittleEndian.PutUint32(data[1:], v)
				case 44:
					data = []byte{151, byte(v)}
				case 56:
					flag := byte(0)
					if w == 1 {
						flag = 1
					}
					data = []byte{214, byte(v), flag}
					ordered = 0
				case 57:
					data = []byte{216, byte(w), byte(v), byte(x), byte(y), byte(y >> 8)}
				case 58:
					data = []byte{217, byte(v), byte(w), byte(w >> 8)}
				case 59:
					data = make([]byte, 6)
					data[0] = 223
					binary.LittleEndian.PutUint32(data[1:], v)
					data[5] = byte(w)
					ordered = 0
				case 62:
					data = []byte{113}
					priority = 0
				case 64:
					data = []byte{54}
				case 65:
					a, b := byte(0), byte(0)
					if v != 0 {
						a = 1
					}
					if w != 0 {
						b = 1
					}
					data = []byte{228, a, b}
				case 67:
					data = []byte{240, 0}
				case 72:
					data = []byte{240, 20}
					ordered = 0
				}
				cases = append(cases, s)
				checks = append(checks, expected{op, byte(recipient), data, ordered, priority})
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		step := r.Callbacks.Shop.Sequence[0]
		c := checks[i]
		if c.op == 44 {
			index := map[byte]int{1: 0, 7: 1, 31: 2}[c.recipient]
			if step.GameplayReports == nil || !bytes.Equal(step.ResourceMessages[index], c.data) || len(step.Packets) != 0 {
				t.Fatalf("direct-list scalar case%d", i)
			}
			continue
		}
		if len(step.Packets) != 1 {
			t.Fatalf("scalar case%d packet count%d", i, len(step.Packets))
		}
		p := step.Packets[0]
		if p.Recipient != c.recipient || p.Ordered != c.ordered || p.A4 != 0 || p.A5 != c.priority {
			t.Fatalf("scalar case%d route/metadata", i)
		}
		got, want := p.Data, c.data
		if c.op == 40 {
			if len(got) != 13 {
				t.Fatalf("timer size case%d", i)
			}
			got = got[:5]
			want = want[:5]
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("scalar case%d defined bytes", i)
		}
	}
	gameplayReportsCapture(t, "scalar-fields", out)
}
func TestGameplayReportsHealthDelta(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type expected struct {
		recipient byte
		id        uint16
		delta     int16
	}
	var checks []expected
	for _, recipient := range []uint32{1, 7, 31} {
		for _, id := range []uint32{0, 1, 32767, 32768, 65535} {
			for _, delta := range []int32{-32768, -257, -1, 0, 1, 32767} {
				cases = append(cases, gameplayReportsBase(19, reportValue(recipient), reportValue(id), reportValue(uint32(delta))))
				checks = append(checks, expected{byte(recipient), uint16(id), int16(delta)})
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		step := r.Callbacks.Shop.Sequence[0]
		c := checks[i]
		if c.delta >= 0 {
			if len(step.Packets) != 0 || step.Return != uint32(c.delta) {
				t.Fatalf("nonnegative health delta case%d", i)
			}
			continue
		}
		want := []byte{66, byte(c.id), byte(c.id >> 8), byte(c.delta), byte(uint16(c.delta) >> 8)}
		if len(step.Packets) != 1 {
			t.Fatalf("negative health delta count case%d", i)
		}
		p := step.Packets[0]
		if p.Recipient != c.recipient || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 || !bytes.Equal(p.Data, want) {
			t.Fatalf("negative health delta case%d", i)
		}
	}
	gameplayReportsCapture(t, "health-delta", out)
}
