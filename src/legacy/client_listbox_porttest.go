//go:build porttest

package legacy

import "github.com/opennox/opennox/v1/client/gui"

func PortTestListboxScrollIndex(d *gui.ScrollListBoxData) int { return uiListIndex(d) }
