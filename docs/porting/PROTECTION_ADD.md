# Protection additive updates — 2026-09-10

This chunk replaces integer addition (`sub_56F920`), mana addition
(`nox_xxx_protectMana_56F9E0`) and level-byte addition (`sub_56F980`). A shared
native helper checks signed handle eligibility, finds the first matching record,
decodes its value, adds modulo 2^32, updates encoded value/checksum, then invokes
rekey/shuffle. Ineligible IDs return unchanged raw bits; eligible misses return
zero; successful calls return the new key.

Integer deltas retain their 32 bits. Mana deltas truncate to signed int16 and
sign-extend; level deltas truncate to uint8 and zero-extend. These conversions
must occur before modular addition. The public Go mana wrapper calls Go directly.

As with setters, these functions' decompiled pointer return types carry scalar
bits. C declarations become uint32_t. Five surrounding decompiled pointer-typed
assignments/returns retain explicit C casts through uintptr_t, preserving their
existing interfaces and bits. No direct caller dereferences the returned value;
unrelated caller signatures remain unchanged.

500 deterministic scenarios exercise three C paths and the public Go mana path
(2,000 calls). The first 300 combine boundary old values/deltas, including zero,
signed limits and UINT_MAX; remaining scenarios supply random values. A widened
signed int64 oracle followed by uint32 conversion checks wrapping arithmetic
independently of the implementation's uint32 addition. Full manager/list/RNG
snapshots also cover invalid/missing handles, duplicates, truncation, counters,
and rekey side effects using the already-validated independent floating oracle.

Final validation outcomes and source counts follow after checks complete.

Original-C tests passed and are recoverable at `a1b5153c`. The C definition of
F980 remains in that reference revision; only its missing header declaration
was added for the baseline adapter. Production C after conversion is **142,130
physical lines**, down **59**, in 153 files; test-reference C is **0**. The delta
includes two now-unused C global declarations in explevel.c. See [C_LOC.md](C_LOC.md).

All accumulated protection tests pass on 386 under default/server/highres. All
three targets build through `go run ./internal/noxbuild -o ../build/port-add/bin`,
and symbol checks confirm the three additive updates are Go-backed C exports.
The five caller casts and native mana wrapper were independently reviewed. Fresh
`add-port` warrior gameplay exits 0 against both preserved screenshots with
overrides disabled. Evidence: build/port-add. The latest full-suite checkpoint
remains rekey's exact match to known failures.
