# Native record storage and callers

Status: qualified on Linux 386/SSE2. **11 production C imports and six test
C imports removed** across 40 source files. Selected production cgo files fall
247→236; C exports remain 1,179. One unused specialized C body retires, leaving
78 embedded callback bodies. Standalone production/test C remains zero.

See the [baseline](go-native-record-storage-baseline.json) and
[batch manifest](go-native-record-storage-batch.json). All four preceding phases
match the current 3,052 source fingerprints exactly, so their frozen captures,
root inventory and production behavior form the accepted original baseline.

## Scope and compatibility

The 40-file connected batch replaces C record types with existing native owners
in shared storage, browser UI, inventory, renderer, particles, map lists and
memory-file callers. The grid retains 128 rows of 128 44-byte cells, each modeled
as eleven raw words. Signed fallback values and list addresses are converted at
lookup. Allocation/free ownership and the retained outer pointer table are unchanged.
Timestamps retain eight uint16 fields and the existing clock/E2E behavior.

Independent Linux 386 probes compare actual C headers and candidate native
records: grid44, timestamp16, inventory148 (84 cells), list12, window404,
animation68, renderer1056, particle52, map-item36 and pointer4 bytes, including
relevant alignments and offsets. The existing native memory-file compile-time
assertion separately enforces its 16-byte layout. Storage remains unmanaged;
raw word aliases and all actual C export signatures remain unchanged.

The browser record's existing Go172/packed-C169 size discrepancy is preserved;
this batch does not change that allocation or packing behavior. libc character
classification and floating-point helpers remain outside this batch.

The unused `go_nox_drawable_call_sprite_func` has no runtime references; its
removal leaves the live drawable C alias and export entrypoints in place.

## Review and qualification

One Luna helper supplied six bounded files and reviewed the primary draft.
Primary independently reconstructed every draft and reviewed layouts, ownership,
callers and preamble closures. Review corrected four fixture type boundaries
before installation, preserving independent C allocation/observation in particle
and tile-worklist tests. Raw storage observations use measured native offsets;
root assertions and frozen captures remain unchanged.

All seven storage captures equal frozen hashes. Complete default/highres roots
each pass 2,426 tests; server passes 2,415. All have only the expected prerequisite
probe skip, with exact root-name sets. Safe build/static checks and all three
production/ABI gates pass. The known suite retains exactly 304 failure events
and 17 passing, two failing, 32 skipped packages. Headless character creation,
save/load/resume passes. All phases have identical source fingerprints and all
1,654 original assets are unchanged. External SDL2/OpenGL/OpenAL bindings remain.
No complete safe runtime suite or safe raw-storage pass is claimed.

See [qualification](go-native-record-storage-qualification.json) and the
[dependency inventory](go-native-record-storage-inventory-after.json).

Local drafts, reconstruction records and probes:
`build/port-go-native-record-storage/`.

The first compile found one incorrectly pruned `common/flags` import: its Go
package is named `noxflags`, not the directory basename. Restoring the original
import fixes the fixture without changing behavior or expectations. This was
primary cleanup tooling, not Luna's six-file implementation. The failed build
remains under `storage/`; corrected-source qualification uses `storage-final/`.
Future import cleanup must resolve declared package names rather than infer them
from paths. No other non-C import was pruned in this candidate's cleanup pass.
