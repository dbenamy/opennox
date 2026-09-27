# Duration-spell and internal audio-stream dispatch

## Scope and accepted baseline

Replace duration-spell and internal audio-stream foreign fallback calls and their
two C observer fixtures. Keep all 53 registered duration handlers and all 20 native
audio cases unchanged, including late hook reads, abort and no-op behavior.
Keep duration KeepAlive calls, layouts and external audio backends unchanged.

Baseline source: qualified section/particle revision `6f03cc39`. Fresh original
runs reused its verified binaries after exact production/test source and environment
checks: **118 owner roots in each of default/server/highres**, plus six separately
repeated focused contracts per profile. Three safe duration contracts passed.
No skips. Production baseline is reused from identical source; all production
gates run again after conversion. See [baseline](single-pointer-dispatch-baseline.json).

Selected owners: duration and sustained spells, spell start/lifecycle/effects,
AI spell contracts, root audio owners, client duration effects and native layouts.
This affected-owner selection is not a full root-corpus claim. Assertions and
frozen captures remain unchanged. Four-file native conversion qualified against baseline `aa7fac01`.

## Producer and semantic review

Duration constructor arguments come from native getter identities or nil in the
spell dispatch and magic-wall paths. Lifecycle fixtures supply those identities
or two isolated observer keys. Preserve independent result/void argument words,
counts, reset behavior and signed 32-bit results. The new void observer's zero
return is discarded by its existing callers. Existing contracts cover nil/live
records, signed boundaries, ignored results, GC, mutation and real lifecycle owners.

Audio descriptor/voice-API blob slots 93956..93988 and 94012..94024 receive 13 bridge
keys. Context/voice constructors supply four stream keys; event setup supplies
three event keys. Root service consumes callback offsets 216/276/280/284. Fixture
owners install fourteen controlled operation keys in those same fields. Translate
those keys together with dispatch, preserving operation order, pointer identity,
signed result conversion, ignored results, nil-owner panic and root/hook restore.
Keep per-iteration closure capture semantics (module Go 1.25) explicit in review.

For later review: unknown duration/audio keys receive explicit panics after the
producer audit; the duration registration API still accepts native callbacks.
Other callback families remain unchanged. Generic C signatures are shared with
those families, so this batch removes two fixture imports and four C call sites while
embedded-body count remains 11.

Monster dispatch was considered but deferred after the audit found foreign
observers in the combat/lifecycle fixtures. Move that connected family together;
do not infer native-only producers from the stable production table alone.

## Evidence, delegation and recovery

Primary implementation/review; Luna quota remains unavailable. Ignored artifacts:
`build/port-single-pointer-dispatch/`, including exact source snapshots, draft,
producer searches, function comparison and run records.

Before native compilation, removed seven superseded UI-completion binaries after
committed-source/replacement/hash/host-use checks: 387,067,904 allocated bytes.
Rebuild `defbdb98` for those historical outputs; logs/metadata remain. Exact paths
and records: `ui-cleanup-approved.json` and `ui-cleanup-deleted.jsonl` in the batch
directory. Current section/particle and modifier artifacts remain intact.

## Native qualification and counts

Matching original/native selections passed **118/118/118 owner roots**, six
separately repeated focused contracts per profile, and three safe contracts.
Safe build/static, three production builds/ABI, exact known-suite outcomes and
fresh headless save/load passed. All 1,654 original asset hashes are unchanged.
Retired fixture symbols are absent. Existing assertions/captures, all 53 duration
registrations, all 20 native audio cases and duration KeepAlive calls are unchanged.
No native source correction was needed during qualification.

Fixture C imports: **91→89**. Embedded production C bodies remain **11**, since
other callback families still use the same signatures. Production cgo remains
**5 client/highres, 6 server** across three project packages. Selected legacy
exports remain zero; headers remain 157 files /2,731 lines. Standalone C remains
**0 production /0 test**. Evidence:
[qualification](single-pointer-dispatch-qualification.json),
[inventory](single-pointer-dispatch-inventory-after.json).

After all qualification jobs finished, removed 1,654 verified duplicate scenario
assets (560,214,016 allocated bytes), preserving originals, saves and results.
Restore before replay: `python3 build/port-artifact-cleanup/restore-recent-scenario.py
build/baseline/runs/single-pointer-dispatch-save/deduplicated-assets.json`.
