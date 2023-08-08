package bippa

import (
	"golang.org/x/exp/slices"
)

const (
	MIN_TEAM_LENGTH = 3
	MAX_TEAM_LENGTH = 6
)

type Team []Pokemon

func (team Team) SelectFighters(indices []int) Fighters {
	y := Fighters{}
	for i, idx := range indices {
		y[i] = team[idx]
	}
	return y
}

func (team Team) IsValidLength() bool {
	n := len(team)
	return MIN_TEAM_LENGTH <= n && n <= MAX_TEAM_LENGTH
}

func (team Team) PokeNames() PokeNames {
	y := make(PokeNames, len(team))
	for i, poke := range team {
		y[i] = poke.Name
	}
	return y
}

func (team Team) Items() Items {
	y := make(Items, len(team))
	for i, poke := range team {
		y[i] = poke.Item
	}
	return y
}

func (team Team) Find(name PokeName) Pokemon {
	idx := slices.Index(team.PokeNames(), name)
	return team[idx]
}

func (team Team) Sort() Team {
	names := team.PokeNames()
	names.Sort()
	y := make(Team, 0, len(team))
	for _, name := range names {
		poke := team.Find(name)
		y = append(y, poke)
	}
	return y
}

func (team1 Team) Equal(team2 Team) bool {
	if len(team1) != len(team2) {
		return false
	}
	team1 = team1.Sort()
	team2 = team2.Sort()
	for i, poke := range team1 {
		if !poke.Equal(&team2[i]) {
			return false
		}
	}
	return true
}
