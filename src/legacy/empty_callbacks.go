package legacy

import "C"

// These no-op callbacks remain distinct C-callable identities. Legacy modifier
// dispatch passes arguments that the callbacks intentionally ignore. Registration
// and effect lookups compare their addresses.

//export nullsub_22
func nullsub_22() {}

//export nullsub_36
func nullsub_36() {}

//export nullsub_38
func nullsub_38() {}

//export nullsub_39
func nullsub_39() {}

//export nullsub_40
func nullsub_40() {}

//export nullsub_41
func nullsub_41() {}

//export nullsub_42
func nullsub_42() {}

//export nullsub_43
func nullsub_43() {}

//export nullsub_44
func nullsub_44() {}
