//go:build porttest

package legacy

// The C route was qualified before its unused production exports were retired.
// Keep the independent numeric contracts against the existing Go getters.
type PortTestBalanceGetterPair struct {
	ScalarGo, IndexGo float64
}

func PortTestBalanceGetters(key string, index int) PortTestBalanceGetterPair {
	// internCStr matches production ownership for these immediate calls.
	k := internCStr(key)
	return PortTestBalanceGetterPair{
		ScalarGo: float64(nox_xxx_gamedataGetFloat_419D40(k)),
		IndexGo:  float64(nox_xxx_gamedataGetFloatTable_419D70(k, int(int32(index)))),
	}
}
