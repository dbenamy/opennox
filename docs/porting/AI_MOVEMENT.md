# AI movement actions — 2026-09-11

Connected batch: random walk (545020 and private 545090), confusion (545140),
face location/object/angle and set angle (545210/545240/545300/545340/5453E0).
Six action registrations currently go from Go to C. All char returns are discarded
at the action interface; internal returns only forward them. Native registration
can remove all eight C bodies/entry points together. Audio and combat capability
callees, and 534120's other C callers, remain outside this batch.

One guarded C-owned object/monster/definition/target fixture and synthetic grid
and direction table cover the entire batch through the real action registry.
The C baseline captures every changed object/monster word, both directions, RNG
indices and stack-change state. Read-only data and guards are checked; global
server binding, grid, direction table and engine flags are restored. Golden
hashes retain those complete per-case state records in compact form. During
conversion the real C dispatcher and native registration also compare every case
directly, reporting the first differing input/state. The C dispatcher is temporary.

Corpus: all 65,536 current/target angle pairs; all 65,536 low-16-bit set-angle
values; all 65,536 signed starting walk directions with running/terrain variants;
5,120 facing point cases; 8,000 confusion seed/capability/stack combinations;
54 start/end/cancel cases; one explicit full-expression rounding discriminator.
The initial baseline passes in 4.2 seconds, including registry-vs-C comparison.
Finite boundary, zero-distance and nonfinite facing vectors are included. Facing
current/target table indices are valid 0..255; arbitrary outside-table addresses
are outside this contract. Set-angle and random-walk wrap are separately exhaustive.

Precision: PC53 double deltas; float32 length and normalized X; normalized Y
remains double for the location turn but spills for the alignment dot test. The
dot test spills its first product across a C table lookup. Actual mutable table
reads remain authoritative. Rounding direction*30 before adding position changes
the explicit terrain case and must fail. Disassembly and logs: build/port-ai-actions.

Current stage: original-C baseline captured; native batch not yet adopted. No
production C removed since 5ef3e3d6 (140,442 lines). One shared qualification cycle
will follow direct parity, instead of eight separate helper cycles.
