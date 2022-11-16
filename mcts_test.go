package bippa

import (
	"testing"
	"fmt"
	//"math"
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
)

func TestRunMCTS(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	p1Fighters := Fighters{NEW_RENTAL_POKEMONS["サンダー"](), NEW_RENTAL_POKEMONS["リザードン"](), NEW_RENTAL_POKEMONS["カメックス"]()}
	p2Fighters := Fighters{NEW_RENTAL_POKEMONS["フシギバナ"](), NEW_RENTAL_POKEMONS["ゲンガー"](), NEW_RENTAL_POKEMONS["パルシェン"]()}
	battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}

// 	battle.P1Fighters[0].CurrentHP = 1
// 	battle.P1Fighters[1].CurrentHP = 0
// 	battle.P1Fighters[2].CurrentHP = 0

// 	battle.P2Fighters[1].CurrentHP = 0
// 	battle.P2Fighters[2].CurrentHP = 0

	randomInstructionTrainer := NewRandomInstructionTrainer(mtRandom)

	playoutNum := 36000
	p1WinCount := 0.0
	for i := 0; i < playoutNum; i++ {
		winner, err := randomInstructionTrainer.Playout(randomInstructionTrainer, battle, mtRandom)
		if err != nil {
			panic(err)
		}
		p1WinCount += WINNER_TO_REWARD[winner]
	}
	fmt.Println(p1WinCount / float64(playoutNum))

// 	randomPlayoutBattleEval := NewPlayoutBattleEval(randomInstructionTrainer, mtRandom)

// 	allNodes, err := RunMCTS(battle, 1280, math.Sqrt(2), NoBattlePolicy, &randomPlayoutBattleEval, mtRandom)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for actionCmd, pucb := range allNodes[0].ActionCmdPUCBs {
// 		fmt.Println(actionCmd, pucb.P, pucb.AccumReward, pucb.Trial)
// 		fmt.Println(pucb.AverageReward())
// 	}
}