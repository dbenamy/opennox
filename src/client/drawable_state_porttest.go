//go:build porttest

package client

import "strings"

// Temporarily remove an owned fixture type from the actual lookup table.
func (c *Client) PortTestDrawableHideType(name string) func() {
	key := strings.ToLower(name)
	old, ok := c.Things.byID[key]
	delete(c.Things.byID, key)
	return func() {
		if ok {
			c.Things.byID[key] = old
		} else {
			delete(c.Things.byID, key)
		}
	}
}
