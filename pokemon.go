package bippa

import (
	"fmt"
	"github.com/sw965/omw/fn"
	omwjson "github.com/sw965/omw/json"
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
)

type PokeData struct {
	Types     Types
	Genders   Genders
	BaseHP    int
	BaseAtk   int
	BaseDef   int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Abilities Abilities
	Learnset  MoveNames
}

func LoadPokeData(name PokeName) (PokeData, error) {
	if _, ok := POKE_NAME_TO_STRING[name]; !ok {
		msg := fmt.Sprintf("%s が POKE_NAME_TO_STRING の中に存在しない", name.ToString())
		return PokeData{}, fmt.Errorf(msg)
	}

	path := POKE_DATA_PATH + POKE_NAME_TO_STRING[name] + omwjson.EXTENSION
	buff, err := omwjson.Load[EasyReadPokeData](path)
	if err != nil {
		return PokeData{}, err
	}
	return buff.From()
}

func (p *PokeData) ToEasyRead() EasyReadPokeData {
	return EasyReadPokeData{
		Types:     p.Types.ToStrings(),
		Genders:   p.Genders.ToStrings(),
		BaseHP:    p.BaseHP,
		BaseAtk:   p.BaseAtk,
		BaseDef:   p.BaseDef,
		BaseSpAtk: p.BaseSpAtk,
		BaseSpDef: p.BaseSpDef,
		BaseSpeed: p.BaseSpeed,
		Abilities: p.Abilities.ToStrings(),
		Learnset:  p.Learnset.ToStrings(),
	}
}

type EasyReadPokeData struct {
	Types     []string
	Genders   []string
	BaseHP    int
	BaseAtk   int
	BaseDef   int
	BaseSpAtk int
	BaseSpDef int
	BaseSpeed int
	Abilities []string
	Learnset  []string
}

func (e *EasyReadPokeData) From() (PokeData, error) {
	types, err := StringsToTypes(e.Types)
	if err != nil {
		return PokeData{}, err
	}

	genders, err := StringsToGenders(e.Genders)
	if err != nil {
		return PokeData{}, err
	}

	abilities, err := StringsToAbilities(e.Abilities)
	if err != nil {
		return PokeData{}, err
	}

	learnset, err := StringsToMoveNames(e.Learnset)
	if err != nil {
		return PokeData{}, err
	}

	return PokeData{
		Types:     types,
		Genders:   genders,
		BaseHP:    e.BaseHP,
		BaseAtk:   e.BaseAtk,
		BaseDef:   e.BaseDef,
		BaseSpAtk: e.BaseSpAtk,
		BaseSpDef: e.BaseSpDef,
		BaseSpeed: e.BaseSpeed,
		Abilities: abilities,
		Learnset:  learnset,
	}, nil
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	d := Pokedex{}
	for _, name := range ALL_POKE_NAMES {
		data, err := LoadPokeData(name)
		if err != nil {
			panic(err)
		}
		d[name] = &data
	}
	return d
}()

func (p Pokedex) ToEasyRead() EasyReadPokedex {
	e := EasyReadPokedex{}
	for k, v := range p {
		e[k.ToString()] = v.ToEasyRead()
	}
	return e
}

type EasyReadPokedex map[string]EasyReadPokeData

type PokeName int

const (
	EMPTY_POKE_NAME PokeName = iota
	GYARADOS                 // ギャラドス
	SNORLAX                  // カビゴン
	SMEARGLE                 // ドーブル
	SALAMENCE                // ボーマンダ
	METAGROSS                // メタグロス
	LATIOS                   // ラティオス
	EMPOLEON                 // エンペルト
	BRONZONG                 // ドータクン
	TOXICROAK                // ドクロッグ
)

func (pn PokeName) ToString() string {
	return POKE_NAME_TO_STRING[pn]
}

type PokeNames []PokeName

var ALL_POKE_NAMES = func() PokeNames {
	buff, err := omwjson.Load[[]string](ALL_POKE_NAMES_PATH)
	if err != nil {
		panic(err)
	}

	pns, err := StringsToPokeNames(buff)
	if err != nil {
		panic(err)
	}
	return pns
}()

func (pns PokeNames) ToStrings() []string {
	ss := make([]string, len(pns))
	for i, pn := range pns {
		ss[i] = pn.ToString()
	}
	return ss
}

type PokeNamesSlice []PokeNames

