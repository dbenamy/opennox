# Complete NXZ map compressor

Status: C baseline qualified in all three targets and an independent repeat.
Go implementation passes focused frozen and independent contracts; production
qualification pending. C baseline commit: **3ee3d454**. Starting size: **71,252 C lines /90 files /zero reference C**.
Scope: `legacy/cnxz/nxz_comp.c` (1,892 lines) and `nxz_common.c` (18 lines).

## Scope and contracts

Complete the codec's Go ownership: sliding dictionary, match selection, adaptive
symbol codes, bit packing, block boundaries and file wrapper. The only production
C entry to compression is the existing Go map wrapper in `legacy/maps.go`, which
calls CompressFile. Compressor-internal symbols have no other callers.

Reuse all decoder contracts and 115 synthetic frozen compressed outputs; add 48
mixed-copy cases around symbol rebuild, 64 KiB window, 521-byte match and
500,000-byte file chunk boundaries. Compare compression with all 50 shipped NXZ
files, using actual filenames. Verify source preservation, destination truncation,
same-path compression/decompression and destination preservation on input errors.
Independent expanded inputs and shipped compressed files supplement C hashes.
Repeat captures in separate processes and qualify default/server/highres.

## Audit notes and review items

The C match search uses five-byte rolling lookahead, including bytes beyond the
current 500,000-byte block. Preserve following file bytes for block-end hashes.
The final file tail can cause C to read a few bytes beyond its allocation during
hash maintenance. Those hashes are unused after the final block; the Go draft uses
explicit zero padding, without changing valid encoded output expectations.

The existing empty-input allocation panic and unchecked file-size arithmetic will
be replaced with explicit empty-file output and addressable-size checks, with
independent contracts. These are bounded, reversible corrections under the user's
standing authorization. Do not regenerate valid-file goldens to absorb differences.

## Evidence

Manifest: `map-encode-batch.json`; affected selection: `map-encode-tests.txt`.
Local initial C checks: `build/port-map-compression/c-focused.log`,
`c-matches.log`. Frozen match capture SHA-256:
`c8c09c83ad7cd5b40da7986ececb7f59fef83ca902d91a1491fcf659bf19728f`.
C gate: `build/port-map-compression/c-qualified/result.json`, all steps passed.
The first native behavioral run passed all frozen expectations (5.071s driver).
Additional native contracts pass (5.520s): empty/sparse-oversized input, every
length 4–521 against 19 distance boundaries, and 12 small block sizes down to one
byte, including repeated table rebuilds. No goldens changed.

Working C size: **69,342 physical lines /88 files /zero reference C**, a reduction
of **1,910 lines**. Initial scope estimate overstated the shared helper by two
lines; the tracked physical-line counter is authoritative. Full gates pending.

## Qualification harness issue

All three production builds/ABI and the exact 1,553 known failure entries passed.
The first normal replay passed its map-byte check and reference frames, then an
optional ignored asset-deduplication helper rejected the new run name. The parent
runner correctly recorded the nonzero process exit; this was cleanup failure, not
a game mismatch. Keep the original attempt and full run copy.

Remove implicit cleanup from the tracked scenario runner. Perform verified local
deduplication separately after validation. Reuse unchanged-source build/ABI/suite
evidence through the existing guarded mechanism, and run both scenarios afresh.
No game/test source or valid-file expectations changed during these checks.

Both fresh scenarios now pass: `map-encode-native-qualified` and
`map-encode-flat-native-qualified`. Each regenerated one warrior map exactly and
matched its reference frames (12 normal, 14 flat-floor). Final integration report:
`build/port-map-compression/native-final/production/production.json`. The guarded
reuse check accepted the original three builds and full suite without rebuilding.
Successful copies were deduplicated separately with verified restoration manifests;
the first cleanup-failed run copy is retained unchanged. Full corpus still running.

## Accumulated-process memory issue

The first default accumulated sweep exhausted the 386 process address space
(about 3.99 GiB already in use) while allocating a 1.85 MiB server fixture in
TestProtectionValidateABI. It completed 826 selected tests; this is a failed,
incomplete sweep, not a qualification. Log: `milestone/default.jsonl`.
The fixture alone passes with a 768 MiB Go memory limit (0.318s package time).
Host memory remained available; this was the 32-bit process limit.

The test driver now defaults to `GOMEMLIMIT=768MiB`, preserves explicit overrides
and records effective runtime settings. Two driver tests cover default/override
propagation. Full reruns are pending in isolated target output directories with
GOMAXPROCS=1 and a 1,800s per-package timeout for concurrent target execution.
All selected cases remain required; the lone opt-in prerequisite skip is unchanged.
No game or fixture source changed.
