# Go MP3 synthesis performance

The synthesis path now accumulates one independent lane at a time, retaining
its original tap order and explicit float32 rounding. For mono it also omits
calculations whose PCM stores were overwritten by the following even-lane
stores. Every history/workspace store remains; stereo keeps all lanes and stores.
No tables, decoder interface, fixture, or audio golden changed.

A CPU profile of 813,000 frames attributed 66.81% cumulative / 54.05% flat time
to synthesis. Luna prepared bounded drafts; primary independently checked the
arithmetic order and PCM index overlap, then qualified the integrated changes.
The optimization assumes disjoint PCM and float workspace, as used by the
production decoder and fixtures. Deliberate unsafe aliasing and floating-point
exception flags are outside that contract. Existing documented unused-QMF and
original private-bit initialization limits still apply (see MP3_FRAME.md).

## Measured scope

Native 386/SSE2, GOMAXPROCS=1, one pinned CPU, memory-resident input, 100 repeats
(81,300 frames) per trial, three rotated-order trials per implementation. The
unchanged Go binary was rebuilt from preceding revision `81e5291d`; C uses the
retained original decoder benchmark. Medians exclude compilation and input I/O.

| Input | Original C µs/frame | Unchanged Go | Optimized Go | Go time reduction |
| --- | ---: | ---: | ---: | ---: |
| Shipped mono 22,050 Hz `Dialog/C1CAP01E.WAV` | 11.098 | 32.849 | 21.771 | 33.7% |
| Synthetic stereo MPEG1 44,100 Hz | 23.372 | 65.282 | 59.923 | 8.2% |

Both Go versions allocate zero bytes per benchmark loop. Go remains 1.96× C on
the mono input and 2.56× on the synthetic stereo input. These are bounded decoder
measurements, not estimates for all assets or whole-game performance. The
stereo stream alternates existing generated coded/escape frame patterns; it is
not shipped music. Benchmark checksums cover first/last samples per frame;
separate exact PCM/state tests provide correctness evidence.

## Qualification

[Machine-readable record](mp3-synthesis-performance-qualification.json) includes
source fingerprints, phase commands/results, benchmark trials and binary hashes.
[Batch manifest](mp3-synthesis-performance-batch.json) records reproducible checks.

- All ten frozen decoder test roots pass default, server, highres, safe and
  cgo-disabled configurations; vet passes. Existing synthesis coverage includes
  197,774 C records and frame coverage includes 2,682 calls in 544 sequences.
- Both audio profiles retain all 1,246 shipped asset observations and historical
  PCM goldens. Guarded capture SHA-256 remains
  `e0688114198dc50b4acf9b5695ea4cd8bf6dbb9f91d53d1e9d550c02ff1e2763`.
- Default/highres raw-address audio consumer tests pass without skips.
- Safe/static and fresh production/ABI checks pass. Full suite matches exactly
  304 existing failure events, with 17 passing / 2 failing / 32 skipped packages.
- Headless character creation/settings and save/load comparisons pass.

Only synthesis.go differs in production source from the preceding qualified
revision. The default decoder check was reused from the exact source used in
the benchmark; four remaining profiles ran in the batch. A manifest path typo
in an unused audio phase was corrected while the decoder phase ran, before the
audio phase started. The finalizer verifies each executed phase specification,
environment, production specification and source fingerprint; it does not claim
that the decoder's earlier snapshot of unused commands was byte-identical.

## Artifacts and recovery

Working artifacts are in `build/port-mp3-performance`: `cpu.pprof`, `top.txt`,
`synth-profile.txt`, `baseline.json`, `synthesis-before.go`, the three benchmark
binaries, `comparison.json`, `comparison-stereo.json`, test outputs and scenario
reports. Profile-instrumentation allocations are not decoder allocations.
The committed qualification embeds the benchmark results for VM-independent
review. Frozen tests, capture generators and production qualification tooling
remain committed; ignored profiling/timing scratch is optional recovery data.

Scenario duplicate assets were removed only after successful comparisons, with
restore manifests. Twenty-one old cache archives (509,121,918 bytes) were removed
after exact stat/hash verification and host-namespace process/open-file checks
(PID1 systemd); record `build/port-artifact-cleanup/go-cache-pre-sep23-round2-removed.json`.
Original assets, frozen fixtures and retained binaries remain intact.

Standalone production and test-reference C remain **zero files / zero lines**.
Production C preamble bodies remain **79** (76 shared dispatchers, three typed
invokers); external libraries and generated bridges remain outside that metric.
Next bounded review: remove redundant local C round trips to already-Go book
callbacks, preserving initialization, integer normalization and exported ABI.