type Pokemon struct {
	Name   PokeName
	Gender Gender
	Level  Level

	Nature  Nature
	Ability Ability
	Item    Item

	LearnedMoveNames MoveNames
	PointUps         PointUps
	Moveset          Moveset

	IndividualStat IndividualStat
	EffortStat     EffortStat
	Stat           PokemonStat

	Types         Types
	StatusAilment StatusAilment
	SleepTurn     int
	RankStat      RankStat

	IsFlinchState             bool
	RemainingTurnTauntState   int
	IsProtectState            bool
	ProtectConsecutiveSuccess int
	SubstituteHP              int

	//場に出てから経過したターンをカウントする。ねこだまし用。
	TurnCount int
	//繰り出そうとしている技を保持する。ふいうち用。
	ThisTurnPlannedUseMoveName MoveName

	//主に、いかくや天候特性などに用いる。
	IsAfterSwitch bool
	IsHost        bool
}

func NewPokemon(name PokeName, gender Gender, level Level, nature Nature, ability Ability, item Item, moveNames MoveNames, pointUps PointUps, iv *IndividualStat, ev *EffortStat) (Pokemon, error) {
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
		LearnedMoveNames: make(MoveNames, 0, MAX_MOVESET_LENGTH),
		PointUps:         make(PointUps, 0, MAX_MOVESET_LENGTH),
	}
	p.Name = name
	err := p.SetGender(gender)
	if err != nil {
		return Pokemon{}, err
	}
	p.Level = level
	p.Nature = nature
	err = p.SetAbility(ability)
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
	p.IndividualStat = *iv
	p.EffortStat = *ev
	p.Types = POKEDEX[name].Types
	err = p.UpdateStat()
	return p, err
}

func (p *Pokemon) SetGender(g Gender) error {
	genders := POKEDEX[p.Name].Genders
	if !slices.Contains(genders, g) {
		msg := fmt.Sprintf("性別が %s の %sは 存在しない", g.ToString(), p.Name.ToString())
		return fmt.Errorf(msg)
	}
	p.Gender = g
	return nil
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
	p.LearnedMoveNames = append(p.LearnedMoveNames, k)
	p.PointUps = append(p.PointUps, up)
	pp := NewPowerPoint(MOVEDEX[k].BasePP, up)
	p.Moveset[k] = &pp

	if len(p.Moveset) > MAX_MOVESET_LENGTH {
		return fmt.Errorf("技は最大で%d個までしか覚えられません。", MAX_MOVESET_LENGTH)
	}
	return nil
}

func (p *Pokemon) UpdateStat() error {
	ev := p.EffortStat
	err := ev.SumError()
	if err != nil {
		return err
	}
	p.Stat = NewPokemonStat(p.Name, p.Level, p.Nature, p.IndividualStat, p.EffortStat)
	return nil
}

func (p Pokemon) Clone() Pokemon {
	p.LearnedMoveNames = slices.Clone(p.LearnedMoveNames)
	p.PointUps = slices.Clone(p.PointUps)
	p.Types = slices.Clone(p.Types)
	p.Moveset = p.Moveset.Clone()
	return p
}

//後でバトルで変化する部分のみをチェックするEqual関数も作る。
func (p *Pokemon) Equal(other *Pokemon) bool {
	if p.Name != other.Name {
		return false
	}

	if p.Gender != other.Gender {
		return false
	}

	if p.Level != other.Level {
		return false
	}

	if p.Nature != other.Nature {
		return false
	}

	if p.Ability != other.Ability {
		return false
	}

	if p.Item != other.Item {
		return false
	}

	if !slices.Equal(p.LearnedMoveNames, other.LearnedMoveNames) {
		return false
	}

	if !slices.Equal(p.PointUps, other.PointUps) {
		return false
	}

	if !p.Moveset.Equal(other.Moveset) {
		return false
	}

	if p.IndividualStat != other.IndividualStat {
		return false
	}

	if p.EffortStat != other.EffortStat {
		return false
	}

	if p.Stat != other.Stat {
		return false
	}

	if !slices.Equal(p.Types, other.Types) {
		return false
	}

	if p.StatusAilment != other.StatusAilment {
		return false
	}

	if p.SleepTurn != other.SleepTurn {
		return false
	}

	if p.RankStat != other.RankStat {
		return false
	}

	if p.SubstituteHP != other.SubstituteHP {
		return false
	}

	if p.IsFlinchState != other.IsFlinchState {
		return false
	}

	if p.IsProtectState != other.IsProtectState {
		return false
	}

	if p.ProtectConsecutiveSuccess != other.ProtectConsecutiveSuccess {
		return false
	}

	if p.RemainingTurnTauntState != other.RemainingTurnTauntState {
		return false
	}

	if p.TurnCount != other.TurnCount {
		return false
	}

	if p.ThisTurnPlannedUseMoveName != other.ThisTurnPlannedUseMoveName {
		return false
	}
	return p.IsHost == other.IsHost
}

