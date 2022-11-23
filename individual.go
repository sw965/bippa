package bippa

type IndividualVal int

const (
	MIN_INDIVIDUAL_VAL = IndividualVal(0)
	MAX_INDIVIDUAL_VAL = IndividualVal(31)
)

func (individualVal IndividualVal) IsValid() bool {
	return (individualVal >= MIN_INDIVIDUAL_VAL) && (individualVal <= MAX_INDIVIDUAL_VAL)
}

type IndividualVals []IndividualVal

var ALL_INDIVIDUAL_VALS = func() IndividualVals {
	length := int(MAX_INDIVIDUAL_VAL + 1)
	result := make(IndividualVals, length)
	for i := 0; i < length; i++ {
		result[i] = IndividualVal(i)
	}
	return result
}()

var ALL_INDIVIDUAL_VALS_LENGTH = len(ALL_INDIVIDUAL_VALS)

type Individual struct {
	HP    IndividualVal
	Atk   IndividualVal
	Def   IndividualVal
	SpAtk IndividualVal
	SpDef IndividualVal
	Speed IndividualVal
}

var ALL_MIN_INDIVIDUAL = Individual{
	HP: MIN_INDIVIDUAL_VAL, Atk: MIN_INDIVIDUAL_VAL, Def: MIN_INDIVIDUAL_VAL,
	SpAtk: MIN_INDIVIDUAL_VAL, SpDef: MIN_INDIVIDUAL_VAL, Speed: MIN_INDIVIDUAL_VAL,
}

var ALL_MAX_INDIVIDUAL = Individual{
	HP: MAX_INDIVIDUAL_VAL, Atk: MAX_INDIVIDUAL_VAL, Def: MAX_INDIVIDUAL_VAL,
	SpAtk: MAX_INDIVIDUAL_VAL, SpDef: MAX_INDIVIDUAL_VAL, Speed: MAX_INDIVIDUAL_VAL,
}
