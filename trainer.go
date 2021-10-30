package bippa

import (
	"bufio"
	"fmt"
	"math/rand"
)

type Trainer func(*Battle) (BattleCommand, error)

func NewRandomInstructionTrainer(random *rand.Rand) Trainer {
	result := func(battle *Battle) (BattleCommand, error) {
		availableBattleCommands := battle.P1Fighters.AvailableBattleCommands()
		return availableBattleCommands.RandomChoice(random), nil
	}
	return result
}

func NewScanTrainer(scanner *bufio.Scanner) Trainer {
	result := func(battle *Battle) (BattleCommand, error) {
		fmt.Println(battle.P1Fighters[0].Name + " はどうする？" )

		for {
			scanner.Scan()
			text := scanner.Text()
			battleCommand := BattleCommand(text)
			if battle.P1Fighters.IsAvailableBattleCommand(battleCommand) {
				return battleCommand, nil
			}
			fmt.Println("不適なコマンドなので、再度入力してね♡")
		}
	}
	return result
}

func (p1Trainer Trainer) OneGame(p2Trainer Trainer, battle Battle, random *rand.Rand) (Battle, error) {
	var err error
	var battleCommand BattleCommand

	if battle.IsGameEnd() {
		return Battle{}, fmt.Errorf("既にゲームが終了している状態でtrainer.OneGame関数を呼び出した")
	}

	for {
		if battle.IsP1Phase() {
			battleCommand, err = p1Trainer(&battle)
		} else {
			p1BattleCommand := battle.P1Command
			battle.P1Command = ""

			battle, err = battle.ReversePlayer()
			if err != nil {
				return Battle{}, err
			}

			battleCommand, err = p2Trainer(&battle)

			battle, err = battle.ReversePlayer()
			if err != nil {
				return Battle{}, err
			}

			battle.P1Command = p1BattleCommand
		}

		if err != nil {
			return Battle{}, err
		}

		battle, err = battle.Run(battleCommand, random)

		if err != nil {
			return Battle{}, err
		}

		if battle.IsGameEnd() {
			break
		}
	}
	return battle, nil
}
