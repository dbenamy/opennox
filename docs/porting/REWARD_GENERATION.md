# Object initialization and reward generation

Next connected batch: 21 address blocks / 1,862 physical C lines in
GAME3_3.c, from 004F0390 through 004F2210. Includes simple object initializers,
generator setup, weighted reward categories and tiers, spell/ability/field-guide
books, armor/weapon modifiers, potions/gold/gems and marker selection/placement.
All candidate production functions remain original C.

Use the existing guarded object/owner fixture and real object allocator, type
registry and RNG. Add guarded reward-marker data, definition and modifier-table
inputs, and enough initialized created-object storage for books and modifiers.
Record full created objects, modifier assignments, placement/deletion requests,
input guards, returned identities, globals and RNG position. Use actual shipped
constant tables where relevant; make definition availability and class masks
explicit fixture inputs. Avoid empty-table-only coverage: require successful
creation for every reward family and each modifier slot, plus absent/filtered
entries, explicit book choices and marker probabilities.

Test tier boundaries and weighted-selection endpoints across seeded runs,
including the original odd-tier table overwrite. Check class restrictions,
duplicate-modifier suppression, masks, creation failures, quest stages, player
counts and marker flags. Separate invalid caller inputs from supported edge
cases. Repeat complete original-C captures in separate processes, lock hashes,
commit/push before any production conversion. Then port the connected batch,
run accumulated variants/builds/full-suite/headless qualification once, update
C_LOC and recovery docs, commit/push and continue.

Local audit: build/port-reward-generation/candidate-scope.json and
candidate-source.txt. Guarded fixtures and the initial corpus are implemented; original-C runs are
in progress. Initialization (288 cases) and tier selection (896 cases) passed.
The first factory run exposed missing fixture ownership of unplaced returned
objects; the dispatcher now registers their actual returned allocation for
capture/teardown without synthesizing placement or mutating object state.
No production reward conversion has started.

## Original-C baseline

All 21 production functions remain C. **6,926 cases / 48 complete capture groups**
repeat byte-for-byte in separate processes and their hashes are locked in
`src/reward_generation_porttest_test.go`. The combined regression passes
**41,306 cases / 363 groups** in 122.877s. Baseline-only locked rerun passes (20.166s).
Production C remains **119,818 lines / 149 files / zero reference C**.

The corpus exercises initializers, tier and category boundaries, explicit book
lists (including bytes other than 1), missing/filtered definitions, every reward
factory, modifier eligibility/fallback and duplicate IDs, wand replenishment,
gold/XP scaling and marker placement. Captures include 288 successful armor
creations, 288 weapon creations, all four modifier slots, 144 Ankh placements
and 16 cases requesting potion deletion. A permanent positive test checks an
actual four-modifier weapon and an explicitly selected ability-book payload.

Fixture storage uses guarded 256-byte marker inputs, guarded created init/use
buffers and the shared 64-byte update payload/16-byte guard. Factory-returned
objects are registered for capture and teardown without inventing placement
calls. Definition/type tables and gold constants are restored after each case.
No C algorithms were copied into the fixture. Sparse modifier masks and equal
identifier bytes on different descriptor pointers distinguish mask filtering
and duplicate suppression from superficially similar implementations.

Local evidence: build/port-reward-generation/c-final-*.json, c-confirm-*.json,
baseline-hashes.json, positive-coverage.json and combined-c-regression.log.
Commit and push this baseline before production conversion. Expected batch C
reduction: 1,862 lines, leaving 117,956 (subject to the final physical count).
