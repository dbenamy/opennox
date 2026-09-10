# Protection float updates — 2026-09-10

The float setter (56F8C0) and additive update (56FA40) now use Go. Both retain
signed handle eligibility, first-match lookup, raw-ID return on ineligible IDs,
zero on missing IDs, checksum updates and the existing rekey/RNG behavior.
The C ABI returns uint32 rather than a fake pointer; all four production callers
ignore the return. Setter conversion truncates float32 to signed int64 and keeps
the low 32 bits. Addition first adds the decoded uint32 value to the float32
operand, then truncates. Nonfinite/out-of-range results explicitly produce zero,
matching the observed x87 integer-indefinite low bits.

A standalone 386 C probe initially suggested 64-significand-bit addition, but
the actual Go-hosted C baseline disproved that model: old=1 plus float bits
5effffff loses the low increment. The live test reads x87 control word 027f
(53-bit precision, round-to-nearest-even). The Go runtime also explicitly sets
double precision in runtime/asm_386.s. Consequently float64 addition matches the
hosted C behavior. The baseline test asserts the relevant precision/rounding
bits; no production FPU control or exception-flag handling was found in the game.
This preserves the hosted game, not a standalone C executable's startup defaults.

Original-C baseline commit: `a55a2ef0`. The two C bodies were still active there;
the new pure helpers were not yet connected to production. 1,800 scenarios run
through both C entries (3,600 calls). Every directed value/old-value combination
reaches a record, with additional randomized lists, missing/ineligible/duplicate
IDs, rekey/shuffle, counter wrap and full state/RNG/link checks. Raw float bits
cross the fixture using memcpy, covering signed zeros, subnormals, fractions,
uint32 and signed-int64 boundaries, infinities, and quiet/signaling NaNs.
Negative powers 2^-70 through 2^-15 and adjacent float32 values exercise rounding
thresholds. Independent math/big Float precision-53 nearest-even models provide
the ABI and pure-test oracles. Pure tests additionally check 100,000 seeded raw
bit cases and run on both 386 and amd64.

Production C after conversion: **141,941 physical lines (−43)** in 153 files;
test-reference C: **0**. Verification artifacts are under `build/port-float/`.

All accumulated protection tests pass on 386 default/server/highres. All three
production binaries build and expose both float entries through Go-backed C
bridges. The `float-port` warrior scenario accepts both preserved screenshots
with overrides disabled. Pure tests pass on 386 and amd64. The full suite was
last run in the immediately preceding object chunk and matched its known
failure set exactly; it was not repeated for this focused conversion.
