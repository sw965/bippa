package team

import (
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
	bp "github.com/sw965/bippa"
)

type Action struct {
	PokeName bp.PokeName
	MoveName bp.MoveName
	Index int
}

type Actions []Action

const (
	MAX = 3
)

type Team []bp.Pokemon

func (team Team) PokeNames() bp.PokeNames {
	ret := make(bp.PokeNames, len(team))
	for i, pokemon := range team {
		ret[i] = pokemon.Name
	}
	return ret
}

func (team Team) Clone() Team {
	clone := make(Team, len(team))
	for i := range clone {
		clone[i] = team[i].Clone()
	}
	return clone
}

func (team Team) Sort() Team {
	emptyCount := 0
	for _, pokemon := range team {
		if pokemon.Name == bp.EMPTY_POKE_NAME {
			emptyCount += 1
		}
	}

	ret := make(Team, 0, MAX - emptyCount)
	for _, pokeName := range bp.ALL_POKE_NAMES {
		for _, pokemon := range team {
			if pokemon.Name == pokeName {
				ret = append(ret, pokemon)
			}
		}
	}
	empty := make(Team, emptyCount)
	return omwslices.Concat(ret, empty)
}

func LegalActions(team *Team) Actions {
	teamV := *team
	//ポケモンを選ぶ行動
	if len(teamV) < MAX {
		actions := make(Actions, 0, len(bp.ALL_POKE_NAMES))
		for _, name := range bp.ALL_POKE_NAMES {
			actions = append(actions, Action{PokeName:name})
		}
		return actions
	}

	//技を選ぶ行動
	for i, pokemon := range teamV {
		learnset := bp.POKEDEX[pokemon.Name].Learnset
		n := omwmath.Min(bp.MAX_MOVESET_NUM, len(learnset))
		if len(pokemon.Moveset) < n {
			actions := make(Actions, 0, len(learnset))
			for _, moveName := range learnset {
				if _, ok := pokemon.Moveset[moveName]; !ok {
					actions = append(actions, Action{MoveName:moveName, Index:i})
				}
				return actions
			}
		}
	}
	return Actions{}
}

func Equal(team1, team2 Team) bool {
	team1 = team1.Sort()
	team2 = team2.Sort()
	for i, pokemon := range team1 {
		if !pokemon.Equal(&team2[i]) {
			return false
		}
	}
	return true
}

func Push(team Team, action *Action) (Team, error) {
	team = team.Clone()

	if action.PokeName != bp.EMPTY_POKE_NAME {
		pokemon := bp.Pokemon{Name:action.PokeName, Level:bp.DEFAULT_LEVEL, Moveset:bp.Moveset{}}
		team = append(team, pokemon)
	}

	if action.MoveName != bp.EMPTY_MOVE_NAME {
		basePP := bp.MOVEDEX[action.MoveName].BasePP
		team[action.Index].Moveset[action.MoveName] = &bp.PowerPoint{Max:basePP, Current:basePP}
	}
	return team, nil
}

func NewGame() sequential.Game {
	return sequential.Game{
		LegalActions:LegalActions,
		Equal:Equal,
		Push:Push,
		IsEnd:IsEnd,
	}
}