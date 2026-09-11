# Floor-rendering eligibility — 2026-09-11

Scope: 475810_draw_B, called only by the Go client draw loop through a legacy
wrapper. It checks EngineNoFloorRendering first, then queries the tile at the
viewport's World.Max coordinate. Exactly tile values -1 and 255 suppress the
floor; all other signed return values allow it. Disabling floor rendering skips
both viewport and grid access. The helper does not mutate state.

The original-C baseline covers 22,083 cases: 21² coordinate pairs × 10 signed
and boundary tile values × 5 flag combinations, all 32 flags individually, and
a disabled-rendering call with null viewport/grid. Coordinates include map
boundaries, negative values, float32 precision transitions and signed extremes.
Other viewport position fields deliberately differ. Tests check the whole
viewport, flags, guarded grid rows and cells, and restored global state.
Grid geometry and list overrides remain covered by the preceding grid corpus;
this fixture focuses on the owner's integer coordinates, flags and result filter.

The original-C baseline c3c073b5 passes all 22,083 cases on 386. The native
implementation now calls tileAtPoint and preserves the flag short circuit and
World.Max coordinates. The C definition and declaration are removed, with no
replacement ABI. Native qualification is complete. Artifacts:
build/port-floor-eligibility. Pre-port C count: 140,455 lines, 153 files, zero
reference C. No C caller needs a replacement export.

The native change removes 13 physical C lines (the function and trailing blank
line): **140,442 lines**, 153 production C files, zero reference C. Default, server and high-resolution
accumulated checks and all three ELF32/80386 production builds pass. Fresh
floor-eligibility-port gameplay passed both preserved screenshot checks with
overrides disabled. The immediately preceding grid
chunk matched the exact full-suite baseline; this small owner change used
accumulated targeted checks and gameplay without repeating that full suite.
