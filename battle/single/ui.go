package single

import (
// 	"fmt"
 	bp "github.com/sw965/bippa"
 	bt "github.com/sw965/bippa/battle"
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

// type DisplayUI struct {
// 	RealSelfPokeName string
// 	RealSelfLevel bp.Level
// 	RealSelfMaxHP int
// 	RealSelfCurrentHP int

// 	RealOpponentPokeName string
// 	RealOpponentLevel bp.Level
// 	RealOpponentMaxHP int
// 	RealOpponentCurrentHP int

// 	Message bt.Message
// }

// func NewDisplayUI(battle *Battle, msg bt.Message) (DisplayUI) {
// 	var realSelfPokemon bp.Pokemon
// 	var realOpponentPokemon bp.Pokemon

// 	if battle.IsRealSelf {
// 		realSelfPokemon = battle.SelfFighters[0]
// 		realOpponentPokemon = battle.OpponentFighters[0]
// 	} else {
// 		realSelfPokemon = battle.OpponentFighters[0]
// 		realOpponentPokemon = battle.SelfFighters[0]
// 	}

// 	return DisplayUI{
// 		RealSelfPokeName:realSelfPokemon.Name.ToString(),
// 		RealSelfLevel:realSelfPokemon.Level,
// 		RealSelfMaxHP:realSelfPokemon.MaxHP,
// 		RealSelfCurrentHP:realSelfPokemon.CurrentHP,

// 		RealOpponentPokeName:realOpponentPokemon.Name.ToString(),
// 		RealOpponentLevel:realOpponentPokemon.Level,
// 		RealOpponentMaxHP:realOpponentPokemon.MaxHP,
// 		RealOpponentCurrentHP:realOpponentPokemon.CurrentHP,

// 		Message:msg,
// 	}
// }

// func (ui *DisplayUI) CurrentHP(isRealSelf bool) int {
// 	if isRealSelf {
// 		return ui.RealSelfCurrentHP
// 	} else {
// 		return ui.RealOpponentCurrentHP
// 	}
// }

// func (ui DisplayUI) DecrementCurrentHP(isRealSelf bool) DisplayUI {
// 	if isRealSelf {
// 		ui.RealSelfCurrentHP -= 1
// 	} else {
// 		ui.RealOpponentCurrentHP -= 1
// 	}
// 	return ui
// }

// func (ui DisplayUI) Conceal(isRealSelf bool) DisplayUI {
// 	if isRealSelf {
// 		ui.RealSelfPokeName = "nil"
// 		ui.RealSelfLevel = -1
// 		ui.RealSelfMaxHP = -1
// 		ui.RealSelfCurrentHP = -1
// 	} else {
// 		ui.RealOpponentPokeName = "nil"
// 		ui.RealOpponentLevel = -1
// 		ui.RealOpponentMaxHP = -1
// 		ui.RealOpponentCurrentHP = -1
// 	}
// 	return ui
// }

// func (ui DisplayUI) Clone() DisplayUI {
// 	return ui
// }

// type DisplayUIs []DisplayUI

// func (uis DisplayUIs) LastElementSlice() DisplayUIs {
// 	ret := make(DisplayUIs, 0, cap(uis))
// 	ret = append(ret, omwslices.End(uis))
// 	return ret
// }

// func (uis DisplayUIs) AppendMessage(msg bt.Message, clearMessage bool) DisplayUIs {
// 	end := omwslices.End(uis).Clone()
// 	if clearMessage {
// 		end.Message = ""
// 	}
// 	for _, m := range msg.Accumulate() {
// 		end.Message = m
// 		uis = append(uis, end)
// 		end = end.Clone()
// 	}
// 	return uis
// }

// func (uis DisplayUIs) AppendDecrementCurrentHP(dmg int, isRealSelf bool) (DisplayUIs, error) {
// 	if dmg < 0 {
// 		msg := fmt.Sprintf("ダメージ = %d になっている。ダメージは0以上にする必要があります。", dmg)
// 		return DisplayUIs{}, fmt.Errorf(msg)
// 	}

// 	end := omwslices.End(uis)
// 	if end.CurrentHP(isRealSelf) <= 0 {
// 		return DisplayUIs{}, fmt.Errorf("現在のHPが既に0以下なのに、DecrimentCurrentHPを実行しようとした。")
// 	}
// 	for i := 0; i < dmg; i++ {
// 		end = end.DecrementCurrentHP(isRealSelf)
// 		uis = append(uis, end)
// 		if end.CurrentHP(isRealSelf) == 0 {
// 			break
// 		}
// 	}
// 	return uis, nil
// }

// type ObserverUI struct {
// 	LastSelfViewBattle Battle
// 	LastOpponentViewBattle Battle

// 	RealSelfLastUsedMoveName bp.MoveName
// 	RealOpponentLastUsedMoveName bp.MoveName

// 	Displays DisplayUIs

// 	RealSelfTrainerName bp.TrainerName
// 	RealOpponentTrainerName bp.TrainerName
// }

// func NewObserverUI(battle *Battle, c int) (ObserverUI, error) {
// 	displays := make(DisplayUIs, 0, c)
// 	displays = append(displays, NewDisplayUI(battle, ""))
// 	if battle.IsRealSelf {
// 		return ObserverUI{
// 			LastSelfViewBattle:*battle,
// 			LastOpponentViewBattle:battle.SwapView(),
// 			Displays:displays,
// 		}, nil
// 	} else {
// 		return ObserverUI{}, fmt.Errorf("NewUIの引数の*Battleは、Battle.IsRealSelf = true でなければならない")
// 	}
// }

// func (ui *ObserverUI) LastBattle(isSelfView bool) Battle {
// 	if isSelfView {
// 		return ui.LastSelfViewBattle
// 	} else {
// 		return ui.LastOpponentViewBattle
// 	}
// }

// func (ui *ObserverUI) TrainerName(isSelf bool) bp.TrainerName {
// 	if isSelf {
// 		return ui.RealSelfTrainerName
// 	} else {
// 		return ui.RealOpponentTrainerName
// 	}
// }

// func (ui *ObserverUI) SetLastUsedMoveName(battle *Battle) error {
// 	lastBattle := ui.LastBattle(battle.IsRealSelf)
// 	lastMoveset := lastBattle.SelfFighters[0].Moveset
// 	usedMoveNames := make(bp.MoveNames, 0, 1)
// 	for moveName, pp := range battle.SelfFighters[0].Moveset {
// 		lastPP, ok := lastMoveset[moveName]
// 		if !ok {
// 			return fmt.Errorf("一つ前の状態に存在しない技が含まれている。")
// 		}
// 		if lastPP.Current > pp.Current {
// 			usedMoveNames = append(usedMoveNames, moveName)
// 		}
// 	}

// 	if len(usedMoveNames) == 1 {
// 		usedMoveName := usedMoveNames[0]
// 		if battle.IsRealSelf {
// 			ui.RealSelfLastUsedMoveName = usedMoveName
// 		} else {
// 			ui.RealOpponentLastUsedMoveName = usedMoveName
// 		}
// 		return nil
// 	} else {
// 		return fmt.Errorf("最後に使用した技が二つ以上あるもしくは、最後に使用した技を特定出来ない。")
// 	}
// }

// func (ui *ObserverUI) LastUsedMoveName(isRealSelf bool) bp.MoveName {
// 	if isRealSelf {
// 		return ui.RealSelfLastUsedMoveName
// 	} else {
// 		return ui.RealOpponentLastUsedMoveName
// 	}
// }

// func (ui *ObserverUI) Observer(battle *Battle, eventType EventType) {
// 	isRealSelf := battle.IsRealSelf
// 	lastBattle := ui.LastBattle(isRealSelf)

// 	switch eventType {
// 		case MOVE_USE_EVENT:
// 			err := ui.SetLastUsedMoveName(battle)
// 			if err != nil {
// 				panic(err)
// 			}
// 			lastUsedMoveName := ui.LastUsedMoveName(isRealSelf)
// 			selfPokeName := battle.SelfFighters[0].Name
// 			msg := bt.NewMoveUseMessage(selfPokeName, lastUsedMoveName, isRealSelf)
// 			ui.Displays = ui.Displays.AppendMessage(msg, true)
// 		case ATTACK_MOVE_DAMAGE_EVENT:
// 			var err error
// 			dmg := lastBattle.OpponentFighters[0].CurrentHP - battle.OpponentFighters[0].CurrentHP
// 			ui.Displays, err = ui.Displays.AppendDecrementCurrentHP(dmg, !isRealSelf)
// 			if err != nil {
// 				panic(err)
// 			}

// 			lastUsedMoveName := ui.LastUsedMoveName(isRealSelf)
// 			lastUsedMoveData := bp.MOVEDEX[lastUsedMoveName]
// 			opponentPokeName := battle.OpponentFighters[0].Name
// 			opponentPokeData := bp.POKEDEX[opponentPokeName]

// 			effectType := lastUsedMoveData.Type.EffectType(opponentPokeData.Types)
// 			effectMsg := bt.NewEffectMessage(effectType)
// 			ui.Displays = ui.Displays.AppendMessage(effectMsg, true)
// 		case SWITCH_EVENT:
// 			trainerName := ui.TrainerName(battle.IsRealSelf)

// 			//○○は○○を引っ込めた！
// 			msg := bt.NewBackMessage(trainerName, lastBattle.SelfFighters[0].Name, isRealSelf)
// 			ui.Displays = ui.Displays.AppendMessage(msg, true)

// 			//ポケモンのUIを隠す
// 			ui.Displays = append(ui.Displays, omwslices.End(ui.Displays).Conceal(isRealSelf))

// 			//○○は○○を繰り出した！
// 			msg = bt.NewGoMessage(trainerName, battle.SelfFighters[0].Name, isRealSelf)
// 			ui.Displays = ui.Displays.AppendMessage(msg, true)

// 			//ポケモンを出現させる
// 			lastMsg := omwslices.End(ui.Displays).Message
// 			ui.Displays = append(ui.Displays, NewDisplayUI(battle, lastMsg))
// 		case SELF_FAINT_EVENT:
// 			trainerName := ui.TrainerName(isRealSelf)
// 			pokeName := battle.SelfFighters[0].Name
// 			msg := bt.NewFaintMessage(trainerName, pokeName, isRealSelf)
// 			ui.Displays = ui.Displays.AppendMessage(msg, true)
// 		case OPPONENT_FAINT_EVENT:
// 			pokeName := battle.OpponentFighters[0].Name
// 			trainerName := ui.TrainerName(!isRealSelf)
// 			msg := bt.NewFaintMessage(trainerName, pokeName, battle.IsRealSelf)
// 			ui.Displays = ui.Displays.AppendMessage(msg, true)
// 		case RECOIL_EVENT:
// 			trainerName := ui.TrainerName(isRealSelf)
// 			msg := bt.NewRecoilMessage(trainerName, lastBattle.SelfFighters[0].Name)
// 			ui.Displays = ui.Displays.AppendMessage(msg, true)
// 	}

// 	if battle.IsRealSelf {
// 		ui.LastSelfViewBattle = battle.Clone()
// 	} else {
// 		ui.LastOpponentViewBattle = battle.Clone()
// 	}
// }

type DisplayUI struct {
	RealSelfLeadPokeNames []string
	RealSelfLevels []bp.Level
	RealSelfMaxHPs []int
	RealSelfCurrentHPs []int

	RealOpponentPokeNames []string
	RealOpponentLevels []bp.Level
	RealOpponentMaxHPs []int
	RealOpponentCurrentHPs []int

	Message bt.Message
}