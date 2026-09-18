//go:build porttest

package blobdata

func PortTestPrefabScriptsTables() map[uintptr][]byte {
	return map[uintptr][]byte{
		197556: append([]byte(nil), data587000[197556:197593]...),
		282608: append([]byte(nil), data587000[282608:282620]...),
	}
}
