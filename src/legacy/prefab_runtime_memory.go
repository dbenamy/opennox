package legacy

import "unsafe"

// MapPrefabFreeNode releases a prefab cache wrapper with the same allocator
// used by its C constructors and removal paths. Group references have their
// own tracked Go allocator and must not use this function.
func MapPrefabFreeNode(p unsafe.Pointer) { mapRoomRelease(p) }
