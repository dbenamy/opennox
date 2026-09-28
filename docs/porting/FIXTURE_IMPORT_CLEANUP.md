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

## Status

Original qualification accepted: verified broad supersets and fresh eight focused
roots per profile, plus two safe contracts, all passed. The draft is not installed. Expected fixture
C-import count: 69 to 51; production counts remain 4 client/highres and 5 server.
Standalone C, production embedded bodies and legacy exports remain zero.

Luna remains unavailable under its existing quota; primary reviewed and prepared
this batch locally. Artifacts: `build/port-fixture-import-cleanup/`.
