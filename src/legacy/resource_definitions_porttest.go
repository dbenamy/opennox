//go:build porttest

package legacy

import "unsafe"

func PortTestResourceParser(name string, input, data unsafe.Pointer) int {
	return resourceParser(name, (*byte)(input), data)
}
