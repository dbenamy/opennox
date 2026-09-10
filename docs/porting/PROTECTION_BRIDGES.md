# Retire unused protection C bridges — 2026-09-10

The reserved-record/handle conversion removed the last C callers of the integer
record constructor. The earlier spell/ability conversion likewise left no C
callers of the one-bit helper. This cleanup removes their obsolete exports and
prototypes (`nox_xxx_protectionCreateStructForInt_56F280`, `sub_56FCB0`).

The public Go integer constructor keeps its API and calls the private constructor
directly. Its previous call to a Go function with C integer parameter types was
already a native Go call, not a C ABI round trip. The float path is unchanged.
Tests drop the retired C-only constructor mode and retain both public Go paths,
including exact signaling-NaN bits. The bit tests call the native helper; the
remaining spell/ability C entry points retain their ABI tests. Checksum and
manager entry points with actual C callers remain exported.

This cleanup removes Go bridge declarations and C header prototypes, not C
function bodies. Production C remains **142,265 physical lines** in 153 files;
chunk delta **0**, with **0** test-reference C lines. See [C_LOC.md](C_LOC.md).
Final validation results follow after checks complete.

All accumulated protection tests pass under default/server/highres on 386.
The constructor fixture retains 8,112 live-path cases (1,014 values × 4 keys ×
2 public Go paths). All three production binaries build and symbol checks confirm
the two obsolete entries are absent while needed spell/ability/checksum entries
remain. Fresh `bridges-port` warrior gameplay exits 0 against both preserved
screenshots with overrides disabled. Evidence: build/port-bridges. The full
suite was last checked at the rekey checkpoint; no new full-suite run is needed
for this declaration/path cleanup.
