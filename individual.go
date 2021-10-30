package bippa

type Individual int

const (
	MIN_INDIVIDUAL = Individual(0)
	MAX_INDIVIDUAL = Individual(31)
)

func (individual Individual) IsValid() bool {
	return individual >= MIN_INDIVIDUAL && individual <= MAX_INDIVIDUAL
}

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

func (individualState *IndividualState) IsAllValid() bool {
	return individualState.HP.IsValid() && individualState.Atk.IsValid() && individualState.Def.IsValid() &&
		individualState.SpAtk.IsValid() && individualState.SpDef.IsValid() && individualState.Speed.IsValid()
}
