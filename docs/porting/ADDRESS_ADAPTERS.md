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

## Qualified direct addresses

The five getter bodies are removed. Every use now takes the address of the same
C function, preserving its unsafe.Pointer type and surrounding control flow.
Generated cgo static imports/link names and all five translated use sites confirm
that mapping. The first ad hoc inspection check assumed an older generated
variable spelling; it was corrected after reading the Go 1.26 output. No code
change was needed for that audit correction.

All 38 selected roots pass in each normal profile with unchanged golden hashes.
Safe build/static, three fresh production builds/ABI, exact known-suite comparison,
headless character creation and explicit save/reload/resumption pass. The suite
still has the established 1,553 failure entries; this does not claim a green suite.
All four binaries retain ten distinct Go-backed empty callback addresses and the
four target function symbols; the five old getter names and PortTest symbols are
absent. Only the four intended source files changed from the qualified baseline.
See [address-adapters-native-qualification.json](address-adapters-native-qualification.json).

Standalone C remains **six lines/one file**, zero standalone reference C. The
same heuristic scanner now finds 81 production preamble bodies, down five: 76
shared dispatchers and five remaining adapters. This excludes generated cgo and
is not an active preprocessed function count. Existing C interop remains.

To maintain disk headroom, ten superseded entry-classifier executables were removed
after hash/qualification checks and checking 159 host processes for references.
This reclaimed 489,077,624 bytes. All source/manifest/log/capture/qualification
records and current binaries remain. Historical entry and dependent safe-baseline
finalizers require rebuilding those old executables before rerunning. The consumed
record is build/port-artifact-cleanup/superseded-entry-removed.json. Successful new
scenario asset copies were separately deduplicated with restoration manifests.
