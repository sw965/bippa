package bippa

import (
	omwmaps "github.com/sw965/omw/maps"
)

type PokeName int

const (
	EMPTY_POKE_NAME PokeName = iota
	BULBASAUR
	CHARMANDER
	SQUIRTLE
	SUICUNE
	GARCHOMP
)

var STRING_TO_POKE_NAME = map[string]PokeName{
	"":EMPTY_POKE_NAME,
	"フシギダネ":BULBASAUR,
	"ヒトカゲ":CHARMANDER,
	"ゼニガメ":SQUIRTLE,
	"スイクン":SUICUNE,
	"ガブリアス":GARCHOMP,
}

var POKE_NAME_TO_STRING = omwmaps.Invert[map[PokeName]string](STRING_TO_POKE_NAME)

func (name PokeName) ToString() string {
	return POKE_NAME_TO_STRING[name]
}

type PokeNames []PokeName

func (names PokeNames) ToStrings() []string {
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = name.ToString()
	}
	return ret
}

type PokeNamess []PokeNames

type Level int

const (
	DEFAULT_LEVEL Level = 50
)

type StatCalculator struct {
	Base int
	IV IV
	EV EV
}

func (c *StatCalculator) HP() int {
	lvi := int(DEFAULT_LEVEL)
	ivi := int(c.IV)
	evi := int(c.EV)
	return ((c.Base*2 + ivi + evi/4) * lvi / 100) + lvi + 10
}

func (c *StatCalculator) HPOther(bonus NatureBonus) int {
	lvi := int(DEFAULT_LEVEL)
	ivi := int(c.IV)
	evi := int(c.EV)
	stat := float64((c.Base*2 + ivi + evi/4) * lvi / 100 + 5)
	return int(stat * float64(bonus))
}

type Pokemon struct {
	Name PokeName

	Level Level
	Nature Nature
	Moveset Moveset
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

func NewPokemon(name PokeName, nature Nature, movesetNames MoveNames, ivStat *IVStat, evStat *EVStat) (Pokemon, error) {
	pokeData := POKEDEX[name]
	natureData := NATUREDEX[nature]

	hpCalc := StatCalculator{Base:pokeData.BaseHP, IV:ivStat.HP, EV:evStat.HP}
	atkCalc := StatCalculator{Base:pokeData.BaseAtk, IV:ivStat.Atk, EV:evStat.Atk}
	defCalc := StatCalculator{Base:pokeData.BaseDef, IV:ivStat.Def, EV:evStat.Def}
	spAtkCalc := StatCalculator{Base:pokeData.BaseSpAtk, IV:ivStat.SpAtk, EV:evStat.SpAtk}
	spDefCalc := StatCalculator{Base:pokeData.BaseSpDef, IV:ivStat.SpDef, EV:evStat.SpDef}
	speedCalc := StatCalculator{Base:pokeData.BaseSpeed, IV:ivStat.Speed, EV:evStat.Speed}

	hp := hpCalc.HP()
	atk := atkCalc.HPOther(natureData.AtkBonus)
	def := defCalc.HPOther(natureData.DefBonus)
	spAtk := spAtkCalc.HPOther(natureData.SpAtkBonus)
	spDef := spDefCalc.HPOther(natureData.SpDefBonus)
	speed := speedCalc.HPOther(natureData.SpeedBonus)

	moveset, err := NewMoveset(name, movesetNames)
	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
		Name:name,
		Level:DEFAULT_LEVEL,
		Nature:nature,
		IVStat:*ivStat,
		EVStat:*evStat,
		MaxHP:hp,
		CurrentHP:hp,
		Atk:atk,
		Def:def,
		SpAtk:spAtk,
		SpDef:spDef,
		Speed:speed,
		Moveset:moveset,
	}, err
}

func (p *Pokemon) Equal(other *Pokemon) bool {
	if p.Name != other.Name {
		return false
	}

	if p.Level != other.Level {
		return false
	}

	if p.MaxHP != other.MaxHP {
		return false
	}

	if p.CurrentHP != other.CurrentHP {
		return false
	}

	if p.Atk != other.Atk {
		return false
	}

	if p.Def != other.Def {
		return false
	}

	if p.SpAtk != other.SpAtk {
		return false
	}

	if p.SpDef != other.SpDef {
		return false
	}

	if p.Speed != other.Speed {
		return false
	}
	return p.Moveset.Equal(other.Moveset)
}

func (p Pokemon) Clone() Pokemon {
	p.Moveset = p.Moveset.Clone()
	return p
}

func (p *Pokemon) HPPercentage() float64 {
	return float64(p.CurrentHP) / float64(p.MaxHP)
}

func (p *Pokemon) IsFaint() bool {
	return p.CurrentHP <= 0
}

func (p *Pokemon) ToEasyRead() EasyReadPokemon {
	return EasyReadPokemon{
		Name:p.Name.ToString(),

		Level:p.Level,
		Nature:p.Nature.ToString(),
		Moveset:p.Moveset.ToEasyRead(),
		UnassignedLearnMoveCount:p.UnassignedLearnMoveCount,

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

func NewTemplateBulbasaur() Pokemon {
	pokemon, err := NewPokemon(BULBASAUR, ADAMANT, MoveNames{TACKLE, VINE_WHIP}, &MAX_IV_STAT, &HA252_S4)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewTemplateCharmander() Pokemon {
	pokemon, err := NewPokemon(CHARMANDER, MODEST, MoveNames{EMBER}, &MAX_IV_STAT, &HC252_S4)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewTemplateSquirtle() Pokemon {
	pokemon, err := NewPokemon(SQUIRTLE, MODEST, MoveNames{WATER_GUN}, &MAX_IV_STAT, &HC252_S4)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewTemplateSuicune() Pokemon {
	pokemon, err := NewPokemon(SUICUNE, BOLD, MoveNames{SURF, ICE_BEAM}, &MAX_IV_STAT, &HB252_S4)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewTemplateGarchomp() Pokemon {
	pokemon, err := NewPokemon(GARCHOMP, JOLLY, MoveNames{STONE_EDGE}, &MAX_IV_STAT, &AS252_B4)
	if err != nil {
		panic(err)
	}
	return pokemon
}

type Pokemons []Pokemon

func (ps Pokemons) Names() PokeNames {
	ret := make(PokeNames, len(ps))
	for i, p := range ps {
		ret[i] = p.Name
	}
	return ret
}

func (ps Pokemons) Clone() Pokemons {
	ret := make(Pokemons, len(ps))
	for i, p := range ret {
		ret[i] = p.Clone()
	}
	return ret
}

func (ps Pokemons) Equal(other Pokemons) bool {
	for i, p1 := range ps {
		p2 := other[i]
		if !p1.Equal(&p2) {
			return false
		}
	}
	return true
}

func (ps Pokemons) IndexByName(name PokeName) int {
	for i, p := range ps {
		if p.Name == name {
			return i
		}
	}
	return -1
}

func (ps Pokemons) IsAllFaint() bool {
	for _, p := range ps {
		if p.CurrentHP > 0 {
			return false
		}
	}
	return true
}

func (ps Pokemons) ToEasyRead() EasyReadPokemons {
	ret := make(EasyReadPokemons, len(ps))
	for i, p := range ps {
		ret[i] = p.ToEasyRead()
	}
	return ret
}