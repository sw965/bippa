package bippa

import (
	"fmt"
	"github.com/sw965/omw"
)

type PokeName string
type PokeNames []PokeName

func (pns PokeNames) Sort() PokeNames {
	index := func(pokeName PokeName) int { return omw.Index(ALL_POKE_NAMES, pokeName) }
	isChenge := func(i, j int) bool { return index(pns[i]) > index(pns[j]) }
	return omw.Sort(pns, isChenge)
}

type Type string

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

type Types []Type

type Level int

const (
	DEFAULT_LEVEL = Level(50)
)

type Gender string

const (
	MALE    = Gender("♂")
	FEMALE  = Gender("♀")
	UNKNOWN = Gender("不明")
)

type Genders []Gender

var ALL_GENDERS = Genders{MALE, FEMALE, UNKNOWN}

func NewVaildGenders(pokeName PokeName) Genders {
	switch POKEDEX[pokeName].Gender {
	case "♂♀両方":
		return Genders{MALE, FEMALE}
	case "♂のみ":
		return Genders{MALE}
	case "♀のみ":
		return Genders{FEMALE}
	default:
		return Genders{UNKNOWN}
	}
}

type Nature string
type Natures []Nature
type NatureBonus float64

const (
	NO_NATURE_BONUS   = NatureBonus(1.0)
	UP_NATURE_BONUS   = NatureBonus(1.1)
	DOWN_NATURE_BONUS = NatureBonus(0.9)
)

type Ability string
type Abilities []Ability

type Item string

const (
	EMPTY_ITEM = Item("なし")
)

func (item Item) IsChoice() bool {
	return item == "こだわりハチマキ" || item == "こだわりメガネ" || item == "こだわりスカーフ"
}

type Items []Item

var BATTLE_ITEMS = func() Items {
	result := Items{EMPTY_ITEM}
	return append(result, ALL_ITEMS...)
}()

type MoveName string

const (
	EMPTY_MOVE_NAME = MoveName("なし")
	STRUGGLE        = MoveName("わるあがき")
)

type MoveNames []MoveName

type PointUp int

var (
	MIN_POINT_UP = PointUp(0)
	MAX_POINT_UP = PointUp(3)
)

type PointUps []PointUp

var ALL_POINT_UPS = omw.MakeSliceRange[PointUp](MIN_POINT_UP, MAX_POINT_UP+1, 1)

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

const (
	MIN_MOVESET_LENGTH = 1
	MAX_MOVESET_LENGTH = 4
)

func NewMoveset(pokeName PokeName, moveNames MoveNames, pointUps []PointUp) (Moveset, error) {
	if !omw.IsUnique(moveNames) {
		return Moveset{}, fmt.Errorf("同じ技 を 覚えさせようとしている")
	}

	pokeData := POKEDEX[pokeName]
	learnset := pokeData.Learnset
	for _, moveName := range moveNames {
		if !omw.Contains(learnset, moveName) {
			return Moveset{}, fmt.Errorf("%v は %v を 覚えない", pokeName, moveName)
		}
	}

	if len(moveNames) != len(pointUps) {
		return Moveset{}, fmt.Errorf("len(moveName) != len(pointUps)")
	}

	pps := make([]PowerPoint, len(moveNames))
	for i, moveName := range moveNames {
		basePP := MOVEDEX[moveName].BasePP
		pointUp := pointUps[i]
		if !omw.Contains(ALL_POINT_UPS, pointUp) {
			return Moveset{}, fmt.Errorf("PointUpが、%v の 範囲外", ALL_POINT_UPS)
		}
		pps[i] = NewPowerPoint(basePP, pointUps[i])
	}

	y := omw.NewMap[MoveName, *PowerPoint](moveNames, omw.ValuesToPointers(pps))
	yLen := len(y)

	if yLen == 0 {
		return Moveset{}, fmt.Errorf("%v が 何も覚えていない", pokeName)
	} else if yLen > MAX_MOVESET_LENGTH {
		return Moveset{}, fmt.Errorf("%vつ以上 の 技 を 覚えさせようとしている", MAX_MOVESET_LENGTH+1)
	} else {
		return y, nil
	}
}

func (ms Moveset) Copy() Moveset {
	ks, vs := omw.Items(ms)
	vvs := omw.PointersToValues(vs)
	pvs := omw.ValuesToPointers(vvs)
	return omw.NewMap(ks, pvs)
}

func (ms1 Moveset) Equal(ms2 Moveset) bool {
	ks1 := omw.Keys(ms1)
	ks2 := omw.Keys(ms2)
	vs1 := omw.PointersToValues(omw.Values(ms1))
	vs2 := omw.PointersToValues(omw.Values(ms2))
	return omw.Equals(ks1, ks2) && omw.Equals(vs1, vs2)
}

type Individual int

const (
	EMPTY_INDIVIDUAL = Individual(-1)
	MIN_INDIVIDUAL   = Individual(0)
	MAX_INDIVIDUAL   = Individual(31)
)

type Individuals []Individual

var ALL_INDIVIDUALS = omw.MakeSliceRange[Individual](MIN_INDIVIDUAL, MAX_INDIVIDUAL+1, 1)

var BIPARAM_INDIVIDUALSS = func() []Individuals {
	v := Individuals{}
	err := omw.LoadJson(&v, BIPARAM_INDIVIDUALS_PATH)
	if err != nil {
		panic(err)
	}
	rng := func(i int) Individuals { return omw.MakeSliceRange(v[i], v[i+1], 1) }
	return omw.MakeSliceFunc(len(v)-1, rng)
}()

