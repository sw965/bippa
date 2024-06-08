package team

import (
	"fmt"
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
	bp "github.com/sw965/bippa"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
)

type ActionType int

const (
	EMPTY_ACTION_TYPE ActionType = iota
	POKEMON_ACTION_TYPE
	MOVE_ACTION_TYPE
	NATURE_ACTION_TYPE

	HP_IV_ACTION_TYPE
	ATK_IV_ACTION_TYPE
	DEF_IV_ACTION_TYPE
	SP_ATK_IV_ACTION_TYPE
	SP_DEF_IV_ACTION_TYPE
	SPEED_IV_ACTION_TYPE

	HP_EV_ACTION_TYPE
	ATK_EV_ACTION_TYPE
	DEF_EV_ACTION_TYPE
	SP_ATK_EV_ACTION_TYPE
	SP_DEF_EV_ACTION_TYPE
	SPEED_EV_ACTION_TYPE
)

type Action struct {
	PokeName bp.PokeName
	MoveName bp.MoveName
	Nature bp.Nature
	IV bp.IV
	EV bp.EV
	Index int
	Type ActionType
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
	teamV := *team
	if len(teamV) < MAX {
		ret := make(Actions, 0, len(bp.ALL_POKE_NAMES)+1)
		pokeNames := team.PokeNames()
		for _, name := range bp.ALL_POKE_NAMES {
			if !slices.Contains(pokeNames, name) {
				ret = append(ret, Action{PokeName:name, Type:POKEMON_ACTION_TYPE})
			}
		}
		if len(teamV) >= 3 {
			ret = append(ret, Action{PokeName:bp.EMPTY_POKE_NAME, Type:POKEMON_ACTION_TYPE})
		}
		return ret
	}

	getMoveNameActions := func(pokemon *bp.Pokemon) Actions {
		if pokemon.UnassignedLearnMoveCount == 0 {
			return Actions{}
		}

		learnset := bp.POKEDEX[pokemon.Name].Learnset
		moveNames := fn.Filter[bp.MoveNames](
			learnset,
			func(moveName bp.MoveName) bool {
				_, ok := pokemon.Moveset[moveName]
				return !ok
			},
		)

		ret := make(Actions, 0, len(moveNames)+1)
		for _, moveName := range moveNames {
			ret = append(ret, Action{MoveName:moveName, Type:MOVE_ACTION_TYPE})
		}
		if len(pokemon.Moveset) >= 1 {
			ret = append(ret, Action{MoveName:bp.EMPTY_MOVE_NAME, Type:MOVE_ACTION_TYPE})
		}
		return ret
	}

	getNatureActions := func(pokemon *bp.Pokemon) Actions {
		if pokemon.Nature != bp.EMPTY_NATURE {
			return Actions{}
		}
		ret := make(Actions, len(bp.ALL_NATURES))
		for i, nature := range bp.ALL_NATURES {
			ret[i] = Action{Nature:nature, Type:NATURE_ACTION_TYPE}
		}
		return ret
	}

	getIVActions := func(pokemon *bp.Pokemon) Actions {
		ivStat := pokemon.IVStat
		var actionType ActionType
		if ivStat.HP == bp.EMPTY_IV {
			actionType = HP_IV_ACTION_TYPE
		} else if ivStat.Atk == bp.EMPTY_IV {
			actionType = ATK_IV_ACTION_TYPE
		} else if ivStat.Def == bp.EMPTY_IV {
			actionType = DEF_IV_ACTION_TYPE
		} else if ivStat.SpAtk == bp.EMPTY_IV {
			actionType = SP_ATK_IV_ACTION_TYPE
		} else if ivStat.SpDef == bp.EMPTY_IV {
			actionType = SP_DEF_IV_ACTION_TYPE
		} else if ivStat.Speed == bp.EMPTY_IV {
			actionType = SPEED_IV_ACTION_TYPE
		} else {
			return Actions{}
		}

		ret := make(Actions, len(bp.ALL_IVS))
		for i, iv := range bp.ALL_IVS {
			ret[i] = Action{IV:iv, Type:actionType}
		}
		return ret
	}

	getEVActions := func(pokemon *bp.Pokemon) Actions {
		evStat := pokemon.EVStat
		var actionType ActionType
		if evStat.HP == bp.EMPTY_EV {
			actionType = HP_EV_ACTION_TYPE
		} else if evStat.Atk == bp.EMPTY_EV {
			actionType = ATK_EV_ACTION_TYPE
		} else if evStat.Def == bp.EMPTY_EV {
			actionType = DEF_EV_ACTION_TYPE
		} else if evStat.SpAtk == bp.EMPTY_EV {
			actionType = SP_ATK_EV_ACTION_TYPE
		} else if evStat.SpDef == bp.EMPTY_EV {
			actionType = SP_DEF_EV_ACTION_TYPE
		} else if evStat.Speed == bp.EMPTY_EV {
			actionType = SPEED_EV_ACTION_TYPE
		} else {
			return Actions{}
		}

		ret := make(Actions, 0, len(bp.EFFECTIVE_EVS))
		sum := pokemon.EVStat.Sum()
		for _, ev := range bp.EFFECTIVE_EVS {
			if (sum + ev) > bp.MAX_SUM_EV {
				break
			}
			ret = append(ret, Action{EV:ev, Type:actionType})
		}
		return ret
	}

	for _, pokemon := range teamV {
		if pokemon.Name == bp.EMPTY_POKE_NAME {
			continue
		}

		actions := getMoveNameActions(&pokemon)
		if len(actions) != 0 {
			return actions
		}

		actions = getNatureActions(&pokemon)
		if len(actions) != 0 {
			return actions
		}

		actions = getIVActions(&pokemon)
		if len(actions) != 0 {
			return actions
		}

		actions = getEVActions(&pokemon)
		if len(actions) != 0 {
			return actions
		}
	}
	return Actions{}
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
	idx := action.Index

	switch action.Type {
		case POKEMON_ACTION_TYPE:
			if action.PokeName == bp.EMPTY_POKE_NAME {
				team = append(team, bp.Pokemon{})
			} else {
				pokeData := bp.POKEDEX[action.PokeName]
				c := omwmath.Min(len(pokeData.Learnset), bp.MAX_MOVESET)
				pokemon := bp.Pokemon{
					Name:action.PokeName,
					Level:bp.DEFAULT_LEVEL,
					Moveset:bp.Moveset{},
					UnassignedLearnMoveCount:c,
					IVStat:bp.EMPTY_IV_STAT,
					EVStat:bp.EMPTY_EV_STAT,
				}
				team = append(team, pokemon)
			}
		case MOVE_ACTION_TYPE:
			team[idx].UnassignedLearnMoveCount -= 1
			if action.MoveName == bp.EMPTY_MOVE_NAME {
				team[idx].Moveset[action.MoveName] = &bp.PowerPoint{}
			} else {
				moveData := bp.MOVEDEX[action.MoveName]
				pp := bp.NewPowerPoint(moveData.BasePP, bp.MAX_POINT_UP)
				team[idx].Moveset[action.MoveName] = &pp
			}
		case NATURE_ACTION_TYPE:
			team[idx].Nature = action.Nature
		case HP_IV_ACTION_TYPE:
			team[idx].IVStat.HP = action.IV
		case ATK_IV_ACTION_TYPE:
			team[idx].IVStat.Atk = action.IV
		case DEF_IV_ACTION_TYPE:
			team[idx].IVStat.Def = action.IV
		case SP_ATK_IV_ACTION_TYPE:
			team[idx].IVStat.SpAtk = action.IV
		case SP_DEF_IV_ACTION_TYPE:
			team[idx].IVStat.SpDef = action.IV
		case SPEED_IV_ACTION_TYPE:
			team[idx].IVStat.Speed = action.IV
		case HP_EV_ACTION_TYPE:
			team[idx].EVStat.HP = action.EV
		case ATK_EV_ACTION_TYPE:
			team[idx].EVStat.Atk = action.EV
		case DEF_EV_ACTION_TYPE:
			team[idx].EVStat.Def = action.EV
		case SP_ATK_EV_ACTION_TYPE:
			team[idx].EVStat.SpAtk = action.EV
		case SP_DEF_EV_ACTION_TYPE:
			team[idx].EVStat.SpDef = action.EV
		case SPEED_EV_ACTION_TYPE:
			team[idx].EVStat.Speed = action.EV
		default:
			return Team{}, fmt.Errorf("不適切なActionType")
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

		if pokemon.UnassignedLearnMoveCount != 0 {
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