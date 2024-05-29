package bippa

import (
	"fmt"
	omwmath "github.com/sw965/omw/math"
	"github.com/sw965/crow/tensor"
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	bp "github.com/sw965/bippa"
	//"math/rand"
	//"github.com/sw965/crow/game/sequential"
	//"github.com/sw965/crow/mcts/uct"
	//"github.com/sw965/crow/ucb"
)

type Action struct {
	PokeName bp.PokeName
	MoveName bp.MoveName
	Index int
}

type Actions []Action

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

func LegalActions(team Team) Actions {
	//ポケモンを選ぶ行動
	if len(team) < MAX_TEAM_NUM {
		actions := make(Actions, 0, len(ALL_POKE_NAMES))
		for _, name := range ALL_POKE_NAMES {
			actions = append(actions, Action{PokeName:name})
		}
		return actions
	}

	//技を選ぶ行動
	for i := range team {
		pokemon := team[i]
		learnset := POKEDEX[pokemon.Name].Learnset
		n := omath.Min(MAX_MOVESET_NUM, len(learnset))
		if len(pokemon.Moveset) < n {
			actions := make(TeamBuildActions, 0, len(learnset))
			for _, moveName := range learnset {
				if _, ok := pokemon.Moveset[moveName]; !ok {
					actions = append(actions, TeamBuildAction{MoveName:moveName, Index:i})
				}
				return actions
			}
		}
	}
	return TeamBuildActions{}
}

func Equal(team1, team2 Team) bool {
	team1 = team1.Sort()
	team2 = team2.Sort()
	for _, pokemon := range team1 {
		if !pokemon.Equal(team2[i]) {
			return false
		}
	}
	return true
}

func Push(team Team, action *Action) (Team, error) {
	team = team.Clone()

	if action.PokeName != EMPTY_POKE_NAME {
		pokemon := Pokemon{Name:action.PokeName, Level:DEFAULT_LEVEL, Moveset:Moveset{}}
		team = append(team, pokemon)
	}

	if action.MoveName != EMPTY_MOVE_NAME {
		basePP := MOVEDEX[action.MoveName].BasePP
		team[action.Index].Moveset[action.MoveName] = &PowerPoint{Max:basePP, Current:basePP}
	}
	return team, nil
}

func TeamEvalFeature(team Team) tensor.D1 {
	makeFeature := func(pokemon *Pokemon) tensor.D1 {
		moveFeature := make(tensor.D1, len(ALL_MOVE_NAMES))
		defFeature := make(tensor.D1, len(ALL_TYPESS))
		spDefFeature := make(tensor.D1, len(ALL_TYPESS))
		if pokemon.Name == EMPTY_POKE_NAME {
			return oslices.Concat(moveFeature, defFeature, spDefFeature)
		}

		for i, moveName := range ALL_MOVE_NAMES {
			if _, ok := pokemon.Moveset[moveName]; ok {
				moveData := MOVEDEX[moveName]
				var statFeature float64
				if moveData.Category == PHYSICS {
					statFeature = float64(pokemon.Atk) / 150.0
				} else if moveData.Category == SPECIAL {
					statFeature = float64(pokemon.Def) / 150.0
				} else {
					statFeature = 1.0
				}
				moveFeature[i] = statFeature
			}
		}

		for i, ts := range ALL_TYPESS {
			pokeTypes := POKEDEX[pokemon.Name].Types
			if slices.Equal(pokeTypes, ts) {
				defFeature[i] = float64(pokemon.Def) / 100.0
				spDefFeature[i] = float64(pokemon.SpDef) / 100.0
			}
		}
		return oslices.Concat(moveFeature, defFeature, spDefFeature)
	}

	ret := make(tensor.D1, 0, 6400)
	for _, pokemon := range team {
		ret = append(ret, makeFeature(&pokemon)...)
	}
	return ret
}