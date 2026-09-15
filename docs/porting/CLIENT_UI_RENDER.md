# Shared UI rendering and progress-bar port

Scope: thirteen routines for clipping/intersection/flag state, pixel fill and
particle borders, renderer bounds, and progress-bar construction/events/drawing.
The selected address blocks occupy 261 C lines; removing their attached border
comment removes 262 lines in total. This connected batch reuses the qualified
renderer, image, font, GUI and clipping owners. Production C starts at 91,740 lines
in 101 files, with zero test-reference C.

## Owned tests and baseline findings

Fixtures call the actual C routines and real window event/draw dispatch, using
owned RenderData, its C current-state pointer, screen bounds, GUI allocation and
parent links, basic font, material-indexed RLE images and software framebuffer.
They capture renderer/cache state, pixels, draw calls, exact fill bytes, rectangle
mutations/return bits, window geometry/flags/owner and progress values. No GUI,
clipping, text or raster algorithms are stubbed.

Intersection cases cover aliasing the destination with either input, disjoint,
edge-touching, empty/inverted and extreme coordinates, with an independent
no-write-on-rejection contract. Fill cases cover signed nonpositive lengths,
unaligned destinations, distinct word bytes, every remainder and surrounding
sentinels. Preserve the original whole-word plus optional halfword behavior:
odd trailing bytes are deliberately untouched. A general memset would differ.

Clipping cases cover raw flags, bounded/empty/overflowed rectangles, horizontal
narrowing, save/restore and bounds return/state. Repeat testing exposed a fixture
mistake: the C rectangle-copy function is declared pointer-returning, but the
underlying noxCopyRect returns scalar true (1). Capturing an address-relative
offset was wrong. The corrected baseline captures exact scalar bits (0/1), and
native code can correct the declaration to int. All remaining C callers ignore
the result; the Go narrowing caller also ignores it. This is a type correction
with unchanged 386 return bits, not a new allocation or changed clipping rule.

Particle-border tests cover clipping, dimensions, colors and positive radii.
Two initial fixture failures exposed the existing rasterizer's preconditions:
radius zero divides by zero during particle creation; unclipped particle images
must remain inside the framebuffer. The production border caller uses radius4.
The final matrix uses positive radii and contained unclipped draws; out-of-buffer
positions still run with clipping enabled. No production fix is included for the
zero-radius helper. Record it for later review separately from this translation.

Progress cases cover missing style, parent offsets, the constructor's cleared flag,
caller-supplied/default owner, image versus colored drawing, optional background,
fill/border/text, smoothing, geometry and valid/invalid percentage updates. Actual
registered callback dispatch and rendered pixels are checked; an independent
contract verifies that changing zero to full progress changes the framebuffer.

## ABI and ownership

Retain four C bridges used by remaining callers: renderer bounds, particle border,
pixel-fill wrapper and rectangle copy. Retire nine private or Go-only entries,
including the progress-bar C callbacks; registered Go GUI callbacks need no C
trampolines. Internal intersection/fill helpers remain production dependencies,
not algorithms kept solely as test references. Clipping uses the actual current
C RenderData owner so existing C UI callers keep their state semantics.

## Frozen C baseline

Six root tests contain **9,880 results in five groups**, plus independent
contracts. Final C captures repeat byte-for-byte before expectations are frozen.
Progress snapshots include the actual text-smoothing flag, with prior false/true
states and text-enabled/disabled paths.

| Capture group | Results | SHA-256 |
| --- | ---: | --- |
| ui-render-borders | 360 | `53a826566438b955d4cd3a4e3a0e80b2203c866501d82c440fd33eb5255f00ed` |
| ui-render-clip-bounds | 896 | `164c78fa31852e3e5ff07cdbc0b707e12413676d80afbf066fbe85d2e0567bfe` |
| ui-render-fill | 1520 | `44a77c3ae962c59d080bcdc2f6b03604fe9a1e2035448a05eb26ccd7767394e2` |
| ui-render-intersection | 192 | `6bda793def3833cc944df4ba986416119a3e7df94106c4076197f7f7870a4934` |
| ui-render-progress | 6912 | `d503d5b21d7648b83286fd91b37e794fd5f2c941481dd19b2c0d6a14fa537700` |

C qualification passed: focus6, affected client/server/highres
80/79/80; every selected test started and completed.
Production C build passed in 4.251s; fresh headless gameplay passed
in 37.552s with null audio and reference comparison enabled.
No production C change was needed to establish this baseline.

Native conversion has not been applied. C remains 91,740 lines /101 files /zero
reference C. Native qualification will include the accumulated port tests.
