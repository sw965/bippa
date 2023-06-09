package bippa

import (
	"github.com/sw965/omw"
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

func (team Team) LegalBuildActions() TeamBuildActions {
	pokeNames := team.PokeNames()
	if len(team.PokeNames()) < MAX_TEAM_LENGTH {
		f := func(pokeName PokeName) bool { return !slices.Contains(pokeNames, pokeName) }
		legalPokeNames := omw.Filter(ALL_POKE_NAMES, f)
		y := make(TeamBuildActions, len(legalPokeNames))
		for i, pokeName := range legalPokeNames {
			y[i] = TeamBuildAction{PokeName: pokeName}
		}
		return y
	}

	for i, pokemon := range team {
		if pokemon.Gender == "" {
			legalGenders := NewValidGenders(pokemon.Name)
			y := make(TeamBuildActions, len(legalGenders))
			for j, gender := range legalGenders {
				y[j] = TeamBuildAction{Gender: gender, Index: i}
			}
			return y
		}
	}

	for i, pokemon := range team {
		if pokemon.Nature == "" {
			y := make(TeamBuildActions, len(ALL_NATURES))
			for j, nature := range ALL_NATURES {
				y[j] = TeamBuildAction{Nature: nature, Index: i}
			}
			return y
		}
	}

	for i, pokemon := range team {
		if pokemon.Ability == "" {
			legalAbilities := POKEDEX[pokemon.Name].AllAbilities
			y := make(TeamBuildActions, len(legalAbilities))
			for j, ability := range legalAbilities {
				y[j] = TeamBuildAction{Ability: ability, Index: i}
			}
			return y
		}
	}

	for i, pokemon := range team {
		if pokemon.Item == "" {
			y := make(TeamBuildActions, 0, len(ALL_ITEMS))
			items := team.Items()
			for _, item := range ALL_ITEMS {
				if slices.Contains(items, item) {
					continue
				}
				y = append(y, TeamBuildAction{Item: item, Index: i})
			}
			return y
		}
	}

	for i, pokemon := range team {
		learnset := POKEDEX[pokemon.Name].Learnset
		n := omw.Min(len(learnset), MAX_MOVESET_LENGTH)
		if len(pokemon.Moveset) < n {
			legalMoveNames := POKEDEX[pokemon.Name].Learnset
			y := make(TeamBuildActions, len(legalMoveNames))
			for j, moveName := range legalMoveNames {
				y[j] = TeamBuildAction{MoveName: moveName, Index: i}
			}
			return y
		}
	}
	return TeamBuildActions{}
}

func (team Team) Push(action *TeamBuildAction) Team {
	y := make(Team, 0, MAX_TEAM_LENGTH)
	for _, pokemon := range team {
		y = append(y, pokemon)
	}

	if action.PokeName != "" {
		pokemon := Pokemon{Name: action.PokeName, Moveset: Moveset{}}
		pokemon.UpdateState()
		y = append(y, pokemon)
		return y
	}

	idx := action.Index
	if action.Gender != "" {
		y[idx].Gender = action.Gender
		return y
	}

	if action.Nature != "" {
		y[idx].Nature = action.Nature
		y[idx].UpdateState()
		return y
	}

	if action.Ability != "" {
		y[idx].Ability = action.Ability
		return y
	}

	if action.Item != "" {
		y[idx].Item = action.Item
		return y
	}

	if action.MoveName != "" {
		pp := NewMaxPowerPoint(action.MoveName)
		y[idx].Moveset[action.MoveName] = &pp
		return y
	}
	return Team{}
}

type TeamBuildAction struct {
	PokeName PokeName
	Gender   Gender
	Nature   Nature
	Ability  Ability
	Item     Item
	MoveName MoveName
	Index    int
}

type TeamBuildActions []TeamBuildAction