func (p *Pokemon) UsableMoveNames() MoveNames {
	mns := fn.Filter[MoveNames](p.LearnedMoveNames.FilterByNotEmpty(), func(mn MoveName) bool { return p.Moveset[mn].Current > 0 })
	if p.IsTauntState() {
		mns = fn.Filter[MoveNames](mns, func(mn MoveName) bool {
			moveData := MOVEDEX[mn]
			return moveData.Category != STATUS
		})
	}
	if len(mns) == 0 {
		mns = MoveNames{STRUGGLE}
	}
	return mns
}

func (p *Pokemon) IsFullHP() bool {
	return p.Stat.IsFullHP()
}

func (p *Pokemon) IsFainted() bool {
	return p.Stat.CurrentHP <= 0
}

func (p *Pokemon) ApplyHealToBody(heal int) error {
	if heal < 0 {
		return fmt.Errorf("回復量は0以上でなければならない")
	}

	heal = omwmath.Min(heal, p.Stat.MaxHP-p.Stat.CurrentHP)
	p.Stat.CurrentHP += heal
	return nil
}

func (p *Pokemon) IsTauntState() bool {
	return p.RemainingTurnTauntState > 0
}

func (p *Pokemon) IsSubstituteState() bool {
	return p.SubstituteHP > 0
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
		Name:   p.Name.ToString(),
		Gender: p.Gender.ToString(),
		Level:  p.Level,

		Nature:  p.Nature.ToString(),
		Ability: p.Ability.ToString(),
		Item:    p.Item.ToString(),

		LearnedMoveNames: p.LearnedMoveNames.ToStrings(),
		PointUps:         p.PointUps,
		Moveset:          p.Moveset.ToEasyRead(),

		IndividualStat: p.IndividualStat,
		EffortStat:     p.EffortStat,
		Stat:           p.Stat,

		Types:         p.Types.ToStrings(),
		StatusAilment: p.StatusAilment.ToString(),
		SleepTurn:     p.SleepTurn,
		RankStat:      p.RankStat,

		IsFlinchState:             p.IsFlinchState,
		RemainingTurnTauntState:   p.RemainingTurnTauntState,
		IsProtectState:            p.IsProtectState,
		ProtectConsecutiveSuccess: p.ProtectConsecutiveSuccess,
		SubstituteHP:              p.SubstituteHP,

		TurnCount:                  p.TurnCount,
		ThisTurnPlannedUseMoveName: p.ThisTurnPlannedUseMoveName.ToString(),

		IsAfterSwitch: p.IsAfterSwitch,
		IsHost:        p.IsHost,
	}
}

type PokemonEachStatCalculator struct {
	BaseStat   int
	Level      Level
	Individual Individual
	Effort     Effort
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
	stat := float64((c.BaseStat*2+iv+ev/4)*level/100 + 5)
	return int(stat * float64(bonus))
}

type PokemonStat struct {
	MaxHP     int
	CurrentHP int
	Atk       int
	Def       int
	SpAtk     int
	SpDef     int
	Speed     int
}

func NewPokemonStat(name PokeName, level Level, nature Nature, iv IndividualStat, ev EffortStat) PokemonStat {
	p := POKEDEX[name]
	n := NATUREDEX[nature]

	hpCalc := PokemonEachStatCalculator{BaseStat: p.BaseHP, Level: level, Individual: iv.HP, Effort: ev.HP}
	hp := hpCalc.HP()

	atkCalc := PokemonEachStatCalculator{BaseStat: p.BaseAtk, Level: level, Individual: iv.Atk, Effort: ev.Atk}
	atk := atkCalc.HPOther(n.AtkBonus)

	defCalc := PokemonEachStatCalculator{BaseStat: p.BaseDef, Level: level, Individual: iv.Def, Effort: ev.Def}
	def := defCalc.HPOther(n.DefBonus)

	spAtkCalc := PokemonEachStatCalculator{BaseStat: p.BaseSpAtk, Level: level, Individual: iv.SpAtk, Effort: ev.SpAtk}
	spAtk := spAtkCalc.HPOther(n.SpAtkBonus)

	spDefCalc := PokemonEachStatCalculator{BaseStat: p.BaseSpDef, Level: level, Individual: iv.SpDef, Effort: ev.SpDef}
	spDef := spDefCalc.HPOther(n.SpDefBonus)

	speedCalc := PokemonEachStatCalculator{BaseStat: p.BaseSpeed, Level: level, Individual: iv.Speed, Effort: ev.Speed}
	speed := speedCalc.HPOther(n.SpeedBonus)

	return PokemonStat{MaxHP: hp, CurrentHP: hp, Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}
}

