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

## C baseline in development

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
full/partial-turn modes. These latter fixtures are still being qualified.

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

The two C guards are installed; no native translation is installed. Baseline
captures are preliminary, not frozen. Evidence and scoped draft plan are under
build/port-color-light. Complete source audit, independent contracts, repeated C
captures, affected targets and gameplay before replacing the six functions.

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
tests-result.json. Repeat, affected checks and real gameplay are now running.
Static final preflight passes.

The C guards add six lines: current **61,071 / 82 / zero reference C**. The six
functions now occupy **268 lines**. A complete source-reference audit also found
three map-classification helpers left unreferenced by the preceding map conversion:
sub_44D040, sub_44D060 and sub_44D090 (28 lines). Retire these with native integration,
including their declarations, and verify all eight retired symbols are absent.
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
All source readers are joined. Commit and push this C correction and recoverable
baseline before native integration. No Go light animation is installed yet.
