package bippa

type Individual int

const (
	EMPTY_INDIVIDUAL = Individual(-1)
	MIN_INDIVIDUAL   = Individual(0)
	MAX_INDIVIDUAL   = Individual(31)
)

type Individuals []Individual

var ALL_INDIVIDUALS = func() Individuals {
	n := int(MAX_INDIVIDUAL - MIN_INDIVIDUAL) + 1
	result := make(Individuals, n)
	for i := 0; i < n; i++ {
		result[i] = Individual(i)
	}
	return result
}()

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