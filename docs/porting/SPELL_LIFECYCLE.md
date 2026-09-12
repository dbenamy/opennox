# Spell casting and buff lifecycle

## Scope and C baseline

29 connected GAME4.c functions (4FC960 through 4FF620), 1,326 removable C lines:
phoneme delivery, book queues, mana and casting admission, spell projectiles,
shock collision, duration cancellation/ray messages, and all buff operations.
Keep the neighboring nox_setImaginaryCaster and sub_57AEE0 declarations.
Production is still C while establishing this baseline.

The optional spell fixture extends the guarded player/controls fixture, uses
actual spell definitions and phoneme trees, actual book allocation and linked
lists, and the actual duration manager. Cast completion invokes the root
PlayerSpell owner on the fixture server. The client phoneme path uses a guarded
Drawable and the real client network-code lookup; it needs no display. Complete
normalized records, objects, packets, callbacks, resources, clocks and RNG state
are captured. No C algorithm is copied into a test implementation.

The corpus currently contains 2,246 cases / 60 groups. It includes all 32 buff
slots, nil/out-of-range accessors, duration/power boundaries and exclusions;
mana subtraction, refund and aggregate admission; all player classes; positive
book insertion/progression and cancellation; nonempty duration lists; ray
messages; projectile creation and inherited enchantment; shock collision;
spell-power and position boundaries; and casting restrictions. Independent
assertions check distinct buff duration/power, mana returns, successful book and
projectile creation, and male/female/nonplayer phoneme sounds.

## Deliberate correction before baseline

The server phoneme helper added 8 to a typed object pointer, advancing 6,176
bytes instead of reading the class byte at offset 8. Nonplayer cases exposed
an out-of-bounds read and crash. Cast before adding the byte offset. The client
branch already reads the correct field. This is a documented bug correction,
not exact preservation of an invalid read. See [DECISIONS.md](DECISIONS.md).

## Conversion review notes

The existing server.Object.EnchantPower Go method reads BuffsDur; C reads
BuffsPower. Correct that owner when consolidating the native implementation and
record the behavior change. The independent duration 1,234 / power 7 contract
must remain valid. Audit and retire unused C exports before final qualification.

Local source audit and evidence: build/port-spell-lifecycle. The tracked fixture
and expected hashes make the baseline reproducible without the local captures.
Production C count before conversion: **115,985 lines / 149 files / zero reference C**.

## Locked validation

Final locked spell corpus plus all 3,205 existing controls cases passes twice
in one process (38.765s). All 60 spell and 58 controls expected hashes match.
Local final captures/log: c-final-*.json and c-final.log. Earlier diagnostic
captures are preserved losslessly as .json.gz. The only normalization correction
was the Magic update record's caster pointer in the outer callback snapshot;
all other data matched across the original repeats. Save the Magic type ID while
the fixture registry is active, since that registry is restored before this
outer snapshot. C count is unchanged. Commit/push this baseline before conversion.
