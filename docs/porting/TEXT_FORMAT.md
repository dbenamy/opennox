# Text formatting, scalar strings and audio directory

## C baseline

Eight roots /668,876 captured cases pass in two independent default runs and
server/highres, with no skips and all static checks. Every complete capture
matches the original probe; all four source manifests agree. Only five tagged
files differ from qualified production2390c78d, whose production/ABI/full-suite/
headless gameplay/save-load evidence is reused. See
[text-format-c-qualification.json](text-format-c-qualification.json).

Coverage includes numeric padding/signedness, UTF-16 units and byte widening,
null/empty strings, bounded writes/guards/returns, float rounding and special
values; exhaustive wide/byte comparisons, saturated decimal parsing and copies;
and21 real audio catalog directory/file/missing/symlink/separator cases.

## Implementation decisions pending integration

Preserve the custom formatter rather than substitute fmt.Sprintf. Typed Go
arguments replace C varargs. Retain literal format processing and the existing
trade named-line second pass. Preserve recipient guards before formatting.
Temporary message text uses Go-owned storage; fixed UI buffers use their actual
capacity and retain NUL termination. In-capacity behavior is unchanged; old buffer
overflows have no defined result to preserve. The low-level bounded formatter
still returns full length and only emits NUL when it fits.

Only the live base10/no-end-pointer decimal API is ported. Preserve ASCII-only
case folding in the current C locale and32-bit signed saturation. The actual audio
catalog caller needs normalized directory presence; use ifs.Stat and retire
unreachable glob iteration instead of recreating an unused file-search subsystem.
Shipped production calls use the exact audio path; stored path behavior remains.

Float precision beyond1074 cannot affect the old31-byte result prefix: every
finite binary64 fraction terminates within1074 decimal places. Preserve signed
NaN and libc inf spelling. Unsupported formats/missing arguments panic; old calls
aborted or read undefined varargs. No production rounding-mode setter was found.

## Delegation and storage

GPT-6 Luna scalar implementation passed655,391 cases; primary handled review and
execution. Scalar caller migration used a guarded unrun installer. Primary caught
an unnecessary cgo import and then a mistakenly removed still-used alloc import
before any compile. This reinforces checking transformed code, not accepting the
helper's import-usage report. Cleanup audit found eight valid candidates but needed
repeated steering and supplied manually transcribed paths with corrections. Primary
created machine-readable audit and independently found ten obsolete root-test
archives /817,011,134bytes, verified symbols/hashes/mtimes and inactivity, then
removed them. No source, assets or captures were deleted. Future cleanup audits
must return machine-readable paths and a bounded search, even when df reports zero
(root reserved space may still permit a small audit file). Keep deletion execution
with primary. No subscription savings have been measured.

Native integration/qualification is pending. C remains1,219 lines in12 files.
