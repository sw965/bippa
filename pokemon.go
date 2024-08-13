package bippa

import (
	"fmt"
	omwjson "github.com/sw965/omw/json"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
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

	Stat PokemonStat

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

		Stat:p.Stat,

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

	Stat PokemonStat
	StatusAilment StatusAilment
	Rank RankStat

	SubstituteHP int
	IsFlinchState bool
	IsProtectState bool
	ProtectConsecutiveSuccess int
	RemainingTurnTauntState int

	TurnCount int
	ThisTurnPlannedUseMoveName MoveName

	Id int
}

func NewPokemon(name PokeName, level Level, nature Nature, ability Ability, item Item, moveNames MoveNames, pointUps PointUps, iv *IndividualStat, ev *EffortStat) (Pokemon, error) {
	if name == EMPTY_POKE_NAME {
		return Pokemon{}, nil
	}

	moveNamesLen := len(moveNames)
	if moveNamesLen == 0 {
		return Pokemon{}, fmt.Errorf("最低でも、%d個の技は覚えさせる必要があります。", MIN_MOVESET_LENGTH)
	}

	if moveNamesLen != len(pointUps) {
		return Pokemon{}, fmt.Errorf("覚えさせる技の数とポイントアップの数が一致しません。")
	}

	p := Pokemon{
		MoveNames:make(MoveNames, 0, MAX_MOVESET_LENGTH),
		PointUps:make(PointUps, 0, MAX_MOVESET_LENGTH),
	}
	p.Name = name
	p.Level = level
	p.Nature = nature
	err := p.SetAbility(ability)
	if err != nil {
		return Pokemon{}, err
	}
	p.Item = item
	p.Moveset = Moveset{}
	for i, moveName := range moveNames {
		err := p.SetInMoveset(moveName, pointUps[i])
		if err != nil {
			return Pokemon{}, err
		}
	}
	p.Individual = *iv
	p.Effort = *ev
	p.Types = POKEDEX[name].Types
	err = p.UpdateStat()
	return p, err
}

func (p *Pokemon) SetAbility(a Ability) error {
	if p.Name == EMPTY_POKE_NAME {
		return fmt.Errorf("Pokemon.SetAbilityを呼び出す場合は、Pokemon.Name != EMPTY_POKE_NAME でなければなりません。")
	}

	if !slices.Contains(POKEDEX[p.Name].Abilities, a) {
		m := fmt.Sprintf("特性：%s の % s は 不適です。", a.ToString(), p.Name.ToString())
		return fmt.Errorf(m)
	}
	p.Ability = a
	return nil
}

func (p *Pokemon) SetInMoveset(k MoveName, up PointUp) error {
	if !slices.Contains(POKEDEX[p.Name].Learnset, k) {
		return fmt.Errorf("%s は %s を 覚えません。", p.Name.ToString(), k.ToString())
	}
	p.MoveNames = append(p.MoveNames, k)
	p.PointUps = append(p.PointUps, up)
	pp := NewPowerPoint(MOVEDEX[k].BasePP, up)
	p.Moveset[k] = &pp

	if len(p.Moveset) > MAX_MOVESET_LENGTH {
		return fmt.Errorf("技は最大で%d個までしか覚えられません。", MAX_MOVESET_LENGTH)
	}
	return nil
}

func (p *Pokemon) UpdateStat() error {
	ev := p.Effort
	err := ev.SumError()
	if err != nil {
		return err
	}
	p.Stat = NewPokemonStat(p.Name, p.Level, p.Nature, p.Individual, p.Effort)
	return nil
}

func (p *Pokemon) Equal(other *Pokemon) bool {
	if p.Name != other.Name {
		return false
	}

	if p.Level != other.Level {
		return false
	}

	if p.Stat != other.Stat {
		return false
	}
	return p.Moveset.Equal(other.Moveset)
}

func (p Pokemon) Clone() Pokemon {
	p.Moveset = p.Moveset.Clone()
	return p
}

