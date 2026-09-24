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
	server.RegisterObjectDamageValueGo("DefaultDamage", C.nox_xxx_damageDefaultProc_4E0B30, damageDefault)
	server.RegisterObjectDamageValueGo("SkeletonDamage", C.sub_4E23C0, damageSkeleton)
	server.RegisterObjectDamageValueGo("PlayerDamage", C.nox_server_handler_PlayerDamage_4E17B0, damagePlayer)
	server.RegisterObjectDamageValueGo("StoneDamage", C.sub_4E24B0, damageDefault)
	server.RegisterObjectDamageValueGo("MechGolemDamage", C.sub_4E24E0, damageMechGolem)
	server.RegisterObjectDamageValueGo("FlammableDamage", C.nox_xxx_damageFlammable_4E2520, damageFlammable)
	server.RegisterObjectDamageValueGo("BlackPowderDamage", C.nox_xxx_damageBlackPowder_4E2560, damageBlackPowder)
	server.RegisterObjectDamageValueGo("ArmorDamage", C.nox_xxx_damageArmor_4E1500, damageArmor)
	server.RegisterObjectDamageValueGo("WeaponDamage", C.sub_4E14B0, damageWeapon)
	server.RegisterObjectDamageValueGo("BallDamage", C.sub_4E14A0, func(_, _, _ *server.Object, _, _ int32) int32 { return 0 })
	server.RegisterObjectDamageValueGo("MonsterGeneratorDamage", C.nox_xxx_damageMonsterGen_4E27D0, damageGenerator)

	server.RegisterObjectDamageSoundGo("DefaultDamageSound", C.nox_xxx_soundDefaultDamageSound_532E20, func(u, other *server.Object) { nox_xxx_soundDefaultDamageSound_532E20(asObjectC(u), asObjectC(other)) })
	server.RegisterObjectDamageSoundGo("PlayerDamageSound", C.nox_xxx_soundPlayerDamageSound_5328B0, func(u, other *server.Object) { nox_xxx_soundPlayerDamageSound_5328B0(asObjectC(u), asObjectC(other)) })
}
