# Wall-edge rasterizer

Status: **qualified Go conversion**, after qualified world-wall conversion
**6b077ec9**. Frozen C baseline: **cbe1f86c**. Removed **284 C lines**; **75,473 C lines /
91 files / zero reference C** remain.

Scope: nox_xxx_edgeDraw_480EF0, **279 C block lines**. Its only production caller
is the Go world-wall renderer and ignores its scratch/pointer-shaped return.
The flags argument is unused. Retire that return and argument in the private Go
API; the fixture continues to check all former flag values have identical effects.
Its light-multiplication and backbuffer-pitch C exports become private Go calls;
the underlying real Go owners remain. No C algorithm will remain solely for tests.

Fixtures reuse world-wall ownership: actual image bag/handles and lazily allocated
pixel data, production pixel-row initialization, pinned Go renderer pixels, actual
clip/configuration words and viewport offset. They cover mixed opaque/transparent
RLE runs, lighting gradients, clipped rows/columns and partial runs, independent
horizontal intervals, height cropping, alternate-row rendering, nil and unsupported
image types, masked type 0x43, 1/255/260-pixel rows and unused flags. Unsupported
malformed RLE buffers are outside the asset contract; the actual image owner rejects
unloadable pixel data before this renderer can receive an empty buffer.

Independent contracts calculate pixels from source colors and widened signed
interpolation, require image/light input immutability, check no-op/excluded pixels,
and compare complete alternate-row output including odd first rows. Distinct
background rows expose transparent-run copying. Preserve the original low-resolution
copy-length behavior: an initial partial transparent run contributes no copied
pixels, while later transparent runs count in full. Do not silently replace this
with an idealized rectangle copy during translation. Review that behavior separately
if a visual correction is wanted.

The shared C pixel-row table, clipping/configuration words and tile callers remain
with their current owners. A broader buffer ownership redesign and nearby mapped
color-codec callbacks are outside this batch. The direct tests cover the same real
owners in default, server and high-resolution configurations. Qualify affected
wall/minimap/object-render/UI/map callers in all three targets; use exact production
fingerprints to reuse 6b077ec9's builds/ABI, full-assets comparison and fresh twelve-
frame gameplay for the unchanged C baseline. The Go conversion received new ones, documented below.

Local evidence: build/port-wall-edge. Applied draft files become stale; use actual
source and PORTING_STATE.md. Original assets/archive remain unchanged.

## Frozen C expectations

Development-a passed three roots / 880 records in 184.551s. Development-b passed
five roots / 1,340 records in 22.536s, adding 432 boundary and 28 alternate-row
contracts. All prior group hashes remained exact. Expectations are now locked;
three-target affected C qualification completed successfully.

| Group | Records | SHA-256 |
| --- | ---: | --- |
| alternate-rows | 28 | `b142b0464cb1c7c6ac78e03df23c86705863f95256708bb863a7d1357a4f4650` |
| boundaries | 432 | `019ff8ac60dcb8d7f55ab32ab32741d047d323013b82ce1311a5a1b24845652a` |
| guards-runs | 140 | `2fa1919602f31736fcddce1e22c52f23a594ae9b22d713e6af3b4810bd60e577` |
| matrix | 720 | `9c5fe9841d7d4c0d6e7df731292b21a4b81c81ee5fb18f5f726fb9a612907e40` |
| pixel-contracts | 20 | `52fbaff56448c3f3f29fe0786ca4458ede98aabf7fb6acca8012492815a5c507` |

Completed C qualification: **251 / 249 / 251 roots**
in **132.964 / 215.611 / 137.833s**. All selected roots started/completed with the expected
optional prerequisite skip. All five hashes match in all targets. Exact production
fingerprints match 6b077ec9, supporting reuse of its three builds/ABI, exact known
full-suite comparison and fresh 12-frame chapter/minimap gameplay. No production
edge algorithm changes were included. All readers joined before baseline commit.

## Go conversion and qualification

The first focused Go run passed all 1,340 frozen records in 185.368s. Review then
moved per-run shading scratch into one reusable array per draw; complete native
qualification covers that final source. No algorithm or expectation correction
was needed.

The Go renderer keeps actual image and pixel-row ownership, shades intersecting
opaque runs with the existing light-multiplication owner, and decodes skipped rows
and runs to preserve stream position. Low-resolution copies retain the original
run-dependent length. The private Go API drops the unused flags argument and
ignored scratch return. The only caller invokes it directly; the edge C interface
and two now-private helper exports are retired. Shared C state remains live for
tile drawing. No C implementation is retained solely for tests.

Affected default/server/highres qualification completed **251 / 249 / 251 roots**
in **227.171 / 213.614 / 139.648s**. Every selected root started/completed, with the expected
optional prerequisite skip. All **1,340 records / five groups** exactly match frozen
C in every target, alongside the earlier wall/minimap/renderer caller contracts.
All three production builds pass ELF32/i386/SSE2/CGO and interface checks: three
retired C interfaces absent, no test helpers. The full-assets failure multiset and
package outcomes remain exactly **1,553 entries; 15 pass /3 fail /32 skip**.
Fresh client-wall-edge-port gameplay matches all **12** chapter/minimap frames
with updates disabled. Source fingerprints stayed unchanged through qualification.
Evidence: build/port-wall-edge/native-qualification.json. Solo gameplay does not
establish remote multiplayer coverage or correctness for malformed image assets.
