# Object-state owners and export retirement

## Original identity audit

Qualified starting revision: `6274a6f3`. Candidate: retire 37 object-state C
exports after moving their fixture and 19 production Go caller files to native
owners. Four live collision exports remain. The conversion is fully qualified.

The original fixture registers 41 C function addresses. A temporary opt-in probe
observed all raw-word normalization paths, including the outer AI fixture's map
alias, Field2, AI-main and spawn-player normalization. Object-address-only lookups
remain separate from raw callback words. It changes no normalized return value.

Two fresh default-profile processes each pass the exact 20 selected roots with
unchanged frozen hashes. Each records 2,811 object-state setups, 3,352 actions,
all 44 operation IDs and 48,699,471 normalization lookups. All 41 function addresses
are distinct; none is consumed by normalization. Coverage and generated-ID
counters agree between runs. Three raw map-size diagnostics differ in the loot
cases, while all frozen outputs agree; raw map sizes are not claimed deterministic.

Map cardinality is nevertheless observable when shop nodes receive IDs from
`90000 + len(p.ids) + reservedFunctionIDs`. Preserve the 37 retired registrations
using the existing reservation counter (2 to 39), then compare frozen outputs.
This addresses both lookup use and map-size effects; a reference search alone was
insufficient. No live callback identity is retired by this decision.

The probe has been removed from working source after its two successful runs.
[Probe results](object-state-identity-probe.json) and the
[exact instrumentation patch](object-state-identity-probe.patch) preserve the
experiment. Apply the patch to `6274a6f3` with `git apply --unidiff-zero`, use the documented Linux 386 environment
and prior native-owner batch asset variables, then run the report's test pattern
with `OPENNOX_OBJECT_STATE_IDENTITY_PROBE` pointing to a fresh JSONL output.
Local build/command/binary records are in `build/port-object-state-owners/identity-original/`.
The temporary probe is test-only; production source was unchanged.

## Qualification

The next affected-owner selection includes 283 roots in each profile: object
state, collision registries, AI combat/callback/main/monster state, damage,
equipment, attack, player death, projectiles, object death, spawn policy,
generators, spell effects, effects, shop and trade. Run that selection on original
and converted source, then safe/static, three fresh production/ABI builds, exact
known-suite comparison and headless creation/save/load/resume. Assertions and
frozen hashes remain unchanged. Source and import/export counts are updated only
after conversion qualifies.

All three fresh original-profile baselines pass exactly 283 roots with no skips.
Source fingerprints match qualified `6274a6f3`. See the
[baseline](object-state-owners-baseline.json), [commands](object-state-owners-batch.json)
and [test selection](object-state-owners-tests.txt). Conversion qualification passes with the same exact root sets and no skips.

All converted gates pass: 283 roots in each profile, safe build/static checks,
three production builds/ABI checks, exact known-suite comparison, and headless
creation/save/load/resume. All 37 retired C symbols are absent; the four live
collision exports remain. Every phase has identical source fingerprints and all
1,654 original asset hashes remain unchanged. No frozen expectations changed.

The final change spans 22 files and removes 37 exports/prototypes plus the
test-only C stateCall dispatcher. Selected cgo files stay 230 (233/463 eliminated);
C exports fall from 1,179 to 1,142 (748/1,890 retired). Production C callback bodies
stay 78. Headers remain 157 files, now 3,865 physical lines. Standalone production
and test-reference C lines remain zero. See
[qualification](object-state-owners-qualification.json) and
[inventory](object-state-owners-inventory-after.json).

Primary reconstructed all 22 files from exact edits/diffs. Review corrected a
fixture draft's Pointf import and explicit float32-to-float64 mass widening before
installation. Public wrapper signatures, temporary point allocations, signed
char results and callback identities are preserved. The equipment sync caller
uses Sub_4E4500(..., true) plus stateSyncEnd; stateSync is not interchangeable
because it computes its boolean argument differently.
