# Direct C callback addresses

Five production C functions only return the address of an existing function:

| Getter | Target | Use |
| --- | --- | --- |
| controlNormalUpdate | nox_xxx_updatePlayer_4F8100 | Restore normal player update when leaving observer mode |
| controlBotUpdateAddress | nox_xxx_updatePlayerMonsterBot_4FAB20 | Recognize bot updates before leaving observer mode |
| controlPlayerUpdateAddress | nox_xxx_updatePlayer_4F8100 | Restore normal update after failed bot-record allocation |
| controlInversionAddress | nox_xxx_inversionEffect_4E03D0 | Recognize the inversion modifier callback |
| orchestrationChestInitAddress | nox_xxx_initChest_4F0400 | Recognize chests during reward generation |

The proposed replacement is `unsafe.Pointer(C.targetFunction)` at each Go use.
Existing headers already declare the functions; all four files already import
unsafe. Keep the same C symbols, pointer types, allocation and dispatch behavior.
This removes five preamble bodies, not standalone C lines. The latter remains
six lines in one MP3 implementation include file, with zero standalone reference C.

## Baseline and acceptance

All 38 existing player-control and server-orchestration roots pass in default,
server and highres with the original getters. Existing golden hashes are unchanged.
Production evidence from the preceding qualified callback conversion is reused
only after exact whole-source fingerprint equality and all four production binary
hashes are checked. See [address-adapters-c-qualification.json](address-adapters-c-qualification.json).

The inversion fixture installs the actual modifier callback and checks behavior;
the reward fixture sets the actual chest initializer and checks resulting inventory.
Control fixtures cover ordinary observer and bot operations, but do not force
bot allocation failure or explicitly seed the observer's bot-identity early return.
They register callback identities without explicitly serializing the player's
callback slot. Do not claim those branches or assignments are frozen observations.
For this mechanical address substitution, exact source mapping and generated cgo
references to the original functions complement the existing behavioral tests.
No duplicate C getters are added as test oracles.

After conversion require the same 38 roots per normal profile, safe build/static,
fresh production/ABI and exact known-suite checks, headless creation and save/load.
The optional-safe renderer-fixture limitation documented in EMPTY_CALLBACKS.md
is unchanged; a safe build is not a claim that the broad safe fixture suite passes.

Luna drafted the five substitutions and traced callers. Primary review confirmed
the substitutions but caught overstated allocation-failure, bot-identity and actor
snapshot coverage; the helper corrected its report before acceptance. This was
useful mechanical drafting, but branch coverage required direct source review.
