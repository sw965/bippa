package bippa

type Individual int

const (
	EMPTY_INDIVIDUAL Individual = 32
	MIN_INDIVIDUAL Individual = 0
	MAX_INDIVIDUAL Individual = 31
)

type Individuals []Individual

var ALL_INDIVIDUALS = func() Individuals {
	ret := make(Individuals, int(MAX_INDIVIDUAL + 1))
	for i := range ret {
		ret[i] = Individual(i)
	}
	return ret
}()

type IndividualStat struct {
	HP Individual
	Atk Individual
	Def Individual
	SpAtk Individual
	SpDef Individual
	Speed Individual
}

var EMPTY_INDIVIDUAL_STAT = IndividualStat{
	HP:EMPTY_INDIVIDUAL,
	Atk:EMPTY_INDIVIDUAL,
	Def:EMPTY_INDIVIDUAL,
	SpAtk:EMPTY_INDIVIDUAL,
	SpDef:EMPTY_INDIVIDUAL,
	Speed:EMPTY_INDIVIDUAL,
}

var MIN_INDIVIDUAL_STAT = IndividualStat{
	HP:MIN_INDIVIDUAL,
	Atk:MIN_INDIVIDUAL,
	Def:MIN_INDIVIDUAL,
	SpAtk:MIN_INDIVIDUAL,
	SpDef:MIN_INDIVIDUAL,
	Speed:MIN_INDIVIDUAL,
}

var MAX_INDIVIDUAL_STAT = IndividualStat{
	HP:MAX_INDIVIDUAL,
	Atk:MAX_INDIVIDUAL,
	Def:MAX_INDIVIDUAL,
	SpAtk:MAX_INDIVIDUAL,
	SpDef:MAX_INDIVIDUAL,
	Speed:MAX_INDIVIDUAL,
}

func NewMaxIndividualStat() IndividualStat {
	return MAX_INDIVIDUAL_STAT.Clone()
}

func (iv *IndividualStat) Clone() IndividualStat {
	return IndividualStat{
		HP:iv.HP,
		Atk:iv.Atk,
		Def:iv.Def,
		SpAtk:iv.SpAtk,
		SpDef:iv.SpDef,
		Speed:iv.Speed,
	}
}

func (iv *IndividualStat) IsAnyEmpty() bool {
	if iv.HP == EMPTY_INDIVIDUAL {
		return true
	}

	if iv.Atk == EMPTY_INDIVIDUAL {
		return true
	}

	if iv.Def == EMPTY_INDIVIDUAL {
		return true
	}

	if iv.SpAtk == EMPTY_INDIVIDUAL {
		return true
	}

	if iv.SpDef == EMPTY_INDIVIDUAL {
		return true
	}

	if iv.Speed == EMPTY_INDIVIDUAL {
		return true
	}
	return false
}