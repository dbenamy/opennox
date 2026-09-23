# Exact damage callback values

## Scope and decision

Remove registered Go→C→Go round trips from five raw damage callers, using a
separate exact-int32 registry for the eleven canonical callbacks. Preserve the
existing boolean APIs and C callback addresses. Unknown callbacks still use C.
Boolean-only registrations do not imply an exact integer implementation: their
raw callers retain the original C callback. Projectile consumers include both
full-word and low-byte conditions; a result of256 must not become1.

The typed method takes object pointers, preserves signed32 argument
words and return values, and keeps all three objects alive across dispatch. Raw
callers retain their configured-callback precondition. Existing boolean nil-slot
behavior remains unchanged. The reviewed eight-file production patch is applied after baseline97f055a5;
conversion qualification is complete.

## Qualified original baseline

The original projectileDamage owner is exposed through a porttest-only wrapper.
Raw forwarding covers1800 combinations of nil/distinct/aliased source and weapon,
signed amount/kind boundaries and eight exact results including-256,255,256.
Temporary boolean-only registrations independently demonstrate that both true
and false overrides affect the boolean API without replacing the raw C result256.

A sibling of the existing registered-owner corpus enters projectileDamage with
actual canonical C addresses. Its2071 cases cover all eleven callbacks, generator
thresholds, source/weapon/game modes, positive HP changes and golem wrapping.
Existing boolean and direct-C expectations remain unchanged. The first original
run passed all seven new roots. Fifteen captures totaling2071 cases are frozen;
all132 roots pass in default/server/highres with no skips and exact frozen results.
Source/dependency fingerprints confirm only five porttest files differ from
death conversion89c91e08. All four preceding binary hashes match and retain the
eleven callback exports. Production qualification is reused for this test-only
baseline; conversion will require fresh binaries/scenarios. See
[baseline qualification](damage-values-c-qualification.json).

Additional caller coverage includes spell projectiles, sustained spatial/lightning,
object collision, sentry contacts, lightning modifiers and temporary cloud owners
and candidate callbacks. Cloud owner fixtures alone use some ineligible classes;
the selected TestTemporaryAreaCallbacks matrix includes class2 targets, both cloud
callbacks and clear/blocked ray cases. Do not infer callback execution solely from
an owner name.

## Delegation review

Luna drafted the production changes. Primary rejected the first patch before
application: it removed imports but omitted all five promised callsite edits.
The original draft and an erratum remain in build/port-damage-values. The corrected
patch composes complete per-file changes and uses typed object-pointer arguments.
Primary owns baseline design, review, integration and qualification; patch
applicability and reported replacement counts alone are not acceptance evidence.

Standalone C remains zero physical files/lines; production C preamble bodies79.

## Artifact recovery

Twelve older safe/headless/server binaries from reward, damage, mixed and pointer
batches were gzip-archived after qualification-hash and host-use checks. Raw
582940300bytes became274093337bytes. Restore with
`python3 build/port-damage-values/archive-old-binaries.py --restore` before rerunning
historical finalizers. All current death binaries remain plain. Two bounded cache
pages found no retired-marker candidates and removed nothing.


The marker-specific cache scans reclaimed nothing. Primary then explicitly chose
ordinary reproducible-cache eviction:40 repository compiler archives older than
2026-09-23T20:00Z,1882456060bytes, removed after independent stat/hash/archive/package
and host FD/executable/mapping checks. This is not a claim that retired C symbols
prove those archives obsolete. Git history, module cache, outputs and current
binaries are preserved; historical rebuilds may take longer. Exact records are in
build/port-damage-values/old-repository-cache-removed.json. Pruning ran while no
Go/build jobs were active, before conversion qualification.

First conversion discovery88073 failed before tests: an unused unsafe import in
world_motion_sentry.go remained after call replacement. Primary and helper static
review missed it. Removed only that import; failed logs retained under
contracts-v1-failed. No golden or production behavior correction was needed.

Corrected conversion89247 passes all132 roots in each of default/server/highres,
with all15 original capture hashes unchanged and no skips. Remaining pipeline47464 joined PASS: safe/static, four fresh386/SSE2/CGO binaries,
retained Go-backed exports, exact known-suite comparison (304 existing failure
events;17 pass/2 fail/32 skip packages), headless creation and save/load. Production
contains no PortTest symbols. See [conversion qualification](damage-values-qualification.json).
Scenario asset duplicates were removed with hash checks and restore manifests;
original assets remain unchanged. All finalizers/cleanup scripts are consumed.
Check Git for conversion commit/push.

Next: assess a larger object-update callback batch using the53-entry inventory.
Six mutable handlers must be resolved at invocation time. Existing algorithm
fixtures need actual registered-owner siblings before production conversion.
Standalone C remains zero files/lines; production C preamble bodies remain79.

The copied preflight display label and production description said object death
registry during execution. Only those descriptive labels were corrected after
qualification; runner snapshots preserve exactly what ran. Selection, source,
binaries and results were unaffected. The documentation script first stopped
because it expected one label occurrence; the two-field metadata correction was
then completed explicitly, without rerunning or duplicating its earlier updates.
