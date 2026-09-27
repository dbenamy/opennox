# Book/UI callback original baseline

Original production revision:1238c985 (unused ABI adapters). Production reuse is
verified by exact source comparison allowing only the two new porttest files.
Baseline accepted; production conversion and its qualification remain pending.

Four new independent contracts exercise the existing image-completion call path:
last-frame entry at a nonzero start sequence and two-tick frame interval; repeated
completion; exact ImageRef argument; mutation of the selected frame cell versus
replacement of ImagesPtr inside the callback; loop/nil/unsupported paths; both
real book completion owners; 16-byte target layout. A Go-backed foreign ABI
observer records effects without retaining a C animation algorithm.

The first book-owner assertion used inert backing-blob words instead of extracted
Go globals. Corrected the fixture to PortTestBookWords before accepted preflight.
The corrected preflight passes all four roots. Existing frozen captures are untouched.

Audited combined selection:256 roots in default/highres,254 server. Two existing
!server rendering tests are excluded by build tags, not skipped at runtime.
All selected roots passed twice in independent processes. The initial selection
and seven audit additions are separate runs of the same source/profile binaries;
exact names, logs and hashes are captured in book-callbacks-baseline.json.

Primary review owns this batch because Luna is quota-limited. The intended change
covers seven actual callback identities, one unused export, five fixture-only
C-typed wrappers, and corresponding callers/prototypes. Expected production
reduction:seven cgo files/eight exports; these are provisional until qualification.
