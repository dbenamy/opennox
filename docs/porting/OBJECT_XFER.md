# Server object serialization

## Scope and status

Next connected batch: eighteen functions / 1,177 physical C lines (1,173 before ownership guards) in GAME3_3.c,
004F3E30 through EOF. This covers common object records, inventory loading,
placement and typed world-object transfer callbacks. Engine code is still C;
no serialization baseline has been frozen yet. Two ownership prerequisite fixes
pass focused qualification, adding four C lines. Current C is **60,779 lines
in 82 production files, zero reference C**.

## Save/load integration prerequisite

[object-xfer-save-load.yaml](object-xfer-save-load.yaml) starts a warrior campaign,
saves with F2, moves, loads with F4 and confirms, then walks after reloading.
Seven checkpoints cover gameplay, before/after save, movement, confirmation,
restored world and resumed play. The confirmation click accounts for the 1.25
window-to-rendered-coordinate scale.

Run tools/porting/run_scenario.py with OPENNOX_REQUIRE_SAVE_LOAD=1,
OPENNOX_UI_SCENARIO pointing to this scenario and OPENNOX_DISPLAY_BINARY pointing
to the qualified binary. Capture first, then compare against that run directory.
The runner uses fresh asset/save copies and headless X with null audio.

The opt-in checker requires an explicit Game Saved event between before_save and
after_save, then server and client reads of the exact saved WORKING map, followed
by both post-reload checkpoints. Nonempty player and map files must agree between
AUTOSAVE and WORKING within that run. Initial campaign autosave alone is insufficient.
Six checker unit tests cover success and missing/misleading evidence.

Evidence under build/baseline/runs:
- object-xfer-save-load-develop: capture exits zero in 46.256s; restored-world
  screenshot visually inspected.
- object-xfer-save-load-repeat: comparison exits zero in 48.569s, all seven frames
  match and the integrated save/load checker passes.
- Both use the qualified colored-light binary, SHA-256
  43cae25fb2e7c2c93ff7ee5acb18059a6138728b517b5244b798a378c353743a.
- Repeat artifacts: Player.plr 1,384 bytes; war01a.map 358,744 bytes. Full run
  results contain hashes. These are observations, not portable golden save hashes.

Complete player saves have differed between development scenarios. Do not hide
this with blanket normalization or assume cross-run save-byte determinism. Freeze
controlled serialization records separately; compare gameplay pixels and actual
within-run save/load artifacts for integration. Unchanged asset copies may be
deduplicated only after successful completion with restoration manifests; original
assets, screenshots, logs and saved outputs remain.

## Qualification plan and audit notes

Use actual cryptfile read/write streams and actual server object/type owners.
Cover old/current versions, stream positions, serialized bytes, flags, script IDs
and handlers, names, inventory links, placement and type-specific payloads. Add
independent field/round-trip contracts and supported allocation-failure cases.
Repeat C captures before freezing; use affected checks during development and a
broader accumulated milestone at this shared serialization boundary. Include
normal gameplay, save/load and flat-renderer map regeneration after translation.

The old common writer writes float coordinates even for arguments whose old
reader expects integer coordinates. Preserve actual behavior; do not invent a
requirement that every historical argument combination round-trips.

TriggerXfer's C declaration takes float, but registration and actual CallXfer
callers pass an object pointer in the first 386 stack slot. The C function initially
interprets those bits as a pointer and later reuses the parameter as scratch.
Exercise the registered callback; never numerically convert an object address to
float. Confirm all callers before choosing the final native export signature.

Reuse the painting fixture's real type registry and object pool where practical.
Its current buffer tracking cannot blindly cover placement rejection: production
FreeObject frees UseData but only clears several other pointers. Track ownership
explicitly. alloc.IsDead only recognizes a sentinel pointer, not whether an
arbitrary allocation has been freed. No production ownership changes are planned
merely to simplify the fixture.

Fixture development found that names allocated by this C reader use libc calloc,
whereas factory-owned buffers use the tracked Go alloc wrapper. The first contract
run caught an incorrect fixture free at the first nonempty name. Keep those
cleanup paths distinct; this is a test ownership correction, not an engine change.

Development progress: 255 independent common-reader cases and 99 exact-byte
writer cases pass against C in c-common-development3.log. This includes signed
version limits, omission/presence, coordinates, flags, teams, inventory counts,
script IDs, extra status, lifetime and 255/256/257-byte name behavior. These are
not yet the final repeated/frozen corpus. Typed callbacks now use real registered
function pointers and actual factory/type allocations; their contracts are in
development.

## Ownership prerequisite

Both new regressions fail against the original C: callback rejection leaves an
unlinked object alive (two instead of one), and disallowed placement with an
inventory faults while following a freed child. The correction returns failed
children to the real pool and clears the already-disposed inventory head before
freeing its parent. See [DECISIONS.md](DECISIONS.md). No source edits occur while
the default/server/highres checks are running.

Typed default and historical checks, invisible-light read/write checks and
positive/partial inventory loading pass before these corrections. Historical
trigger scripts require the actual legacy file-handle registry; the fixture now
initializes it and closes registered handles after each case. Initial omission
of that setup was a fixture failure, not a serialization defect.

After-change qualification: default/server/highres each pass ten roots with
681 subcases plus the two ownership regressions, no skips. Driver seconds:
6.766 / 119.484 / 170.263. Static checking passes (static-ownership.log).
This qualifies the prerequisite correction and current contracts, not the final
serialization baseline or a native conversion. Current-source production builds,
replays, remaining boundary tests and frozen capture comparisons are still due.
The tracked object-xfer-batch.json manifest records repeatable qualification
commands; c-ownership-qualified repeats all three successfully in 17.173s and records
unchanged source fingerprints. ownership-index-proof.json checks that tested
source matches the staged source archive.

