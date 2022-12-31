package bippa

type Effort int

var (
	EMPTY_EFFORT = Effort(-1)
	MIN_EFFORT     = Effort(0)
	MAX_EFFORT     = Effort(252)
	MAX_SUM_EFFORT = Effort(510)
)

type Efforts []Effort

var ALL_EFFORTS = func() Efforts {
	length := int(MAX_EFFORT + 1)
	result := make(Efforts, length)
	for i := 0; i < length; i++ {
		result[i] = Effort(i)
	}
	return result
}

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