package bippa

type PokeName string

const (
	EMPTY_POKE_NAME = PokeName("なし")
)

func IsEmptyPokeName(pokeName PokeName) bool {
	return pokeName == EMPTY_POKE_NAME
}

type PokeNames []PokeName

var ALL_TEAM_POKE_NAMES = func() PokeNames {
	result := make(PokeNames, 0, len(ALL_POKE_NAMES)+1)
	result = append(result, ALL_POKE_NAMES...)
	result = append(result, EMPTY_POKE_NAME)
	return result
}()

func (pokeNames PokeNames) Filter(f func(PokeName) bool) PokeNames {
	result := make(PokeNames, 0, len(pokeNames))
	for _, pokeName := range pokeNames {
		if f(pokeName) {
			result = append(result, pokeName)
		}
	}
	return result
}

func (pokeNames PokeNames) Index(pokeName PokeName) int {
	for i, iPokeName := range pokeNames {
		if iPokeName == pokeName {
			return i
		}
	}
	return -1
}

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

func (pokeNames PokeNames) In(pokeName PokeName) bool {
	for _, iPokeName := range pokeNames {
		if iPokeName == pokeName {
			return true
		}
	}
	return false
}

func (pokeNames PokeNames) Sort() PokeNames {
	result := make(PokeNames, 0, len(pokeNames))
	for _, pokeName := range ALL_TEAM_POKE_NAMES {
		if pokeNames.In(pokeName) {
			for i := 0; i < pokeNames.Count(pokeName); i++ {
				result = append(result, pokeName)
			}
		}
	}
	return result
}
