package legacy

/*
#include "client__gui__window.h"
#include "client__gui__guicon.h"
*/
import "C"
import (
	"github.com/opennox/libs/console"
	"unicode/utf16"
)

var (
	GetConsole                     func() *console.Console
	Nox_gui_console_flagXxx_451410 func() int
	Nox_gui_console_Hide_4512B0    func() int
)

//export nox_gui_console_flagXxx_451410
func nox_gui_console_flagXxx_451410() int { return Nox_gui_console_flagXxx_451410() }

//export nox_gui_console_Hide_4512B0
func nox_gui_console_Hide_4512B0() int { return Nox_gui_console_Hide_4512B0() }

func Nox_xxx_consoleEsc_49B7A0() {
	nox_xxx_consoleEsc_49B7A0()
}

func textFormatConsole(color byte, format *uint16, args ...textFormatArgument) int {
	GetConsole().Print(console.Color(color), string(utf16.Decode(textFormatUnits(format, args...))))
	return 1
}
