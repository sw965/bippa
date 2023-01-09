package bippa

import (
	"github.com/sw965/omw"
)

type PokeName string
type PokeNames []PokeName

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
	for _, pokeName := range ALL_POKE_NAMES {
		if pokeNames.In(pokeName) {
			for i := 0; i < pokeNames.Count(pokeName); i++ {
				result = append(result, pokeName)
			}
		}
	}
	return result
}

func (pokeNames PokeNames) Access(indices []int) PokeNames {
	result := make(PokeNames, len(indices))
	for i, index := range indices {
		result[i] = pokeNames[index]
	}
	return result
}

func (pokeNames PokeNames) Permutation(r int) ([]PokeNames, error) {
	n := len(pokeNames)
	permutationTotalNum := omw.PermutationTotalNum(n, r)

	permutationNumbers, err := omw.PermutationNumbers(n, r)
	if err != nil {
		return []PokeNames{}, err
	}

	result := make([]PokeNames, permutationTotalNum)
	for i, indices := range permutationNumbers {
		result[i] = pokeNames.Access(indices)
	}
	return result, nil
}
