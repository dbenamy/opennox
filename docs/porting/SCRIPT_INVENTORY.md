# Script inventory commands — C baseline qualified

Qualified parent: **5a755193**, 10,953 physical C lines /48 files /zero reference C.
Select the three remaining routines in `server__script__builtin.c`: carry-capacity
fallback/drop selection, startup cleanup, and SetHalberd replacement. The only
external callers are Go wrappers and the actual legacy builtin table. Two C
counters track reserved slots and whether the capacity notice has been sent.
Root pickup/reset callers use setters/increment; these can move with the batch.

Reuse the complete player/inventory controls fixture for lists, item factory,
RNG, equipment and reports. Boundary callbacks record drop/deletion requests;
independent contracts check the command's selection and order. Existing inventory
and equipment contracts cover the dependencies. Add fresh gameplay/save-load
qualification after conversion. Do not treat boundary callbacks as end-to-end
coverage of deletion or physical playback.

Startup's120 C cases pass: eight item-class combinations, five list orders,
three initial chapter values, journal mask0xE and preservation of other players.
The initial fixture activated only one player; the host slot31 precondition
requires all three sparse fixture players. Correcting the fixture fixed that
failure without production changes. The first capture is
`85dea812ebfb2d36ed1b2557432eb58c7b4aca6ea5a5dafd9a665342b688f369`;
repeat and freeze it with the complete C baseline before translating.

Carry/drop's864 cases pass: exclusions, first minimum/ties,
999999 price cutoff, capacity/reservations, two repeated calls, notice values and
successful/failed drop callbacks. The combined run now adds432 stack-capacity and112 halberd cases; its results
are pending. Startup repeated identically in the selection run. The misleading ItemIsDroppable name must not reverse this command's
exclusion test. Halberd lookup uses shipped table offsets247336 and its four
strings; verify nonzero shipped data before capture.

Evidence and advisory drafts: `build/port-script-inventory`. Startup's failed
initial run is retained. Installed advisory carry/halberd drafts are CONSUMED. No inventory production
code has changed, and no new
capture is frozen yet. The native script-binding conversion remains qualified.

## Final fixture review

The first112 halberd cases pass with a controlled equipment-admission fixture.
The stronger final set has224 cases: real zero-strength modifier definitions,
class admission accepted/rejected, actual insertion at the recorded placement
boundary, actual equip logic, initializer count/arguments, shipped names, and
first-match deletion requests. Deletion remains an observed boundary here; full
inventory deletion behavior belongs to the inherited inventory/integration gates.
The previous fixture only recorded placement and returned a result; it did not
insert the created item, so an independent holder assertion correctly rejected
that setup. The placement adapter now checks arguments1,1 and delegates to real
insertion. No production change was needed.

A Go interface-embedding name collision was fixed using the existing alias
pattern. A later compile was cancelled and fully joined when review spotted a
missing test-state field; its output is not qualification. Final focused C run
is `final-contracts-c-fixed.log`; static-final-c.log passes.

Production-source comparison finds only nine porttest file differences from
5a755193. All three reused binaries' hashes and selected11 C function/owner
symbols are verified. C baseline gates will reuse that parent's production,
full-suite and two scenario qualification, then run fresh affected target sweeps.
Native qualification will rebuild production and run fresh scenarios.

The two completed script-binding-native scenario copies have been deduplicated,
reclaiming1,112,747,701 bytes with hash-verified restoration manifests. That
cleanup mode is CONSUMED. Original assets/archive and all scenario evidence remain.

## Frozen focused C contracts

All four focused roots pass1,640 cases: startup120, carry selection864,
stack capacity432 and halberd224. Hashes are now literal test assertions and
manifest checks. Prior startup/carry captures repeat identically. Final hashes:

- startup: `85dea812ebfb2d36ed1b2557432eb58c7b4aca6ea5a5dafd9a665342b688f369`
- carry selection: `ee653fca8c51b397b2729f32e25af94bed91b24b0b8d400f94c35af9311e6bf2`
- stack capacity: `a59c3196c203e0adb5d5ce5db568c55734bf322b2bce74b0e0cdcae87e64c62b`
- halberd: `c33f3ac928be5f759e425de97f1a4a024882befdc26c83bceed4f299047b7a3c`

All-target affected qualification is running under c-affected-* with the tracked
manifest/selection. The selected three C bodies total94 lines. The native draft
and installer remain advisory and NOT CONSUMED; no production conversion yet.
Six additional now-private bridge exports can retire with the three commands:
capacity lookup, journal mask removal and four script stack adapters. Remaining
player-list C callers keep their live exports. Historical compatibility comments
are not executable references; that implementation remains unchanged.

## C baseline qualified

All jobs are joined. Primary selection passes279 default/highres roots and276
server roots, no skips. The shop supplement passes31 roots each. The initial
selection used ShopUI instead of the actual ClientShopUI prefix; the separately
qualified supplement closes that caller-coverage gap without repeating unrelated
checks. Native gates must run both selections.

The full frozen capture inventory is281 primary files on clients,278 on server,
and25 shop files on each target. Default/highres inventories are identical;
all server captures match their client counterparts. Three hover/world-selection
captures belong to files explicitly tagged `!server`; their absence is expected.
Runners enforced four new hashes plus inherited literals; full inventories were
compared and frozen afterward. All six source fingerprints agree. Static passes.

Production/full-suite/scenario evidence is reused from5a755193 after proving all
nine source differences are porttest-only and checking all three binary hashes.
The selected nine function symbols and two owners are present in every parent
binary. See script-inventory-c-qualification.json. No C correction or production
change was needed. C remains10,953 physical lines /48 files /zero reference C.
