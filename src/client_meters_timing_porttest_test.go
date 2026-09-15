//go:build porttest

package opennox

import (
	"testing"
)

func TestClientMetersHeartbeatMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	id := 0
	for _, fps := range []uint32{1, 30, 60, 100, 144} {
		o.c.srv.SetTickRate(fps)
		for _, health := range []uint32{0, 1, 2, 10, 39, 40, 41, 100} {
			for _, frame := range []uint32{0, 1, 19, 20, 21, 30, 60, 120, 0xffffffff} {
				o.plain(t)
				o.c.srv.SetFrame(frame)
				*o.meters.NamedWord("nox_player_netCode_85319C") = 7
				o.meters.Records[0].Current, o.meters.Records[0].Maximum = health, 100
				for step := 0; step < 2; step++ {
					out = append(out, o.invoke(t, id, step, 14, nil, 0, 0, 0, 0))
				}
				id++
			}
		}
	}
	meterCapture(t, "meters-heartbeat", out, len(out), "16a58b33bb40f081745ec6d1853a87611bd21b74ecdaf39a772c2935dfa7309f")
}

func TestClientMetersInputMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	id := 0
	for _, op := range []int{15, 21, 25} {
		for _, event := range []int{0, 1, 5, 7, 8, 12, 16, 17, 16384} {
			o.plain(t)
			for step := 0; step < 2; step++ {
				out = append(out, o.invoke(t, id, step, op, o.meters.Records[0].Window, event, 0, 0, 0))
			}
			id++
		}
	}
	meterCapture(t, "meters-input", out, len(out), "8d23d000cc3f1ad99844cbb9c13d9f8f272ede0248ee814dcf1fc4301fa36cc1")
}
