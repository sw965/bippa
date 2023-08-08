package bippa

import (
	"github.com/sw965/omw/fn"
)

type Effort int

var (
	EMPTY_EFFORT     = Effort(-1)
	MIN_EFFORT       = Effort(0)
	MAX_EFFORT       = Effort(252)
	MAX_SUM_EFFORT   = Effort(510)
	EFFECTIVE_EFFORT = Effort(4)
)

type Efforts []Effort

var ALL_EFFORTS = func() Efforts {
	n := int(MAX_EFFORT - MIN_EFFORT) + 1
	result := make(Efforts, n)
	for i := 0; i < n; i++ {
		result[i] = Effort(i)
	}
	return result
}()

var EFFECTIVE_EFFORTS Efforts = fn.Filter(ALL_EFFORTS, func(effrot Effort) bool { return effrot%4 == 0 } )

type EffortState struct {
	HP    Effort
	Atk   Effort
	Def   Effort
	SpAtk Effort
	SpDef Effort
	Speed Effort
}

func (es *EffortState) Sum() Effort {
	hp := es.HP
	atk := es.Atk
	def := es.Def
	spAtk := es.SpAtk
	spDef := es.SpDef
	speed := es.Speed
	return hp + atk + def + spAtk + spDef + speed
}