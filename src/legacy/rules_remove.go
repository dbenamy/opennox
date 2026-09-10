package legacy

import "C"

import "github.com/opennox/libs/ifs"

//export sub_57A9F0
func sub_57A9F0(mapName, fileName *C.char) C.int {
	return C.int(bool2int(ifs.Remove("maps\\"+GoString(mapName)+"\\"+GoString(fileName)) == nil))
}
