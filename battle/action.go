package battle

import (
	bp "github.com/sw965/bippa"
	"math/rand"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/omw/fn"
)

const DOUBLE_BATTLE = 2

type SoloAction struct {
	MoveName bp.MoveName
	SrcIndex int
	TargetIndex int
	IsSelfLeadTarget bool
	Speed int
	IsSelfView bool
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
	a.IsSelfView = !a.IsSelfView
}

type SoloActions []SoloAction

func (as SoloActions) SortByOrder(r *rand.Rand) {
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
				return true
			} else if a1Speed < a2Speed {
				return false
			} else {
				return omwrand.Bool(r)
			}
		}
	})
}

func (as SoloActions) ToggleIsSelf() {
	for i := range as {
		as[i].ToggleIsSelf()
	}
}

func (as SoloActions) ToActions() Actions {
	actions := make(Actions, 0, len(as) * len(as))
	for i, a1 := range as {
		for j, a2 := range as {
			action := Action{a1, a2}
			if i != j && !actions.Contains(action) {
				actions = append(actions, action)
			}
		} 
	}

	ret := fn.Filter(actions, func(action Action) bool {
		first := action[0]
		second := action[1]

		if first.SrcIndex == second.SrcIndex {
			return false
		}

		isFirstSwitch := !first.IsMove()
		isSecondSwitch := !second.IsMove()

		if isFirstSwitch && isSecondSwitch {
			if first.TargetIndex == second.TargetIndex {
				return false
			}
		}
		return true
	})

	if len(ret) == 0 {
		ret = make(Actions, len(as))
		for i, a := range as {
			ret[i][0] = a
		}
		return ret
	} else {
		return ret
	}
}

type SoloActionsSlice []SoloActions

type Action [DOUBLE_BATTLE]SoloAction
type Actions []Action

func (as Actions) Contains(a Action) bool {
	for _, e := range as {
		if e == a {
			return true
		}
	}
	return false
}

func (as Actions) ToSoloActions() SoloActions {
	ret := make(SoloActions, 0, len(as) * DOUBLE_BATTLE)
	for _, a := range as {
		for _, solo := range a {
			ret = append(ret, solo)
		}
	}
	return ret
}

type ActionsSlice []Actions