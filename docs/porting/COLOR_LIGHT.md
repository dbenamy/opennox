# Colored-light animation and viewport updates

Starting point: map drawable conversion **b593c7ae**, committed and pushed;
**61,065 C lines / 82 files / zero reference C**.

Scope: six connected functions in GAME3_2.c, 004CE390–004CEBA0 exclusive,
**262 original C lines**. This completes animation of fields loaded by the
previous map-reader batch. Reuse the real drawable/light owners and setters.
The callback is registered through GAME_data_init.go and blob initialization;
five animation helpers have no remaining external production callers.

This is a smaller coherent batch than the usual 1,000–3,000-line target. Combining
unrelated resource UI would add a different owner and integration surface without
helping this translation. Continue to use the larger target when dependencies fit.

## C baseline and independent contracts

Interpolation covers RGB, intensity and penumbra: 5,040 records over counts
0/1/2/3/16/128/255, periods 0/1/2/3/255/256/257/65535, cyclic/ping-pong modes,
endpoints and high unsigned frame values. Counts are signed chars in the C helper;
negative byte representations disable interpolation. Only valid array indices are
used by enabled cases. Independent contracts cover disabled cases, initial entries,
and unchanged ownership/RNG.

Viewport update has 150 cases across inclusive 100-pixel margins, empty/single/
multiple entry arrays, editor suppression and actual update-list removal.
Direction covers actual static drawable lookup, missing targets, quadrants and
near-axis positions. Rotation covers signed speed, frame boundaries, FPS and
full/partial-turn modes.

All eight independent coincident-target/zero-arc regressions fail against the
unchanged C implementation: the angle changes and the setter clears the light
mode. Two guards now preserve the full current light state where no direction
change is defined. The eight regressions pass after this six-line prerequisite.
This reversible behavior choice is recorded for review; do not capture NaN-to-
integer artifacts as compatibility requirements.

An added whole-drawable mutation contract initially assumed the intensity helper
changed only float intensity and radius. It correctly exposed that the actual
C adapter also updates fixed-point intensity at offset 148. The fixture's allowed
fields and native draft were corrected by tracing the adapter, and an independent
fixed-point consistency contract was added. This was a fixture/draft correction,
not another production behavior change.

Rotation fixtures exclude extreme frame/speed/FPS combinations where the legacy
cycle narrows beyond int32 and its later floating-to-int conversion is outside
C's defined range. Retain high frame cases where these operations remain defined;
do not freeze undefined conversion results. Disabled-property checks separately
verify the penumbra/rotation mode gates. No fixture uses invalid array indices.

Evidence is under build/port-color-light. Development captures precede the final
frozen set below; use the final phase directories for qualification evidence.

## Frozen capture set and final qualification

The frozen set contains **8,853 records / five groups**: 5,040 interpolation,
150 viewport, 40 target direction, 3,047 defined-range rotation and 576 sequential
combined updates. Additional independent cases cover eight degenerate directions
and eight property-gate checks. Sequence cases retain actual update-list ownership
and verify target tracking takes precedence over rotation across frames.

Pre-sequence C checks passed in default/server/highres. They are development
preflights, not the final seven-root baseline. The final seven-root target phases
pass in c-final-default/server/highres, no skips, with all five hashes matching.
Driver seconds: 29.213 / 27.508 / 36.737; exact timings remain in
tests-result.json.
Static final preflight passes.

The C guards added six lines: baseline **61,071 / 82 / zero reference C**. The six
functions then occupied **268 lines**. A complete source-reference audit also found
three map-classification helpers left unreferenced by the preceding map conversion:
sub_44D040, sub_44D060 and sub_44D090 (28 lines). These and their declarations were
removed with native integration; all eight retired symbols are audited.
This is removal of unreachable leftovers, not a separate behavior translation.

Final repeat passes all seven roots/five hashes (6.009s), and the affected
selection passes **34 roots**, no skips (18.738s). It includes map readers,
particle light properties, drawable updates and object rendering. Real C gameplay
matches all **41 frames**, repeats them exactly, and flat gameplay matches all
**14 frames** with exact warrior map regeneration. C gameplay driver steps:
client build 74.635s, first replay 60.536s, repeat 62.903s, flat 49.755s.
Exact process times/binary/frame hashes are in color-light-replay.json.

All six final phases report unchanged source. c-index-proof.json checks **1,871
staged source files** against default/repeat/server/highres/affected/gameplay.
All C source readers were joined before the baseline was committed and pushed as
215e515a, followed by native integration.

## Native integration

C prerequisite and baseline **215e515a** are committed and pushed. The six
animation callbacks/helpers are now Go, sharing the common interpolation logic
and reusing actual particle-light setters. The five private animation C interfaces
and three unused map-classification helpers are removed; the registered update
callback remains. C is **60,775 lines / 82 files / zero reference C**, **296 fewer**
than the corrected baseline (268 animation, 28 cleanup).

The first native focused run passes all seven roots and all **8,853 frozen
records / five hashes**, including independent contracts (112.302s). No native
compiler or behavior corrections were needed after installation. Static checks
pass. Source review verified signed counts, byte interpolation steps, floating
operation order, the required float32 direction narrowing, fixed-point intensity,
callback order and actual list/lookup ownership. No frozen expected results changed.

Affected default/server/highres checks pass **34 / 33 / 34 roots**, no skips, all
hashes matching (18.032s / 114.955s / 46.106s). The server build excludes the
rendering-only TestClientObjectRenderOcclusion. All target readers are joined.
Production qualification passes all three builds and interface checks, including
absence of test helpers and the eight retired names. The full asset suite matches
exactly **1,553 known failure entries**, with **15 pass / 3 fail / 32 skip** package
outcomes. All **41 gameplay frames** and **14 flat frames** match C; exact warrior
map regeneration also passes. Production driver took **282.834s**, followed by
flat gameplay **49.092s**. Exact replay timings and hashes are tracked in
color-light-replay.json; these timings are not controlled performance benchmarks.

native-index-proof.json checks all **1,873 staged source files** against the three
affected phases and production. All four phases report unchanged source, and all
readers are joined. Final C: **60,775 / 82 files / zero reference C**.

Next candidate: the connected server map-object readers/writers and world-object
transfer callbacks, eighteen functions / about 1,173 C lines. Include actual
save/load integration and field/byte contracts, since initial map loading alone
would not qualify the write paths. Candidate audit and plan are ignored under
build/port-object-xfer; no next-batch source changes are part of this conversion.
