# Client inventory display and feedback

Status: all sixteen routines are native and fully qualified against pushed C
baseline `18f9cfb4`. Production C is 84,237 / 97 files / zero reference C. The
preceding transaction conversion is `2f34230e` (pushed).

## Scope

Sixteen routines / 1,162 C lines cover item identification, player stats, tray and
weapon drawing, modifier/durability values, tooltips, names and secondary feedback.
Main inventory construction/event dispatch, paper-doll drawing and journal
composition remain for another batch. Three modifier/durability helpers have no
C callers outside this batch; thirteen interfaces still have real C consumers.

| Routine | C lines |
| --- | ---: |
| `sub_4625D0` | 29 |
| `sub_4626C0` | 24 |
| `sub_462700` | 30 |
| `sub_463370` | 14 |
| `sub_4633B0` | 25 |
| `sub_463420` | 9 |
| `nox_xxx_guiDrawInventoryTray_4643B0` | 171 |
| `sub_465D50_draw` | 22 |
| `nox_xxx_inventoryDrawProc_466580` | 26 |
| `sub_466F50` | 61 |
| `sub_4627F0` | 384 |
| `nox_client_makePlayerStatsDlg_463880` | 207 |
| `sub_466660` | 59 |
| `sub_466E20` | 30 |
| `nox_xxx_inventoryNameSignInit_4671E0` | 33 |
| `sub_467750` | 38 |

## Baseline and contracts

Reuse actual drawable allocation/deletion, inventory/equipment, modifier, GUI,
string-manager, net-list, player, renderer and font owners. Identify panels use
real static-text and listbox widgets. Capture actual pixels, text/draw metadata,
renderer state, known pointer identities, cursor, widgets, inventory state and
game requests. Existing meter/inventory fixture constructors accept optional
thing names; previous callers pass none and retain their frozen expectations.

Sixteen tests produce 2,478 frozen results / twelve groups, plus 735
independent coordinate, scalar and raw/scaled-durability cases. Captured matrices
also enforce independent behavioral contracts:

- Element callback identity, ignored classes, two-slot precedence, missing values.
- Durability truncation with fractional, zero and negative multipliers, 16-bit
  health limits and callback/class gates.
- Tray scroll boundaries, visible/hidden rows, stack quantities, charge signedness,
  equipped/alternate decoration, repair overlays and durability thresholds.
- Actual identify rows, cached redraw, selection clearing, armor/damage modifiers,
  combined attributes, two explicit flavor descriptions and icon material colors.
- Mouse regions/cursor modes, tooltip boundaries and mutation of selected item codes.
- Three player classes, XP, armor, carried weight including hidden-row items,
  speed effects and three language-layout inputs.
- Signed player level versus clamped special-mode rank, Unicode names, null player.
- Secondary feedback across visible/hidden cells, old selection, pending codes,
  success, rejection, restoration requests and other status values.

The asset audit confirms two decimal places in damage labels. Fixtures use that
precision. Known modifier pointers in static-image item snapshots are normalized
only at their four declared slots, with ownership checks. C-c/C-d matched all
2,478 records. C-e/C-f also matched after the final precision adjustment; their
twelve hashes are now frozen in the tagged tests.

## Headless inventory scenario

[The tracked scenario](inventory-display.yaml) checks gameplay, open inventory,
player stats and shirt identification. Input accounts for 1024×768 rendering
inside Xvfb's 1280×960 window. Use a fresh asset/save copy, null audio and
GODEBUG=randautoseed=0. Capture references with the committed C baseline; actual
comparison must use NOX_E2E_OVERRIDE=false.

The qualified transaction binary still implements display in C. Fresh runs
client-inventory-display-seeded-{capture,repeat} pass all four complete screens;
the capture was visually inspected. Their result manifests record the binary hash,
seed setting and command. Production source is unchanged, allowing reuse of the
transaction's three qualified binaries and exact known full-suite result for the
C checkpoint; fresh affected tests on all three targets remain required. Native
qualification will rebuild all three targets and rerun the full asset comparison.

