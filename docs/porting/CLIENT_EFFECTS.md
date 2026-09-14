# Client effects

Connected scope: **46 routines /2,289 physical C section lines**, all still C.
Five whole files (client__draw__fx, drawrays, lightning, plasma, glowdraw),
GAME3_1 plasma setup 4BA670, lightning initialization 4BAB30, orbit update 4CA720,
four curve helpers 4BE800–4BEDE0, and GAME3 spark/orb helpers 4B6880, 4B6970,
4B69F0 and 4B6B80. Current production C: **98,505 lines /147 files /zero
reference C**. Two explicit prerequisite repairs precede the Go conversion;
see [DECISIONS.md](DECISIONS.md).

## Reference fixture

The asset-independent owner uses production drawable allocation, type lookup,
list/spatial ownership, construction/deletion, renderer, viewport, RenderData,
and a 96×96 RGB5551 framebuffer. It starts no display. It saves and restores
borrowed globals, initializes bounded production coefficient/trigonometry/
distance tables, and releases only owned resources. Named C globals and mapped
historical addresses are separate storage and must remain so through this port.
Synthetic types have no unrelated draw callbacks or asset dependencies.

Captured state includes raw framebuffer hashes, complete RenderData and
512-byte C drawable prefixes, normalized list/cache/callback pointers, scratch
arrays, both RNG indices, factory calls, allocation failures and deletion order.
Screen-particle tests use the production pool at capacities 0, 1, 7 and 128,
including disabled allocation and tail reuse. No C algorithm is copied for tests.
Independent assertions complement repeated C hashes:

| Group | Recorded results | Additional checks |
| --- | ---: | --- |
| Orb probe | 1 frame | Nonempty pixels, radius and counter |
| Orb lifetime | 6,070 snapshots /1,056 inputs | Clip equality, byte wrap, shrink/delete and exact mutation |
| Hermite callbacks | 462 cases | Guards, continuity, constant curve, endpoints and midpoint |
| Spark bounce | 1,920 transitions /480 inputs | Signed arithmetic, parity/frame wrap, exact mutation, unchanged RNG |
| Particle construction | 224 cases | Seven producers, failure ordering, positions, callback and ownership |
| Screen particles | 112 cases | Pool/link/tail reuse, colors/velocity, 100 attempts and 400 RNG calls |
| Particle callbacks | 64 cases | Coordinate copy before external mutation, record guards |
| Glow/sparks | 1,716 snapshots /608 inputs | Nineteen paths, TTL, colors, clipping, motion, child failure and RNG |
| Lightning setup | 16 snapshots | Eight prior-cache masks, repeated initialization |
| Lightning recursion | 192 cases | Signed packed geometry, exact recursion/RNG count and unwind |
| Lightning passes | 176 cases | Length thresholds, modes and pass combinations |
| Ray dispatch | 396 cases | Packet guards, cache capacity, exact lookup/append/return and RNG |
| Orbit updates | 224 cases | Bounds/age expiry and unchanged RNG |
| Curve raster | 560 cases | Guards, empty steps, coefficient rounding, settings and linear endpoint |
| Plasma segments | 126 cases | Cursor/start/end, frame wrap, viewport transform and emission |
| Plasma geometry | 720 snapshots /120 inputs | Distance/angle thresholds, segment caps, first points and phase wrap |
| Ray drawing | 317 snapshots /160 inputs | Coordinate/static/dynamic/missing bindings, pause, mouse and expiry |
| Stored orbit callback | 60 snapshots | Invoke actual stored C callback, lifetime/frame wrap and RNG |

Separate contracts check moving-orb arrival/distance, plasma missing endpoints,
and equivalent chain-lightning particle endpoints across binding modes.

## Baseline corrections and current status

The early three-test fixture is pushed as **69512647**. Its stationary-orb probe
passes standard/server/highres; orb lifetime and Hermite callback references
remain valid. The full 21-test effects baseline is now qualified.

Independent checks caught two missing startup tables in fixture construction:
Hermite coefficients and the integer-distance lookup. Both now use bounded data
from the production embedded blobs and restore borrowed state. The distance
table correction changes moving-orb and plasma outputs. Earlier glow reference
(1,617 snapshots) and early ray-drawing hashes are withdrawn, not port oracles.
Corrected glow (1,716 snapshots) and plasma geometry match repeated captures.

Two production C fixes are intentional: missing plasma endpoints now skip
rendering, and coordinate chain lightning now initializes both particle target
coordinates. The latter originally left Y uninitialized, producing different
particle/RNG results in eight cases between runs. The endpoint contract compares
coordinate and object bindings for identical positions. These repairs are
separate from equivalence-preserving Go conversion.

All corrected references are locked after identical repeats. The guarded full
accumulated standard baseline passes all 654 selected/executed/completed root
tests (373.911s wall). Focused server and highres each pass all 21 effects tests
(161.309s /25.696s). Metadata: build/port-client-effects/baseline.json.
No effects algorithms have moved to Go yet. The complete native draft is staged
locally as client_effects*.go.stage and awaits integration and qualification.
Diagnostic PNGs in c-diagnostics were visually reviewed (curve, energy bolt,
lightning, orb, plasma); exact framebuffer words remain the oracle. Historical
sprite-color goldens are not regenerated; preserve the current renderer policy.

## Conversion and qualification

The provisional ABI audit retains 32 entry points and retires 14. Twenty-one
draw callbacks have Go address references. Retain 4CA720 while stored update
pointers use the C ABI. Plasma's internal synchronous callback can become a
typed Go callback. Curve evaluator 4BEDE0 has an outside GAME2.c caller and
keeps a C adapter for external callbacks. Retarget the sole Go caller of
lightning initialization 4BAB30. Recheck callers and prototypes before deletion.

Once the baseline is qualified, commit/push it, then convert all 46 routines.
Run focused comparisons during implementation. Because this batch adds shared
client/renderer ownership, complete the full accumulated native matrix under
all three tags, all production builds/ELF32/SSE2/ABI checks, exact known-failure
full-suite comparison, and a fresh unchanged headless gameplay replay.
Use tools/porting/run_tests.py for selected matrices; it rejects zero selection
and requires every selected root test to run and complete. Do not edit source
while a build/test runs. Update C LOC, decisions, this log and handoff; commit,
push and continue under the user's standing authorization.

Local evidence is under build/port-client-effects. baseline-sections.c.txt is
the repaired C reference; original-sections.c.txt predates repairs/scope expansion.
Completed captures may be gzip archived with verified SHA256/restoration entries
in completed-artifact-archives.json. Preserve the original extracted assets and
7z. Raw logs/captures remain local; tracked hashes and fixture code recover the
baseline from Git.
