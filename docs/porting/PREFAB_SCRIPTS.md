# Prefab script merging and generation setup

## Scope and status

The preceding prefab runtime conversion **d5f460cc is committed and pushed**.
This connected candidate covers **20 C functions / 1,363 original body lines**:
all binary helpers in server__script__file.c; object callback setters and prefab
name adjustment; pending-object reference remapping; generation initialization,
selection copying and bounds ordering. [Selection](prefab-scripts-selection.json)
records original bodies and source locations. No selected algorithm is ported yet.

## Baseline strategy

Use actual FILE/binfile handles and temporary script files. Observe output bytes,
input/output positions, return values and remapping counters. Cover each instruction,
wide and signed operands, both merge inputs, reserved globals/functions, builtin
callback adjustments, section markers, complete merge and file rewrite. Parse
complete results through the existing script reader; matching captured C bytes is
not the only correctness criterion. Reuse real prefab/object/waypoint owners for
names, callbacks, pending references and generation state. Qualify the connected
corpus on default/server/highres and fresh production after any production repair.
Freeze and repeat captures before replacing algorithms; retain no C reference code.

## Independent contracts in progress

Two initial contracts test short distinct suffix names and instruction copying
with relocation disabled. The latter asks the actual asm decoder to accept each
fixture and then requires its 32-bit operand to survive copying unchanged. C uses
one-byte locals for several operands despite the serialized/VM word width. The
name helpers construct a suffix, then unconditionally overwrite it with
ERROR_NAME_TOO_LONG!; tests expect the suffix consumed by the script reader.
The first original-C run (`c-contract-before`) reproduced both defects: `Door`
with instance 1 became the error placeholder, and copying integer 128 produced
−128. All fixtures passed the existing asm decoder. No oracle is frozen.

The installed pre-baseline corrections retain 32-bit operand locals, format names
within the actual 256-byte shared buffer and use the existing placeholder only
when the result does not fit, and restore a 4,096-byte name buffer in the reserved
function merger with explicit read bounds. An actual generated nc.obj contains
six-byte GLOBAL names, confirming that the prior four-byte scalar is insufficient.
These are deliberate reversible correctness repairs under the standing policy.
Working C is **36,924 lines / 73 files** (+7 temporary lines before conversion).

The expanded four-root C run (`c-contract-after`) passes. It exercises every
opcode, signed/wide operands, relocation overflow and builtin callback rewrites,
plus complete merged scripts with varied function counts checked by asm.ReadScript.
Input/output bytes, cursor positions, returns and remapping counters are observed.
No changed behavior has been concealed by regenerating a frozen expectation.

Further audit concerns: the reserved-function merger reads a variable-length name
into a four-byte scalar (the actual reader documents GLOBAL names); generation
initialization later copies a potentially nil path; root merge uses an existing
output-file convention and creates wrapper handles around caller-owned files.
Confirm ownership, intended limits and real call contracts before changes. Record
any reversible correction with evidence rather than concealing it in new goldens.


## File lifecycle corrections

`c-files-before` reproduced three connected wrapper defects with ordinary valid
inputs: no output when the destination was absent, 4,096 output bytes instead of
146 because an old tail survived, and three file registrations left after return.
Open the destination with create/truncate semantics. The C bridge now removes only
handles it created, leaving caller-owned streams and preexisting handles intact.
`c-files-after` passes all seven roots, including partial scalar reads, file merges
with existing/missing destinations, empty-input paths and usable borrowed handles.
The instruction/complete-merge captures remain unchanged. These are reversible
pre-baseline repairs; fresh production qualification is required before freezing.

Callback assignment coverage is now compiling/running using actual pooled objects,
script storage, VM lookup and getter, with field canaries across editor, prefab and
runtime modes. Whole-buffer contracts check that only the selected callback field
changes; captures store resulting buffer digests and scalar indices.

Disk cleanup reclaimed **1,660,044,319 bytes (about 1.55 GiB)** from the three
completed prefab C-baseline scenario asset copies after original-file SHA-256
comparison. Changed files, saves, reports and binaries are preserved. Each run has
a restoration manifest. build/port-prefab-scripts/deduplicate-prefab-c-assets.py
--apply is consumed; never repeat its deletion mode. Original assets/archive and
latest native scenario copies remain available.