func (p *Pokemon) IsFullHP() bool {
	return p.Stat.IsFullHP()
}

func (p *Pokemon) IsFainted() bool {
	return p.Stat.CurrentHP <= 0
}

func (p *Pokemon) IsSubstituteState() bool {
	return p.SubstituteHP > 0
}

func (p *Pokemon) IsTauntState() bool {
	return p.RemainingTurnTauntState > 0
}

func (p *Pokemon) UsableMoveNames() MoveNames {
	ns := fn.Filter[MoveNames](p.MoveNames, func(n MoveName) bool { return p.Moveset[n].Current > 0 })
	if p.IsTauntState() {
		ns = fn.Filter[MoveNames](p.MoveNames, func(n MoveName) bool {
			moveData := MOVEDEX[n]
			return moveData.Category != STATUS
		})
	}
	if len(ns) == 0 {
		ns = MoveNames{STRUGGLE}
	}
	return ns
}

func (p *Pokemon) ApplyHealToBody(heal int) error {
	if heal < 0 {
		return fmt.Errorf("回復量は0以上でなければならない")
	}

	heal = omwmath.Min(heal, p.Stat.MaxHP - p.Stat.CurrentHP)
	p.Stat.CurrentHP += heal
	return nil
}

func (p *Pokemon) ApplyDamageToBody(dmg int) error {
	if dmg < 0 {
		return fmt.Errorf("ダメージは0以上でなければならない")
	}
	dmg = omwmath.Min(dmg, p.Stat.CurrentHP)
	p.Stat.CurrentHP -= dmg
	return nil
}

func (p *Pokemon) ApplyDamageToSubstitute(dmg int) error {
	if dmg < 0 {
		return fmt.Errorf("ダメージは0以上でなければならない")
	}
	dmg = omwmath.Min(dmg, p.SubstituteHP)
	p.SubstituteHP -= dmg
	return nil
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
		Stat:p.Stat,
	}
}

type Pokemons []Pokemon

func (ps Pokemons) IsAnyFainted() bool {
	for _, p := range ps {
		if p.IsFainted() {
			return true
		}
	}
	return false
}

func (ps Pokemons) Names() PokeNames {
	ret := make(PokeNames, len(ps))
	for i, p := range ps {
		ret[i] = p.Name
	}
	return ret
}

func (ps Pokemons) Levels() Levels {
	lvs := make(Levels, len(ps))
	for i, p := range ps {
		lvs[i] = p.Level
	}
	return lvs
}

func (ps Pokemons) MaxHPs() []int {
	hps := make([]int, len(ps))
	for i, p := range ps {
		hps[i] = p.Stat.MaxHP
	}
	return hps
}

func (ps Pokemons) CurrentHPs() []int {
	hps := make([]int, len(ps))
	for i, p := range ps {
		hps[i] = p.Stat.CurrentHP
	}
	return hps
}

