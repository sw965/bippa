package bippa

import (
	omwmath "github.com/sw965/omw/math"
	"golang.org/x/exp/slices"
	"math/rand"
	"fmt"
)

func IsHit(p int, r *rand.Rand) bool {
	return p > r.Intn(100)
}

type PokeName string
type PokeNames []PokeName

func (pns PokeNames) Sort() {
	isSwap := func(name1, name2 PokeName) bool {
		return slices.Index(ALL_POKE_NAMES, name1) > slices.Index(ALL_POKE_NAMES, name2)
	}
	slices.SortFunc(pns, isSwap)
}

type Level int

const (
	DEFAULT_LEVEL = Level(50)
)

type Ability string
type Abilities []Ability

type stateCalculator struct{}
var StateCalculator = stateCalculator{}

func (sc *stateCalculator) HP(base BaseState, iv Individual, ev Effort) State {
	lv := int(DEFAULT_LEVEL)
	y := ((int(base)*2)+int(iv)+(int(ev)/4))*lv/100 + lv + 10
	return State(y)
}

func (sc *stateCalculator) OtherThanHP(base BaseState, iv Individual, ev Effort, bonus NatureBonus) State {
	lv := int(DEFAULT_LEVEL)
	y := ((int(base)*2)+int(iv)+(int(ev)/4))*lv/100 + 5
	return State(float64(y) * float64(bonus))
}

type State int
type States []State

type Pokemon struct {
	Name    PokeName
	Types   Types
	Level   Level
	Gender  Gender
	Nature  Nature
	Ability Ability
	Item    Item
	Moveset Moveset

	IndividualState IndividualState
	EffortState     EffortState

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
	//やどりぎのタネ
	IsLeechSeed bool
	//こらえるの連続成功数
	EndureConsecutiveSuccessCount int
	//とんぼがえり・ボルトチェンジなどの攻撃後に交代する技
	AfterUTurn bool
}

