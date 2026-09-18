//go:build porttest

package netstr

import "net/netip"

// PortTestStatisticsPeers owns real connection records without opening sockets.
func (s *Streams) PortTestStatisticsPeers() func() {
	old := s.streams
	for i := range s.streams {
		s.streams[i] = &Conn{g: s, ind: i, id: i, addr: netip.AddrPortFrom(netip.AddrFrom4([4]byte{192, 0, 2, byte(i + 1)}), uint16(10000+i))}
	}
	return func() { s.streams = old }
}