func (ps Pokemons) isAnyFainted() bool {
	for _, p := range ps {
		if p.IsFainted() {
			return true
		}
	}
	return false
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

func (ps Pokemons) IsAllFainted() bool {
	for _, p := range ps {
		if !p.IsFainted() {
			return false
		}
	}
	return true
}

func (ps Pokemons) ToEasyRead() EasyReadPokemons {
	es := make(EasyReadPokemons, len(ps))
	for i, p := range ps {
		es[i] = p.ToEasyRead()
	}
	return es
}

func (ps Pokemons) NotFaintedIndices() []int {
	return omwslices.IndicesFunc(ps, func(p Pokemon) bool { return !p.IsFainted() })
}

func (ps Pokemons) ToPointers() PokemonPointers {
	pps := make(PokemonPointers, len(ps))
	for i := range ps {
		pps[i] = &ps[i]
	}
	return pps
}

func (ps Pokemons) Ids() []int {
	ids := make([]int, len(ps))
	for i, p := range ps {
		ids[i] = p.Id
	}
	return ids
}

func (ps Pokemons) ById(id int) (Pokemon, error) {
	for _, p := range ps {
		if p.Id == id {
			return p, nil
		}
	}
	return Pokemon{}, fmt.Errorf("該当するIdのポケモンが見つからなかった。")
}

type PokemonPointers []*Pokemon

func (ps PokemonPointers) SortBySpeed() {
	slices.SortFunc(ps, func(p1, p2 *Pokemon) bool {
		return p1.Stat.Speed > p2.Stat.Speed
	})
}

func (ps PokemonPointers) NotFainted() PokemonPointers {
	return fn.Filter(ps, func(p *Pokemon) bool { return !p.IsFainted() })
}

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
	level := int(c.Level)
	iv := int(c.Individual)
	ev := int(c.Effort)
	stat := float64((c.BaseStat*2 + iv + ev/4) * level / 100 + 5)
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

	hpCalc := PokemonEachStatCalculator{BaseStat:p.BaseHP, Level:level, Individual:iv.HP, Effort:ev.HP}
	hp := hpCalc.HP()

	atkCalc := PokemonEachStatCalculator{BaseStat:p.BaseAtk, Level:level, Individual:iv.Atk, Effort:ev.Atk}
	atk := atkCalc.HPOther(n.AtkBonus)

	defCalc := PokemonEachStatCalculator{BaseStat:p.BaseDef, Level:level, Individual:iv.Def, Effort:ev.Def}
	def := defCalc.HPOther(n.DefBonus)

	spAtkCalc := PokemonEachStatCalculator{BaseStat:p.BaseSpAtk, Level:level, Individual:iv.SpAtk, Effort:ev.SpAtk}
	spAtk := spAtkCalc.HPOther(n.SpAtkBonus)

	spDefCalc := PokemonEachStatCalculator{BaseStat:p.BaseSpDef, Level:level, Individual:iv.SpDef, Effort:ev.SpDef}
	spDef := spDefCalc.HPOther(n.SpDefBonus)

	speedCalc := PokemonEachStatCalculator{BaseStat:p.BaseSpeed, Level:level, Individual:iv.Speed, Effort:ev.Speed}
	speed := speedCalc.HPOther(n.SpeedBonus)

	return PokemonStat{MaxHP:hp, CurrentHP:hp, Atk:atk, Def:def, SpAtk:spAtk, SpDef:spDef, Speed:speed}
}

func (s *PokemonStat) HPPercentage() float64 {
	return float64(s.CurrentHP) / float64(s.MaxHP)
}

func (s *PokemonStat) IsFullHP() bool {
	return s.MaxHP == s.CurrentHP
}

//ギャラドス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Gyarados() Pokemon {
	p, err := NewPokemon(
		GYARADOS, STANDARD_LEVEL, JOLLY, INTIMIDATE, WACAN_BERRY,
		MoveNames{WATERFALL, STONE_EDGE, THUNDER_WAVE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP:164, Atk:92, Speed:MAX_EFFORT},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//カビゴン
func NewMoruhu2007Snorlax() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Speed = MAX_INDIVIDUAL
	p, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, FIRE_PUNCH, BELLY_DRUM, SUBSTITUTE},
		MAX_POINT_UPS,
		&iv, &HB252_D4,
	)
	if err != nil {
		panic(err)
	}
	return p
}

//カビゴン
func NewMoruhu2008Snorlax() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Speed = MIN_INDIVIDUAL
	p, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, FIRE_PUNCH, BELLY_DRUM, PROTECT},
		MAX_POINT_UPS,
		&iv, &HB252_D4,
	)
	if err != nil {
		panic(err)
	}
	return p
}

//カビゴン
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Snorlax() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Speed = MIN_INDIVIDUAL
	p, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP:204, Atk:52, Def:156, SpDef:96},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//カビゴン
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Snorlax() Pokemon {
	p, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP:204, Atk:52, Def:156, SpDef:60, Speed:36},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//ドーブル
