package battle

import (
 	"fmt"
 	bp "github.com/sw965/bippa"
// 	omwslices "github.com/sw965/omw/slices"
)

type EventType int

const (
	MOVE_USE_EVENT EventType = iota
	ATTACK_MOVE_DAMAGE_EVENT
	SWITCH_EVENT
	SELF_FAINT_EVENT
	OPPONENT_FAINT_EVENT
	RECOIL_EVENT
)

type DisplayUI struct {
	P1LeadPokeNames []string
	P1LeadLevels []bp.Level
	P1LeadMaxHPs []int
	P1LeadCurrentHPs []int

	P2LeadPokeNames []string
	P2LeadLevels []bp.Level
	P2LeadMaxHPs []int
	P2LeadCurrentHPs []int

	Message Message
}

type ObserverUI struct {
	LastP1ViewManager Manager
	LastP2ViewManager Manager
}

func (ui *ObserverUI) LastManager(isPlayer1View bool) Manager {
	if isPlayer1View {
		return ui.LastP1ViewManager
	} else {
		return ui.LastP2ViewManager
	}
}

func (ui *ObserverUI) LastUsedMoveName(current *Manager, leadIdx int) (bp.MoveName, error) {
	last := ui.LastManager(current.IsPlayer1View)
	lastMoveset := last.SelfLeadPokemons[leadIdx].Moveset
	currentMoveset := current.SelfLeadPokemons[leadIdx].Moveset
	for k, v := range lastMoveset {
		if v.Current > currentMoveset[k].Current {
			return k, nil
		}
	}
	return bp.EMPTY_MOVE_NAME, fmt.Errorf("LastUsedMoveNameが特定出来なかった。")
}

func (ui *ObserverUI) Observer(current *Manager, event EventType, leadIdx int) {
	//lastManager := ui.LastManager(m.IsPlayer1View)
	switch event {
		case MOVE_USE_EVENT:
			fmt.Println("チェック")
			lastUsedMoveName, err := ui.LastUsedMoveName(current, leadIdx)
			if err != nil {
				panic(err)
			}
			msg := fmt.Sprintf("%s の %s", current.SelfLeadPokemons[leadIdx].Name.ToString(), lastUsedMoveName.ToString())
			fmt.Println(msg)
	}
}