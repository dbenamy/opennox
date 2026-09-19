//go:build porttest

package blobdata

// PortTestPlayerActionTable returns the shipped weapon-to-animation mapping.
func PortTestPlayerActionTable() []byte { return append([]byte(nil), data587000[215824:215932]...) }

func PortTestPlayerAbilityTable() []byte { return append([]byte(nil), data587000[206108:206148]...) }
func PortTestPlayerWeight() []byte       { return append([]byte(nil), data581450[10216:10224]...) }

func PortTestPlayerCorpsePoints() []byte { return append([]byte(nil), data587000[280376:281168]...) }

// Shipped four-item SetHalberd name storage. Pointers are relocated by the owner.
func PortTestScriptInventoryNames() []byte { return append([]byte(nil), data587000[247352:247416]...) }
