# Glyph and item eligibility — 2026-09-11

Scope: glyph-selection predicate 57B400 and item eligibility 57B450. Both have
independent lazy Glyph type caches. Type lookup runs before all absence gates;
a zero lookup is retried, while a nonzero result survives later registry changes.
57B400 rejects missing player, and only wizards may select Glyph. 57B450 also
requires a drawable and local player; glyph restriction precedes the cheat bypass.
Only after these gates does it read the class-mask dependency.

## Original-C baseline

8,556 calls cover all 256 class bytes against eight masks, all 256 masks for the
three legitimate classes with glyph/non-glyph types, absence/cheat/glyph gates,
separate cache fill/reuse/zero retry, signed-bit type IDs, and the Go item wrapper.
Fixtures use a minimal real client type lookup and C-owned player/drawable/type
storage. They check exact results, cache values, lookup counts, class-mask callback
arguments and readonly player/drawable bytes, restoring all globals afterward.

The current 386 C machine code masks class shift counts to five bits before
narrowing to a byte. Baseline tests cover that observed ABI even for invalid
class bytes; this is not a claim that out-of-range signed C shifts are portable.
Native code must compute the class mask before calling Sub_57B370, as C does.

All cases pass against original C before replacement. Production C before this
chunk: **141,082 physical lines**, 153 files, zero reference C. Local artifacts:
`build/port-glyph-eligibility/`.

57B400 retains a C caller in GAME2_2.c. 57B450 has only its Go wrapper and may
retire its C bridge. Keep C nox_cheat_allowall: GAME3_3.c still uses it. The two
cache words have no other references and can move to private Go state.

## Go conversion

Original-C baseline: `d6d7c136`. Both predicates now execute Go. Only 57B400
retains its live C entry point; the item wrapper calls the private helper directly.
The two lazy caches moved from otherwise-unused blob words to private Go uint32
state. The live C cheat flag remains shared with its other C caller.

The native item predicate computes the byte class mask before the class-mask
callback and keeps the observed 386 shift-count masking. Cache lookup still
precedes all early returns. Nil item wrapper input remains supported: Drawable.C
already accepted a nil receiver in the original wrapper.

Production C: **141,042 physical lines (−40)** in 153 files; reference C: **0**.
All accumulated protection/network/waypoint/rules/spell-class/ping/glyph tests
pass on 386 default/server/highres. All three production binaries build. Fresh
glyph-eligibility-port gameplay passes both preserved screenshots with overrides
disabled. Full suite was last repeated at the immediately preceding alias
milestone, where its 1,553 failure entries exactly matched the known baseline.
This small predicate chunk used targeted accumulated tests, builds and gameplay.
