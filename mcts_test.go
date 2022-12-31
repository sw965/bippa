package bippa

import (
	"testing"
	"fmt"
	"math"
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
)

func Putera() Pokemon {
	pokemon, err := NewPokemon(
		"プテラ", "ようき", "プレッシャー", "♂", "なし",
		MoveNames{"いわなだれ", "こおりのキバ"},
		PointUps{MAX_POINT_UP, MAX_POINT_UP}, &IndividualState{}, &EffortState{},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func Gaburiasu() Pokemon {
	pokemon, err := NewPokemon(
		"ガブリアス", "ようき", "さめはだ", "♂", "なし",
		MoveNames{"ドラゴンクロー", "じしん"},
		PointUps{MAX_POINT_UP, MAX_POINT_UP}, &IndividualState{}, &EffortState{},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func Sanda() Pokemon {
	pokemon, err := NewPokemon(
		"サンダー", "ひかえめ", "プレッシャー", "不明", "なし",
		MoveNames{"10まんボルト", "はねやすめ"},
		PointUps{MAX_POINT_UP, MAX_POINT_UP}, &IndividualState{}, &EffortState{},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func Rapurasu() Pokemon {
	pokemon, err := NewPokemon(
		"ラプラス", "いじっぱり", "ちょすい", "♀", "なし",
		MoveNames{"れいとうビーム", "こおりのつぶて"},
		PointUps{MAX_POINT_UP, MAX_POINT_UP}, &IndividualState{}, &EffortState{},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func TestRunMCTS(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	p1Fighters := Fighters{Putera(), Gaburiasu(), Sanda()}
	p2Fighters := Fighters{Gaburiasu(), Rapurasu(), Sanda()}

	battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}

	battle.P1Fighters[0].CurrentHP = 1
	battle.P1Fighters[1].CurrentHP = 0
	battle.P1Fighters[2].CurrentHP = 0

	battle.P2Fighters[0].CurrentHP = 100
	battle.P2Fighters[1].CurrentHP = 100
	battle.P2Fighters[2].CurrentHP = 100

	randomInstructionTrainer := NewRandomInstructionTrainer(mtRandom)
	randomPlayoutBattleEval := NewPlayoutBattleEval(randomInstructionTrainer, mtRandom)

	var allNodes BattleNodes
	var err error

	for i := 0; i < 1; i++ {
		allNodes, err = RunMCTS(battle, 25600, math.Sqrt(2), NoBattlePolicy, &randomPlayoutBattleEval, mtRandom)
		if err != nil {
			panic(err)
		}
		//fmt.Println(len(allNodes))
	}

	for actionCmd, pucb := range allNodes[0].ActionCmdPUCBs {
		fmt.Println(actionCmd, pucb.P, pucb.AccumReward, pucb.Trial)
		fmt.Println(pucb.AverageReward())
	}
}
