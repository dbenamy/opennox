# C source-size checkpoints

Update this table after every completed conversion chunk. Run
`python3 tools/porting/c_loc.py` from the repository root, or pass a Git revision
to reproduce a historical count.

Count physical lines, including blanks and comments, in tracked `src/**/*.c`.
Exclude headers, C embedded in Go preambles, external dependencies and generated
build outputs. Files beginning with `//go:build porttest` are test references and
are counted separately. This measures source size, not active code, semantic
reachability or remaining porting effort. See [the inventory](C_INVENTORY.md) for
target build selection and linker evidence.

| Completed chunk | Production .c files | Production C lines | Change | Test-reference C lines |
| --- | ---: | ---: | ---: | ---: |
| Before conversions (`0e9d2e1f`) | 153 | 142,665 | — | 0 |
| Protection checksum (`66fa7bd4`, 2026-09-10) | 153 | 142,637 | −28 | 33 |
| Retire checksum C test reference (2026-09-10) | 153 | 142,637 | 0 | 0 |
| Protection record lookup/index/swap (2026-09-10) | 153 | 142,570 | −67 | 0 |
| Protection spell/ability bitsets (2026-09-10) | 153 | 142,503 | −67 | 0 |
| Protection integer/float construction (2026-09-10) | 153 | 142,458 | −45 | 0 |
| Protection deletion/cleanup (2026-09-10) | 153 | 142,393 | −65 | 0 |
| Protection randomized insertion (2026-09-10) | 153 | 142,351 | −42 | 0 |
| Protection reserved records/handles (2026-09-10) | 153 | 142,327 | −24 | 0 |
| Protection rekey/shuffle (2026-09-10) | 153 | 142,265 | −62 | 0 |
| Retire unused protection C bridges (2026-09-10) | 153 | 142,265 | 0 | 0 |
| Protection integer/byte/word setters (2026-09-10) | 153 | 142,189 | −76 | 0 |
| Protection additive updates (2026-09-10) | 153 | 142,130 | −59 | 0 |
| Protection buffer validation (2026-09-10) | 153 | 142,115 | −15 | 0 |
| Protection object checksum/toggles (2026-09-10) | 153 | 141,984 | −131 | 0 |
| Protection float updates (2026-09-10) | 153 | 141,941 | −43 | 0 |
| Protection initialization (2026-09-10) | 153 | 141,914 | −27 | 0 |
| Protection floating RNG/state (2026-09-10) | 153 | 141,844 | −70 | 0 |
| Client unit-code/bit helpers (2026-09-10) | 153 | 141,821 | −23 | 0 |
| Dynamic unit-code/extent lookup (2026-09-10) | 153 | 141,786 | −35 | 0 |
| Waypoint allocation/link/predicates (2026-09-10) | 153 | 141,745 | −41 | 0 |
| Map-rule loading/parsing (2026-09-10) | 153 | 141,455 | −290 | 0 |
| Rule writing/online buffer fix (2026-09-10) | 153 | 141,351 | −104 | 0 |
| Rule-file deletion (2026-09-10) | 153 | 141,340 | −11 | 0 |
| Command-rule loading/dispatch (2026-09-10) | 153 | 141,215 | −125 | 0 |
| Spell-class eligibility / unused chat predicate (2026-09-10) | 153 | 141,180 | −35 | 0 |
| Player-ping minimum/average (2026-09-10) | 153 | 141,126 | −54 | 0 |
| Network alias table / exhaustion fix (2026-09-11) | 153 | 141,082 | −44 | 0 |
| Glyph/item eligibility and caches (2026-09-11) | 153 | 141,042 | −40 | 0 |
| Collision reflection / containment (2026-09-11) | 153 | 141,000 | −42 | 0 |
| Line projection / clamping (2026-09-11) | 153 | 140,903 | −97 | 0 |
| Durability classification (2026-09-11) | 153 | 140,879 | −24 | 0 |
| Waypoint link insertion (2026-09-11) | 153 | 140,845 | −34 | 0 |
| Tile selection state (2026-09-11) | 153 | 140,785 | −60 | 0 |

The checksum removes two C function definitions from GAME5_2.c; C ABI entry
points remain as generated bridges into Go. Translation-unit counts do not fall
because the file still contains other functions. The 33 test-reference lines were subsequently removed after successful
differential validation; they remain recoverable from Git at `66fa7bd4`.
