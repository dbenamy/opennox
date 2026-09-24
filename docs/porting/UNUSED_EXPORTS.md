# Unused internal export bridges

Status: qualified. Original baseline at `92029ddd`, frozen in `5fa49336`;
the reviewed removal passes the completed-batch gates. This is a reachability-based cleanup,
not an algorithm replacement.

Retire 65 unused C export bridges: twelve complete forwarding files, the player
departure adapter beside `scriptLog`, and six class-allocation adapters beside
`DeadWord`. Keep that live logger and sentinel unchanged. Delete exactly 63 simple
header prototypes; preserve all other header types, macros and declarations.
The existing game implementations, registries and callback identities are untouched.

## Reachability and review

Primary independently scanned every tracked source/header/C/assembly line,
including cgo preambles, for all 65 names. Outside their export/definition and
prototype lines, four matches required review: two historical comments, an
unrelated private root-package Go function with the same name as a geometry
export, and a simple prototype with a trailing `idb` comment. None is a caller of
the retired C bridge. Preserve the unrelated Go function. Evidence and exact
declarations are in [unused-exports-reachability.json](unused-exports-reachability.json).

The repository builds `./cmd/opennox` executables for default/highres/server;
no supported external shared-library API or dynamic symbol lookup was found in
the source/build surface. Unknown external consumers cannot be disproved by a
repository scan. This is a reversible internal API retirement within the agreed
engine-glue milestone. Fresh ABI qualification explicitly requires all 65 C
symbols to be absent and preserves the remaining required symbol set.

All twelve deleted Go files contain only these exported functions, imports and
preambles, with no other Go declarations or C build directives. Primary verified
the two retained files preserve their complete live declarations, and each header
draft is byte-identical to its original after removing only the listed prototypes.

## Original qualification

All 164 selected owner roots pass without skips in default, server and highres.
The selection covers entry/listbox UI, respawn, polygons/population, match roster,
monster control, quest progress, targeting, world geometry and lightning owners.
It deliberately excludes the standalone map-population diagnostic probe.
Allocator and server package checks pass in all profiles; `internal/binfile` has
no direct test files, so its check establishes package compilation only. Fresh
headless map/save/load supplies integration coverage after removal. No fixture is
added solely to preserve a dead wrapper.

Original source is identical to the freshly qualified string-boundary chunk, so
that production evidence is reused only for the baseline. After removal, repeat
owner/package checks, safe build/static, fresh production builds/ABI, exact known
suite comparison and headless save/load. Root-name sets and source fingerprints
must match the recorded scope. Existing test files and expectations stay unchanged.

The owner audit added 103 qualified existing roots to the accumulated selector.
A systematic audit of remaining selector gaps is underway; missing names do not
imply their earlier focused qualifications were never run.

## Delegation

Luna supplied the export inventory and deletion draft. Primary rejected the first
inventory's omission of C preambles and caught a live logger in a proposed whole
file removal. The corrected inventory includes preambles and full header contexts;
primary separately rechecked this entire 65-symbol cohort. Luna then produced
bounded deletion manifests and exact prototype removals. Primary retained scope,
baseline, integration and qualification ownership. No measured cost savings are
claimed; the broad inventory needed more review than the bounded draft.

Artifacts: `build/port-unused-exports/`; original evidence:
[unused-exports-c-qualification.json](unused-exports-c-qualification.json).

## Completed qualification

All 164 original owner roots pass in each production profile without skips, with
identical selected-name sets. Package checks, safe build/static checks, three fresh
production builds/ABI checks and headless character creation/save/load/resume pass.
The package suite matches all 304 known failure events and package outcomes exactly.
All phases share identical source fingerprints; 1,733 test/fixture files are unchanged.
Safe owner contracts were not part of this batch's scope.

Selected project cgo files fall from 413 to 399 in each profile: 64/463 removed on
net since this phase began. Three project packages still use cgo. Selected legacy
exports fall from 1,887 to 1,822; 157 headers now contain 4,479 physical lines.
Embedded C callback bodies remain 79; standalone production/test C remains zero.
External native dependency selections are unchanged.

The systematic selector audit found 899 omitted default/highres roots and 892
server roots. Repair and full-corpus qualification follow this chunk before further
source conversion. Earlier focused qualifications remain valid for their recorded
scope. See [unused-exports-qualification.json](unused-exports-qualification.json).
