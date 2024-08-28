package battle

import (
	bp "github.com/sw965/bippa"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/math/rand"
)

type SoloAction struct {
	MoveName bp.MoveName
	SrcIndex int
	TargetIndex int
	IsSelfLeadTarget bool
	Speed int
	IsSelf bool
}

func (a *SoloAction) IsMove() bool {
	return a.MoveName != bp.EMPTY_MOVE_NAME
}

func (a *SoloAction) Priority() int {
	if a.IsMove() {
		moveData := bp.MOVEDEX[a.MoveName]
		return moveData.PriorityRank
	} else {
		return 999
	}
}

func (a *SoloAction) ToggleIsSelf() {
	a.IsSelf = !a.IsSelf
}

type SoloActions []SoloAction

func (as SoloActions) SortByOrder(m *Manager) {
	slices.SortFunc(as, func(a1, a2 SoloAction) bool {
		a1Priority := a1.Priority()
		a2Priority := a2.Priority()
		if a1Priority > a2Priority {
			return true
		} else if a1Priority < a2Priority {
			return false
		} else {
			a1Speed := a1.Speed
			a2Speed := a2.Speed
			if a1Speed > a2Speed {
				return !m.IsTrickRoomState()
			} else if a1Speed < a2Speed {
				return m.IsTrickRoomState()
			} else {
				return omwrand.Bool(GlobalContext.Rand)
			}
		}
	})
}

func (as SoloActions) ToggleIsSelf() {
	for i := range as {
		as[i].ToggleIsSelf()
	}
}

type SoloActionsSlice []SoloActions

type Action [DOUBLE_BATTLE_NUM]SoloAction
type Actions []Action

func (as Actions) ToSoloActions() SoloActions {
	sas := make(SoloActions, 0, len(as) * DOUBLE_BATTLE_NUM)
	for _, a := range as {
		for _, soloAction := range a {
			sas = append(sas, soloAction)
		}
	}
	return sas
}

type ActionsSlice []Actions