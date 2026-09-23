# Damage callback registry

## Scope and compatibility

The proposed conversion adds a typed Go dispatch registry for eleven existing
object damage callbacks. Their registered C addresses, names and object layout
remain unchanged. Existing RegisterObjectDamage stays available; unregistered
callbacks still use the five-word C dispatcher. DamageFunc takes three typed
object pointers and signed32 amount/kind, returning normalized bool. No damage
algorithm or registration synchronization changes are intended.

Three special C-export wrappers share typed Go helpers with the registry:
MechGolem doubles amount for kinds9/17 with int32 wrapping; Flammable substitutes
9999999 for kinds1/7/12; BlackPowder rejects other kinds except0/1/2/12 and substitutes
999999 for accepted kinds. The BlackPowder C export retains its early rejection
before integer-address conversion. Ball remains the existing always-false callback;
its zero-argument C export remains available. All eleven C exports must remain
Go-backed in all four production binaries. Damage-sound registrations are unchanged.

The initial nil-slot gate suppresses argument adapters. The configured original
compiler evaluates source then weapon adapters, once each, before selecting the
callback from the object's slot. Preserve that observed ordering. Typed pointers
remain alive through raw fallback with runtime.KeepAlive; storage ownership and
pointer layouts are unchanged. Target remains386/SSE2/CGO; no64-bit claim or measured
performance claim.

## Baseline and independent contracts

The existing damage fixture's C-switch route and all its frozen expectations remain.
An additive Registry option installs the real registered callback address for the
whole case, verifies its exact name/address against the C fixture, and calls the
actual Object.CallDamage owner. It does not temporarily replace/restore callback
slots to force agreement with the previous route's hashes.

The new frozen corpus has2071 actual-registry owner cases:1800 general cases across
ten callbacks,144 generator boundary cases,108 source/weapon/mode cases,4 positive
health/return cases and15 signed golem wrapping cases. Generator uses valid existing
update storage and frame/health boundaries. Nine further rejected BlackPowder kinds
exercise the original C export with nil objects. All16 captures total2080 cases.
Independent positive cases require health50→38 and true; recorder tests check3125
pointer/amount/kind/return combinations, typed nils and custom object adapters.
Nil-slot and mutation tests check adapter suppression, source-before-weapon order,
once-only conversion and late callback selection.

The selected92 roots include the full Damage and Attack families, AI combat,
collision/projectile consumers, shock/area/shield spells, and object/player death.
Primary added TestSpellLifecycleProjectilesAndShock after Luna's bounded audit did
not identify its shock cases. This is focused owner coverage, not a claim that every
poison/harpoon/death-ball branch executes. Full known-suite and headless scenarios
remain conversion gates. Test-only baseline production evidence can be reused only
after all non-test fingerprints and four preceding reward binary hashes match.

Original capture73267 passed on unchanged production. Sixteen new hashes were frozen
before strict three-profile qualification. Original capture39792 failed because the
first independent ordering test assumed an early callback read: argument adapters
cleared the slot before the actual read, yielding a nil C call. The failure is kept
in ignored artifacts. Primary's early-capture review request was incorrect; Luna's
initial late-selection draft was right. The corrected test clears the slot in the
source adapter and installs callbackB in the weapon adapter, requiringB without
intentionally dispatching nil. No production fix or old golden change was made.

## Delegation and review

Luna drafted the implementation and raw forwarding fixture, audited owner tests and
identified obsolete cache archives. Primary owns the actual-owner baseline, final
selection and acceptance. Review caught BlackPowder's early gate and a fixture type/
function name collision. Execution also corrected primary's callback timing advice.
This supports bounded delegation plus independent validation, not blanket approval
of either model's reasoning. Model cost/speed savings remain unmeasured.

## Cleanup and C count

Two bounded cache pages found55 obsolete archives totaling2,488,347,446bytes.
Primary independently checked single-link regular-file metadata, archive headers,
hashes, retired bridge markers, marker absence from source/current four binaries,
and host FD/executable/mapped-file references before removal. Per-page records under
build/port-damage-registry prevent repeats. Source, assets and current binaries are
preserved. Future audit output publication should use atomic replacement.

Standalone C remains zero physical files/lines; production C preamble bodies remain79.
This batch removes registered Go→C→Go round trips, not those fallback dispatcher bodies.

## Qualification status

Baseline20b211e2 is committed/pushed. Conversion15863 and full pipeline5214 joined
PASS: all92 roots per profile,16 frozen captures2080cases and3125raw forwarding
cases plus nil/order contracts; no skips, changed expectations or post-baseline
source corrections. Safe/static, four fresh386/SSE2/CGO binaries, all11 retained
Go-backed C exports, exact known-suite comparison (304 existing failure events;
17pass/2fail/32skip packages), headless creation and save/load pass. Production
contains no PortTest symbols. See [baseline](damage-registry-c-qualification.json)
and [conversion qualification](damage-registry-qualification.json).

The copied preflight display label still said reward-use during execution. Only
that descriptive manifest label was corrected afterward; original commands and
records remain unchanged and the correction is explicit in qualification JSON.
Completed scenario asset copies were deduplicated with hash/host-use checks and
restore manifests. Application/finalizers/dedup scripts are consumed.

Next: review an unapplied object-death registry draft (14 names/four callers).
It requires a stable actual-owner baseline before conversion; no field or address
retirement is proposed. Damage-sound has no identified in-repository Go field
invocation, so it is not the next optimization target.
