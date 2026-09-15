# Shared object renderer port

Scope: ten routines / 567 physical C lines in `GAME3_1.c` and `GAME2_3.c`:
shared object rendering, ghost opacity, shiny animation, beam colors / queue /
clipping / rasterization, and renderer clipping-state save/restore. Shared
host/client team lookup and the renderer-bounds setter remain dependencies.
No prerequisite production-C correction was needed for this batch.

## Owned tests

Twelve client tests contain nine capture groups (38,677 results) and three
independent probes. Fixtures own the actual drawable pool, network lookup, player
array and team IDs, named animation list, image handles, software renderer,
world lighting, scanline spans, sight polygon/intersections, mapped tables and
legacy globals. The clipping fixture owns the actual C current-RenderData pointer,
which may differ from the Go renderer's current state. No rendering, team lookup,
polygon clipping or allocation algorithm is replaced with a test stub.

Matrices cover material and enchantment priority, signed height and repeated
below-floor cropping, ordinary/ghost/semitransparent fades, frame wrap, every
fade age through the lifetime at twelve rates, local/remote/missing/observer
players, local-player presence, team relations, invisibility and see-invisible,
movement and elapsed-frame thresholds, all 32 directions and six occlusion-span
layouts, nested clipping saves and restores, raw flags and extreme coordinates,
shiny animation periods/counts and cached lookup, beam capacity/exhaustion/reset,
five sight polygons, real intersections, alpha, clipping and raster pixels.
Results include full render state, normalized drawable/cache state, image calls,
pixels, queue/saved-state words, return values and RNG indices. Independent
probes verify observer suppression versus visible pixels, actual lazy Ghost type
lookup and opacity, and reuse of the named ShinySpot animation cache.

The fade precision test pins its OS thread and checks the existing hosted x87
control word (53-bit precision, nearest rounding). Original C captures are
repeated in separate processes before expectations are frozen. Eight groups
matched c-f/c-g; the final precision group matched c-h/c-i. Raw captures and
qualification logs are local in `build/port-client-object-render`; hashes below
and a recoverable C baseline commit are tracked. No C algorithm is kept solely
as a test reference after conversion.

## Decisions and limits for later review

The existing scanline helper can accept short nonhorizontal paths even with
empty spans. The C renderer uses its acceptance result but ignores returned
horizontal bounds. Preserve this behavior; an early assertion that every empty
span rejects was too broad and was corrected before freezing captures.

Occlusion helpers deliberately panic in the server-only target. The first server
fixture run exposed that unsupported path; the occlusion matrix is now selected
only for client/highres builds. Shared renderer tests still run on the server
build. This is a fixture build constraint, with no change to production behavior.

Preserve early invisibility returns after cache/render-state writes, raw clipping
flags and first-save-wins nesting, unsigned coordinate arithmetic, signed height,
32-bit frame wrap, beam queue bytes across resets and named-animation caching.
Only four C entry points must remain for current C callers: beam append/reset
and clipping save/restore. Six internal rendering/ghost/shiny/beam entries can be
retired. The shared team lookup keeps its existing host/client routing; this batch
does not port all host-side team management or the separate monster renderer.

## Frozen captures

| Capture group | Results | SHA-256 |
| --- | ---: | --- |
| object-render-beam-queue | 3340 | `992c65aab7901000f68e229c0cea1104e9009b11a740c51153fdb34a20ec2a77` |
| object-render-beam-raster | 162 | `f82a5f22319c14b9e2b164c5d8b25ce82d03941c521d93bf5e47be9f1c2fe0fd` |
| object-render-clip-state | 75 | `575ba34078ad52f3bffbd22abda55c928c606cd21a98bea724979d26a71053ec` |
| object-render-ghost-fade | 6144 | `b5c83febbb7e0aadb9221326e67892bb22f1e40794b8097cb5a3c26da1ea4049` |
| object-render-material-crop | 5040 | `764d3fae9e5c55e8deb071e8af27ab98ea8295bf294a03dfb638c5bcc8e50737` |
| object-render-occlusion | 6912 | `b10776df02ca123de504afd8c6022350c482a3fb73a0181aec53b273f8b5816c` |
| object-render-player-invisibility | 10752 | `6d938e4401b4a7eebd9e7e80279c08ef8f8c9982359f7cdefa0f5b5109b8fc0d` |
| object-render-shiny-periods | 3072 | `941ab7b2b50e71a88fea225a990320f2d464b8be5044e58203f1adb36a4f3f64` |
| object-render-fade-precision | 3180 | `b7172234edc7ad017e466d0f6a79510d27456404d1533a9cd293ea278941b882` |

## C qualification

The final C focus passed all 12 tests; affected server/highres passed 73/74.
The accumulated C suite passed all 706 selected tests (one expected optional
skip) before the additional root-only precision test. Final affected checks
include that test; native accumulated qualification will include all 707.
Production C build passed in 4.238s and fresh headless gameplay
passed in 38.159s with null audio and reference comparison
enabled. All qualified selections started and completed. The earlier server
occlusion fixture failure is explained above and was not accepted as qualification.

Native conversion has not been applied. Production C remains 92,307 physical
lines in 101 files, with zero test-reference C.
