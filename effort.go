package bippa

import (
	"github.com/sw965/omw"
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

var ALL_EFFORTS Efforts = omw.MakeIntegerRange[Efforts](MIN_EFFORT, MAX_EFFORT+1, 1)
var EFFECTIVE_EFFORTS Efforts = omw.Filter(ALL_EFFORTS, omw.IsRemainderZero(EFFECTIVE_EFFORT))

type EffortState struct {
	HP    Effort
	Atk   Effort
	Def   Effort
	SpAtk Effort
	SpDef Effort
	Speed Effort
}

var INIT_EFFORT_STATE = EffortState{HP: -1, Atk: -1, Def: -1, SpAtk: -1, SpDef: -1, Speed: -1}

func (es *EffortState) Sum() Effort {
	hp := es.HP
	atk := es.Atk
	def := es.Def
	spAtk := es.SpAtk
	spDef := es.SpDef
	speed := es.Speed
	return hp + atk + def + spAtk + spDef + speed
}