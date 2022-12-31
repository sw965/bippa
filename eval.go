package bippa

import (
	"math/rand"
)

type BattleEval struct {
	Func    func(*Battle) (float64, error)
	Reverse func(float64) float64
}

func NewPlayoutBattleEval(trainer Trainer, random *rand.Rand) BattleEval {
	evalFunc := func(battle *Battle) (float64, error) {
		var winner Winner
		var err error

		if battle.IsGameEnd() {
			winner, err = battle.Winner()
		} else {
			battleV := *battle
			winner, err = trainer.Playout(trainer, battleV, random)
			if err != nil {
				return 0.0, err
			}
		}
		return WINNER_TO_REWARD[winner], err
	}

	reverse := func(evalY float64) float64 {
		//自分が勝ち(1)なら相手は負け(0)
		//自分が負け(0)なら相手は勝ち(1)
		//引き分けなら互いに(0.5)
		return 1 - evalY
	}
	return BattleEval{Func: evalFunc, Reverse: reverse}
}

type PokemonEval func(*Pokemon) (float64, error)
type TeamEval func(Team) (float64, error)