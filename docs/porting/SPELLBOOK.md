# Spellbook UI — next batch audit

Status: C baseline development; production unchanged. Nine roots pass
with no skips (`c-render-owner`, 122.493 seconds); reward and animation
coverage is being expanded. No capture is frozen. Starting C: 68,597 lines in
88 files, zero reference C.

Selected scope: GAME2.c from 0045ABC0 through 0045D9B0, plus all five function
bodies in client__gui__guibook.c. This is 36 functions / 1,745 C block lines:
book sorting, initialization, visibility, page/window callbacks, reward/knowledge
updates and the animated addition of a book entry to the quickbar. Keep the
quickbar owner (0045D9D0 onward and guispell.c) as the following connected batch;
keep its actual callbacks/ownership during book testing. Do not add unrelated
sprite helpers before the sort comparator to pad scope.

Initial static word-reference audit is saved in `build/port-book/callers.json`
and `remaining-callers.json`. Fourteen book interfaces still have C references; 22
appear internal to selected C/Go users, but callback address uses and wrappers
must be rechecked before retirement. Existing GUI ports show how to replace
selected window callbacks with native Go callbacks while keeping actual GUI
ownership. obj_5d4594_1046620 is defined in guibook.c and used only by selected
book bodies; it can become private Go state when its C file is removed. Memmap's
metadata name is not evidence of a remaining C access.

Use actual GUI/window/image/font, player spell/ability/creature metadata, sound
observer, cursor, game flags and animation/RNG owners. Existing reusable fixtures:
client_inventory_window_owner_porttest_test.go (actual GUI and image references),
client_window_owner_porttest_test.go (deterministic real font), tooltip owner
(actual spell/ability/guide definitions), and shop UI owner (resource failure and
actual-asset loading patterns). Do not replace sorting/drawing/list owners with
facsimiles or expect a short warrior replay to cover spell/guide pages.

Independent contracts should cover empty/one/page-boundary/max lists, all three
classes, known/unknown entries, quest versus other game modes, guide view, hidden
spells, sorting case/order and reserved spell34; exact windows/capture/focus,
page turns and selection state; partial image-load failures and repeated lifecycle;
reward duplication/rank transitions; drawing pixels and font/width boundaries;
animation counters/coordinates/path insertion at 0/19/20 and screen-width tiers.
Freeze complete normalized state/effects plus decoded pixels for nonempty pages.
Use real shipped assets for initialization/rendering and extend replay input to
open/turn/close the warrior ability book (plus wizard/conjurer checks if practical).

Audit numerical edges before freezing: sort-list page divisor depends on font
height; guide and spell tables have fixed capacities; path callback writes up to
20 coordinate pairs; animation uses x87 intermediate precision and normalization.
Preserve successful original behavior and identify any actual prerequisite fix
with an independent contract before capture. No bug correction or baseline has been made for this scope yet.

The original archive/assets remain untouched. This is only a recoverable audit,
not a native implementation or test result.

Further audit: initialization loads eleven still images plus forward/backward
ImageRef animations, then creates the real book/child/arrow/icon windows. Use
separate owned ImageRefAnim records; copying only ImageRef aliases the mutable
OnEnd callback/timing state. Preserve every early image failure's partial globals.
The shared player arrays from objectRenderOwner are suitable; explicitly bind the
book's player pointer and restore all 32 named globals, vector and blob regions.
Destroy owned windows before freeing image references or restoring old globals.

Sort-list special/hidden spell flags can make a synthetic special-count exceed the
visible count if arbitrary contradictory definitions are supplied. Audit actual
spell flags and supported metadata combinations before selecting boundary inputs;
do not run a huge wrapped C qsort as an oracle. Use the real font owner and valid
page geometry (the divisor is 2*(141/(fontHeight+2))-2). Preserve duplicate-name
ordering observed in C. The existing string comparator/normalizer remains a live
owner; do not replace byte/UTF16 semantics with broad Unicode transformations.

The installed facade exposes all 36 actual C calls, callback identities, 32 named
word addresses and the two-float animation vector. Ignored facade generators and
drafts are stale; do not rerun them over the reviewed installed source.

The actual startup guide-family table at blob587000 offset132100 is
[24,7,8,25,26,0], followed by the pointer at132124 to that row and a null terminator
at132128. Read blob_587000.dat; copy the 32-byte region and rebase its pointer to
the actual memmap row in fixtures. This gives an independent reward/removal case
for guide24 propagating to 7/8/25/26. Creature enumeration uses real rows1..40 at
blob5D4594 offsets740076+28*id, with valid word+4 and image word+16. Do not use a
synthetic table when the actual startup family is available.

Manual facade review also added the selection word nox_xxx_aNox_cfg_0_587000_132136,
C dimensions and local player ID, bringing exposed named words to 32. The generated
facade has been enhanced; draft-facade.py is now stale and would overwrite those
additions. Bind players[0].Active/NetCodeVal and the local ID together to exercise
the actual active-player lookup in ability rewards. Only normalize pointer fields,
not arbitrary scalar/float bit patterns that happen to resemble an address.

Initial fixture evidence: `c-initial` failed before test discovery because the
local-player ID declaration used a cgo-inconsistent integer typedef; the actual
C declaration is unsigned int. After correction, `c-owner` passed all three roots
in 106.334 seconds. Initialization covers success and all thirteen missing resources;
the path contract covers 18 boundary/coordinate cases with independent float32
word expectations. Palette setup has its own root.

`c-sorting` exposed an incorrect fixture assumption about font height. The actual
renderer uses cap height: the selected font gives 11 pixels, a 13-pixel row and
18 entries per page. Corrected boundary inputs and expectations pass all four
roots in `c-sorting-capheight` (22.903 seconds). Spell sorting covers 198 cases,
including empty/max tables, hidden/special entries, duplicate names and game modes.
These are unfrozen development captures. No production correction was needed.
All applied initial/sorting scripts are stale.

Further unfrozen coverage: `c-metadata` passes five roots in 22.778 seconds,
adding 144 guide/ability list cases using actual metadata owners. `c-pages`
passes six roots in 32.455 seconds, including sequential forward/backward turns,
contents/detail transitions, direct selection and animation guards. Snapshots now
include known-entry arrays, capture/focus and decoded pixel hashes. Rendering,
visibility and rewards remain under development; this is not a complete baseline.

Recovery checkpoint: `c-render-owner` passes all nine selected roots without skips
in 122.493 seconds. This includes 1,182 page-transition records, 40 visibility
records, 30 rendering records and 120 spell-family knowledge cases. The rendering
fixture uses actual text/image drawing with wrapped names/descriptions, three
classes, guide/spell/ability pages and size variants. Full original-asset and
animation qualification is still outstanding; captures are deliberately unfrozen.

Development failures remain preserved. `c-visibility` selected the book itself
when looking for another capture window; the test now uses its known arrow child.
`c-render` found the tiny fixture images entirely clipped at the default book
position; the test moves the real window onscreen. `c-render-position` then exposed
a missing server backlink in the lightweight spell owner during mana-cost lookup;
`PortTestBookSpellOwner` connects the actual owner and restores it afterward.
No production bug was found or corrected in these attempts. The complete nine-root
run includes every test omitted after that earlier rendering panic.
