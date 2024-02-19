package bippa

import (
	"fmt"
	"golang.org/x/exp/slices"
	omwmath "github.com/sw965/omw/math"
)

type PokeName int

const (
	NO_POKE_NAME PokeName = iota
	NYA_O_HA
)

func StringToPokeName(s string) PokeName {
	return STRING_TO_POKE_NAME[s]
}

type PokeNames []PokeName

func (names PokeNames) Sort() {
	isSwap := func(name1, name2 PokeName) bool {
		return slices.Index(ALL_POKE_NAMES, name1) > slices.Index(ALL_POKE_NAMES, name2)
	}
	slices.SortFunc(names, isSwap)
}

type Gender int

const (
	MALE Gender = iota
	FEMALE
	UNKNOWN
)

func NewGender(s string) (Gender, error) {
	switch s {
		case "♂":
			return MALE, nil
		case "♀":
			return FEMALE, nil
		case "不明":
			return UNKNOWN, nil
		default:
			return -1, fmt.Errorf("不適なgender")
	}
}

type Genders []Gender

var ALL_GENDERS = Genders{MALE, FEMALE, UNKNOWN}

func NewGenders(ss []string) (Genders, error) {
	return fn.MapError[Genders](ss, NewGender)
}

type Nature string
type Natures []Nature
type NatureBonus float64

const (
	NO_NATURE_BONUS   = NatureBonus(1.0)
	UP_NATURE_BONUS   = NatureBonus(1.1)
	DOWN_NATURE_BONUS = NatureBonus(0.9)
)

type NatureBonuses []NatureBonus

var ALL_NATURE_BONUSES = NatureBonuses{NO_NATURE_BONUS, UP_NATURE_BONUS, DOWN_NATURE_BONUS}

type Item int

const (
	NO_ITEM Item = iota
	LIFE_ORB //いのちのたま
	SITRUS_BERRY //オボンのみ
	BLACK_SLUDGE //くろいヘドロ
	CHOICE_SCARF //こだわりスカーフ
	CHOICE_BAND //こだわりハチマキ
	CHOICE_SPECS //こだわりメガネ
	ROCKY_HELMET //ゴツゴツメット
	WHITE_HERB //しろいハーブ
	LEFTOVERS //たべのこし
	ASSAULT_VEST //とつげきチョッキ
)

func NewItem(s string) (Item, error) {
	switch s {
		case "なし":
			return NO_ITEM, nil
		case "いのちのたま":
			return LIFE_ORB, nil
		case "オボンのみ":
			return SITRUS_BERRY, nil
		case "くろいヘドロ":
			return BLACK_SLUDGE, nil
		case "こだわりスカーフ":
			return CHOICE_SPECS, nil
		case "こだわりハチマキ":
			return CHOICE_BAND, nil
		case "こだわりメガネ":
			return CHOICE_SPECS, nil
		case "ゴツゴツメット":
			return ROCKY_HELMET, nil
		case "しろいハーブ":
			return WHITE_HERB, nil
		case "たべのこし":
			return LEFTOVERS, nil
		case "とつげきチョッキ":
			return ASSAULT_VEST, nil
		default:
			return -1, fmt.Errorf("不適なitem")
	}
}

func (item Item) IsChoice() bool {
	return slices.Contains(CHOICE_ITEMS, item)
}

func LifeOrbEffect(bt *Battle) {
	dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 10.0)
	bt.Damage(dmg)
}

type Items []Item

var ALL_ITEMS = func() Items {
	d, err := omwjson.Load[[]string](ALL_ITEMS_PATH)
	if err != nil {
		panic(err)
	}
	y, err := fn.MapError[Items](d, NewItem)
	if err != nil {
		panic(err)
	}
	return y
}()

var CHOICE_ITEMS = Items{CHOICE_BAND, CHOICE_SPECS, CHOICE_SCARF}

type IndividualVal int

