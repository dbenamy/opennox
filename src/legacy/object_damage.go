package legacy

/*
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4_3.h"
int nox_xxx_soundPlayerDamageSound_5328B0(nox_object_t* a1, nox_object_t* a2);
*/
import "C"
import "github.com/opennox/opennox/v1/server"

func init() {
	server.RegisterObjectDamageGo("DefaultDamage", C.nox_xxx_damageDefaultProc_4E0B30, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageDefault(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("SkeletonDamage", C.sub_4E23C0, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageSkeleton(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("PlayerDamage", C.nox_server_handler_PlayerDamage_4E17B0, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damagePlayer(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("StoneDamage", C.sub_4E24B0, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageDefault(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("MechGolemDamage", C.sub_4E24E0, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageMechGolem(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("FlammableDamage", C.nox_xxx_damageFlammable_4E2520, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageFlammable(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("BlackPowderDamage", C.nox_xxx_damageBlackPowder_4E2560, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageBlackPowder(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("ArmorDamage", C.nox_xxx_damageArmor_4E1500, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageArmor(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("WeaponDamage", C.sub_4E14B0, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageWeapon(u, source, weapon, amount, kind) != 0
	})
	server.RegisterObjectDamageGo("BallDamage", C.sub_4E14A0, func(_, _, _ *server.Object, _, _ int32) bool { return false })
	server.RegisterObjectDamageGo("MonsterGeneratorDamage", C.nox_xxx_damageMonsterGen_4E27D0, func(u, source, weapon *server.Object, amount, kind int32) bool {
		return damageGenerator(u, source, weapon, amount, kind) != 0
	})

	server.RegisterObjectDamageSound("DefaultDamageSound", C.nox_xxx_soundDefaultDamageSound_532E20)
	server.RegisterObjectDamageSound("PlayerDamageSound", C.nox_xxx_soundPlayerDamageSound_5328B0)
}
