# Slider widgets

Status: conversion fully qualified against C baseline **7df42a4b**.
All 26,692 original-C results remain exact.

## Scope and ownership

Twelve routines in GAME3.c, from 004B4860 through 004B5640, cover horizontal and
vertical sliders: construction, child buttons, callback setup, input, value/range
updates, focus, destruction and drawing. The selected block is 599 physical C
lines. Before conversion, production C remains 91,478 lines / 101 files; there is
no test-reference C.

The caller audit retains two C interfaces: the constructor, used by list boxes
and GAME2_1.c, and the vertical input callback used by GAME2_1.c. Ten private
interfaces can retire. Existing Go constructor callers will call Go directly.
The actual GUI owns the windows and thumb buttons; each slider owns its copied
16-byte SliderData. Tests exercise real allocation, callbacks, focus and renderer
owners, including destruction between cases and final GUI cleanup.

## Baseline and contracts

Four original-C capture groups repeat exactly in separate processes:

| Group | Results |
| --- | ---: |
| Construction and programmatic values | 11,088 |
| Input sequences and notifications | 8,320 |
| Numeric state and range changes | 4,212 |
| Drawing states | 3,072 |
| Total | 26,692 |

Tracked expectations are SHA-256 hashes of structured results. Captures include
all 101 words of each window and thumb, input and owned SliderData, notification
order and raw payloads, focus, pixels and renderer state. Only identified pointer
fields are normalized; scalar notification values retain their exact bits.
Full captures and qualification logs live in `build/port-client-sliders`.

Cases cover orientations, image/color modes, root/child windows, explicit/default
notification owners, enabled and no-focus flags, ordinary/equal/reversed/extreme
ranges, dimensions, accepted/rejected values, mouse edges and outside positions,
recognized/unrecognized keys, exact key states 0–3, focus transitions and repeated
range updates. Drawing checks cover transparent color sentinels, missing images,
highlight and thumb states. Extreme/wrapped coordinates use state-only captures
to avoid unrelated raster preconditions.

Independent contracts verify actual child ownership, copied data, the required
constructor flag, scale and thumb position, visible pixels, rejected values and
focus notifications. Invalid styles must leave constructor inputs and parent
children unchanged. Both orientation bits select the horizontal implementation.

The first fixture compile rejected unsigned 0xffffffff as a 386 int event code;
using -1 preserves the intended raw event bits. No production behavior changed.
The initial capture failures were unfrozen hash placeholders, not engine failures.

## Compatibility decisions to review later

Preserve the original nonzero-minimum quirks: horizontal programmatic positioning
uses value × scale, keyboard positioning uses (value − minimum) × scale, and drag
conversion does not add the minimum. The vertical image draw callback is a no-op;
its thumb still draws. Tab-navigation helpers called by the input handlers are
also existing no-ops. These are compatibility choices, not new UI fixes.

Keep arithmetic at the hosted 386 precision: stored scales are float32, operations
use 53-bit intermediates, and float-to-position conversion passes through int64
before narrowing to 32 bits. Equal ranges and overflow cases retain the baseline
state. The remaining C input bridge must preserve raw key-state arguments rather
than normalize all states at least two to a pressed boolean.

Headless warrior gameplay checks integration, not every slider interaction. The
actual-widget fixtures supply slider branch coverage; physical mouse/display
feel remains a manual release check.

## Native implementation review

The first native run caught a missing explicit copy into owned SliderData:
`alloc.New` allocates zeroed storage and uses its argument only for the type/size.
The independent ownership contract and all four capture groups detected this.
The implementation now assigns the input value after allocating. Review also
keeps focus dispatch on the client GUI, matching the original C bridge, rather
than choosing the window's GUI. Frozen expectations are unchanged.

The second native run matched construction, drawing and numeric groups. Input
captures isolated a color-field mistake: `nox_xxx_wndSetRectColor2MB_46AFE0`
changes the background, not the enabled color. Full window words and pixels
identified this without changing notification/value expectations. The native
input handler uses the actual background setter for both hover transitions.

## Native qualification

All twelve routines now run in Go. Two C exports remain for live callers; ten
private interfaces and their prototypes are removed. Go constructor callers and
window callbacks dispatch directly to Go. The original 599-line C block is gone;
production C is **90,879 lines / 101 files**, with zero test-reference C.

Six focused tests pass against unchanged expectations and contracts. The complete
accumulated port corpus and affected server/highres variants pass, with selected
tests verified to start and finish. All three production builds pass ELF32/i386,
SSE2/CGO and symbol audits. The full asset suite matches the exact known 1,553
failure entries and package outcomes (15 pass / 3 fail / 32 skip). Fresh headless
warrior gameplay passes with null audio and reference comparisons enabled.

Accumulated frozen coverage is **375,818 results / 1,103 groups**, plus independent
contracts. Local qualification.json records tests, binaries, captures and gameplay.

Qualification counts: accumulated 720, affected server 86 / highres 87. Fresh gameplay completed in 35.874s. All 1,460 source fingerprints remained unchanged during qualification.