## Development findings and retained evidence

The initial panel repeat differed only in a randomly selected localized shirt
description. String-manager Entry.Value uses process-global math/rand. Fixing the
test process's seed made the full comparison repeatable; production localization,
text and compared pixel area are unchanged. Preserve panels-repeat2's failure and
got/diff images. Earlier attempts that missed buttons are development evidence.

Initial fixture failures identified: a uint32/uint SetID mismatch; missing System
message format, item/class names and qualified localization IDs; a raw label that
needed the static-text event handler; and renderer cache reset between case
objects. These were fixture corrections, with production untouched. C-a/C-b
matched ten groups; the other two had 1,280 differences exclusively in declared
modifier pointer addresses. Preserve raw captures and the pointer diff. The
normalization correction yielded the matching C-c/C-d comparison.

Ignored evidence lives in build/port-client-inventory-render: candidate scope,
production callers, state/callee audit, asset format audit, development logs,
captures/comparisons, source fingerprints and qualification scripts. Full logs
stay local. No test-only C algorithm is retained.

## C qualification

Client/default: 204 tests passed in 90.908s. Server: 203 in 190.776s. Highres:
204 in 95.930s. Every selected test started and finished; all 1,552 source
fingerprints stayed unchanged. Production source matches `2f34230e`; each reused
binary matches its qualified SHA256 and was independently rechecked as
ELF32/i386/SSE2/CGO with all sixteen display interfaces and no test helpers.
The previously qualified exact full-suite failure set is reused because only
porttest files changed. The new seeded headless display comparison passed.
Evidence: c-qualification.json, c-binary-verification.json,
c-production-source-reuse.json, c-source-fingerprints.json, baseline-hashes.json
and c-e-c-f-comparison.json. C remains 85,399 / 97 files / zero reference C.

## Native conversion and development

All sixteen display routines are translated in four Go files; 411 lines leave
GAME2_1.c and 751 leave client__gui__guiinv.c. Working C is 84,237 / 97 files /
zero reference C. Three modifier/durability helpers become private Go functions;
thirteen real C interfaces remain. Fixed-signature adapters invoke the existing
production variadic formatter so localized numeric formatting remains unchanged.
Focused comparisons pass; broader qualification is not yet complete.

Native-a passed fifteen of sixteen focused tests; the button capture identified
thirty trace differences from skipping empty key-label draw calls. The original
key-name wrapper returns a non-null empty string, so preserve the renderer call.
Native-b then passed all sixteen tests and matched all 2,478 frozen records.

Review also found that identification looked up the item name before its label.
The original C does the reverse, which matters for random localization variants.
An isolated child-process contract with a real multi-variant string manager failed
before the correction (seed1 chose Unknown-C rather than Unknown-A) and passes
afterwards across 32 seeds, including the next random value. It does not mutate
the parent suite's global random state. Native-c passes all seventeen tests and
still matches every frozen record; hashes are unchanged.

Native header assembly uses the existing bounded UTF16 copy into its 256-unit
mapped buffer. An additional contract checks the 255-unit payload boundary,
long ASCII/Greek/emoji names, termination and unchanged adjacent storage. The old
C's unterminated temporary/unbounded concatenation on oversized inputs is not
executed as an oracle. Native-d passes all eighteen tests in 27.070s, including this addition.

## Qualification interruption: spell-force callback

The first accumulated run selected 839 tests, started 757 and completed 756 before
runtime GC terminated TestSpellEffectsForce. The invalid pointer was 0x3f8ccccd
(the raw float word for distance 1.1) in CallVoidPtr3's argument frame. This is an
incorrect numeric-to-pointer conversion in the earlier spell-effect port, not an
oracle mismatch. The complete log/result and original source fingerprints remain
in native-qualification-first-failure. All three concurrent production builds
finished successfully and were joined; no source edits overlapped them.

