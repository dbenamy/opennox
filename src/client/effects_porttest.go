//go:build porttest

package client

import (
	"log/slog"
	"strings"
	"sync/atomic"

	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// PortTestEffectsClient owns real drawable allocation, indexing and type lookup
// without registering unrelated GUI, network or server callbacks.
func PortTestEffectsClient(s *server.Server, names []string) (*Client, func()) {
	c := &Client{Log: slog.Default(), Server: s, r: noxrender.NewRender(slog.Default(), s)}
	c.handle = atomic.AddUintptr(&clientLast, 1)
	clients.Store(c.handle, c)
	c.Objs.init(c)
	c.Objs.Init(128)
	c.Things.init(s.Strings())
	c.Things.byInd = []*ObjectType{nil}
	c.Things.byID = make(map[string]*ObjectType)
	vp, freeVP := alloc.New(noxrender.Viewport{})
	c.vp = vp
	var frees []func()
	for i, name := range names {
		typ, free := alloc.New(ObjectType{})
		text, freeText := alloc.CString(name)
		typ.Name, typ.Field_1c = text, int32(i+1)
		c.Things.byInd = append(c.Things.byInd, typ)
		c.Things.byID[strings.ToLower(name)] = typ
		frees = append(frees, free, freeText)
	}
	return c, func() {
		c.Objs.Free()
		c.r.Part.Free()
		// No asset bag or fonts were opened. Bag.Free requires an open file,
		// and Fonts.Free would close shared fonts this owner did not load.
		freeVP()
		for i := len(frees) - 1; i >= 0; i-- {
			frees[i]()
		}
		clients.Delete(c.handle)
	}
}
