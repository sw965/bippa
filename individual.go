package bippa

type Individual_ int

const (
	MIN_INDIVIDUAL_ = Individual_(0)
	MAX_INDIVIDUAL_ = Individual_(31)
)

func (individual_ Individual_) IsValid() bool {
	return individual_ >= MIN_INDIVIDUAL_ && individual_ <= MAX_INDIVIDUAL_
}

type Individual struct {
	HP    Individual_
	Atk   Individual_
	Def   Individual_
	SpAtk Individual_
	SpDef Individual_
	Speed Individual_
}

var ALL_MIN_INDIVIDUAL = Individual{
	HP: MIN_INDIVIDUAL_, Atk: MIN_INDIVIDUAL_, Def: MIN_INDIVIDUAL_,
	SpAtk: MIN_INDIVIDUAL_, SpDef: MIN_INDIVIDUAL_, Speed: MIN_INDIVIDUAL_,
}

var ALL_MAX_INDIVIDUAL = Individual{
	HP: MAX_INDIVIDUAL_, Atk: MAX_INDIVIDUAL_, Def: MAX_INDIVIDUAL_,
	SpAtk: MAX_INDIVIDUAL_, SpDef: MAX_INDIVIDUAL_, Speed: MAX_INDIVIDUAL_,
}

func (individual *Individual) IsAllValid() bool {
	return individual.HP.IsValid() && individual.Atk.IsValid() && individual.Def.IsValid() &&
		individual.SpAtk.IsValid() && individual.SpDef.IsValid() && individual.Speed.IsValid()
}
