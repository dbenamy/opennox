//go:build porttest && 386

// Package fpenv provides processor-state operations for qualified Linux 386
// port fixtures. Callers must pin their OS thread while changing its state.
package fpenv

// Environment contains a 28-byte x87 environment followed by MXCSR.
type Environment [8]uint32

func Control() uint16
func SetControl(value uint16)
func MXCSR() uint32
func SetMXCSR(value uint32)

//go:noescape
func save(out *Environment)

//go:noescape
func load(in *Environment)

func Read() (out Environment) { save(&out); return }

//go:nosplit
func Begin() Environment {
	old := Read()
	SetControl((uint16(old[0]) &^ 0x0f00) | 0x0200)
	SetMXCSR(old[7] &^ 0x6000)
	return old
}

// Restore matches the qualified Linux i386 libc fesetenv behavior. It restores
// exception/control fields while retaining the current stack tags and TOP.
//
//go:nosplit
func Restore(old Environment) {
	current := Read()
	current[0] = (current[0] &^ 0x0f3f) | (old[0] & 0x0f3f)
	current[1] = (current[1] &^ 0x003f) | (old[1] & 0x003f)
	current[3] = 0
	current[4] &= 0xf8000000
	current[5] = 0
	current[6] &= 0xffff0000
	current[7] = old[7]
	load(&current)
}
