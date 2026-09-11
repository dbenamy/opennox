# Modifier effects and weapon use (planned)

Connected next batch: **41 functions / 977 physical C lines**. Production is
still original C for this family. Equipment is the preceding batch; starting
count after equipment is 128,081 C lines / 149 files / zero reference C.

## Scope

- src/legacy/GAME3_2.c: 004DFB50, 004DFB80, 004DFBB0, 004DFC30, 004DFCA0, 004DFD10, 004DFD40, 004DFD80, 004DFDB0, 004DFDE0, 004DFE10, 004DFE40, 004DFF40, 004E0040, 004E0140, 004E0170, 004E01D0, 004E02C0, 004E0370, 004E0380, 004E03D0, 004E03F0, 004E0480, 004E04C0, 004E04D0, 004E0640, 004E0670, 004E06F0, 004E0740, 004E07C0, 004E0850, 004E08E0, 004E0960, 004E09B0.
- src/legacy/GAME4_3.c: 0053C520, 0053C940, 0053F290, 0053F480, 0053F4F0, 0053F670, 0053F8E0.

The first family covers item engagement flags, inventory modifier lookup,
speed/protection, regeneration/replenishment, armor/damage modifiers, status
and resource effects, readiness and projectile speed. The second covers item
recharge/rate, fireball/wand projectiles, wand spell acceptance, fire-wand sparks
and use callback dispatch. Keep damage resolution and spell implementations as
production dependencies. The unused grip search has a keep.go reference under
the disabled `none` tag; preserve its behavior within the family for now.

## Testing plan

Extend the existing equipment/inventory/resource/shop fixture with an optional
EffectsUse input and snapshot; keep all existing 12,009 contracts unchanged.
Use action IDs starting at 500, a thin original-C dispatcher, and stable function
pointer identities. Capture actual modifier descriptors, guarded object/update/
use/health memory, buff state, spell acceptance arguments and outcomes, created
projectile/Spark state, packets, audio, RNG state and raw return bits.

Cover all modifier slots and inventory filters; null/class guards; speed and
protection rounding/clamp boundaries; frame-wrap timing; regeneration HP and
holder eligibility; replenishment byte wrap; signed effect arithmetic; status
application and notification order; mana/HP limits; readiness callback identity;
recharge saturation; wand empty-charge/cooldown checks; accept/reject outcomes;
player targeting versus nonplayer coordinates; single/triple fireball directions;
allocation failure; projectile velocity and spark RNG; use callback return values.
Keep original valid-input preconditions explicit (including nonzero divisors).

Repeat and compare full original-C captures, lock hashes, verify existing
contracts, then commit/push the baseline before conversion. Convert the whole
family and qualify once at its boundary: accumulated default/server/highres
corpora, three production builds, unchanged full-suite failure multiset, fresh
headless gameplay. Record actual C_LOC, commit/push, summarize and continue.

Detailed local audit: build/port-equipment/next-effects-use-scope.json and
build/port-effects-use/fixture-plan.md. No user question is pending.
