# Collision reflection and containment — 2026-09-11

Scope: live C reflection 57B810 and point containment 57B850. Preserve their C
entry points, exact pointer/result bits, and input aliasing behavior. Existing
Go reflection in object_death_ball.go and analogous server/object.go containment
use different intermediate precision and are outside this C-caller conversion.

## Original-386 baseline

Disassembly of the production binary confirms x87 PC53 arithmetic. Reflection
loads old velocity X before comparing the double-width normal product. On the
nonpositive branch it copies old Y as raw bits into X and stores old X through
x87 into Y; on the other branch both operands pass through x87 and change sign.
Signaling NaNs therefore quiet asymmetrically on the swap branch. Signed zero,
NaNs, infinity and positive products that would underflow float32 matter.

73,984 reflection cases cross sixteen raw IEEE edge values, plus identical and
partial input overlap. An independent sign/classification and raw-bit oracle
checks every word, returned pointer bits and outer guards. It does not multiply
float values to decide the expected branch.

Containment keeps its coordinate sums/subtractions in x87 registers without
float32 spills and loads the unsuffixed 0.70709997 constant as a double. Tests
cover 3,675 analytic diamond-grid cases (including large positions), 25 floating
boundary/extreme cases and 50,000 deterministic raw-word cases, including input
aliasing. Every input word and both guards must remain unchanged; returns are
exactly 0/1. The 53,700 original-C results are packed into a 6,713-byte fixture,
with 2,589 true cases. SHA-256:
`cb769ec14e512cd58b2523769616b98903b7eb74065c1de9cdbac0fe86fa25cf`.

The bitset is asset-free original-C output, not a retained C reference program.
The deterministic inputs remain in the Go test source. Initial baseline capture
used the still-original C implementation; a second run verified the saved bitset
and both complete test groups before conversion. Local artifacts/disassembly:
`build/port-collision-primitives/`.

All **127,684 operations** pass against original C before replacement.
Production C before conversion: **141,042 physical lines**, 153 files, zero
reference C. Do not replace the containment expression with float32 locals or
algebraically reassociate its additions merely because source declarations say
float; the actual compiled precision and operation order are the contract here.

## Go conversion

Original-C baseline: `f85e37ee`. Both live C entries now execute private Go
helpers. Reflection explicitly preserves x87 signaling-NaN quieting, raw swap
bits, overlap behavior and the original pointer return. Containment uses float64
intermediates in the original order, with the double constant retained.

All 127,684 targeted operations pass with native Go. An independent Python
precision probe over the 50,000 random cases found zero mismatches using float64
and 42 mismatches using float32 intermediates, confirming that the saved baseline
actually detects this precision regression. Probe: local check-precision.py.

Production C: **141,000 physical lines (−42)** in 153 files; reference C: **0**.
All accumulated protection/network/waypoint/rules/spell-class/ping/glyph/collision
tests pass on 386 default/server/highres. All three binaries build. Fresh
collision-primitives-port gameplay passes both preserved screenshots with
overrides disabled. Full suite exactly matches the known baseline: 15 passing,
3 known failing, 32 skipped/no-test packages and the same 1,553 failure entries,
with none added or removed. Comparison metadata is under the local artifact dir.
