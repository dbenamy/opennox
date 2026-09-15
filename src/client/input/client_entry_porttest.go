//go:build porttest

package input

import "github.com/opennox/libs/client/keybind"

func (h *Handler) PortTestEntryState() []uint32 {
	caps := uint32(0)
	if h.capsState {
		caps = 1
	}
	out := []uint32{caps, uint32(h.modKey), uint32(h.inputSeq), uint32(h.nox_input_seq), uint32(h.nox_input_seq_prev)}
	for key := keybind.Key(0); key < 256; key++ {
		st := h.k.nox_input_map_byKey[key]
		flag := uint32(0)
		if st.field2 {
			flag = 1
		}
		out = append(out, uint32(st.state), flag, uint32(st.seq))
	}
	return out
}
