//go:build porttest

package legacy

/*
#include "defs.h"
_Static_assert(sizeof(int) == 4, "protection ABI requires 32-bit int");
_Static_assert(sizeof(unsigned int) == 4, "protection ABI requires 32-bit unsigned int");
*/
import "C"
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
// algorithm copies are kept here; remaining C types belong to fixture bindings.

func nox_client_copyRect_49F6F0(x, y, w, h C.int) C.int {
	return C.int(uiRenderCopyRect(int(x), int(y), int(w), int(h)))
}

func nox_client_newScreenParticle_431540(kind, x, y, vx, vy, gravity C.int, size, timer, phase, mode C.char) *C.nox_screenParticle {
	return (*C.nox_screenParticle)(unsafe.Pointer(screenParticleCreate(int(kind), int(x), int(y), int(vx), int(vy), int(gravity), byte(size), byte(timer), byte(phase), byte(mode))))
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

func nox_server_questMaybeWarp_4E8F60() C.bool { return C.bool(worldQuestMaybeWarp()) }

func nox_xxx_clientAskInfoMb_4BF050(dr *nox_drawable) *C.wchar2_t {
	return (*C.wchar2_t)(unsafe.Pointer(uiItemTooltip((*client.Drawable)(unsafe.Pointer(dr)))))
}

func nox_xxx_client_57B400(ptr C.int) C.int {
	return C.int(glyphSelectionAllowed((*client.Drawable)(unsafe.Pointer(uintptr(uint32(ptr))))))
}

func nox_xxx_collideReflect_57B810(normal *C.float, velocity C.int) C.int {
	collisionReflect((*types.Pointf)(unsafe.Pointer(normal)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(velocity)))))
	return velocity
}

func nox_xxx_createSpark_54FD80(a1 C.float, a2 C.float, a3 C.int, a4 C.int, a5 C.float, a6 C.float, a7 C.float, a8 C.int) *C.float {
	return (*C.float)(temporarySpark(types.Ptf(float32(a1), float32(a2)), types.Ptf(float32(a5), float32(a6)), int32(a3), int32(a4), float32(a7), objectFromInt(int32(a8))).CObj())
}

func nox_xxx_mapGenEdge_543EB0(index, edge C.int) C.int {
	return C.int(generateBorderEdge(int32(index), int32(edge)))
}

func nox_xxx_map_57B850(pos *C.float2, shape *C.float, point *C.float2) C.int {
	return C.int(bool2int(collisionContains((*types.Pointf)(unsafe.Pointer(pos)), (*[11]float32)(unsafe.Pointer(shape)), (*types.Pointf)(unsafe.Pointer(point)))))
}

func nox_xxx_mathPointOnTheLine_57C8A0(line *C.float4, point, out *C.float2) C.int {
	return C.int(bool2int(projectLine((*[4]float32)(unsafe.Pointer(line)), (*types.Pointf)(unsafe.Pointer(point)), (*types.Pointf)(unsafe.Pointer(out)))))
}

func nox_xxx_monsterGetSoundSet_424300(u *nox_object_t) unsafe.Pointer {
	return resourceMonsterSound(asObjectS(u))
}

func nox_xxx_netGetUnitByExtent_4ED020(value C.int) C.int {
	obj := GetServer().S().Objs.GetObjectByInd(int(uint32(value)))
	// Preserve the raw 32-bit address ABI used by remaining decompiled C callers.
	return C.int(uintptr(unsafe.Pointer(obj)))
}

func nox_xxx_netSavePlayer_41CE00() C.int { return C.int(playerFileSaveRequest()) }

func nox_xxx_netSendPacket1_4E5390(to, data, size, related, priority C.int) C.int {
	return C.int(reliableEnqueue(int(to), unsafe.Slice((*byte)(unsafe.Pointer(uintptr(uint32(data)))), int(size)), (*server.Object)(unsafe.Pointer(uintptr(uint32(related)))), int(priority), 1))
}

func nox_xxx_netSendPointFx_522FF0(code C.char, pos *C.float2) C.int {
	return C.int(visibilityFXPoint(byte(code), *(*types.Pointf)(unsafe.Pointer(pos))))
}

func nox_xxx_playerCheckSpellClass_57AEA0(class, ind C.int) C.int {
	return C.int(playerSpellClassCheck(int32(class), int32(ind)))
}

func nox_xxx_playerResetImportantCtr_4E4F40(to C.int) C.int { return C.int(reliableResetRate(int(to))) }

func nox_xxx_protectPlayerHPMana_56F870(id C.int, value C.ushort) C.uint32_t {
	return C.uint32_t(setProtectionRecord(int32(id), uint32(value)))
}

func nox_xxx_protectionStringCRCLen_56FAE0(data *C.int, size C.uint) C.int {
	return C.int(protectionStringChecksum(unsafe.Pointer(data), uint32(size)))
}

func nox_xxx_sendArrowTrapFX_5238A0(pos *C.float, extra C.char) {
	visibilityFXArrowTrap(*(*types.Pointf)(unsafe.Pointer(pos)), byte(extra))
}

func nox_xxx_serverHandleClientConsole_443E90(pl *nox_playerInfo, action C.char, text *C.wchar2_t) C.int {
	return C.int(consoleCommandRemote(asPlayerS(pl), byte(action), GoWStringP(unsafe.Pointer(text))))
}

func nox_xxx_tileAllocTileInCoordList_5040A0(a0 C.int, a1 C.int, a2 C.float) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabTileNew(int32(uint32(a0)), int32(uint32(a1)), math.Float32bits(float32(a2)))))
}

