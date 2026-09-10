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

The checksum removes two C function definitions from GAME5_2.c; C ABI entry
points remain as generated bridges into Go. Translation-unit counts do not fall
because the file still contains other functions. The 33 reference lines include
build constraints/provenance comments and are absent from ordinary binaries.
