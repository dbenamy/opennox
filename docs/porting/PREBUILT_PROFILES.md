# Prebuilt profile test trial

Status: bounded trial passed on the qualified `4a0ab0dc` engine/test source.
Use the optional two-process mode for the next complete root-corpus sweep.
Builds, broader package suites and headless gameplay scenarios remain sequential.

`tools/porting/run_profiles.py` builds all three root porttest binaries sequentially,
then runs at most two through the existing discovery/execution accounting driver.
Each binary record binds its path, SHA256, source fingerprints, supplemental
tracked source inputs (including assembly), root package and tags. Embedded Go
build metadata must match Linux/386/SSE2/cgo and the requested profile. Source
must remain unchanged throughout. This is a local frozen-build check, not a
general build cache: dependencies, toolchain and assets must also stay fixed.

The default `run_tests.py` mode remains available. Both modes now reject any
failure event, including package/subtest failures, even if a wrapper returns zero.
The controller joins all submitted runs before returning, records failed runs,
and refuses nonempty capture/diagnostic output variables. Each profile has its
own binaries, records, discovery output, runtime log and result file.

The source audit found root fixture writes in process-local temporary directories,
no fixed root-package listeners, and separate-package network tests outside this
scope. Process-global state is isolated between binaries. Native host services and
resources remain shared; this is not permission to parallelize `go test ./...`.
The VM has four CPUs and about 7.6 GiB RAM. A sequential full-profile sample peaked
near 1.35 GiB RSS, leaving room for a bounded two-process trial. That sample does
not prove a maximum for all possible fixtures.

All 43 porting-tool Python tests passed, including mocked stale binary/source/tag
rejections, complete execution accounting, failure reporting, build-before-run
ordering, the two-process cap, and joining after failure. The native trial passed
all 96 selected roots in each profile with no failures or skips and exact expected
name sets. Log timestamps verify overlapping server/highres execution. The
selection exercises audio, file/map fixtures, session entry and player reset.
This establishes the bounded trial, not full-corpus isolation or a general speedup.
Keep the next full sweep's failures visible and investigate before retrying.

Primary designed the orchestration and reviewed the implementation. Luna drafted
the prebuilt runner and mocked checks. Primary caught failure accounting that
initially ignored non-root events and required supplemental source fingerprints;
Luna caught the shared diagnostic-output variable. Both corrections preceded the
native trial. No measured subscription savings are claimed.

Evidence: [prebuilt-profiles-pilot.json](prebuilt-profiles-pilot.json).
Reproducible selection and environment:
[prebuilt-profile-pilot-batch.json](prebuilt-profile-pilot-batch.json).
Run with `run_batch.py ... --phase pilot --out <fresh-directory>` after sourcing
the port environment. Original local artifacts: `build/port-complete-corpus/prebuilt-pilot/`.
