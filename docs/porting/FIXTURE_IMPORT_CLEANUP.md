# Fixture import cleanup

## Scope and review

Remove C imports from eighteen porttest fixture files: fifteen include-only
preambles with no C selectors, definitions, exports or compiler directives, and
three self-contained scalar bridges (AI state, network code and world geometry).
Replace C.int/C.float with int32/float32, retaining casts, arithmetic, argument
order, pointer bit patterns, captures and assertions. Production code and allocation
ownership are unchanged. The 28-header transitive project include closure contains
no startup hooks. AST review: 14 changed functions, 50 unchanged, none added/removed.
AI private wrapper callers are confined to the same fixture file.

## Qualification plan

Reuse the preceding callback batch's broad original results only after exact
source, supplemental source, environment, binary/log hash and root-name verification.
Run eight scalar-owner roots freshly in each profile, plus two safe network-code
contracts. Commit the accepted original baseline before installing the draft.

After conversion, compile all root tests in default/server/highres. Execute the
broad default selection (2,482 roots) and the eight scalar-owner roots in each
other profile; repeat the eight roots in all three profiles. The earlier broad
server/highres results are original supersets, not native coverage claims. Preserve
the preceding eight broad exclusions. Run safe contracts/build, static mapped-state
checks, all production builds/ABI, exact known-suite comparison and fresh save/load
against unchanged references. No extra GUI preflight is needed for test-only edits.

This asymmetric scope follows the actual change: removing unused includes requires
compile/link coverage; explicit-width scalar substitutions are exercised across
all targets. Existing boundary contracts provide independent checks; no mirrored
implementation tests or regenerated captures are introduced.

## Qualified result

All gates passed without source corrections. Native roots passed without skips:
2,482 default, 8 server, 8 highres; eight focused repeats per profile and two safe
contracts also passed. Static checks, safe build, three production builds/ABI,
exact known-suite outcomes and fresh save/load passed. Original assets and all
root test assertions/captures are unchanged.

Fixture C imports fell from 69 to 51. Production imports remain 4 client/highres
and 5 server across two project cgo packages; legacy exports, embedded production
C bodies and standalone production/test C remain zero. Headers remain 157 files /
2,731 physical lines. This does not qualify other architectures or remove external
native backends.

Evidence: [baseline](fixture-import-cleanup-baseline.json),
[qualification](fixture-import-cleanup-qualification.json),
[inventory](fixture-import-cleanup-inventory-after.json).

Luna remained unavailable under its existing quota; primary reviewed and completed
this batch locally. Artifacts: `build/port-fixture-import-cleanup/`.

## Artifact recovery and disk cleanup

Removed seven superseded damage/monster test/safe/production binaries after
verifying source revision `945c190b`, their qualified callback replacements and
host open-file/mapping checks: 386,945,024 allocated bytes reclaimed. Rebuild that
revision to recover them; logs and source records remain. Current baseline and
conversion binaries are retained.

The completed scenario's 1,654 unchanged asset copies were verified against the
originals and removed, reclaiming 560,476,160 allocated bytes. Restore with:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/fixture-import-cleanup-save/deduplicated-assets.json
```

Original assets, scenario outputs and restoration manifests remain.

The seven root/safe/production binaries for this batch were later removed after
qualified replacements and committed source were verified. Rebuild 6b51667b;
[fixture constants cleanup](FIXTURE_CONSTANTS.md#artifact-cleanup-and-recovery)
records the checks. Reports and test logs remain.
