package bippa

import (
	"github.com/sw965/omw/fn"
)

type EasyReadPokeData struct {
	Types []string
	BaseHP int
	BaseAtk int
	BaseDef int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Abilities []string
	Learnset []string
}

func (p *EasyReadPokeData) From() (PokeData, error) {
	types, err := StringsToTypes(p.Types)
	if err != nil {
		return PokeData{}, err
	}

	abilities, err := StringsToAbilities(p.Abilities)
	if err != nil {
		return PokeData{}, err
	}

	learnset, err := StringsToMoveNames(p.Learnset)
	if err != nil {
		return PokeData{}, err
	}

	return PokeData{
		Types:types,
		BaseHP:p.BaseHP,
		BaseAtk:p.BaseAtk,
		BaseDef:p.BaseDef,
		BaseSpAtk:p.BaseSpAtk,
		BaseSpDef:p.BaseSpDef,
		BaseSpeed:p.BaseSpeed,
		Abilities:abilities,
		Learnset:learnset,
	}, nil
}

type EasyReadPokedex map[string]EasyReadPokeData

type EasyReadMoveData struct {
    Type         string
    Category     string
    Power        int
    Accuracy     int
    BasePP       int
	IsContact    bool
	PriorityRank int
	CriticalRank int
	Target       string
}

func (m *EasyReadMoveData) From() (MoveData, error) {
	t, err := StringToType(m.Type)
	if err != nil {
		return MoveData{}, err
	}

	category, err := StringToMoveCategory(m.Category)
	if err != nil {
		return MoveData{}, err
	}

	target, err := StringToTargetRange(m.Target)
	if err != nil {
		return MoveData{}, err
	}

	return MoveData{
		Type:t,
		Category:category,
		Power:m.Power,
		Accuracy:m.Accuracy,
		BasePP:m.BasePP,
		IsContact:m.IsContact,
		PriorityRank:m.PriorityRank,
		CriticalRank:m.CriticalRank,
		Target:target,
	}, nil
}

type EasyReadMovedex map[string]EasyReadMoveData

type EasyReadDefTypeData map[string]float64
type EasyReadTypedex map[string]EasyReadDefTypeData

type EasyReadNaturedex map[string]NatureData

type EasyReadMoveset map[string]PowerPoint

func (m EasyReadMoveset) From() (Moveset, error) {
	ret := Moveset{}
	for k, v := range m {
		moveName, err := StringToMoveName(k)
		if err != nil {
			return Moveset{}, err
		}
		pp := PowerPoint{Max:v.Max, Current:v.Current}
		ret[moveName] = &pp
	}
	return ret, nil
}

type EasyReadPokemon struct {
	Name string
	Level Level
	Nature string

	MoveNames []string
	PointUps PointUps
	Moveset EasyReadMoveset

	IndividualStat IndividualStat
	EffortStat EffortStat

	MaxHP int
	CurrentHP int
	Atk int
	Def int
	SpAtk int
	SpDef int
	Speed int
}

func (p *EasyReadPokemon) From() (Pokemon, error) {
	pokeName, err := StringToPokeName(p.Name)
	if err != nil {
		return Pokemon{}, err
	}

	nature, err := StringToNature(p.Nature)
	if err != nil {
		return Pokemon{}, err
	}

	moveNames, err := fn.MapWithError[MoveNames](p.MoveNames, StringToMoveName)
	if err != nil {
		return Pokemon{}, err
	}

	moveset, err := p.Moveset.From()
	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
		Name:pokeName,
		Level:p.Level,
		Nature:nature,

		MoveNames:moveNames,
		PointUps:p.PointUps,
		Moveset:moveset,

		IndividualStat:p.IndividualStat,
		EffortStat:p.EffortStat,

		MaxHP:p.MaxHP,
		CurrentHP:p.CurrentHP,
		Atk:p.Atk,
		Def:p.Def,
		SpAtk:p.SpAtk,
		SpDef:p.SpDef,
		Speed:p.Speed,
	}, nil
}

type EasyReadPokemons []EasyReadPokemon

func (ps EasyReadPokemons) From() (Pokemons, error) {
	ret := make(Pokemons, len(ps))
	for i, p := range ps {
		easy, err := p.From()
		if err != nil {
			return Pokemons{}, err
		}
		ret[i] = easy
	}
	return ret, nil
}