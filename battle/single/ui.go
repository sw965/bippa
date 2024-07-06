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
	RealSelfPokeName string
	RealSelfLevel bp.Level
	RealSelfMaxHP int
	RealSelfCurrentHP int

	RealOpponentPokeName string
	RealOpponentLevel bp.Level
	RealOpponentMaxHP int
	RealOpponentCurrentHP int

	Message bt.Message
}

func NewDisplayUI(battle *Battle, msg bt.Message) (DisplayUI) {
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
		RealSelfPokeName:realSelfPokemon.Name.ToString(),
		RealSelfLevel:realSelfPokemon.Level,
		RealSelfMaxHP:realSelfPokemon.MaxHP,
		RealSelfCurrentHP:realSelfPokemon.CurrentHP,

		RealOpponentPokeName:realOpponentPokemon.Name.ToString(),
		RealOpponentLevel:realOpponentPokemon.Level,
		RealOpponentMaxHP:realOpponentPokemon.MaxHP,
		RealOpponentCurrentHP:realOpponentPokemon.CurrentHP,

		Message:msg,
	}
}

func (ui DisplayUI) Conceal(isRealSelf bool) DisplayUI {
	if isRealSelf {
		ui.RealSelfPokeName = "nil"
		ui.RealSelfLevel = -1
		ui.RealSelfMaxHP = -1
		ui.RealSelfCurrentHP = -1
	} else {
		ui.RealOpponentPokeName = "nil"
		ui.RealOpponentLevel = -1
		ui.RealOpponentMaxHP = -1
		ui.RealOpponentCurrentHP = -1
	}
	return ui
}

func (ui DisplayUI) Clone() DisplayUI {
	return ui
}

type DisplayUIs []DisplayUI

func (uis DisplayUIs) LastElementSlice() DisplayUIs {
	ret := make(DisplayUIs, 0, cap(uis))
	ret = append(ret, omwslices.End(uis))
	return ret
}

func (uis DisplayUIs) AppendMessage(msg bt.Message, clearMessage bool) DisplayUIs {
	end := omwslices.End(uis).Clone()
	if clearMessage {
		end.Message = ""
	}
	for _, m := range msg.Accumulate() {
		end.Message = m
		uis = append(uis, end)
		end = end.Clone()
	}
	return uis
}

func (uis DisplayUIs) AppendRealOpponentDamageOrHeal(diff int) DisplayUIs {
	var sign int
	if diff >= 0 {
		sign = 1
	} else {
		sign = -1
	}

	ui := omwslices.End(uis).Clone()
	n := sign * diff
	for i := 0; i < n; i++ {
		ui.RealOpponentCurrentHP += sign
		uis = append(uis, ui)
		ui = ui.Clone()
	}
	return uis
}

type ObserverUI struct {
	LastSelfViewBattle Battle
	LastOpponentViewBattle Battle
	Displays DisplayUIs
	SelfTrainerName string
	OpponentTrainerName string
}

func NewObserverUI(battle *Battle, c int) (ObserverUI, error) {
	displays := make(DisplayUIs, 0, c)
	displays = append(displays, NewDisplayUI(battle, ""))
	if battle.IsRealSelf {
		return ObserverUI{
			LastSelfViewBattle:*battle,
			LastOpponentViewBattle:battle.SwapView(),
			Displays:displays,
		}, nil
	} else {
		return ObserverUI{}, fmt.Errorf("NewUIの引数の*Battleは、Battle.IsRealSelf = true でなければならない")
	}
}

func (ui *ObserverUI) LastBattle(isSelfView bool) Battle {
	if isSelfView {
		return ui.LastSelfViewBattle
	} else {
		return ui.LastOpponentViewBattle
	}
}

func (ui *ObserverUI) TrainerName(isSelf bool) string {
	if isSelf {
		return ui.SelfTrainerName
	} else {
		return ui.OpponentTrainerName
	}
}

