package bippa

import (
	"github.com/sw965/omw"
)

type MoveName string

const (
	EMPTY_MOVE_NAME = MoveName("なし")
	STRUGGLE        = MoveName("わるあがき")
)

type MoveNames []MoveName

func (moveNames MoveNames) Count(moveName MoveName) int {
	result := 0
	for _, iMoveName := range moveNames {
		if iMoveName == moveName {
			result += 1
		}
	}
	return result
}

func (moveNames MoveNames) In(moveName MoveName) bool {
	for _, iMoveName := range moveNames {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func (moveNames MoveNames) InAll(moveName ...MoveName) bool {
	for _, iMoveName := range moveName {
		if !moveNames.In(iMoveName) {
			return false
		}
	}
	return true
}

func (moveNames MoveNames) Sort() MoveNames {
	result := make(MoveNames, 0, len(moveNames))
	for _, moveName := range ALL_MOVE_NAMES {
		if moveNames.In(moveName) {
			count := moveNames.Count(moveName)
			for i := 0; i < count; i++ {
				result = append(result, moveName)
			}
		}
	}

	for i := 0; i < moveNames.Count(EMPTY_MOVE_NAME); i++ {
		result = append(result, EMPTY_MOVE_NAME)
	}
	return result
}

func (moveNames MoveNames) Index(moveName MoveName) int {
	for i, iMoveName := range moveNames {
		if iMoveName == moveName {
			return i
		}
	}
	return -1
}

func (moveNames1 MoveNames) Equal(moveNames2 MoveNames) bool {
	//index out of range 対策
	if len(moveNames1) != len(moveNames2) {
		return false
	}

	for i, moveName1 := range moveNames1 {
		moveName2 := moveNames2[i]
		if moveName1 != moveName2 {
			return false
		}
	}
	return true
}

func (moveNames MoveNames) Access(indices []int) MoveNames {
	result := make(MoveNames, len(indices))
	for i, index := range indices {
		result[i] = moveNames[index]
	}
	return result
}

func (moveNames MoveNames) Combination(r int) ([]MoveNames, error) {
	n := len(moveNames)
	combinationTotalNum, err := omw.CombinationTotalNum(n, r)
	if err != nil {
		return []MoveNames{}, err
	}
	
	combinationNumbers, err := omw.CombinationNumbers(n, r)
	if err != nil {
		return []MoveNames{}, err
	}

	result := make([]MoveNames, combinationTotalNum)
	for i, indices := range combinationNumbers {
		result[i] = moveNames.Access(indices)
	}
	return result, nil

}

type MoveNameWithFloat64 map[MoveName]float64

func (moveNameWithFloat64 MoveNameWithFloat64) KeysAndValues() (MoveNames, []float64) {
	length := len(moveNameWithFloat64)
	keys := make(MoveNames, 0, length)
	values := make([]float64, 0, length)
	for k, v := range moveNameWithFloat64 {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
