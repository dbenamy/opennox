# Live empty callback identities

Status: actual-C baseline qualified after the tagged spell-fixture repair; no
callback conversion is applied. See
[empty-callbacks-c-qualification.json](empty-callbacks-c-qualification.json). The preceding safe-export conversion is pushed as `e663a971`. Standalone
production C remains 25 physical lines in three files; no reference C is added.

Ten empty callbacks remain live. Four damage hooks use five arguments, three
defense hooks use six, two update hooks use three, and Energy Bolt destruction
passes one argument. The legacy definitions ignore those arguments; the duration
callback explicitly declares void(void). Any conversion must preserve the current
386 calling convention and all ten distinct C-callable identities. Readiness and
Replenishment addresses are also compared as effect identifiers; they are not
merely pointer tokens, since attack dispatch invokes them too.

The declaration-only tagged fixture exposes the actual C pointers. A separate
read-only helper exposes the existing modifier registration maps. The root test
independently checks names, effects, arities, exact registration/getter identity,
non-nil and pairwise distinct addresses. It invokes all 338 combinations of
nil/non-nil arguments through the existing ccall paths and checks complete guarded
foreign records remain unchanged. No callback implementation is copied into tests.
Three actual-C processes produce the same capture hash:
`848b76f173284c29edddd5d637863061d643428004ed8f51098f51863ea44b77`.

The selected suite contains that contract plus 23 existing roots covering effect
identity/dispatch, melee effect callbacks, damage/equipment slots and sustained
spell lifecycle/cleanup. Existing consumer roots complement the direct callback
cases; they do not each prove an exact empty callback was invoked. The final baseline passes all 69 broadened roots in default/server/highres and
all 24 callback consumers under safe, without skips.

The benchmark measures the production Go → ccall → C indirect-call path, with
the same closure overhead before and after. Initial five-sample C medians ranged
from roughly 73 to 151 ns/call. Those preliminary results do not predict game-frame
cost or actual callback frequency. A retained benchmark binary with source
fingerprint supplies the qualification comparison. No timing threshold is added
to ordinary correctness tests. Do not accept a callback conversion merely to
reduce the C line count; inspect absolute overhead and variation first.

Luna supplied a useful reference audit, then drafted bounded identity/registry
helpers. Primary verified mappings and wrote the independent root contract and
benchmark. Primary also corrected the timing proposal: a C-only loop would isolate
a boundary but omit the Go-to-C dispatch present in production. Only one helper
is used, and primary owns application, source freeze, qualification and acceptance.

## Safe-profile fixture prerequisite

The initial 24-root suite passed default/server/highres. Safe passed eight roots,
then `TestSpellLifecycleDurationLists` panicked in fixture setup: its 18-word
cache range beginning at 1569676 crossed the registered duration allocator/list
slots at 1569724, 1569728 and 1569736. This happened with the original C callbacks.
The remaining fifteen sustained-spell roots did not run in that failed process.

The fixture now explicitly saves, clears and restores its fourteen type-cache
slots, excluding the duration-state gap. Production and golden hashes are
unchanged. The first setup repair reached an equivalent stale snapshot loop and
failed there; that second access is now corrected too. Four historical zero
columns remain in the capture schema for the retired-state gap, while the live
duration list and records are still captured separately. Neither failed safe run
counts as qualification. Follow-up qualification expands to
69 roots: the callback suite plus every root in the eight files that explicitly
construct spell-lifecycle fixtures. Preserve the failed run under `c-safe`;
reviewed runs use new `c-reviewed-*` directories and the reviewed test pattern.
Only tagged fixture sources differ from prior production; baseline evidence reuse
must verify the modified fixture, as well as the three new fixtures, is excluded
from every production build selection.

The expanded safe run then stopped earlier in `TestClientPresentationChant`:
`PortTestNewEffectsEnvironment` accesses registered offset 1313532 through memmap
at `client_effects_environment_porttest.go:76`. That unrelated renderer fixture is
outside this callback change; its failed run remains recorded, and no broad safe
success is claimed. Acceptance requires all 69 roots in default/server/highres
and the original 24 callback consumers under safe. This limitation should be
reviewed when extending optional-safe support for the renderer fixtures.

## Accepted baseline evidence

The reviewed runs pass 69/69/69 roots in the three normal profiles and 24 under
safe, with all discovered roots completed and no skips. The new 338-case capture
and every existing selected golden remain unchanged. All production source files
are unchanged: only three new tagged fixtures and the repaired tagged fixture
differ from the preceding source manifest. Official Go file selections exclude
all four from default/server/highres/safe production builds. The preceding four
production binary hashes and ten callback addresses per binary were rechecked
before explicitly reusing their ABI, known-suite and gameplay/save-load evidence.

The retained C benchmark binary and its source fingerprint are recorded. Compare
it with a retained Go-export binary using interleaved fresh processes after other
agent builds and archives finish. Do not label a measured increase performance
neutral, or infer real game-frame costs from these boundary timings.

## Qualified Go exports

Ten distinct Go exports replace the C empty bodies without changing their names,
no-argument C signatures, registrations or existing argument-passing dispatch.
All 338 guarded cases match the frozen C hash. Default/server/highres pass all
69 consumer roots each, and safe passes all 24 callback roots without skips.
Safe build/static, three fresh production builds/ABI, headless character creation
and explicit save/reload/resumption pass. The full suite matches the established
1,553 failure entries and package outcomes exactly; it is not a green full suite.
All four production binaries contain ten distinct callback addresses backed by
Go exports, and production builds contain no PortTest symbols. No golden changed.
See [empty-callbacks-native-qualification.json](empty-callbacks-native-qualification.json).

Five fresh processes per version, in alternating order, compared retained C and
Go benchmark binaries after unrelated builds/archives finished. Median per-call
cost increased by 83.6–139.8 ns across the ten callbacks (1.57–2.27 times the C
cost). Primary accepts this small, reversible compatibility cost; the conversion
is not performance-neutral. Actual call frequency and game-frame impact remain
unmeasured. If profiling identifies a material cost, consider bypassing known
empty callbacks within Go while preserving compatibility addresses, or revert
this small batch. Neither follow-up optimization is implemented.

Removing GAME5_2.c and common__object__modifier.c reduces standalone production C
from 25 lines/three files to **six lines/one file**, with zero standalone reference
C. The remaining file includes minimp3; this is not a C-free build. A lexical
inventory also finds 86 production C preamble bodies: 76 generic dispatchers and
ten small adapters. Its 51 decoder-header bodies include inactive conditional
branches and must not be interpreted as 51 compiled functions. Generated bridges,
C types and external libraries remain outside the standalone LOC metric.

The optional-safe renderer-fixture limitation above remains unchanged. The tagged
spell-fixture repair stays covered in the normal profiles and selected safe roots.
Scenario asset copies were deduplicated only after successful runs and matching
original hashes; per-run restoration manifests remain under build/baseline/runs.
