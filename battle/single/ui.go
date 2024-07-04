package single

import (
	"fmt"
	bp "github.com/sw965/bippa"
	bt "github.com/sw965/bippa/battle"
	omwslices "github.com/sw965/omw/slices"
)

type EventType int

const (
	MOVE_USE_EVENT EventType = iota
	OPPONENT_DAMAGE_EVENT
	SWITCH_EVENT
	SELF_FAINT_EVENT
	OPPONENT_FAINT_EVENT
	RECOIL_EVENT
)

type DisplayUI struct {
	RealSelfPokeName bp.PokeName
	RealSelfLevel bp.Level
	RealSelfMaxHP int
	RealSelfCurrentHP int

	RealOpponentPokeName bp.PokeName
	RealOpponentLevel bp.Level
	RealOpponentMaxHP int
	RealOpponentCurrentHP int

	Message bt.Message
}

func NewDisplayUI(battle *Battle, msg bt.Message) DisplayUI {
	var realSelfPokemon bp.Pokemon
	var realOpponentPokemon bp.Pokemon

	if battle.IsRealSelf {
		realSelfPokemon = battle.SelfFighters[0]
		realOpponentPokemon = battle.OpponentFighters[0]
	} else {
		realSelfPokemon = battle.OpponentFighters[0]
		realOpponentPokemon = battle.SelfFighters[0]
	}

	return DisplayUI{
		RealSelfPokeName:realSelfPokemon.Name,
		RealSelfLevel:realSelfPokemon.Level,
		RealSelfMaxHP:realSelfPokemon.MaxHP,
		RealSelfCurrentHP:realSelfPokemon.CurrentHP,

		RealOpponentPokeName:realOpponentPokemon.Name,
		RealOpponentLevel:realOpponentPokemon.Level,
		RealOpponentMaxHP:realOpponentPokemon.MaxHP,
		RealOpponentCurrentHP:realOpponentPokemon.CurrentHP,

		Message:msg,
	}
}

type DisplayUIs []DisplayUI

func (uis DisplayUIs) 

type ObserverUI struct {
	LastSelfViewBattle Battle
	LastOpponentViewBattle Battle
	Displays DisplayUIs
	SelfTrainerName string
	OpponentTrainerName string
}

func NewObserverUI(battle *Battle) (ObserverUI, error) {
	if battle.IsRealSelf {
		return ObserverUI{
			LastSelfViewBattle:*battle,
			LastOpponentViewBattle:battle.SwapView(),
		}, nil
	} else {
		return ObserverUI{}, fmt.Errorf("NewUIの引数の*Battleは、Battle.IsRealSelf = true でなければならない")
	}
}

func (o *ObserverUI) LastBattle(isSelfView bool) Battle {
	if isSelfView {
		return o.LastSelfViewBattle
	} else {
		return o.LastOpponentViewBattle
	}
}

func (o *ObserverUI) TrainerName(isSelf bool) string {
	if isSelf {
		return o.SelfTrainerName
	} else {
		return o.OpponentTrainerName
	}
}

func (o *ObserverUI) LastUsedMoveName(battle *Battle) (bp.MoveName, error) {
	lastBattle := o.LastBattle(battle.IsRealSelf)
	lastMoveset := lastBattle.SelfFighters[0].Moveset
	usedMoveNames := make(bp.MoveNames, 0, 1)
	for moveName, pp := range battle.SelfFighters[0].Moveset {
		lastPP, ok := lastMoveset[moveName]
		if !ok {
			var moveName bp.MoveName
			return moveName, fmt.Errorf("一つ前の状態に存在しない技が含まれている。")
		}
		if lastPP.Current > pp.Current {
			usedMoveNames = append(usedMoveNames, moveName)
		}
	}

	if len(usedMoveNames) == 1 {
		return usedMoveNames[0], nil
	} else {
		var moveName bp.MoveName
		return moveName, fmt.Errorf("最後に使用した技が二つある判定になっている。")
	}
}

func (o *ObserverUI) OpponentCurrentHPDiff(battle *Battle) int {
	lastBattle := o.LastBattle(battle.IsRealSelf)
	diff := battle.OpponentFighters[0].CurrentHP - lastBattle.OpponentFighters[0].CurrentHP
	return diff
}

func (o *ObserverUI) Observer(battle *Battle, eventType EventType) {
	switch eventType {
		case MOVE_USE_EVENT:
			lastUsedMoveName, err := o.LastUsedMoveName(battle)
			if err != nil {
				panic(err)
			}
			for _, msg := range bt.NewMoveUseMessage(battle.SelfFighters[0].Name, lastUsedMoveName, battle.IsRealSelf).Accumulate() {
				o.Displays = append(o.Displays, NewDisplayUI(battle, msg))
			}
		case OPPONENT_DAMAGE_EVENT:
			endDisplay := omwslices.End(o.Displays)
			lastMsg := endDisplay.Message
			dmg := (o.OpponentCurrentHPDiff(battle) * -1)
			if dmg < 0 {
				panic("OPPONENT_DAMAGE_EVENTで回復している。")
			}
			lastBattle := o.LastBattle(battle.IsRealSelf).Clone()
			for i := 0; i < dmg; i++ {
				lastBattle.OpponentFighters[0].CurrentHP -= 1
				o.Displays = append(o.Displays, NewDisplayUI(&lastBattle, lastMsg))
				lastBattle = lastBattle.Clone()
			}
		case SWITCH_EVENT:
			lastBattle := o.LastBattle(battle.IsRealSelf)
			lastPokeName := lastBattle.SelfFighters[0].Name
			trainerName := o.TrainerName(battle.IsRealSelf)
			var endDisplay DisplayUI
			if len(o.Displays) == 0 {
				endDisplay = NewDisplayUI(&lastBattle, "")
			} else {
				endDisplay = omwslices.End(o.Displays).Clone()
			}
			for _, msg := range bt.NewBackMessage(trainerName, lastPokeName, battle.IsRealSelf).Accumulate() {
				endDisplay.Message = msg
				o.Displays = append(o.Displays, endDisplay)
				endDisplay = endDisplay.Clone()
			}

			endDisplay = 
			for _, msg := range bt.NewGoMessage(trainerName, battle.SelfFighters[0].Name, battle.IsRealSelf).Accumulate() {
			}
		case SELF_FAINT_EVENT:
			fmt.Println("自分は倒れた")
		case OPPONENT_FAINT_EVENT:
			fmt.Println("相手は倒れた")
		case RECOIL_EVENT:
			fmt.Println("反動")
	}
}