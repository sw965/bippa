package bippa

import (
	"github.com/sw965/omw"
)

const (
	SINGLE_BATTLE_MIN_TEAM_NUM = 3
	MAX_TEAM_NUM = 6
)

type TeamBuildAction struct {
	PokeName PokeName
	MoveName MoveName
	Index int
}

type TeamBuildActions []TeamBuildAction

type Team []Pokemon

func (team Team) Clone() Team {
	clone := make(Team, len(team))
	for i := range clone {
		clone[i] = team[i].Clone()
	}
	return clone
}

func LegalTeamBuildActions(team Team) TeamBuildActions {
	//ポケモンを選ぶ行動
	if len(team) < MAX_TEAM_NUM {
		actions := make(TeamBuildActions, 0, len(ALL_POKE_NAMES))
		for _, name := range ALL_POKE_NAMES {
			actions = append(actions, TeamBuildAction{PokeName:name})
		}
		return actions
	}

	//技を選ぶ行動
	for i := range team {
		pokemon := team[i]
		learnset := POKEDEX[pokemon.Name].Learnset
		n := omw.Min(MAX_MOVESET_NUM, len(learnset))
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

func PushTeam(team Team, action *TeamBuildAction) (Team, error) {
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