type IndividualState struct {
	HP    Individual
	Atk   Individual
	Def   Individual
	SpAtk Individual
	SpDef Individual
	Speed Individual
}

var ALL_MIN_INDIVIDUAL_STATE = IndividualState{
	HP: MIN_INDIVIDUAL, Atk: MIN_INDIVIDUAL, Def: MIN_INDIVIDUAL,
	SpAtk: MIN_INDIVIDUAL, SpDef: MIN_INDIVIDUAL, Speed: MIN_INDIVIDUAL,
}

var ALL_MAX_INDIVIDUAL_STATE = IndividualState{
	HP: MAX_INDIVIDUAL, Atk: MAX_INDIVIDUAL, Def: MAX_INDIVIDUAL,
	SpAtk: MAX_INDIVIDUAL, SpDef: MAX_INDIVIDUAL, Speed: MAX_INDIVIDUAL,
}

type Effort int

var (
	EMPTY_EFFORT   = Effort(-1)
	MIN_EFFORT     = Effort(0)
	MAX_EFFORT     = Effort(252)
	MAX_SUM_EFFORT = Effort(510)
)

func IsEffectiveEffort(effort Effort) bool {
	return effort%4 == 0
}

type Efforts []Effort

var ALL_EFFORTS = omw.MakeSliceRange[Effort](MIN_EFFORT, MAX_EFFORT+1, 1)
var EFFECTIVE_EFFORTS = omw.Filter(ALL_EFFORTS, IsEffectiveEffort)

var BIPARAM_EFFORTSS = func() []Efforts {
	v := Efforts{}
	err := omw.LoadJson(&v, BIPARAM_EFFORTS_PATH)
	if err != nil {
		panic(err)
	}
	rng := func(i int) Efforts { return omw.MakeSliceRange(v[i], v[i+1], 1) }
	return omw.MakeSliceFunc(len(v)-1, rng)
}()

type EffortState struct {
	HP    Effort
	Atk   Effort
	Def   Effort
	SpAtk Effort
	SpDef Effort
	Speed Effort
}

func (es *EffortState) Sum() Effort {
	hp := es.HP
	atk := es.Atk
	def := es.Def
	spAtk := es.SpAtk
	spDef := es.SpDef
	speed := es.Speed
	return hp + atk + def + spAtk + spDef + speed
}

func (es *EffortState) SumError() error {
	if es.Sum() > MAX_SUM_EFFORT {
		return fmt.Errorf("努力値 の 合計値 は、%v を 超えてはならない", MAX_SUM_EFFORT)
	}
	return nil
}

func CalcHp(baseHP int, individual Individual, effort Effort) int {
	intLevel := int(DEFAULT_LEVEL)
	result := ((baseHP*2)+int(individual)+(int(effort)/4))*intLevel/100 + intLevel + 10
	return result
}

func CalcState(baseState int, individual Individual, effort Effort, natureBonus NatureBonus) int {
	result := ((baseState*2)+int(individual)+(int(effort)/4))*int(DEFAULT_LEVEL)/100 + 5
	return int(float64(result) * float64(natureBonus))
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

type Rank int

const (
	MIN_RANK = Rank(-6)
	MAX_RANK = Rank(6)
)

func (rank Rank) ToBonus() RankBonus {
	if rank >= 0 {
		result := (float64(rank) + 2.0) / 2.0
		return RankBonus(result)
	} else {
		abs := float64(rank) * -1
		result := 2.0 / (abs + 2.0)
		return RankBonus(result)
	}
}

type RankState struct {
	Atk   Rank
	Def   Rank
	SpAtk Rank
	SpDef Rank
	Speed Rank
}

func (rs1 *RankState) Add(rs2 *RankState) RankState {
	atk := rs1.Atk + rs2.Atk
	def := rs1.Def + rs2.Def
	spAtk := rs1.SpAtk + rs2.SpAtk
	spDef := rs1.SpDef + rs2.SpDef
	speed := rs1.Speed + rs2.Speed
	return RankState{Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}
}

func (rs RankState) Regulate() RankState {
	rs.Atk = omw.Min(rs.Atk, MAX_RANK)
	rs.Def = omw.Min(rs.Def, MAX_RANK)
	rs.SpAtk = omw.Min(rs.SpAtk, MAX_RANK)
	rs.SpDef = omw.Min(rs.SpDef, MAX_RANK)
	rs.Speed = omw.Min(rs.Speed, MAX_RANK)

	rs.Atk = omw.Max(rs.Atk, MIN_RANK)
	rs.Def = omw.Max(rs.Def, MIN_RANK)
	rs.SpAtk = omw.Max(rs.SpAtk, MIN_RANK)
	rs.SpDef = omw.Max(rs.SpDef, MIN_RANK)
	rs.Speed = omw.Max(rs.Speed, MIN_RANK)
	return rs
}

func (rs *RankState) ContainsDown() bool {
	if rs.Atk < 0 {
		return true
	}

	if rs.Def < 0 {
		return true
	}

	if rs.SpAtk < 0 {
		return true
	}

	if rs.SpDef < 0 {
		return true
	}

	return rs.Speed < 0
}

func (rs RankState) ResetDown() RankState {
	rs.Atk = omw.Max(rs.Atk, 0)
	rs.Def = omw.Max(rs.Def, 0)
	rs.SpAtk = omw.Max(rs.SpAtk, 0)
	rs.SpDef = omw.Max(rs.SpDef, 0)
	rs.Speed = omw.Max(rs.Speed, 0)
	return rs
}

type RankBonus float64
