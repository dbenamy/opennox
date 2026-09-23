# Exact damage callback values

## Scope and decision

Remove registered Go→C→Go round trips from five raw damage callers, using a
separate exact-int32 registry for the eleven canonical callbacks. Preserve the
existing boolean APIs and C callback addresses. Unknown callbacks still use C.
Boolean-only registrations do not imply an exact integer implementation: their
raw callers retain the original C callback. Projectile consumers include both
full-word and low-byte conditions; a result of256 must not become1.

The proposed typed method takes object pointers, preserves signed32 argument
words and return values, and keeps all three objects alive across dispatch. Raw
callers retain their configured-callback precondition. Existing boolean nil-slot
behavior remains unchanged. No production patch has been applied yet.

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