func (s *PokemonStat) IsFullHP() bool {
	return s.MaxHP == s.CurrentHP
}

func (s *PokemonStat) CurrentHPRatio() float64 {
	return float64(s.CurrentHP) / float64(s.MaxHP)
}

type EasyReadPokemon struct {
	Name   string
	Gender string
	Level  Level

	Nature  string
	Ability string
	Item    string

	LearnedMoveNames []string
	PointUps         PointUps
	Moveset          EasyReadMoveset

	IndividualStat IndividualStat
	EffortStat     EffortStat
	Stat           PokemonStat

	Types         []string
	StatusAilment string
	SleepTurn     int
	RankStat      RankStat

	IsFlinchState             bool
	RemainingTurnTauntState   int
	IsProtectState            bool
	ProtectConsecutiveSuccess int
	SubstituteHP              int

	TurnCount                  int
	ThisTurnPlannedUseMoveName string

	IsAfterSwitch bool
	IsHost        bool
}

func (e *EasyReadPokemon) From() (Pokemon, error) {
	pokeName, err := StringToPokeName(e.Name)
	if err != nil {
		return Pokemon{}, err
	}

	gender, err := StringToGender(e.Gender)
	if err != nil {
		return Pokemon{}, err
	}

	nature, err := StringToNature(e.Nature)
	if err != nil {
		return Pokemon{}, err
	}

	ability, err := StringToAbility(e.Ability)
	if err != nil {
		return Pokemon{}, err
	}

	item, err := StringToItem(e.Item)
	if err != nil {
		return Pokemon{}, err
	}

	learnedMoveNames, err := fn.MapWithError[MoveNames](e.LearnedMoveNames, StringToMoveName)
	if err != nil {
		return Pokemon{}, err
	}

	moveset, err := e.Moveset.From()
	if err != nil {
		return Pokemon{}, err
	}

	types, err := StringsToTypes(e.Types)
	if err != nil {
		return Pokemon{}, err
	}

	statusAilment, err := StringToStatusAilment(e.StatusAilment)
	if err != nil {
		return Pokemon{}, err
	}

	thisTurnPlannedUseMoveName, err := StringToMoveName(e.ThisTurnPlannedUseMoveName)
	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
		Name:   pokeName,
		Gender: gender,
		Level:  e.Level,

		Nature:  nature,
		Ability: ability,
		Item:    item,

		LearnedMoveNames: learnedMoveNames,
		PointUps:         e.PointUps,
		Moveset:          moveset,

		IndividualStat: e.IndividualStat,
		EffortStat:     e.EffortStat,
		Stat:           e.Stat,

		Types:         types,
		StatusAilment: statusAilment,
		SleepTurn:     e.SleepTurn,
		RankStat:      e.RankStat,

		IsFlinchState:             e.IsFlinchState,
		RemainingTurnTauntState:   e.RemainingTurnTauntState,
		IsProtectState:            e.IsProtectState,
		ProtectConsecutiveSuccess: e.ProtectConsecutiveSuccess,
		SubstituteHP:              e.SubstituteHP,

		TurnCount:                  e.TurnCount,
		ThisTurnPlannedUseMoveName: thisTurnPlannedUseMoveName,
		IsHost:                     e.IsHost,
	}, nil
}

type Pokemons []Pokemon

func (ps Pokemons) Names() PokeNames {
	names := make(PokeNames, len(ps))
	for i, p := range ps {
		names[i] = p.Name
	}
	return names
}

func (ps Pokemons) CurrentHPs() []int {
	hps := make([]int, len(ps))
	for i, p := range ps {
		hps[i] = p.Stat.CurrentHP
	}
	return hps
}

