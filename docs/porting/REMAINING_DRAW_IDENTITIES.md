# Remaining registered drawable callbacks

Scope: 12 drawing exports, seven root-owner forwarders and five legacy screen
owners. Extend the existing native draw registry with distinct keys for debug,
monster, vector, released soul, animation state, player, NPC, harpoon, rope,
undead killer, maiden and player waypoint. Preserve the original 56 key indices.
The 25 update exports and screen-particle callbacks remain outside this batch.

## Behavior and baseline

The root-forwarding closures must read their hook variables at call time. Legacy
initialization precedes root hook assignment, and tests can replace those hooks.
Released soul shares the vector owner but needs its own identity. Harpoon shares
the slave renderer, and maiden calls the monster owner after palette handling;
their keys also remain distinct. Preserve signed 32-bit results and the viewport
word conversion on the supported Linux386 target. The debug key remains the
object-type default. No renderer algorithms or record layouts change.

The new `TestClientDrawRootCallbackDispatch` enters through actual callback keys
on the original C path. It checks seven distinct keys, the debug default, pointer
arguments, signed32 result boundaries, runtime hook replacement, int-observing and
void-discard callers, and nil arguments. It saves/restores all modified hooks and
uses owned native viewport/drawable allocations. Existing real renderer/pixel
contracts remain necessary; a stub hook contract only establishes dispatch behavior.

The conservative affected selection starts with selected identity getters, the
two affected fixture API families, and root renderer/parser owners. It follows
root test helper calls and direct seed references, recording all edges. It selects
763 default/highres roots and 753 server roots; all ten exclusions carry explicit
`porttest && !server` source constraints. The common draw dispatcher is unchanged
in this batch, so it is not an additional selection seed. Shared test helpers still
make this a deliberately broad selection, not a formal dynamic call graph.
Primary review and original-path acceptance are complete: exact name sets passed
in every profile without failures/skips, and the new contract passed twice per
profile on identical verified source/binaries.

The original baseline reused production qualification from `5687bc68`: production
source was identical, with only the new root test added. The new contract passed
separately and again within the complete affected selection on verified binaries
before the baseline was accepted and the conversion installed.

After conversion run the full affected selection in all profiles, a fresh client
preview, safe/static checks and all production/ABI/known-suite/save-load gates.
The immediately preceding 66-export batch qualified the common dispatch boundary
with the complete 2,460-root corpus. This follow-on only registers additional
owners, so focused qualification is appropriate unless new findings broaden the
change. No frozen captures or existing assertions may change.

## Delegation and recovery

Luna owns an ignored legacy-only migration overlay and exact change manifest.
Primary owns the new contract, selection, original-path acceptance, integration
and final qualification. Artifacts: `build/port-remaining-draw-identities/`.
The frozen scouting report is under the prior drawing batch's `luna/next/`.
Baseline `9f91acda` and the conversion preserve recovery; conversion fully qualified.

Primary reviewed all 12 owner mappings, the unchanged first 56 registrations and
indices, retained exports, and all selected identity consumers. The draft needed
two mechanical corrections before installation: preserve the adjacent `//export`
directive on the retained screen-particle function, and remove the remaining
animation-state header prototype. The corrected overlay contains 13 source paths;
all 68 identities are distinct. The conversion was installed after baseline acceptance and commit.


## Qualified result

All 763 default/highres and 753 server focused roots pass with exact baseline names,
no failures and no skips. Original-path root-dispatch contracts passed twice in each
profile; post-conversion contracts preserve dynamic hooks, arguments and results.
Safe/static checks, three production ABI checks, both fresh headless save/load
scenarios and exact known asset-suite comparison pass. All 1,654 original asset
hashes remain unchanged. See [qualification](remaining-draw-identities-qualification.json).

Selected cgo files fall from 161 to 159 in clients and 162 to 160 on server;
legacy exports fall from 399 to 387. Two production files and one fixture no longer
import C. The 157 headers now contain 3,149 physical lines;
77 embedded production C bodies remain, and standalone production/test C stay zero.
External native dependencies are unchanged. The previous full-corpus qualification
is explicitly retained as historical boundary evidence, not rerun or claimed for
this source.

Luna's bounded overlay was useful after the two documented primary corrections.
Every owner/key mapping, retained export directive and retired prototype was
checked before builds. No subscription savings estimate is inferred.
