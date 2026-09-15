# Health, mana, potion, weapon and charge meters

Status: corrected C baseline qualified and ready for conversion. Production is
still C. All3,858 results /11 groups are frozen and match isolated and wider
runs;12 independent contracts also pass.

## C baseline qualification

| Check | Result |
| --- | --- |
| Affected client tests | 160 passed,78.651s |
| Server variant | 159 passed,182.282s |
| Highres variant | 160 passed,84.290s |
| Production client | ELF32/i386, CGO, SSE2;37 meter symbols; no test symbols |
| Full asset suite | Exact1,553-entry known failure multiset;15pass/3fail/32skip packages |
| Fresh warrior gameplay | Passed,36.070s; reference comparison enabled |
| Fresh wizard HUD repeat | Passed,39.371s; inspected health/mana reference |
| Source integrity | 1,513 Go/C/header fingerprints unchanged during qualification |

Captured groups cover scalar/unaligned player stores, labels and visible glyphs,
bar visibility, bubble state/RNG, charge-row deduplication, class-specific window
creation/recreation, heartbeat timing/sound, input, potion priority/counts/use,
and weapon durability/rendering/tooltips. The wizard scenario and reference
hashes are tracked beside this report. Runtime performance was not measured by
these qualification times.

No C algorithms are retained solely for tests. The sound observer records IDs
and volumes, then continues the actual implementation; it is absent from the
production binary. C count is87,888 /98 files /zero reference C. Conversion will
also move the private seven-record array into Go: the caller audit found no
remaining C consumer outside this batch.

Raw evidence is under build/port-client-meters: baseline-hashes.json,
c-repeat-comparison.json, c-qualification.json, c-binary-verification.json and
caller-audit.json. The initial qualification failure and the owner repair are
recorded below; the final baseline is c-o/c-p, not the earlier poisoned-tube
capture with a missing image.

Scope:37 routines /1,069 physical C lines. Whole389-line
client__gui__guimeter.c plus680 GAME2_1.c lines from00470A90 through before004724E0.
These share seven meter records, windows, embedded potion drawables, bubble
arrays, item/tooltips, renderer state, cooldowns and client sound requests.
Keep their initialization, updates, input and drawing in one connected batch.

## Development notes (historical)

### Initial owners and coverage

Reuse the actual window, renderer, font, player and input-cooldown owners. The
existing font-reference helper accepts optional aliases so the constructor's
small font name resolves through the real renderer lookup. A37-operation probe
calls the production C ABI and exposes callback addresses. The meter environment
owns/restores named globals, seven20-byte records, mapped meter/potion state,
charge-raster rows and cursor storage. Embedded potion drawables have distinct
ownership from the ordinary sprite pool.

Initial tests cover scalar updates, cooldown interactions, numeric labels and
health/mana bars; matrices capture state and real renderer pixels. Constructor,
item/inventory, bubble/RNG, callback and sound-request coverage still need to be
completed before freezing. Do not confuse a disabled client-audio backend with
observing requested sound IDs and volumes. The existing server-audio fixture
covers another owner and is not a substitute.

The first fixture discovery failed because Cgo requires the existing spelling of
three unsigned global declarations. The next discovery exposed the existing typedef spelling of the bubble flag;
that is corrected too. At that checkpoint c-c was starting with six initial tests; final qualification
is recorded above. Actual source supersedes historical ignored drafts in build/port-client-meters.

## Prerequisites before freezing

The charge label used wchar2_t WideCharStr[4] for a signed32-bit decimal count.
That fits only three characters plus the terminator. Expand to12 units, enough
for the entire signed32-bit range. The matrix includes999/1000,65535 and signed
limits. Original overflow evidence is the source capacity/format contract;
no out-of-allocation stack write is executed as a golden.

The small health/mana bar divided by maximum unconditionally, while the tube
already handles a zero maximum. Select zero fill height when maximum is zero,
with an explicit empty-meter contract through both draw paths. All ordinary
nonzero-maximum arithmetic remains unchanged. These two corrections each change
one existing line, so C remains87,884 /98 files /zero reference C. They are
confident reversible choices recorded for later review in DECISIONS.md.

Next: finish real-owner coverage, audit/repeat captures, qualify affected targets
and integration, commit/push the C baseline, then translate and qualify the whole
connected batch. No user decision is pending.


