package legacy

/*
#include "server__object__die__die.h"
#include "GAME4_3.h"
#include "GAME5.h"

*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
)

func init() {
	server.RegisterObjectDeathGo("PlayerDie", C.nox_xxx_diePlayer_54D2B0, func(u *server.Object) { playerDeath(u) }, 0)
	server.RegisterObjectDeathGo("PotionDie", C.nox_xxx_diePotion_54CBB0, objectDeathPotion, 0)
	server.RegisterObjectDeathGo("ImpEggDie", C.nox_xxx_dieImpEgg_54CAE0, func(u *server.Object) { objectDeathImpEgg(u) }, 0)
	server.RegisterObjectDeathGo("GlyphDie", C.nox_xxx_dieGlyph_54DF30, func(u *server.Object) { Nox_xxx_dieGlyph_54DF30(u) }, 0)
	server.RegisterObjectDeathGo("BarrelDie", C.nox_xxx_dieBarrel_54DFA0, objectDeathBarrel, 0)
	server.RegisterObjectDeathGo("CreateObjectDie", C.nox_xxx_dieCreateObject_54E010, func(u *server.Object) { objectDeathCreate(u, false) }, 132)
	server.RegisterObjectDeathGo("SpawnObjectDie", C.nox_xxx_dieSpawnObject_54E070, func(u *server.Object) { objectDeathCreate(u, true) }, 132)
	server.RegisterObjectDeathGo("PolypDie", C.nox_xxx_diePolyp_54CB10, diePolyp, 0)
	server.RegisterObjectDeathGo("MarkerDie", C.nox_xxx_dieMarker_54E460, objectDeathMarker, 0)
	server.RegisterObjectDeathGo("WeaponDie", C.nox_xxx_dieWeapon_54E370_obj_die, objectDeathWeapon, 0)
	server.RegisterObjectDeathGo("ArmorDie", C.nox_xxx_dieArmor_54E170_obj_die, objectDeathArmor, 0)
	server.RegisterObjectDeathGo("BoulderDie", C.nox_xxx_dieBoulder_54E4B0, objectDeathBoulder, 0)
	server.RegisterObjectDeathGo("GameBallDie", C.nox_xxx_dieGameBall_54E620, func(u *server.Object) { objectiveBallReset(u) }, 0)
	server.RegisterObjectDeathGo("MonsterGeneratorDie", C.nox_xxx_dieMonsterGen_54E630, generatorDeath, 0)

	server.RegisterObjectDeathParse("CreateObjectDie", resourceObjectParser("death", "spawn"))
	server.RegisterObjectDeathParse("SpawnObjectDie", resourceObjectParser("death", "spawn"))
}