func NewPokemon(pokeName PokeName, gender Gender, nature Nature, ability Ability, item Item,
	moveNames MoveNames, ppups PowerPointUps, ivState *IndividualState, evState *EffortState) (Pokemon, error) {
	pokeData, ok := POKEDEX[pokeName]
	if !ok {
		var msg string
		if pokeName == "" {
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

type CriticalRank int

// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%99%E3%81%B0%E3%82%84%E3%81%95#.E8.A9.B3.E7.B4.B0.E3.81.AA.E4.BB.95.E6.A7.98
type SpeedBonus int

const (
	INIT_SPEED_BONUS = SpeedBonus(4096)
)

func NewSpeedBonus(bt *Battle) SpeedBonus {
	y := int(INIT_SPEED_BONUS)
	if bt.P1Fighters[0].Item == CHOICE_SCARF {
		y = FiveOrMoreRounding(float64(y) * 6144.0 / 4096.0)
	}
	return SpeedBonus(y)
}

type FinalSpeed float64

func NewFinalSpeed(bt *Battle) FinalSpeed {
	speed := bt.P1Fighters[0].Speed
	rankBonus := bt.P1Fighters[0].RankState.Speed.ToBonus()
	bonus := NewSpeedBonus(bt)

	y := FiveOrMoreRounding(float64(speed) * float64(rankBonus))
	y = FiveOverRounding(float64(y) * float64(bonus) / 4096.0)
	return FinalSpeed(y)
}

type Action string

func (a Action) IsMoveName() bool {
	_, ok := MOVEDEX[MoveName(a)]
	return ok
}

func (a Action) IsPokeName() bool {
	_, ok := POKEDEX[PokeName(a)]
	return ok
}

func (a Action) Priority() int {
	if a == Action(STRUGGLE) {
		return 0
	} else if a.IsMoveName() {
		return MOVEDEX[MoveName(a)].PriorityRank
	} else {
		return 999
	}
}

type Actions []Action

type Winner struct {
	IsPlayer1 bool
	IsPlayer2 bool
}

var (
	WINNER_PLAYER1 = Winner{IsPlayer1: true, IsPlayer2: false}
	WINNER_PLAYER2 = Winner{IsPlayer1: false, IsPlayer2: true}
	DRAW           = Winner{IsPlayer1: false, IsPlayer2: false}
)

var WINNER_REWARD = map[Winner]float64{WINNER_PLAYER1: 1.0, WINNER_PLAYER2: 0.0, DRAW: 0.5}

// 小数点以下が0.5以上ならば、繰り上げ
func FiveOrMoreRounding(x float64) int {
	afterTheDecimalPoint := float64(x) - float64(int(x))
	if afterTheDecimalPoint >= 0.5 {
		return int(x + 1)
	}
	return int(x)
}

// 小数点以下が0.5より大きいならば、繰り上げ
func FiveOverRounding(x float64) int {
	afterTheDecimalPoint := float64(x) - float64(int(x))
	if afterTheDecimalPoint > 0.5 {
		return int(x + 1)
	}
	return int(x)
}

type PhysicsAttackBonus int

const (
	INIT_PHYSICS_ATTACK_BONUS = PhysicsAttackBonus(4096)
)

func NewPhysicsAttackBonus(poke *Pokemon) PhysicsAttackBonus {
	y := int(INIT_PHYSICS_ATTACK_BONUS)
	if poke.Item == CHOICE_BAND {
		y = FiveOrMoreRounding(float64(y) * 6144.0 / 4096.0)
	}
	return PhysicsAttackBonus(y)
}

type SpecialAttackBonus int

const (
	INIT_SPECIAL_ATTACK_BONUS = SpecialAttackBonus(4096)
)

func NewSpecialAttackBonus(poke *Pokemon) SpecialAttackBonus {
	y := int(INIT_SPECIAL_ATTACK_BONUS)
	if poke.Item == CHOICE_SPECS {
		y = FiveOrMoreRounding(float64(y) * 6144.0 / 4096.0)
	}
	return SpecialAttackBonus(y)
}

type AttackBonus int

func NewAttackBonus(poke *Pokemon, moveName MoveName) AttackBonus {
	moveData := MOVEDEX[moveName]
	switch moveData.Category {
	case PHYSICS:
		bonus := NewPhysicsAttackBonus(poke)
		return AttackBonus(bonus)
	case SPECIAL:
		bonus := NewSpecialAttackBonus(poke)
		return AttackBonus(bonus)
	default:
		return -1
	}
}

type FinalAttack int

func NewFinalAttack(poke *Pokemon, moveName MoveName, isCrit bool) FinalAttack {
	moveData := MOVEDEX[moveName]

	var atk State
	var rank Rank

	switch moveData.Category {
	case PHYSICS:
		atk = poke.Atk
		rank = poke.RankState.Atk
	case SPECIAL:
		atk = poke.SpAtk
		rank = poke.RankState.SpAtk
	}

	atkBonus := NewAttackBonus(poke, moveName)

	if rank < 0 && isCrit {
		rank = 0
	}

	rankBonus := rank.ToBonus()

	y := int(float64(atk) * float64(rankBonus))
	y = FiveOverRounding(float64(y) * float64(atkBonus) / 4096.0)
	return omwmath.Max(FinalAttack(y), 1)
}

type DefenseBonus int

const (
	INIT_DEFENSE_BONUS = DefenseBonus(4096)
)

func NewDefenseBonus(poke *Pokemon) DefenseBonus {
	y := INIT_DEFENSE_BONUS
	if poke.Item == ASSAULT_VEST {
		v := FiveOrMoreRounding(float64(y) * (6144.0 / 4096.0))
		y = DefenseBonus(v)
	}
	return y
}

type FinalDefense int

func NewFinalDefense(poke *Pokemon, moveName MoveName, isCrit bool) FinalDefense {
	moveData := MOVEDEX[moveName]

	var def State
	var rank Rank

	switch moveData.Category {
	case PHYSICS:
		def = poke.Def
		rank = poke.RankState.Def
	case SPECIAL:
		def = poke.SpDef
		rank = poke.RankState.SpDef
	}

	if rank > 0 && isCrit {
		rank = 0
	}

	bonus := rank.ToBonus()
	y := int(float64(def) * float64(bonus))
	return omwmath.Max(FinalDefense(y), 1)
}

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type PowerBonus int

const (
	INIT_POWER_BONUS = PowerBonus(4096)
)

type FinalPower int

func NewFinalPower(moveName MoveName) FinalPower {
	moveData := MOVEDEX[moveName]
	power := moveData.Power
	bonus := INIT_POWER_BONUS

	y := FiveOverRounding(float64(power) * float64(bonus) / 4096.0)
	return FinalPower(y)
}

type CriticalBonus float64

var (
	CRITICAL_BONUS    = CriticalBonus(6144.0 / 4096.0)
	NO_CRITICAL_BONUS = CriticalBonus(4096.0 / 4096.0)
)

var CRITICAL_N = map[CriticalRank]int{0: 24, 1: 8, 2: 2, 3: 1}

func NewCriticalBonus(isCrit bool) CriticalBonus {
	if isCrit {
		return CRITICAL_BONUS
	} else {
		return NO_CRITICAL_BONUS
	}
}

type SameTypeAttackBonus float64

const (
	SAME_TYPE_ATTACK_BONUS    = SameTypeAttackBonus(6144.0 / 4096.0)
	NO_SAME_TYPE_ATTACK_BONUS = SameTypeAttackBonus(4096.0 / 4096.0)
)

func NewSameTypeAttackBonus(isSTAB bool) SameTypeAttackBonus {
	if isSTAB {
		return SAME_TYPE_ATTACK_BONUS
	} else {
		return NO_SAME_TYPE_ATTACK_BONUS
	}
}

type EffectivenessBonus float64

// https://latest.pokewiki.net/%E3%83%80%E3%83%A1%E3%83%BC%E3%82%B8%E8%A8%88%E7%AE%97%E5%BC%8F
type RandomDamageBonus float64
type RandomDamageBonuses []RandomDamageBonus

var RANDOM_DAMAGE_BONUSES = RandomDamageBonuses{
	0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1.0,
}
var MAX_RANDOM_DAMAGE_BONUS = omwmath.Max(RANDOM_DAMAGE_BONUSES...)
var MEAN_RANDOM_DAMAGE_BONUS = omwmath.Mean(RANDOM_DAMAGE_BONUSES...)

type DamageBonus int

const (
	INIT_DAMAGE_BONUS = DamageBonus(4096)
)

func NewDamageBonus(poke *Pokemon) DamageBonus {
	y := INIT_DAMAGE_BONUS
	if poke.Item == LIFE_ORB {
		v := FiveOrMoreRounding(float64(y) * 5324.0 / 4096.0)
		y = DamageBonus(v)
	}
	return y
}

type FinalDamage int

func NewFinalDamage(attacker, defender *Pokemon, moveName MoveName, isCrit bool, randBonus RandomDamageBonus) FinalDamage {
	power := NewFinalPower(moveName)
	atk := NewFinalAttack(attacker, moveName, isCrit)
	def := NewFinalDefense(defender, moveName, isCrit)

	critBonus := NewCriticalBonus(isCrit)
	stab := attacker.SameTypeAttackBonus(moveName)
	effeBonus := defender.EffectivenessBonus(moveName)
	dmgBonus := NewDamageBonus(attacker)

	y := int(DEFAULT_LEVEL)*2/5 + 2
	y = int(float64(y) * float64(power) * float64(atk) / float64(def))
	y = y/50 + 2
	y = FiveOverRounding(float64(y) * float64(critBonus))
	y = int(float64(y) * float64(randBonus))
	y = FiveOverRounding(float64(y) * float64(stab))
	y = int(float64(y) * float64(effeBonus))
	y = FiveOverRounding(float64(y) * float64(dmgBonus) / 4096.0)
	return FinalDamage(y)
}