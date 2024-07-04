package single

import (
	"fmt"
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

type UI struct {
	LastSelfViewBattle single.Battle
	LastOpponentViewBattle single.Battle
}

func (ui *UI) LastBattle(isSelfView bool) Battle {
	
}

func (ui *UI) LastUsedMoveName(battle *Battle) {
	var lastBattle 

	battle.SelfFighters[0].Moveset 

}

func (ui *UI) Observer(battle *Battle, eventType EventType) {
	switch eventType {
		case MOVE_USE_EVENT:
			fmt.Println("攻撃")
		case OPPONENT_DAMAGE_EVENT:
			fmt.Println("相手はダメージを受けた")
		case SWITCH_EVENT:
			fmt.Println("交代する")
		case SELF_FAINT_EVENT:
			fmt.Println("自分は倒れた")
		case OPPONENT_FAINT_EVENT:
			fmt.Println("相手は倒れた")
		case RECOIL_EVENT:
			fmt.Println("反動")
	}
}