# Fixture constants and local scalar types

## Scope and review

Seven fixture files replace local C scalar types and three constant/layout-only
imports with native equivalents. Rules, subtile lookup, record removal and float
record fixtures retain explicit-width casts and argument order. All six private
scalar wrappers are confined to their own files. The float fixture's control-word
observer remains C. Allocator adapters, algorithms, captures and assertions in root tests remain
unchanged. Six more files lose their C import.

The inventory guard now checks the actual Go owner: a 148-byte cell and 84-element
backing array. Console red uses the existing console.ColorRed value 6. The session
protocol fixture retains 0x000F039A: video_highres.go has no build constraint and
sets NOX_HIGH_RES for default, server and highres alike. Do not infer a different
protocol value from the profile names. Target compiler probes independently confirm
all constants, C layouts and scalar widths in each profile. Fifteen headers in the
removed include closure have no startup hooks. AST review: 14 changed functions,
eight unchanged, none added or removed.

## Qualification plan

Run original owner selections freshly using only the preceding batch's corrected,
verified binaries: 369 default, 366 server and 369 highres. Server omits two hover
tests and the inventory world-selection test by build tags. Repeat 14 focused roots
per profile and six safe contracts; commit accepted original evidence first.

After conversion, start with safe contracts, then focused preflight and exact
owner selections in all profiles, focused repeats, static checks, safe build,
production builds/ABI, exact known-suite outcomes and fresh save/load. All target
test binaries compile every root; runtime scope follows affected fixture callers.

## Status

Original baseline accepted: 369/366/369 owners, 14 focused repeats per profile
and six safe contracts passed without skips. Draft not installed. Expected fixture C imports: 20→14.
Production remains 4 client/highres and 5 server, with zero embedded production C
bodies, legacy exports and standalone C. Headers remain 157 /2,731 physical lines.
Luna quota remains unavailable; primary handles review and qualification.
Artifacts: `build/port-fixture-constants/`.
