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

type EasyReadPokemon struct {
	Name string

	Level Level
	Nature string
	Moveset EasyReadMoveset
	UnassignedLearnMoveCount int

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

type EasyReadPokemons []EasyReadPokemon