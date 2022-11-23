package bippa

import (
	"fmt"
)

const (
	MIN_TEAM_LENGTH = 3
	MAX_TEAM_LENGTH = 6
)

type Team []Pokemon

func NewTeam(pokemons []Pokemon) (Team, error) {
	team := Team(pokemons)
	pokeNames := team.PokeNames()

	if !pokeNames.Filter(IsEmptyPokeName).IsUnique() {
		return Team{}, fmt.Errorf("同じポケモンを、チームに入れる事は出来ない")
	}

	if team.IsValidLength() {
		return team, nil
	} else {
		return Team{}, fmt.Errorf("チームは、%v匹～%v匹で構成されていなければならない", MIN_TEAM_LENGTH, MAX_TEAM_LENGTH)
	}
}

func (team Team) IsValidLength() bool {
	length := len(team)
	return MIN_TEAM_LENGTH <= length && length <= MAX_TEAM_LENGTH
}

func (team Team) PokeNames() PokeNames {
	result := make(PokeNames, len(team))
	for i, pokemon := range team {
		result[i] = pokemon.Name
	}
	return result
}

func (team Team) Find(pokeName PokeName) (Pokemon, error) {
	if pokeName == EMPTY_POKE_NAME {
		return Pokemon{}, fmt.Errorf("Team.Find関数において EMPTY_POKE_NAME は 引数として不適")
	}

	index := team.PokeNames().Index(pokeName)
	if index == -1 {
		errMsg := fmt.Sprintf("ポケモン名 : %v は チームに存在しない", pokeName)
		return Pokemon{}, fmt.Errorf(errMsg)
	} else {
		return team[index], nil
	}
}

func (team Team) Sort() Team {
	result := make(Team, 0, len(team))
	pokeNames := team.PokeNames().Sort()
	for _, pokeName := range pokeNames {
		if pokeName == EMPTY_POKE_NAME {
			result = append(result, NewEmptyPokemon())
		} else {
			pokemon, err := team.Find(pokeName)
			if err != nil {
				panic(err)
			}
			result = append(result, pokemon)
		}
	}
	return result
}