const (
	MIN_INDIVIDUAL_VAL   = IndividualVal(0)
	MAX_INDIVIDUAL_VAL   = IndividualVal(31)
)

type IndividualVals []IndividualVal

var ALL_INDIVIDUAL_VALS = func() IndividualVals {
	n := int(MAX_INDIVIDUAL_VAL - MIN_INDIVIDUAL_VAL) + 1
	y := make(IndividualVals, n)
	for i := 0; i < n; i++ {
		y[i] = IndividualVal(i)
	}
	return y
}()

type Individual struct {
	HP    IndividualVal
	Atk   IndividualVal
	Def   IndividualVal
	SpAtk IndividualVal
	SpDef IndividualVal
	Speed IndividualVal
}

var ALL_MIN_INDIVIDUAL = Individual{
	HP: MIN_INDIVIDUAL_VAL, Atk: MIN_INDIVIDUAL_VAL, Def: MIN_INDIVIDUAL_VAL,
	SpAtk: MIN_INDIVIDUAL_VAL, SpDef: MIN_INDIVIDUAL_VAL, Speed: MIN_INDIVIDUAL_VAL,
}

var ALL_MAX_INDIVIDUAL = Individual{
	HP: MAX_INDIVIDUAL_VAL, Atk: MAX_INDIVIDUAL_VAL, Def: MAX_INDIVIDUAL_VAL,
	SpAtk: MAX_INDIVIDUAL_VAL, SpDef: MAX_INDIVIDUAL_VAL, Speed: MAX_INDIVIDUAL_VAL,
}

type EffortVal int

var (
	MIN_EFFORT_VAL       = Effort(0)
	MAX_EFFORT_VAL       = Effort(252)
	MAX_SUM_EFFORT_VAL   = Effort(510)
	EFFECTIVE_EFFORT_VAL = Effort(4)
)

type EffortVals []EffortVal

var ALL_EFFORT_VALS = func() EffortVals {
	n := int(MAX_EFFORT_VAL - MIN_EFFORT_VAL) + 1
	y := make(Efforts, n)
	for i := 0; i < n; i++ {
		y[i] = Effort(i)
	}
	return y
}()

var EFFECTIVE_EFFORT_VALS Efforts = fn.Filter(ALL_EFFORT_VALS, func(v EffortVal) bool { return v%4 == 0 } )

type Effort struct {
	HP    EffortVal
	Atk   EffortVal
	Def   EffortVal
	SpAtk EffortVal
	SpDef EffortVal
	Speed EffortVal
}

func (e *Effort) Sum() EffortVal {
	hp := e.HP
	atk := e.Atk
	def := e.Def
	spAtk := e.SpAtk
	spDef := e.SpDef
	speed := e.Speed
	return hp + atk + def + spAtk + spDef + speed
}

type RankVal int

const (
	MIN_RANK_VAL = RankVal(-6)
	MAX_RANK_VAL = RankVal(6)
)

func (v RankVal) ToBonus() float64 {
	if v >= 0 {
		return (float64(v) + 2.0) / 2.0
	} else {
		abs := float64(v) * -1
		return 2.0 / (abs + 2.0)
	}
}

type Rank struct {
	Atk   RankVal
	Def   RankVal
	SpAtk RankVal
	SpDef RankVal
	Speed RankVal
}

func (rank1 *Rank) Add(rank2 *Rank) Rank {
	atk := rank1.Atk + rank2.Atk
	def := rank1.Def + rank2.Def
	spAtk := rank1.SpAtk + rank2.SpAtk
	spDef := rank1.SpDef + rank2.SpDef
	speed := rank1.Speed + rank2.Speed
	return Rank{Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}
}

