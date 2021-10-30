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
	if !team.IsUnique() {
		return Team{}, fmt.Errorf("同じポケモンをチームに入れる事は出来ない")
	}

	if !team.IsValidLength() {
		return Team{}, fmt.Errorf("チームは%v匹～%v匹で構成されていなければならない", MIN_TEAM_LENGTH, MAX_TEAM_LENGTH)
	}

	return team, nil
}

func (team Team) PokeNames() PokeNames {
	result := make(PokeNames, len(team))
	for i, pokemon := range team {
		result[i] = pokemon.Name
	}
	return result
}

func (team Team) IsUnique() bool {
	return team.PokeNames().IsUnique()
}

func (team Team) IsValidLength() bool {
	length := len(team)
	return length >= MIN_TEAM_LENGTH && length <= MAX_TEAM_LENGTH
}

func (team Team) SelectFighters(indices *[FIGHTERS_LENGTH]int) (Fighters, error) {
	result := Fighters{}
	for i, index := range indices {
		result[i] = team[index]
	}

	if !result.IsUnique() {
		return Fighters{}, fmt.Errorf("同じポケモンを選出する事は出来ない")
	}
	return result, nil
}
