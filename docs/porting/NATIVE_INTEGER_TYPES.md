# Native integer owner and caller types

Status: original baseline qualified on `7e698a49`; draft not installed.

Twenty production files use C only for integer types and header includes. Replace
`C.int`, `C.uint`, `C.short` and `C.uint32_t` with explicit `int32`, `uint32`,
`int16` and `uint32`. Keep signed narrowing, unsigned bit patterns, pointer-word
conversion and memory layout on the qualified Linux 386 target. Character,
floating-point, struct and actual libc boundaries are outside this batch.

The draft covers 34 private typed functions and their callers, initially 37 files.
AST inspection found no non-call references to these functions. All 20 candidate
preambles contain includes/comments only: no C bodies or compiler directives.
Three calls use their existing Go owner directly: two modifier informational
messages and the mouse-mode setter. Preserve their original narrowing/data
addresses and ignored return behavior. Other arithmetic, branches, ordering,
allocation and ownership stay unchanged.

Existing independent contracts include exhaustive network low-word/high-bit
checks, signed boundary/handle wrap cases, pointer identities, list ownership,
message serialization, state changes and frozen captures. No mirror tests are
added solely to assert Go's integer widths. The focused selection includes 300
roots from 140 test files; server omits exactly the two client hover roots and
client object-render occlusion, each guarded by `!server`.

Original qualification reuses the exact-source default full corpus for the first
300-root observation, and 30/32 previously qualified focused roots in server/high
resolution. New runs cover the missing 267/268 roots and then the complete
selection once/profile. Every selected root therefore has two original
observations with matching source, verified binaries and runtime environment.
All pass without skips. See [the baseline](native-integer-types-baseline.json).
Production baseline evidence is reused only because production source is unchanged.

After conversion, run the 37-root integer/handle/pointer/message preflight and
static check, all three complete focused profiles, safe build, fresh production
builds/ABI checks, exact known-suite comparison, final headless character
creation/save/load/resume, asset integrity and C/cgo counts. The just-qualified
full default corpus is retained as baseline evidence; this bounded representation
change uses complete caller/owner coverage rather than another full sweep unless
a failure or newly found dependency broadens the scope. A separate GUI preview
adds little for this type-only change; the fresh final gameplay scenario remains.
Frozen expectations are not regenerated.

Primary works locally because Luna remains quota-blocked. Draft/review artifacts
and run records are under `build/port-native-integer-types/`. Seven superseded
`519ce712` production/safe/test binaries were deleted after qualified `7e698a49`
replacement/source/hash/host-use checks, reclaiming 387,088,384 allocated bytes.
Rebuild them from `519ce712`; source/logs and current binaries/assets remain.
Cleanup journal: `superseded-binaries-{approved.json,deleted.jsonl}`; script consumed.
