package bippa

type EasyReadPokeData struct {
	Types []string
	BaseHP int
	BaseAtk int
	BaseDef int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Learnset []string
}

type EasyReadPokedex map[string]EasyReadPokeData

type EasyReadMoveData struct {
    Type        string
    Category    string
    Power       int
    Accuracy    int
    BasePP      int
}

type EasyReadMovedex map[string]EasyReadMoveData

type EasyReadDefTypeData map[string]float64
type EasyReadTypedex map[string]EasyReadDefTypeData

type EasyReadNaturedex map[string]NatureData

type EasyReadMoveset map[string]PowerPoint

func (m EasyReadMoveset) From() Moveset {
	ret := Moveset{}
	for k, v := range m {
		moveName := StringToMoveName(k)
		pp := PowerPoint{Max:v.Max, Current:v.Current}
		ret[moveName] = &pp
	}
	return ret
}

type EasyReadPokemon struct {
	Name string
	Level Level
	Nature string

	MoveNames []string
	PointUps PointUps
	Moveset EasyReadMoveset

	IVStat IVStat
	EVStat EVStat

	MaxHP int
	CurrentHP int
	Atk int
	Def int
	SpAtk int
	SpDef int
	Speed int
}

func (p *EasyReadPokemon) From() Pokemon {
	return Pokemon{
		Name:StringToPokeName(p.Name),
		Level:p.Level,
		Nature:StringToNature(p.Nature),

		MoveNames:StringsToMoveNames(p.MoveNames),
		PointUps:p.PointUps,
		Moveset:p.Moveset.From(),

		IVStat:p.IVStat,
		EVStat:p.EVStat,

		MaxHP:p.MaxHP,
		CurrentHP:p.CurrentHP,
		Atk:p.Atk,
		Def:p.Def,
		SpAtk:p.SpAtk,
		SpDef:p.SpDef,
		Speed:p.Speed,
	}
}

type EasyReadPokemons []EasyReadPokemon

func (ps EasyReadPokemons) From() Pokemons {
	ret := make(Pokemons, len(ps))
	for i, p := range ps {
		ret[i] = p.From()
	}
	return ret
}