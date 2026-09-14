//go:build porttest

package opennox

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameplayReportsJournalPadding(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{53, 54, 55} {
		s := gameplayReportsBase(op)
		p := s.Callbacks.Shop
		sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.Reports
		sp.Records = []legacy.PortTestGameplayReportRecord{{Text: strings.Repeat("L", 63), Flags: 0x1234}, {Text: "x", Flags: 0x5678}}
		sp.Calls = map[uint32][5]legacy.PortTestGameplayReportArg{
			0: {reportValue(7), {Kind: "record", Ref: 1}},
			1: {reportValue(7), {Kind: "record", Ref: 2}},
		}
		p.Sequence = []legacy.PortTestShopAction{{Op: 1800 + op}, {Op: 1800 + op, Value: 1}}
		cases = append(cases, s)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		if !r.Intact || !r.Callbacks.Shop.Intact {
			t.Fatalf("journal guards %d", i)
		}
		found := 0
		for _, p := range r.Callbacks.Shop.Sequence[1].Packets {
			if len(p.Data) != 68 || p.Data[0] != 213 || p.Data[1] != byte(i+1) || p.Data[2] != 'x' {
				continue
			}
			found++
			if p.Recipient != 7 || p.Data[3] != 0 {
				t.Errorf("journal defined fields case%d", i)
			}
			end := 66
			if i == 1 {
				end = 68
			} else if binary.LittleEndian.Uint16(p.Data[66:]) != 0x5678 {
				t.Errorf("journal flags case%d", i)
			}
			for j := 4; j < end; j++ {
				if p.Data[j] != 0 {
					t.Errorf("journal case%d has nonzero unused byte at%d", i, j)
					break
				}
			}
		}
		if found != 1 {
			t.Errorf("journal short-message count case%d: %d", i, found)
		}
	}
	if path := os.Getenv("OPENNOX_REPORTS_JOURNAL_CAPTURE"); path != "" {
		data, err := json.Marshal(out)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
