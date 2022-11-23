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

func CalcHp(baseHP int, individualVal IndividualVal, effortVal EffortVal) int {
	intLevel := int(DEFAULT_LEVEL)
	result := ((baseHP*2)+int(individualVal)+(int(effortVal)/4))*intLevel/100 + intLevel + 10
	return result
}

func CalcState(baseState int, individualVal IndividualVal, effortVal EffortVal, natureBonus NatureBonus) int {
	result := ((baseState*2)+int(individualVal)+(int(effortVal)/4))*int(DEFAULT_LEVEL)/100 + 5
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

type RankVal int

const (
	MIN_RANK_VAL = RankVal(-6)
	MAX_RANK_VAL = RankVal(6)
)

func (rankVal RankVal) ToBonus() RankBonus {
	if rankVal >= 0 {
		result := (float64(rankVal) + 2.0) / 2.0
		return RankBonus(result)
	} else {
		abs := float64(rankVal) * -1
		result := 2.0 / (abs + 2.0)
		return RankBonus(result)
	}
}

type Rank struct {
	Atk   RankVal
	Def   RankVal
	SpAtk RankVal
	SpDef RankVal
	Speed RankVal
}

func (rank1 *Rank) Add(rank2 *Rank) Rank {
	atk := rank1.Atk + rank2.Atk
	def := rank1.Def + rank2.Def
	spAtk := rank1.SpAtk + rank2.SpAtk
	spDef := rank1.SpDef + rank2.SpDef
	speed := rank1.Speed + rank2.Speed
	return Rank{Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}
}

func (rank Rank) Regulate() Rank {
	if rank.Atk > MAX_RANK_VAL {
		rank.Atk = MAX_RANK_VAL
	}

	if rank.Def > MAX_RANK_VAL {
		rank.Def = MAX_RANK_VAL
	}

	if rank.SpAtk > MAX_RANK_VAL {
		rank.SpAtk = MAX_RANK_VAL
	}

	if rank.SpDef > MAX_RANK_VAL {
		rank.SpDef = MAX_RANK_VAL
	}

	if rank.Speed > MAX_RANK_VAL {
		rank.Speed = MAX_RANK_VAL
	}

	if rank.Atk < MIN_RANK_VAL {
		rank.Atk = MIN_RANK_VAL
	}

	if rank.Def < MIN_RANK_VAL {
		rank.Def = MIN_RANK_VAL
	}

	if rank.SpAtk < MIN_RANK_VAL {
		rank.SpAtk = MIN_RANK_VAL
	}

	if rank.SpDef < MIN_RANK_VAL {
		rank.SpDef = MIN_RANK_VAL
	}

	if rank.Speed < MIN_RANK_VAL {
		rank.Speed = MIN_RANK_VAL
	}
	return rank
}

func (rank *Rank) InDown() bool {
	if rank.Atk < 0 {
		return true
	}

	if rank.Def < 0 {
		return true
	}

	if rank.SpAtk < 0 {
		return true
	}

	if rank.SpDef < 0 {
		return true
	}

	return rank.Speed < 0
}

func (rank *Rank) ResetDown() Rank {
	result := Rank{}

	if rank.Atk < 0 {
		result.Atk = 0
	} else {
		result.Atk = rank.Atk
	}

	if rank.Def < 0 {
		result.Def = 0
	} else {
		result.Def = rank.Def
	}

	if rank.SpAtk < 0 {
		result.SpAtk = 0
	} else {
		result.SpAtk = rank.SpAtk
	}

	if rank.SpDef < 0 {
		result.SpDef = 0
	} else {
		result.SpDef = rank.SpDef
	}

	if rank.Speed < 0 {
		result.Speed = 0
	} else {
		result.Speed = rank.Speed
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
