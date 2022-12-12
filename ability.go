package bippa

type Ability string

func (ability Ability) IsValid(pokeName PokeName) bool {
	for _, iAbility := range POKEDEX[pokeName].AllAbilities {
		if iAbility == ability {
			return true
		}
	}
	return false
}

type Abilities []Ability

func (abilities Abilities) In(ability Ability) bool {
	for _, iAbility := range abilities {
		if iAbility == ability {
			return true
		}
	}
	return false
}

func (abilities Abilities) InAll(ability ...Ability) bool {
	for _, iAbility := range ability {
		if !abilities.In(iAbility) {
			return false
		}
	}
	return true
}

func (abilities Abilities) Index(ability Ability) int {
	for i, iAbility := range abilities {
		if iAbility == ability {
			return i
		}
	}
	return -1
}

type AbilityWithFloat64 map[Ability]float64

func (abilityWithFloat64 AbilityWithFloat64) KeysAndValues() (Abilities, []float64) {
	length := len(abilityWithFloat64)
	keys := make(Abilities, 0, length)
	values := make([]float64, 0, length)

	for k, v := range abilityWithFloat64 {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
