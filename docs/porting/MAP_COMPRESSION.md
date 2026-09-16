# Complete NXZ map compressor

Status: **qualified, including final loop-bound correction**. C baseline **3ee3d454**, Go implementation **2699f7fc**,
cleanup separation **a2fc6ae5**, test memory accounting **fd976cf1**.
Starting size: **71,252 C lines /90 files /zero reference C**.
Final size: **69,342 /88 files /zero reference C**, removing **1,910 lines**.
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
The four capture groups contain 222 frozen records. Independent expanded inputs
and shipped compressed files supplement C hashes.
Repeat captures in separate processes and qualify default/server/highres.

## Audit notes and review items

The C match search uses five-byte rolling lookahead, including bytes beyond the
current 500,000-byte block. Preserve following file bytes for block-end hashes.
The final file tail can cause C to read a few bytes beyond its allocation during
hash maintenance. Those hashes are unused after the final block; the Go implementation uses
explicit zero padding, without changing valid encoded output expectations.

The existing empty-input allocation panic and unchecked file-size arithmetic are
replaced with explicit empty-file output and addressable-size checks, with
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

Qualified C size: **69,342 physical lines /88 files /zero reference C**, a reduction
of **1,910 lines**. Initial scope estimate overstated the shared helper by two
lines; the tracked physical-line counter is authoritative. All gates below passed.

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
the first cleanup-failed run copy is retained unchanged. The full corpus subsequently passed in all three targets.

## Accumulated-process memory issue

The first default accumulated sweep exhausted the 386 process address space
(about 3.99 GiB already in use) while allocating a 1.85 MiB server fixture in
TestProtectionValidateABI. It completed 826 selected tests; this is a failed,
incomplete sweep, not a qualification. Log: `milestone/default.jsonl`.
The fixture alone passes with a 768 MiB Go memory limit (0.318s package time).
Host memory remained available; this was the 32-bit process limit.

The test driver now defaults to `GOMEMLIMIT=768MiB`, preserves explicit overrides
and records effective runtime settings. Two driver tests cover default/override
propagation. Full reruns passed in isolated target output directories with
GOMAXPROCS=1 and a 1,800s per-package timeout for concurrent target execution.
All selected cases remain required; the lone opt-in prerequisite skip is unchanged.
No game or fixture source changed.

## Performance observation for review

The all-50-map compression test took 0.17s in the final C default capture and
0.40–0.45s in the initial Go focused runs. Synthetic mixed-match checks took
0.87s C and 0.96–1.02s Go. These include file I/O and are not isolated benchmarks;
do not infer an overall port-speed factor. The aggregate real-map difference is
under 0.3s in this VM, but the writer is measurably slower in these observations.
Preserve this evidence for later profiling rather than claiming performance parity.

## Final qualification

All 15 selected codec roots pass in each target without skips. Three production
builds pass ELF32/i386/SSE2/CGO and retired-symbol audits. The full asset suite
retains exactly 1,553 failure entries and 15 pass /3 fail /32 skip packages. Both
forced-map gameplay modes match all 26 frames and original expanded map bytes.

The completed accumulated corpus uses GOMAXPROCS=1 and GOMEMLIMIT=768MiB:

| Target | Selected/completed tests | Root-package tests | Driver seconds |
| --- | ---: | ---: | ---: |
| default | 1064 | 1047 | 608.950 |
| server | 1060 | 1043 | 772.499 |
| highres | 1064 | 1047 | 637.666 |

Each target has only the existing opt-in TestMapPopulationPrerequisiteProbe skip.
Evidence: `build/port-map-compression/milestone-bounded-{default,server,highres}`.
Targets ran concurrently in isolated directories; their times are not additive.
Source fingerprints stayed unchanged and every reader joined before the next
source batch. Seventeen tooling tests pass. Retain the earlier failed cleanup and
address-space attempts; neither was relabeled as a pass.

The complete codec package now has no C source or cgo import. Existing live map
wrappers use its Go API. No C algorithm remains solely as a test reference.

## Final loop-bound review

After the complete corpus passed and every reader joined, review found that a
fixed 500,000 increment could overflow int after the last block of a near-limit
386 file. Advance by the actual consumed chunk instead. This reaches srcSz exactly
and leaves all normal compressed bytes unchanged. Re-run all affected codec tests
and production builds/ABI, known suite and both fresh scenarios. The accumulated
corpus above qualifies 2699f7fc; its scope outside the codec is unchanged by this
small wrapper correction, so a second complete sweep is not scheduled.
All final affected and production gates pass in
`build/port-map-compression/native-loop-final`. Final scenarios are
`map-encode-native-final` and `map-encode-flat-native-final`; both regenerate the
map exactly and match every reference frame. Source fingerprints remain unchanged
through this final run, and all readers joined. This checkpoint commits the final
wrapper correction with its qualification evidence.
