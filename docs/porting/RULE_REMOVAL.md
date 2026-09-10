# Rule-file deletion — 2026-09-10

Scope: 57A9F0 constructs `maps\\<map>\\<file>` and removes that path through the
existing case-insensitive filesystem implementation. Its live C callers include
server-options reset and temporary rule-file cleanup. Keep 57A950 separate: it
passes a derived path to the broader 4D0550 rule/config loader.

The original-C fixture creates an isolated temporary tree, saves/restores the
working directory and snapshots every file and directory. Ten cases cover exact
0/1 returns, missing and existing files, case-insensitive/backslash paths, nested
file paths, sibling preservation, empty/nonempty directory removal, empty map
or filename, NUL termination of both arguments and a 238-byte filename within
the original buffer limit. Successful file removal leaves its directory intact;
failed removal must leave the entire tree intact. Expected trees are specified
independently rather than computed by the production path/removal helper.

Production C before conversion: **141,351 physical lines**, 153 files, zero
reference C lines. Local artifacts: `build/port-rule-remove/`.

All ten original-C cases pass on 386 before replacement. No production changes
are part of the baseline checkpoint.

## Go conversion

Original-C baseline: `b861ab46`. The live 57A9F0 bridge now performs native Go
path concatenation and calls the same ifs.Remove implementation. C strings still
terminate at NUL; paths retain the original backslashes and case-insensitive
lookup behavior. There is no temporary fixed-size C path buffer or test-only
C implementation. 57A950 and the command-rule loader remain outside this chunk.

Production C: **141,340 physical lines (−11)** in 153 files; test-reference C: **0**.
All accumulated protection/network/waypoint/rule tests pass on 386 default,
server and highres. Symbol inspection confirms the Go-backed live C entry.
The latest full-suite milestone remains the writer conversion, whose 1,553
failure entries exactly matched the known baseline; this isolated deletion
helper does not require another full-suite run immediately afterward.

All three production binaries build. `rule-remove-port` passes the preserved
headless gameplay scenario with both screenshot checks and overrides disabled.
