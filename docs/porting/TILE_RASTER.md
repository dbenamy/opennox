# Tile texture/fill callbacks and raster helpers

Status: **qualified frozen C baseline**.
Production remains **75,473 C lines / 91 files / zero reference C**, from qualified
wall-edge conversion **51230d6e**.

Five connected C routines / **1,445 function-block lines**: texture callback
nox_xxx_tileDraw_4815E0, fill callback sub_481770, their private unrolled raster
helpers sub_4831C0/sub_484450, and wrap-threshold setup nox_xxx_tile_486060.
The two unrolled helpers account for 1,300 lines. Static store-footprint analysis
confirms a 46-row diamond, 2,116 packed input bytes and the exact fill-word phase.
Runtime original-C tests remain the behavioral baseline.

## Owners and contracts

Fixtures use the real Client tile buffer initializer, stride/ring extent, sight
shape tables and wrap-threshold initializer. They own and restore actual shared
C words, tile definitions, image bag/handles, configuration and callback slots.
The buffer is allocated/freed through its production owner. The callbacks are
selected by the ordinary video-option function and invoked through the actual
three-argument C dispatch slot.

Frozen **797 records / four groups** cover:

- 240 primitive cases: four byte strides, five aligned/unaligned origins, six
  full-width color patterns and texture/fill operations. Independent per-pixel
  expectations check all diamond bytes and untouched surroundings; source bytes
  remain unchanged. The two color halves deliberately differ in several cases.
- 512 callback/ring cases: normal, boundary and wrapped origins; both selected
  callbacks, four colors, four definition flag patterns and upper-word tile-index
  inputs. Expectations compute complete buffer bytes and configuration effects.
- 24 nil-image/scroll-origin cases: signed origins, source/destination offsets,
  actual callbacks, unchanged definitions/words and preservation of prior flags.
- 21 setup cases: stride multiplication overflow, pixel shifts and exact dirty/
  wrap-threshold state. Threshold is 45×stride + (46<<pixelShift), distinct from
  the 46 per-row shape entries; the normal renderer uses pixelShift=1.

The fast fill preserves the original 32-bit word pattern, with a low halfword at
an odd starting column. The wrapped fill restarts its pattern at each split. Do
not silently replace either with a uniform uint16 fill. Actual game image and
buffer ownership keeps source and destination separate; malformed images and
invalid tile indices remain outside the valid asset contract.

## Interfaces and planned conversion

All six remaining C composition calls use the same void, three-argument callback
slot. Their old selected functions have decompiled char/char-pointer return types
and differing argument lists; no production caller consumes those returns. Use
that actual void/three-argument dispatch ABI for both retained Go-backed exports.
Texture ignores the third argument; fill preserves its uint16 tile-index narrowing.
The upper-word fixture cases establish this before normalization.

Retire the two private primitive interfaces, the wrap-setup C interface and the
now-private textured-floors getter export. Keep the getter's actual Go owner and
shared C state. In particular, the flat-floor flag definition is embedded in the
fill function's source block and must remain when that block is removed. No C
algorithm will remain solely for tests.

## Baseline evidence and gameplay

Development-a/b discovery failed in **44.125 / 43.747s** because the fixture did
not match the existing C typedef spellings for two adjacent flags. Each driver
and compiler was joined before correction; production was unchanged. With the
first actual declarations matched, development-c passed all four roots / 797
records in **182.243s**. Expectations are frozen for three-target affected checks.

A proposed fresh-file TexturedFloors=0 gameplay setup was ineffective: E2E's config
reader deliberately uses fixed e2eInputConf and ignores nox.cfg. Both runs matched
all twelve ordinary textured frames. Their flat-mode claims were invalidated in
build/port-tile-raster/ineffective-config-gameplay.json; they are not flat-floor
coverage. Use the actual video-options checkbox (ID2033), capture it and gameplay,
and repeat in an independent process before accepting that additional baseline.
This changes scenario input, not engine code or test expectations.

Qualification covers affected wall/edge/minimap/object-render/UI/map/tile owners
in default, server and high-resolution configurations. Exact production fingerprints
can reuse 51230d6e's production builds, ABI audit, known full-suite comparison and
normal twelve-frame gameplay for the unchanged C baseline. The Go conversion
requires new builds/ABI, the exact known-failure comparison and fresh gameplay in
both normal and GUI-selected flat-floor scenarios.

Evidence and unapplied implementation drafts: build/port-tile-raster. Applied
fixture/freeze drafts are stale; use actual source and PORTING_STATE.md. Original
assets/archive remain unchanged.

## Frozen expectation hashes

| Group | Records | SHA-256 |
| --- | ---: | --- |
| callback-ring | 512 | `8b2d59847b69464a4b12bb8a6966980f920c0bd16a63de3b2a6f7e582df8f857` |
| nil-origins | 24 | `5dbfe27adaddadff79aa3874cc1f2a571bb86dc293be746068ac98d97ba6b943` |
| primitives | 240 | `226ad32567debd8455361abf42a4a0e69810d5e163206618479aa548163992b7` |
| wrap-setup | 21 | `50c2681ec7661f5d919bdfe9a7b2e63a98e3a620ec8d2df633846724cf812129` |

## Completed C qualification

Affected default/server/highres checks completed **259 / 257 / 259 roots**
in **132.071 / 214.990 / 139.289s**. Every selected root started and completed, with the
expected optional prerequisite skip. All 797 records match exactly in every target;
source fingerprints stayed unchanged. Production exactly matches51230d6e, validating
reuse of its three builds/ABI, full-suite failure comparison and normal12-frame
chapter/minimap replay.

The actual options control requires a hover tick before press/release. The first
coordinate trial used the wrong scale, and the first instantaneous checkbox click
did not toggle it; those captures are development evidence, not flat-floor coverage.
The corrected scenario uses the existing slow-interaction action after hovering,
moves away, and checks the option before/after. Exactly101 pixels in checkbox2033
change; all fourteen other checkbox regions remain identical. Visual inspection
confirms Textured Floors is unchecked. No engine input behavior was changed.

Fresh client-tile-flat-gui2-c captures **14 frames**; independent
client-tile-flat-gui2-c-repeat matches all14 with updates disabled. The tracked
[tile-raster-flat.yaml](tile-raster-flat.yaml) and
[pixel manifest](tile-raster-flat-pixels.json) preserve the scenario and decoded
NRGBA hashes. This includes the option transition plus chapter/minimap gameplay.
The engine may temporarily restore textures for flagged tile definitions; that
existing fallback remains live, while direct fixtures independently cover both
callbacks and flag states. E2E ignores nox.cfg, so recovery must use this GUI input
sequence. Production remains original C, with no algorithm edits in this baseline.
