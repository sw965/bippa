package bippa

import (
	"math/rand"
)

type Eval struct {
	Func    func(*Battle) (float64, error)
	Reverse func(float64) float64
}

func NewRandomPlayoutEval(trainer Trainer, random *rand.Rand) Eval {
	evalFunc := func(battle *Battle) (float64, error) {
		var gameWinner Winner
		var err error

		if battle.IsGameEnd() {
			gameWinner, err = battle.Winner()
		} else {
			battleV := *battle
			gameEndBattle, err := trainer.OneGame(trainer, battleV, random)
			if err != nil {
				return 0.0, err
			}
			gameWinner, err = gameEndBattle.Winner()
		}
		return WINNER_TO_SIGMOID_REWARD[gameWinner], err
	}

	reverse := func(evalY float64) float64 {
		//自分が勝ち(1)なら相手は負け(0)
		//自分が負け(0)なら相手は勝ち(1)
		//引き分けなら互いに(0.5)
		return 1 - evalY
	}
	return Eval{Func: evalFunc, Reverse: reverse}
}

type NotBad struct {
	W []float64
	B float64
}

func NewInitNotBad(random *rand.Rand) NotBad {
	notBad := NotBad{}
	w := make([]float64, (3 * FIGHTERS_LENGTH * FIGHTERS_LENGTH) * 8 * 8 * FIGHTERS_LENGTH * FIGHTERS_LENGTH)
	for i := 0; i < len(w); i++ {
		w[i] = random.Float64()
	}
	notBad.W = w
	notBad.B = 1.0
	return notBad
}

func (notBad *NotBad) Sum(x []float64) float64 {
	result := 0.0
	for i, ele := range x {
		result += (ele * notBad.W[i])
	}
	return result + notBad.B
}

func (notBad *NotBad) Output(x []float64) float64 {
	return Sigmoid(notBad.Sum(x))
}

func (notBad *NotBad) Eval() Eval {
	evalFunc := func(battle *Battle) (float64, error) {
		notBadEvalX, err := battle.NotBadEvalX()
		return notBad.Output(notBadEvalX), err
	}

	reverse := func(evalY float64) float64 {
		return 1.0 - evalY
	}
	return Eval{Func:evalFunc, Reverse:reverse}
}

func (notBad *NotBad) Gradients(x []float64, t float64) []float64 {
	result := make([]float64, len(x))
	y := notBad.Output()
	for i, xElement := range x {
		result[i] = SigmoidMeanSquaredErrorDerivative(y, t, xElement)
	}
	return result
}

func (notBad *NotBad) Learning(x []float64, t float64, learningRate float64) float64 {
	y := notBad.Output(x)
	gradiends := notBad.Gradients(x, y, t)
	for i, gradiend := range gradiends {
		gradient := SigmoidMeanSquaredErrorDerivative(x, y, t)
		notBad.W -= (gradiend * learningRate)
	}
	notBad.B -= (SigmoidMeanSquaredErrorDerivative(y, t, 1.0) * learningRate)
}
