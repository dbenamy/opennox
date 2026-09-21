# Client resources and lifecycle

## Scope

Port the remaining client resource teardown, frame timing, player colors, image
caches, modal window and menu drawing helpers. The selected files contain 851 C
lines before preserving their shared variable definitions: `GAME1_3.c`,
`GAME2_1.c`, `GAME2_2.c`, `client__drawable__drawdb.c` and
`client__shell__mainmenu.c`. Production remains unchanged while qualifying the C
baseline. Current total: **3,543 physical C lines in 28 files**, zero reference C.

## Contracts

Nine test roots produce eight captures /6,415 cases:

| Behavior | Captured cases | Coverage |
|---|---:|---|
| Sprite teardown | 66 | Every kind, sparse/full graphs, all 55 states, 54 slots and eight directions; each allocation freed once and root last |
| Frame samples | 1,728 | Clock narrowing, wrap/subtraction, low-word sample index, full-width cursor update, count wrap |
| Frame average | 60 | Signed count boundaries, clamp60/default33, wrapped 64-bit sums |
| Player colors | 1,024 | Every channel value, channel order, packed halves, mutable white word, unchanged source colors |
| Active-player colors | 5 | Empty/full/sparse/last-slot traversal; inactive records unchanged |
| Animation caches | 16 | Both lookup failures, repeated loading and clearing, exact lookup names |
| Modal window | 60 | Sizes/colors, repeated create/destroy, real GUI ownership and destruction flags |
| Main menu draw | 3,456 | Background styles, real images/pixels, record sentinels, timers/underflow, seeds, RNG consumption |

The ninth root checks scalar getters/setters and all 256 binding-count values,
including unchanged neighboring bytes. Independent assertions complement frozen
C output hashes. No original algorithm is copied into test code.

## Compatibility and ownership review

- Sprite data originates from raw libc `calloc` in the Go parser. Its teardown
  must use matching libc `free`, rather than the tracked Go allocator. Vector
  slot4 is metadata; kind6 traverses 55 states with 264-byte stride and 54 owned
  groups per state. Children precede their containing group/root in release order.
- The clock callback is Go uint64 but its existing C interface returns unsigned
  32-bit ticks. Narrow before widening for 64-bit history subtraction. Cursor
  indexing uses its low word; cursor advancement uses both words modulo60.
- Player RGB5551 output repeats both 16-bit halves. The final color copies the
  mutable C white word exactly, even when it differs from canonical white.
- Menu image/name sentinel words gate animation and drawing separately. Preserve
  timer post-decrement (including zero underflow), transition order and RNG calls.
  The rectangle uses the low16 color bits; transparency tests the full sentinel
  0x80000000. Highlight selection reads draw.Field0 bit2.
- Repeated modal creation overwrites the current slot while the previous window
  remains GUI-owned. Preserve existing behavior; owner cleanup releases leftovers.

## Reachability audit

The two packed-color conversion function addresses are stored only at
0x973F18+5248 and +7716, with no readers. The video/material callback is registered
only into 0x5D4594+1193504, also without readers. Its private table/count have no
consumers. Remove those unused registrations and four private helper bodies during
conversion. Keep live binding-count and shell-state getters, including their Go
callers in the binding editor and server browser. Move Go player-file callers too.
Shared globals in deleted translation units retain their definitions and values.

## Baseline preparation findings

The initial fixture compile exposed a 386 integer constant overflow; use explicit
uint32 subtraction. The player traversal fixture also needed actual slot indices;
without them the real iterator repeats an entry. Timing contracts initially missed
the C clock narrowing. These fixture corrections preceded capture freezing.

A separate free wrapper conflicted with existing test instrumentation. Reusing the
map-theme Go callbacks then passed probes but stalled intermittently at a GC stop
after the teardown root in an independent run. The corrected instrumentation adds
a resource-only thread-local event buffer to the existing free wrapper, records
addresses without callbacks into Go, and normalizes after observation stops.
It restores thread state and does not change existing map-theme/grid observers.
The failed-run stack and intermediate probes remain in the ignored build folder.

## Qualified original-C baseline

Default repeat, server and highres each pass all nine roots, no skips, eight
matching captures and static checks. The observer regression runs 41
resource/map-theme/world-grid roots three times with GOGC=10; all pass, no skips.
Production source is byte-identical to qualified a8fbf8a8: only seven new tagged
test files and the existing tagged allocator wrapper differ. Reuse that revision's
production builds/ABI, exact known-suite comparison, gameplay and save/load evidence.
See [the baseline manifest](client-resources-c-batch.json) and
[qualification report](client-resources-c-qualification.json).

Native conversion is next. No original C body has been removed by this baseline.
