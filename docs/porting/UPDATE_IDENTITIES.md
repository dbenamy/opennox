# Object update callback identities

Scope: replace 53 named object-update callback addresses and three special
identities with stable native Go keys. Retire their 56 C export bridges and
migrate direct owner calls and identity consumers. Keep the 386 layout, resource
names/data sizes, original return widths in direct-call fixtures, capture IDs,
mutable hook lookup and unknown raw-C callback fallback.

Named registrations already use the Go update cache. A thin anonymous registration
entrypoint will bind the player-observer, scheduled-spell and monster-bot identities
to that same cache without inventing resource names or changing dispatch logic.
Root player/observer/pixie/harpoon fast paths retain their existing meaning.

## Baseline and coverage

The original baseline adds three independent contracts: lifetime and uniqueness
of all 56 keys through foreign object storage and GC, late replacement of the two
special mutable hooks, and full normalized state equality between direct and
stored-callback bot updates across all player classes and four frame values.
The bot fixture's original default behavior is unchanged; the two added routes
both discard the result exactly as Object.CallUpdate does. Existing direct-owner
captures continue to check the original result separately.

The 279-root selection includes update-registry captures and raw fallback/hook
contracts; temporary effects; world mechanisms, motion, geometry and collisions;
objectives; generators; unit gameplay; player controls; server orchestration;
and spell-start/lifecycle owners. Freeze exact selected names in all three
profiles before conversion. Existing assertions and captures must not change.

Production source is identical to the preceding qualified transfer conversion;
reuse that production baseline. After conversion rerun focused profiles,
safe/static checks, all three production/ABI builds, the exact known asset suite,
and fresh headless character creation with explicit save/load/resume. Reassess
full-corpus qualification if implementation changes shared dispatch beyond the
thin registration entrypoint or a failure leaves the affected scope uncertain.

One Luna helper prepares a bounded ignored overlay. Primary reviews each owner,
conversion width, key consumer and fixture mapping before installing it.

The first baseline attempt exposed a setup error in the new bot contract: using
the generic timing runner omitted AI owner initialization. Use the existing
controlsRun environment for both direct and stored-callback cases. This preserves
the actual bot/AI owner and requires no change to production or frozen captures.

## Frozen original result

All 279 selected roots pass in default/server/high-resolution profiles without
skips under `baseline-fixed/`. Exact discovered/completed root sets and source
fingerprints match. Production source is unchanged from qualified commit
`1b14d1a3`; only the two documented porttest files differ. The initial fixture
setup failure is retained under `baseline-original/` and is not accepted evidence.

[Original baseline](update-identities-baseline.json),
[test selection](update-identities-tests.txt),
[qualification manifest](update-identities-batch.json).

Before baseline compilation, verified unused create/init and damage executables
were removed (14 files, 786,116,608 allocated bytes). The completed transfer run's
1,654 duplicate asset files were also removed (559,902,720 bytes), with originals
preserved and a restoration manifest. See PORTING_STATE.md for recovery paths.

## Qualified conversion

All 279 focused roots pass in each default/server/high-resolution profile with
no skips. Safe/static checks, all three production/ABI builds, exact known asset
suite comparison and fresh headless character creation/save/load/resume pass.
The known suite remains 304 failure events, 17 passing/two failing/32 skipped
packages. All 1,654 original asset hashes are unchanged.

All accepted phases share identical source fingerprints. All 29 changed/new/
deleted files match the reviewed overlay; all 56 selected C exports are absent,
and retained exported signatures/bodies are unchanged. No root assertions or
frozen captures changed. Primary independently checked all 56 owner calls, all
53 registry rows and all 43 sparse direct-operation identities. Signed-byte,
integer and pointer-word returns preserve their original fixture semantics.

Production cgo files fall 202→198 (265/463 eliminated); legacy C exports fall
795→739 (1,151/1,890 retired). Four additional fixture files drop cgo. Headers
remain 157 files with 3,471 physical lines; embedded production callback bodies
remain 77. Standalone production/test-reference C remain zero. External native
library bindings are unchanged.

The existing named registration and dispatch/cache algorithms remain unchanged.
The three anonymous callback keys use a thin registration entrypoint into that
same cache. Root fast paths, mutable hooks and unknown raw-C fallback remain.
The preceding item full default corpus passed 2,441 roots plus one diagnostic
skip. This batch uses focused owner/raw-fallback/lifetime/hook contracts and
production qualification; it does not introduce a shared dispatch algorithm.

Luna prepared the bounded overlay; primary reviewed the owner/consumer/fixture
mapping and requested the thin anonymous API without an additional nil policy.
The baseline setup correction is described above. Conversion qualification
required no behavior or capture changes.

[Qualification](update-identities-qualification.json),
[inventory](update-identities-inventory-after.json).