The correction keeps the distance word and opaque callback argument as uintptr
through CallVoidUPtr3, including record loading and direct Go callers. A new
forced-GC variant runs the existing force matrix with collection inside the C
callback; its expected captures remain the original C hashes. Full qualification
and three builds will rerun on the corrected source. See SPELL_EFFECTS.md.

The corrected accumulated run passed all 839 tests in 430.745s. Its explicit
selection predated the new ForceGC test; the focused 18-test spell-effects run
already passed that regression on the same source in 175.179s (including package
compilation). Add ForceGC to the tracked accumulated pattern for future runs;
server/highres selection includes it now. Report these runs separately rather
than claiming one 840-test accumulated invocation. Qualification is continuing.

Highres passed 279 tests in 125.110s, and all three production interfaces verified.
The full asset suite then found TestCodeStatic in common/memmap/nox: the new
bounds fixture used PtrUint32 at the font global's former mapped address. This
was intended as adjacent-buffer guard storage, but the font is now relocated.
The additional test/package failure made 1,555 entries rather than the known
1,553. Preserve native-qualification-static-failure, including its exact log.

The fixture now explicitly borrows the 256-unit header plus two raw guard units
(the mapped hole after the header), and separately checks the actual font global
through its owned pointer. This preserves the physical-boundary check and adds
actual-global preservation. Static-map rules remain unchanged. Only this tagged
test file differs from the completed accumulated/server/highres source; production
and binary hashes are unchanged. Recheck all eighteen display tests in all three
targets, rerun the full asset comparison, verify the existing production binaries
and run fresh gameplay. The fixture-only delta and both source manifests document
why production builds and broader targeted runs need not repeat.

A focused-command typo used a nonexistent focused-pattern.txt; it failed before
Go discovery. The final recheck uses the existing focus-pattern.txt. No skipped
or unstarted run is counted as verification.

## Final qualification

All 2,478 frozen results / twelve groups match without hash changes. Independent
case coverage grows from 735 to 773 with 32 localization seeds and six long-label
inputs. Accumulated default passes 839 tests in 430.745s; focused spell-effects
passes 18 in 175.179s, adding the new GC regression for 840 distinct default tests
across those runs. Server passes 278 in 201.582s; highres 279 in 125.110s, both
including that regression. The tracked accumulated selection now includes it.

After the sole bounds-fixture correction, all eighteen display tests pass again:
default 27.179s, server 26.343s, highres 27.236s. Final production source is
identical to that used by the earlier broad tests and three builds. Final binary
hashes/interfaces were reverified; builds took 67.187s / 9.688s / 66.359s.
All three are ELF32/i386/SSE2/CGO, retain thirteen display interfaces, retire three
and contain neither general porttest helpers nor the new collection callback.

The full asset suite again has exactly the established **1,553 failure entries**,
with **15 pass / 3 fail / 32 skip** packages and no added/removed failures. Fresh
seeded gameplay matches all four C-reference screens in **38.846s**, with reference
override disabled. All 1,559 final source fingerprints remain unchanged. Every
job, including failed attempts, was joined before source edits.

C is **84,237 / 97 files / zero reference C** (−1,162). Cumulative frozen capture
coverage is **453,359 results / 1,156 groups**; repeated forced-GC runs do not
inflate those counts. The earlier hallway mismatch did not recur; its unresolved
cause and automatic failure capture remain documented in CLIENT_INVENTORY.md.

Final evidence: native-qualification.json, default-coverage-union.json,
fixture-only-qualification-delta.json, native-final-source-fingerprints.json,
native-binary-verification.json and baseline/runs/client-inventory-display-port.
The completed run has a verified asset deduplication/restoration manifest. The
completed-development-compression.json manifest preserves 114 older captures as
verified gzip files, reclaiming 348,006,688 bytes. Original assets are unchanged.
