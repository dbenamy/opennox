# Text formatting, scalar strings and audio directory

Status: qualified native conversion. C falls 1,219→635 lines /12→8 files; 28 C
interfaces retired, zero reference C. C baseline: d443c297.

## Native qualification

Default/server/highres pass 375/373/375 affected roots with no skips, eight frozen
captures /668,876 cases, two native consumer contracts and static checks. Three
production builds/ABI pass; full-suite failures and package results exactly match
the known baseline. Fresh headless gameplay and explicit save/load pass. All final
source manifests agree and preflight/final default binary hashes match. See
[text-format-native-qualification.json](text-format-native-qualification.json).
The chronology below preserves findings and failed/interrupted attempts; only
native-default3, native-preflight2, native-server1, native-highres1 and
native-production supply final evidence. All jobs joined; drafts/installers are
consumed. No C algorithm is retained for these tests.

The server profile excludes two client-hover tests via their existing `!server`
build constraint; all selected tests executed and passed.

## C baseline

Eight roots /668,876 captured cases pass in two independent default runs and
server/highres, with no skips and all static checks. Every complete capture
matches the original probe; all four source manifests agree. Only five tagged
files differ from qualified production 2390c78d, whose production/ABI/full-suite/
headless gameplay/save-load evidence is reused. See
[text-format-c-qualification.json](text-format-c-qualification.json).

Coverage includes numeric padding/signedness, UTF-16 units and byte widening,
null/empty strings, bounded writes/guards/returns, float rounding and special
values; exhaustive wide/byte comparisons, saturated decimal parsing and copies;
and 21 real audio catalog directory/file/missing/symlink/separator cases.

## Implementation decisions

Preserve the custom formatter rather than substitute fmt.Sprintf. Typed Go
arguments replace C varargs. Retain literal format processing and the existing
trade named-line second pass. Preserve recipient guards before formatting.
Temporary message text uses Go-owned storage; fixed UI buffers use their actual
capacity and retain NUL termination. In-capacity behavior is unchanged; old buffer
overflows have no defined result to preserve. The low-level bounded formatter
still returns full length and only emits NUL when it fits.

Only the live base10/no-end-pointer decimal API is ported. Preserve ASCII-only
case folding in the current C locale and 32-bit signed saturation. The actual audio
catalog caller needs normalized directory presence; use ifs.Stat and retire
unreachable glob iteration instead of recreating an unused file-search subsystem.
Shipped production calls use the exact audio path; stored path behavior remains.

Float precision beyond 1074 cannot affect the old 31-byte result prefix: every
finite binary64 fraction terminates within 1074 decimal places. Preserve signed
NaN and libc inf spelling. Unsupported formats/missing arguments panic; old calls
aborted or read undefined varargs. No production rounding-mode setter was found.

## Delegation and storage

GPT-6 Luna scalar implementation passed 655,391 cases; primary handled review and
execution. Scalar caller migration used a guarded unrun installer. Primary caught
an unnecessary cgo import and then a mistakenly removed still-used alloc import
before any compile. This reinforces checking transformed code, not accepting the
helper's import-usage report. Cleanup audit found eight valid candidates but needed
repeated steering and supplied manually transcribed paths with corrections. Primary
created machine-readable audit and independently found ten obsolete root-test
archives /817,011,134 bytes, verified symbols/hashes/mtimes and inactivity, then
removed them. No source, assets or captures were deleted. Future cleanup audits
must return machine-readable paths and a bounded search, even when df reports zero
(root reserved space may still permit a small audit file). Keep deletion execution
with primary. No subscription savings have been measured.

The notes below record integration findings before final qualification.

## Native integration notes

First native compile exposed a missed live inline header adapter:
`client__gui__guicon.h` still formatted console text. The original shortened
symbol-search regex missed `nox_vsnwprintf`; exact retired-symbol search found it.
The adapter and its typed test bridge now call the native formatter. The existing
console contract includes4,096-character inputs; it and gameplay message tests
pass in native probe2. All eight new captures /668,876 cases match C byte for byte;
the native terminated-buffer/long-text/embedded-NUL/early-guard tests also pass.
Probe1 failed compile only and is joined; probe2 is joined PASS. Full qualification
is still pending. Physical C is635 /8files, a584-line reduction.

Completed repeated C captures now share verified canonical probe3 captures through
relative symlinks (241,343,004 logical bytes). Older completed message captures
also share verified identical files; audit lists exact canonical targets. Preserve
those targets. Logical byte totals do not imply physical savings when files
already have other hard links. Cleanup must check actual free space after mutation.

The broad default run passed375 normal roots but the runner rejected the selected
standalone TestMapPopulationPrerequisiteProbe skip. Its four cases are already
executed by TestMapPopulationPrerequisiteRegressions. The explicit affected-root
list now excludes only that standalone opt-in helper, retaining the regression
root; full zero-skip qualification is rerunning. No source changed for this fix.
An additional38 obsolete pre-format legacy-package cache archives were independently
symbol/hash/mtime-verified and removed after all jobs joined (1,622,664,172bytes).

Default2 passed375 roots/no skips and all contract/static gates. During review,
primary identified the tagged gameplay adapter materializing numeric ABI words as
Go pointers. Qualification was stopped/joined during preflight compilation, before
scenario execution. The adapter now decodes each argument's role before making a
pointer; the new test asserts integer extremes remain in word fields. Two tagged
files changed, production unchanged; final gates restart consistently rather than
mix source manifests. This was a review correction before a GC-dependent failure,
not an observed test mismatch. Native-default3/preflight2 are the final attempt.

Completed preflight asset copies were verified and deduplicated: 556,388,715 bytes.
Restore via `build/port-final-formatting/deduplicate-format-preflight.py --restore
text-format-native` with Python. Final save-run assets remain intact.
