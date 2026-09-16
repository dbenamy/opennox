# Floor and edge asset readers

Status: **qualified frozen C baseline; no Go conversion applied**.
Previous rendering milestone: **642fba50**, committed and pushed.
Current C size: **73,453 physical lines /91 files /zero reference C**, including
one three-line prerequisite correction. The connected conversion scope is nine
readers/helpers, 513 corrected C block lines (510 before the prerequisite).

## Ownership and compatibility

The six readers load definitions, skip records in another loading pass and bind
images to existing definitions. Root `things.go` calls their `legacy/things.go`
wrappers. Two helpers release the image arrays; another matches facade names.
The actual owners are MemFile, 176 tile definitions, 64 edge definitions, shared
counts, the mapped startup facade table and the client's image bag.

Floor definition/image readers leave the following END word; the floor skip
reader consumes it. Edge readers consume END unless the format byte is exactly
1, which returns early. Preserve partial definition writes on rejection and the
count behavior on a retry. Other format-byte values follow the normal path.
Names consume their declared byte length, while comparisons and copies retain C
string termination. Floor names clear 31 bytes; edge names write through the NUL
and leave the remaining bytes alone.

Keep raw calloc/free compatibility for image arrays. The edge binder allocates
five bytes per image but stores contiguous four-byte handles, leaving a zero
tail. Preserve that extent during this translation. Free helpers clear pointers
within the active count and leave other metadata/counts alone. Do not send these
untracked allocations to `alloc.FreePtr`.

Inline image references use the actual external-image resolver, which currently
returns nil; indexed images use the real image bag. No substitute resolver or C
algorithm is retained for testing. Normalize image handles to owned indices,
never compare allocation addresses.

## Prerequisite correction for later review

The original floor image binder used a pointer-derived index when no definitions
were loaded. It now returns the existing missing-definition failure after reading
the name, before accessing a definition. The edge binder already fails in this
case. This adds three C lines and has an independent cursor/no-side-effects
contract. Normal loading order supplies definitions first. This is a narrow,
reversible correction under the standing authorization; it does not attempt
broader malformed-file handling. Qualify fresh production binaries and gameplay
for the corrected C baseline before translation.

## Contracts and qualification

Direct actual-owner matrices cover capacities, names, color conversion/sentinel,
field products, zero/max byte dimensions, indexed/inline/mixed image references,
format/END failures, first matching name, exact cursor and scratch effects, whole
metadata state, repeat/count-limited free and failed definition retry.

The first C development run completed in 186.753s and failed because the fixture
had not initialized the actual startup facade-name table. Binding, first-match,
failure and free contracts passed; the facade test then panicked, preventing the
skip test from running. The correction loads the original embedded table bytes,
relocates its six real name pointers and restores all touched bytes on cleanup.
It changes fixture setup only. Additional review cases check embedded NUL names
and format 255. The corrected run passed seven roots in 187.891s; 4,396 records
in seven groups are frozen. Default/server/highres affected tests, production interfaces, full-suite comparison
and both gameplay replays pass.

Affected selection extends the rendering/definition/map owner corpus with these
readers. Default/server/highres each require discovery and execution of every
selected root, repeated exact C captures, production builds/ABI, exact known
full-suite failures and fresh normal 12 plus GUI flat-floor 14 frame comparisons.
The full accumulated rendering corpus was qualified immediately before this batch;
use affected qualification here unless new evidence widens the scope.

Out of scope: truncated buffers, definition names beyond the original fixed
buffer contract and allocation exhaustion. Solo replay does not establish remote
multiplayer coverage. Original assets/archive remain unchanged.

## Local recovery

Ignored evidence/drafts: `build/port-floor-assets`. `prepare-c.py` and
`review-tests.py` have already been applied and must not be rerun. The native
implementation remains an unexecuted `.stage` draft. Join the entire active test
driver before editing Go/C/header source. Qualify/freeze/commit/push C before
applying the native draft.

## Frozen C checkpoint

| Group | Records | SHA-256 |
| --- | ---: | --- |
| binding | 630 | `b1e4d0288b1fb06bcf9bb28aeb20ba9ee02e4f34b48c277b160c3874c11a76be` |
| definitions | 3360 | `e72c5717e6ed60608fd938d346df873b456a998dcbabd1a39fa892e34500690a` |
| facades | 22 | `c46dc7c97e5d008e43ea19709c8e79cbc9a72547e421b1e5a02427388ad7ad93` |
| failures | 6 | `a3c38136006b82d590e109a0bfc9fe493788f5910042176934f0214cdc6a8528` |
| first-match | 2 | `714b6b32de6da860eb869b5088f1f24f0b10875c6bc8d54e30286c813f617a72` |
| free | 16 | `c1b5d4dec07d07c54b189294a39bcb58fe59abe7504d60ebeb190dddd4650d41` |
| skipping | 360 | `0ed871e6198beaeedef4ba84a056545cb216308a3c3af0d15dd7d3bf43ddebc7` |

Affected roots **298 / 296 / 298** all completed in **159.933 / 237.897 / 160.450s**
(default/server/highres). All 4,396 results repeat exactly in separate processes
and all three targets. Three fresh production binaries pass ELF32/i386/SSE2/CGO
and all eleven original C interfaces are present, with no test helpers. The full
asset suite retains exactly 1,553 failure entries and 15 pass/3 fail/32 skip
packages. Fresh normal 12 and GUI flat-floor14 frames match the preceding qualified
composition replay. Source fingerprints are unchanged throughout qualification.
Evidence: build/port-floor-assets/c-qualification.json and c-evidence/.
