# Native record storage and callers

Status: original baseline accepted at `fe44bab3`; reviewed conversion awaits
installation and runtime qualification. Expected removal: 11 production and six
test C imports, plus one unused specialized C adapter body. These are candidate
counts until measured after qualification.

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

After installation, run all seven storage captures, complete default/server/
highres root suites, safe build/static checks, three production/ABI gates,
exact known-suite comparison and headless character creation/save/load/resume.
Verify exact test names/skips, unchanged source, original assets and dependency
counts. Keep external SDL2/OpenGL/OpenAL bindings.

Local drafts, reconstruction records and probes:
`build/port-go-native-record-storage/`.
