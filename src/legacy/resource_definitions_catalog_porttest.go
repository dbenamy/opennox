//go:build porttest

package legacy

/*
int sub_4F0640();
*/
import "C"

func PortTestResourceLinkCatalogs() int { return int(C.sub_4F0640()) }
