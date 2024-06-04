package team

import (
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
	bp "github.com/sw965/bippa"
)

type Action struct {
	PokeName bp.PokeName
	MoveName bp.MoveName
	Nature bp.Nature
	IVStat bp.IVStat
	EVStat bp.EVStat
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
	emptyCount := omwslices.Count(team.PokeNames(), bp.EMPTY_POKE_NAME)
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
	ret := make(Actions, 0, 1600)

	for _, name := range bp.ALL_POKE_NAMES {
		ret = append(ret, Action{PokeName:name})
	}

	teamV := *team
	getMoveNameAction := func(idx int) Actions {
		pokemon := teamV[idx]
		if len(pokemon.Moveset) == MAX_MOVESET_NUM {
			return Actions{}
		}

		learnset := bp.POKEDEX[pokemon.Name].Learnset

		moveNames = omwslices.Filter[bp.MoveNames](
			moveNames,
			func(moveName bp.MoveName) bool {
				_, ok := pokemon.Moveset[moveName]
				return !ok
			}
		)

		for _, pokemon := range pokeData.Learnset {

		}
	}

	f := func(i int) {
		pokemon := teamV[i]
		pokeData := bp.POKEDEX[pokemon.Name]
		for _, moveName := range pokeData.Learnset {
			ret = append(ret, Action{MoveName:moveName})
		}

		for _, nature := range bp.ALL_NATURES {
			ret = append(ret, Action{Nature:nature})
		}

		for _, iv := range bp.ALL_IVS {
			stat := bp.IVStat{HP:iv}
			ret = append(ret, Action{IVStat:stat})
		}

		for _, iv := range bp.ALL_IVS {
			stat := bp.IVStat{Atk:iv}
			ret = append(ret, Action{IVStat:stat})
		}

		for _, iv := range bp.ALL_IVS {
			stat := bp.IVStat{Def:iv}
			ret = append(ret, Action{IVStat:stat})
		}

		for _, iv := range bp.ALL_IVS {
			stat := bp.IVStat{SpAtk:iv}
			ret = append(ret, Action{IVStat:stat})
		}

		for _, iv := range bp.ALL_IVS {
			stat := bp.IVStat{SpDef:iv}
			ret = append(ret, Action{IVStat:stat})
		}

		for _, iv := range bp.ALL_IVS {
			stat := bp.IVStat{Speed:iv}
			ret = append(ret, Action{IVStat:stat})
		}

		for _, ev := range bp.EFFECTIVE_EVS {
			stat := bp.EVStat{HP:ev}
			ret = append(ret, Action{EVStat:stat})
		}

		for _, ev := range bp.EFFECTIVE_EVS {
			stat := bp.EVStat{Atk:ev}
			ret = append(ret, Action{EVStat:stat})
		}

		for _, ev := range bp.EFFECTIVE_EVS {
			stat := bp.EVStat{Def:ev}
			ret = append(ret, Action{EVStat:stat})
		}

		for _, ev := range bp.EFFECTIVE_EVS {
			stat := bp.EVStat{SpAtk:ev}
			ret = append(ret, Action{EVStat:stat})
		}

		for _, ev := range bp.EFFECTIVE_EVS {
			stat := bp.EVStat{SpDef:ev}
			ret = append(ret, Action{EVStat:stat})
		}

		for _, ev := range bp.EFFECTIVE_EVS {
			stat := bp.EVStat{Speed:ev}
			ret = append(ret, Action{EVStat:stat})
		}
	}

	for i := range teamV {
		f(i)
	}
	return ret
}

func Equal(team1, team2 *Team) bool {
	team1V := *team1
	team1V = team1V.Sort()
	team2V := *team2
	team2V = team2V.Sort()

	for i, pokemon := range team1V {
		if !pokemon.Equal(&team2V[i]) {
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

func IsEnd(team *Team) bool {
	teamV := *team
	if len(teamV) < MAX {
		return false
	}

	for _, pokemon := range teamV {
		if pokemon.Nature == bp.EMPTY_NATURE {
			return false
		}

		pokeData := bp.POKEDEX[pokemon.Name]
		n := omwmath.Min(len(pokeData.Learnset), bp.MAX_MOVESET_NUM)
		if len(pokemon.Moveset) < n {
			return false
		}

		if pokemon.IVStat.IsAnyEmpty() {
			return false
		}

		if pokemon.EVStat.IsAnyEmpty() {
			return false
		}
	}
	return true
}