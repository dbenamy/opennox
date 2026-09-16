# Wall-edge rasterizer

Status: **qualified frozen C baseline**, after qualified world-wall conversion
**6b077ec9**. Production remains **75,757 C lines / 91 files / zero reference C**.

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
frame gameplay for the unchanged C baseline. The Go conversion requires new ones.

Local evidence: build/port-wall-edge. Applied draft files become stale; use actual
source and PORTING_STATE.md. Original assets/archive remain unchanged.

## Frozen C expectations

Development-a passed three roots /880 records in184.551s. Development-b passed
five roots /1,340 records in22.536s, adding432 boundary and28 alternate-row
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
fingerprints match6b077ec9, supporting reuse of its three builds/ABI, exact known
full-suite comparison and fresh12-frame chapter/minimap gameplay. No production
edge algorithm changes were included. All readers joined before baseline commit.
