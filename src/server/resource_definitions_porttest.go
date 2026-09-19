//go:build porttest

package server

// PortTestResourceRegisteredParser resolves the same registration used by object loading.
func PortTestResourceRegisteredParser(kind, name string, typ *ObjectType, args []string) error {
	var f ObjectParseFunc
	switch kind {
	case "use":
		f = useParseFuncs[name]
	case "update":
		f = updateParseFuncs[name]
	case "death":
		f = deathParseFuncs[name]
	case "collide":
		f = collideParseFuncs[name]
	}
	if f == nil {
		panic("missing resource parser: " + kind + "/" + name)
	}
	return f(typ, args)
}
