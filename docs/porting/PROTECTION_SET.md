# Protection integer/byte/word setters — 2026-09-10

This chunk replaces `sub_56F780`, `nox_xxx_playerResetProtectionCRC_56F7D0`,
`sub_56F820` and `nox_xxx_protectPlayerHPMana_56F870`. They share one native Go
implementation; all four keep C exports because C callers remain. The public
Go reset wrapper calls Go directly.

An ineligible signed handle (<657757279) returns its original raw bits without
changing state. An eligible missing handle returns zero. A match updates the
first matching record's encoded value and checksum, calls rekey/shuffle, and
returns the new key. Integer setters preserve all 32 bits; byte/word arguments
zero-extend after their ABI-width truncation. Duplicate IDs retain first-match
behavior.

The decompiled C return declarations were pointers even though every path
returns scalar handle/key bits or zero. The caller audit found no returned-pointer
dereferences or function-pointer uses. All results are ignored except one HP/mana
setter result converted to unsigned int in the mana-clamp path. Declarations now
use uint32_t, preserving the observed 386 return bits and avoiding fabricated Go
pointers. The pre-port test adapter converts old pointer results to integer bits
inside C; it also works with the corrected scalar declarations.

Tests reuse the rekey fixture's independent floating-RNG oracle and C-allocated
list setup. 500 deterministic scenarios × five call paths compare full decoded
record order, original node identities/links/endpoints, count/sequence/checksum,
return bits, both server RNG indices, floating RNG raw state/range, and wrapping
swap/rekey counters. Cases cover empty and populated lists, duplicates, zero and
high-bit keys, signed eligibility boundaries, misses, and byte/word truncation.
Misses/ineligible handles must leave all manager and RNG state unchanged.

Final pre-port revision, post-port checks and C count follow after validation.

The final original-C setter/rekey baseline passes; local evidence is in
build/port-setters/c-before.log. Primary-owned tests and Terra's bounded
implementation draft received reciprocal review before replacement.

Pre-port tests are recoverable at `8ec1b116`. Production C after conversion is
**142,189 physical lines**, down **76**, in 153 files; reference C is **0**.
See [C_LOC.md](C_LOC.md).

All accumulated protection tests pass under default/server/highres on 386,
including the unchanged original-C setter expectations. All three production
targets build through `go run ./internal/noxbuild -o ../build/port-setters/bin`;
symbol checks confirm all four setters now use Go export bridges. Fresh
`setters-port` warrior gameplay exits 0 against both preserved screenshots with
overrides disabled. Logs/binaries are under build/port-setters. No new full-suite
run was needed after the rekey milestone's unchanged failure set.
