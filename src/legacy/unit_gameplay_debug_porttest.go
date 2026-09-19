//go:build porttest

package legacy

func PortTestUnitDebug(kind int, frame, code uint32, name string) { unitDebug(kind, frame, code, name) }
