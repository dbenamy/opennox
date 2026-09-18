//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"reflect"
	"testing"
	"unsafe"
)

func statisticsRandomOwner(t *testing.T) (map[string]*uint32, []byte) {
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	state := serverConfigOwnBytes(t, 0x5D4594, 741384, 276)
	serverConfigOwnBytes(t, 0x581450, 8368, 8)
	*memmap.PtrFloat64(0x581450, 8368) = 1e-9
	clear(state)
	return words, state
}
func TestGameStatisticsRandomState(t *testing.T) {
	words, state := statisticsRandomOwner(t)
	type row struct {
		A, B   uint32
		X, Y   int32
		Scale  uint64
		Result []float64
		Seed   int32
		State  []byte
	}
	var rows []row
	for _, a := range []uint32{0, 1, 30, 54, 55} {
		for _, b := range []uint32{0, 1, 30, 54, 55} {
			if a == b {
				continue
			}
			for _, pair := range [][2]int32{{0, 0}, {1, 0}, {0, 1}, {999999999, 0}, {0, 999999999}, {123456789, 987654321}} {
				for _, scale := range []float64{1e-9, 0.5e-9} {
					clear(state)
					*memmap.PtrUint32(0x5D4594, 741656) = 1
					*words["random-a"], *words["random-b"] = a, b
					nextA, nextB := a+1, b+1
					if nextA == 56 {
						nextA = 1
					}
					if nextB == 56 {
						nextB = 1
					}
					*memmap.PtrInt32(0x5D4594, 741384+uintptr(4*nextA)) = pair[0]
					*memmap.PtrInt32(0x5D4594, 741384+uintptr(4*nextB)) = pair[1]
					*memmap.PtrFloat64(0x581450, 8368) = scale
					before := bytes.Clone(state)
					got, seed := legacy.PortTestStatisticsRandom(1, 1)
					want := pair[0] - pair[1]
					if nextA == nextB {
						want = 0
					}
					if want < 0 {
						want += 1000000000
					}
					if len(got) != 1 || got[0] != float64(want)*scale || seed != 1 || *words["random-a"] != nextA || *words["random-b"] != nextB {
						t.Fatal("random state step", a, b, pair, got, want)
					}
					*(*int32)(unsafe.Pointer(&before[4*nextA])) = want
					if !bytes.Equal(before, state) {
						t.Fatal("random changed unrelated state")
					}
					rows = append(rows, row{a, b, pair[0], pair[1], math.Float64bits(scale), got, seed, bytes.Clone(state)})
				}
			}
		}
	}
	spellbookCapture(t, "game-statistics-random-state", rows, "7d67170ddecd66a3ebf2d084c8792c72a6362aa3d925a46fcaf3c9c67ad218f3")
}
func TestGameStatisticsRandomSequences(t *testing.T) {
	words, state := statisticsRandomOwner(t)
	type row struct {
		Seed   int32
		Count  int
		Floats []float64
		Bytes  []byte
		State  []byte
		A, B   uint32
	}
	var rows []row
	for _, seed := range []int32{0, 1, 2, 23, 161803398, 999999999, 2147483647, -2147483648} {
		for _, count := range []int{0, 1, 2, 54, 55, 56, 110, 257} {
			reset := func() { clear(state); *words["random-a"], *words["random-b"] = 0, 0 }
			reset()
			input := seed
			if input > 0 {
				input = -input
			}
			got, end := legacy.PortTestStatisticsRandom(input, count+1)
			expectedState := bytes.Clone(state)
			a, b := *words["random-a"], *words["random-b"]
			if end != 1 {
				t.Fatal("seed not initialized", end)
			}
			reset()
			data := legacy.PortTestStatisticsRandomBytes(seed, count)
			if !bytes.Equal(expectedState, state) || *words["random-a"] != a || *words["random-b"] != b {
				t.Fatal("byte helper consumes wrong number of random values")
			}
			for i, v := range data {
				n := int(got[i] * 255)
				if n < 0 {
					n = -n
				}
				if v != byte(n) {
					t.Fatal("byte scaling", i, v, n)
				}
			}
			reset()
			again, _ := legacy.PortTestStatisticsRandom(input, count+1)
			if !reflect.DeepEqual(got, again) {
				t.Fatal("seed reproducibility")
			}
			if count >= 54 {
				unique := map[float64]bool{}
				for _, v := range got {
					unique[v] = true
				}
				if len(unique) < 20 {
					t.Fatal("degenerate random fixture")
				}
			}
			rows = append(rows, row{seed, count, got, data, expectedState, a, b})
		}
	}
	spellbookCapture(t, "game-statistics-random-sequences", rows, "81aa19a062e9affe225e6e491b4dde1cc4a9857268dba12e6289fb9837bcc12c")
}
