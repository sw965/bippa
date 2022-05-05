package bippa

import (
  "math/rand"
)

type Eval struct {
  Func func(*Battle) (float64, error)
  ReverseFunc func(float64) float64
}

func NewRandomPlayoutEval(trainer Trainer, random *rand.Rand) Eval {
  evalFunc := func(battle *Battle) (float64, error) {
    var gameWinner Winner
    var err error

    if battle.IsGameEnd() {
      gameWinner, err = battle.Winner()
    } else {
      battleV := *battle
      gameEndBattle, err := trainer.OneGame(battleV, random)
      if err != nil {
        return 0.0, err
      }
      gameWinner, err = gameEndBattle.Winner()
    }
    return WINNER_TO_SIGMOID_REWARD[gameWinner], err
  }

  reverseFunc := func(evalY float64) float64 {
    //自分が勝ち(1)なら相手は負け(0)
    //自分が負け(0)なら相手は勝ち(1)
    //引き分けなら互いに(0.5)
    return 1 - evalY
  }
  return Eval{Func:evalFunc, ReverseFunc:reverseFunc}
}
