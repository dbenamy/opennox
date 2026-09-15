//go:build porttest

package noxrender

func (r *NoxRender) PortTestUITextSmoothing() bool { return r.text.smooth }
