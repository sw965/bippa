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
	ITEM_USE_EVENT
	CURRENT_HP_EVENT
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

func (ui *ObserverUI) LastUsedMoveName(currentManager *Manager) (bp.MoveName, bp.PokeName, error) {
	lastManager := ui.LastManager(currentManager.IsPlayer1View)

	getLastUsedMoveName := func(current, last bp.Moveset) bp.MoveName {
		n := len(current)
		if n != len(last) {
			return bp.EMPTY_MOVE_NAME
		}

		mns := make(bp.MoveNames, 0, n)
		for k, v := range current {
			pp, ok := last[k]
			if !ok {
				return bp.EMPTY_MOVE_NAME
			}

			if v.Current < pp.Current {
				mns = append(mns, k)
			}
		}

		if len(mns) != 1 {
			return bp.EMPTY_MOVE_NAME
		}
		return mns[0]
	}

	for i, currentPokemon := range currentManager.SelfLeadPokemons {
		lastPokemon := lastManager.SelfLeadPokemons[i]
		if currentPokemon.Id != lastPokemon.Id {
			continue
		}
		lastUsedMoveName := getLastUsedMoveName(currentPokemon.Moveset, lastPokemon.Moveset)
		if lastUsedMoveName != bp.EMPTY_MOVE_NAME {
			return lastUsedMoveName, currentPokemon.Name, nil
		}
	}
	return bp.EMPTY_MOVE_NAME, bp.EMPTY_POKE_NAME, fmt.Errorf("最後に使った技を特定出来なかった")
}

func (ui *ObserverUI) SwitchPokeName(currentManager *Manager) (bp.PokeName, bp.PokeName, error) {
	lastManager := ui.LastManager(currentManager.IsPlayer1View)
	currentLeadIds := currentManager.SelfLeadPokemons.Ids()
	lastLeadIds := lastManager.SelfLeadPokemons.Ids()
	var beforeId int
	var afterId int

	for i, currentId := range currentLeadIds {
		lastLeadId := lastLeadIds[i]
		if currentId != lastLeadId {
			beforeId = lastLeadId
			afterId = currentId
			break
		}
	}

	beforePokemon, err := currentManager.SelfBenchPokemons.ById(beforeId)
	if err != nil {
		return bp.EMPTY_POKE_NAME, bp.EMPTY_POKE_NAME, err
	}

	afterPokemon, err := currentManager.SelfLeadPokemons.ById(afterId)
	return beforePokemon.Name, afterPokemon.Name, err
}

func (ui *ObserverUI) CurrentHPDiff(currentManager *Manager) (int, int, bool, error) {
	lastManager := ui.LastManager(currentManager.IsPlayer1View)

	get := func(currentPokemons, lastPokemons bp.Pokemons) (int, int) {
		for i, currentPokemon := range currentPokemons {
			lastPokemon := lastPokemons[i]
			currentHP := currentPokemon.Stat.CurrentHP
			lastCurrentHP := lastPokemon.Stat.CurrentHP
			if currentPokemon.Id == lastPokemon.Id && currentHP != lastCurrentHP {
				diff := currentHP - lastCurrentHP
				return diff, i
			}
		}
		return 0, -1
	}

	selfDiff, selfIdx := get(currentManager.SelfLeadPokemons, lastManager.SelfLeadPokemons)
	opponentDiff, opponentIdx := get(currentManager.OpponentLeadPokemons, lastManager.OpponentLeadPokemons)

	if selfDiff != 0 && opponentDiff != 0 {
		return 0, 0, false, fmt.Errorf("現在のHPの差分が2匹のポケモンで検知された為、差分計算が出来ません。")
	}

	if selfDiff != 0 {
		return selfDiff, selfIdx, true, nil
	}

	if opponentDiff != 0 {
		return opponentDiff, opponentIdx, false, nil
	}
	return 0, 0, false, fmt.Errorf("現在のHPの差分が特定出来ませんでした。")
}

func (ui *ObserverUI) LastUsedItem(currentManager *Manager) (bp.Item, bp.PokeName, bool, error) {
	lastManager := ui.LastManager(currentManager.IsPlayer1View)

	get := func(currentPokemons, lastPokemons bp.Pokemons) (bp.Item, bp.PokeName) {
		for i, currentPokemon := range currentPokemons {
			lastPokemon := lastPokemons[i]
			if currentPokemon.Id == lastPokemon.Id && currentPokemon.Item != lastPokemon.Item {
				return lastPokemon.Item, currentPokemon.Name
			}
		}
		return bp.EMPTY_ITEM, bp.EMPTY_POKE_NAME	
	}

	selfItem, selfPokeName := get(currentManager.SelfLeadPokemons, lastManager.SelfLeadPokemons)
	opponentItem, opponentPokeName := get(currentManager.OpponentLeadPokemons, lastManager.OpponentLeadPokemons)
	if selfItem != bp.EMPTY_ITEM && opponentItem != bp.EMPTY_ITEM {
		msg := fmt.Sprintf("最後に使ったアイテムが2つ検知された。 %s %s", selfItem.ToString(), opponentItem.ToString())
		return bp.EMPTY_ITEM, bp.EMPTY_POKE_NAME, false, fmt.Errorf(msg)
	}
	if selfItem != bp.EMPTY_ITEM {
		isSelf := true
		return selfItem, selfPokeName, isSelf, nil
	}

	if opponentItem != bp.EMPTY_ITEM {
		isSelf := false
		return opponentItem, opponentPokeName, isSelf, nil
	}
	return bp.EMPTY_ITEM, bp.EMPTY_POKE_NAME, false, fmt.Errorf("最後に使ったアイテムを特定出来なかった。")
}

func (ui *ObserverUI) Observer(current *Manager, event EventType) {
	lastManager := ui.LastManager(current.IsPlayer1View)
	switch event {
		case MOVE_USE_EVENT:
			lastUsedMoveName, pokeName, err := ui.LastUsedMoveName(current)
			if err != nil {
				panic(err)
			}
			fmt.Println(fmt.Sprintf("%s の %s!", pokeName.ToString(), lastUsedMoveName.ToString()))
		case SWITCH_EVENT:
			beforePokeName, afterPokeName, err := ui.SwitchPokeName(current)
			if err != nil {
				panic(err)
			}
			fmt.Println("戻れ！", beforePokeName.ToString())
			fmt.Println("行け！", afterPokeName.ToString())
		case CURRENT_HP_EVENT:
			diff, leadIdx, isSelf, err := ui.CurrentHPDiff(current)
			if err != nil {
				panic(err)
			}
			leadPokemons := current.LeadPokemons(isSelf)
			fmt.Println(diff, leadPokemons[leadIdx].Name.ToString(), isSelf)
		case ITEM_USE_EVENT:
			item, pokeName, isSelf, err := ui.LastUsedItem(current)
			if err != nil {
				panic(err)
			}
			fmt.Println(item.ToString(), isSelf, pokeName.ToString())
	}

	if lastManager.IsPlayer1View {
		ui.LastP1ViewManager = lastManager.Clone()
	} else {
		ui.LastP2ViewManager = lastManager.Clone()
	}
}