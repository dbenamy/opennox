# Generic object death callbacks

Scope: eight registered callbacks, GAME5 54DFA0–54E620 and both functions
in server__object__die__die.c. Together these occupy 265 physical C lines,
including the removable translation unit. MonsterGeneratorDie remains with
its generator/script/score family; PlayerDie remains a later owner.

## Original-C baseline

1,704 locked cases: eight smoke, 200 creation contracts, 32 marker contracts,
360 barrel contracts, 384 boulder contracts and 720 armor/weapon contracts.
Hashes live in src/object_death_porttest_test.go. A second independent run
matches all hashes and existing callback, initialization and penalty regressions
pass (7.659s). No production C is removed in the baseline commit.

The shared fixture uses real object allocation, barrel drop selection, placement,
force/decay, audio queues, player memory, string manager, modifier/armor lookup,
wide-string formatting and packet encoding. Guarded C-owned DeathData and InitData,
full owner update memory, source changes, created objects, audio, packets,
caches, and Logic/Other indices are captured. All replaced services and blob
regions are restored. Raw player and position addresses are normalized to stable
IDs, including generic armor's unusual position-pointer audio ID.

Discriminating cases cover inline DeathData names (empty, missing, 127 bytes),
nonzero sound and Spawn-only flag changes; first marker match including duplicates;
barrel caches/drop thresholds/counts; boulder rotation, missing first/second debris
types and 64 seeds; material priority and plural names across four languages,
with/without player holders. Real messages distinguish every selected string key.
GameBall's single-line delegation is exercised with no start object; the retained
ball mechanics are outside this conversion's behavioral scope.

Boulder allocation failure leaves the actor undeleted. Both random APIs use
Logic, not Other. The owner-chain lookup can return a terminal non-player object;
Marker preserves raw four-slot access. Armor generic passes position pointer bits
as an audio ID; deferred audio captures this behavior. The Create/Spawn record
contains inline bytes, not a name pointer.

## Native conversion and qualification

Original-C baseline `29780767` was committed and pushed before conversion.
All eight callbacks now live in src/legacy/object_death_native.go, with generated
ABI entries for registration. server__object__die__die.c is removed entirely.
The remaining static C helper only forwards a variadic message call; it contains
no translated policy or test oracle. All 1,704 locked cases match (1.335s).

Accumulated default/server/highres tests pass (60.128s/49.661s/49.157s).
Three production binaries build and verify ELF32/i386/GO386=sse2. Full-suite
failure entries match exactly (1,553; 15 pass/3 fail/32 skip packages).
Fresh `object-death-port` headless gameplay exits 0 against preserved screenshots,
overrides disabled, null audio.

Production C: **134,304 physical lines**, **152 files**, zero reference C:
**265 lines and one translation unit removed**. Artifacts: build/port-object-death.
The retained barrel loot engine's uppercase-name table branch is outside this
policy conversion; these fixtures exercise its other table. Both engine table
branches should be distinguished when that shared loot engine is ported.
