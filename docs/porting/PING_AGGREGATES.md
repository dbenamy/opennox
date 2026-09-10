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
