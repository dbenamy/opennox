package legacy

func sub_495060(code int32, current, maximum int16) int32 {
	return int32(bool2int(combatAllyAdd(uint32(code), uint16(current), uint16(maximum))))
}

func sub_4950C0(code int32) int32 { return int32(bool2int(combatAllyRemove(uint32(code)))) }

func sub_4950F0(code int32, flag int8) int32 {
	return int32(bool2int(combatAllyFlag(uint32(code), byte(flag))))
}

func nox_xxx_unitSpriteCheckAlly_4951F0(code int32) int32 {
	return int32(bool2int(combatAllyLookup(uint32(code)) != nil))
}
