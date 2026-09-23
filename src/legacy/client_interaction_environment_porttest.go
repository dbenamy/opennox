//go:build porttest

package legacy

import "unsafe"

// PortTestClientInteractionWords snapshots the live owners, separately from mapped blobs.
func PortTestClientInteractionWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_5d4594_1064856":                    (*uint32)(unsafe.Pointer(&interactionChatRootWord)),
		"dword_5d4594_1064860":                    (*uint32)(unsafe.Pointer(&interactionChatEditWord)),
		"dword_5d4594_1064864":                    (*uint32)(unsafe.Pointer(&interactionChatDataWord)),
		"dword_5d4594_1064868":                    (*uint32)(unsafe.Pointer(&interactionChatActive)),
		"dword_5d4594_1096636":                    (*uint32)(unsafe.Pointer(&interactionHoverDepth)),
		"dword_5d4594_1096640":                    (*uint32)(unsafe.Pointer(&interactionHoverDrawable)),
		"dword_5d4594_1123520":                    (*uint32)(unsafe.Pointer(&interactionConversationActive)),
		"dword_5d4594_1123524":                    (*uint32)(unsafe.Pointer(&interactionConversationRoot)),
		"dword_5d4594_1193712":                    (*uint32)(unsafe.Pointer(&interactionObserverIcon)),
		"dword_5d4594_1303452":                    (*uint32)(unsafe.Pointer(&interactionGameOverRootWord)),
		"dword_5d4594_1305680":                    (*uint32)(unsafe.Pointer(&interactionHelpRootWord)),
		"dword_5d4594_1319056":                    (*uint32)(unsafe.Pointer(&interactionKeyState)),
		"dword_5d4594_1319060":                    (*uint32)(unsafe.Pointer(&interactionKeyRootWord)),
		"dword_5d4594_1321216":                    (*uint32)(unsafe.Pointer(&interactionVoteIcon)),
		"dword_5d4594_805820":                     (*uint32)(unsafe.Pointer(&interactionPrimaryKey)),
		"dword_5d4594_811904":                     (*uint32)(unsafe.Pointer(&interactionDrawToggle)),
		"dword_5d4594_825736":                     (*uint32)(unsafe.Pointer(&interactionMessageHead)),
		"dword_5d4594_825744":                     (*uint32)(unsafe.Pointer(&interactionChatIcon)),
		"dword_8531A0_2576":                       (*uint32)(unsafe.Pointer(&dword_8531A0_2576)),
		"nox_client_spriteUnderCursorXxx_1096644": (*uint32)(unsafe.Pointer(&interactionUsableCursorDrawable)),
		"nox_xxx_useAudio_587000_80772":           (*uint32)(unsafe.Pointer(&interactionCursorMode)),
	}
	old := make(map[string]uint32)
	for n, p := range words {
		old[n] = *p
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

// PortTestClientInteractionGlyphCache binds the existing eligibility cache to the
// fixture's actual Things table and restores the preceding owner's value.
func PortTestClientInteractionGlyphCache(id uint32) func() {
	old := glyphClientType
	glyphClientType = id
	return func() { glyphClientType = old }
}
