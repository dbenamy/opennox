//go:build porttest

package legacy

import "unsafe"

// Static tokens stand in for retired callback addresses in snapshot normalization only.
// They are never passed to a C dispatcher or invoked.
var clientUICallbackIdentitySlots [28]byte

const (
	clientUICallbackID_nox_client_toggleSpellbook_45AC70    = 0
	clientUICallbackID_nox_client_trapSetSelect_4604B0      = 1
	clientUICallbackID_nox_xxx_abilityReward_45D290         = 2
	clientUICallbackID_nox_xxx_bookFillAll_45D570           = 3
	clientUICallbackID_nox_xxx_bookHideMB_45ACA0            = 4
	clientUICallbackID_nox_xxx_bookSetForward_45D200        = 5
	clientUICallbackID_nox_xxx_book_45DBE0                  = 6
	clientUICallbackID_nox_xxx_buttonsGetSelectedRow_45E180 = 7
	clientUICallbackID_nox_xxx_clientUpdateButtonRow_45E110 = 8
	clientUICallbackID_nox_xxx_quickBarClose_4606B0         = 9
	clientUICallbackID_sub_45CFC0                           = 10
	clientUICallbackID_sub_45D500                           = 11
	clientUICallbackID_sub_45D870                           = 12
	clientUICallbackID_sub_45D9B0                           = 13
	clientUICallbackID_sub_4602F0                           = 14
	clientUICallbackID_sub_4604E0                           = 15
	clientUICallbackID_sub_460660                           = 16
	clientUICallbackID_sub_460940                           = 17
	clientUICallbackID_sub_460EA0                           = 18
	clientUICallbackID_sub_462740                           = 19
	clientUICallbackID_sub_467650                           = 20
	clientUICallbackID_sub_467BB0                           = 21
	clientUICallbackID_sub_467C10                           = 22
	clientUICallbackID_sub_467C80                           = 23
	clientUICallbackID_sub_478040                           = 24
	clientUICallbackID_sub_479590                           = 25
	clientUICallbackID_sub_4795A0                           = 26
	clientUICallbackID_sub_4C05F0                           = 27
)

func clientUICallbackKey(id int) unsafe.Pointer {
	return unsafe.Pointer(&clientUICallbackIdentitySlots[id])
}
