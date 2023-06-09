package bippa

import (
	"github.com/sw965/omw"
)

type Individual int

const (
	EMPTY_INDIVIDUAL = Individual(-1)
	MIN_INDIVIDUAL   = Individual(0)
	MAX_INDIVIDUAL   = Individual(31)
)

type Individuals []Individual

var ALL_INDIVIDUALS Individuals = omw.MakeIntegerRange[Individuals](MIN_INDIVIDUAL, MAX_INDIVIDUAL+1, 1)

type IndividualState struct {
	HP    Individual
	Atk   Individual
	Def   Individual
	SpAtk Individual
	SpDef Individual
	Speed Individual
}

var ALL_MIN_INDIVIDUAL_STATE = IndividualState{
	HP: MIN_INDIVIDUAL, Atk: MIN_INDIVIDUAL, Def: MIN_INDIVIDUAL,
	SpAtk: MIN_INDIVIDUAL, SpDef: MIN_INDIVIDUAL, Speed: MIN_INDIVIDUAL,
}

var ALL_MAX_INDIVIDUAL_STATE = IndividualState{
	HP: MAX_INDIVIDUAL, Atk: MAX_INDIVIDUAL, Def: MAX_INDIVIDUAL,
	SpAtk: MAX_INDIVIDUAL, SpDef: MAX_INDIVIDUAL, Speed: MAX_INDIVIDUAL,
}

var INIT_INDIVIDUAL_STATE = IndividualState{HP: -1, Atk: -1, Def: -1, SpAtk: -1, SpDef: -1, Speed: -1}