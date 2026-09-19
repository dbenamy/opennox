//go:build porttest

package legacy

func PortTestServerBrowserLabel(key string, list bool) {
	if list {
		browserShowList()
	} else {
		browserLabel(key)
	}
}
