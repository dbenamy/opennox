//go:build porttest && 386

#include "textflag.h"

TEXT ·Control(SB), NOSPLIT, $0-2
 LEAL ret+0(FP), AX
 FSTCW (AX)
 RET

TEXT ·SetControl(SB), NOSPLIT, $0-2
 LEAL value+0(FP), AX
 FLDCW (AX)
 RET

TEXT ·MXCSR(SB), NOSPLIT, $0-4
 LEAL ret+0(FP), AX
 STMXCSR (AX)
 RET

TEXT ·SetMXCSR(SB), NOSPLIT, $0-4
 LEAL value+0(FP), AX
 LDMXCSR (AX)
 RET

TEXT ·save(SB), NOSPLIT, $0-4
 MOVL out+0(FP), AX
 FSTENV (AX)
 // FLDENV (AX): Go 1.26 classifies this as a destination operand,
 // but its encoder expects a source operand. Emit the standard D9 /4 form.
 BYTE $0xd9
 BYTE $0x20
 STMXCSR 28(AX)
 RET

TEXT ·load(SB), NOSPLIT, $0-4
 MOVL in+0(FP), AX
 // FLDENV (AX): Go 1.26 classifies this as a destination operand,
 // but its encoder expects a source operand. Emit the standard D9 /4 form.
 BYTE $0xd9
 BYTE $0x20
 LDMXCSR 28(AX)
 RET

