//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func waypointGraphCorpus() []legacy.PortTestWaypointGraphSpec {
	var out []legacy.PortTestWaypointGraphSpec
	for seed := uint32(1); seed <= 2048; seed++ {
		n := int(seed%30) + 2
		sp := legacy.PortTestWaypointGraphSpec{Nodes: make([]legacy.PortTestWaypointGraphNode, n), Start: 0, End: n - 1, Capacity: int(seed % 33), Epoch: []uint32{0, 1, 0xfffffffe, 0xffffffff}[seed%4], OneShot: 17}
		rng := seed
		for i := range sp.Nodes {
			nd := legacy.PortTestWaypointGraphNode{Flags: 1, Flags2: 128, Parent: -1, Next: -1}
			if seed%5 == 0 && i%3 == 0 {
				nd.Epoch = sp.Epoch + 1
			}
			if seed%7 == 0 && i%4 == 0 {
				nd.Flags = 0
			}
			if seed%11 == 0 && i%5 == 0 {
				nd.Flags2 = 1
			}
			for j := 0; j < int(seed%5); j++ {
				rng = rng*1664525 + 1013904223
				nd.Edges = append(nd.Edges, int(rng%uint32(n)))
			}
			if i+1 < n && seed%3 != 0 {
				nd.Edges = append(nd.Edges, i+1)
			}
			sp.Nodes[i] = nd
		}
		if seed%13 == 0 {
			sp.Start = -1
		}
		if seed%17 == 0 {
			sp.End = -1
		}
		if seed%19 == 0 {
			sp.End = sp.Start
		}
		out = append(out, sp)
	}
	for _, n := range []int{1, 2, 15, 16, 17, 18, 32, 33, 255, 256, 257, 258} {
		for _, cap := range []int{0, 1, 2, 15, 16, 17, 32} {
			sp := legacy.PortTestWaypointGraphSpec{Nodes: make([]legacy.PortTestWaypointGraphNode, n), Start: 0, End: n - 1, Capacity: cap, Epoch: 1, OneShot: 17}
			for i := range sp.Nodes {
				sp.Nodes[i] = legacy.PortTestWaypointGraphNode{Flags: 1, Flags2: 128, Parent: -1, Next: -1}
				if i+1 < n {
					sp.Nodes[i].Edges = []int{i + 1}
				}
			}
			out = append(out, sp)
		}
	}
	return out
}

// Independent layer traversal preserves reversed discovery order and the C
// epoch collision rule; it models the returned route without pointer storage.
func waypointGraphRoute(sp legacy.PortTestWaypointGraphSpec) ([]int32, bool) {
	eligible := func(i int) bool { return i >= 0 && sp.Nodes[i].Flags&1 != 0 && sp.Nodes[i].Flags2&128 != 0 }
	if !eligible(sp.Start) || !eligible(sp.End) {
		return nil, false
	}
	epoch := sp.Epoch + 1
	seen := make([]uint32, len(sp.Nodes))
	parent := make([]int, len(sp.Nodes))
	for i, n := range sp.Nodes {
		seen[i] = n.Epoch
		parent[i] = -1
	}
	seen[sp.Start] = epoch
	layer := []int{sp.Start}
	for len(layer) != 0 {
		var next []int
		for _, v := range layer {
			if v == sp.End {
				var backwards []int32
				for v >= 0 && len(backwards) < 256 {
					backwards = append(backwards, int32(v+1))
					v = parent[v]
				}
				route := make([]int32, len(backwards))
				for i, v := range backwards {
					route[len(route)-1-i] = v
				}
				return route, true
			}
			for _, to := range sp.Nodes[v].Edges {
				if seen[to] != epoch && eligible(to) {
					seen[to] = epoch
					parent[to] = v
					next = append([]int{to}, next...)
				}
			}
		}
		layer = next
	}
	return nil, false
}
func TestAIPathWaypointGraphBaseline(t *testing.T) {
	specs := waypointGraphCorpus()
	got, restored := legacy.PortTestWaypointGraph(specs)
	if !restored {
		t.Fatal("graph globals not restored")
	}
	for i, r := range got {
		if !r.GuardsOK || !r.OnlyTraversalState {
			t.Fatalf("graph %d guard/read-only mutation", i)
		}
		route, found := waypointGraphRoute(specs[i])
		wantReturn, wantFlag := 0, uint32(2)
		if found {
			wantReturn = min(len(route), specs[i].Capacity)
			wantFlag = 0
			if wantReturn != len(route) {
				wantFlag = 1
			}
			for j := 0; j < min(len(route), specs[i].Capacity+1); j++ {
				if r.OutputWords[j] != route[j] {
					t.Fatalf("graph %d route[%d]=%d want=%d", i, j, r.OutputWords[j], route[j])
				}
			}
		}
		if r.Return != wantReturn || r.OneShot != wantFlag {
			t.Fatalf("graph %d count/flag=%d/%d want=%d/%d", i, r.Return, r.OneShot, wantReturn, wantFlag)
		}
	}
	b, _ := json.Marshal(got)
	hash := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("cases=%d complete-state-sha256=%s", len(got), hash)
	const baseline = "9aaae7e2192fbc77655d061d7399e5f24dc2e45bccbd348db98ec2d54a09f3d2"
	if baseline != "" && hash != baseline {
		t.Fatal("original C graph state differs", hash)
	}
}
