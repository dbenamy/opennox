# Legacy online session state — Go conversion qualified

Qualified parent: **cc21d5db**. Production remains **11,898 C lines /50 files**,
zero reference C. The provisional scope has28 bodies /475 body lines; the
selection and textual caller inventory are tracked alongside this report.
The inventory's root-package `sub_420100` is an independent Go implementation,
not a reference to its same-named C function. That C function and sub_41D1B0
have no live callers.

Service list534808 and channel list529316 have no production writers beyond
zero initialization; user list531648 is only initialized/reset to zero. Account
index587000_60044 remains−1. Their former population code is absent. Whole-source
symbol/address audits are being completed before simplifying the unreachable
list/parser branches. The shipped12-record queue table at587000+58128 contains
zero pointers and counts. Its historical callback-address column is data, not
an active registration. No production queue-pointer writer has been found.

Live retry timers, session-status callbacks, connection-failure state, briefing
and map-transition words must retain their behavior. Contracts exercise28,800
start/attempt/retry/reset combinations across frame zero/wrap, exact1 versus other
nonzero flags, inclusive3600×FPS expiry and120×FPS retry deadlines. They also check
the exact108-byte login scratch clearing extent, status callbacks for all12 queue
states, empty-service behavior and signed/full-width word setters. Tests use the
shipped empty queues and actual fixed owners; they do not create unreachable list
contents solely to preserve dead production code.

Initial compilation caught a fixture-only cgo type-spelling mismatch for the
status owner; its declaration now matches the existing unsigned-int spelling.
No production algorithm change and no frozen new capture is claimed yet.

The first corrected C run passes both roots. The connected scope now also includes
the shared C network-log formatter: its only remaining callers are the three
selected retry paths. Contracts now check exact log ordering and signed frame
formatting, allowing the formatter and its private buffer/export to retire too.
The120-root affected selection includes browser, session-entry/dialog, briefing,
character-creation and game-statistics tests; inherited capture hashes stay frozen.

## Frozen C qualification

All120 affected roots pass in default/server/highres; inherited captures remain
unchanged. Both new captures match byte-for-byte across all targets and an
independent default repeat. Static memory checks pass. See
[baseline record](online-session-c-qualification.json) and
[affected manifest](online-session-batch.json).

Production-source fingerprint comparison against cc21d5db finds only the two new
porttest files. Its freshly qualified binaries/ABI, exact known full-suite outcome
and four scenarios are therefore reused for this C baseline. Native qualification
will require fresh builds and scenarios. All baseline jobs are joined.

Retain the six state words already exposed through live root-package setters or
getters, including the connection-failure marker; the wider Go-only session API
can be simplified after its callers are migrated. Closed, unreferenced counters
and fixed list/service/account owners do not need replacement Go storage.

## Native implementation review (qualification pending)

Baseline **f73daa4e** is pushed. Native code preserves the live session status,
retry state/timing, briefing/map markers and login scratch reset. Retry scheduling
retains the separate frame/FPS reads for its returned deadline and stored deadline;
logging keeps signed32-bit frame formatting and original call order.

The immutable empty lists, unselected account and empty legacy queues permit
retiring their unreachable traversal/parsing/allocation-dependent branches. Their
live failure/status callbacks remain. Statistics serialization is unchanged; only
its always-failing service lookup and unreachable abort path are removed. The
root-package queue/service wrappers retain their existing entrypoints with no-op
behavior where appropriate. The account-name getter remains an empty string.

Six state owners move to Go. Eleven closed or fixed owners are removed without
replacement storage. Only sub_41D1A0 remains as a selected C export for its actual
client-decoder caller;27 selected private interfaces and the former C log callback
are retired. The C log formatter's final three callers now format directly in Go,
so common__log.c and its private buffer disappear. All remaining-C retired-symbol
searches are clean. Working C is **11,409 lines /49 files /zero reference C**,−489.
Fresh native all-target and production qualification is pending.

### Full-suite tooling correction

The first native production run built all three binaries and passed ABI checks,
but the full suite exposed a stale `TestReadMemmap` minimum of1,396 variables.
Retiring17 mapped owners correctly leaves1,379, so that assertion can no longer
measure parser correctness. Replace it with semantic round-trip preservation of
the actual mappings and an independent four-entry fixture covering address order,
zero offsets, widths, disabled entries and custom comments. Writes remain confined
to temporary source trees. Focused tooling tests pass.

This changes only `src/internal/blobs/memmap_test.go` after the successful native
120-root sweeps. Fingerprint comparisons prove production and affected-test source
are identical, so those sweeps remain valid; rerun the full production/full-suite
qualification on the corrected tooling test. No C or native behavior and no frozen
capture expectation changes. The failed run remains explicit evidence, not a
qualified checkpoint.

## Final native qualification

[Native qualification](online-session-native-qualification.json) records120 roots
and67 exact manifest captures per target, with no skips. The affected source is
identical across runs; the subsequent isolated blob-tooling test change is
recorded explicitly. Static checks and the focused tooling tests pass.

Fresh default/highres/server binaries pass ELF32/i386/SSE2/CGO, retained/retired
symbols and absence-of-test-helper audits. The full asset suite exactly matches
all1,553 known failure entries and15 pass /3 fail /32 skip package outcomes. Fresh
dialog, gameplay and save/load scenarios pass against qualified parent references.
All jobs are joined. Production C is **11,409 lines /49 files /zero reference C**,
**−489** from baseline f73daa4e. No remaining native correction or user question.
