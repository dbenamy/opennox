# Placement trigonometry

## Scope

Production baseline: `79b7e9e3`. This batch replaces the sin/cos calls in generator
radial placement and inventory random placement. The latter also serves
teleporting, flag drops, AI, scripts and server placement. Original baseline
`9dcf1b9a` is committed/pushed; the installed conversion is fully qualified.

The conversion removes two production C imports and one fixture C import. It
preserves the original double sum passed to cos, the stored float32 angle passed
to sin, radius updates, RNG consumption, owner calls and coordinate arithmetic.
No allocator, callback, layout or external native-library change is included.

## Rounding decision

A direct replacement with Go's standard math functions is rejected: an isolated
386 comparison found 132,938 different double results among 524,288 sin/cos values.
Those differences reached final float32 coordinates in two of 16,777,216 ordinary
coordinate cases and 163,231 of 4,194,304 cancellation stress cases. Both ordinary
differences were in the tested game-coordinate range. These are deterministic
probe counts, not estimates of gameplay frequency.

Instead, the conversion adapts the small-argument IBM/glibc polynomial and table
implementation to Go. Provenance is pinned in
[placement-trig-upstream.json](placement-trig-upstream.json); notices and the
[LGPL-2.1 text](../licenses/GLIBC-MATH-LGPL-2.1.txt) accompany the adaptation.
The pinned upstream snapshot is not claimed to be the exact distribution source:
its native adaptation is compared directly against this VM's installed 32-bit
libc. All 440 table entries were checked by independently decoding both upstream
endian arrays, then matching the compiled Go table against that raw-bit hash.
The exact native draft also passes all 70 frozen boundary cases in a standalone
386 build with cgo disabled.

The two engine owners start from the finite RNG table's angle in [-pi, pi], add
1.8849558, store float32, and perform at most 64 attempts (32 for the generator).
Every reachable trig argument has magnitude below 128. The native prototype
matches **all 524,288 trig results for every one of the 4,096 RNG starting values
across 64 attempts**, repeated in separate processes. The combined raw-double
SHA256 is `0e09c10d0957707d8160e4ab4c5fbec4f008f77c428a66bba4a762418def66e9`.
It also matches the broader random-angle/cancellation probe without differences.

This is a private placement compatibility helper, not a general libm replacement.
The helper uses standard Go math beyond the upstream small-argument reduction
range; those inputs are unreachable from these owners and are not claimed to be
libc-bit-compatible. Requalify the domain before expanding its callers. Keep
binary64 evaluation boundaries and the qualified 386/SSE2 target.

## Contracts and review

New original-path contracts pass for the exhaustive RNG-angle hash and 70 frozen
boundary cases around signed zero, subnormals, polynomial transitions, quadrant
changes and the qualified range endpoints. Expectations were captured twice from
libc before production conversion. The test-only observers use the same libc
entrypoints as the production owners; the conversion switches them to the
native helpers and removes their C import. No C algorithm remains solely as a
reference implementation.

The source review identifies four changed, nine unchanged and seven added
functions across five files. Existing owner assertions and goldens remain fixed.
The caller graph plus explicit owner families selects 992 candidate root names
before profile filtering. Private numerical contracts run separately in
`./legacy`; full asset/runtime environment and exact-name/source checks apply.

Luna remains quota-limited. The primary performed the bounded prototype, capture
work and review without substituting a model. Local probes, drafts, manifests and
runtime results are under `build/port-placement-trig/`.

## Artifact headroom

Seven superseded native-layout test/safe/production executables were removed after
source fingerprints matched `6e9681f2`, qualified replacements matched `79b7e9e3`,
and artifact hashes/stat/host-use checks passed. Reclaimed 387,059,712 allocated
bytes. All current libc-batch qualified binaries remain.
Plan/journal: `build/port-placement-trig/layout-cleanup-{approved.json,deleted.jsonl}`.

## Accepted original baseline

All 992 default, 989 server and 992 high-resolution selected roots passed twice,
with exact discovered-name sets, no skips and matching source fingerprints.
Both private contracts passed twice in default/server/highres/safe. At that
baseline, production source was identical to `79b7e9e3`; its production
qualification was reused for the test-only baseline. The completed conversion qualification is recorded below.
See [placement-trig-baseline.json](placement-trig-baseline.json).

Six older original-baseline executables were subsequently removed after matching
source revisions `ca1d9fdf` and `e1926909`, hashes, current accepted replacements
and host-use checks: 400,789,504 allocated bytes. Rebuild those revisions using
retained profile commands/source maps; captures and logs remain. The first cleanup
attempt stopped before deletion because one historical result path differed;
verification now uses its committed accepted baseline record.
Plan/journal: `build/port-placement-trig/old-baselines-{approved.json,deleted.jsonl}`.
Four unused managed Go cache archives were also removed after hash/stat and
host checks with all build/test jobs joined: 130,002,944 allocated bytes.
Rebuild normally; records: `build/port-placement-trig/cache-headroom-*`.

A further 15 original-baseline executables from collision/death/create-init/damage/
item batches were removed after exact source/hash/replacement and host-use checks:
1,013,202,944 allocated bytes. Rebuild `3ad2c7d5`, `0f5e68ea`, `f951e4a2`,
`555c7d1a` and `e546fed3` using retained commands and source maps. All original
captures/logs and current placement baseline binaries remain. Evidence:
`build/port-placement-trig/historical-baselines-{approved.json,deleted.jsonl}`.

## Completed conversion qualification

Native preflight/static checks and safe build/static checks pass. Converted
992 default, 989 server and 992 high-resolution roots match the exact original
name sets without skips. Both private numerical contracts pass twice in all four
profiles. All three production builds/ABI, exact known-suite failure and package
outcomes, and a fresh headless save/load/resume pass. The collector verifies
installed/source/binary hashes and that existing assertions and frozen data are
unchanged. All 1,654 original assets retain their hashes.

Selected production cgo files: **9/10 → 7/8 client/server**. Three project cgo
packages, 20 embedded bodies, zero selected exports and 157 headers/2,731 lines
remain; external bindings are unchanged. Standalone production/test-reference C
remain **0/0 lines**. Evidence: [qualification](placement-trig-qualification.json)
and [inventory](placement-trig-inventory-after.json).

After qualification, 1,654 verified scenario asset duplicates were removed,
reclaiming 559,898,624 allocated bytes. Originals, saves/results and a restore
manifest remain. Restore with
`python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/placement-trig-save/deduplicated-assets.json`.
Evidence: `build/port-placement-trig/final-cleanup/`.
