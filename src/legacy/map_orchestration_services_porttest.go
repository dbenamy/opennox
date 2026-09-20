//go:build porttest

package legacy

// Outer orchestration substitutes only serialization; the default uses real IO.
var orchestrationTestSave func(string, int) int

func init() {
	worldMapSaveService = func(path string, flags int) bool {
		if orchestrationTestSave != nil {
			return orchestrationTestSave(path, flags) != 0
		}
		return worldMapSave(path, flags)
	}
}
