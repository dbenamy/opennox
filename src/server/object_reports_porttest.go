//go:build porttest

package server

// Own actual animation metadata consumed by the shared production range lookup.
func (s *Server) PortTestObjectReportAnimations(frames, delay int) func() {
	old := s.Types.playerAnimFrames
	s.Types.playerAnimFrames = make([][2]int, 64)
	for i := range s.Types.playerAnimFrames {
		s.Types.playerAnimFrames[i] = [2]int{frames, delay}
	}
	return func() { s.Types.playerAnimFrames = old }
}

// Own the actual audio allocator and catalog used by observer delivery tests.
func (s *Server) PortTestObjectReportAudio() func() {
	old := s.Audio
	s.Audio = serverAudio{}
	s.Audio.Init(s)
	s.Audio.bySound[1].Field12 = 1
	s.Audio.bySound[2].Field12 = 1
	return func() { s.Audio.Free(); s.Audio = old }
}
