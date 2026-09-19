//go:build porttest

package client

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
)

// PortTestBookGuideTypes owns real client type lookup entries used by the guide
// loader. It does not replace lookup or loader functions.
func (c *Client) PortTestBookGuideTypes(names []string, missing int) func() {
	old := c.Things
	c.Things = clientObjTypes{byID: make(map[string]*ObjectType)}
	var frees []func()
	for i, name := range names {
		if i == 0 || i == missing {
			continue
		}
		typ, free := alloc.New(ObjectType{})
		frees = append(frees, free)
		canonical := name
		if i == 7 {
			canonical = "Bomber"
		}
		ptr, release := alloc.CString(canonical)
		frees = append(frees, release)
		typ.Name = (*byte)(ptr)
		typ.Field_1c = int32(1000 + i)
		typ.ObjSubClass = object.SubClass(i % 4)
		c.Things.byID[strings.ToLower(name)] = typ
	}
	return func() {
		c.Things = old
		for i := len(frees) - 1; i >= 0; i-- {
			frees[i]()
		}
	}
}
