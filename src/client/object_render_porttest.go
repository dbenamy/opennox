//go:build porttest

package client

import "image"

// Supply actual sight-polygon vertices; intersection and visibility algorithms
// remain the production implementation.
func (c *Client) PortTestObjectRenderSight(points []image.Point) (func() []image.Point, func()) {
	old := c.Sight
	if len(points) > len(c.Sight.sightPointsArr) {
		panic("sight fixture exceeds capacity")
	}
	c.Sight.sightPointsCnt = len(points)
	clear(c.Sight.sightPointsArr[:])
	copy(c.Sight.sightPointsArr[:], points)
	c.Sight.dword_5d4594_1217452 = 0
	clear(c.Sight.arr_5d4594_1212068[:])
	return func() []image.Point {
		n := c.Sight.dword_5d4594_1217452
		if n < 0 || n > len(c.Sight.arr_5d4594_1212068) {
			panic("invalid sight intersection count")
		}
		return append([]image.Point(nil), c.Sight.arr_5d4594_1212068[:n]...)
	}, func() { c.Sight = old }
}
