package bippa

type PokeName string

func (pokeName PokeName) IsValid() bool {
	_, ok := POKEDEX[pokeName]
	return ok
}

type PokeNames []PokeName

func (pokeNames PokeNames) Count(pokeName PokeName) int {
	result := 0
	for _, iPokeName := range pokeNames {
		if iPokeName == pokeName {
			result += 1
		}
	}
	return result
}

func (pokeNames PokeNames) IsUnique() bool {
	for _, pokeName := range pokeNames {
		if pokeNames.Count(pokeName) != 1 {
			return false
		}
	}
	return true
}
