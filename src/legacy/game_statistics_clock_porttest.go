//go:build porttest

package legacy

func PortTestStatisticsClock() (set func(uint32), restore func()) {
	old := statisticsTime
	return func(value uint32) { statisticsTime = func() uint32 { return value } }, func() { statisticsTime = old }
}
