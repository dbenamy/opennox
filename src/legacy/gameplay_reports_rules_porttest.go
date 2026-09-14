//go:build porttest

package legacy

import (
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestGameplayReportRules struct {
	StartFrame, Notified, Limit uint32
	TeamMode, CountdownActive   bool
	Teams                       int
}
type PortTestGameplayReportCountdown struct {
	Seconds int
	Text    string
}
type PortTestGameplayReportRulesResult struct {
	Notified     uint32
	PlayerStatus [3]uint32
	Countdown    []PortTestGameplayReportCountdown
}
type portTestGameplayReportRules struct {
	membersReady bool
	spec         *PortTestGameplayReportRules
	countdown    []PortTestGameplayReportCountdown
}

func (p *portTestShopPools) gameplayReportRulesPrepare(sp *PortTestGameplayReportRules) func() {
	p.reportRules = nil
	if sp == nil {
		return func() {}
	}
	// These text-resource slots are empty in the headless fixture. Seed the ID
	// used by the existing Go countdown handler in server.go, then restore them.
	textOffsets := []uintptr{198872, 198928}
	textBefore := make([][]byte, len(textOffsets))
	for i, off := range textOffsets {
		b := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), 56)
		textBefore[i] = append([]byte(nil), b...)
		clear(b)
		copy(b, "Settings.c:SuddenDeathImminent")
	}
	st := &portTestGameplayReportRules{spec: sp}
	p.reportRules = st
	offsets := []uintptr{3520, 3536, 3476}
	values := []uint32{sp.StartFrame, sp.Notified, sp.Limit}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		g := memmap.PtrUint32(0x5D4594, off)
		old[i] = *g
		*g = values[i]
	}
	flags := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	if sp.TeamMode {
		noxflags.SetGamePlay(noxflags.GameplayFlag4)
	}
	core := p.proxy.core
	teams := append([]server.Team(nil), core.Teams.Arr...)
	count := core.Teams.ActiveCnt
	if sp.Teams < 0 || sp.Teams >= len(teams) {
		panic("report team count")
	}
	clear(core.Teams.Arr)
	core.Teams.ActiveCnt = sp.Teams
	for i := 1; i <= sp.Teams; i++ {
		tm := &core.Teams.Arr[i]
		tm.IDVal = server.TeamID(i)
		tm.ColorInd = server.TeamColor(i)
		// Initialize the owned team's slot and active fields; actual team iteration is exercised.
		*(*byte)(unsafe.Add(tm.C(), 58)) = byte(i)
		*(*uint32)(unsafe.Add(tm.C(), 64)) = 1
	}
	return func() {
		for i, off := range textOffsets {
			copy(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), 56), textBefore[i])
		}
		copy(core.Teams.Arr, teams)
		core.Teams.ActiveCnt = count
		noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
		noxflags.SetGamePlay(flags)
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
		p.reportRules = nil
	}
}
func (s *portTestRoamOwnerServer) gameplayReportRules() *portTestGameplayReportRules {
	if s.callbacks == nil || s.callbacks.shop == nil || s.callbacks.shop.pools == nil {
		return nil
	}
	return s.callbacks.shop.pools.reportRules
}
func (s *portTestRoamOwnerServer) GetFlag3592() bool {
	if st := s.gameplayReportRules(); st != nil {
		return st.spec.CountdownActive
	}
	return s.portTestRandomServer.GetFlag3592()
}
func (s *portTestRoamOwnerServer) ServStartCountdown(seconds int, text strman.ID) {
	if st := s.gameplayReportRules(); st != nil {
		st.countdown = append(st.countdown, PortTestGameplayReportCountdown{seconds, string(text)})
		return
	}
	s.portTestRandomServer.ServStartCountdown(seconds, text)
}

func (p *portTestShopPools) gameplayReportRulesMembers() {
	st := p.reportRules
	if st == nil || st.membersReady {
		return
	}
	st.membersReady = true
	for i := range p.proxy.life.players {
		u := &p.proxy.life.players[i]
		p.identify(u.TeamVal.C(), 620000+uint32(i))
		u.TeamVal.Field0 = 0
		if tm := p.proxy.core.Teams.ByID(u.TeamVal.ID); tm != nil {
			head := (*uint32)(unsafe.Add(tm.C(), 44))
			u.TeamVal.Field0 = *head
			*head = uint32(uintptr(u.TeamVal.C()))
		}
	}
}
