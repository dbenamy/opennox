#ifndef NOX_SERVER_SCRIPT_INTERNAL_H
#define NOX_SERVER_SCRIPT_INTERNAL_H

#include "defs.h"

void nox_script_push(int val);
void nox_script_pushf(float val);
int nox_script_pop();
float nox_script_popf();

#endif // NOX_SERVER_SCRIPT_INTERNAL_H
