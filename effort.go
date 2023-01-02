package bippa

import (
	"encoding/json"
	"io/ioutil"
)

type Effort int

var (
	EMPTY_EFFORT = Effort(-1)
	MIN_EFFORT     = Effort(0)
	MAX_EFFORT     = Effort(252)
	MAX_SUM_EFFORT = Effort(510)
)

type Efforts []Effort

var ALL_VALID_EFFORTS = func() Efforts {
	length := int(MAX_EFFORT / 4) + 1
	result := make(Efforts, 0, length)
	for i := 0; i < int(MAX_EFFORT + 1); i++ {
		if i%4 == 0 {
			result = append(result, Effort(i))
		}
	}
	return result
}()

var ALL_UPPER_LIMIT_EFFORTS = func() Efforts {
	length := len(ALL_VALID_EFFORTS)
	result := make(Efforts, length)
	for i, v := range ALL_VALID_EFFORTS {
		result[i] = v + 1
	}
	return result
}()

var LOWER_LIMIT_EFFORTS = func() Efforts {
	filePath := LOWER_LIMIT_PATH + "effort.json"
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	result := Efforts{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		panic(err)
	}
	return result
}()

var UPPER_LIMIT_EFFORTS = func() Efforts {
	filePath := UPPER_LIMIT_PATH + "effort.json"
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	result := Efforts{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		panic(err)
	}
	return result
}()

func (efforts Efforts) In(effort Effort) bool {
	for _, v := range efforts {
		if v == effort {
			return true
		}
	}
	return false
}

type EffortState struct {
	HP    Effort
	Atk   Effort
	Def   Effort
	SpAtk Effort
	SpDef Effort
	Speed Effort
}

func (effortState *EffortState) Sum() Effort {
	return effortState.HP + effortState.Atk + effortState.Def + effortState.SpAtk + effortState.SpDef + effortState.Speed
}