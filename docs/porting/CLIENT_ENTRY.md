# Text-entry widgets

Status: converted and fully qualified against C baseline **4087008d**.

## Scope and owners

Connected entry-field constructor, keyboard and mouse events, text/composition,
focus, colored/image drawing, and editing-context lifecycle. Initial GAME2_2.c
scope is ten routines /601 physical C lines; two thin GAME1_2.c context callers
and the two private context globals complete the scope. Only the input callback
and colored draw callback have remaining C callers; the other ten function
interfaces can retire. The corrected entry section is 612 lines, the two thin
wrappers nine lines, plus four C declaration/definition lines for context globals.

Fixtures use the real GUI, window pool, renderer, fonts and RLE image owners,
real language-bearing StringManager and real input.Handler. A small seat.Input
adapter delivers keyboard/composition/text events and observes text-input mode.
Composition languages create their actual listbox, buttons and slider. Captures
include window/data words, listbox item bytes, input state, notifications, focus,
image lookups and renderer/pixels where relevant. Pointer normalization uses
identified owned objects only. Raw UTF-16 units, including unpaired surrogates,
are preserved in text-state tests.

## Existing behavior corrected before freezing

The original C implementation failed `TestClientEntryFocusedDestroyContract`:
destroying a focused entry left text input enabled and the active-entry global
set after its data was freed. Original keyboard/composition-owner and constructor
contracts passed. The guarded run selected, started and completed all three
contracts in 169.786 seconds; the ownership regression was its only failed test.
Evidence: `build/port-client-entry/original-contracts*`.

The entry destructor now clears active input ownership and disables text input
when this entry is the active owner. It does not change global GUI focus ordering
or disturb a different active entry. This is a local, reversible correction for
later review, authorized by the standing workflow.

Two bounds corrections follow direct source/owner audit. InputHandler stores
composition strings without a length limit, while the entry has a fixed 256-unit
composition buffer. Both C copy paths now copy at most 255 UTF-16 units and add a
terminator, preserving the original untouched bytes after short strings. Review
caught and removed unnecessary zero-padding in the first bounded-copy draft;
a dedicated contract checks both paths with nonzero trailing bytes. The long-composition contract exercises both paths and verifies the
prefix, count, terminator and adjacent entry fields. We did not execute the
original unbounded copy with oversized input.

Both drawing routines previously advanced through the main-text buffer until
main plus composition fit. If composition alone exceeded the available width,
that condition could never become true. Scrolling now also stops at the main
text terminator. Narrow-width and composition cases exercise the corrected
behavior. The original out-of-bounds loop is not used as a rendering oracle.
Colored versus image width/clipping differences otherwise remain intentional.

## C baseline qualification

Four groups match repeated C captures in separate processes: keyboard 6,144;
limits/composition 1,314; drawing 3,200; events/context 768. Total **11,426 results /
four groups**, plus six independent contracts. The initial bounded-copy revision
passed 736 accumulated tests (one optional skip) in 404.499 seconds and 102 server
tests in 171.740 seconds. Review then refined the copy to preserve unused bytes;
all four hashes remained unchanged and the new trailing-byte contract passed.

Final source qualification: ten focused tests (13.125 seconds), affected client
104 (54.050), server103 (168.356), highres104 (68.180). All selected tests started
and finished. The C client built in 56.208 seconds; fresh-asset warrior gameplay
passed in 34.592 seconds with existing pixel references and override disabled.
The full asset suite in 48.682 seconds exactly matched the known 1,553 failure
entries and 15 pass /3 fail /32 skip package outcomes. All 1,473 Go/C/header source
fingerprints remained unchanged during these final checks.

[The typing scenario](entry-input-smoke.yaml) also exercises the real character
creation name entry: clear the default name, type `abc`, delete one unit, append
`d`. The C capture visibly showed the blank field and then `abd`; a second C
process passed pixel comparison with override disabled (11.187 /10.988 seconds).
It uses separate new reference images, never overwrites the warrior references,
and runs under Xvfb with null audio. [Image hashes](entry-input-reference.json)
record the oracle; PNGs and raw logs remain ignored local artifacts.

To recover the typing reference, build the C baseline commit, use fresh assets
and an empty save directory, copy the scenario into a new run directory with an
empty `testdata` child, and run with `NOX_E2E` pointing to that scenario. Use
`NOX_E2E_OVERRIDE=true` only for this new C reference capture, then repeat with
`false` and verify the image hashes. The Go conversion must use `false`.

C `calloc` allocations are not registered by the Go allocation observer. The
baseline makes no allocation-liveness claim based on that observer. A separate
native ownership contract checks the new Go-owned entry allocation through
deferred GUI cleanup, without reading freed memory; it passed on the first native run.

Physical production C at this corrected baseline: **90,594 /101 files /zero
reference C**, an 11-line increase for the prerequisite corrections. Local
artifacts and drafts: `build/port-client-entry`. The C baseline is committed and pushed.

## Go conversion

Twelve routines moved together; ten private function interfaces and two C globals
are retired. Only the entry input callback and colored draw callback retain C
bridges. Go callers invoke Go directly. Editing context is now a boolean and a
typed window pointer: the old four-byte allocation was only a presence token.
Context teardown still leaves focus handling to window lifecycle. The remaining
listbox constructor/init dependency is C and owns its own allocations.

Preserve signed16 maximum length, exact raw key-state comparisons, full-width key
switching before scan-code truncation, digit-before-alphanumeric filter priority,
locale-dependent CRT character classification, UTF-16 code-unit deletion and
password masking, constructor half-buffer initialization and notification timing.
Colored and image rendering retain their different width/clipping behavior.

The first native guarded run passed all eleven focused tests in 168.124 seconds.
All 11,426 results /four capture groups exactly matched C, without native fixes or
changed goldens. The new allocation-liveness contract also passed. Final qualification also passed:

- Accumulated default: 738 (403.240 seconds), including the existing optional skip.
- Affected server: 104 (172.067 seconds); highres: 105 (68.098 seconds).
- Three production builds: client 57.188, highres 8.616,
  server 53.998 seconds. ELF32/i386, SSE2/CGO, both required C bridges,
  twelve retired function/global symbols and absent test helpers verified.
- Full asset suite: exact known 1,553 failure entries and package outcomes.
- Fresh warrior gameplay: 37.030 seconds; character-name typing:
  11.107 seconds. Both compare existing C pixels with override disabled.
- All 1,476 Go/C/header source fingerprints remained unchanged during qualification.

Accumulated frozen coverage is now **403,118 results /1,109 groups**, plus
independent contracts. No native corrective patch or golden update was needed.
Physical production C is now **89,969 /101 files /zero reference C**, down 625
lines from the corrected baseline (614 net from before the prerequisites).
