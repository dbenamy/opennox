# Native item use, drop and pickup identities

## Original baseline and scope

All 203 selected roots pass without skips in default, server and high-resolution
profiles against `8c9b3fc6` plus three test-only files. Production fingerprints
match that qualified checkpoint exactly; its production, ABI, known-suite and
save/load evidence supplies the production baseline.

The first 181 roots were supplemented by 22 PlayerControls roots after review
found the team-flag return route calling the CrownPickup wrapper. Both runs have
identical source fingerprints; all 203 names are required after conversion.

Replace 41 distinct C addresses: fourteen Use, eleven Drop and sixteen Pickup
registrations. Retire their 41 exports and two adjacent Go-only inventory fixture
adapters. Shared owners retain distinct registration identities. Preserve data
sizes, nil AmmoUse/BowUse registrations, ten mutable hooks and six script-facing
Use getters, plus the existing Drop getter.

New contracts cover all keys stored in foreign ObjectType/Object records across
GC/stack growth, late replacement of all mutable hooks, getter identities and
full signed Use return words through the actual item-use path. The configurable
independent C observer defaults to its previous result of one. Existing assertions
and captures are unchanged. Focused roots include inventory, item effects, world
collisions, book awards, scripts, player files, unit gameplay, orchestration and
resource definitions.

The typed Use.Get API returns bool, but effectsUse consumes the full C integer.
Add a native int32 result registry and explicit result/discard methods. Unknown
addresses retain the original integer or void C calling convention respectively.
Food pickup still discards the result; original caller gates and nil behavior stay
unchanged. Keep object layout and the existing bool cache. This reversible API
choice requires a full default port corpus sweep after focused qualification.

Primary owns API design, baseline acceptance and fixture migration; one GPT-6
Luna helper drafts the bounded production overlay. Review exact keys, owners,
input narrowing, result widths, sparse capture IDs and remaining references before
qualification. Add independent native/raw dispatch contracts after conversion.
Run all focused profiles, full default port corpus, safe/static checks, three
production/ABI builds, exact known-suite comparison, headless save/load, source
fingerprints and original-asset hashes before accepting the conversion.

[Original evidence](item-identities-baseline.json),
[commands](item-identities-batch.json), [selection](item-identities-tests.txt).
Local artifacts: `build/port-item-identities/`.


## Qualified conversion

All 205 focused roots pass in each default/server/high-resolution profile without
skips: the 203 frozen original roots plus two independent native/raw Use contracts.
The complete default port corpus passes 2441 roots, with only the established
TestMapPopulationPrerequisiteProbe diagnostic skip (2442 selected roots total).
Safe/static checks, three production/ABI builds and headless character creation,
save/load/resume pass. The ordinary asset suite matches the exact known results:
304 failure events, 17 passing/two failing/32 skipped packages. Existing assertions
and frozen captures remain unchanged; original 1,654 asset hashes match.

All accepted phases have identical source fingerprints; all 22 changed/new/deleted
source files match review. All 41 exports are absent; retained exported function
signatures and bodies are unchanged. Cgo files fall 213→208 (255/463 eliminated);
selected legacy exports fall 866→825 (1,065/1,890 retired). Embedded production
callback bodies remain 77; headers remain 157 files with 3,554 physical lines.
Standalone production/test-reference C remain zero. External bindings are unchanged.

Luna supplied the bounded production overlay and reviewed primary fixtures/API.
Primary corrected GoldPickup's bool handling, removed an unused import and the
now-empty resource export file, and found the additional player-control caller.
That caller's original tests were added before installation. The first compile
caught primary fixture cleanup removing a local used by an existing pointer-return
check; it was restored without changing expectations. Accepted qualification uses
contracts-fixed. Review remains necessary; no measured cost savings are claimed.

[Qualification](item-identities-qualification.json),
[inventory](item-identities-inventory-after.json).
