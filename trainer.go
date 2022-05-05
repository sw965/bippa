package bippa

import (
	"fmt"
	"math/rand"
)

type Trainer func(*Battle) (Action, error)

func NewRandomInstructionTrainer(random *rand.Rand) Trainer {
	result := func(battle *Battle) (Action, error) {
		legalActions := battle.LegalActions()
		return legalActions.RandomChoice(random), nil
	}
	return result
}

func (trainer Trainer) OneGame(battle Battle, random *rand.Rand) (Battle, error) {
	if battle.IsGameEnd() {
		return Battle{}, fmt.Errorf("既にゲームが終了している状態でtrainer.OneGame関数を呼び出した")
	}

	for {
		action, err := trainer(&battle)
		if err != nil {
			return Battle{}, err
		}

		battle, err = battle.Push(&action, random)
		if err != nil {
			return Battle{}, err
		}

		if battle.IsGameEnd() {
			break
		}
	}
	return battle, nil
}
