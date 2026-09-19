# Client audio streams, buffer cache and driver queues — baseline in progress

Qualified parent: **eeaf1030**, client world/spell/item presentation.
Parent production C: **23,648 physical lines /66 files /zero reference C**.
Current prerequisite corrections: **23,661 /66 files /zero reference C** (+13).

The connected batch contains 70 C bodies /1,109 original body lines
(1,122 after prerequisite corrections): sample metadata and
bag/WAV reading, driver/context/voice queues, timers and callback dispatch, fixed
block pools, cached sample chunks and references, and voice buffer traversal.
Scope: GAME2_2 addresses 486640–487D60 and GAME3_1 addresses 4BD280–4BDC00.
See client-audio-streams-selection.json and client-audio-streams-callers.json.

Whole-repository caller/callback audit finds 25 external roots and all 70 bodies
reachable. Roots include client audio initialization, AIL adapters, frame service,
audio events and cleanup. No orphan pruning is proposed. The first provisional
parser included three indented call sites as definitions; the corrected selection
requires a declaration at column zero and asserts unique names before reachability.
The discarded provisional graph must not be used.

Reuse actual audioStructXxx/36-byte sample entries, legacy file handles, intrusive
lists, timers and allocation owners. Independent contracts should cover pool
exhaustion/reuse, reference narrowing/clamping, sample format and partial reads,
WAV override/fallback and close ownership, cache eviction/refcounts, buffer lists,
driver/voice capacity and ordering, callback failure cleanup and timer wrap.
Use deterministic callback observations and real ownership; physical playback
quality remains outside the headless qualification.

The test-only dispatcher for all 70 actual C functions is installed, along with
pool allocation/reuse and reference-count contracts (initial 2 roots /43 records
passed). Format, ownership, constructor failure, buffer and real file-read contracts
are being added. The first owner build exposed a test-only callback address linker
issue; a C accessor now returns the static callbacks. owner-io-second is running. No goldens are frozen and production source is unchanged. Reuse the
parent production baseline only after verifying that identity; any prerequisite
production fix requires fresh qualification.

Before starting these tests, 164 regenerable Go compiler-cache artifacts older
than six hours were removed with no compilers active, reclaiming 8,618,371,198
bytes. The manifest is build/port-client-audio-streams/removed-stale-go-cache.json.
The cleanup is consumed; ordinary builds regenerate missing compiler artifacts.
Original assets, archive, qualification captures and current production binaries
remain untouched.


## Prerequisite corrections under test

The original `sub_486B60` used only 36 bytes of its local array for the override
path. `owner-io-second` passed the bag read, then aborted on an ordinary longer
fixture path with the compiler's buffer-length check. The independent path test
is retained. The corrected C uses a separate 280-byte path and bounded formatting,
covering the catalog's directory, sample name and extension; excessive paths fall
back to the bag. The same parser now requires a complete format chunk and nonzero
channel count before using format metadata. Missing, short and zero-channel format
fixtures require fallback with the override closed and original bag metadata intact.
These are reversible compatibility corrections to review, not frozen old behavior.
Because production C changed, parent production qualification cannot be reused.

`owner-initial` failed to link the test-only static callback addresses; the C
accessor fixed that bridge. `owner-io-second` passed buffer traversal before the
path abort. Both sessions are joined. `corrected-c-initial` includes pool, formats,
buffers, real reads, driver/context/voice lifecycle and failure cleanup, sample
cache eviction, voice callback/selection contracts and 64-bit tick boundaries.
All ignored `*-draft.go` files installed so far are consumed; tracked source wins.
No original-C expectations have been frozen yet.


The original C clock adapter returns `C.uint`; its result is narrowed to 32 bits
before the audio context's 64-bit elapsed arithmetic. This is preserved, including
high-word and wrap cases. The first tick fixture controlled only `timer.PlatformTicks`;
controlling `legacy.PlatformTicks` separately fixed that fixture failure. All twelve
then-installed focused roots passed in `clock-owner-fixed`; the additional device
queue and bulk allocation contracts passed in `queues-initial`, and signed/short-read
contracts passed in `read-bounds-initial`. The final dispatch/format sweep is active.

The cache's unsuccessful fill can retain the active catalog stream until its next
open or close. Tests record this existing ownership behavior; catalog cleanup and
subsequent reuse close it. No extra cache behavior change is proposed. Bulk voice
creation with a zero request fills available capacity, due to the existing do/while
contract; valid positive requests are used by production initialization, and both
behaviors are covered. Allocation counts for low-level pools remain strictly positive,
as required by actual callers. Artificial multi-gigabyte allocations are not used
for signed-length tests: the offered storage covers each effective read request.

Affected qualification includes the connected audio asset and stream roots,
object-report/world-collision audio contracts, intrusive lists/player groups and
existing timer package tests. Full production qualification compares the exact
known full-suite results and fresh audio-enabled gameplay/save/load scenarios.


All three focused targets now pass with identical source fingerprints and identical
16 captures /927 records. The expectations are frozen in the fixtures and indexed
in client-audio-streams-captures.json. All sessions are joined. freeze.py and the
installed drafts are consumed. This is a recoverable C checkpoint; affected target
sweeps and fresh C production remain before native implementation.
