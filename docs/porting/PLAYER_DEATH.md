# Player death, scoring and corpse creation

## Scope

Seven C bodies are now Go: PlayerDie, arena/elimination/King of the Realm scoring,
death notification, corpse cache initialization and corpse spawning. The registered
PlayerDie callback retains its C signature; six private C interfaces are removed.
The live respawn caller invokes Go directly. No C algorithms remain solely for tests.

The corrected C baseline contains 679 body lines. The conversion removes **693
physical C lines**, including adjacent obsolete headings and blanks, leaving
**26,145 lines in 67 production C files**, with zero reference C.

## Behavior and contracts

Fourteen focused capture groups cover 1,101 cases through real player/team,
object-factory, inventory, RNG, audio, ability, message-queue and statistics owners.
The lifecycle tests invoke the actual registered PlayerDie callback.

- Arena and elimination: assigned/unassigned/friendly teams, environmental and
  self deaths, assists, signed score and unsigned death-count wrap, exact message
  bytes and ordering, and match-statistics actor/target attribution.
- King of the Realm: owned crowns, fractional point conversion, team scoring,
  friendly/self penalties, held/unheld crown transfer and transfer disablement.
- Recent assists: strict ten-second expiry, frame wrap, missing/inactive players,
  unavailable network codes, and exclusion of the victim and killer.
- Death lifecycle: online source/weapon reporting, owned projectiles and monsters,
  state/mana/casting cleanup, active abilities/cooldowns/enchantments, male/female
  and special-cause sounds, and cooperative pending-character-load cancellation.
- Quest: remaining lives, exhaustion, fractional starting lives, gold/statistics
  penalties, reset messages, per-player slots and exact RNG consumption.
- Corpses: shipped suffix/direction/position tables, missing object definitions,
  untouched center slot, early allocation failure, position rounding, flags and
  exactly one decay-duration draw per created object.

Preserve the historical scoring asymmetries. Arena self/environment deaths without
an assist subtract a lesson without incrementing the death count. King of the Realm
with no player killer leaves the victim's death count alone. Unteamed enemy kills
can assign a dropped crown to the killer; the teamed scoring branch does not.
The crown-drop adapter's second argument is a pending-owner pointer despite its
old parameter name `stamp`; captures normalize only that identified pointer.

Preserve the low-byte-only reset of TrapSpellsCnt, game-mode precedence and
queue ordering. The callback's historical integer return has no reader in the
actual caller/dispatch graph; the Go export returns zero with the same C signature.

## Corrections and review notes

Independent arena contracts exposed an original-C null-team dereference for an
unteamed killer and teamed victim. The baseline adds a two-line guard around that
team-score update. The player still gains the lesson. This is a small, reversible
prerequisite correction, documented for later review in DECISIONS.md. Original
failure evidence: build/port-player-death/original-teams/tests.jsonl.

Fixture corrections preserved a valid player roster by allocating separate monster
and projectile sources, explicitly initialized zeroed allocations, restored temporary
class changes before typed snapshots, and initialized the real ability manager.
They did not change production behavior or frozen expectations.

Review added 48 death-driven statistics contracts before conversion, because the
initial scoring fixtures had disabled recording. The first native run then caught
six valid-assist failures: Players.ByID uses network codes, whereas the C adapter
uses Players.ByInd for player slots. The native implementation was corrected;
no golden was regenerated to hide the difference.

## Qualification

C baseline **4c3191c1** is committed and pushed. Independent focused captures repeat
identically in separate processes on all three targets. The full affected selection
covers death/controls, inventory/respawn, object state, roster/team/gameplay reports,
quest/spell lifecycle, statistics, objectives/rewards, session and orchestration.

Both C and native pass **318 roots /24,386 entries per target**, no skips. All
**149 captures /33,295 records** match across targets and between implementations.
Each implementation's target and production gates use identical source fingerprints.
Static checks, three fresh 386/SSE2/CGO production binaries and ABI checks, exact
known full-suite failures, headless gameplay and explicit save/load pass.

See [C qualification](player-death-c-qualification.json),
[native qualification](player-death-native-qualification.json),
[capture index](player-death-captures.json) and
[selection/caller audit](player-death-selection.json). Raw evidence is under
build/port-player-death/c-final-* and native-*.

## Recovery and disk

All installer/freezing scripts and copied drafts are consumed. Actual source and
committed expectations take precedence; do not replay ignored mutation scripts.
Original assets and the archive remain unchanged and outside Git.

Verified cleanup removed twelve superseded binaries (584,677,176 bytes) and
72 older compressed binaries (1,648,354,499 bytes), preserving successful build/ABI
reports and manifests. Completed C scenario copies were verified and deduplicated,
reclaiming 1,112,747,701 bytes. Restore manifests preserve file hashes, modes and
timestamps; restore mode remains in deduplicate-player-death-c-assets.py. Audit
and apply modes are consumed. Latest production binaries and raw evidence remain.

Completed native scenario copies were also verified/deduplicated, reclaiming
1,112,747,701 bytes. Their restore mode remains in
build/port-player-death/deduplicate-player-death-native-assets.py; audit/apply are
consumed. All build/test/cleanup sessions are joined.
