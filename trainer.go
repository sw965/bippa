package bippa

import (
	"fmt"
	"math/rand"
)

type Trainer func(*Battle) (Action, error)

func NewRandomInstructionTrainer(random *rand.Rand) Trainer {
	result := func(battle *Battle) (Action, error) {
		legalActions := battle.P1Fighters.LegalActions()
		return legalActions.RandomChoice(random), nil
	}
	return result
}

// func NewGoodTrainer(random *rand.Rand) Trainer {
// 	result := func(battle *Battle) (Action, error) {
// 		if battle.P1Fighters[0].IsLeechSeed {
//
// 		}
// 	}
// }

func (p1Trainer Trainer) OneGame(p2Trainer Trainer, battle Battle, random *rand.Rand) (Battle, error) {
	if battle.IsGameEnd() {
		return Battle{}, fmt.Errorf("既にゲームが終了している状態でtrainer.OneGame関数を呼び出した")
	}

	var err error
	var action Action

	for {
		if battle.IsP1Phase() {
			action, err = p1Trainer(&battle)
		} else {
			p1Command := battle.P1Command
			battle.P1Command = ""
			battle = battle.Reverse()
			action, err = p2Trainer(&battle)
			battle = battle.Reverse()
			battle.P1Command = p1Command
		}

		if err != nil {
			return Battle{}, err
		}

		battle, err = battle.Push(action, random)
		if err != nil {
			return Battle{}, err
		}

		if battle.IsGameEnd() {
			break
		}
	}
	return battle, nil
}
