//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestReliableReportsResets(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []struct {
		Name   string
		Return uint32
		State  legacy.PortTestReliableReportState
	}
	defer func() {
		spellbookCapture(t, "reliable-reports-resets", rows, "f0c64ff12ec99a556527231149c44c55af96d093a9b716f8dc518fc334d1daf1")
	}()
	for _, op := range []int{0, 1, 2, 3, 4, 5} {
		for _, to := range []int{0, 1, 7, 31} {
			for _, mode := range []uint32{0, 1, 2} {
				for _, rate := range []uint32{1, 3, 30, 255} {
					for _, fps := range []uint32{1, 30, 60, 255} {
						name := fmt.Sprintf("op%d/to%d/mode%d/rate%d/fps%d", op, to, mode, rate, fps)
						t.Run(name, func(t *testing.T) {
							o.reset()
							*o.words["rateMode"] = mode
							*o.words["rate"] = rate
							o.s.SetTickRate(fps)
							rates := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565124)), 384)
							seq := (*[32]uint16)(memmap.PtrOff(0x5D4594, 1565524))
							for i := 0; i < 32; i++ {
								seq[i] = uint16(65535 - i)
								p := rates[i*12:]
								p[0] = byte(i%2 + 1)
								p[1] = byte(i%6 + 1)
								p[2] = 99
								p[3] = 0xa5
								binary.LittleEndian.PutUint32(p[4:], 0xdeadface)
								binary.LittleEndian.PutUint32(p[8:], 0x12345678)
							}
							if legacy.PortTestReliableReports(8, 31, 0, 0, []byte{77}, nil, 0, false) != 1 {
								t.Fatal("enqueue")
							}
							before := o.state()
							rv := legacy.PortTestReliableReports(op, to, 0, 0, nil, nil, 0, false)
							after := o.state()
							divisor := uint32(1)
							if mode == 1 {
								divisor = rate
							}
							wantReturn := uint32(0)
							switch op {
							case 0, 3:
								wantReturn = 2 * (fps / divisor)
							case 1:
								p := before.Rates[12*to:]
								wantReturn = uint32(p[0]) * uint32(p[1]) * (fps / divisor)
							case 4:
								wantReturn = uint32(to)
							case 5:
								wantReturn = 2 * (fps / divisor)
							}
							if rv != wantReturn {
								t.Fatalf("return%d want%d", rv, wantReturn)
							}
							for i, v := range after.Sequence {
								want := before.Sequence[i]
								if op == 0 || op == 2 || op == 4 && i == to {
									want = 0
								}
								if v != want {
									t.Fatalf("sequence%d=%d want%d", i, v, want)
								}
							}
							for i := 0; i < 32; i++ {
								p := after.Rates[i*12:][:12]
								old := before.Rates[i*12:][:12]
								want := [12]byte{}
								copy(want[:], old)
								if op == 0 || op == 3 {
									clear(want[:])
								}
								if op == 0 || op == 3 || op == 5 && i == to {
									want[0] = 1
									want[1] = 2
									want[2] = byte(rate)
								}
								if op == 0 || op == 3 || (op == 1 || op == 5) && i == to {
									hi := uint32(want[0]) * uint32(want[1]) * (fps / divisor)
									lo := uint32(0)
									if want[1] > 2 {
										lo = uint32(want[0]) * uint32(want[1]-1) * (fps / divisor)
									}
									binary.LittleEndian.PutUint32(want[4:], hi)
									binary.LittleEndian.PutUint32(want[8:], lo)
								}
								if *(*[12]byte)(p) != want {
									t.Fatalf("rate%d=%v want%v", i, p, want)
								}
							}
							if op == 0 {
								if after.Pool || len(after.Nodes) != 0 {
									t.Fatal("init retained queue")
								}
							} else if !after.Pool || len(after.Nodes) != 1 {
								t.Fatal("reset damaged queue")
							}
							rows = append(rows, struct {
								Name   string
								Return uint32
								State  legacy.PortTestReliableReportState
							}{name, rv, after})
						})
					}
				}
			}
		}
	}
}
