//go:build !porttest

package legacy

import "unsafe"

func mapRoomBeforeRelease(unsafe.Pointer) {}
