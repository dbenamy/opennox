# Client effects

Converted **46 routines /2,289 physical C lines** to Go: five whole C files
(client__draw__fx, drawrays, lightning, plasma, glowdraw), GAME3_1 plasma setup
4BA670, lightning initialization 4BAB30, orbit update 4CA720, four curve helpers
4BE800–4BEDE0, and GAME3 spark/orb helpers 4B6880, 4B6970, 4B69F0 and 4B6B80.
Production C is **96,216 lines /142 files /zero reference C**, down 2,289 from
the qualified C baseline. Two explicit prerequisite repairs preceded conversion;
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

## C baseline and prerequisite repairs

The early three-test fixture is pushed as 69512647; the complete baseline is
pushed as **b554a327**. The C baseline passes 654 selected accumulated standard
root tests (373.911s), plus all 21 effects tests in server/highres (161.309s /
25.696s). Stationary orb and Hermite references from the early checkpoint remain
valid. All 18 final capture groups have identical repeated C results.

Independent checks caught two missing startup tables in fixture construction:
Hermite coefficients and integer distance. The fixture now loads bounded data
from the production embedded blobs and restores borrowed state. The distance
correction changed moving-orb/plasma outputs; the early 1,617-snapshot glow
reference and early ray-drawing hashes are withdrawn, not port oracles.

Two production C repairs are deliberate: skip missing plasma endpoints, and
initialize both particle target coordinates for coordinate chain lightning.
The latter originally left Y uninitialized, producing run-dependent particle
positions/RNG consumption in eight cases. A binding-equivalence contract covers
the repair. These fixes are distinct from the subsequent equivalent Go port.

## Native conversion and qualification

Seven legacy/client_effects Go files replace all 46 algorithms. Fourteen private
C entry points are retired; 32 retained exports serve real C callers or stored
callbacks and are exercised directly by the fixture. The orbit callback remains
C-callable. Internal plasma curve callbacks are typed Go calls; 4BEDE0 adapts its
remaining external C callback. The Go lightning initializer calls its native
helper. No old C algorithm is retained solely for tests.

All **13,344 captured results /18 groups**, plus independent contracts, match
unchanged references. Initial differences in spark clipping and curve/plasma
rounding were traced to C signedness and compiler storage boundaries, corrected
in Go, and documented in DECISIONS.md. Native-c passes all 21 focused tests.
Representative diagnostic PNGs were reviewed; exact framebuffer words remain
the oracle and historical sprite-color goldens were not regenerated.

Full accumulated standard/server/highres matrices pass with **654/653/654**
selected/executed/completed root tests (377.692s/444.741s/383.076s wall). Each has one intentionally
skipped optional TestMapPopulationPrerequisiteProbe; all actual regression
contracts pass. Accumulated coverage is 109,486 captured results /1,057 groups
plus contracts. Every selected suite uses the guarded runner.

All three production binaries build and pass ELF32/i386, SSE2 and C ABI audits:
14 retired symbols absent, 32 required exports present, no fixture helpers.
The asset-backed full suite matches all 1,553 known failure entries exactly and
all package outcomes (15 pass /3 fail /32 skip), with normal suite exit status.
Fresh unchanged headless warrior gameplay passes in **35.750s**, using Xvfb,
null audio, original assets and override disabled.

Local evidence: build/port-client-effects/qualification.json, baseline.json,
variant/full-suite logs and binary-verification.json. The reproducible accumulated
selection is tracked in accumulated-test-pattern.txt; see RECOVERY.md for commands.
Completed raw captures/binaries may be gzip archived with SHA256/restoration
manifests. The original assets and 7z remain untouched. Repaired source reference
is baseline-sections.c.txt; the older original-sections.c.txt is not the oracle.

Next connected batch: 27 drawable-update routines /1,069 C lines, reusing the
qualified owner. Its audit and fixture drafts are staged under
build/port-client-updates; no source from that batch was applied during this
qualification. Retarget the remaining Go energy-spark wrapper in that batch.
