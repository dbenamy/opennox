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
candidate-source.txt. Baseline implementation has not started.
