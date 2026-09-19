# Client audio events and playback — baseline in progress

Qualified parent: **50f0711d**, client audio streams/cache/driver queues.
Qualified parent C is **22,400 physical lines / 66 files / zero reference C**.
The current test-only device adapter adds ten physical C lines (22,410); no
production algorithm has been converted in this batch.

The connected batch contains 62 C bodies / 1,182 body lines: audio event setup,
selection, scheduling and playback; music/dialog controls; AIL sample refill and
format mapping. Scope spans GAME2.c, GAME1_3.c and client__audio__audevent.c.
Adjacent GUI overlay helpers were excluded after inspecting their behavior.
Whole-repository caller/callback audit finds 35 external roots and all 62 functions
reachable. See client-audio-events-selection.json and client-audio-events-callers.json.

Reuse the real stream/cache/driver owner, audio asset tables, file handles, timer
and music modules. Build independent contracts for default event data, flags and
format mapping, music queue limits/restore order, event ownership and pool reuse,
sample selection/RNG, scheduling and spatial volume, and sample refill boundaries
with callback/PCM observations. Preserve device API ownership when reusing Go AIL;
physical playback remains outside headless qualification.

Production algorithms are unchanged. The test-only adapter changes the literal C
source; perform fresh C production qualification before conversion.
No expectations are frozen and no conversion has started.


The 62-function test dispatcher and owners for 13 actual C globals are installed.
Initial format/volume/pan/switch contracts passed; music queue/default contracts
also passed. A test-only device observer now controls readiness and records loaded
PCM buffers while the actual C refill/EOS bodies and real voice/chunk owners run.
It falls back to the existing AIL backend outside an owned case. Its C aliases are
inside `NOX_PORT_TEST_AUDIO_EVENTS`; production behavior is unchanged, but the C
file has a new conditional adapter block, so baseline source identity needs a
fresh audit before any production reuse.

The first refill sweep passed normal boundaries and scratch guards; an interior
empty chunk ended playback after the preceding data. The production cache emits
nonempty chunks. Preserve this boundary convention, including the initial empty
chunk advancing once before the next result is checked; do not change playback
semantics to make concatenation include later chunks. The fixture now asserts it.
The sample-selection fixture uses the actual random generator, checks loop limits,
shuffle ordering, and consumption of the Other stream without changing Logic RNG.

No expectations are frozen. Production algorithms are unchanged. Installed
`build-dispatcher.py`, music/selection/owner/manager/callback drafts are consumed;
use tracked source.

Event ownership fixtures now exercise the actual 200-event pool, metadata tables,
cache/catalog, stream contexts, voices and callbacks. `lifecycle-third` passes all
11 root tests, including pool exhaustion/automatic reclaim, per-sound limits,
manager frames, stale handles and scheduling (32-bit clock narrowing, unsigned
64-bit deadline comparisons, exact-deadline waiting). `static-baseline.log` passes.

The recording device does not deliver asynchronous completion. The first lifecycle
sweep reused voices without delivering end notification and failed its third start;
the fixture now invokes the real stream end notification after stop, before reuse.
Production behavior was not changed. The manager fixture's first compile used the
wrong memmap accessor signature; corrected to PtrOff before running the contracts.

Serial zero is both the first allocation's serial and the cleared serial. A handle
to that first stopped event still validates until reuse changes its identity; this
is captured as existing behavior rather than silently changed during translation.
Callback transitions, cache reload/shuffle modes and priority eviction contracts
are being added. Expectations are still unfrozen.

## Frozen C checkpoint

All 17 focused captures /820 records match byte-for-byte across separate default,
server and highres runs from identical source. They are frozen in the tests and
indexed in client-audio-events-captures.json. `freeze.py` is consumed.

The affected sweeps `c-final-{default,server,highres}` pass 51 root tests plus five
timer tests each. All 36 captured artifacts /2,353 records match, and all 2,583
source fingerprints are identical. This includes the preceding stream/cache,
audio-asset, intrusive-list and affected game caller contracts.

The public playback fixture confirms immediate volume updates but pan target
updates take effect on the next timer update. Its first assertion expected pan
to change immediately; it was corrected against the existing timer adapter before
freezing. No production change was made. All fixture drafts, including priority,
playback and public-entry drafts, are consumed.

Fresh C production qualification is running in `c-final-production`; this recovery
checkpoint does not yet claim production qualification. Native type/selection drafts
under build are not installed or qualified. The interface audit currently retains
18 actual C entry points plus three stored voice callbacks (21 total); 41 can
become private Go functions. Review the complete adapter return conventions too:
the existing Go timer adapters return nil/zero for raw/interpolated setters.
