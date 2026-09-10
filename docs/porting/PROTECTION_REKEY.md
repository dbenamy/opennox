# Protection record rekey/shuffle — 2026-09-10

This chunk replaces `nox_xxx_protectData_56F5C0`, retaining its C export for
remaining protected-value mutation callers. The existing public Go wrapper
calls Go directly. Rekeying also removes the last production C calls to index
lookup (`sub_56F6F0`) and payload swap (`sub_56F720`), so those obsolete exports
and declarations are retired. Their useful record/counter tests continue through
native Go helpers; decoded-ID lookup still has C callers and keeps its export.

The new key is the server frame XOR one draw from the unchanged legacy floating
RNG. Each of floor(count/4) shuffle iterations draws two Logic indices with the
original inclusive bounds: [0,count/2] then [count/2+1,count-1]. Only payloads
swap; nodes, links and endpoints retain their identity. Each successful swap
increments its wrapping counter. Rekey XORs oldKey XOR newKey into both words of
each record, rebuilds the checksum from the complement of the new key, increments
the rekey counter, installs the key and returns its raw 32-bit value. Other RNG,
record count and handle sequence remain unchanged.

Floating-point RNG implementation and arithmetic remain in C. Tests obtain an
independent expected draw and post-state from that unchanged helper, restore its
pre-state, then compare the complete rekey operation against the expected key
and RNG state. This validates consumption and integration without asserting that
the historical RNG algorithm is correct or portable to a different architecture.

Primary-owned pure helper tests cover 1,000 deterministic lists with nil lists,
zero/unchanged/high-bit keys, exact encoded payloads, checksum and untouched links.
Terra prepares the original-C state fixtures and the primary reviews them before
production replacement. Final validation outcomes and source counts follow.

The original-C harness is recoverable at `8d59cf6b`. Its 400 seeded scenarios
run through both C and Go-wrapper paths (800 runs), covering 0/1/2/3-record
no-shuffle boundaries, 4/5/7/8 records and larger lists through 257. Tests check
exact payload order and original node identities, both RNG systems, zero keys
and frames, high-bit keys/returns and counter wrapping. Corrupt lists and
inconsistent counts remain outside the valid-state contract.

All accumulated protection tests pass on 386 under `porttest`, `server porttest`
and `highres porttest`; pure helpers pass on 386/amd64. All three production
targets build through `go run ./internal/noxbuild -o ../build/port-rekey/bin`.
Symbol checks confirm the rekey C export is Go-backed and both retired record
bridges are absent. Fresh `rekey-port` gameplay exits 0 against both preserved
screenshots with overrides disabled.

The asset-backed full default suite finishes with 15 passing packages, 3 known
failing packages and 32 without tests. Its exact failing package/test entries
match build/port-create/accepted-suite.jsonl: no added or resolved failures.
The known failures remain sprite rendering, obsolete blob parsing and audio PCM
hashes. Raw full-suite logs stay local because they can contain environment
credentials. Evidence is under build/port-rekey.

Production C: **142,265 physical lines**, down **62**, in 153 files; reference C:
**0**. See [C_LOC.md](C_LOC.md).
