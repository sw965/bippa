package bippa

type EffortVal int

var (
	MIN_EFFORT_VAL     = EffortVal(0)
	MAX_EFFORT_VAL     = EffortVal(252)
	MAX_SUM_EFFORT_VAL = EffortVal(510)
)

type EffortVals []EffortVal

var ALL_EFFORT_VALS = func() EffortVals {
	length := int(MAX_EFFORT_VAL + 1)
	result := make(EffortVals, length)
	for i := 0; i < length; i++ {
		result[i] = EffortVal(i)
	}
	return result
}()

type Effort struct {
	HP    EffortVal
	Atk   EffortVal
	Def   EffortVal
	SpAtk EffortVal
	SpDef EffortVal
	Speed EffortVal
}

func (effort *Effort) Sum() EffortVal {
	return effort.HP + effort.Atk + effort.Def + effort.SpAtk + effort.SpDef + effort.Speed
}