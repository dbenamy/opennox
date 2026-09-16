# Summon-creature panel port

## Qualified native result

The complete summon-creature panel is native: **36 C functions**, **994 C lines
removed**, **63,323 lines / 85 files / zero reference C** remaining. Six C adapters
remain and thirty private C interfaces retire. All **971 summon records in 14
groups**, plus the unchanged book/quickbar records, pass: **7,954 / 63 groups** in
each of default/server/highres, no skips or changed goldens.

Three production builds and ABI audits pass; the full asset suite has exactly its
known failures. The fresh conjurer replay matches all 20 frozen frames, and the
flat replay matches while regenerating the warrior map exactly (51 copied maps
removed, one regenerated; 51.202 driver seconds). All readers are joined. Final
source is identical across the three target and production manifests; the staged
source proof is build/port-summon/native-index-proof.json (1,821 files).

Review items: the intentional menu-bottom-clamp prerequisite, preserved unsigned
opening animation and incidental ignored add return, and corrected global-coordinate
hit testing. No unresolved blocker remains. The remaining sections record baseline
and development history; provisional statements describe those earlier stages.

## Original C baseline scope

Quickbar is qualified and pushed in **657c6dc4**. The next connected candidate
covers 36 functions / roughly 994 physical C lines: GAME3_1.c from 004C1CA0 up
to 004C3390, plus client__gui__guisumn.c. Production remains C while the baseline
is constructed. Current production count: **64,317 lines / 86 files / zero
reference C**. Baseline wrappers call the actual C implementation.

The owner reuses spellbook/inventory GUI, renderer, strings, netlist and drawable
fixtures. A variadic extra-name argument extends the spellbook fixture to supply
creature metadata; existing callers retain identical behavior. Raw blob state
and extracted named globals are owned separately. Pointer identities are
normalized only for known records, windows, images and callbacks.

## Planned coverage

- All 16 active-record masks: first/next traversal, duplicate-code lookup,
  first-free allocation, full capacity and byte-exact deactivation semantics.
- All 16 occupancy masks, all valid positions for footprints 1/2/4, grid lookup
  bounds, painting and clear return-address behavior.
- All 256 assignments of absent/small/medium/large to four slots, including
  over-capacity layouts; independent occupied-cell contracts for fitting layouts.
- Real constructor, selection/events, command messages, tooltip and text behavior,
  opening/closing animation, sprite highlights, health bars and frame pixels.
- Add/duplicate/full/remove lifecycle, metadata classification and plant command
  eligibility; actual drawable owners for effects that depend on live creatures.
- Repeated C captures, affected target configurations, focused native comparisons,
  production ABI and an appropriate headless integration scenario.

Do arithmetic/signedness and pointer review before final long gates. In
particular the C slide coordinate is an unsigned named word despite negative
initial coordinates; preserve observed comparison semantics rather than inferring
an intended animation. The menu's bottom clamp uses screen width in the C code;
exercise nonsquare screens before choosing whether any prerequisite fix is needed.

## Review notes

The add-creature function has a `char` return that sometimes contains the low byte
of a record address. Its sole production caller (cdecode.c) ignores the return.
Test state, messages, sounds and visibility; do not freeze an allocation-dependent
byte as a portable expectation. Other meaningful pointer returns are normalized
against actual owned addresses. No production behavior has changed yet.

No subagents are active. Ignored exploration/drafts are under build/port-summon;
prepare_baseline.py is an initial generator, not a safe resume action after edits.

## Baseline development findings

The first eight C groups pass. Extending menu drawing initially looked up the
command via Window.ID; C stores it in word 8 (widget data), so the fixture now
finds the real command child using that field. No production change was needed.
Health records, guide name pointers and guide images are owned explicitly, which
allows real health-bar drawing and both icon/fallback branches without replacing
lookup or rendering algorithms. Live drawable list ownership is restored before
fixture cleanup; highlight tests verify preservation of all unrelated flags.

The corner contract proves the old bottom-clamp bug: at 640×480, the menu becomes
(0,562)–(148,639). Changed `nox_win_width` to `nox_win_height` in the original C
constructor, then ran all fourteen groups successfully (7.355 test seconds).
This intentionally fixes an existing behavior before baseline freeze; see the
decision log. Twelve corner cases now pass across landscape and portrait screens.
The final capture candidates contain **971 records / 14 groups**. Hashes are
installed in the tests and summon-captures.json; repeated/target qualification is
running. No baseline commit or native conversion has occurred yet.

Completed replay asset maintenance reclaimed 2,189,186,472 bytes from four quickbar
runs. Only files matching the original extracted asset by SHA-256 were removed;
run-local deduplicated-assets.json preserves paths/hashes/modes/timestamps and a
restore command. Screenshots, scenario files, logs and modified/generated assets
remain. This maintenance does not alter gameplay results or source qualification.


## Frozen C contract checkpoint

All fourteen groups / 971 records match in default, an independent repeat,
server and highres: phase seconds 30.875 / 12.797 / 109.557 / 39.222. The combined
book/quickbar/summon run also passes **63 groups / 7,954 records**, including every
unchanged book and quickbar hash (70.622 step seconds). All source readers are
joined. The original C implementation, with the documented one-line menu fix,
remains installed; a baseline contract commit protects this state while the
conjurer integration scenario is being developed. Do not call the gameplay
baseline qualified yet or replace C before that scenario is captured/repeated.

