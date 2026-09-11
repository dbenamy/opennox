# Health, poison, mana and gold

## Connected batch

25 functions / 646 physical C lines:

- GAME3_3 4E4560 through before 4E4670: health setter and synchronization.
- GAME3_3 4EE460 through before 4EED40: health changes/history, damage/death
  dispatch, poison activation/update/removal, and mana operations.
- GAME4 4FA590 through before 4FA700: gold arithmetic and guarded accessors.
- server__object__pickdrop__pickup.c: gold pickup and its ordinary-item fallback.

Production C remains intact. The prior trade engine is committed and pushed as
`56c79668`. Current production count is 131,120 lines / 151 files / zero
reference C. Scope data is in build/port-resources/scope.json.

## Locked original-C baseline

3,591 cases in eleven groups repeat byte-for-byte, including complete captures.
All 4,986 prior shop and trade cases pass with their existing hashes alongside
the locked resource cases (28.933s). Production C is intact. Hashes are in
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
