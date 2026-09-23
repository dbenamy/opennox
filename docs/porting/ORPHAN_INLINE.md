# Orphan inline and empty callback cleanup

**Qualified** against pushed storage baseline `f5970121`. See [native qualification](orphan-inline-native-qualification.json).

The patch removes four unused header helpers (`memset32`, `COERCE_FLOAT`, `noxCopyRect`, and `noxSetRect`) from `src/legacy/defs.h`; three unused GUI C wrappers (`nox_window_call_draw_func_go`, `nox_window_call_func_go`, and `nox_window_call_tooltip_func`) from the cgo preamble in `src/legacy/gui_window.go`; and the unused `iswalpha_go` wrapper from `src/legacy/input_c.go`. In `src/legacy/GAME5_2.c` and `src/legacy/GAME5_2.h`, it removes the definitions and declarations for `nullsub_9`, `nullsub_10`, `nullsub_24`, `nullsub_30`, and `nullsub_31`, along with the definition and declaration of `nullsub_35`. It also removes only the two `C.nullsub_35` calls in `objectiveObelisk` in `src/legacy/objectives_update.go`; both surrounding `NeedSync()` calls and the branch logic remain.

The focused consumer selection is the four roots in [orphan-inline-tests.txt](orphan-inline-tests.txt): `TestObjectivesObeliskIdle`, `TestObjectivesObeliskTransfer`, `TestObjectivesObeliskWands`, and `TestObjectivesObeliskEligibility`. The idle and transfer tests dispatch through op808 into `objectiveObelisk`, exercising the idle energy-boundary and spent-energy `%8` paths. Wand and eligibility roots cover related variants.

All ten registered C callback identities remain separate and retained: `nullsub_22`, `nullsub_29`, `nullsub_36`, `nullsub_38`, `nullsub_39`, `nullsub_40`, `nullsub_41`, `nullsub_42`, `nullsub_43`, and `nullsub_44`. The symbol gate checks their presence and distinct addresses in each production and safe binary.

The tracked `.c` inventory falls from 51 to **45 physical lines across four production C files** (−6). That count covers `.c` files only; headers and C embedded in Go preambles are excluded. All three legacy contracts and all four Obelisk roots pass without skips in default/server/highres, with unchanged frozen storage hashes. Static checks, safe build, all three production binaries/ABI, exact known-suite comparison, headless gameplay and explicit save/load pass. The preflight binary matches production and all source fingerprints agree. The ten retained callback addresses remain distinct in the preceding storage binaries, current production binaries and safe binary. Safe runtime was not exercised; known-suite failures are unchanged.

Metadata note: the preflight command inherited the preceding storage batch’s
`OPENNOX_DISPLAY_IMPLEMENTATION` descriptive label. The scenario is correctly
`orphan-inline-native`; accepted source fingerprints and binary hashes identify
the cleanup code. The original archived manifest is retained unchanged.

The three-file Go AST audit permits exactly the two removed calls in `objectiveObelisk`; arguments are side-effect-free address/float-bit conversions. Both `NeedSync()` calls and their branches remain unchanged. The audit uses a fresh FileSet when printing to ignore source line positions, which otherwise introduced a false formatting difference. No production adjustment was needed after application.
