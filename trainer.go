package bippa

import (
	"fmt"
	"github.com/sw965/omw"
	"math/rand"
)

type Trainer func(*Battle) (ActionCmd, error)

func NewRandomInstructionTrainer(random *rand.Rand) Trainer {
	result := func(battle *Battle) (ActionCmd, error) {
		legalActionCmds := battle.P1Fighters.LegalActionCmds()
		return omw.RandomChoice(legalActionCmds, random), nil
	}
	return result
}

func (p1Trainer Trainer) Playout(p2Trainer Trainer, battle Battle, random *rand.Rand) (Winner, error) {
	if battle.IsGameEnd() {
		return Winner{}, fmt.Errorf("既にゲームが終了している状態でtrainer.OneGame関数を呼び出した")
	}

	var err error
	var actionCmd ActionCmd

	for {
		if battle.IsP1Phase() {
			actionCmd, err = p1Trainer(&battle)
		} else {
			p1ActionCmd := battle.P1ActionCmd
			battle.P1ActionCmd = ""
			battle = battle.Reverse()
			actionCmd, err = p2Trainer(&battle)
			battle = battle.Reverse()
			battle.P1ActionCmd = p1ActionCmd
		}

		if err != nil {
			return Winner{}, err
		}

		battle, err = battle.Push(actionCmd, random)
		if err != nil {
			return Winner{}, err
		}

		if battle.IsGameEnd() {
			break
		}
	}
	return battle.Winner()
}