func NewMoruhu2007Smeargle() Pokemon {
	p, err := NewPokemon(
		SMEARGLE, MIN_LEVEL, BRAVE, OWN_TEMPO, FOCUS_SASH,
		MoveNames{FAKE_OUT, FOLLOW_ME, DARK_VOID, ENDEAVOR},
		MAX_POINT_UPS,
		&MIN_INDIVIDUAL_STAT, &EffortStat{},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//ドーブル
func NewMoruhu2008Smeargle() Pokemon {
	return NewMoruhu2007Smeargle()
}

//ボーマンダ
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagiSalamence2009() Pokemon {
	p, err := NewPokemon(
		SALAMENCE, STANDARD_LEVEL, MODEST, INTIMIDATE, SITRUS_BERRY,
		MoveNames{DRACO_METEOR, HEAT_WAVE, RAIN_DANCE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:20, SpAtk:236, Speed:252},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//メタグロス
func NewMoruhu2007Metagross() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Speed = MIN_INDIVIDUAL
	p, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, BRAVE, CLEAR_BODY, LUM_BERRY,
		MoveNames{EARTHQUAKE, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP:MAX_EFFORT, Def:128, SpDef:128},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//メタグロス
func NewMoruhu2008Metagross() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Speed = MIN_INDIVIDUAL
	p, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, BRAVE, CLEAR_BODY, LUM_BERRY,
		MoveNames{HAMMER_ARM, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP:MAX_EFFORT, Def:128, SpDef:128},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//メタグロス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Metagross() Pokemon {
	p, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, ADAMANT, CLEAR_BODY, LUM_BERRY,
		MoveNames{COMET_PUNCH, BULLET_PUNCH, EARTHQUAKE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:236, Atk:36, Def:4, SpDef:172, Speed:60},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//ラティオス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Latios() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Atk = MIN_INDIVIDUAL
	p, err := NewPokemon(
		LATIOS, STANDARD_LEVEL, TIMID, LEVITATE, FOCUS_SASH,
		MoveNames{DRACO_METEOR, THUNDERBOLT, RAIN_DANCE, PROTECT},
		MAX_POINT_UPS,
		&iv, &CS252_H4,
	)
	if err != nil {
		panic(err)
	}
	return p
}

//エンペルト
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Empoleon() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Atk = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		EMPOLEON, STANDARD_LEVEL, MODEST, TORRENT, WACAN_BERRY,
		MoveNames{HYDRO_PUMP, SURF, ICY_WIND, PROTECT},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP:68, Def:12, SpAtk:252, SpDef:4, Speed:172},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドータクン
func NewMoruhu2007Bronzong() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Speed = MIN_INDIVIDUAL
	p, err := NewPokemon(
		BRONZONG, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
		MoveNames{GYRO_BALL, EXPLOSION, TRICK_ROOM, HYPNOSIS},
		MAX_POINT_UPS,
		&iv, &HD252_B4,
	)
	if err != nil {
		panic(err)
	}
	return p
}

//ドータクン
func NewMoruhu2008Bronzong() Pokemon {
	iv := MAX_INDIVIDUAL_STAT.Clone()
	iv.Atk = MIN_INDIVIDUAL
	iv.Speed = MIN_INDIVIDUAL
	p, err := NewPokemon(
		BRONZONG, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
		MoveNames{PSYCHIC, EXPLOSION, TRICK_ROOM, HYPNOSIS},
		MAX_POINT_UPS,
		&iv, &HD252_B4,
	)
	if err != nil {
		panic(err)
	}
	return p
}

//ドクロッグ
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Toxicroak() Pokemon {
	p, err := NewPokemon(
		TOXICROAK, STANDARD_LEVEL, ADAMANT, DRY_SKIN, FOCUS_SASH,
		MoveNames{CROSS_CHOP, SUCKER_PUNCH, FAKE_OUT, TAUNT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &AS252_B4,
	)
	if err != nil {
		panic(err)
	}
	return p
}