# Collision callback dispatch

## Scope under qualification

The candidate binds all 53 named collision registrations to their full existing Go
export wrappers. Callback addresses, sizes, parser order and three aliases of the
default callback stay intact. Three replaceable handlers remain resolved when called.
Production conversion is not installed yet; new tests are being captured against
the original dispatch at lifecycle checkpoint `72b9c14f`.

Keep the existing integer-word CallCollide API, including its nil no-op and unknown
callback fallback. A separate pointer-argument entrypoint will preserve the four
queue/projectile dispatch calls, their raw three-pointer fallback and their original
configured-slot precondition. Normal-vector storage must remain live and escape Go
stacks before conversion to integer words; require compiler escape evidence as well
as stack-growth/GC and mutation contracts. The UndeadKiller third argument is a
nonzero condition, while several other wrappers interpret it as a normal pointer.
Only convert the arguments each complete wrapper declares.

The temporaryMagicMissile owner consumes an integer return. It stays on its existing
raw call in this batch; a void callback API cannot supply that contract. The adjacent
two damage-sound registrations stay for a subsequent batch with their own nil-default
behavior and owner coverage. These are reversible scope decisions under the standing
authorization, not claims that either remaining path is unnecessary.

## Baseline design

Retain the existing algorithm tests and expectations. Add registered-route siblings
for projectile, world, objective, state, trigger, pickup and spell-award owners.
Record actual registered invocation counts so a selected matrix cannot silently
exercise only unrelated helpers. The new void-owner observations discard only
returns the owner ignores; original direct-helper return tests remain unchanged.

Independently assert all registration addresses and sizes, default aliases, raw
integer/pointer forwarding, cleared/nil slots, successive handler replacements,
and object/normal identity and mutation during GC and stack growth. Existing
collision-queue and projectile-trace fixtures cover ordering, reciprocal normals,
wall contacts and raw observation callbacks. Freeze original-path observations
before installing production conversion. Final affected-family selection and
qualification results are pending.

## Delegation and review

One GPT-6 Luna helper produced the ignored 53-binding draft and execution-coverage
map. Primary independently parsed every current export signature and checked each
argument conversion and ordered name/address/size triple. Primary requested removal
of an unnecessary boolean normalization from UndeadKiller; the wrapper already
performs that check. The corrected mapping passes static review, not yet runtime
qualification. The coverage audit distinguishes real registry calls from direct
export calls and does not count captured callback addresses as execution.

Luna drafted 14 bounded projectile siblings without changing allocation ownership.
Primary checked all 23 op/name mappings against the original C switch and compared
each sibling body with its original matrix. One Positive capture call still used the
old hash map and bypassed invocation-count reporting; primary corrected it before
compilation. Original-path capture 15368 completed all 43 roots; only the primary-added
SpellWall sibling failed because its invocation counter was drained twice. The
state assertions passed. Primary removed the redundant counter check, retaining
all behavior assertions and rejected artifacts. Corrected capture 46948 passed all 43 roots. All 72 capture hashes and 4,163
invocation counts matched the first run. Frozen expectations cover 4,116 owner
observations; 257 exact affected roots are now qualifying in three profiles.
Primary owns lifecycle/pointer design, new contracts, integration and qualification.
No measured subscription-cost or time saving is claimed. Drafts and review evidence
are under ignored build/port-lifecycle-registry/next-collision-*; captures and new
qualification artifacts use build/port-collision-registry.

Standalone C remains zero files/physical lines. Production C preamble bodies remain
79; the candidate still retains raw callback ABI glue and C-visible identities.

## Workspace recovery

With no Go jobs active, primary verified Luna's capped inventory of 30 old repository
compiler archives and reclaimed 1,297,507,410 bytes. Every entry was a regular,
single-link archive of at least 8 MiB, older than 2026-09-23 22:00 UTC, with archive
magic and the repository module marker. Stat/hash checks and host FD/executable/
mapped-file checks passed. This is reproducible cache eviction, not unreachable-code
proof; older rebuilds may take longer. Source, module cache, assets, qualified
binaries, captures and logs remain. Exact plan/removal records are in
build/port-collision-registry. Removal script is consumed by its completion record.

Primary also corrected one coverage-audit omission: TestObjectStateOwnership selects
operation 43, reaching the direct DefaultCollide wrapper. The new alias contracts
still add actual registered-name coverage for Default/Elevator/Telekinesis. Keep
fixture-operation tracing explicit rather than inferring coverage from root names.

## Qualified original baseline

All 257 selected roots pass without skips in default/server/highres. All 72 frozen
capture groups repeat 4,116 observations, with exactly 4,163 counted invocations
across 44 names. Three mutable handlers have successive-replacement contracts;
existing registered-path roots cover Barrel, AudioEvent, Pentagram, Sign, TrapDoor,
Teleport and MonsterGenerator (TrapDoor also has a new sibling). Together these
exercise all 53 registration names, including the three no-op aliases.

Nineteen new porttest files and three existing test adapters are the only source
changes. All other source/dependency fingerprints and all four preceding lifecycle
binary hashes match; all 51 callback export identities remain in those qualified
binaries. Production evidence is reused only for this test-only checkpoint. See
[original baseline qualification](collision-registry-c-qualification.json).
The typed production patch remains unapplied; pointer-API compiler escape checks
are still required during conversion.
