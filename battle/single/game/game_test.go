package game_test

import (
	"testing"
	"fmt"
	"github.com/sw965/bippa/battle/single"
	// "github.com/sw965/bippa/battle/single/game"
	bp "github.com/sw965/bippa"
	"github.com/sw965/omw/fn"
)

type Action struct {
	MoveName bp.MoveName
	SrcIndex int
	TargetIndex int
	Speed int
	IsAllyTarget bool
}

func (a *Action) ToEasyRead() EasyReadAction {
	return EasyReadAction{
		MoveName:a.MoveName.ToString(),
		SrcIndex:a.SrcIndex,
		TargetIndex:a.TargetIndex,
		Speed:a.Speed,
		IsAllyTarget:a.IsAllyTarget,
	}
}

type EasyReadAction struct {
	MoveName string
	SrcIndex int
	TargetIndex int
	Speed int
	IsAllyTarget bool
}

type Actions []Action

func LegalFunc(b *single.Battle) Actions {
	ret := make(Actions, 0, 100)
	for i, selfPokemon := range b.SelfLeadPokemons {
		if selfPokemon.IsFainted() {
			continue
		}
		speed := selfPokemon.Speed
		moveset := selfPokemon.Moveset
		for j, opponentPokemon := range b.OpponentLeadPokemons {
			if opponentPokemon.IsFainted() {
				continue
			}
			for moveName, _ := range moveset {
				ret = append(ret, Action{
					MoveName:moveName,
					SrcIndex:i,
					TargetIndex:j,
					Speed:speed,
					IsAllyTarget:false,
				})
			}
		}

		for j, selfTargetPokemon := range b.SelfLeadPokemons {
			if selfTargetPokemon.IsFainted() {
				continue
			}
			for moveName, _ := range moveset {
				ret = append(ret, Action{
					MoveName:moveName,
					SrcIndex:i,
					TargetIndex:j,
					Speed:speed,
					IsAllyTarget:true,
				})
			}
		}

		for moveName, _ := range moveset {
			//対象指定なし
			ret = append(ret, Action{
				MoveName:moveName,
				SrcIndex:i,
				TargetIndex:-1,
				Speed:speed,
				IsAllyTarget:false,
			})
		}
	}

	ret = fn.Filter(ret, func(action Action) bool {
		moveData := bp.MOVEDEX[action.MoveName]
		fmt.Println("target", moveData.Target)
		if action.IsAllyTarget {
			if action.SrcIndex == action.TargetIndex {
				return false
			}
		}

		switch moveData.Target {
			case bp.NORMAL_TARGET:
				//単体への攻撃は、必ず対象を指定しなければならない(-1は対象指定なし)
				if action.TargetIndex == -1 {
					return false
				}

				//味方への攻撃
				if action.IsAllyTarget {
					//自分自身への攻撃は出来ない。
					return action.SrcIndex != action.TargetIndex
				} else {
					return true
				}
			case bp.OPPONENT_TWO_TARGET:
				//「いわなだれ」のような技は、対象指定は出来ない。(action.TargetIndex != -1 ならば 対象指定している)
				if action.TargetIndex != -1 {
					return false
				}

				//味方への攻撃は出来ない。
				return !action.IsAllyTarget
			case bp.SELF_TARGET:
				/*
				    「まもる」のように、自分に対象指定する技は、対象指定出来ないものとする。
					自分への対象指定と考える事も出来るが、
					対象指定の選択肢が、「自分への1つしかない」のであれば、
					自分の意思で対象指定をしている訳ではないので、
					対象指定は出来ないと見なしても、プログラム的には問題ない。
				*/

				if action.TargetIndex != -1 {
					return false
				}

				//味方への対象指定は出来ない。(action.IsAllyTargetは、自分への対象はfalseと見なす)
				return !action.IsAllyTarget
			case bp.OTHERS_TARGET:
				//「大爆発」や「地震」のような技

				if action.TargetIndex != -1 {
					return false
				}
				return !action.IsAllyTarget
			case bp.ALL_TARGET:
				//「あまごい」や「トリックルーム」のような技

				if action.TargetIndex != -1 {
					return false
				}
				return !action.IsAllyTarget
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				//「わるあがき」

				if action.TargetIndex != -1 {
					return false
				}
				return !action.IsAllyTarget
		}
		return true
	})
	return ret
}

func (as Actions) ToEasyRead() EasyReadActions {
	ret := make(EasyReadActions, len(as))
	for i, a := range as {
		ret[i] = a.ToEasyRead()
	}
	return ret
}

type EasyReadActions []EasyReadAction

func TestSingleBattleLegalSeparateActions(t *testing.T) {
	battle := single.Battle{
		SelfLeadPokemons:bp.Pokemons{bp.NewRomanStan2009Gyarados(), bp.NewRomanStan2009Metagross()},
		SelfBenchPokemons:bp.Pokemons{bp.NewRomanStan2009Latios()},
		OpponentLeadPokemons:bp.Pokemons{bp.NewKusanagi2009Empoleon(), bp.NewKusanagi2009Snorlax()},
		OpponentBenchPokemons:bp.Pokemons{bp.NewKusanagi2009Toxicroak()},
	}
	fmt.Println(LegalFunc(&battle).ToEasyRead())
}

func TestDoubleBattleLegalSeparateActions(t *testing.T) {

}