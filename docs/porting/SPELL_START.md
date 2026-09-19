# Spell creation and start callbacks

Qualified parent `1dc9e771` is pushed. The original-C baseline is qualified and frozen; production remains unchanged. Remaining C:
10,184 lines /42 files /zero reference C.

The two callbacks previously deferred from unit gameplay cover Pixie spawning
and teleport start (118 body lines). Move the charm-control owner and its Go
consumers alongside them. Pixie creation has only a private Go caller; teleport
start remains a duration callback and needs a Go-backed C export. The tracked
selection records the current C bodies.

Reuse existing sustained-spell, factory, RNG, wall/tile, audio and player-message
owners. Preserve all earlier frozen spell expectations. Teleport contracts cover
player/nonplayer caster, explicit/default recipient, all five forbidden tile
names, blocked/open rays, cooperative/ordinary delay, fractional and negative
delays and uint32 frame wrap. Pixie contracts must cover count/cache behavior,
object allocation failures, RNG consumption, positions/directions, owner/target
and expiry state, and emitted audio. Audit shipped direction constants and C
floating store boundaries before freezing.

Fresh three-target native production/ABI, exact known full-suite comparison and
headless gameplay/save-load follow conversion. Original-C production evidence may
reuse the parent only after verifying all baseline source changes are porttest-only.

## Initial C evidence

The first teleport-only run passes1,728 cases (an earlier conversation estimate
of3,456 was incorrect). The expanded second run passes2,520 cases:1,728 teleport,
24 explicit rounding boundaries and768 Pixie combinations. The teleport fixtures
now install an owned localized string and assert exact recipient/message bytes
and deferred audio; the first exploratory hash is superseded, not frozen.
Pixie contracts predict RNG consumption and created-object position/direction,
expiry, owner record, cache and audio using actual shipped direction constants,
the existing object factory and map traces. The next run adds64 nil/player-owner
cases. These are still initial results; freezing and cross-target gates follow.

Compiled C audit (`build/port-spell-start/pixie-c-disassembly.txt`) confirms the
radius is stored as float32, direction multiplication/addition stay wide, then
coordinates are stored as float32. The native draft is ignored and NOT installed.
The fixture records the real creation callback's owner argument; the established
creation fixture does not run the entire production world insertion sequence.
Inherited ownership/target-selection contracts and fresh integration supplement
this boundary, as in preceding spell batches.

Selection audit found that neighboring test names use `TestTemporary...` and
`TestProjectileCollision...`, not the plural prefixes in the first primary
selector. Add the31 neighboring roots as a separate recorded sweep on each
target; keep the97 completed primary roots and their source fingerprints. Native
qualification must include both selections. This corrects coverage without
rerunning the already qualified primary cases.

## Qualified C baseline

All2,616 focused cases pass in separate final/repeat processes, including the24
near-integer Pixie limits and8 actual charm setter checks. Five capture hashes
are frozen. Each target passes97 primary roots /357 captures plus31 neighboring
roots /65 captures, with no skips. The two root selections are disjoint, and all
422 capture hashes are identical across default/server/highres. Source
fingerprints match all six sweeps and the reviewed checkout; static-current passes.
Primary drivers enforced five new hashes and inherited literals; neighbor drivers
enforced inherited literals. Full inventories were frozen after comparison.

All four source changes are porttest-only. The three `1dc9e771` production
binaries rehash correctly, and both selected C functions plus the charm global
remain present as C interfaces. That parent's production/ABI, exact known
full-suite and gameplay/save-load evidence is explicitly reused; no fresh C
production/scenario run is claimed. See `spell-start-c-qualification.json`.
All test/build jobs are joined. Next: integrate the ignored native draft and run
both affected selections and fresh production qualification before committing.