func (ps Pokemons) Clone() Pokemons {
	c := make(Pokemons, len(ps))
	for i, p := range ps {
		c[i] = p.Clone()
	}
	return c
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

func (ps Pokemons) IsAnyFainted() bool {
	for _, p := range ps {
		if p.IsFainted() {
			return true
		}
	}
	return false
}

func (ps Pokemons) IsAllFainted() bool {
	for _, p := range ps {
		if !p.IsFainted() {
			return false
		}
	}
	return true
}

func (ps Pokemons) FaintedIndices() []int {
	return omwslices.IndicesFunc(ps, func(p Pokemon) bool { return p.IsFainted() })
}

func (ps Pokemons) NotFaintedIndices() []int {
	return omwslices.IndicesFunc(ps, func(p Pokemon) bool { return !p.IsFainted() })
}

func (ps Pokemons) SortBySpeed() {
	slices.SortFunc(ps, func(p1, p2 Pokemon) bool {
		return p1.Stat.Speed > p2.Stat.Speed
	})
}

func (ps Pokemons) ToPointers() PokemonPointers {
	pps := make(PokemonPointers, len(ps))
	for i := range ps {
		pps[i] = &ps[i]
	}
	return pps
}

func (ps Pokemons) ToEasyRead() EasyReadPokemons {
	es := make(EasyReadPokemons, len(ps))
	for i, p := range ps {
		es[i] = p.ToEasyRead()
	}
	return es
}

type PokemonPointers []*Pokemon

func (ps PokemonPointers) SortBySpeed() {
	slices.SortFunc(ps, func(p1, p2 *Pokemon) bool {
		return p1.Stat.Speed > p2.Stat.Speed
	})
}

func (ps PokemonPointers) FilterByNotFainted() PokemonPointers {
	return fn.Filter(ps, func(p *Pokemon) bool { return !p.IsFainted() })
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

//ギャラドス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Gyarados() Pokemon {
	p, err := NewPokemon(
		GYARADOS, MALE, STANDARD_LEVEL, JOLLY, INTIMIDATE, WACAN_BERRY,
		MoveNames{WATERFALL, STONE_EDGE, THUNDER_WAVE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP: 164, Atk: 92, Speed: MAX_EFFORT},
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
		SNORLAX, FEMALE, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
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
		SNORLAX, FEMALE, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
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
		SNORLAX, FEMALE, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP: 204, Atk: 52, Def: 156, SpDef: 96},
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
		SNORLAX, FEMALE, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP: 204, Atk: 52, Def: 156, SpDef: 60, Speed: 36},
	)
	if err != nil {
		panic(err)
	}
	return p
}

//ドーブル
func NewMoruhu2007Smeargle() Pokemon {
	p, err := NewPokemon(
		SMEARGLE, FEMALE, MIN_LEVEL, BRAVE, OWN_TEMPO, FOCUS_SASH,
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
		SALAMENCE, FEMALE, STANDARD_LEVEL, MODEST, INTIMIDATE, SITRUS_BERRY,
		MoveNames{DRACO_METEOR, HEAT_WAVE, RAIN_DANCE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP: 20, SpAtk: 236, Speed: 252},
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
		METAGROSS, UNKNOWN, STANDARD_LEVEL, BRAVE, CLEAR_BODY, LUM_BERRY,
		MoveNames{EARTHQUAKE, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP: MAX_EFFORT, Def: 128, SpDef: 128},
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
		METAGROSS, UNKNOWN, STANDARD_LEVEL, BRAVE, CLEAR_BODY, LUM_BERRY,
		MoveNames{HAMMER_ARM, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP: MAX_EFFORT, Def: 128, SpDef: 128},
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
		METAGROSS, UNKNOWN, STANDARD_LEVEL, ADAMANT, CLEAR_BODY, LUM_BERRY,
		MoveNames{COMET_PUNCH, BULLET_PUNCH, EARTHQUAKE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP: 236, Atk: 36, Def: 4, SpDef: 172, Speed: 60},
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
		LATIOS, MALE, STANDARD_LEVEL, TIMID, LEVITATE, FOCUS_SASH,
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
		EMPOLEON, FEMALE, STANDARD_LEVEL, MODEST, TORRENT, WACAN_BERRY,
		MoveNames{HYDRO_PUMP, SURF, ICY_WIND, PROTECT},
		MAX_POINT_UPS,
		&iv, &EffortStat{HP: 68, Def: 12, SpAtk: 252, SpDef: 4, Speed: 172},
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
		BRONZONG, UNKNOWN, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
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
		BRONZONG, UNKNOWN, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
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
		TOXICROAK, FEMALE, STANDARD_LEVEL, ADAMANT, DRY_SKIN, FOCUS_SASH,
		MoveNames{CROSS_CHOP, SUCKER_PUNCH, FAKE_OUT, TAUNT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &AS252_B4,
	)
	if err != nil {
		panic(err)
	}
	return p
}
