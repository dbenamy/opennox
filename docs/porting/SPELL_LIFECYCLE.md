# Spell casting and buff lifecycle

## Scope and C baseline

29 connected GAME4.c functions (4FC960 through 4FF620), 1,326 removable C lines:
phoneme delivery, book queues, mana and casting admission, spell projectiles,
shock collision, duration cancellation/ray messages, and all buff operations.
Keep the neighboring nox_setImaginaryCaster and sub_57AEE0 declarations.
Production remained C while establishing this baseline.

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

## Native implementation and first comparison

Baseline c6435866 was pushed before conversion. All 29 bodies are replaced by
spell_lifecycle_{buffs,mana,casting,books,duration,exports}.go. Go callers use the
native owners directly; 15 unused C exports and header declarations are retired
before qualification, with 14 production ABIs retained. The fixture dispatches
native helpers and retains exactly the locked expected hashes.

All 2,246 cases / 60 groups match in 9.329s. The first native run matched 59 groups;
projectile Y positions differed by one ULP for two directions. The compiled C
keeps the intermediate at x87 53-bit precision through inherited velocity
addition. Keeping the Go intermediate float64 until the final store reproduces
all 56 projectile cases exactly. No expected hashes changed for conversion.

The existing Go EnchantPower accessor now reads BuffsPower; see the separately
recorded decision. The book record remains 60 bytes in the real allocator, and
its index/queue overlap and raw next-spell read preserve the original layout.
No reference C algorithms remain. Current production C: **114,659 lines / 149
files / zero reference C**, a reduction of 1,326. Final qualification is recorded below.

Final C captures are now losslessly compressed c-final-*.json.gz. Native final
captures are native-final-*.json; local comparison logs are preserved.

## Final qualification

Accumulated port checks, including all 46,805 focused cases, pass in default /
server / highres: 190.330s / 193.123s / 196.334s.
Three builds verified ELF32/i386/SSE2/CGO; all 15 retired symbols are absent.
Asset-backed full-suite failure multiset unchanged (1,553 entries; 15 pass, 3 fail,
32 skip). Fresh unchanged repeat-a headless gameplay passes in 38.474s.
Evidence: build/port-spell-lifecycle and baseline/runs/spell-lifecycle-port.
