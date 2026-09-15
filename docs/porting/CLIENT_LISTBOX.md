# Listbox widgets

Status: converted and fully qualified against C baseline **43ac4785**.

The connected batch covers eleven callbacks/helpers in GAME3.c from 004A28E0
through before 004A4840, plus the constructor translation unit. Original scope:
1,226 +171 physical C lines /twelve routines. The following class-menu code is
outside this batch. Reuse the real GUI, renderer, fonts, input and image owners
from the entry tests, plus the production palette initializer with complete
save/restore of its seventeen cells and pointer-table entries.

## Confirmed insertion bug

`TestClientListboxMiddleInsertionContract` failed against original C in a guarded
run that selected, started and completed its one test (19.035 seconds). Starting
with `one`, `two`, `three` and inserting `X` at index1 should yield `one`, `X`,
`two`, `three`; instead row2 remained `three`. The source treats a byte offset
as a uint32-pointer displacement during shifting, multiplying it by four.

The original regression allocates sixteen real rows, keeping the mistaken
source/destination positions inside owned storage (rows4/5/8/9). It proves the
order error without relying on an out-of-allocation access. The correction uses
byte addressing. Evidence: `build/port-client-listbox/original-insertion*`.

## Other prerequisite corrections

- Bound row/label terminator indices to255/63 after the existing bounded copy.
  Preserve the original short-string padding and unused final unit. The original
  source already marks these unbounded terminator indices as TODOs.
- Allocate and initialize one extra multi-selection slot for the sentinel.
  Selecting every row and the removal helper both need that slot. Use memmove
  for the overlapping selection shift.
- Reject negative head-removal counts before row addressing.
- Stop row lookup before reading past allocated capacity. Preserve reads of an
  existing spare row when logical count is lower than capacity.
- Stop auto-scroll when advancing no longer changes the scroll position, so one
  wrapped row taller than its viewport does not keep the callback looping.
- Stop single-line clipping when the temporary text becomes empty, including
  viewports narrower than the widget padding. Preserve stored row text.
- Guard the second sentinel write after an empty-area click with every row
  selected; the existing extra slot is sufficient without a write beyond it.
- Represent the selection union word as uint32 in the Go layout. It contains
  either a scalar index/-1 or a C array address; treating all those bits as a Go
  pointer is misleading. Its 386 offset48 and size56 remain unchanged, and no
  Go listbox caller accesses the old pointer field by name.

These are local, reversible corrections before freezing the baseline. Seven
independent contracts now cover insertion, text bounds, full selection/removal,
scroll bounds/termination, narrow drawing, full-selection empty clicks and ABI. The original bounds failures are established
by source/ownership audit; they are not used as invalid-memory golden captures.
Both C and native qualification are recorded below.

Preserve other defined behavior, including the unusual existing selection
adjustment on head removal, until a separate review chooses to change it. Do not
silently replace its comparison rules with conventional list-selection behavior.
Fixtures will keep row count/capacity and memory bounds explicit.

Corrected-baseline production C was89,979 /101 files /zero reference C: ten prerequisite
lines added since the fully qualified entry conversion. Local drafts, scope and
raw evidence are under `build/port-client-listbox`.

The seven contracts passed together in `c-c` (166.844 seconds). Initial capture
comparison identified raw returned row-text addresses in nested render results
and window addresses in double-click notifications. Normalize only these known
owned references. The drawing capture already repeated unchanged. Expanded
constructor/control/resize boundary coverage is included in the frozen baseline.

## Frozen C baseline

Four groups contain **13,708 results**: rows/capacity1,200; input/selection8,640;
drawing2,048; controls/constructor boundaries1,820. Captures include full window,
row and data storage, selection sentinel, callback-time notices, input, named
images and renderer/pixel state. Row/input/drawing groups matched `c-c`→`c-d`;
all four matched the subsequent `c-e` run. Expected counts and SHA-256 digests
are in the committed tests. No generated C algorithm is retained for testing.

Seven independent contracts and four matrices passed together (11 tests,
26.822s). Affected client116 /60.573s, server115 /177.096s and highres116
/76.003s passed, with selected tests started and completed. The client built in
58.928s and passed fresh warrior gameplay in34.900s, null audio and reference
comparison enabled. The full asset suite exactly preserved all1,553 failure
entries and15pass/3fail/32skip package outcomes. Source fingerprints taken before
the full suite remained unchanged. Artifacts: `build/port-client-listbox/c-final-*`,
`c-qualification.json`, `c-full-suite-*`, `c-gameplay-qualification.json`.

C caller audit retains four interfaces: single-selection input, pre-events,
scroll-index lookup and constructor. Eight other routines became private Go
helpers, including the entry composition-listbox initializer. Complete twelve-
routine replacement removed1,407 corrected-baseline C lines. No C algorithms
remain solely for listbox testing; the index probe calls Go directly.


## Native qualification

`client_ui_listbox.go` owns row/selection allocations, events, scroll/slider
updates, constructor and both draw modes. Go entry construction now calls its
initializer directly. Eight private interfaces retired; four C caller interfaces
remain. Short-string padding, spare-row storage, low-byte row heights, signed
scroll calculations and the existing selection quirks remain observable.

The first build found a C export pointer-type mismatch; correcting its signature
preserved the header ABI. Review also corrected font lookup before a constructor
has established its default notification owner. The first executing native run
then matched all **13,708 results /four groups byte for byte**, with no golden
changes or subsequent behavioral corrections. All eight independent contracts
passed, including new copied-options, default-owner and deferred-allocation
cleanup checks. Focused12 tests completed in173.467s. Matrix result counts are
also fixed explicitly, so dropping cases cannot silently reduce coverage.

- Accumulated port tests: 750 /415.437s, all selected tests ran and finished.
- Affected server: 116 /174.847s; highres: 117 /74.566s.
- Three production builds passed; ELF32/i386/SSE2/CGO verified, four required C
  interfaces present, eight retired symbols and test helpers absent.
- Full asset suite preserved the exact1,553 failure entries and15pass/3fail/32skip
  package outcomes. Fresh warrior gameplay passed in37.944s.
- Character-name typing compared the existing blank/edited screenshots in
  11.389s. Both scenarios used null audio and reference override=false.
- 1483 Go/C/header fingerprints remained unchanged throughout qualification.

Production C is now **88,572 lines /100 files /zero test-reference C**. The
conversion removes1,407 lines from the corrected baseline (net1,397 since entry).
Accumulated frozen coverage is **416,826 results /1,113 groups**, plus independent
contracts. See `build/port-client-listbox/qualification.json`, capture comparison,
binary verification, source verification and scenario results for local evidence.
Older captures are preserved in verified gzip archives with restoration metadata.

The captures cover positive row capacity and audited allocated storage; they do
not establish behavior for arbitrary corrupted counts or invalid caller pointers.
The next connected candidate is window geometry, visibility, flags and draw-data
helpers (360 C lines /26 routines), reusing these real GUI owners.