func (ui *ObserverUI) LastUsedMoveName(battle *Battle) (bp.MoveName, error) {
	lastBattle := ui.LastBattle(battle.IsRealSelf)
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

func (ui *ObserverUI) RealOpponentCurrentHPDiff(battle *Battle) int {
	lastBattle := ui.LastBattle(battle.IsRealSelf)
	diff := battle.OpponentFighters[0].CurrentHP - lastBattle.OpponentFighters[0].CurrentHP
	return diff
}

func (ui *ObserverUI) Observer(battle *Battle, eventType EventType) {
	switch eventType {
		case MOVE_USE_EVENT:
			lastUsedMoveName, err := ui.LastUsedMoveName(battle)
			if err != nil {
				panic(err)
			}
			msg := bt.NewMoveUseMessage(battle.SelfFighters[0].Name, lastUsedMoveName, battle.IsRealSelf)
			ui.Displays = ui.Displays.AppendMessage(msg, true)
		case OPPONENT_DAMAGE_EVENT:
			diff := ui.RealOpponentCurrentHPDiff(battle)
			if diff > 0 {
				panic("OPPONENT_DAMAGE_EVENTで、HPが回復している。")
			}
			ui.Displays = ui.Displays.AppendRealOpponentDamageOrHeal(diff)
		case SWITCH_EVENT:
			lastBattle := ui.LastBattle(battle.IsRealSelf)
			lastPokeName := lastBattle.SelfFighters[0].Name
			trainerName := ui.TrainerName(battle.IsRealSelf)

			//○○は○○を引っ込めた！
			msg := bt.NewBackMessage(trainerName, lastPokeName, battle.IsRealSelf)
			ui.Displays = ui.Displays.AppendMessage(msg, true)

			//ポケモンのUIを隠す
			ui.Displays = append(ui.Displays, omwslices.End(ui.Displays).Conceal(battle.IsRealSelf))

			//○○は○○を繰り出した！
			msg = bt.NewGoMessage(trainerName, battle.SelfFighters[0].Name, battle.IsRealSelf)
			ui.Displays = ui.Displays.AppendMessage(msg, true)

			//ポケモンを出現させる
			lastMsg := omwslices.End(ui.Displays).Message
			ui.Displays = append(ui.Displays, NewDisplayUI(battle, lastMsg))
		case SELF_FAINT_EVENT:
			lastBattle := ui.LastBattle(battle.IsRealSelf)
			trainerName := ui.TrainerName(battle.IsRealSelf)
			pokeName := lastBattle.SelfFighters[0].Name
			if pokeName != battle.SelfFighters[0].Name {
				panic("SELF_FAINT_EVENTで、直前のポケモン名と現在のポケモン名が異なっている。")
			}
			msg := bt.NewFaintMessage(trainerName, pokeName, battle.IsRealSelf)
			ui.Displays = ui.Displays.AppendMessage(msg, true)
		case OPPONENT_FAINT_EVENT:
			lastBattle := ui.LastBattle(battle.IsRealSelf)
			pokeName := lastBattle.OpponentFighters[0].Name
			trainerName := ui.TrainerName(!battle.IsRealSelf)
			if pokeName != battle.OpponentFighters[0].Name {
				panic("OPPONENT_FAINT_EVENTで、直前のポケモン名と現在のポケモン名が異なっている。")
			}
			msg := bt.NewFaintMessage(trainerName, pokeName, battle.IsRealSelf)
			ui.Displays = ui.Displays.AppendMessage(msg, true)
		case RECOIL_EVENT:
			trainerName := ui.TrainerName(battle.IsRealSelf)
			lastBattle := ui.LastBattle(battle.IsRealSelf)
			msg := bt.NewRecoilMessage(trainerName, lastBattle.SelfFighters[0].Name)
			ui.Displays = ui.Displays.AppendMessage(msg, true)
	}

	if battle.IsRealSelf {
		ui.LastSelfViewBattle = battle.Clone()
	} else {
		ui.LastOpponentViewBattle = battle.Clone()
	}
}