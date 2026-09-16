# Options panels

The next connected batch covers main-menu and in-game options, shared volume
previews and settings controls: twelve functions / approximately 889 physical C
lines. C is still installed. Starting checkpoint: **62,571 lines / 83 files /
zero reference C**, after qualified binding conversion 0649e67a.

## Scope and planned contracts

- GAME3.c: 004AA6B0–004AB260 exclusive, 004AD9B0–004AEE30 exclusive;
  client__shell__options.c in full. The already-Go advanced-video panel remains
  its actual dependency. The unrelated disconnect dialog is outside this batch.
- Real GUI parser, sliders, radio/checkbox controls, modal visibility, menu
  animations and close-panel rendering; existing entry/listbox fixtures supply
  reusable input and rendering ownership.
- Volume changes and release notifications, effect/dialog/music enable state,
  zero/nonzero transitions, absent button callbacks, current slider ownership,
  preview requests and real timer state. Exercise actual dialog and audio owners.
- Gamma and sensitivity endpoints, intermediate values and floating-point bits.
- Menu/in-game construction, missing resources, supported window sizes, apply,
  cancel, visibility and cleanup; temporary configuration output as appropriate.
- Extend the qualified binding replay with option interactions and return paths.
  Compare repeated C captures, then native Go against frozen records and frames.

## Audit notes to preserve or review

- Main-menu refresh currently maps controls 331/332 to windowed/fullscreen. The
  disabled original 8/16-bit code is not active behavior.
- Menu viewport refresh sends event 16392 with argument 1; in-game show uses 0.
- Volume-change event 16393 suppresses effect preview while the current input
  window belongs to that slider. Completion event 16396 does not.
- Options centering converts an unsigned subtraction to signed int before
  division. The binding editor's unsigned division is different; do not share
  that arithmetic merely because both windows are centered.
- The approximate 889-line scope is intentionally not padded to meet a LOC goal.
  Use connected affected checks after the binding accumulated milestone, plus
  production/interface/known-suite and relevant gameplay qualification.

Evidence and working drafts: build/port-options. The C baseline is frozen and passes all five fixture phases. No native options
implementation is installed.

## Baseline development

The numeric fixture currently passes 1,320 records through both actual C event
handlers, including all 0–100 values, extremes and unrelated events. Gamma has
independent clamping/config-dirty checks; sensitivity has independent curve and
midpoint checks. No hashes are frozen yet.

Arithmetic finding: C rounds the complete sensitivity exponent to float32 at the
powf call, retaining extra precision for division/subtraction. Rounding each Go
float32 operation independently differs in 41 recorded cases. A single exponent
rounding matches all current records. Preserve this arithmetic explicitly and
verify native output bits against the frozen baseline.

The direct C checkbox function-pointer calls bypass Go window extension handlers.
The real-checkbox regression fails in all six menu/channel combinations against
unchanged C (c-checkbox-before2.log). A scoped correction replaces twenty direct
callback paths with nox_window_call_field_93. It removes 96 lines of callback
pointers and redundant absent-handler branches, leaving 62,475 C lines / 83 files
provisionally. This is a prerequisite C fix, not the later Go conversion.
The after-fix numeric, checkbox and 1,440-case volume matrix passes in
c-checkbox-after2.log. Missing callbacks remain covered. Three-target and gameplay
qualification still remain before committing the correction.

The volume matrix also preserves music-enable behavior: its checkbox handler
sets the volume target from the timer’s current value, superseding the preceding
slider target. Audio readiness and checkbox state can differ when no device is
ready; the fixtures check both independently.

Construction uses the real parser and Go enhancement/advanced-video owners.
The first synthetic resource omitted radio GROUPs and incorrectly cleared sibling
checkbox selection. Correcting the fixture to the asset groups resolves that
ownership mistake; it is not another production behavior change.

## Frozen corrected-C baseline

The baseline has **3,130 frozen records / ten capture groups** and six additional
real-checkbox regressions. All eleven selected tests complete without skips in
default, repeat, server and highres. All 56 affected options/binding/radio/slider/
listbox/entry tests pass without skips. The baseline hashes are tracked in
options-captures.json; phase commands are in options-batch.json.

| Capture group | Records |
| --- | ---: |
| Numeric sliders | 1,320 |
| Volume transitions | 1,440 |
| Constructors | 36 |
| Dialog preview sequencing | 180 |
| Close/apply/cancel visibility | 9 |
| Menu animation completion | 1 |
| Overlay rendering | 32 |
| Actual settings-file round trips | 48 |
| Viewport refresh boundaries | 28 |
| Basic clicks and hover dispatch | 36 |

Settings round trips use the real legacy main-section writer, project its audio
and sensitivity fields into temporary config files, and parse them with the real
legacy parser. Gamma passes through the modern Viper writer and reader. Captures
record only the options under test. Apply/cancel lifecycle tests separately
observe write requests; fresh gameplay exercises the full return path with the
usual E2E write suppression.

Driver times: c-default 28.498s, c-repeat 6.207s, c-server 110.533s,
c-highres 35.606s, c-affected 28.702s. Each phase records unchanged source identity.
The real options gameplay baseline is being developed; fixture qualification alone
does not yet close this prerequisite/baseline commit.

## Qualified C gameplay and prerequisite

The corrected C replay options-c and options-c-repeat matches all **41 frames**:
both options panels, volume-to-zero/mute and re-enable, dialog previews, music,
gamma/sensitivity, both binding-editor transitions and resuming gameplay. Process
times: 56.735s capture and 61.344s repeat; c-gameplay phase 123.646s. The scenario
and frame metadata are tracked in options-panels.yaml and options-replay.json.
The first menu-only development run also passed and verified pointer coordinates.

Review note: game_numeric shows a transient clipping/redraw artifact already in
C, similar to the earlier binding wheel frame. Both C replays agree; later frames
redraw normally and game_resumed is visually complete. Preserve it for translation
and investigate the renderer separately. This is not a new native-Go finding.

The scoped checkbox fix is now qualified. It removes **96 C plumbing lines**,
leaving **62,475 / 83 / zero reference C**. This reduction is a C correction, not
Go conversion. The remaining native scope is **793 physical C lines / twelve
functions**. Static memory-access preflight passes. c-index-proof.json matches
the staged source to all five fixture phases and the C gameplay phase (1,845
source files). No C algorithm is retained solely for tests after translation.
