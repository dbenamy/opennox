//go:build porttest

package handles

// PortTestInit gives a sequential fixture an isolated handle arena, preserving
// any caller-owned arena and allocation cursor.
func PortTestInit() func() {
	oldData, oldBase, oldEnd, oldCur, oldFree := data, base, end, cur, free
	cur = 0
	Init()
	return func() {
		Release()
		data, base, end, cur, free = oldData, oldBase, oldEnd, oldCur, oldFree
	}
}
