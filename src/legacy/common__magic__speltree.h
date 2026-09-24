#ifndef NOX_PORT_COMMON_MAGIC_SPELTREE
#define NOX_PORT_COMMON_MAGIC_SPELTREE

#include "defs.h"

int nox_xxx_spellNameToN_4243F0(char* id);
int nox_xxx_spellGetAud44_424800(int ind, int a2);
int nox_xxx_spellManaCost_4249A0(int ind, int a2);
wchar2_t* nox_xxx_spellDescription_424A30(int ind);
unsigned int nox_xxx_spellFlags_424A70(int ind);
void* nox_xxx_spellIcon_424A90(int ind);
void* nox_xxx_spellIconHighlight_424AB0(int ind);
int nox_xxx_spellFirstValid_424AD0();
int nox_xxx_spellNextValid_424AF0(int ind);
bool nox_xxx_spellIsValid_424B50(int ind);
bool nox_xxx_spellIsEnabled_424B70(int ind);

#endif // NOX_PORT_COMMON_MAGIC_SPELTREE
