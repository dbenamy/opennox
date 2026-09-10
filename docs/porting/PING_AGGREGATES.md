# Player-ping aggregates

Scope: 554290 minimum and 554300 signed average, both called only through Go
wrappers. Keep timing callback 554240 because it has other C callers.

267 datasets (11 explicit plus 256 deterministic generated), each runs four
isolated ABI/wrapper calls and a sequence of four calls: 2,136 aggregate results.
Verify raw uint32 result, exact callback trace and immutability of C-owned player
storage. A synthetic 32-slot active-player list uses the real first/next bridge.
Cover host-only/empty, sparse order, first-read zero/negative filtering, second
read zero/negative/different values, unsigned minimum, truncating signed division,
signed wrap in both directions, and callback sequence exhaustion.

The C baseline uses -fno-strict-overflow. Average oracle sums in int64 then wraps
to int32 before signed division, independent of the production accumulator loop.
Only active non-host players are visited; eligible players read timing twice.

Original-C 386 baseline passes all 2,136 aggregate results before conversion.
Production C before conversion: **141,180 physical lines**, 153 files, zero
reference C. Local artifacts: `build/port-ping-aggregate/`.

## Go conversion

Original-C baseline: `fea6ca7b`. The existing Go wrappers now call private native
helpers. Both aggregate C symbols and declarations are retired; 554240 remains
live for other C callers. Explicit int32 conversion preserves the signed timing
check and wrapping accumulator, and the public results remain uint32.

Production C: **141,126 physical lines (−54)** in 153 files; reference C: **0**.
All accumulated protection/network/waypoint/rules/spell-class/ping tests pass on
386 default, server and highres. All three production binaries build. The fresh
ping-aggregate-port headless scenario passes both preserved screenshot checks
with overrides disabled. Full suite was last repeated at the command-rule
milestone, with exactly the same 1,553 known failures; these two small leaves
used targeted accumulated tests, all builds and gameplay checks.
