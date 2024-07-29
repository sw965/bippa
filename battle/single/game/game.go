package game

import (
	//"fmt"
	//"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle/single"
    //bp "github.com/sw965/bippa"
)

// func Equal(b1, b2 *single.Battle) bool {
// 	return b1.SelfFighters.Equal(b2.SelfFighters) && b1.OpponentFighters.Equal(b2.OpponentFighters) && b1.Turn == b2.Turn
// }

// func IsEnd(b *single.Battle) (bool, []float64) {
// 	isP1AllFaint := b.SelfFighters.IsAllFaint()
// 	isP2AllFaint := b.OpponentFighters.IsAllFaint()

// 	if isP1AllFaint && isP2AllFaint {
// 		return true, []float64{0.5, 0.5}
// 	} else if isP1AllFaint {
// 		return true, []float64{0.0, 1.0}
// 	} else if isP2AllFaint {
// 		return true, []float64{1.0, 0.0}
// 	} else {
// 		return false, []float64{}
// 	}
// }

func LegalSeparateActionsSlice(b *single.Battle) single.ActionsSlice {
	moveSoloActionsSlice := b.LegalSeparateMoveSoloActionsSlice()
	ret := make(single.ActionsSlice, len(moveSoloActionsSlice))
	for i, action := range moveSoloActionsSlice {
		ret[i] = action.ToActions()
	}
	return ret
}

func NewPushFunc(context *single.Context) func(single.Battle, single.Actions) (single.Battle, error) {
	return func(battle single.Battle, actions single.Actions) (single.Battle, error) {
		battle = battle.Clone()
		for _, soloAction := range actions.ToSoloActions() {
			if soloAction.IsSelfView {
				battle.SoloAction(&soloAction, context)
			} else {
				battle.SwapView()
				battle.SoloAction(&soloAction, context)
				battle.SwapView()
			}
		}
		battle.Turn += 1
		return battle, nil
	}
}

// func New(context *single.Context) simultaneous.Game[single.Battle, single.ActionSlices, single.Actions, single.Action] {
//     gm := simultaneous.Game[single.Battle, single.ActionSlices, single.Actions, single.Action]{
//         Equal:                Equal,
//         IsEnd:                IsEnd,
//         LegalSeparateActions: LegalSeparateActions,
//         Push:                 NewPushFunc(context),
//     }
//     return gm
// }