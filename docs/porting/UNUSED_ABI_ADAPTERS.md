# Unused engine ABI adapters

## Scope

Retire 56 unused C export wrappers and 33 private C-typed wrappers after a whole
source audit. Their remaining references are declarations, comments describing
native owners, transpiler mappings or token-test input strings. Followed macro
aliases have no runtime callers. Keep native owners, live callbacks, allocator
hooks, shared records/layouts and test fixtures unchanged.

The draft removes 54 prototype occurrences, four macro aliases, 11 unnecessary
cgo imports and one empty Go file. It also cleans stale comments and blank lines
in touched headers. An AST comparison confirms that 668 retained function bodies
and signatures are unchanged. Root functions with legacy names are distinct
native implementations and remain; no C algorithms are kept solely for tests.

## Baseline and acceptance

The original source is the qualified animation revision `f8332e3d`. Source/test
fingerprints and environment settings must match exactly before reusing results.
The selected 357 default roots come from its full 2,475-root corpus, excluding the
established prerequisite skip. For server/highres, reuse the 68 just-qualified
GUI/animation roots and run the 289 additional owner roots twice per profile.
The selection conservatively includes owner families and related state/serialization
cases; it is broader than the exact reachable scope, since deleted wrappers have
no live callers. No new artificial fixture calls to those wrappers are added.

After installation require all 357 roots in each profile, safe/static checks,
three production builds and retired-export ABI checks, the established full-suite
result, headless preview and final save/load/resume, and original asset hashes.
This deletion-only batch does not alter shared dispatch, layout or retained
function code; the prior full corpus can remain the latest broad qualification,
with its earlier-source limitation stated explicitly after conversion.

## Review and delegation

Luna was unavailable after its usage-limit response. Primary prepared and reviewed
this batch while the preceding animation corpus ran, using an ignored draft and
keeping active test source frozen. No replacement model or extra helper was used.

The audit includes complete Go identifier references and C preambles, including
one-line bodies. It separately follows four macro aliases and excludes 14 exports
with live fixture/allocator users. Imports are pruned only in touched files;
other production imports keep their packages initialized. All remaining C types,
compiler flags, allocator observers, operation-name strings and captures stay.

Local evidence and draft: `build/port-unused-abi-adapters/`. No conversion is
installed at this original baseline. This report will be updated after qualification.

Both additional original selections passed twice per profile with exact expected
names and no failures/skips. Recorded [baseline](unused-abi-adapters-baseline.json),
[manifest](unused-abi-adapters-batch.json), and
[test selection](unused-abi-adapters-tests.txt).

Nine obsolete project cache archives (420,343,808 allocated bytes) and four
superseded GUI production/safe binaries (187,723,776 bytes; rebuild `2bca8420`)
were removed after source/hash and host-use checks. Latest animation binaries,
all baseline evidence, source and original assets remain. Journals are under
`build/port-unused-abi-adapters/`.
