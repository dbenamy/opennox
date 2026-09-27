# Native string boundaries

## Scope and decision

Production baseline is `81cf750f`; accepted original baseline `7034a7e4` is pushed.
The reviewed 11-file conversion is installed and fully qualified; artifacts are
under `build/port-string-boundaries/`. Remove the legacy string file's
C import: native int8/uint16 signatures, byte copying and terminated-string scan.
Update six fixture call sites to pass pointers through existing native adapters.

Centralize CString allocation in the existing allocator package, retaining raw
malloc/free normally and tracked Malloc/FreePtr in safe builds. The normal raw
entry still uses cgo's special C.malloc wrapper, preserving its process-fatal
failure. Safe retains its recoverable `cannot allocate` panic. This consolidates
allocation at the boundary that will be ported later; it does not claim that the
allocator is already C-free. Preserve the existing length-overflow check, every
input byte including embedded NUL, and the final terminator. Do not simplify the
browser's lone production CString round trip in this batch.

GoString scans/copies natively without adding allocator memlog observations.
Wide-string decoding and copy algorithms remain unchanged; their alias becomes
uint16. No allocation layout, external backend or callback behavior changes.

## Original preflight and selection

Four new private contracts pass on original default and safe builds: byte content,
NUL/bounds/copy independence, guarded bounded copying, UTF-16/malformed surrogate
handling, ownership/release and allocation-failure disposition. The failure test
uses a child process and an impossible 32-bit malloc size; no large Go string is
allocated. Normal exits2 via runtime fatal allocation failure; safe recovers and
exits with the fixture marker. Original expectations remain frozen for conversion.

The unscoped reference graph conflated same-named helpers across packages. A
package-aware graph still selects 2,484 candidate roots across 1,212 assertion
files because shared string/initialization paths have broad consumers. It includes
method references and wchar2_t signature/body uses. Profile discovery filters the
actual compiled names; the existing opt-in map-population diagnostic is excluded.

Repeat the new private contracts twice in default/server/highres/safe. Run each
broad original profile once, then the identical broad converted selections; no
need to repeat every established root a second time solely for this batch.
This retains repeated new original-path evidence and independent contracts while
avoiding duplicate full regression runs. Do not reuse differing source fingerprints.
Production qualification is reused from `81cf750f` for the test-only baseline; all
production/ABI/known-suite and fresh headless save/load gates run after conversion.
Run an early converted headless save/load after native preflight, before the broad
converted sweep; preserve reference screens/captures and repeat final integration.

Luna remains quota-limited. Primary owns the bounded draft, caller/type review and
qualification. No substitute model. The original baseline and corrected native conversion are accepted.

## Artifact headroom

Removed seven superseded libc-helper qualified executables after source/replacement
checks (`79b7e9e3` → `81cf750f`) and host-use verification: 387,010,560 allocated
bytes. Six transfer/duration original-baseline executables were also verified
against `66d3a57f` / `12a8ed98` and removed: 404,799,488 bytes. Retained commands
and source maps reproduce these builds; captures and logs remain. Journals are
`build/port-string-boundaries/libc-cleanup-*` and `historical-baselines-*`.
An update-baseline source mismatch stopped that candidate's removal; it and the
unresolved server-fixture candidate remain untouched. No source/assets were deleted.

Nine older modifier/audio/remaining-draw original-baseline executables were removed
after exact committed-source, replacement hashes and host-use checks: 605,700,096
allocated bytes. Rebuild `09f15464`, `9f6b2b46` and `9f91acda` with retained commands
and source maps; logs/captures remain. Evidence:
`build/port-string-boundaries/more-baselines-{approved.json,deleted.jsonl}`.

## Accepted original baseline

The exact 2,483 default, 2,472 server and 2,483 high-resolution selected root names
passed once each, with no skips and matching source/environment records. Four
private contracts passed twice in default/server/highres/safe. Only the two new
porttest files differ from qualified production `81cf750f`; its production gates
are reused for this original baseline only. See
[string-boundaries-baseline.json](string-boundaries-baseline.json).

## Native preflight correction

The first native preflight stopped at compilation: the map-catalog fixture passed
`internCStr`'s native `*int8` directly to `GoStringP`, which accepts `unsafe.Pointer`.
Add that explicit pointer conversion. No test had run; assertions and frozen
expectations are unchanged. Original attempt/source metadata are retained under
`build/port-string-boundaries/preflight-initial/`; the correction and initial/final
draft hashes are recorded in `preflight-correction.json` and the manifests.
Corrected native preflight/static and early fresh save/load pass. Broad converted
regressions and remaining production qualification also pass.

The preview's 1,654 verified duplicate assets were removed after completion:
559,878,144 allocated bytes. Original assets, saves/results and restore records
remain. Restore with
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/string-boundaries-preview-save/deduplicated-assets.json`.
Evidence: `build/port-string-boundaries/preview-cleanup/`.

Twelve completed placement/string original or qualification JSONL logs are now
losslessly compressed, reclaiming 522,993,664 allocated bytes. Committed log hashes
were checked before compression and against decompressed output; host-use checks
passed. Restore with `gzip -dk FILE.jsonl.gz`. Exact paths, hashes and commands:
`build/port-string-boundaries/log-archive.json`. Current converted logs remain live.

## Completed qualification

Converted 2,483 default, 2,472 server and 2,483 high-resolution roots match the
exact original selections without skips. Four private contracts pass twice in
all four profiles. Corrected preflight/static, early headless save/load, safe
build/static, all three production builds/ABI, exact known-suite failure/package
outcomes and a final fresh headless save/load/resume pass. All 1,654 original
asset hashes remain unchanged. The collector verifies the one-line fixture cast
correction, final draft/source hashes, and unchanged assertions/captures.

Selected production cgo files: **7/8 → 6/7 client/server**. Three project cgo
packages, 20 embedded production bodies, zero selected C exports and 157 headers /
2,731 physical lines remain. Standalone production/test-reference C remains
**0/0 lines**. External bindings are unchanged. Evidence: [qualification](string-boundaries-qualification.json)
and [inventory](string-boundaries-inventory-after.json).

The completed string final scenario's 1,654 verified duplicate assets
were removed, reclaiming 559,984,640 allocated bytes. Originals,
saves/results and the restore manifest remain. Restore with
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/string-boundaries-save/deduplicated-assets.json`.
Evidence: `build/port-string-boundaries/final-cleanup/`.
