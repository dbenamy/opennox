# Numeric fixture processor state

## Scope

Replace eight embedded C helpers in five numeric fixtures with a porttest-only,
cgo-free 386 assembly leaf. Keep the existing thread pinning, control-word masks,
world-numeric setup/restoration, numeric inputs and frozen assertions/captures.
Room and painting fixtures retain their aligned allocation/abort calls for a
later ownership batch. Production implementation and external backends are unchanged.

The proposed conversion reduces fixture C imports from 10 to 7. Production stays
at 4 client/highres files, 5 server files, 2 direct project cgo packages, zero
embedded C bodies and zero C exports in those profiles. Eleven safe-only
allocator/memory exports also remain and need retirement in the ownership phase.
Standalone production/test C remain zero;
157 headers / 2,731 physical lines remain. These are separate measures of scope.

## Baseline correction

Safe-mode checks on unchanged `bafd3e21` failed with `incorrect free` for both
`TestMapRoomsRounding` and `TestMapPaintingCoordinates`: `C.aligned_alloc` inputs
bypass the safe allocation tracker, but fixture cleanup used tracked `legacyFree`.
Before freezing the numeric baseline, mark these private fixture regions as raw
and release them with `alloc.RawFree`. Engine-created records retain their existing
tracked cleanup. Both original-C safe checks then passed. Alignment, record bytes,
production code and frozen expectations are unchanged. This reversible correction
is a review item, not a change to allocator semantics.

The first full safe selection also exposed a raw-blob redirect target in the world
numeric fixture accessed through the live-variable accessor. Its historical address
remains in the extraction registry, although no engine user remains. Use explicit
backing-blob access for this fixture-owned storage, preserving the exact address,
bytes and restoration. Keep runtime memory checks enabled. The preceding normal
owner/focused runs passed, but fresh qualification uses the corrected exact source.

## Behavior review

The qualified VM's i386 libc saves MXCSR in an overlapping `fenv_t` field.
Its `fesetenv` restores control/exception bits while preserving current x87 TOP
and stack tags, and clears instruction/data diagnostic fields. A blind reload of
the saved x87 environment would differ. The native leaf uses its own 32-byte
layout, preserves those masks and clears those same fields. This is qualified
Linux 386 behavior; it makes no portability claim about other libc implementations.

An isolated C/native probe passed 1,536 combinations in each of three fresh
processes, comparing control-word read/write and complete setup/restoration
snapshots. It varies x87 precision/rounding, MXCSR rounding/FTZ, sticky exception
flags, current TOP/tags and diagnostic fields. Artificial state is installed,
exercised, captured and reset within a single call before ordinary Go runtime work.
An earlier probe left that state live across runtime calls and produced varying
mismatches; those attempts are diagnostic evidence only. Probe-only entrypoints
are excluded from the installed leaf.

Go 1.26's assembler operand classification rejects ordinary `FLDENV (AX)`.
The leaf encodes that standard instruction as its two bytes, with an explanatory
comment. Disassembly confirmed the intended instructions; the isolated leaf passed
assembly declaration checks. The actual project leaf must also compile and pass
vet with cgo disabled. No toolchain change is required.

Caller review includes receiver-method consumers: gameplay-report fixtures inspect
the control word through `(*portTestShopPools).gameplayReportsAction`. Adding their
typed owners expanded the candidate selection to 215 roots. The optional
`TestMapPopulationPrerequisiteProbe` is excluded; its four cases are covered by
selected `TestMapPopulationPrerequisiteRegressions`. All three profiles contain the
remaining 214 roots. Fourteen focused numeric roots and eight safe roots supplement
the owner sweep. An AST comparison identifies eight changed existing Go functions,
36 unchanged functions and nine added leaf functions/declarations; assembly is
reviewed separately.

## Qualification and recovery

The corrected original baseline passed 214 owners and 14 focused repeats in each
profile, plus eight safe roots. No numeric conversion is installed.
The preceding tile-grid production qualification is reusable for the original
baseline because only fixture cleanup changed; all production gates will rerun
after conversion. Acceptance requires exact discovered/run/pass name sets,
unchanged source fingerprints (including assembly), safe ownership checks, original
asset hashes and unchanged frozen expectations.

Ignored working evidence: `build/port-fixture-fpu/`, including processor/library
hashes, disassembly, probe records, rejected probe diagnostics, ownership failure
and correction results, original-source copies and the uninstalled draft.
Committed expectations and the original baseline revision recover behavior if
those local artifacts are lost. Completed scripts are single-use.

Primary performed the work; Luna remains unavailable because its quota is exhausted.
No substitute helper was used.
