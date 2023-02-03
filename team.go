package bippa

import (
	"fmt"
	"github.com/sw965/omw"
)

const (
	MIN_TEAM_LENGTH = 3
	MAX_TEAM_LENGTH = 6
)

type Team []Pokemon

func NewTeam(pokemons []Pokemon) (Team, error) {
	team := Team(pokemons)
	pokeNames := team.PokeNames()

	if !omw.IsUnique(pokeNames) {
		return Team{}, fmt.Errorf("同じポケモンを、同じチームに入れようとした")
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

func (team Team) Natures() Natures {
	result := make(Natures, len(team))
	for i, pokemon := range team {
		result[i] = pokemon.Nature
	}
	return result
}

func (team Team) Abilities() Abilities {
	result := make(Abilities, len(team))
	for i, pokemon := range team {
		result[i] = pokemon.Ability
	}
	return result
}

func (team Team) Items() Items {
	result := make(Items, len(team))
	for i, pokemon := range team {
		result[i] = pokemon.Item
	}
	return result
}

func (team Team) Find(pokeName PokeName) (Pokemon, error) {
	index := omw.Index(team.PokeNames(), pokeName)
	if index == -1 {
		errMsg := fmt.Sprintf("ポケモン名 : %v は チームに存在しない", pokeName)
		return Pokemon{}, fmt.Errorf(errMsg)
	} else {
		return team[index], nil
	}
}

func (team Team) Sort() (Team, error) {
	result := make(Team, 0, len(team))
	pokeNames := team.PokeNames().Sort()
	for _, pokeName := range pokeNames {
		pokemon, err := team.Find(pokeName)
		if err != nil {
			return Team{}, err
		}
		result = append(result, pokemon)
	}
	return result, nil
}
