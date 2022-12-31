package bippa

import (
	"fmt"
)

func InLearnset(pokeName PokeName, moveName MoveName) bool {
	pokeData := POKEDEX[pokeName]
	for _, iMoveName := range pokeData.Learnset {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func CalcHp(baseHP int, individual Individual, effort Effort) int {
	intLevel := int(DEFAULT_LEVEL)
	result := ((baseHP*2)+int(individual)+(int(effort)/4))*intLevel/100 + intLevel + 10
	return result
}

func CalcState(baseState int, individual Individual, effort Effort, natureBonus NatureBonus) int {
	result := ((baseState*2)+int(individual)+(int(effort)/4))*int(DEFAULT_LEVEL)/100 + 5
	return int(float64(result) * float64(natureBonus))
}

type Level int

const (
	DEFAULT_LEVEL = Level(50)
)

type StatusAilment string

const (
	NORMAL_POISON = StatusAilment("どく")
	BAD_POISON    = StatusAilment("もうどく")
	SLEEP         = StatusAilment("ねむり")
	BURN          = StatusAilment("やけど")
	PARALYSIS     = StatusAilment("まひ")
	FREEZE        = StatusAilment("こおり")
)

type Rank int

const (
	MIN_RANK = Rank(-6)
	MAX_RANK = Rank(6)
)

func (rank Rank) ToBonus() RankBonus {
	if rank >= 0 {
		result := (float64(rank) + 2.0) / 2.0
		return RankBonus(result)
	} else {
		abs := float64(rank) * -1
		result := 2.0 / (abs + 2.0)
		return RankBonus(result)
	}
}

type RankState struct {
	Atk   Rank
	Def   Rank
	SpAtk Rank
	SpDef Rank
	Speed Rank
}

func (rankState1 *RankState) Add(rankState2 *RankState) RankState {
	atk := rankState1.Atk + rankState2.Atk
	def := rankState1.Def + rankState2.Def
	spAtk := rankState1.SpAtk + rankState2.SpAtk
	spDef := rankState1.SpDef + rankState2.SpDef
	speed := rankState1.Speed + rankState2.Speed
	return RankState{Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}
}

func (rankState RankState) Regulate() RankState {
	if rankState.Atk > MAX_RANK {
		rankState.Atk = MAX_RANK
	}

	if rankState.Def > MAX_RANK {
		rankState.Def = MAX_RANK
	}

	if rankState.SpAtk > MAX_RANK {
		rankState.SpAtk = MAX_RANK
	}

	if rankState.SpDef > MAX_RANK {
		rankState.SpDef = MAX_RANK
	}

	if rankState.Speed > MAX_RANK {
		rankState.Speed = MAX_RANK
	}

	if rankState.Atk < MIN_RANK {
		rankState.Atk = MIN_RANK
	}

	if rankState.Def < MIN_RANK {
		rankState.Def = MIN_RANK
	}

	if rankState.SpAtk < MIN_RANK {
		rankState.SpAtk = MIN_RANK
	}

	if rankState.SpDef < MIN_RANK {
		rankState.SpDef = MIN_RANK
	}

	if rankState.Speed < MIN_RANK {
		rankState.Speed = MIN_RANK
	}
	return rankState
}

func (rankState *RankState) InDown() bool {
	if rankState.Atk < 0 {
		return true
	}

	if rankState.Def < 0 {
		return true
	}

	if rankState.SpAtk < 0 {
		return true
	}

	if rankState.SpDef < 0 {
		return true
	}

	return rankState.Speed < 0
}

func (rankState *RankState) ResetDown() RankState {
	result := RankState{}

	if rankState.Atk < 0 {
		result.Atk = 0
	} else {
		result.Atk = rankState.Atk
	}

	if rankState.Def < 0 {
		result.Def = 0
	} else {
		result.Def = rankState.Def
	}

	if rankState.SpAtk < 0 {
		result.SpAtk = 0
	} else {
		result.SpAtk = rankState.SpAtk
	}

	if rankState.SpDef < 0 {
		result.SpDef = 0
	} else {
		result.SpDef = rankState.SpDef
	}

	if rankState.Speed < 0 {
		result.Speed = 0
	} else {
		result.Speed = rankState.Speed
	}

	return result
}

type RankBonus float64

const (
	MIN_MOVESET_LENGTH = 1
	MAX_MOVESET_LENGTH = 4
)

type PointUp int

var (
	MIN_POINT_UP = PointUp(0)
	MAX_POINT_UP = PointUp(3)
)

type PointUps []PointUp

var ALL_POINT_UPS = PointUps{0, 1, 2, 3}
var ALL_POINT_UPS_LENGTH = len(ALL_POINT_UPS)

func NewMaxPointUps(length int) PointUps {
	result := make(PointUps, length)
	for i := 0; i < length; i++ {
		result[i] = MAX_POINT_UP
	}
	return result
}

func (pointUps PointUps) Index(pointUp PointUp) int {
	for i, iPointUP := range pointUps {
		if iPointUP == pointUp {
			return i
		}
	}
	return -1
}

type PowerPoint struct {
	Max     int
	Current int
}

var EMPTY_POWER_POINT = PowerPoint{Max: -1, Current: -1}

func NewPowerPoint(basePP int, pointUp PointUp) PowerPoint {
	v := (5.0 + float64(pointUp)) / 5.0
	max := int(float64(basePP) * v)
	return PowerPoint{Max: max, Current: max}
}

type PowerPoints []PowerPoint

type Moveset map[MoveName]*PowerPoint

var EMPTY_MOVESET = Moveset{EMPTY_MOVE_NAME: &EMPTY_POWER_POINT}

func NewMoveset(pokeName PokeName, moveNames MoveNames, pointUps []PointUp) (Moveset, error) {
	for _, moveName := range moveNames {
		if !InLearnset(pokeName, moveName) {
			return Moveset{}, fmt.Errorf("%v は %v を 覚えない", pokeName, moveName)
		}
	}

	if len(moveNames) != len(pointUps) {
		return Moveset{}, fmt.Errorf("len(moveName) != len(pointUps)")
	}

	powerPoints := make([]PowerPoint, len(moveNames))
	for i, moveName := range moveNames {
		basePP := MOVEDEX[moveName].BasePP
		pointUp := pointUps[i]
		if !(MIN_POINT_UP <= pointUp && MAX_POINT_UP >= pointUp) {
			return Moveset{}, fmt.Errorf("ポイントアップが、%v～%vの範囲外", MIN_POINT_UP, MAX_POINT_UP)
		}
		powerPoints[i] = NewPowerPoint(basePP, pointUps[i])
	}

	moveset := Moveset{}
	for i, moveName := range moveNames {
		moveset[moveName] = &powerPoints[i]
	}

	movesetLength := len(moveset)

	if MIN_MOVESET_LENGTH <= movesetLength && MAX_MOVESET_LENGTH >= movesetLength {
		return moveset, nil
	} else {
		return Moveset{}, fmt.Errorf("覚えさせる技の数が、%v～%vの範囲外", MIN_MOVESET_LENGTH, MAX_MOVESET_LENGTH)
	}
}

func (moveset Moveset) Keys() MoveNames {
	result := make(MoveNames, 0, len(moveset))
	for k, _ := range moveset {
		result = append(result, k)
	}
	return result
}

func (moveset Moveset) PaddingKeys() MoveNames {
	result := make(MoveNames, 0, MAX_MOVESET_LENGTH)
	for k, _ := range moveset {
		result = append(result, k)
	}

	padNum := MAX_MOVESET_LENGTH - len(moveset)
	for i := 0; i < padNum; i++ {
		result = append(result, EMPTY_MOVE_NAME)
	}
	return result
}

func (moveset Moveset) Copy() Moveset {
	result := Moveset{}
	for moveName, powerPoint := range moveset {
		copyPowerPoint := PowerPoint{Max: powerPoint.Max, Current: powerPoint.Current}
		result[moveName] = &copyPowerPoint
	}
	return result
}

func (moveset1 Moveset) Equal(moveset2 Moveset) bool {
	for moveName1, powerPoint1 := range moveset1 {
		powerPoint2, ok := moveset2[moveName1]
		if !ok {
			return false
		}
		if powerPoint1.Max != powerPoint2.Max {
			return false
		}

		if powerPoint1.Current != powerPoint2.Current {
			return false
		}
	}
	return true
}