A narrow porttest-only observer now records client sound ID/volume requests at
entry to the existing sound function, then lets its actual lookup/allocation/
playback path continue. It does not replace the audio implementation. The hook is
compiled only by the tagged fixture's CFLAGS; production builds must have no
observer symbol. Four conditional probe lines in GAME2.c raise the physical C
count to87,888 /98 files /zero reference C. They contain no retained test-only
algorithm. An independent heartbeat/toggle contract checks IDs896/901, volume,
and repeat suppression through the actual input-cooldown owner.


The first executing pass c-c completed all six initial tests. Its three independent
contracts passed; the scalar/label/bar matrices produced1,081 results and failed
only their deliberately unset expectations. Those captures are still unfrozen.
A renderer wrapper now observes numeric-label text while forwarding the actual
DrawString call, allowing an independent full signed-count text contract. Run c-d
includes that observation; its capture shape intentionally changes before freeze.

An initial production client build is available for a fresh wizard scenario. The
first build command used the repository root instead of the src module directory;
rerunning from src succeeded without source changes. The wizard trial selects
class ID602 (right-hand class choice) using the supplied SelClass.wnd geometry
and existing scenario coordinates. Its first screenshot is a candidate C reference,
not accepted until inspected and repeated with reference override=false.


The wizard dialogue-only trial was insufficient for HUD coverage. Extend through
dialogue dismissal and a short walk; visually inspect the resulting frame with
both health and mana tubes visible. The accepted C reference contains dialogue
and unobscured-HUD frames. Fresh capture passed in39.918s, and a fresh repeat passed
in39.371s with reference override=false. See meter-wizard-smoke.yaml and
meter-wizard-reference.json for recovery; the PNGs remain ignored local artifacts.
The first renderer-observation run c-d passed all four contracts and produced
updated pre-freeze captures. Remaining matrix failures are unset expectations.


Constructor pass c-e completed all nine selected tests in171.495s. All five
independent contracts passed, including class-specific window creation and the
missing-player case. The four unfrozen matrices produced1,115 results; comparison
failures are expected until coverage and independent-process repeats are complete.
The constructor matrix exercises real potion definitions, localized strings,
class-dependent window trees, drawing, recreation and separate weapon/charge
window constructors. Both C and Go screen dimensions are owned at256×256.

Inventory fixture audit: the actual client cell is148 bytes and the search uses
four columns with21 rows per column (84 cells), despite the UI exposing20 rows.
Use the declared C type and assert its size; older scratch notes describing144
bytes or an80-cell allocation were wrong. The fixture owns/restores the actual
inventory and equipped-item arrays, and observes the real queued use messages.

The fresh wizard HUD reference was visually inspected and then reproduced with
reference comparison enabled. See meter-wizard-reference.json for hashes and
recovery details. Dialogue and gameplay HUD frames both matched in the fresh run.


Run c-f failed discovery on a fixture drawable field name and an unused import;
corrected without production changes. Run c-g completed18 selected tests in
93.693s. Seven independent contracts passed, including bubble lifetime/RNG,
constructor ownership, potion priority, scalar, sound and text contracts. Two
contracts failed usefully: charge drawing had an all-zero geometry table because
headless owners skip full blob initialization, and quick-potion tests set only
the Go cursor while the C gate reads the mapped legacy cursor cell. The fixture
now copies the exact bounded production charge table and durability threshold,
restores them afterward, and sets the actual legacy cursor input. Keep these
assertions; do not freeze the earlier incomplete captures. Run c-h adds actual
weapon definitions/equipment/inventory lookup, durability and tooltip coverage.


Run c-k passed12 independent contracts and produced3,858 results across11 groups,
covering all37 probe operations. Run c-l repeated every capture byte-for-byte.
The initial wider160-test qualification found one bubble-group difference;
all other tests passed. A dedicated meter capture prefix avoids dumping unrelated
large corpora during diagnosis. c-m was cancelled after a fixture rewrite syntax
error prevented installing that prefix; c-n reproduced the wider mismatch.

The first differing fields were solely image-observer output positions after a
nil poison-image draw (0,0 versus31,28). Pixels, bubble state and RNG were identical.
The plain poisoned-tube case had omitted the actual image owner. It now supplies
the owned poison overlay, so this case exercises a real image draw. c-o changes
only that group; its new input/capture is being checked in wider run c-p. The
original k/l captures remain available; this is an explicit fixture-input repair,
not a change to C or a golden update hiding a port discrepancy.
