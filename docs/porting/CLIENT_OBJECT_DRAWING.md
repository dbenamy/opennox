# Object drawing port

The connected batch covers 18 drawing callbacks (door, arrows/tails, glyph,
summon, equipment/base/flag, generator, shield, pressure plate, trigger and powder),
the door parser, and eight connected material/team/countdown helpers: 27 routines.
The corrected C baseline contains 912 lines in this scope, including 12 whole
files. Shared drawObject and draw-data cleanup remain production C dependencies.
This does not qualify all of drawObject's player/team/invisibility/clipping paths.

## Corrected C prerequisites

Three reversible corrections precede the conversion and define its baseline:

- Failed optional arrow-tail allocation draws the main arrow, preserves its old
  tail anchor for retry, and avoids nil endpoint/lifetime/list writes. Both arrow
  variants have allocation-failure and successful-retry contracts.
- Glyph paths without distance-derived alpha explicitly use 255. Previously they
  used the drawable pointer's low byte, making pixels depend on allocation address.
  The contract compares real rendered pixels with an opaque sprite across padding.
- Generator Random selects 0 through count-1, matching the inclusive RNG contract
  and declared image count. One-frame state arrays have distinguishable adjacent
  sentinel images; 128 seeds/state verify selected image and RNG consumption.

These changes add seven C lines (93,775 to 93,782;118 files;zero test-reference C).
Review later: optional-tail omission/retry and opaque-default glyph behavior are
intentional corrections, not behavior changes introduced by the Go translation.
Generator OneShot is valid during the destruction countdown or completed state;
its original provisional index is overwritten by that lifecycle. The actual
network update initializes this countdown on entering flag0x400. Completion sets
0x800, clamps the final image, and clears lighting flags. Tests exercise these
states instead of constructing an invalid free-running OneShot generator.

## Owned baseline

14 focused tests comprise ten matrices (13,298 results) and four independent
contracts/probes. The owner uses actual drawable pool/list/spatial indexing,
modifier definitions/team arrays, C-backed image handles and data, material-indexed
RLE rendering, real world-light interpolation and an isolated default font.
Summon children use the actual unlinked constructor and raw-delete path. The
shared snapshot accepts an explicit unlinked list and validates its count and
absence from the spatial index; existing snapshots retain their original behavior.

The production DrawImageAt wrapper replaces its low-level image hook. The fixture
therefore delegates the real wrapper and observes its resulting metadata/image at
the interface boundary. It does not replace rasterization. Capture rows include
pixels, full drawable/render state, ordered images, named-image loads, actual C
caches/colors, RNG indices, spawn/deletion history and explicitly normalized owned
pointers. Flags/countdowns, parent animation restoration, child dispatch/deletion,
network-owner following and allocation retry also have independent assertions.

C captures c-e/c-f repeat byte-for-byte. Earlier fixture corrections were renderer
hook placement, exact door byte offsets, static base image data and Go integer
casts; none required additional production changes. Qualification manifests and
raw captures are local under build/port-client-draw-objects. Only validated hashes
and expectations are committed; no C algorithm is retained solely for tests.

| Capture group | Results | SHA-256 |
| --- | ---: | --- |
| object-drawing-arrows | 1458 | `c7cedab0fa2d041920816e900230bb2e5df0fb0d04af3a18313fa6400b4576c1` |
| object-drawing-door-parser | 15 | `363fdd6579402ea968fd231b79e37180f6b0093a0522cd8a57c50914dc1a3bfe` |
| object-drawing-door | 432 | `cad98de7e6c05e4f3cdd9a42c5cb763bee1c30683a89dfb2dd60d48cb346c21f` |
| object-drawing-equipment | 1728 | `93efdb84ec27ae3124beb33f0006089b31a04e1e7cddc4ff01568c710a177879` |
| object-drawing-generator | 7488 | `22a5ccdafb27a8fe752425ed4b60276767b40a11aaa821febe2b4e0cd07d1eee` |
| object-drawing-geometry | 1260 | `9a2075112d571e8505b385373147974d0742aac82785085de10708a64e9e0370` |
| object-drawing-glyph | 315 | `ea5fad0b65bbb9ab218b2173abc4ae7e368f910ece2adf78dafef4bb3fd0fb62` |
| object-drawing-helpers | 26 | `6d350ea30a670458baf1396910e5d219e4ee057a9364a89ac5592ac89da29c00` |
| object-drawing-shield | 90 | `fbd4a997d5f8f1afd70ba2f4771d496016c0ef2ffe70825e08f91743b7d6c71e` |
| object-drawing-summon | 486 | `28d66af6b2db0802b7e6810d7a5430332d4ac5b885a8d32595ecc7c41752d560` |

The corrected C baseline is fully qualified; no native conversion is applied yet.

## Qualified C checkpoint

Standard:686 selected/completed (685 pass, one optional skip), 370.247s. Server/highres:
53 pass each (160.677s/34.943s). Focused14 pass, 20.682s. Corrected C production build
passes; fresh headless gameplay passes in 38.281s. All13,298 captured results repeat
byte-for-byte. Local c-qualification.json and c-gameplay-qualification.json record
the commands/results. Native drafts remained outside the build for these checks.
