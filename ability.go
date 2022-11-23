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

func (abilities Abilities) Index(ability Ability) int {
	for i, iAbility := range abilities {
		if iAbility == ability {
			return i
		}
	}
	return -1
}
