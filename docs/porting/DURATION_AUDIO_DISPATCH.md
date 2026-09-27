# Duration-spell and internal audio-stream dispatch

## Scope and accepted baseline

Replace duration-spell and internal audio-stream foreign fallback calls and their
two C observer fixtures. Keep all 53 registered duration handlers and all20 native
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
frozen captures remain unchanged. Four-file conversion is reviewed but not yet
installed at this baseline checkpoint.

## Producer and semantic review

Duration constructor arguments come from native getter identities or nil in the
spell dispatch and magic-wall paths. Lifecycle fixtures supply those identities
or two isolated observer keys. Preserve independent result/void argument words,
counts, reset behavior and signed32-bit results. The new void observer's zero
return is discarded by its existing callers. Existing contracts cover nil/live
records, signed boundaries, ignored results, GC, mutation and real lifecycle owners.

Audio descriptor/voice-API blob slots93956..93988 and94012..94024 receive13 bridge
keys. Context/voice constructors supply four stream keys; event setup supplies
three event keys. Root service consumes callback offsets216/276/280/284. Fixture
owners install fourteen controlled operation keys in those same fields. Translate
those keys together with dispatch, preserving operation order, pointer identity,
signed result conversion, ignored results, nil-owner panic and root/hook restore.
Keep per-iteration closure capture semantics (module Go1.25) explicit in review.

For later review: unknown duration/audio keys receive explicit panics after the
producer audit; the duration registration API still accepts native callbacks.
Other callback families remain unchanged. Generic C signatures are shared with
those families, so this batch is expected to remove two fixture imports and four
C call sites while embedded-body count remains11.

Monster dispatch was considered but deferred after the audit found foreign
observers in the combat/lifecycle fixtures. Move that connected family together;
do not infer native-only producers from the stable production table alone.

## Evidence, delegation and recovery

Primary implementation/review; Luna quota remains unavailable. Ignored artifacts:
`build/port-single-pointer-dispatch/`, including exact source snapshots, draft,
producer searches, function comparison and run records.

Before native compilation, removed seven superseded UI-completion binaries after
committed-source/replacement/hash/host-use checks:387,067,904 allocated bytes.
Rebuild `defbdb98` for those historical outputs; logs/metadata remain. Exact paths
and records: `ui-cleanup-approved.json` and `ui-cleanup-deleted.jsonl` in the batch
directory. Current section/particle and modifier artifacts remain intact.
