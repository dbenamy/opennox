//go:build porttest

package server

// PortTestOrchestrationRewardTypes adds a bow subclass to the real test factory.
func (s *Server) PortTestOrchestrationRewardTypes() func() {
	restore := s.PortTestRewardTypes([]string{"RewardMarker", "RewardMarkerPlus", "RewardChest", "RewardOther", "Quiver", "Ankh", "RedPotion", "QuestGoldPile", "QuestGoldChest", "RubyGem", "EmeraldGem", "DiamondGem", "PortTestRewardWeapon"}, nil, true, 0, 4)
	s.Types.ByID("PortTestRewardWeapon").subclass = 4
	return restore
}

// PortTestOrchestrationActivators seeds unresolved IDs in the actual timer list.
func (s *Server) PortTestOrchestrationActivators(ids [][2]uint32) (func() ([][2]uint32, []ActivatorArgs), func()) {
	old := s.Activators
	s.Activators = serverActivators{}
	for _, pair := range ids {
		s.Activators.append(&activator{triggerID: pair[0], callerID: pair[1], arg: ActivatorArgs{Callback: 7, Arg: 9}})
	}
	return func() ([][2]uint32, []ActivatorArgs) {
		var refs [][2]uint32
		var args []ActivatorArgs
		for p := s.Activators.head; p != nil; p = p.next {
			refs = append(refs, [2]uint32{p.triggerID, p.callerID})
			args = append(args, p.arg)
		}
		return refs, args
	}, func() { s.Activators = old }
}
