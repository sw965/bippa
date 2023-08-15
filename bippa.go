package bippa

import (
	"math/rand"
)

func IsHit(p int, r *rand.Rand) bool {
	return p > r.Intn(100)
}

type Level int

const (
	DEFAULT_LEVEL = Level(50)
)

type Ability string
type Abilities []Ability

type stateCalculator struct{}
var StateCalculator = stateCalculator{}

func (sc *stateCalculator) HP(base BaseState, iv Individual, ev Effort) State {
	lv := int(DEFAULT_LEVEL)
	y := ((int(base)*2)+int(iv)+(int(ev)/4))*lv/100 + lv + 10
	return State(y)
}

func (sc *stateCalculator) OtherThanHP(base BaseState, iv Individual, ev Effort, bonus NatureBonus) State {
	lv := int(DEFAULT_LEVEL)
	y := ((int(base)*2)+int(iv)+(int(ev)/4))*lv/100 + 5
	return State(float64(y) * float64(bonus))
}

type State int
type States []State

type CriticalRank int

// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%99%E3%81%B0%E3%82%84%E3%81%95#.E8.A9.B3.E7.B4.B0.E3.81.AA.E4.BB.95.E6.A7.98
type SpeedBonus int

const (
	INIT_SPEED_BONUS = SpeedBonus(4096)
)

func NewSpeedBonus(bt *Battle) SpeedBonus {
	y := int(INIT_SPEED_BONUS)
	if bt.P1Fighters[0].Item == CHOICE_SCARF {
		y = FiveOrMoreRounding(float64(y) * 6144.0 / 4096.0)
	}
	return SpeedBonus(y)
}

type FinalSpeed float64

func NewFinalSpeed(bt *Battle) FinalSpeed {
	speed := bt.P1Fighters[0].Speed
	rankBonus := bt.P1Fighters[0].RankState.Speed.ToBonus()
	bonus := NewSpeedBonus(bt)

	y := FiveOrMoreRounding(float64(speed) * float64(rankBonus))
	y = FiveOverRounding(float64(y) * float64(bonus) / 4096.0)
	return FinalSpeed(y)
}

type Action struct {
	MoveName MoveName
	PokeName PokeName
}

func NewMoveUseAction(moveName MoveName) Action {
	return Action{MoveName:moveName}
}

func NewSwitchAction(pokeName PokeName) Action {
	return Action{PokeName:pokeName}
}

func (a Action) Priority() int {
	if a.MoveName == WA_RU_A_GA_KI {
		return 0
	} else if a.MoveName != NO_MOVE_NAME {
		return MOVEDEX[a.MoveName].PriorityRank
	} else {
		return 999
	}
}

type Actions []Action

type Winner struct {
	IsPlayer1 bool
	IsPlayer2 bool
}

var (
	WINNER_PLAYER1 = Winner{IsPlayer1: true, IsPlayer2: false}
	WINNER_PLAYER2 = Winner{IsPlayer1: false, IsPlayer2: true}
	DRAW           = Winner{IsPlayer1: false, IsPlayer2: false}
)

var WINNER_REWARD = map[Winner]float64{WINNER_PLAYER1: 1.0, WINNER_PLAYER2: 0.0, DRAW: 0.5}

// 小数点以下が0.5以上ならば、繰り上げ
func FiveOrMoreRounding(x float64) int {
	afterTheDecimalPoint := float64(x) - float64(int(x))
	if afterTheDecimalPoint >= 0.5 {
		return int(x + 1)
	}
	return int(x)
}

// 小数点以下が0.5より大きいならば、繰り上げ
func FiveOverRounding(x float64) int {
	afterTheDecimalPoint := float64(x) - float64(int(x))
	if afterTheDecimalPoint > 0.5 {
		return int(x + 1)
	}
	return int(x)
}