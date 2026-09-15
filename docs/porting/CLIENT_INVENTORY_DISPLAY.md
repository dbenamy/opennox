# Client inventory display and feedback

Status: C baseline fully qualified; conversion follows this recoverable checkpoint.
Production remains 85,399 C lines / 97 files / zero reference C, at transaction
conversion `2f34230e` (pushed). No production display code has changed.

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