func nox_xxx_tileCheckByte3_544070(index C.int) C.int {
	return C.int(bool2int(selectBorderPrimary(int32(index))))
}

func nox_xxx_tileCheckByte4_5440A0(variation C.int) C.int {
	return C.int(bool2int(selectBorderVariation(int32(variation))))
}

func nox_xxx_tileCheckImageVari_51D570(variation C.int) C.int {
	return C.int(bool2int(selectTileVariation(int32(variation))))
}

func nox_xxx_tileCheckImage_51D540(index C.int) C.int {
	return C.int(bool2int(selectTileImage(int32(index))))
}

func nox_xxx_tile_51D5C0(value C.int) C.int {
	return C.int(bool2int(setTileFlag(int32(value))))
}

func nox_xxx_toxicCloudPoison_53D9D0(a1 C.int, a2 C.int) {
	temporaryCloudCandidate(objectFromInt(int32(a1)), objectFromInt(int32(a2)), true)
}

func nox_xxx_waterBarrel_53CC30(a1 *C.float, a2 C.int) {
	temporaryWaterCandidate((*server.Object)(unsafe.Pointer(a1)), *(*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))))
}

func nox_xxx_waypointNext_579870(a1 C.int) C.int {
	if a1 == 0 {
		return 0
	}
	return C.int(waypointRaw(waypointFromRaw(int32(a1)).WpNext))
}

func nox_xxx_wndDraw_49F7F0() { objectRenderSaveClip() }

func sub_411490(index, edge C.int) C.int {
	return C.int(normalizeBorderEdge(int32(index), int32(edge)))
}

func sub_41CEE0(info unsafe.Pointer, all C.int) C.int {
	return C.int(playerFileClientWrite(info, int(all)))
}

func sub_436550() C.int { return C.int(interactionFrameGate()) }

func sub_43AF30() C.int { return C.int(browserUI.hosting) }

func sub_45A010(dr *nox_drawable) *nox_drawable { return (*nox_drawable)(asDrawable(dr).Field_104.C()) }

func sub_48C6B0(x, y C.int) C.uint { return C.uint(screenDistance(int32(x), int32(y))) }

func sub_49F860() C.int { return C.int(objectRenderRestoreClip()) }

func sub_4C5050() { objectRenderBeamReset() }

func sub_4E4F30(to C.int) C.int { return C.int(ReliableResetSequence(ntype.PlayerInd(to))) }

func sub_4E55F0(to C.uchar) C.int { return C.int(reliableRemoveRecipient(byte(to))) }

func sub_4E8E50() *C.uchar { return (*C.uchar)(unsafe.Pointer(worldQuestPending())) }

func sub_4E8E60() C.int { return C.int(worldQuestCountdown()) }

func sub_4E9010() C.int { return C.int(bool2int(worldQuestExitReady())) }

func sub_504290(a0 C.char, a1 C.char) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabWallNew(byte(uint32(a0)), byte(uint32(a1)))))
}

func sub_51D2C0(source, target C.int) C.int {
	return C.int(bool2int(appendWaypointLink(waypointFromRaw(int32(source)), waypointFromRaw(int32(target)), int8(memmap.Uint8(0x973F18, 35972)))))
}

func sub_51D300(source, target C.int, kind C.char) C.int {
	return C.int(bool2int(appendWaypointLink(waypointFromRaw(int32(source)), waypointFromRaw(int32(target)), int8(kind))))
}

func sub_51DD50(x, y, flags, key C.int) {
	pushTileFill(int32(x), int32(y), int32(flags), int32(key))
}

func sub_51DE30(x, y, flags *C.uint32_t) C.int {
	return C.int(bool2int(popTileFill((*uint32)(unsafe.Pointer(x)), (*uint32)(unsafe.Pointer(y)), (*uint32)(unsafe.Pointer(flags)))))
}

func sub_53BD10(a1 C.int, a2 C.int) {
	temporaryAntiCandidate(objectFromInt(int32(a1)), objectFromInt(int32(a2)))
}

func sub_53D8C0(a1 C.int, a2 C.int) {
	temporaryCloudCandidate(objectFromInt(int32(a1)), objectFromInt(int32(a2)), false)
}

func sub_543E60(record, category C.int) C.int {
	p := (*[4]uint32)(unsafe.Pointer(uintptr(uint32(record))))
	return C.int(bool2int(mergeBorderEdge(p, int32(category))))
}

func sub_544020(name *C.char) C.int {
	return C.int(bool2int(selectBorderName((*byte)(unsafe.Pointer(name)))))
}

func sub_56F780(id, value C.int) C.uint32_t {
	return C.uint32_t(setProtectionRecord(int32(id), uint32(value)))
}

func sub_56F820(id C.int, value C.uchar) C.uint32_t {
	return C.uint32_t(setProtectionRecord(int32(id), uint32(value)))
}

func sub_5798A0(a1 C.int) C.int {
	if a1 == 0 {
		return 0
	}
	return C.int(waypointRaw(waypointFromRaw(int32(a1)).WpNext))
}

func sub_579E70() *C.uint32_t {
	wp, _ := alloc.New(server.Waypoint{})
	wp.Flags |= 0x1000000
	return (*C.uint32_t)(unsafe.Pointer(wp))
}

func sub_57C790(line *C.float4, point, out *C.float2, length C.float) {
	projectLineClamped((*[4]float32)(unsafe.Pointer(line)), (*types.Pointf)(unsafe.Pointer(point)), (*types.Pointf)(unsafe.Pointer(out)), float32(length))
}