The initial conjurer attempt removed its starting map to force decompression,
which prevented the campaign from starting: checkHasSoloMaps explicitly checks
for the expanded con01a.map file. Keep that file for the conjurer scenario.
This is a scenario prerequisite, not a summon bug; no unrelated map-availability
change is needed. Continue using the separate warrior/flat replay to exercise
forced map regeneration. The failed attempt's screenshots and logs are preserved.

External caller audit suggests four required C interfaces after conversion:
sub_4C1CA0, nox_xxx_cliSummonCreat_4C2E50,
nox_xxx_cliSummonOnDieOrBanish_4C3140 and sub_4C3260. Three other entry points
have only Go callers and can move with their callers. Audit final symbols and
callbacks after translation rather than treating this preliminary list as proof.


## Gameplay baseline qualified

C-contract checkpoint **f40d94dc** is pushed. The new conjurer scenario visibly
opens the creature guide, assigns Summon Bat, casts it, shows the cage and health
bar, orders Guard, opens the individual bat menu, and banishes it until the cage
closes. The original capture and fresh repeat both pass with exact frame matches.
See summon-replay.json for all frame hashes, scenario and binary identity, and
measured timings. No source changed after the qualified C-contract checkpoint.
The source index remains 1,818 files, matching all C qualification manifests.

All C tests/build/scenario readers are joined. The baseline is now ready for Go
replacement. Final audit identifies **six** interfaces to keep: the four external
C entries above plus sub_4C2C20 and sub_4C2CE0, since the GUI still invokes tooltip
callbacks through C pointers. Keep those adapters; the other thirty functions can
be private Go functions once their callers move.


## Native implementation installed; qualification pending

The 36 selected C functions are now implemented in four gui_summon_*.go files.
Six interfaces remain for outside C calls and tooltip callbacks; thirty private
interfaces retire, and the three Go caller wrappers invoke Go directly. Removed
739 lines from GAME3_1.c and all 255 lines of guisumn.c: **994**, provisionally
leaving **63,323 / 85 files / zero reference C**. No C algorithm is kept for tests.

Before starting native comparisons, reviewed record layout (32 bytes), signed
coordinate division, uint32 animation comparisons, byte-only state transitions,
menu command storage (window word 8), exact painter/iterator return values and
callback ownership. The menu-bottom prerequisite remains intentionally fixed.
The draft had used quickbar's different widget-data accessor; corrected to offset
32 before installation. Drawable flags use ObjFlags at 120, not Flags70 at 280.

The native qualification manifest runs all 63 affected groups in each target,
then three production builds/ABI, the known full asset suite, the 20-frame conjurer
replay and a separate flat replay with forced warrior-map regeneration. The
complete accumulated port corpus already passed at the immediately preceding
quickbar milestone; no shared production infrastructure changes require another
complete corpus here. Broaden if failures leave the affected scope uncertain.


The first native run passes twelve of fourteen groups. Two groups differ only in
hover behavior: child-window hit tests used Window.PointIn, whose coordinates are
local, whereas the original C entry uses global position plus size. The native
owner now calls that existing Go legacy hit-test implementation directly through
a small typed adapter. Highlights, hover sounds and hover colors are rechecked
against unchanged captures. This is a translation correction, not a golden update.
The focused native retry uses the c-default command from the manifest but runs
the installed Go owner; its directory is native-focused-02, not C evidence.


The corrected native focused run passes all **971 records / 14 groups** with
unchanged hashes (117.302 step seconds, including compilation). All three affected
target phases now run independently; highres uses GOMAXPROCS=1. A comment-only
clarification in the control-event test names window word 8 as widget data;
final target manifests include it. Finish these and production/replays before
accepting the conversion or advancing C_LOC.


Final affected checks pass all **7,954 records / 63 groups**, no skips, all hashes
unchanged: default 87.996 driver seconds, server 168.256, highres 97.015.
All three readers are joined. Production qualification is the only remaining
source reader; the default binary and its six-retained/thirty-retired ABI audit
pass. Finish the other production gates before acceptance.


All three production builds and interface audits pass: default 115.714 seconds,
highres 12.636, server 125.014. The final gate continues with the full asset-suite
comparison and fresh gameplay. These timings include cold CGO compilation and
are not comparable to warm focused test execution alone.


The production core gate passes: the asset suite is exactly the known **1,553
failure entries**, with 15 passing / 3 failing / 32 skipped packages. The fresh
conjurer replay matches all **20 C-reference frames** (58.829 process seconds),
including the guarding response and final disappearance of the cage. The separate
flat replay is the final remaining gate. All final source manifests match 1,821
files. The additional successful C scenario asset copies were deduplicated with
the same hash-checked, restorable maintenance procedure (2,225,554,860 bytes).

Preserved behavior for later review: the unsigned opening-coordinate comparison
snaps the negative initial position to the open target on the first frame; closing
still follows its existing stepped movement. The port preserves that C behavior.
The menu-bound correction is the only intentional gameplay behavior change here.

The final flat replay also passes; production phase completion is recorded in
build/port-summon/production/result.json. This completes summon qualification.
