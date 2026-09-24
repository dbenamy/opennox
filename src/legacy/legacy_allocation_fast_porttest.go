//go:build porttest && !safe

package legacy

func PortTestAllocationUsesTracker() bool { return false }
