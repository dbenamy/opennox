# Health, poison, mana and gold

## Connected batch

25 functions / 646 physical C lines:

- GAME3_3 4E4560 through before 4E4670: health setter and synchronization.
- GAME3_3 4EE460 through before 4EED40: health changes/history, damage/death
  dispatch, poison activation/update/removal, and mana operations.
- GAME4 4FA590 through before 4FA700: gold arithmetic and guarded accessors.
- server__object__pickdrop__pickup.c: gold pickup and its ordinary-item fallback.

The completed batch removes 646 physical C lines and one C file. Production
count is 130,474 lines / 150 files / zero reference C. Scope data is reproducible
from original-C baseline `0393b17c` and the ranges above.

## Locked original-C baseline

3,591 cases in eleven groups repeat byte-for-byte, including complete captures.
All 4,986 prior shop and trade cases pass with their existing hashes alongside
the locked resource cases (28.933s). Production C was intact for baseline capture. Hashes are in
src/resources_porttest_test.go; captures are reproducible using porttest and
OPENNOX_CALLBACK_CAPTURE. Ignored artifacts: build/port-resources.

The fixture extends the shop/player/callback harness with guarded full health
and monster-update records. Each action captures object and complete player
bytes, health/update data, queued packets, protection records, RNG and counters,
death callback, harpoon and pickup requests. Temporary state is restored.

Coverage: gold arithmetic and pickup; mana signed/unsigned boundaries, wrapping,
limits, old values and god mode; health setters, history, restoration, disabled
healing, inventory checksums and owned-monster synchronization; repeated lethal
and nonlethal damage, death overrides and deletion; poison mutation, reporting,
resistance, immunity and random seeds. Missing records and nil inputs are tested
where the C contract supports them. Full protection algorithms execute against
seven guarded real records. Gold pickup uses actual allocation/deletion,
localized messages and audio requests.

Harpoon breaking and ordinary pickup are retained services recorded at their
existing Go hooks, including arguments/results. Damage runs the actual retained
monster death wrapper, but does not cover every monster-specific death callback
or XP reward chain; those algorithms remain outside this conversion. Poison
resistance uses the retained calculation and real balance tables/RNG.

Nonplayer mana-add returns are first checked against the low sixteen bits of
the actual input pointer, then normalized to stable identity. Monster definition
and sound-set pointers are normalized by identity, as are existing fixture
pointers. No algorithm is copied into a test oracle. No hashes were obtained
from Go replacements.

Caller audit includes C/Go uses, headers and GoldPickup registration. Keep live
ABI roots until retained C callers are ported; route native Go callers directly.

## Native implementation

Original-C baseline: `0393b17c`. The 25 functions move to legacy/resources.go,
with 23 C ABI roots in resources_exports.go for retained callers and the pickup
registration. Private owner-health reporting and object-gold adjustment bridges
are retired; native callers use the Go implementations directly. GoldPickup has
both native registration and the retained callback address. The sole embedded C
adapter calls the existing variadic line-message service with its integer value.

The complete eleven native capture files are byte-identical to original C;
3,591 cases pass in 9.392s with unchanged hashes.
The conversion removes 646 physical C lines and the entire pickup C
file: 130,474 production lines / 150 files / zero test-reference C.

Preserved details include signed health getters and signed/unsigned overflow,
health-holder checksum toggles before and after mutation, all 32 monster sync
words, single death dispatch across repeated lethal calls, poison byte wrapping
with status/timer decisions based on the full integer input, resistance's float
spill, and owner/cooperative poison packet routing. Mana subtraction compares
the updated mana for its protection delta; refresh adds the signed maximum to
protection rather than setting it. Gold subtraction saturates the balance while
applying the full negative protection delta. These historical behaviors are
covered by the locked C captures rather than silently corrected during porting.


The accumulated tests also exposed a prior createWeapon write-barrier bug on
uninitialized C modifier slots. That initializer now stores the C descriptor's
address bits without scanning the previous slot as a Go pointer. See
OBJECT_CREATION.md; the existing original-C corpus remains unchanged.

## Qualification

Accumulated default/server/highres port tests pass in 114.993s / 87.221s / 88.604s.
All three production builds are ELF32/Intel 80386 with GO386=sse2. The full suite
matches exactly the original 1,553 failure entries (15 pass / 3 fail / 32 skip
packages). Fresh resources-port gameplay exits zero against unchanged repeat-a
goldens, with overrides disabled, Xvfb and null audio. The unchanged creation
corpus passes three repetitions with GOGC=20 (25.862s). No baseline hashes changed.
