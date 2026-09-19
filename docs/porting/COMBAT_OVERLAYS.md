# Client combat overlays — C baseline in progress

Qualified parent: **f588409c**, native speech bubbles. Production C remains
**25,368 lines /67 files /zero reference**. No conversion in this batch yet.

Scope: 32 connected client presentation bodies /871 original C body lines:
ally slots, friend membership, kill-feed ring/drawing and console notifications,
drawable-attached effects, and floating health changes. The tracked selection
records live decoder, renderer and initialization/cleanup callers. The nearby
ordered-message queue, browser and audio streams are separate subsystems.

## Independent contracts

Nineteen focused groups /707 records pass in all three c-repeat targets.
Captures and source fingerprints match exactly, including the positive icon-pixel
contract. Goldens are frozen; wider qualification remains. Coverage includes:

- Ally capacity 32, duplicate updates, cleared-field padding, full rejection,
  removal/reuse, missing output preservation and byte/word boundaries.
- Friend capacity 128, prepend/duplicate order, first-match removal and lifecycle.
- Health capacity 32, signed amounts including −32768, input sequence timestamps,
  removal/lifecycle, expiry at 30/31 frames, local-player colors and actual pixels.
- Feed literal player-name glyphs; 27 participant-presence combinations through
  actual insertion/console formatting; 205 insertions across ring wrap; 75 expiry
  cases including uint32 frame wrap and the actual four-row limit at 480px.
- Initialization caches, projectile-to-weapon substitutions, real spell/ability/
  thing icons and fallback selection, checked through actual image pixels.
- Global and per-drawable effect chains, head/middle/tail/sole removal, real Go
  membership/cleanup callers, and 150 trail-history frames through sprite/line
  rendering with position, count, reserved-byte and alpha postconditions.
- A supplementary speech-bubble contract distinguishes equal-width text through
  glyph traces and pixels without changing any earlier frozen expectation.

Owners reuse the real renderer, fonts, players/teams, drawable and allocation
classes. Every mapped/global value is restored. Original mapped console-format
and direction tables are supplied explicitly. Captures normalize owned pointers.

## Prerequisite repairs for review

1. Three C feed-name paths treated player names as format strings, collapsing
   percent pairs. Independent glyph checks reproduced it; use literal UTF-16
   copies. The assist prefix retains its fixed format string.
2. A missing victim reused the preceding console notification's name in 18 cases.
   Initialize that C buffer to empty before formatting.
3. Go DrawableFX.Next occupied byte 16, which is movement history; C stores the
   actual per-drawable link at 64. Real-allocation contracts reproduced failed
   membership and incomplete cleanup. Move Next to 64, assert its offset, and
   remove the obsolete small-pointer workaround. Both contracts now pass.

All are reversible pre-baseline fixes; none changes production C LOC. Fresh
production qualification is required because production source changed.

## Fixture diagnostics and evidence

Earlier health pixels used transparent colors; corrected with the color
constructor. The first console fixture omitted the original format string,
producing empty output before it exposed the separate missing-victim bug.

The sprite-trail fixture passed nonzero initialization values to alloc.New,
which uses its argument only to infer the allocation type. Its image handle
therefore remained zero. Explicitly initialize allocated draw data and retain a
real direct-sprite prerequisite assertion. Earlier lighting/opacity adjustments
alone did not fix this; the shared owner already supplies fixed lighting.

All historical diagnostic runs remain under build/port-combat-overlays. Frozen captures are indexed in combat-overlays-captures.json. All repeat
sessions are joined.
Every copied fixture draft is consumed; actual source wins. Wider affected
selection includes 176 roots before target constraints. Static checks, broader affected tests and fresh production qualification
remain before translation.

## Disk maintenance

Removed 131 older regenerable Go compiler artifacts with all builds joined,
reclaiming 6,444,559,234 bytes. removed-stale-go-cache.json records the paths and
reason. Go regenerates these on demand. Original assets/archive, source,
production binaries and captures remain. Removal is consumed.
