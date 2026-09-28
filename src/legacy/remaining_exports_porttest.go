//go:build porttest

package legacy

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

// These Go-only fixture adapters preserve historical argument widths and capture
// entrypoints. Production uses the native owners directly. No C exports or C
// algorithm copies are kept here; explicit-width Go types retain the fixture ABI.

func nox_client_copyRect_49F6F0(x, y, w, h int32) int32 {
	return int32(uiRenderCopyRect(int(x), int(y), int(w), int(h)))
}

func nox_drawable_next_45A070(dr *nox_drawable) *nox_drawable {
	if dr == nil {
		return nil
	}
	return (*nox_drawable)(asDrawable(dr).NextPtr.C())
}

func nox_getHostPlayerUnit() *nox_object_t {
	return asObjectC(GetServer().S().Players.HostUnit())
}

func nox_server_questMaybeWarp_4E8F60() bool { return bool(worldQuestMaybeWarp()) }

func nox_xxx_clientAskInfoMb_4BF050(dr *nox_drawable) *uint16 {
	return (*uint16)(unsafe.Pointer(uiItemTooltip((*client.Drawable)(unsafe.Pointer(dr)))))
}

func nox_xxx_client_57B400(ptr int32) int32 {
	return int32(glyphSelectionAllowed((*client.Drawable)(unsafe.Pointer(uintptr(uint32(ptr))))))
}

func nox_xxx_collideReflect_57B810(normal *float32, velocity int32) int32 {
	collisionReflect((*types.Pointf)(unsafe.Pointer(normal)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(velocity)))))
	return velocity
}

func nox_xxx_mapGenEdge_543EB0(index, edge int32) int32 {
	return int32(generateBorderEdge(int32(index), int32(edge)))
}

func nox_xxx_map_57B850(pos *[2]float32, shape *float32, point *[2]float32) int32 {
	return int32(bool2int(collisionContains((*types.Pointf)(unsafe.Pointer(pos)), (*[11]float32)(unsafe.Pointer(shape)), (*types.Pointf)(unsafe.Pointer(point)))))
}

func nox_xxx_mathPointOnTheLine_57C8A0(line *[4]float32, point, out *[2]float32) int32 {
	return int32(bool2int(projectLine((*[4]float32)(unsafe.Pointer(line)), (*types.Pointf)(unsafe.Pointer(point)), (*types.Pointf)(unsafe.Pointer(out)))))
}

func nox_xxx_monsterGetSoundSet_424300(u *nox_object_t) unsafe.Pointer {
	return resourceMonsterSound(asObjectS(u))
}

func nox_xxx_netGetUnitByExtent_4ED020(value int32) int32 {
	obj := GetServer().S().Objs.GetObjectByInd(int(uint32(value)))
	// Preserve the raw 32-bit address representation used by fixture callers.
	return int32(uintptr(unsafe.Pointer(obj)))
}

func nox_xxx_netSavePlayer_41CE00() int32 { return int32(playerFileSaveRequest()) }

func nox_xxx_netSendPacket1_4E5390(to, data, size, related, priority int32) int32 {
	return int32(reliableEnqueue(int(to), unsafe.Slice((*byte)(unsafe.Pointer(uintptr(uint32(data)))), int(size)), (*server.Object)(unsafe.Pointer(uintptr(uint32(related)))), int(priority), 1))
}

func nox_xxx_netSendPointFx_522FF0(code int8, pos *[2]float32) int32 {
	return int32(visibilityFXPoint(byte(code), *(*types.Pointf)(unsafe.Pointer(pos))))
}

func nox_xxx_playerCheckSpellClass_57AEA0(class, ind int32) int32 {
	return int32(playerSpellClassCheck(int32(class), int32(ind)))
}

func nox_xxx_playerResetImportantCtr_4E4F40(to int32) int32 { return int32(reliableResetRate(int(to))) }

func nox_xxx_protectPlayerHPMana_56F870(id int32, value uint16) uint32 {
	return uint32(setProtectionRecord(int32(id), uint32(value)))
}

func nox_xxx_protectionStringCRCLen_56FAE0(data *int32, size uint32) int32 {
	return int32(protectionStringChecksum(unsafe.Pointer(data), uint32(size)))
}

func nox_xxx_sendArrowTrapFX_5238A0(pos *float32, extra int8) {
	visibilityFXArrowTrap(*(*types.Pointf)(unsafe.Pointer(pos)), byte(extra))
}

func nox_xxx_serverHandleClientConsole_443E90(pl *nox_playerInfo, action int8, text *uint16) int32 {
	return int32(consoleCommandRemote(asPlayerS(pl), byte(action), GoWStringP(unsafe.Pointer(text))))
}

