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

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

// Populate only after auditing original C results. Every capture is mandatory.
var gameplayReportHashes = map[string]string{
	"aggregate-coop":         "a7de634b90db6de3e830c0913adf5638520e32d997a23f7a3a99b22daabcbf52",
	"aggregate-float-cache":  "30f56e5ac00f3c38cd9c77af8f00630ef63f68bc3019cdb13312092bb112506c",
	"aggregate-local":        "d97bc87d062062d90bba3fd63f97db075eee8615916cd6c09c07196ec23f8ffa",
	"creature-resources":     "4ef33964e3d4c7f3d23309c00e5dbb30314b8181e0c2ab17ff645a4253b3ca22",
	"earthquake":             "b4f385d01286c9436d8058f8c7b5ecfbe78571ca6fefdf931ba8c1f2900f364e",
	"elimination-frame-wrap": "71915d9a711f188119a48646cbcc7b1aeb0d8d6d54bb78d54e270560d16f2a51",
	"elimination-rules":      "e094031e8ace9546c20ceec6f8693c0c1ad79b9e2037dc27fc5d859da750f64a",
	"health-delta":           "e27f68756fb497f89266eb6c913a3ca3459af364eed62158def73e777ea39ef5",
	"health-fields":          "5e4468131bd9468f3fbb6e1a8cb83cc4240f8cbc24b3acb4c920c0e19e003892",
	"height-fields":          "f4f047bca1fc0bbd5c12ec335be8aac84ce4d5881af49d0bb9f93d66d736d9b7",
	"hidden":                 "39677c28ac8d983bf33396eccf86b3dc457559f9c6a64879cb8808d682038baa",
	"interesting-fanout":     "858fd4ed5f387b5a57eec8527cff997395c15f78b1e832a6eafd900a26a007c9",
	"inventory-class-masks":  "7ff2e495ce78b03609d57ea6ea4b802dbff3a4b44395cefff5604b0a812332a3",
	"inventory-fields":       "78becf60dd91e20a034f86f9e9904fa4d609ba6e07d59cbe8669fa03d2e66f4d",
	"journal-boundaries":     "9260d04c61be33bb02c2bf5237bbaaf944cac05c9a3a4733696593ab5de4e33c",
	"non-players":            "eb7549432165153c9e52350830250ea78edaa2a14014cf1642b2f1c577747047",
	"notifications":          "de5ce71f6e6cb922edb30fd83e3f8ca1f4dfd718df5bc39a4ab47930019c8d81",
	"npc":                    "45a71ea98f5f0de594c0de76828998094a09937d1566a65a289facb3f16da184",
	"object-fields":          "4eda9691ce8fbaac7d1a32dc18098c626061b5c292abea6ced08e8343b5db3eb",
	"player-fields":          "cbc6f037e28e7028ca799b9cde1bbda942436391b7935842bf4fbd6ca056e2af",
	"rate":                   "6f58062569b6fb2645bf3c9c3f2e0c5231e5145e7762f3be2103fd4606d71eb2",
	"scalar-fields":          "1c1bd2819428de0849b808f19dbba0d23b55c3d85ef0521236bac101d9bdefb7",
	"score-changes":          "7af4d365abdb5adc14168f5544c2c3e31a5dbbfa6dfddb0ca78f8c64689d364b",
	"stats-scavenger":        "fc9d6060ca436b8875c881bd0a53bab4540a72430a21b8b74c9bdc5af5815bb4",
	"team-base-cache":        "9df4a41bcede82346e3b3e77a8f499dffbd0f5dc6c5baea2c94dc967eed0211f",
	"team-health":            "69f6974ee51efec03b5df34d43c17dfccbdc606461101845ef75fa30a0744d4d",
	"winners":                "6f76af90723b3a0744eeb04973106d7a021329a8b78cd0d3accce5b72bc5b3e6",
	"zero-max-health":        "58847363506593d5ca775fd85f5d1c3307160a01bd65e235a07df19cb2aa044a",
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
	deadline := *memmap.PtrUint32(0x5D4594, 3468)
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
					binary.LittleEndian.PutUint32(data[5:], deadline-10000)
					binary.LittleEndian.PutUint32(data[9:], s.Owner.Frame)
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
