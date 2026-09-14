//go:build porttest

package opennox

// Complete corrected-C captures, repeated in default/server/highres.
var hallwayCExpected = map[string]string{
	"fallback":     "35290c699083c7ec2639a45e1f54d0e08fd466e021d886cf41b9b562e0e7e206",
	"obstructions": "96a741b29d6e102b5f975c8a7f1b28bf4afccece37b8d0d0f12d6695192b3d11",
	"routes":       "712a46430f47c5d55a5f1bcbe9ef4dfc638fae8e205ddab827b6266c7a1bfdbc",
}
