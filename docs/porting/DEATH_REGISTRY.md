# Object death callback registry

## Scope

Proposed typed dispatch covers14 resource registrations and four existing callers:
projectile collisions, temporary projectile lifetime, chest collision and resource
health depletion. Keep all names, C addresses, optional132-byte data sizes, parser
entries and raw object layout. Existing registration API remains; unregistered
callbacks retain the void-pointer C fallback. Each caller retains its original
nil/else behavior. The new shared method would no-op on a nil slot.

Most registrations invoke existing Go helpers. Potion/ImpEgg wrappers extract tiny
shared helpers with unchanged audio, flags and deletion effects; ImpEgg retains its
exported flags return. Player, Spawn and GameBall also retain exported return values
that these void callers discard. Glyph uses a closure reading its mutable Go
function variable at invocation time. Receiver lifetime stays explicit through the
raw fallback. No address or export retirement, ownership or algorithm redesign.

## Original baseline

An additive porttest DeathRegistry option installs the real C callback slot in each
specialized object-death/creation/generator owner. It validates the pointer and size
against14 original C symbols, records stable callback identity for snapshots, and
enters the existing projectileDie caller. The same entrypoint will reach the new
registry after conversion; it contains no copied raw invocation that could keep
testing the old path accidentally. It sets the original caller's dead flag before
dispatch. Existing direct-export routes and all their goldens remain unchanged.

Nine frozen captures contain1864 cases:1860 across eight object/creation/generator
captures and four player lifecycle cases. Coverage includes spawn names/sizes,
barrel chance/cache/RNG, marker owner slots, boulder RNG, equipment descriptions/
materials/ownership, potion/imp/cloud behavior and generator quest/killer state.
Eight independent raw pointer/eligibility cases cover exact receiver forwarding,
once-only callback execution and rejected targets. A separate GlyphDie contract
changes its handler between two calls and checks the current handler receives the
same object exactly once. Player cases independently require death reports, audio,
flags/state, cleared mana/casting and preserved padding. Marker and boulder independent
assertions are retained too. All14 registrations execute through the actual owner.

GameBall uses the existing no-start branch; no claim of positive ball-reset coverage.
Glyph checks handler dispatch identity, not trap gameplay. The affected117-root
selection also retains damage/attack, creation/death, generator, lifetime, resource,
projectile, chest and spell consumers. Existing projectile/lifetime test callbacks
alone were insufficient to prove registered dispatch; they remain useful for caller
branches and raw fallback alongside the new real-registration fixtures.

Original86830 failed before freezing: primary mistakenly used alloc.New's argument
as an initializer, but this allocator accepts only a type exemplar and zeroes memory.
The target class remainedzero, rejecting dispatch. Marker/player/glyph independent
assertions exposed the missing execution. Explicit target class assignment fixed the
fixture; original37503 then passed all12 new roots. Failed logs/captures are retained.
No production correction or changed old golden was needed.

## Delegation and cleanup review

Luna drafted production, callback mappings and player/glyph sibling tests; primary
owns actual-owner routing, other matrices, raw forwarding and qualification. Review
corrected two unused imports and a named object.Flags-to-uint32 return conversion
before applying the production draft. The player sibling matches the existing test
apart from its name/comment, dispatch route and distinct capture label/hash.

A cache audit incorrectly included15files later than its advertised16:00 cutoff.
Primary validation rejected it before deletion. The original audit plus erratum
remain, and only8independently verified eligible obsolete archives were removed:
398794162bytes. Source/current four binary marker absence, archive/hash/stat and
host FD/executable/mapping checks passed. Original assets/current binaries preserved.
Audit timestamp criteria require independent verification; model prose is not proof.

## Qualification and count

All117 baseline roots pass in default/server/highres, no skips. All nine frozen
captures match in fresh processes; raw/glyph contracts pass. Seven porttest-only
files differ from damage conversion5e8dd778; all other source/dependency fingerprints
and four production binary hashes match. All14 C exports exist. Baseline44588 joined
PASS. See [baseline qualification](death-registry-c-qualification.json).

An extra unchanged TestProjectileCollisionDeath capture (10598 PASS) confirms three
eligible nil-slot cases emit DelayedDelete(actor70000); the other21 do not. Its
original24-case hash remains unchanged. Thus the explicit ineligible-nil raw test
does not stand in for eligible nil-slot coverage. No duplicate fixture was needed.

The cache cutoff was subsequently explicitly broadened to18:00 by primary review;
15additional obsolete archives648228414bytes then passed fresh source/four-binary,
stat/hash/header and host-use checks and were removed. This is separately recorded,
not acceptance of the erroneous earlier pre16 claim.

Production patch remains unapplied at this checkpoint. After conversion require
fresh contracts, safe/static, four production/ABI binaries, exact known-suite
comparison, headless creation and save/load.

Standalone C remains zero files/physical lines; production C preamble bodies remain79.
This removes registered round trips while preserving external callback compatibility.
