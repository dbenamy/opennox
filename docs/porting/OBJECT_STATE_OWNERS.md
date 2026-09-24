# Object-state owners and export retirement

## Original identity audit

Qualified starting revision: `6274a6f3`. Candidate: retire 37 object-state C
exports after moving their fixture and 19 production Go caller files to native
owners. Four live collision exports remain. Conversion is not installed yet.

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

## Pending qualification

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
and [test selection](object-state-owners-tests.txt). Conversion qualification remains pending.
