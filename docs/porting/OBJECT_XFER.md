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
