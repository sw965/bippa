package bippa

import (
	"fmt"
)

type State_ int

func NewHpState_(baseHP int, individual_ Individual_, effort_ Effort_) State_ {
	intLevel := int(MAX_LEVEL)
	result := ((baseHP*2)+int(individual_)+(int(effort_)/4))*intLevel/100 + intLevel + 10
  return State_(result)
}

func NewState_(baseState int, individual_ Individual_, effort_ Effort_, natureBonus float64) State_ {
	result := ((baseState*2)+int(individual_)+(int(effort_)/4))*int(MAX_LEVEL)/100 + 5
	return State_(float64(result) * natureBonus)
}

type State struct {
  MaxHP State_
  CurrentHP State_
  Atk State_
  Def State_
  SpAtk State_
  SpDef State_
  Speed State_
}

func NewState(pokeName PokeName, nature Nature, individual *Individual, effort *Effort, pokeData *PokeData) (State, error) {
	if !individual.IsAllValid() {
		errMsg := fmt.Sprintf("個体値は%v～%vでなければならない", MIN_INDIVIDUAL_, MAX_INDIVIDUAL_)
		return State{}, fmt.Errorf(errMsg)
	}

	if !effort.IsAllValid() {
		errMsg := fmt.Sprintf("努力値は%v～%vでなければならない", MIN_EFFORT_, MAX_EFFORT_)
		return State{}, fmt.Errorf(errMsg)
	}

	if !effort.IsValidSum() {
		errMsg := fmt.Sprintf("努力値の合計値は%vを超えてはならない", MAX_SUM_EFFORT_)
		return State{}, fmt.Errorf(errMsg)
	}

	natureData := NATUREDEX[nature]
	hp := NewHpState_(pokeData.BaseHP, individual.HP, effort.HP)
	atk := NewState_(pokeData.BaseAtk, individual.Atk, effort.Atk, natureData.AtkBonus)
	def := NewState_(pokeData.BaseDef, individual.Def, effort.Def, natureData.DefBonus)
	spAtk := NewState_(pokeData.BaseSpAtk, individual.SpAtk, effort.SpAtk, natureData.SpAtkBonus)
	spDef := NewState_(pokeData.BaseSpDef, individual.SpDef, effort.SpDef, natureData.SpDefBonus)
	speed := NewState_(pokeData.BaseSpeed, individual.Speed, effort.Speed, natureData.SpeedBonus)
	return State{MaxHP: hp, CurrentHP: hp, Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}, nil
}
