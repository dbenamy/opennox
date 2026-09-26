//go:build porttest

package gui

// PortTestOwnAnimations isolates animation ticking from unrelated fixture owners.
// Restore frees any remaining owned nodes before restoring the previous list.
func PortTestOwnAnimations() func() {
	oldList, oldState, oldSpeed := animList, animGlobalState, AnimSpeed
	animList, animGlobalState, AnimSpeed = nil, AnimInDone, 1
	return func() {
		for animList != nil {
			animList.Free()
		}
		animList, animGlobalState, AnimSpeed = oldList, oldState, oldSpeed
	}
}
