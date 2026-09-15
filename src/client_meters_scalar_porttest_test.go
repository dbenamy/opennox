//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientMetersSetterContract(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	if got := legacy.PortTestMeterCall(4, nil, 73, 101, 0, 0); got != 73 {
		t.Fatal("health setter return")
	}
	if r := o.meters.Records[0]; r.Current != 73 || r.Maximum != 101 {
		t.Fatal("health meter values")
	}
	if got := legacy.PortTestMeterCall(8, nil, 29, 83, 0, 0); got != 29 {
		t.Fatal("mana setter return")
	}
	if r := o.meters.Records[1]; r.Current != 29 || r.Maximum != 83 {
		t.Fatal("mana meter values")
	}
	if *o.meters.NamedWord("dword_5d4594_1096260") != 32 {
		t.Fatal("meter animation duration")
	}
	legacy.PortTestMeterCall(10, nil, 20, 90, 0, 0)
	if inputKeyTimeoutsOld[17] != 120 {
		t.Fatal("charge change did not use actual input cooldown owner")
	}
}
func TestClientMetersScalarMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	values := []int{-2147483648, -1, 0, 1, 99, 255, 65535, 2147483647}
	id := 0
	for _, a := range values {
		for _, b := range values {
			o.plain(t)
			for step, op := range []int{4, 5, 6, 7, 8, 9, 13, 10, 23, 12, 11, 0, 1, 3, 26, 2} {
				out = append(out, o.invoke(t, id, step, op, nil, a, b, 0, 0))
			}
			id++
		}
	}
	meterCapture(t, "meters-scalars", out, len(out), "1e072b1122721961dce51fef95d5b553b58dfbbc58998b98074151756e2c384d")
}

func TestClientMetersSoundContract(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	o.meters.Records[0].Current = 10
	o.meters.Records[0].Maximum = 100
	*o.meters.NamedWord("nox_player_netCode_85319C") = 7
	if got := legacy.PortTestMeterCall(14, nil, 0, 0, 0, 0); got != 4 {
		t.Fatal("heartbeat did not set cooldown")
	}
	if len(o.sounds) != 1 || o.sounds[0] != [2]int{896, 82} {
		t.Fatalf("heartbeat requests %v", o.sounds)
	}
	legacy.PortTestMeterCall(14, nil, 0, 0, 0, 0)
	if len(o.sounds) != 1 {
		t.Fatal("heartbeat ignored cooldown")
	}
	legacy.PortTestMeterCall(25, o.meters.Records[0].Window, 7, 0, 0, 0)
	if len(o.sounds) != 2 || o.sounds[1] != [2]int{901, 100} {
		t.Fatalf("meter toggle requests %v", o.sounds)
	}
}