## Full C corpus

The expanded corpus selects **1,174 complete state records / twelve groups** in
object-xfer-captures.json, alongside independent stream, field and ownership
contracts. Every byte of the 772-byte legacy object ABI and its type-owned data
is captured; declared pointer slots record presence, while relationship identity
and order have separate contracts. Go server handles and allocation addresses do
not become golden data. The groups include 270 door geometry cases with signed
positions, historical formats, trigger dimensions, glyph names, linked mover and
transporter records, and pickup script names.

A final added stale-TOC case exposed the existing C factory adapter bypassing the
nil-type guard. Original final-default/server/highres attempts fail on that case
and are not qualification evidence. Routing the adapter through the guarded
server factory makes the focused corpus pass. All existing capture expectations
are unchanged. Corrected final target runs have suffix 2. Default, independent repeat, server
and highres each pass 21 roots / 1,297 leaf cases without skips and match all
twelve frozen capture hashes. Batch-driver seconds: 26.715 / 6.395 / 104.386 /
34.848. Static memory-access checking also passes (static-c-final.log).

c-gameplay succeeds with unchanged source fingerprints in 322.763s: a fresh
production client matches 41 gameplay frames, two actual save/load runs match
seven frames each and satisfy the explicit save/reload artifact checks, and flat
rendering matches 14 frames with exact map regeneration. These runs include the
ownership and guarded-factory corrections. The broader C accumulated regression
milestone succeeds: 1,198 selected top-level checks execute and finish (1,181 in
the root package), with only the existing opt-in TestMapPopulationPrerequisiteProbe
skip. Driver time is 801.459s; every batch records unchanged source fingerprints.
The staged-source proof checks all six qualification phases against 1,890 source
files (source archive tree ead6528fef353fe2164831c54926341f8cfbc648).

The production C client SHA-256 is
703e627243a0424786d3a9c05f3008965a930f4f32039f74c6df3b3825d8f577.
Actual save/load player files differ between repeated runs (1,376 / 1,384 bytes),
while the saved map matches; each run independently proves that its saved files
were loaded. This reinforces the decision to compare controlled serialization
bytes and within-run saved artifacts, rather than assume whole saves are identical
across runs. Frozen expectations are now recoverable in Git; C remains 60,779
lines / 82 production files / zero reference C until the native conversion.

Translation review caught an important stream detail before installation: old
integer coordinates are read as one eight-byte operation, while float coordinates
use two four-byte operations. The existing stream checksum depends on operation
boundaries, so preserve the original grouping even when decoded values and stream
positions would be identical. The preparatory Go draft was corrected accordingly;
no C expectation changed.

## Native conversion in progress

The C baseline is committed and pushed as 62a006d2. Four Go implementation/export
files replace all eighteen scoped functions. Seventeen existing interfaces remain
for C callers/registered callbacks, while the private historical reader is retired;
Go common/inventory/placement wrappers invoke Go directly. The trigger pointer
signature is documented in DECISIONS.md. Working C is 59,602 lines / 82 files.

The first installed native-focused phase passes 21 roots / 1,297 leaf cases and all
twelve frozen capture hashes in 114.871s, including compilation. No implementation
or expectation changes were required after installation. Static checking passes.
Accumulated default/server checks and full production qualification are running;
accumulated highres is still due. Do not treat this as final qualification yet.

All three native production binaries now pass their ABI audits (17 retained
Go-backed symbols, one retired private symbol, no test helpers). The full asset
suite matches exactly: 1,553 known failure entries and package outcomes 15 pass /
3 fail / 32 skip. Native normal gameplay matches 41 frames and the explicit
save/load scenario matches seven frames, proves saved-map reload, and resumes.
Its saved map also matches C's SHA-256
f4249247a267beaa3e6fb9e426769c9f0bd31a6d0fd1e95f22f22ad07b1ddc99.
Flat rendering also matches all 14 frames with exact map regeneration. The whole
production phase succeeds in 465.107s with unchanged source fingerprints. Client
SHA-256: 56ac204ae1e7f3035aacffa414b70bf35ec9b341ac535e276c8c55de68b3bf2a.
The final accumulated sweeps remain pending.

## Final native qualification

All five native phases pass with identical, unchanged fingerprints across
1,894 source files. The accumulated milestone executes and finishes **1198 / 1194 / 1198**
selected checks in default/server/highres, with only the existing opt-in
TestMapPopulationPrerequisiteProbe skip in each. All twelve frozen groups match
without changing expectations. Static checking also passes.

| Phase | Batch-driver seconds |
| --- | ---: |
| native-focused | 114.871 |
| native-default | 809.116 |
| native-server | 901.829 |
| native-highres | 738.743 |
| native-production | 465.107 |

The conversion removes **1,177 C lines / eighteen functions**, leaving **59,602
physical C lines in 82 production files**, zero reference C. The remaining file
prefix, including its existing boundary blank line, is preserved. Seventeen
Go-backed interfaces retain callback/caller compatibility; the private historical
reader is retired. No C algorithm is retained solely for testing.

Review points: the baseline deliberately fixes failed-child and rejected-parent
ownership and the stale-type adapter; general side-buffer cleanup and replaced
name-buffer policy remain outside this batch. Old coordinate read grouping is
preserved because stream checksums depend on call boundaries. Trigger uses the
actual pointer callback signature. Repeated full saves are not assumed byte-identical;
controlled serialization records, actual save/load artifacts and gameplay pixels
provide the corresponding evidence. See DECISIONS.md.