func (rank Rank) Regulate() Rank {
	rank.Atk = omwmath.Min(rank.Atk, MAX_RANK)
	rank.Def = omwmath.Min(rank.Def, MAX_RANK)
	rank.SpAtk = omwmath.Min(rank.SpAtk, MAX_RANK)
	rank.SpDef = omwmath.Min(rank.SpDef, MAX_RANK)
	rank.Speed = omwmath.Min(rank.Speed, MAX_RANK)

	rank.Atk = omwmath.Max(rank.Atk, MIN_RANK)
	rank.Def = omwmath.Max(rank.Def, MIN_RANK)
	rank.SpAtk = omwmath.Max(rank.SpAtk, MIN_RANK)
	rank.SpDef = omwmath.Max(rank.SpDef, MIN_RANK)
	rank.Speed = omwmath.Max(rank.Speed, MIN_RANK)
	return rank
}

func (rank *Rank) ContainsDown() bool {
	if rank.Atk < 0 {
		return true
	}

	if rank.Def < 0 {
		return true
	}

	if rank.SpAtk < 0 {
		return true
	}

	if rank.SpDef < 0 {
		return true
	}
	return rank.Speed < 0
}

func (rank Rank) ResetDown() Rank {
	rank.Atk = omwmath.Max(rank.Atk, 0)
	rank.Def = omwmath.Max(rank.Def, 0)
	rank.SpAtk = omwmath.Max(rank.SpAtk, 0)
	rank.SpDef = omwmath.Max(rank.SpDef, 0)
	rank.Speed = omwmath.Max(rank.Speed, 0)
	return rank
}

type Pokemon struct {
	Name    PokeName
	Types   Types
	Level   Level
	Gender  Gender
	Nature  Nature
	Ability Ability
	Item    Item
	Moveset Moveset

	Individual Individual
	Effort     Effort

	MaxHP     State
	CurrentHP State
	Atk       State
	Def       State
	SpAtk     State
	SpDef     State
	Speed     State

	StatusAilment        StatusAilment
	BadPoisonElapsedTurn int
	RankState            RankState
	ChoiceMoveName       MoveName

	//ひるみ
	IsFlinch bool
	//こらえる
	IsEndure bool
	//みがわり
	SubstituteDollHP int
	//やどりぎのタネ
	IsLeechSeed bool
	//こらえるの連続成功数
	EndureConsecutiveSuccessCount int
	//とんぼがえり・ボルトチェンジなどの攻撃後に交代する技
	AfterUTurn bool
	//ソーラービーム
	IsSolarBeamCharge bool
	//ちょうはつ
	TauntTurn int
}

func NewPokemon(pokeName PokeName, gender Gender, nature Nature, ability Ability, item Item,
	moveNames MoveNames, ppups PowerPointUps, ivState *IndividualState, evState *EffortState) (Pokemon, error) {
	pokeData, ok := POKEDEX[pokeName]
	if !ok {
		var msg string
		if pokeName == NO_POKE_NAME {
			msg = "ポケモン名 が ゼロ値 に なっている"
		} else {
			msg = fmt.Sprintf("%v という ポケモン名 は 存在しない", pokeName)
		}
		return Pokemon{}, fmt.Errorf(msg)
	}

	if !slices.Contains(ALL_GENDERS, gender) {
		return Pokemon{}, fmt.Errorf("性別 が 不適")
	}

	if !slices.Contains(ALL_NATURES, nature) {
		return Pokemon{}, fmt.Errorf("性格 が 不適")
	}

	if !slices.Contains(pokeData.Abilities, ability) {
		return Pokemon{}, fmt.Errorf("特性 が 不適")
	}

	if !slices.Contains(ALL_ITEMS, item) {
		return Pokemon{}, fmt.Errorf("アイテム が 不適")
	}

	moveset, err := NewMoveset(pokeName, moveNames, ppups)
	if err != nil {
		return Pokemon{}, err
	}

	pokemon := Pokemon{
		Name:pokeName, Types:pokeData.Types, Level:DEFAULT_LEVEL,
		Gender:gender, Nature:nature, Ability:ability, Item:item,
		Moveset:moveset, IndividualState:*ivState, EffortState:*evState,
	}
	
	natureData := NATUREDEX[nature]

	pokemon.MaxHP = StateCalculator.HP(pokeData.BaseHP, ivState.HP, evState.HP)
	pokemon.CurrentHP = pokemon.MaxHP
	pokemon.Atk = StateCalculator.OtherThanHP(pokeData.BaseAtk, ivState.Atk, evState.Atk, natureData.AtkBonus)
	pokemon.Def = StateCalculator.OtherThanHP(pokeData.BaseDef, ivState.Def, evState.Def, natureData.DefBonus)
	pokemon.SpAtk = StateCalculator.OtherThanHP(pokeData.BaseSpAtk, ivState.SpAtk, evState.SpAtk, natureData.SpAtkBonus)
	pokemon.SpDef = StateCalculator.OtherThanHP(pokeData.BaseSpDef, ivState.SpDef, evState.SpDef, natureData.SpDefBonus)
	pokemon.Speed = StateCalculator.OtherThanHP(pokeData.BaseSpeed, ivState.Speed, evState.Speed, natureData.SpeedBonus)
	return pokemon, nil
}

