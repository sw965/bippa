package bippa

import (
	"fmt"
	"github.com/sw965/omw"
)

func InLearnset(pokeName PokeName, moveName MoveName) bool {
	pokeData := POKEDEX[pokeName]
	for _, iMoveName := range pokeData.Learnset {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func CalcHp(baseHP int, individualVal IndividualVal, effortVal EffortVal) int {
	intLevel := int(DEFAULT_LEVEL)
	result := ((baseHP*2)+int(individualVal)+(int(effortVal)/4))*intLevel/100 + intLevel + 10
	return result
}

func CalcState(baseState int, individualVal IndividualVal, effortVal EffortVal, natureBonus NatureBonus) int {
	result := ((baseState*2)+int(individualVal)+(int(effortVal)/4))*int(DEFAULT_LEVEL)/100 + 5
	return int(float64(result) * float64(natureBonus))
}

var (
	MAX_HP = CalcHp(MAX_BASE_HP, MAX_INDIVIDUAL_VAL, MAX_EFFORT_VAL)
	MAX_ATK = CalcState(MAX_BASE_ATK, MAX_INDIVIDUAL_VAL, MAX_EFFORT_VAL, UP_NATURE_BONUS)
	MAX_DEF = CalcState(MAX_BASE_DEF, MAX_INDIVIDUAL_VAL, MAX_EFFORT_VAL, UP_NATURE_BONUS)
	MAX_SP_ATK = CalcState(MAX_BASE_SP_ATK, MAX_INDIVIDUAL_VAL, MAX_EFFORT_VAL, UP_NATURE_BONUS)
	MAX_SP_DEF = CalcState(MAX_BASE_SP_DEF, MAX_INDIVIDUAL_VAL, MAX_EFFORT_VAL, UP_NATURE_BONUS)
	MAX_SPEED = CalcState(MAX_BASE_SPEED, MAX_INDIVIDUAL_VAL, MAX_EFFORT_VAL, UP_NATURE_BONUS)
)

type PokeName string

const (
	EMPTY_POKE_NAME = PokeName("なし")
)

func (pokeName PokeName) OnehotBinary() []int {
	var index int
	if pokeName != EMPTY_POKE_NAME {
		index = ALL_POKE_NAMES.Index(pokeName) + 1
	}
	return omw.OnehotBinary(index, ALL_POKE_NAMES_LENGTH + 1)
}
 
type PokeNames []PokeName

func (pokeNames PokeNames) Index(pokeName PokeName) int {
	for i, iPokeName := range pokeNames {
		if iPokeName == pokeName {
			return i
		}
	}
	return -1
} 

func (pokeNames PokeNames) Count(pokeName PokeName) int {
	result := 0
	for _, iPokeName := range pokeNames {
		if iPokeName == pokeName {
			result += 1
		}
	}
	return result
}

func (pokeNames PokeNames) IsUnique() bool {
	for _, pokeName := range pokeNames {
		if pokeNames.Count(pokeName) != 1 {
			return false
		}
	}
	return true
}

type Type string
type Types []Type

func (types Types) In(type_ Type) bool {
	for _, iType := range types {
		if iType == type_ {
			return true
		}
	}
	return false
}

const (
	NORMAL   = Type("ノーマル")
	FIRE     = Type("ほのお")
	WATER    = Type("みず")
	GRASS    = Type("くさ")
	ELECTRIC = Type("でんき")
	ICE      = Type("こおり")
	FIGHTING = Type("かくとう")
	POISON   = Type("どく")
	GROUND   = Type("じめん")
	FLYING   = Type("ひこう")
	PSYCHIC  = Type("エスパー")
	BUG      = Type("むし")
	ROCK     = Type("いわ")
	GHOST    = Type("ゴースト")
	DRAGON   = Type("ドラゴン")
	DARK     = Type("あく")
	STEEL    = Type("はがね")
	FAIRY    = Type("フェアリー")
)

type Level int

const (
	DEFAULT_LEVEL = Level(50)
)

type Nature string

func (nature Nature) OnehotBinary() []int {
	index := ALL_NATURES.Index(nature)
	return omw.OnehotBinary(index, ALL_NATURES_LENGTH)
}

type Natures []Nature

func (natures Natures) Index(nature Nature) int {
	for i, iNature := range natures {
		if iNature == nature {
			return i
		}
	}
	return -1
}

type NatureBonus float64

const (
	NO_NATURE_BONUS = NatureBonus(1.0)
	UP_NATURE_BONUS = NatureBonus(1.1)
	DOWN_NATURE_BONUS = NatureBonus(0.9)
)

type Ability string

func (ability Ability) IsValid(pokeName PokeName) bool {
	for _, iAbility := range POKEDEX[pokeName].AllAbilities {
		if iAbility == ability {
			return true
		}
	}
	return false
}

func (ability Ability) OnehotBinary() []int {
	index := ALL_ABILITIES.Index(ability)
	return omw.OnehotBinary(index, ALL_ABILITIES_LENGTH)
}

type Abilities []Ability

func (abilities Abilities) In(ability Ability) bool {
	for _, iAbility := range abilities {
		if iAbility == ability {
			return true
		}
	}
	return false
}

func (abilities Abilities) Index(ability Ability) int {
	for i, iAbility := range abilities {
		if iAbility == ability {
			return i
		}
	}
	return -1
}

type IndividualVal int

const (
	MIN_INDIVIDUAL_VAL = IndividualVal(0)
	MAX_INDIVIDUAL_VAL = IndividualVal(31)
)

func (individualVal IndividualVal) IsValid() bool {
	return individualVal >= MIN_INDIVIDUAL_VAL && individualVal <= MAX_INDIVIDUAL_VAL
}

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
	MIN_EFFORT_VAL     = EffortVal(0)
	MAX_EFFORT_VAL     = EffortVal(252)
	MAX_SUM_EFFORT_VAL = EffortVal(510)
)

type Effort struct {
	HP    EffortVal
	Atk   EffortVal
	Def   EffortVal
	SpAtk EffortVal
	SpDef EffortVal
	Speed EffortVal
}

func (effort *Effort) Sum() EffortVal {
	return effort.HP + effort.Atk + effort.Def + effort.SpAtk + effort.SpDef + effort.Speed
}

type Gender string

const (
	MALE    = Gender("♂")
	FEMALE  = Gender("♀")
	UNKNOWN = Gender("不明")
)

func (gender Gender) IsValid(pokeName PokeName) bool {
	genderData := POKEDEX[pokeName].Gender

	if genderData == "♂♀両方" {
		return gender == MALE || gender == FEMALE
	}

	if genderData == "♂のみ" {
		return gender == MALE
	}

	if genderData == "♀のみ" {
		return gender == FEMALE
	}

	return gender == UNKNOWN
}

func (gender Gender) OnehotBinary() []int {
	index := ALL_GENDERS.Index(gender)
	return omw.OnehotBinary(index, ALL_GENDERS_LENGTH)
}

type Genders []Gender

var ALL_GENDERS = Genders{MALE, FEMALE, UNKNOWN}
var ALL_GENDERS_LENGTH = len(ALL_GENDERS)

func (genders Genders) Index(gender Gender) int {
	for i, iGender := range genders {
		if iGender == gender {
			return i
		}
	}
	return -1
}

type Item string

const (
	EMPTY_ITEM = Item("なし")
)

func (item Item) IsValid() bool {
	for _, iItem := range ALL_ITEMS {
		if iItem == item {
			return true
		}
	}
	return item == EMPTY_ITEM
}

func (item Item) IsChoice() bool {
	return item == "こだわりハチマキ" || item == "こだわりメガネ" || item == "こだわりスカーフ"
}

func (item Item) OnehotBinary() []int {
	var index int
	if item != EMPTY_ITEM {
		index = ALL_ITEMS.Index(item) + 1
	}
	return omw.OnehotBinary(index, ALL_ITEMS_LENGTH + 1)
}

type Items []Item

func (items Items) Index(item Item) int {
	for i, iItem := range items {
		if iItem == item {
			return i
		}
	}
	return -1
}

type StatusAilment string

const (
	NORMAL_POISON = StatusAilment("どく")
	BAD_POISON    = StatusAilment("もうどく")
	SLEEP         = StatusAilment("ねむり")
	BURN          = StatusAilment("やけど")
	PARALYSIS     = StatusAilment("まひ")
	FREEZE        = StatusAilment("こおり")
)

type RankVal int

const (
	MIN_RANK_VAL = RankVal(-6)
	MAX_RANK_VAL = RankVal(6)
)

func (rankVal RankVal) ToBonus() RankBonus {
	if rankVal >= 0 {
		result := (float64(rankVal) + 2.0) / 2.0
		return RankBonus(result)
	} else {
		abs := float64(rankVal) * -1
		result := 2.0 / (abs + 2.0)
		return RankBonus(result)
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
	if rank.Atk > MAX_RANK_VAL {
		rank.Atk = MAX_RANK_VAL
	}

	if rank.Def > MAX_RANK_VAL {
		rank.Def = MAX_RANK_VAL
	}

	if rank.SpAtk > MAX_RANK_VAL {
		rank.SpAtk = MAX_RANK_VAL
	}

	if rank.SpDef > MAX_RANK_VAL {
		rank.SpDef = MAX_RANK_VAL
	}

	if rank.Speed > MAX_RANK_VAL {
		rank.Speed = MAX_RANK_VAL
	}

	if rank.Atk < MIN_RANK_VAL {
		rank.Atk = MIN_RANK_VAL
	}

	if rank.Def < MIN_RANK_VAL {
		rank.Def = MIN_RANK_VAL
	}

	if rank.SpAtk < MIN_RANK_VAL {
		rank.SpAtk = MIN_RANK_VAL
	}

	if rank.SpDef < MIN_RANK_VAL {
		rank.SpDef = MIN_RANK_VAL
	}

	if rank.Speed < MIN_RANK_VAL {
		rank.Speed = MIN_RANK_VAL
	}
	return rank
}

func (rank *Rank) InDown() bool {
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

func (rank *Rank) ResetDown() Rank {
	result := Rank{}

	if rank.Atk < 0 {
		result.Atk = 0
	} else {
		result.Atk = rank.Atk
	}

	if rank.Def < 0 {
		result.Def = 0
	} else {
		result.Def = rank.Def
	}

	if rank.SpAtk < 0 {
		result.SpAtk = 0
	} else {
		result.SpAtk = rank.SpAtk
	}

	if rank.SpDef < 0 {
		result.SpDef = 0
	} else {
		result.SpDef = rank.SpDef
	}

	if rank.Speed < 0 {
		result.Speed = 0
	} else {
		result.Speed = rank.Speed
	}

	return result
}

type RankBonus float64

type MoveName string

const (
	EMPTY_MOVE_NAME = MoveName("なし")
	STRUGGLE = MoveName("わるあがき")
)

func (moveName MoveName) OnehotBinary() []int {
	var index int
	if moveName != EMPTY_MOVE_NAME {
		index = ALL_MOVE_NAMES.Index(moveName) + 1
	}
	return omw.OnehotBinary(index, ALL_MOVE_NAMES_LENGTH + 1)
}

type MoveNames []MoveName

func (moveNames MoveNames) Count(moveName MoveName) int {
	result := 0
	for _, iMoveName := range moveNames {
		if iMoveName == moveName {
			result += 1
		}
	}
	return result
}

func (moveNames MoveNames) In(moveName MoveName) bool {
	for _, iMoveName := range moveNames {
		if iMoveName == moveName {
			return true
		}
	}
	return false
}

func (moveNames MoveNames) Sort() MoveNames {
	result := make(MoveNames, 0, len(moveNames))
	for _, moveName := range ALL_MOVE_NAMES {
		if moveNames.In(moveName) {
			result = append(result, moveName)
		}
	}

	for i := 0; i < moveNames.Count(EMPTY_MOVE_NAME); i++ {
		result = append(result, EMPTY_MOVE_NAME)
	}
	return result
}

func (moveNames MoveNames) Index(moveName MoveName) int {
	for i, iMoveName := range moveNames {
		if iMoveName == moveName {
			return i
		}
	}
	return -1
}

const (
	MIN_MOVESET_LENGTH = 1
	MAX_MOVESET_LENGTH = 4
)

type PointUp int

var (
	MIN_POINT_UP = PointUp(0)
	MAX_POINT_UP = PointUp(3)
)

func (pointUp PointUp) Binary() []int {
	result := make([]int, ALL_POINT_UPS_LENGTH)
	result[int(pointUp)] = 1
	return result
}

type PointUps []PointUp

var ALL_POINT_UPS = PointUps{0, 1, 2, 3}
var ALL_POINT_UPS_LENGTH = len(ALL_POINT_UPS)

func NewMaxPointUps(length int) PointUps {
	result := make(PointUps, length)
	for i := 0; i < length; i++ {
		result[i] = MAX_POINT_UP
	}
	return result
}

type PowerPoint struct {
	Max     int
	Current int
}

func NewPowerPoint(basePP int, pointUp PointUp) PowerPoint {
	v := (5.0 + float64(pointUp)) / 5.0
	max := int(float64(basePP) * v)
	return PowerPoint{Max: max, Current: max}
}

type PowerPoints []PowerPoint

type Moveset map[MoveName]*PowerPoint

func NewMoveset(pokeName PokeName, moveNames MoveNames, pointUps []PointUp) (Moveset, error) {
	for _, moveName := range moveNames {
		if !InLearnset(pokeName, moveName) {
			return Moveset{}, fmt.Errorf("%v は %v を 覚えない", pokeName, moveName)
		}
	}

	if len(moveNames) != len(pointUps) {
		return Moveset{}, fmt.Errorf("len(moveName) != len(pointUps)")
	}

	powerPoints := make([]PowerPoint, len(moveNames))
	for i, moveName := range moveNames {
		basePP := MOVEDEX[moveName].BasePP
		pointUp := pointUps[i]
		if !(MIN_POINT_UP <= pointUp && MAX_POINT_UP >= pointUp) {
			return Moveset{}, fmt.Errorf("ポイントアップが、%v～%vの範囲外", MIN_POINT_UP, MAX_POINT_UP)
		}
		powerPoints[i] = NewPowerPoint(basePP, pointUps[i])
	}

	moveset := Moveset{}
	for i, moveName := range moveNames {
		moveset[moveName] = &powerPoints[i]
	}

	movesetLength := len(moveset)

	if MIN_MOVESET_LENGTH <= movesetLength && MAX_MOVESET_LENGTH >= movesetLength {
		return moveset, nil
	} else {
		return Moveset{}, fmt.Errorf("覚えさせる技の数が、%v～%vの範囲外", MIN_MOVESET_LENGTH, MAX_MOVESET_LENGTH)
	}
}

func (moveset Moveset) Keys() MoveNames {
	result := make(MoveNames, 0, len(moveset))
	for k, _ := range moveset {
		result = append(result, k)
	}
	return result
}

func (moveset Moveset) PaddingKeys() MoveNames {
	result := make(MoveNames, 0, MAX_MOVESET_LENGTH)
	for k, _ := range moveset {
		result = append(result, k)
	}

	padNum := MAX_MOVESET_LENGTH - len(moveset)
	for i := 0; i < padNum; i++ {
		result = append(result, EMPTY_MOVE_NAME)
	}
	return result
}

func (moveset Moveset) Copy() Moveset {
	result := Moveset{}
	for moveName, powerPoint := range moveset {
		copyPowerPoint := PowerPoint{Max: powerPoint.Max, Current: powerPoint.Current}
		result[moveName] = &copyPowerPoint
	}
	return result
}

func (moveset1 Moveset) Equal(moveset2 Moveset) bool {
	for moveName1, powerPoint1 := range moveset1 {
		powerPoint2, ok := moveset2[moveName1]
		if !ok {
			return false
		}
		if powerPoint1.Max != powerPoint2.Max {
			return false
		}

		if powerPoint1.Current != powerPoint2.Current {
			return false
		}
	}
	return true
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

	MaxHP     int
	CurrentHP int
	Atk       int
	Def       int
	SpAtk     int
	SpDef     int
	Speed     int

	StatusAilment  StatusAilment
	BadPoisonElapsedTurn int
	Rank           Rank
	ChoiceMoveName MoveName

	IsLeechSeed bool
}

func NewPokemon(pokeName PokeName, nature Nature, ability Ability, gender Gender, item Item,
	moveNames MoveNames, pointUps PointUps, individual *Individual, effort *Effort) (Pokemon, error) {

	if pokeName == "" {
		return Pokemon{}, fmt.Errorf("ポケモン名が、ゼロ値になっている")
	}

	if nature == "" {
		return Pokemon{}, fmt.Errorf("性格が、ゼロ値になっている")
	}

	if ability == "" {
		return Pokemon{}, fmt.Errorf("特性が、ゼロ値になっている")
	}

	if gender == "" {
		return Pokemon{}, fmt.Errorf("性別が、ゼロ値になっている")
	}

	if item == "" {
		return Pokemon{}, fmt.Errorf("アイテムが、ゼロ値になっている。何も持たせない場合は、EMPTY_ITEMを使って。")
	}

	pokeData, ok := POKEDEX[pokeName]

	if !ok {
		return Pokemon{}, fmt.Errorf("ポケモン名 %v は 不適", pokeName)
	}

	natureData, ok := NATUREDEX[nature]

	if !ok {
		return Pokemon{}, fmt.Errorf("性格 %v は 不適", nature)
	}

	if !ability.IsValid(pokeName) {
		return Pokemon{}, fmt.Errorf("特性 %v の %v は不適", ability, pokeName)
	}

	if !gender.IsValid(pokeName) {
		return Pokemon{}, fmt.Errorf("性別 %v の %v は不適", gender, pokeName)
	}

	if !item.IsValid() {
		return Pokemon{}, fmt.Errorf("アイテム %v は 不適", item)
	}

	moveset, err := NewMoveset(pokeName, moveNames, pointUps)

	if err != nil {
		return Pokemon{}, err
	}

	hp := CalcHp(pokeData.BaseHP, individual.HP, effort.HP)
	atk := CalcState(pokeData.BaseAtk, individual.Atk, effort.Atk, natureData.AtkBonus)
	def := CalcState(pokeData.BaseDef, individual.Def, effort.Def, natureData.DefBonus)
	spAtk := CalcState(pokeData.BaseSpAtk, individual.SpAtk, effort.SpAtk, natureData.SpAtkBonus)
	spDef := CalcState(pokeData.BaseSpDef, individual.SpDef, effort.SpDef, natureData.SpDefBonus)
	speed := CalcState(pokeData.BaseSpeed, individual.Speed, effort.Speed, natureData.SpeedBonus)

	return Pokemon{Name: pokeName, Nature: nature, Ability: ability, Gender: gender, Item: item, Moveset: moveset,
		Individual: *individual, Effort: *effort,
		MaxHP: hp, CurrentHP: hp, Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed,
		Types: pokeData.Types, Level: DEFAULT_LEVEL}, nil
}

func NewEmptyPokemon() Pokemon {
	result := Pokemon{}
	result.Name = EMPTY_POKE_NAME
	return result
}

func (pokemon1 *Pokemon) Equal(pokemon2 *Pokemon) bool {
	if pokemon1.Name != pokemon2.Name {
		return false
	}

	if pokemon1.Nature != pokemon2.Nature {
		return false
	}

	if pokemon1.Ability != pokemon2.Ability {
		return false
	}

	if pokemon1.Gender != pokemon2.Gender {
		return false
	}

	if pokemon1.Item != pokemon2.Item {
		return false
	}

	if !pokemon1.Moveset.Equal(pokemon2.Moveset) {
		return false
	}

	if pokemon1.MaxHP != pokemon2.MaxHP {
		return false
	}

	if pokemon1.CurrentHP != pokemon2.CurrentHP {
		return false
	}

	if pokemon1.Atk != pokemon2.Atk {
		return false
	}

	if pokemon1.Def != pokemon2.Def {
		return false
	}

	if pokemon1.SpAtk != pokemon2.SpAtk {
		return false
	}

	if pokemon1.SpDef != pokemon2.SpDef {
		return false
	}

	if pokemon1.Speed != pokemon2.Speed {
		return false
	}

	if pokemon1.Individual != pokemon2.Individual {
		return false
	}

	if pokemon1.Effort != pokemon2.Effort {
		return false
	}

	for _, pokeType := range pokemon1.Types {
		if !pokemon2.Types.In(pokeType) {
			return false
		}
	}

	if pokemon1.Rank != pokemon2.Rank {
		return false
	}

	if pokemon1.StatusAilment != pokemon2.StatusAilment {
		return false
	}

	if pokemon1.BadPoisonElapsedTurn != pokemon2.BadPoisonElapsedTurn {
		return false
	}

	if pokemon1.ChoiceMoveName != pokemon2.ChoiceMoveName {
		return false
	}

	if pokemon1.IsLeechSeed != pokemon2.IsLeechSeed {
		return false
	}

	return true
}

func (pokemon *Pokemon) IsFullHP() bool {
	return pokemon.MaxHP == pokemon.CurrentHP
}

func (pokemon *Pokemon) IsFaint() bool {
	return pokemon.CurrentHP <= 0
}

func (pokemon *Pokemon) IsFaintDamage(damage int) bool {
	return damage >= int(pokemon.CurrentHP)
}

func (pokemon *Pokemon) CurrentDamage() int {
	return int(pokemon.MaxHP - pokemon.CurrentHP)
}

func (pokemon *Pokemon) SameTypeAttackBonus(moveName MoveName) SameTypeAttackBonus {
	moveType := MOVEDEX[moveName].Type
	inType := pokemon.Types.In(moveType)
	return BOOL_TO_SAME_TYPE_ATTACK_BONUS[inType]
}

func (pokemon *Pokemon) EffectivenessBonus(moveName MoveName) EffectivenessBonus {
	result := 1.0
	moveType := MOVEDEX[moveName].Type
	for _, pokeType := range pokemon.Types {
		result *= TYPEDEX[moveType][pokeType]
	}
	return EffectivenessBonus(result)
}

func (pokemon *Pokemon) BadPoisonDamage() int {
	damage := int(float64(pokemon.MaxHP) * float64(pokemon.BadPoisonElapsedTurn) / 16.0)
	if damage < 1 {
		return 1
	} else {
		return damage
	}
}

func (pokemon *Pokemon) MaxHPOnehotBinary() []int {
	return omw.OnehotBinary(pokemon.MaxHP, MAX_HP)
}

func (pokemon *Pokemon) AtkOnehotBinary() []int {
	return omw.OnehotBinary(pokemon.Atk, MAX_ATK)
}

func (pokemon *Pokemon) DefOnehotBinary() []int {
	return omw.OnehotBinary(pokemon.Def, MAX_DEF)
}

func (pokemon *Pokemon) SpAtkOnehotBinary() []int {
	return omw.OnehotBinary(pokemon.SpAtk, MAX_SP_ATK)
}

func (pokemon *Pokemon) SpDefOnehotBinary() []int {
	return omw.OnehotBinary(pokemon.SpDef, MAX_SP_DEF)
}

func (pokemon *Pokemon) SpeedOnehotBinary() []int {
	return omw.OnehotBinary(pokemon.Speed, MAX_SPEED)
}