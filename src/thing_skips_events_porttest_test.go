//go:build porttest

package opennox

import "testing"

func thingSkipEventPayload(b []byte, tag byte, n int) []byte {
	b = append(b, tag)
	switch tag {
	case 1, 2, 3, 4, 5:
		b = thingSkipFill(b, 1, 0x81)
	case 6, 9, 10:
		b = thingSkipFill(b, 2, 0x82)
	case 7:
		for j := 0; j < 3; j++ {
			b = thingSkipName(b, max(1, n), 0x41)
		}
		b = append(b, 0)
	case 8:
		b = thingSkipFill(b, 8, 0x83)
	}
	return b
}
func TestThingSkipsEvents(t *testing.T) {
	var rows []thingSkipResult
	for _, wrapped := range []bool{false, true} {
		op := 5
		if wrapped {
			op = 4
		}
		for _, n := range []int{0, 1, 31, 255} {
			for tag := 0; tag < 256; tag++ {
				var b []byte
				if wrapped {
					b = thingSkipName(b, n, 0x61)
				}
				b = thingSkipEventPayload(b, byte(tag), n)
				ret := 0
				if tag == 0 {
					ret = 1
				} else if tag <= 10 {
					b = append(b, 0)
					ret = 1
				}
				rows = append(rows, thingSkipInvoke(t, op, tag%8, b, ret, len(b), nil))
			}
		}
	}
	thingSkipCapture(t, "event-tags", rows, "eaf550105afc29ecbf1579f004c32119341eace14399dd9044cd87770bf70b07")
}
func TestThingSkipsEventSequences(t *testing.T) {
	var rows []thingSkipResult
	for _, wrapped := range []bool{false, true} {
		op := 5
		if wrapped {
			op = 4
		}
		for _, n := range []int{0, 1, 255} {
			for _, repeats := range []int{1, 3, 64} {
				for _, reverse := range []bool{false, true} {
					for _, ending := range []byte{0, 11, 128, 255} {
						var b []byte
						if wrapped {
							b = thingSkipName(b, n, 0x61)
						}
						for i := 0; i < repeats; i++ {
							for j := 1; j <= 10; j++ {
								tag := j
								if reverse {
									tag = 11 - j
								}
								b = thingSkipEventPayload(b, byte(tag), n)
							}
						}
						b = append(b, ending)
						ret := 0
						if ending == 0 {
							ret = 1
						}
						rows = append(rows, thingSkipInvoke(t, op, repeats%8, b, ret, len(b), nil))
					}
				}
			}
		}
	}
	thingSkipCapture(t, "event-sequences", rows, "69c6b2f8281f179c14a944cf5916efbb6b301e5acb57141c8c9a8cf9c51b84d0")
}
