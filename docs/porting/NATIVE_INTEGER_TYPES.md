# Native integer owner and caller types

Status: qualified. Production baseline `7e698a49`; integer baseline `308da9e7`.

Twenty production files use C only for integer types and header includes. Replace
`C.int`, `C.uint`, `C.short` and `C.uint32_t` with explicit `int32`, `uint32`,
`int16` and `uint32`. Keep signed narrowing, unsigned bit patterns, pointer-word
conversion and memory layout on the qualified Linux 386 target. Character,
floating-point, struct and actual libc boundaries are outside this batch.

The conversion covers 34 private typed functions and their callers in 37 files.
AST inspection found no non-call references to these functions. All 20 candidate
preambles contain includes/comments only: no C bodies or compiler directives.
Three calls use their existing Go owner directly: two modifier informational
messages and the mouse-mode setter. Preserve their original narrowing/data
addresses and ignored return behavior. Other arithmetic, branches, ordering,
allocation and ownership stay unchanged.

Existing independent contracts include exhaustive network low-word/high-bit
checks, signed boundary/handle wrap cases, pointer identities, list ownership,
message serialization, state changes and frozen captures. No mirror tests are
added solely to assert Go's integer widths. The final focused selection includes 317
roots; server omits exactly the two client hover roots and
client object-render occlusion, each guarded by `!server`.

Original qualification reuses the exact-source default full corpus for its first
observation and matching focused observations from server/high resolution. Missing
coverage and complete repeats passed in separate processes. A final caller audit
added 17 temporary-object/waypoint-append roots reached through shared fixture
adapters. Verified original binaries supplied their missing/repeat observations;
their source records were checked against committed `308da9e7`, while fixture
data was unchanged. This deliberately uses original compiled code, not the
converted Go workspace, as the reference. See [the initial baseline](native-integer-types-baseline.json)
and [the extension](native-integer-types-baseline-extension.json).
Every selected original root has two passing observations.

The first compile caught four fixture type mismatches. The handle snapshot/local
now uses int32; two still-C-typed waypoint fixture adapters explicitly convert
native return words to C.int. These are boundary conversions, not changed
assertions or captures. The corrected 37-root preflight and static check pass.

Converted qualification passes 317 default/high-resolution and 314 server roots,
with exact names and no failures/skips. Evidence combines the initial 300/297/300
runs with the 17-root transitive supplement on identical converted source/binaries.
Safe/static checks, all three production/ABI builds, exact known-suite comparison,
fresh final headless character creation/save/load/resume and all 1,654 original
asset hashes pass. The focused scope retains the just-qualified full default
corpus as original evidence; no full converted sweep or separate GUI preview was
needed for this representation-only change. Future failures/new dependencies
still require expanding qualification. Frozen expectations are unchanged.

| Metric | Before | After |
| --- | ---: | ---: |
| Selected project cgo files, client/server | 72 / 73 | 52 / 53 |
| Selected legacy C exports | 0 | 0 |
| Embedded production C bodies | 20 | 20 |
| Headers / physical lines | 157 / 2,731 | 157 / 2,731 |
| Standalone production / test-reference C lines | 0 / 0 | 0 / 0 |

Three newly empty fixture C imports are also removed. AST review records 77
changed and 349 unchanged functions, with none added/removed. All root contract
assertion files are byte-identical. External native-library dependencies remain
unchanged. See [qualification](native-integer-types-qualification.json) and
[dependency inventory](native-integer-types-inventory-after.json).

Primary works locally because Luna remains quota-blocked. Draft/review artifacts
and run records are under `build/port-native-integer-types/`. Seven superseded
`519ce712` production/safe/test binaries were deleted after qualified `7e698a49`
replacement/source/hash/host-use checks, reclaiming 387,088,384 allocated bytes.
Rebuild them from `519ce712`; source/logs and current binaries/assets remain.
Cleanup journal: `superseded-binaries-{approved.json,deleted.jsonl}`; script consumed.

Seven obsolete root/legacy cache archives predating `308da9e7` were removed after
hash/stat/host/compiler-use checks (340,516,864 bytes). The final scenario's 1,654
verified duplicate asset copies were removed (560,017,408 allocated bytes), with
original assets and saves/results retained. Its restore manifest is under
`build/baseline/runs/native-integer-types-save/`; all cleanup scripts are consumed.
