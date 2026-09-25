# Native drawable callbacks and parsers

Selected scope: 66 C exports from the effect, object, particle and sprite bridge
files: 56 draw callbacks, seven parsers and three fixture-only light helpers.
Keep the existing Go algorithms, 386 layouts and external SDL/OpenGL backends.

## Boundary and decisions

Use distinct static byte identities for the 56 drawable callbacks, including
separate identities when two wrappers share an owner. A client-owned typed
registry avoids an import cycle with legacy. Preserve signed int32 results and
separate result/discard dispatch: normal drawable drawing observes an integer;
shield drawing uses the void convention. Unknown C callbacks retain their
existing raw call convention. Both `Drawable.CallDraw` and the legacy draw wrapper
must use this boundary, as must shield drawing and its direct fixture call.

Migrate every selected address registration, parser assignment, identity getter,
comparison and fixture getter together. The draw registry and object/drawable
records retain pointer-sized fields; no serialization or ownership layout changes.
Unselected player/monster/screen callbacks continue through the C fallback.

The seven parsers need no callback identities. Their production adapter currently
copies the attribute text into the supplied buffer, calls the parser with a
256-byte view, ignores its boolean result and returns nil. Preserve that behavior;
fixtures independently observe the boolean. Preserve sparse fixture operation
numbers, signed results, pointer words, float32 inputs and 64-bit light-angle
results. The private C-scalar light-intensity fixture wrapper can call its existing
Go owner directly after its sole caller is verified.

Keep 25 update exports outside this batch: all 23 update-file exports, the orbit
update in the effect file and the monster-generator update in the object file.
The named UPDATE table feeds the root int-returning update loop. One entry,
`nox_xxx_updDrawMagic_4CDD80`, has a void signature, so its observable effect on
secondary-callback gating needs original-path investigation before migration.
Its fixture currently discards the result and cannot establish that behavior.
This is distinct from an int callback invoked as void, whose result can simply
be discarded. The drawing/parser batch has a coherent boundary without changing
that update behavior.

## Original-path baseline

The conservative caller selection starts with all APIs in the five affected
fixture files, eight production identity getters and both draw-dispatch entry
names. It follows calls through root test helpers, including shared UI owners and
snapshot helpers. Generic method-name matches can overselect; this is a coverage
aid, not a complete typed or dynamic call graph. Primary independently matched
all 522 symbol occurrences from the 91-export scouting set, then selected the
66-export draw group. The current scan covers 3,064 existing tracked Go/C/header/
assembly files; earlier helper prose overstated that file count.

The selected baseline contains 1,035 default/highres roots and 1,025 server roots;
ten tests with explicit `porttest && !server` constraints are excluded on server.
Reuse the just-qualified audio full default corpus for all default roots, and 64
previous focused roots in each other profile. Run the remaining 961 server and
971 highres roots with the exact existing source/binary fingerprints and runtime
environment. No new test or production source is needed before this baseline.
Baseline accepted: all expected roots passed without skips or failures, with exact
name-set and source/environment checks. No conversion is installed.

After conversion require the complete focused selection in each profile, safe/
static checks, a fresh default-client scenario before the wider production sweep,
three production ABI checks, exact known asset-suite results and full default
corpus comparison. Keep existing assertions and frozen captures unchanged.

## Delegation and recovery

GPT-6 Luna is drafting the legacy registration/getter/parser/fixture migration in
an ignored overlay. Primary owns the client registry and raw dispatch consumers,
selection/baseline acceptance, integration and qualification. The initial scouting
missed the concrete UPDATE table route and emphasized the wrong return-convention
case; primary traced the table and the helper supplied a verified addendum.

Local selection edges, original functions, overlays and commands are under
`build/port-client-draw-identities/`. Scouting references and decisions are under
`build/port-after-audio/`. Drafts are not accepted source or qualification evidence.
