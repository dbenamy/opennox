//go:build porttest

package legacy

// Fixtures may supply decoded caches; all other calls use the actual loader.
var populationTestLoad func(int32) (int32, bool)

func populationLoadPrefab(index int32) uint32 {
	if populationTestLoad != nil {
		if result, ok := populationTestLoad(index); ok {
			return uint32(result)
		}
	}
	return prefabSelect(index)
}
