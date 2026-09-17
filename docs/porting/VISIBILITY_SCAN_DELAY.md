# Supplemental visibility scan-delay baseline

This test-only checkpoint extends `6e8dccf2` while retaining its original production
C. Native review identified missing direct coverage for the no-enemy scheduling
branch, so it was exercised in an isolated checkout before accepting the Go port.

The added 864 cases cover empty and populated real spatial indices, back-facing
enemies, distances through and beyond 1,000, quest/default/custom vision radii,
zero/30/60 frame rates and 32-bit frame wrap. Independent integer arithmetic checks
the interpolated deadline; qualified RNG state checks enforce exactly the expected
random draw. The stored nearest distance and scan timestamps are also checked.

The fixtures now give indexed objects real circular bounds. Previously the simple
owned objects had no shape, leaving bounds at the origin; that happened to lie
inside the old small test region but excluded objects from the new distant region.
The original twelve spatial-scan captures remain unchanged after that correction.
No production algorithm or golden was changed to accommodate it.

New group: `visibility-effects-scan-delay`, 864 records, SHA-256
`0db5f83c5cf17aaf3c9552fc3b924247ea22c628453e5a9617362b066355b50d`.
Run `go test -tags porttest -run '^TestVisibilityEffects(ScanDelay|SpatialScan)$'
-count=1 .` from `src` in the documented 386/SSE2/CGO environment. The supplemental
original-C evidence is default-target, repeated in separate processes; final native
qualification checks these expectations in default/server/highres.

Artifacts: build/port-visibility-effects/c-extra2-tests.jsonl and
c-extra-repeat-tests.jsonl, with corresponding captures. The isolated checkout is
build/port-visibility-effects/c-extra-tree. Its original-C production behavior is
covered by the preceding full baseline qualification; these changes affect tests
only. Native qualification is tracked in [VISIBILITY_EFFECTS.md](VISIBILITY_EFFECTS.md).
