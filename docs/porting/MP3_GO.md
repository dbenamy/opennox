# Production Go MP3 integration

Status: qualified; see PORTING_STATE.md and the machine-readable
[qualification report](mp3-go-qualification.json). The complete Go decoder is now selected by the production audio
wrapper. The six-line C implementation file and 1,890-line minimp3 header have
been removed from the working tree. This leaves zero standalone production C
files, but does not remove the remaining C preamble adapters or external cgo
libraries.

## Behavior and source recovery

The audio wrapper retains its 48 KiB input buffer, refill/compaction rules,
per-channel sample count, and reset-on-seek behavior. Decoder workspace belongs
to each stream and is cleared before each Layer III frame. Persistent reservoir,
MDCT overlap, QMF and header state retain their original lifetime.

The generated frame contracts and original initialization findings are described
in [MP3_FRAME.md](MP3_FRAME.md). In particular, Go deterministically initializes
scratch; it does not reproduce the original uninitialized private-bit/SCFSI
behavior. The legacy 8 kHz mixed-block workspace extent remains explicit and
bounded in Go. No audio golden has been updated for this integration.

Historical capture tools recover the retired header into ignored capture output
from Git commit `6702100ba617a26cf952b2f8ba629a6f4256426e`, verifying SHA-256
`e03e87c847cbdfd79cd3bf74a3b1fe24bc8c590245e92ce8153e87b23e17507f`.
A checkout used to regenerate captures must contain that history; frozen Go
fixtures do not require C source recovery.

## Performance review item

A controlled benchmark uses memory-resident `Dialog/C1CAP01E.WAV` MP3 data,
100 repetitions (81,300 frames), three trials per implementation, alternating
execution order and one pinned CPU. It excludes compilation, startup and file
I/O. Both implementations produce the same frame count and sample checksum.

Initial Go decoding allocated 16 KiB of scratch per frame. Moving scratch onto
each decoder removed all allocations during the timed loops and reduced the
observed median from 37.29 to 33.29 microseconds per frame. The corresponding
post-change C median was 11.92 microseconds: Go remained 2.79 times slower.
These are measurements on one mono 22,050 Hz asset in the shared VM, not a
whole-game performance result. Further optimization is a reversible follow-up;
exact sample/state contracts must continue to pass.

Ignored benchmark source, binaries, input identity and before/after measurements
are under `build/port-mp3-go/bench`. Production acceptance includes fresh
full build, consumer, known-suite and headless scenario checks.

## Review and disk space

Luna independently reviewed the caller buffer/seek semantics, decoder scratch
lifetime and historical C recovery. No integration regression was found; its
recommendation to use the full Git commit ID was applied. Primary retains
responsibility for acceptance and qualification claims.

To make room for qualification, 76 old reproducible Go build-cache archives
were removed after stat/hash revalidation, recovering
1,622,847,320 bytes. The removal record is
`build/port-artifact-cleanup/go-cache-sep19-removed.json`. Assets, frozen fixtures
and retained production binaries were preserved. Subsequent builds may recreate
the removed cache entries.

Cleanup process deviation: the cache open-file scan saw sandbox processes only.
All known build/benchmark jobs were joined, but this did not meet the documented
host-namespace check. Later scenario deduplication uses the host process namespace.
The removed cache outputs are reproducible; no sole-copy evidence was deleted.

## Completed qualification

- All ten decoder test roots pass default/server/highres/safe and cgo-disabled
  configurations; vet passes.
- All 1,246 shipped MP3 guarded observations match the original C SHA exactly,
  in default/highres. Historical PCM goldens pass unchanged in both profiles.
- Default/highres consumers and safe/static checks pass.
- Fresh client/high-resolution/server builds and all four binary ABI checks pass;
  C decoder symbols are absent and client/safe binaries contain the Go decoder.
- Full-suite comparison matches 304 known failure events exactly, with
  17 passing, two failing and 32 skipped packages.
- Headless character creation/settings and save/load match existing references.
- Recovering original C from Git and rerunning integer capture reproduces the
  committed compressed fixture byte for byte.

No expected audio or existing failure result was changed. The suite manifest
adds only the new MP3 package's passing result. Standalone C falls six lines to
zero; retiring the 1,890-line header is an additional change outside that count.