func (p1 *Pokemon) Equal(p2 *Pokemon) bool {
	if p1.Name != p2.Name {
		return false
	}

	if p1.Nature != p2.Nature {
		return false
	}

	if p1.Ability != p2.Ability {
		return false
	}

	if p1.Gender != p2.Gender {
		return false
	}

	if p1.Item != p2.Item {
		return false
	}

	if !p1.Moveset.Equal(p2.Moveset) {
		return false
	}

	if p1.MaxHP != p2.MaxHP {
		return false
	}

	if p1.CurrentHP != p2.CurrentHP {
		return false
	}

	if p1.Atk != p2.Atk {
		return false
	}

	if p1.Def != p2.Def {
		return false
	}

	if p1.SpAtk != p2.SpAtk {
		return false
	}

	if p1.SpDef != p2.SpDef {
		return false
	}

	if p1.Speed != p2.Speed {
		return false
	}

	if p1.IndividualState != p2.IndividualState {
		return false
	}

	if p1.EffortState != p2.EffortState {
		return false
	}

	for _, pokeType := range p1.Types {
		if !slices.Contains(p2.Types, pokeType) {
			return false
		}
	}

	if p1.RankState != p2.RankState {
		return false
	}

	if p1.StatusAilment != p2.StatusAilment {
		return false
	}

	if p1.BadPoisonElapsedTurn != p2.BadPoisonElapsedTurn {
		return false
	}

	if p1.ChoiceMoveName != p2.ChoiceMoveName {
		return false
	}

	if p1.IsLeechSeed != p2.IsLeechSeed {
		return false
	}

	return true
}

func (p *Pokemon) IsFullHP() bool {
	return p.MaxHP == p.CurrentHP
}

func (p *Pokemon) IsFaint() bool {
	return p.CurrentHP <= 0
}

func (p *Pokemon) IsFaintDamage(dmg int) bool {
	return dmg >= int(p.CurrentHP)
}

func (p *Pokemon) CurrentDamage() int {
	return int(p.MaxHP - p.CurrentHP)
}

func (p *Pokemon) SameTypeAttackBonus(moveName MoveName) SameTypeAttackBonus {
	moveType := MOVEDEX[moveName].Type
	isSameType := slices.Contains(p.Types, moveType)
	return NewSameTypeAttackBonus(isSameType)
}

func (p *Pokemon) EffectivenessBonus(moveName MoveName) EffectivenessBonus {
	y := 1.0
	moveType := MOVEDEX[moveName].Type
	for _, pokeType := range p.Types {
		y *= TYPEDEX[moveType][pokeType]
	}
	return EffectivenessBonus(y)
}

func (p *Pokemon) BadPoisonDamage() int {
	dmg := int(float64(p.MaxHP) * float64(p.BadPoisonElapsedTurn) / 16.0)
	return omwmath.Max(dmg, 1)
}