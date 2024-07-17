package bippa

import (
	"fmt"
	omwmaps "github.com/sw965/omw/maps"
	"golang.org/x/exp/slices"
	omwmath "github.com/sw965/omw/math"
)

type PokeName int

const (
    EMPTY_POKE_NAME PokeName = iota
    GYARADOS   // ギャラドス
    SNORLAX    // カビゴン
    SMEARGLE   // ドーブル
    SALAMENCE  // ボーマンダ
    METAGROSS  // メタグロス
    LATIOS     // ラティオス
    EMPOLEON   // エンペルト
    BRONZONG   // ドータクン
    TOXICROAK  // ドクロッグ
)

var STRING_TO_POKE_NAME = map[string]PokeName{
    "":          EMPTY_POKE_NAME,
    "ギャラドス": GYARADOS,
    "カビゴン":   SNORLAX,
    "ドーブル":   SMEARGLE,
    "ボーマンダ": SALAMENCE,
    "メタグロス": METAGROSS,
    "ラティオス": LATIOS,
    "エンペルト": EMPOLEON,
    "ドータクン": BRONZONG,
    "ドクロッグ": TOXICROAK,
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

type PokeNamesSlice []PokeNames

type Level int

const (
	MIN_LEVEL Level = 1
	STANDARD_LEVEL Level = 50
	MAX_LEVEL Level = 100
)

type StatCalculator struct {
	BaseStat int
	Level Level
	Individual Individual
	Effort Effort
}

func (c *StatCalculator) HP() int {
	lvi := int(c.Level)
	ivi := int(c.Individual)
	evi := int(c.Effort)
	return ((c.BaseStat*2 + ivi + evi/4) * lvi / 100) + lvi + 10
}

func (c *StatCalculator) HPOther(bonus NatureBonus) int {
	lvi := int(c.Level)
	ivi := int(c.Individual)
	evi := int(c.Effort)
	stat := float64((c.BaseStat*2 + ivi + evi/4) * lvi / 100 + 5)
	return int(stat * float64(bonus))
}

type Pokemon struct {
	Name PokeName
	Level Level

	Nature Nature
	Ability Ability
	Item Item

	MoveNames MoveNames
	PointUps PointUps
	Moveset Moveset

	Individual IndividualStat
	Effort EffortStat

	MaxHP int
	CurrentHP int
	Atk int
	Def int
	SpAtk int
	SpDef int
	Speed int

	StatusAilment StatusAilment
	Rank RankStat

	IsFlinch bool

	ThisTurnCommandMoveName MoveName
	IsThisTurnActed bool
	
	SubstituteHP int
}

func NewPokemon(name PokeName, lv Level, nature Nature, ability Ability, item Item, movesetNames MoveNames, pointUps PointUps, ivStat *IndividualStat, effortStat *EffortStat) (Pokemon, error) {
	if name == EMPTY_POKE_NAME {
		return Pokemon{}, fmt.Errorf("ポケモン名が空になっている。")
	}

	pokeData := POKEDEX[name]
	if !slices.Contains(pokeData.Abilities, ability) {
		msg := fmt.Sprintf("%s は 特性：%s を持つ事は出来ない", name.ToString(), ability.ToString())
		return Pokemon{}, fmt.Errorf(msg)
	}

	moveset, err := NewMoveset(name, movesetNames)
	if err != nil {
		return Pokemon{}, err
	}

	if !effortStat.IsValidSum() {
		return Pokemon{}, GetSumEffortError(name)
	}

	ret := Pokemon{
		Name:name,
		Level:lv,
		Nature:nature,
		Ability:ability,
		Item:item,
		Individual:*ivStat,
		Effort:*effortStat,
		Moveset:moveset,
	}
	ret.UpdateStat()
	return ret, nil
}

func (p *Pokemon) UpdateStat() error {
	effortStat := p.Effort
	if !effortStat.IsValidSum() {
		return GetSumEffortError(p.Name)
	}

	if p.Name == EMPTY_POKE_NAME {
		return fmt.Errorf("ポケモン名が空になっている")
	}

	pokeData := POKEDEX[p.Name]
	lv := p.Level
	natureData := NATUREDEX[p.Nature]
	ivStat := p.Individual

	hpCalc := StatCalculator{BaseStat:pokeData.BaseHP, Level:lv, Individual:ivStat.HP, Effort:effortStat.HP}
	atkCalc := StatCalculator{BaseStat:pokeData.BaseAtk, Level:lv, Individual:ivStat.Atk, Effort:effortStat.Atk}
	defCalc := StatCalculator{BaseStat:pokeData.BaseDef, Level:lv, Individual:ivStat.Def, Effort:effortStat.Def}
	spAtkCalc := StatCalculator{BaseStat:pokeData.BaseSpAtk, Level:lv, Individual:ivStat.SpAtk, Effort:effortStat.SpAtk}
	spDefCalc := StatCalculator{BaseStat:pokeData.BaseSpDef, Level:lv, Individual:ivStat.SpDef, Effort:effortStat.SpDef}
	speedCalc := StatCalculator{BaseStat:pokeData.BaseSpeed, Level:lv, Individual:ivStat.Speed, Effort:effortStat.Speed}

	hp := hpCalc.HP()
	atk := atkCalc.HPOther(natureData.AtkBonus)
	def := defCalc.HPOther(natureData.DefBonus)
	spAtk := spAtkCalc.HPOther(natureData.SpAtkBonus)
	spDef := spDefCalc.HPOther(natureData.SpDefBonus)
	speed := speedCalc.HPOther(natureData.SpeedBonus)

	p.MaxHP = hp
	p.CurrentHP = hp
	p.Atk = atk
	p.Def = def
	p.SpAtk = spAtk
	p.SpDef = spDef
	p.Speed = speed
	return nil
}

func (p *Pokemon) AddCurrentHP(a int) error {
	if a < 0 {
		return fmt.Errorf("Pokemon.AddCurrentHPに渡す引数は、0以上でなければならない。")
	}
	p.CurrentHP += omwmath.Min(a, p.MaxHP-p.CurrentHP)
	return nil
}

func (p *Pokemon) SubCurrentHP(a int, isSubstituteHPPriority int ) error {
	if a < 0 {
		return fmt.Errorf("Pokemon.SubCurrentHPに渡す引数は、0以上でなければならない。")
	}
	p.CurrentHP -= omwmath.Min(a, p.CurrentHP)
	return nil
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

func (p *Pokemon) IsFullHP() bool {
	return p.MaxHP == p.CurrentHP
}

func (p *Pokemon) IsFainted() bool {
	return p.CurrentHP <= 0
}

func (p *Pokemon) ToEasyRead() EasyReadPokemon {
	return EasyReadPokemon{
		Name:p.Name.ToString(),
		Level:p.Level,
		Nature:p.Nature.ToString(),

		MoveNames:p.MoveNames.ToStrings(),
		PointUps:p.PointUps,
		Moveset:p.Moveset.ToEasyRead(),

		Individual:p.Individual,
		Effort:p.Effort,

		MaxHP:p.MaxHP,
		CurrentHP:p.CurrentHP,
		Atk:p.Atk,
		Def:p.Def,
		SpAtk:p.SpAtk,
		SpDef:p.SpDef,
		Speed:p.Speed,
	}
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
	for i, p := range ps {
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

func (ps Pokemons) ToPointers() PokemonPointers {
	ret := make(PokemonPointers, len(ps))
	for i := range ps {
		ret[i] = &ps[i]
	}
	return ret
}

type PokemonPointers []*Pokemon

func (ps PokemonPointers) SortBySpeed() {
	slices.SortFunc(ps, func(p1, p2 *Pokemon) bool {
		return p1.Speed > p2.Speed
	})
}

type StatusAilment int

const (
	EMPTY_STATUS_AILMENT StatusAilment = iota
	BURN //やけど
	FREEZE //こおり
	PARALYSIS //まひ
	NORMAL_POISON //どく
	BAD_POISON //もうどく
	SLEEP //ねむり
)