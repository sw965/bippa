package game

import (
    "fmt"
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle/single"
    bp "github.com/sw965/bippa"
)

func Equal(b1, b2 *single.Battle) bool {
	return b1.SelfFighters.Equal(b2.SelfFighters) && b1.OpponentFighters.Equal(b2.OpponentFighters) && b1.Turn == b2.Turn
}

func IsEnd(b *single.Battle) (bool, []float64) {
	isP1AllFaint := b.SelfFighters.IsAllFaint()
	isP2AllFaint := b.OpponentFighters.IsAllFaint()

	if isP1AllFaint && isP2AllFaint {
		return true, []float64{0.5, 0.5}
	} else if isP1AllFaint {
		return true, []float64{0.0, 1.0}
	} else if isP2AllFaint {
		return true, []float64{1.0, 0.0}
	} else {
		return false, []float64{}
	}
}

func LegalSeparateActions(b *single.Battle) single.ActionSlices {
	separateMoveNames := b.SeparateCommandableMoveNames()
	separatePokeNames := b.SeparateSwitchablePokeNames()
	ret := make(single.ActionSlices, single.PLAYER_NUM * single.LEAD_NUM)
	isSelfs := []bool{true, false}
	for playerI := range ret {
		isSelf := isSelfs[playerI]
		moveNames := separateMoveNames[playerI]
		pokeNames := separatePokeNames[playerI]
		actions := make(single.Actions, 0, len(moveNames) + len(pokeNames))
		for _, name := range moveNames {
			actions = append(actions, single.Action{CmdMoveName:name, IsSelf:isSelf})
		}
		for _, name := range pokeNames {
			actions = append(actions, single.Action{SwitchPokeName:name, IsSelf:isSelf})
		}
		if len(actions) == 0 {
			actions = append(actions, single.Action{IsSelf:isSelf})
		}
		ret[playerI] = actions
	}
	return ret
}

func NewPushFunc(context *single.Context) func(single.Battle, single.Actions) (single.Battle, error) {
	return func(battle single.Battle, actions single.Actions) (single.Battle, error) {
		if len(actions) != 2 {
			return single.Battle{}, fmt.Errorf("len(actions) != 2 (NewPushFunc)")
		}

		for actions[0].IsSelf == actions[1].IsSelf {
			return single.Battle{}, fmt.Errorf("プレイヤー1もしくはプレイヤー2が連続で行動しようとした。(actions[0].IsSelf == actions[1].IsSelf)")
		}
		if actions.IsAllEmpty() {
			fmt.Println("エラー前battle", battle.ToEasyRead())
			return single.Battle{}, fmt.Errorf("両プレイヤーのActionがEmptyになっているため、Pushできません。Emptyじゃないようにするには、Action.CmdMoveNameかAction.SwitchPokeNameのいずれかは、ゼロ値以外の値である必要があります。")
		}

		var err error
		sorted := battle.SortActionsByOrder(&actions[0], &actions[1], context.Rand)
		for i := range sorted {
			action := sorted[i]
			if action.CmdMoveName == bp.EMPTY_MOVE_NAME && action.SwitchPokeName == bp.EMPTY_POKE_NAME {
				continue
			}

			if action.IsSelf {
				battle, err = battle.Action(action, context)
			} else {
				battle = battle.SwapView()
				battle, err = battle.Action(action, context)
				battle = battle.SwapView()
			}

			if err != nil {
				return single.Battle{}, err
			}

			if battle.SelfFighters[0].IsFaint() {
				context.Observer(&battle, single.SELF_FAINT_EVENT)
			}

			if battle.OpponentFighters[0].IsFaint() {
				context.Observer(&battle, single.OPPONENT_FAINT_EVENT)
			}
		}
		battle.Turn += 1
		return battle, nil
	}
}

func New(context *single.Context) simultaneous.Game[single.Battle, single.ActionSlices, single.Actions, single.Action] {
    gm := simultaneous.Game[single.Battle, single.ActionSlices, single.Actions, single.Action]{
        Equal:                Equal,
        IsEnd:                IsEnd,
        LegalSeparateActions: LegalSeparateActions,
        Push:                 NewPushFunc(context),
    }
    return gm
}