package bippa

import (
	"fmt"
	"math/rand"
	omwjson "github.com/sw965/omw/json"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
	omwmath "github.com/sw965/omw/math"
	omwrand "github.com/sw965/omw/math/rand"
)

type EasyReadPokemon struct {
	Name string
	Level Level
	Nature string

	MoveNames []string
	PointUps PointUps
	Moveset EasyReadMoveset

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

		Individual:p.Individual,
		Effort:p.Effort,

		MaxHP:p.MaxHP,
		CurrentHP:p.CurrentHP,
		Atk:p.Atk,
		Def:p.Def,
		SpAtk:p.SpAtk,
		SpDef:p.SpDef,
		Speed:p.Speed,

		StatusAilment:p.StatusAilment,
		Rank:p.Rank,
	}, nil
}

type EasyReadPokemons []EasyReadPokemon

func (es EasyReadPokemons) From() (Pokemons, error) {
	ps := make(Pokemons, len(es))
	for i, e := range es {
		p, err := e.From()
		if err != nil {
			return Pokemons{}, err
		}
		ps[i] = p
	}
	return ps, nil
}

type EasyReadPokedex map[string]EasyReadPokeData

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

type PokeData struct {
	Types Types
	BaseHP int
	BaseAtk int
	BaseDef int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Abilities Abilities
	Learnset MoveNames
}

func (p *PokeData) ToEasyRead() EasyReadPokeData {
	return EasyReadPokeData{
		Types:p.Types.ToStrings(),
		BaseHP:p.BaseHP,
		BaseAtk:p.BaseAtk,
		BaseDef:p.BaseDef,
		BaseSpAtk:p.BaseSpAtk,
		BaseSpDef:p.BaseSpDef,
		BaseSpeed:p.BaseSpeed,
		Abilities:p.Abilities.ToStrings(),
		Learnset:p.Learnset.ToStrings(),
	}
}

func LoadPokeData(pokeName PokeName) (PokeData, error) {
	if _, ok := POKE_NAME_TO_STRING[pokeName]; !ok {
		msg := fmt.Sprintf("%s が POKE_NAME_TO_STRING の中に存在しない", pokeName.ToString())
		return PokeData{}, fmt.Errorf(msg)
	}

	path := POKE_DATA_PATH + POKE_NAME_TO_STRING[pokeName] + omwjson.EXTENSION
	buff, err := omwjson.Load[EasyReadPokeData](path)
	if err != nil {
		return PokeData{}, err
	}
	return buff.From()
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	ret := Pokedex{}
	for _, name := range ALL_POKE_NAMES {
		data, err := LoadPokeData(name)
		if err != nil {
			panic(err)
		}
		ret[name] = &data
	}
	return ret
}()

func (p Pokedex) ToEasyRead() EasyReadPokedex {
	ret := EasyReadPokedex{}
	for k, v := range p {
		ret[k.ToString()] = v.ToEasyRead()
	}
	return ret
}

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

func (n PokeName) ToString() string {
	return POKE_NAME_TO_STRING[n]
}

type PokeNames []PokeName

var ALL_POKE_NAMES = func() PokeNames {
	buff, err := omwjson.Load[[]string](ALL_POKE_NAMES_PATH)
	if err != nil {
		panic(err)
	}

	ret, err := StringsToPokeNames(buff)
	if err != nil {
		panic(err)
	}
	return ret
}()

func (ns PokeNames) ToStrings() []string {
	ss := make([]string, len(ns))
	for i, n := range ns {
		ss[i] = n.ToString()
	}
	return ss
}

type PokeNamesSlice []PokeNames

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
	Types Types

	MaxHP int
	CurrentHP int
	Atk int
	Def int
	SpAtk int
	SpDef int
	Speed int

	StatusAilment StatusAilment
	Rank RankStat

	IsFlinchState bool
	SubstituteHP int
}

func NewPokemon(name PokeName, level Level, nature Nature, ability Ability, item Item, moveNames MoveNames, pointUps PointUps, iv *IndividualStat, ev *EffortStat) (Pokemon, error) {
	if name == EMPTY_POKE_NAME {
		return Pokemon{}, nil
	}

	p := Pokemon{
		Name:name,
		Nature:nature,
		Item:item,
		Individual:iv,
		Effort:ev,
		MoveNames:moveNames,
		PointUps:PointUps,
	}

	err := p.SetAbility(ability)
	if err != nil {
		return Pokemon{}, err
	}

	moveset, err := NewMoveset(name, moveNames)
	if err != nil {
		return Pokemon{}, err
	}

	p := Pokemon{
		Name:name,
		Level:lv,
		Nature:nature,
		Ability:ability,
		Item:item,
		Individual:*ivStat,
		Effort:*effortStat,
		MoveNames:moveNames,
		PointUps:pointUps,
		Moveset:moveset,
		Types:pokeData.Types,
	}
	err := p.UpdateStat()
	return p, err
}

