package bippa

import (
	"fmt"
)

const (
	MIN_MOVESET_LENGTH = 1
	MAX_MOVESET_LENGTH = 4
)

type Moveset map[MoveName]*PowerPoint

func NewMoveset(pokeName PokeName, moveNames MoveNames, pointUps []PointUp) (Moveset, error) {
	for _, pointUp := range pointUps {
		if !pointUp.IsValid() {
			errMsg := fmt.Sprintf("pointUpは%v～%vでなければならない", MIN_POINT_UP, MAX_POINT_UP)
			return Moveset{}, fmt.Errorf(errMsg)
		}
	}

	for _, moveName := range moveNames {
		if !CanLearn(pokeName, moveName) {
			errMsg := fmt.Sprintf("%vは%vを覚えない", pokeName, moveName)
			return Moveset{}, fmt.Errorf(errMsg)
		}
	}

	if len(moveNames) != len(pointUps) {
		return Moveset{}, fmt.Errorf("pointUpsのlengthは、moveNamesのlengthと一致していなければならない")
	}

	powerPoints := make([]PowerPoint, len(moveNames))
	for i, moveName := range moveNames {
		basePP := MOVEDEX[moveName].BasePP
		pointUp := pointUps[i]
		if !pointUp.IsValid() {
			errMsg := fmt.Sprintf("pointUpは%v～%vでなければならないが、%vが入力された", MIN_POINT_UP, MAX_POINT_UP, pointUp)
			return Moveset{}, fmt.Errorf(errMsg)
		}
		powerPoints[i] = NewPowerPoint(basePP, pointUps[i])
	}

	result := Moveset{}
	for i, moveName := range moveNames {
		result[moveName] = &powerPoints[i]
	}

	if !result.IsValidLength() {
		errMsg := fmt.Sprintf("movesetのlengthは%v～%vでなければならない", MIN_MOVESET_LENGTH, MAX_MOVESET_LENGTH)
		return Moveset{}, fmt.Errorf(errMsg)
	}
	return result, nil
}

func (moveset Moveset) IsValidLength() bool {
	length := len(moveset)
	return length >= MIN_MOVESET_LENGTH && length <= MAX_MOVESET_LENGTH
}

func (moveset Moveset) Copy() Moveset {
	result := Moveset{}
	for moveName, powerPoint := range moveset {
		copyPowerPoint := PowerPoint{Max:powerPoint.Max, Current:powerPoint.Current}
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
		if powerPoint1 != powerPoint2 {
			return false
		}
	}
	return true
}

func (moveset Moveset) NewMoveNames() MoveNames {
	result := make(MoveNames, 0, len(moveset))
	for moveName, _ := range moveset {
		result = append(result, moveName)
	}
	return result
}
