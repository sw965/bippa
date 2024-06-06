package team

import (
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
	bp "github.com/sw965/bippa"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
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
	teamV := *team
	if len(teamV) < MAX {
		ret := make(Actions, 0, len(bp.ALL_POKE_NAMES))
		currentPokeNames := team.PokeNames()
		for _, name := range bp.ALL_POKE_NAMES {
			if slices.Contains(currentPokeNames, name) {
				ret = append(ret, Action{PokeName:name})
			}
		}
		return ret
	}

	getMoveNameActions := func(pokemon *bp.Pokemon) Actions {
		if len(pokemon.Moveset) == bp.MAX_MOVESET_NUM {
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
			ret = append(ret, Action{MoveName:moveName})
		}
		return ret
	}

	getNatureActions := func(pokemon *bp.Pokemon) Actions {
		if pokemon.Nature != bp.EMPTY_NATURE {
			return Actions{}
		}
		ret := make(Actions, len(bp.ALL_NATURES))
		for i, nature := range bp.ALL_NATURES {
			ret[i] = Action{Nature:nature}
		}
		return ret
	}

	getHPIV := func(pokemon *bp.Pokemon) bp.IV {
		return pokemon.IVStat.HP
	}

	setHPIV := func(ivStat *bp.IVStat, iv bp.IV) {
		ivStat.HP = iv
	}

	getAtkIV := func(pokemon *bp.Pokemon) bp.IV {
		return pokemon.IVStat.Atk
	}

	setAtkIV := func(ivStat *bp.IVStat, iv bp.IV) {
		ivStat.Atk = iv
	}

	getDefIV := func(pokemon *bp.Pokemon) bp.IV {
		return pokemon.IVStat.Def
	}

	setDefIV := func(ivStat *bp.IVStat, iv bp.IV) {
		ivStat.Def = iv
	}

	getSpAtkIV := func(pokemon *bp.Pokemon) bp.IV {
		return pokemon.IVStat.SpAtk
	}

	setSpAtkIV := func(ivStat *bp.IVStat, iv bp.IV) {
		ivStat.SpAtk = iv
	}

	getSpDefIV := func(pokemon *bp.Pokemon) bp.IV {
		return pokemon.IVStat.SpDef
	}

	setSpDefIV := func(ivStat *bp.IVStat, iv bp.IV) {
		ivStat.SpDef = iv
	}

	getSpeedIV := func(pokemon *bp.Pokemon) bp.IV {
		return pokemon.IVStat.Speed
	}

	setSpeedIV := func(ivStat *bp.IVStat, iv bp.IV) {
		ivStat.Speed = iv
	}

	getIVActions := func(pokemon *bp.Pokemon, get func(*bp.Pokemon) bp.IV, set func(*bp.IVStat, bp.IV)) Actions {
		if get(pokemon) != bp.EMPTY_IV {
			return Actions{}
		}

		ret := make(Actions, len(bp.ALL_IVS))
		for i, iv := range bp.ALL_IVS {
			ivStat := bp.IVStat{}
			set(&ivStat, iv)
			ret[i] = Action{IVStat:ivStat}
		}
		return ret
	}

	getHPEV := func(pokemon *bp.Pokemon) bp.EV {
		return pokemon.EVStat.HP
	}

	setHPEV := func(evStat *bp.EVStat, ev bp.EV) {
		evStat.HP = ev
	}

	getAtkEV := func(pokemon *bp.Pokemon) bp.EV {
		return pokemon.EVStat.Atk
	}

	setAtkEV := func(evStat *bp.EVStat, ev bp.EV) {
		evStat.Atk = ev
	}

	getDefEV := func(pokemon *bp.Pokemon) bp.EV {
		return pokemon.EVStat.Def
	}

	setDefEV := func(evStat *bp.EVStat, ev bp.EV) {
		evStat.Def = ev
	}

	getSpAtkEV := func(pokemon *bp.Pokemon) bp.EV {
		return pokemon.EVStat.SpAtk
	}

	setSpAtkEV := func(evStat *bp.EVStat, ev bp.EV) {
		evStat.SpAtk = ev
	}

	getSpDefEV := func(pokemon *bp.Pokemon) bp.EV {
		return pokemon.EVStat.SpDef
	}

	setSpDefEV := func(evStat *bp.EVStat, ev bp.EV) {
		evStat.SpDef = ev
	}

	getSpeedEV := func(pokemon *bp.Pokemon) bp.EV {
		return pokemon.EVStat.Speed
	}

	setSpeedEV := func(evStat *bp.EVStat, ev bp.EV) {
		evStat.Speed = ev
	}

	getEVActions := func(pokemon *bp.Pokemon, get func(*bp.Pokemon) bp.EV, set func(*bp.EVStat, bp.EV)) Actions {
		if get(pokemon) != bp.EMPTY_EV {
			return Actions{}
		}

		ret := make(Actions, 0, len(bp.EFFECTIVE_EVS))
		sum := pokemon.EVStat.Sum()
		for _, ev := range bp.EFFECTIVE_EVS {
			if (sum + ev) > bp.MAX_SUM_EV {
				break
			}
			evStat := bp.EVStat{}
			set(&evStat, ev)
			ret = append(ret, Action{EVStat:evStat})
		}
		return ret
	}

	getIVFuncs := []func(*bp.Pokemon) bp.IV{
		getHPIV, getAtkIV, getDefIV, getSpAtkIV, getSpDefIV, getSpeedIV,
	}

	setIVFuncs := []func(*bp.IVStat, bp.IV){
		setHPIV, setAtkIV, setDefIV, setSpAtkIV, setSpDefIV, setSpeedIV,
	}

	getEVFuncs := []func(*bp.Pokemon) bp.EV{
		getHPEV, getAtkEV, getDefEV, getSpAtkEV, getSpDefEV, getSpeedEV,
	}

	setEVFuncs := []func(*bp.EVStat, bp.EV){
		setHPEV, setAtkEV, setDefEV, setSpAtkEV, setSpDefEV, setSpeedEV,
	}

	for _, pokemon := range teamV {
		actions := getMoveNameActions(&pokemon)
		if len(actions) != 0 {
			return actions
		}

		actions = getNatureActions(&pokemon)
		if len(actions) != 0 {
			return actions
		}

		for i, f := range getIVFuncs {
			actions = getIVActions(&pokemon, f, setIVFuncs[i])
			if len(actions) != 0 {
				return actions
			}
		}

		for i, f := range getEVFuncs {
			actions = getEVActions(&pokemon, f, setEVFuncs[i])
			if len(actions) != 0 {
				return actions
			}
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