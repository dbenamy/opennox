# Collision callback dispatch

## Scope and compatibility

All 53 named collision registrations now bind complete existing Go export wrappers.
Callback addresses, sizes, parser order and three default aliases remain intact;
three replaceable handlers are still resolved at invocation. The original baseline
was captured at lifecycle checkpoint `72b9c14f` and committed as `ee9ab4d0` before this
production conversion.

The integer-word CallCollide API keeps its nil no-op and unknown callback fallback.
CallCollideWith preserves the four queue/projectile pointer-call owners, their raw
three-pointer fallback and configured-slot precondition. Object/target/normal escape
to heap according to exact compiler diagnostics; KeepAlive and stack-growth/GC
contracts qualify lifetime through conversion to callback words. Only the arguments
each complete wrapper declares are converted. UndeadKiller's third word retains its
nonzero-condition meaning; normal consumers retain their pointer interpretations.

The temporaryMagicMissile owner consumes an integer return and stays raw. Two
damage-sound callbacks remain for separate nil-default/owner qualification. These
are reversible scope decisions under standing authorization. No layouts or callback
identities are retired by this batch.

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
before installing production conversion. Final selection and results are recorded
below.

## Delegation and review

One GPT-6 Luna helper produced the ignored 53-binding draft and execution-coverage
map. Primary independently parsed every current export signature and checked each
argument conversion and ordered name/address/size triple. Primary requested removal
of an unnecessary boolean normalization from UndeadKiller; the wrapper already
performs that check. The corrected mapping passed static review before runtime
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
observations; all 257 exact affected roots subsequently passed in three profiles.
Primary owns lifecycle/pointer design, new contracts, integration and qualification.
No measured subscription-cost or time saving is claimed. Drafts and review evidence
are under ignored build/port-lifecycle-registry/next-collision-*; captures and new
qualification artifacts use build/port-collision-registry.

Standalone C remains zero files/physical lines. Production C preamble bodies remain
79; the conversion still retains raw callback ABI glue and C-visible identities.

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
At that baseline checkpoint the production patch was unapplied. Compiler escape
checks and conversion qualification were completed afterward, as recorded below.

## Completed conversion qualification

The five production files now bind all 53 registrations and route four queue/trace
pointer calls through CallCollideWith. Generic CallCollide retains its integer-word
API and nil/raw behavior. Only the test-only original pointer bridge switches to
the new method; all cases and frozen expectations remain unchanged. Compiler
escape diagnostics establish heap lifetime for object, target and normal arguments;
the exact diagnostics and source hashes are recorded in qualification metadata.
Dynamic-handler tests retain identity and normal mutation through stack growth/GC.

All 257 affected roots pass in each profile with all 72 frozen capture groups and
exact 4,163 counted calls unchanged. Safe/static, four fresh production binaries,
retained Go-backed exports, exact known-suite comparison, headless creation and
explicit save/load/resume pass. See [conversion qualification](collision-registry-qualification.json).
Known-suite outcomes remain 304 failure events and 17 pass/2 fail/32 skip packages.
The configured target remains 386/SSE2/CGO; no 64-bit or performance claim is made.
Standalone C stays zero files/lines; 79 production preamble bodies remain.


The later [native collision identity batch](COLLISION_IDENTITIES.md) adds a typed
result route and migrates `temporaryMagicMissile` through it. Declared native
return bits and the existing independent raw-C return contracts are preserved;
the earlier requirement to leave this caller raw is superseded.
