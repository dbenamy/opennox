# Player admission, status and respawn

Qualified repaired C baseline after console conversion **a72dd6c6**.
Candidate selection: 22 routines /619 body lines in GAME1.c, covering player/team
eligibility, sudden-death admission state, name lookup, minimap tracking/fanout,
status messages, client respawn and displayed equipment. Current production C
is **33,268 lines /69 files /zero reference C** after the one-line prerequisite fix.

Two candidate scalar helpers (`sub_40AA20`, `sub_40AA30`) occur only in their C
definitions and header declarations across src; remove them as orphaned bodies
rather than translating/testing them. Audit other wrappers and callback reachability
before finalizing scope. C clients remain in the decoder and server lifecycle code;
Go callers should move directly to native helpers at conversion.

Reuse actual sparse-player, update-data, team membership and report queue owners.
The initial competitor-count matrix distinguishes player records from player units,
occupied teams from empty teams, observer/status0x20 semantics, missing units and
the dedicated host slot. Create membership through the real team owner, not by
setting a team ID alone. The original active-competitor team loop omitted its membership check. After fixing
an initial fixture omission (the shared owner has an active unitless fourth player),
the third C run reproduced 35 incorrect team-mode results and no non-team failures.
The C prerequisite correction adds the same real membership predicate used by the
neighboring per-team count. Empty teams now contribute zero; occupied teams count
once. This reversible behavior correction is recorded for review before conversion.

The sixth focused C run passes seven roots. Status contracts check mask0x423,
low16 object IDs, host gating, exact seven-byte payloads and broadcast recipients.
Timer contracts distinguish host/chat/participant/minute/previous-state gates and
advance time before a second observer clear to verify the deadline is not restarted.
Admission covers quest threshold6, nil callers, re-entry, positive signed scores,
frame wrap, selected special-mode bytes, remembered groups and team/flag limits.

The eighth C run passes 14 roots. Equipment contracts exercise both slot arrays,
full capacity, empty slots, unknown player/modifier IDs, mask narrowing and every
respawn mask bit across classes and game modes. Slots retain their sixth word.
Minimap contracts use real allocated circular lists, active-player fanout, repeated
mark/merged flags, partial removal, middle/head/final unlink, invalid indices and
nil objects. Lesson reset includes the unitless active player and emits messages
only for players with units. Whole-roster status includes every active player.

Queue assertions account for the actual reliable owner: insertion is at the head,
and the stored connection mask is distinct from the per-message destination. Two
initial fixture assumptions about those properties were corrected after inspecting
the queue implementation; no production change was needed.

Scope stays with these connected 20 live routines and two orphan helpers, rather
than adding unrelated GUI lifecycle code solely to reach a line-count target.
Completed console-native scenario copies were hash-verified against original assets
and deduplicated, recovering 1,660,044,319 bytes. Per-run restoration manifests
preserve paths, hashes and metadata; changed files, reports and original assets stay.

Final focused C qualification repeats exactly in two processes: 16 roots /18,542
tests including subtests; 16 captures /18,779 records. Frozen expectations are in
the tests. Added final boundaries include signed flag capacity, all re-entry word
bits, retained modifier pointers during unequip, missing-color early return, and
non-ASCII name comparison in the existing C locale. Static mapped-memory passes.
Three-target and fresh-production baseline qualification pass.

Default/server/highres each pass 316 roots /38,763 tests without skips. All185
captures /50,798 records are identical across targets, including the 16 focused C
captures. The three gates share 2,382 source files. Fresh production passes: three ELF32/SSE2/CGO binaries and ABI checks; exact known
full-suite results (1,553 failure entries; 15 pass /3 fail /32 skip packages);
options/gameplay, save/load and flat-map scenarios. All four gates share identical
source. No native implementation is installed yet. See
[player-state-c-qualification.json](player-state-c-qualification.json).
