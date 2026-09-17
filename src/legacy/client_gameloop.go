//go:build !server

package legacy

/*
#include "GAME1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "GAME3_3.h"
#include "GAME5.h"
#include "GAME5_2.h"

extern unsigned int nox_client_gui_flag_1556112;
extern unsigned int dword_5d4594_2650652;
*/
import "C"
import "github.com/opennox/opennox/v1/client/noxrender"

func Sub_43CEB0() {
	C.sub_43CEB0()
}
func Sub_40A710(a1 int) uint32 {
	return uint32(serverConfigConnectionRate(int32(a1)))
}
func Nox_client_screenParticlesDraw_431720(vp *noxrender.Viewport) {
	screenParticlesDraw(vp)
}
func Nox_client_newScreenParticle_431540(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 int) {
	screenParticleCreate(a1, a2, a3, a4, a5, a6, byte(a7), byte(a8), byte(a9), byte(a10))
}
