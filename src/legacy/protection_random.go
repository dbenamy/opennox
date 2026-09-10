package legacy

import "github.com/opennox/opennox/v1/internal/protection"

// The floating RNG is private to protection startup/rekey. Its former C and
// blob fields have no external readers or writers.
var protectionRandom protection.Random