func nox_xxx_tileAllocTileInCoordList_5040A0(a0 int32, a1 int32, a2 float32) *uint32 {
	return (*uint32)(mapRoomPointer(prefabTileNew(int32(uint32(a0)), int32(uint32(a1)), math.Float32bits(float32(a2)))))
}

func nox_xxx_tileCheckByte3_544070(index int32) int32 {
	return int32(bool2int(selectBorderPrimary(int32(index))))
}

func nox_xxx_tileCheckByte4_5440A0(variation int32) int32 {
	return int32(bool2int(selectBorderVariation(int32(variation))))
}

func nox_xxx_tileCheckImageVari_51D570(variation int32) int32 {
	return int32(bool2int(selectTileVariation(int32(variation))))
}

func nox_xxx_tileCheckImage_51D540(index int32) int32 {
	return int32(bool2int(selectTileImage(int32(index))))
}

func nox_xxx_tile_51D5C0(value int32) int32 {
	return int32(bool2int(setTileFlag(int32(value))))
}

func nox_xxx_waypointNext_579870(a1 int32) int32 {
	if a1 == 0 {
		return 0
	}
	return int32(waypointRaw(waypointFromRaw(int32(a1)).WpNext))
}

func nox_xxx_wndDraw_49F7F0() { objectRenderSaveClip() }

func sub_411490(index, edge int32) int32 {
	return int32(normalizeBorderEdge(int32(index), int32(edge)))
}

func sub_436550() int32 { return int32(interactionFrameGate()) }

func sub_43AF30() int32 { return int32(browserUI.hosting) }

func sub_45A010(dr *nox_drawable) *nox_drawable { return (*nox_drawable)(asDrawable(dr).Field_104.C()) }

func sub_49F860() int32 { return int32(objectRenderRestoreClip()) }

func sub_4C5050() { objectRenderBeamReset() }

func sub_4E4F30(to int32) int32 { return int32(ReliableResetSequence(ntype.PlayerInd(to))) }

func sub_4E55F0(to uint8) int32 { return int32(reliableRemoveRecipient(byte(to))) }

func sub_4E8E50() *uint8 { return (*uint8)(unsafe.Pointer(worldQuestPending())) }

func sub_4E8E60() int32 { return int32(worldQuestCountdown()) }

func sub_4E9010() int32 { return int32(bool2int(worldQuestExitReady())) }

func sub_504290(a0 int8, a1 int8) *uint32 {
	return (*uint32)(mapRoomPointer(prefabWallNew(byte(uint32(a0)), byte(uint32(a1)))))
}

func sub_51D2C0(source, target int32) int32 {
	return int32(bool2int(appendWaypointLink(waypointFromRaw(int32(source)), waypointFromRaw(int32(target)), int8(memmap.Uint8(0x973F18, 35972)))))
}

func sub_51D300(source, target int32, kind int8) int32 {
	return int32(bool2int(appendWaypointLink(waypointFromRaw(int32(source)), waypointFromRaw(int32(target)), int8(kind))))
}

func sub_51DD50(x, y, flags, key int32) {
	pushTileFill(int32(x), int32(y), int32(flags), int32(key))
}

func sub_51DE30(x, y, flags *uint32) int32 {
	return int32(bool2int(popTileFill((*uint32)(unsafe.Pointer(x)), (*uint32)(unsafe.Pointer(y)), (*uint32)(unsafe.Pointer(flags)))))
}

func sub_543E60(record, category int32) int32 {
	p := (*[4]uint32)(unsafe.Pointer(uintptr(uint32(record))))
	return int32(bool2int(mergeBorderEdge(p, int32(category))))
}

func sub_544020(name *int8) int32 {
	return int32(bool2int(selectBorderName((*byte)(unsafe.Pointer(name)))))
}

func sub_56F780(id, value int32) uint32 {
	return uint32(setProtectionRecord(int32(id), uint32(value)))
}

func sub_56F820(id int32, value uint8) uint32 {
	return uint32(setProtectionRecord(int32(id), uint32(value)))
}

func sub_5798A0(a1 int32) int32 {
	if a1 == 0 {
		return 0
	}
	return int32(waypointRaw(waypointFromRaw(int32(a1)).WpNext))
}

func sub_579E70() *uint32 {
	wp, _ := alloc.New(server.Waypoint{})
	wp.Flags |= 0x1000000
	return (*uint32)(unsafe.Pointer(wp))
}

func sub_57C790(line *[4]float32, point, out *[2]float32, length float32) {
	projectLineClamped((*[4]float32)(unsafe.Pointer(line)), (*types.Pointf)(unsafe.Pointer(point)), (*types.Pointf)(unsafe.Pointer(out)), float32(length))
}