## Pending-reference loop repair

The first combined pending-reference/VM-suffix run (`c-vm-before`) stalled on a
missing target. A SIGQUIT stack from the isolated test process located the loop
inside the real `serverObjects.PendingByScriptID`. Its increment expression called
`it.Next()` without assigning the returned pointer, so any nonmatching first
record repeated forever. The fixture's links were valid. Assign the next object
back to the loop variable. This one-line Go correction is inside the selected C
routine's actual lookup service. Missing/duplicate target, wraparound and coordinate
contracts are rerunning (`c-pending-after`). The VM suffix test was not reached in
the stopped process; its parser correction is not installed or claimed yet.


## VM coordinate suffix correction

`c-pending-after` completed all thirteen roots: twelve passed, including all sixteen
pending-reference cases. Only the new VM coordinate contract failed. For example,
`OnEnter%7%-46%92` produced offset `(92,0)` instead of `(-46,92)`. The dependency
reader assigns both coordinates to X. The local `NoxScriptVM.ReadScript` consumer
now restores both signed axes from the suffix, without modifying dependencies.
This is a reversible prerequisite correction; no script captures are frozen yet.
The expanded object naming tests include populated even-numbered callbacks, signed
coordinates, and unchanged callbacks on a second pass. Complete file rewrites are
compared with independently serialized expected scripts and read through the VM.


The 14-root `c-rewrite` run passed, with 24 complete file-rewrite cases (including
missing inputs) and closed-file/backup-removal contracts. Static checks pass.
`c-generation` completed sixteen roots; fifteen passed. Its sole failure was the
fixture expecting a full-width clock: the real `nox_platform_get_ticks` adapter
returns `C.uint`, so `0x10000002a` becomes `42` before storage. The fixture now
asserts that actual ABI behavior, which the native port must preserve. Generation
cases cover missing/bad/valid libraries, metadata, deletion, selection, paint
settings, path construction, counters and closure. The shared pending lookup also
passes a real group-reference regression for first/later/duplicate/missing IDs.

The broader selection is 155 root tests: all 136 qualified prefab-runtime roots,
the sixteen new script contracts, and three connected VM/collision consumers.
Fresh production qualification is required because of the documented repairs.

## Qualified repaired C baseline

Sixteen focused roots pass in two independent processes. Twelve frozen captures
contain **4,182 records**. All **155 selected roots / 2,056 including subtests**
pass without skips on default/server/highres (package times 23.336s/23.627s/24.492s).
All **91 connected captures / 34,959 records** match across the three targets;
the source identity is unchanged. Static mapped-memory checks pass.

Fresh production qualification passes all three ELF32/i386/SSE2/CGO builds and
ABI/interface checks, the exact existing full-suite result (1,553 failure entries;
15 passing, 3 failing and 32 skipped packages), headless gameplay, save/load and
flat-map regeneration. See [qualification](prefab-scripts-c-qualification.json),
[captures](prefab-scripts-captures.json), [manifest](prefab-scripts-batch.json) and
[test selection](prefab-scripts-tests.txt). Local evidence is under
`build/port-prefab-scripts/c-{default,server,highres,production}`.

A final reachability audit found no production caller, callback or registration
for `nox_xxx_tileInitdataClear_4D3C50` or `sub_4D3C70`. The native step will remove
these two orphan bodies and their dedicated fixture instead of preserving code
solely for tests. Eighteen live algorithms remain to translate; all twenty C
interfaces can retire because their external callers are Go. The selected/instance
counters also lose their last C consumers and can move into existing Go state.
This is a reversible scope refinement, recorded for review.

Older disposable Go cache entries were removed only after all gates joined.
The first pass encountered a cache directory and stopped; the completed follow-up
records an additional **6,839,608,799 bytes** removed in `cache-cleanup.json`.
Total free disk is now about **13 GiB**. Source, original assets/archive and
qualification artifacts are preserved. Do not replay completed deletion passes.