func (p *Pokemon) SetAbility(a Ability) error {
	if p.Name == EMPTY_POKE_NAME {
		return fmt.Errorf("Pokemon.SetAbilityを呼び出す場合は、Pokemon.Name != EMPTY_POKE_NAME でなければなりません。")
	}

	if !slices.Contains(POKEDEX[p.Name].Abilities, a) {
		m := fmt.Sprintf("特性：%s の % s は 不適です。", a.ToString(), p.Name.ToString())
		return Pokemon{}, fmt.Errorf(m)
	}
	return nil
}

func (p *Pokemon) Put

func (p *Pokemon) UpdateStat() error {
	ev := p.Effort
	err := ev.SumError()
	if err != nil {
		return err
	}
	p.Stat = NewPokemonStat(p.Name, p.Level, p.Nature, p.Individual, p.Effort)
	return nil
}

func (p *Pokemon) AddCurrentHP(a int) error {
	if a < 0 {
		return fmt.Errorf("Pokemon.AddCurrentHPに渡す引数は、0以上でなければならない。")
	}
	p.CurrentHP += omwmath.Min(a, p.MaxHP-p.CurrentHP)
	return nil
}

func (p *Pokemon) SubCurrentHP(dmg int) error {
	if dmg < 0 {
		return fmt.Errorf("Pokemon.SubCurrentHPに渡す引数は、0以上でなければならない。")
	}
	p.CurrentHP -= omwmath.Min(dmg, p.CurrentHP)
	return nil
}

func (p *Pokemon) SubSubstituteHP(dmg int) error {
	if dmg < 0 {
		return fmt.Errorf("Pokemon.SubSubstituteHPに渡す引数は、0以上でなければならない。")
	}
	p.SubstituteHP -= omwmath.Min(dmg, p.SubstituteHP)
	return nil
}

func (p *Pokemon) SetStatusAilment(status StatusAilment, percentage int, r *rand.Rand) error {
	//https://wiki.xn--rckteqa2e.com/wiki/%E3%81%93%E3%81%8A%E3%82%8A_(%E7%8A%B6%E6%85%8B%E7%95%B0%E5%B8%B8)
	if status == FREEZE && slices.Contains(p.Types, ICE) {
		return nil
	}

	ok, err := omwrand.IsPercentageMet(percentage, r)
	if ok {
		p.StatusAilment = status
	}
	return err
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

func (p *Pokemon) IsSubstituteState() bool {
	return p.SubstituteHP > 0
}

func (p *Pokemon) RankFluctuation(rank *RankStat, percentage int, isClearBodyValid bool, r *rand.Rand) error {
	if isClearBodyValid && p.Ability == CLEAR_BODY {
		return nil
	}
	ok, err := omwrand.IsPercentageMet(percentage, r)
	if ok {
		p.Rank = p.Rank.Fluctuation(rank)
	}
	return err
}

func (p *Pokemon) SetIsFlinchState(percentage int, r *rand.Rand) error {
	ok, err := omwrand.IsPercentageMet(percentage, r)
	if ok {
		p.IsFlinchState = ok
	}
	return err
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

type PokemonEachStatCalculator struct {
	BaseStat int
	Level Level
	Individual Individual
	Effort Effort
}

func (c *PokemonEachStatCalculator) HP() int {
	lv := int(c.Level)
	i := int(c.Individual)
	e := int(c.Effort)
	return ((c.BaseStat*2 + i + e/4) * lv / 100) + lv + 10
}

func (c *PokemonEachStatCalculator) HPOther(bonus NatureBonus) int {
	lv := int(c.Level)
	i := int(c.Individual)
	e := int(c.Effort)
	stat := float64((c.BaseStat*2 + i + e/4) * lv / 100 + 5)
	return int(stat * float64(bonus))
}

type PokemonStat struct {
	MaxHP int
	CurrentHP int
	Atk int
	Def int
	SpAtk int
	SpDef int
	Speed int
}

func NewPokemonStat(name PokeName, level Level, nature Nature, iv IndividualStat, ev EffortStat) PokemonStat {
	p := POKEDEX[name]
	n := NATUREDEX[nature]
	hp := PokemonEachStatCalculator{BaseStat:p.BaseHP, Level:level, Individual:iv.HP, Effort:ev.HP}.HP()
	atk := PokemonEachStatCalculator{BaseStat:p.BaseAtk, Level:level, Individual:iv.Atk, Effort:ev.Atk}.HPOther(n.AtkBonus)
	def := PokemonEachStatCalculator{BaseStat:p.BaseDef, Level:level, Individual:iv.Def, Effort:ev.Def}.HPOther(n.DefBonus)
	spAtk := PokemonEachStatCalculator{BaseStat:p.BaseSpAtk, Level:level, Individual:iv.SpAtk, Effort:ev.SpAtk}.HPOther(n.SpAtkBonus)
	spDef := PokemonEachStatCalculator{BaseStat:p.BaseSpDef, Level:level, Individual:iv.SpDef, Effort:ev.SpDef}.HPOther(n.SpDefBonus)
	speed := PokemonEachStatCalculator{BaseStat:p.BaseSpeed, Level:level, Individual:iv.Speed, Effort:ev.Speed}.HPOther(n.SpeedBonus)
	return PokemonStat{MaxHP:hp, CurrentHP:hp, Atk:atk, Def:def, SpAtk:spAtk, SpDef:spDef, Speed:speed}
}