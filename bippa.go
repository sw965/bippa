package bippa

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"math/rand"
	"fmt"
)

type PokeName string
type PokeNames []PokeName

func (pns PokeNames) Sort() {
	isSwap := func(name1, name2 PokeName) bool {
		return slices.Index(ALL_POKE_NAMES, name1) > slices.Index(ALL_POKE_NAMES, name2)
	}
	slices.SortFunc(pns, isSwap)
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

func NewValidGenders(pn PokeName) Genders {
	switch POKEDEX[pn].Gender {
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
	y := Items{EMPTY_ITEM}
	return append(y, ALL_ITEMS...)
}()

func GetHP(base BaseState, iv Individual, ev Effort) State {
	lv := int(DEFAULT_LEVEL)
	y := ((int(base)*2)+int(iv)+(int(ev)/4))*lv/100 + lv + 10
	return State(y)
}

func GetState(base BaseState, iv Individual, ev Effort, bonus NatureBonus) State {
	lv := int(DEFAULT_LEVEL)
	y := ((int(base)*2)+int(iv)+(int(ev)/4))*lv/100 + 5
	return State(float64(y) * float64(bonus))
}

const (
	EMPTY_STATE = State(-1)
)

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
		y := (float64(rank) + 2.0) / 2.0
		return RankBonus(y)
	} else {
		abs := float64(rank) * -1
		y := 2.0 / (abs + 2.0)
		return RankBonus(y)
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

	IsLeechSeed bool
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

	if !slices.Contains(pokeData.AllAbilities, ability) {
		return Pokemon{}, fmt.Errorf("特性 が 不適")
	}

	if item == "" {
		return Pokemon{}, fmt.Errorf("アイテム が ゼロ値 に なっている (何も持たせない場合は、EMPTY_ITEM を 使って)")
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

	pokemon.UpdateState()
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
	return omw.Max(dmg, 1)
}

func (p *Pokemon) UpdateState() {
	pokeData := POKEDEX[p.Name]
	natureData, ok := NATUREDEX[p.Nature]
	if !ok {
		natureData = NATUREDEX["てれや"]
	}

	p.MaxHP = GetHP(pokeData.BaseHP, p.IndividualState.HP, p.EffortState.HP)
	p.CurrentHP = p.MaxHP
	p.Atk = GetState(pokeData.BaseAtk, p.IndividualState.Atk, p.EffortState.Atk, natureData.AtkBonus)
	p.Def = GetState(pokeData.BaseDef, p.IndividualState.Def, p.EffortState.Def, natureData.DefBonus)
	p.SpAtk = GetState(pokeData.BaseSpAtk, p.IndividualState.SpAtk, p.EffortState.SpAtk, natureData.SpAtkBonus)
	p.SpDef = GetState(pokeData.BaseSpDef, p.IndividualState.SpDef, p.EffortState.SpDef, natureData.SpDefBonus)
	p.Speed = GetState(pokeData.BaseSpeed, p.IndividualState.Speed, p.EffortState.Speed, natureData.SpeedBonus)
}

type CriticalRank int

const (
	FIGHTERS_LENGTH = 3
)

type Fighters [FIGHTERS_LENGTH]Pokemon

func (fg1 *Fighters) Equal(fg2 *Fighters) bool {
	for i, poke := range fg1 {
		if !poke.Equal(&fg2[i]) {
			return false
		}
	}
	return true
}

func (fg *Fighters) PokeNames() PokeNames {
	y := make(PokeNames, FIGHTERS_LENGTH)
	for i, poke := range fg {
		y[i] = poke.Name
	}
	return y
}

func (fg *Fighters) IsAllFaint() bool {
	for _, poke := range fg {
		if !poke.IsFaint() {
			return false
		}
	}
	return true
}

func (fg *Fighters) LegalMoveNames() MoveNames {
	if fg[0].IsFaint() {
		return MoveNames{}
	}

	isPPZeroOver := func(moveName MoveName) bool { return fg[0].Moveset[moveName].Current > 0 }
	y := omw.Filter(maps.Keys(fg[0].Moveset), isPPZeroOver)

	if fg[0].ChoiceMoveName != "" {
		if slices.Contains(y, fg[0].ChoiceMoveName) {
			y = MoveNames{fg[0].ChoiceMoveName}
		}
	} else if fg[0].Item == "とつげきチョッキ" {
		isNotStatusMove := func(moveName MoveName) bool { return MOVEDEX[moveName].Category != STATUS }
		y = omw.Filter(y, isNotStatusMove)
	}

	if len(y) == 0 {
		return MoveNames{STRUGGLE}
	}
	return y
}

func (fg *Fighters) LegalPokeNames() PokeNames {
	y := make([]PokeName, 0, 2)
	for _, poke := range fg[1:] {
		if !poke.IsFaint() {
			y = append(y, poke.Name)
		}
	}
	return y
}

func (fg *Fighters) LegalBattleActions() BattleActions {
	moveNames := fg.LegalMoveNames()
	pokeNames := fg.LegalPokeNames()
	y := make(BattleActions, 0, len(moveNames)+len(pokeNames))
	y = append(y, omw.MapFunc[MoveNames, BattleActions](moveNames, omw.StrTildeToStrTilde[MoveName, BattleAction])...)
	y = append(y, omw.MapFunc[PokeNames, BattleActions](pokeNames, omw.StrTildeToStrTilde[PokeName, BattleAction])...)
	return y
}

// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%99%E3%81%B0%E3%82%84%E3%81%95#.E8.A9.B3.E7.B4.B0.E3.81.AA.E4.BB.95.E6.A7.98
type SpeedBonus int

const (
	INIT_SPEED_BONUS = SpeedBonus(4096)
)

func NewSpeedBonus(bt *Battle) SpeedBonus {
	y := int(INIT_SPEED_BONUS)
	if bt.P1Fighters[0].Item == "こだわりスカーフ" {
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

type BattleAction string

func (a BattleAction) IsMoveName() bool {
	_, ok := MOVEDEX[MoveName(a)]
	return ok
}

func (a BattleAction) IsPokeName() bool {
	_, ok := POKEDEX[PokeName(a)]
	return ok
}

func (a BattleAction) Priority() int {
	if a == BattleAction(STRUGGLE) {
		return 0
	} else if a.IsMoveName() {
		return MOVEDEX[MoveName(a)].PriorityRank
	} else {
		return 999
	}
}

type BattleActions []BattleAction

type Battle struct {
	P1Fighters Fighters
	P2Fighters Fighters
}

func (bt Battle) Reverse() Battle {
	return Battle{P1Fighters: bt.P2Fighters, P2Fighters: bt.P1Fighters}
}

func (bt1 *Battle) Equal(bt2 *Battle) bool {
	return bt1.P1Fighters.Equal(&bt2.P1Fighters) && bt1.P2Fighters.Equal(&bt2.P2Fighters)
}

func (bt *Battle) Accuracy(moveName MoveName) int {
	return MOVEDEX[moveName].Accuracy
}

func (bt *Battle) CriticalN(moveName MoveName) int {
	rank := MOVEDEX[moveName].CriticalRank
	return CRITICAL_N[rank]
}

func (bt *Battle) IsCritical(moveName MoveName, r *rand.Rand) bool {
	//https://wiki.xn--rckteqa2e.com/wiki/%E6%80%A5%E6%89%80
	n := bt.CriticalN(moveName)
	return r.Intn(n) == 0
}

func (bt Battle) Damage(dmg int) Battle {
	if bt.P1Fighters[0].IsFaint() {
		return bt
	}
	dmg = omw.Min(dmg, int(bt.P1Fighters[0].CurrentHP))
	bt.P1Fighters[0].CurrentHP -= State(dmg)
	return bt.SitrusBerryHeal()
}

func (bt Battle) Heal(heal int) Battle {
	if bt.P1Fighters[0].IsFaint() {
		return bt
	}
	heal = omw.Max(heal, bt.P1Fighters[0].CurrentDamage())
	bt.P1Fighters[0].CurrentHP += State(heal)
	return bt
}

func (bt Battle) SitrusBerryHeal() Battle {
	if bt.P1Fighters[0].Item != "オボンのみ" {
		return bt
	}

	if bt.P1Fighters[0].IsFaint() {
		return bt
	}

	max := bt.P1Fighters[0].MaxHP
	current := bt.P1Fighters[0].CurrentHP

	if int(current) <= int(float64(max)*1.0/2.0) {
		bt.P1Fighters[0].Item = EMPTY_ITEM
		heal := int(float64(max) * (1.0 / 4.0))
		bt = bt.Heal(heal)
	}
	return bt
}

func (bt Battle) RankStateFluctuation(stateP *RankState) Battle {
	if bt.P1Fighters[0].IsFaint() {
		return bt
	}

	state := stateP.Add(stateP)

	if bt.P1Fighters[0].Item == "しろいハーブ" && state.ContainsDown() {
		bt.P1Fighters[0].Item = EMPTY_ITEM
		state = state.ResetDown()
	}

	bt.P1Fighters[0].RankState = state.Regulate()
	return bt
}

func (bt Battle) AfterContact() Battle {
	if bt.P2Fighters[0].Ability == "てつのトゲ" {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 8.0)
		bt = bt.Damage(dmg)
	}

	if bt.P2Fighters[0].Item == "ゴツゴツメット" {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 6.0)
		bt = bt.Damage(dmg)
	}
	return bt
}

// https://latest.pokewiki.net/%E3%83%90%E3%83%88%E3%83%AB%E4%B8%AD%E3%81%AE%E5%87%A6%E7%90%86%E3%81%AE%E9%A0%86%E7%95%AA
func (bt Battle) MoveUse(moveName MoveName, r *rand.Rand) Battle {
	if bt.P1Fighters[0].IsFaint() {
		return bt
	}

	if moveName == STRUGGLE {
		if bt.P2Fighters[0].IsFaint() {
			return bt
		}
		bt.P1Fighters[0].CurrentHP = 0
		return bt
	}

	moveData := MOVEDEX[moveName]
	moveset := bt.P1Fighters[0].Moveset.Copy()
	moveset[moveName].Current -= 1
	bt.P1Fighters[0].Moveset = moveset

	if bt.P2Fighters[0].IsFaint() {
		return bt
	}

	if bt.P1Fighters[0].Item.IsChoice() {
		bt.P1Fighters[0].ChoiceMoveName = moveName
	}

	accuracy := moveData.Accuracy
	if accuracy != -1 {
		if r.Intn(100) >= accuracy {
			return bt
		}
	}

	if moveData.Category == STATUS {
		move, ok := STATUS_MOVES[moveName]
		if !ok {
			return bt
		}
		return move(bt, r)
	}

	attackNum := omw.RandInt(moveData.MinAttackNum, moveData.MaxAttackNum+1, r)
	for i := 0; i < attackNum; i++ {
		isCrit := bt.IsCritical(moveName, r)
		dmg := bt.NewFinalDamage(moveName, isCrit, omw.RandChoice(RANDOM_DAMAGE_BONUSES, r))

		if dmg == 0 {
			return bt
		}

		bt = bt.Reverse()
		bt = bt.Damage(int(dmg))
		bt = bt.Reverse()

		if moveData.Contact == "接触" {
			bt = bt.AfterContact()
		}

		if bt.P1Fighters[0].IsFaint() || bt.P2Fighters[0].IsFaint() {
			break
		}
	}

	if bt.P1Fighters[0].Item == "いのちのたま" {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 10.0)
		bt = bt.Damage(dmg)
	}
	return bt
}

func (bt Battle) Switch(pokeName PokeName) Battle {
	idx := slices.Index(bt.P1Fighters.PokeNames(), pokeName)

	bt.P1Fighters[0].RankState = RankState{}
	bt.P1Fighters[0].BadPoisonElapsedTurn = 0
	bt.P1Fighters[0].ChoiceMoveName = ""
	bt.P1Fighters[0].IsLeechSeed = false

	fg := bt.P1Fighters

	if idx == 1 {
		bt.P1Fighters[0] = fg[1]
		bt.P1Fighters[1] = fg[0]
		bt.P1Fighters[2] = fg[2]
	} else {
		bt.P1Fighters[0] = fg[2]
		bt.P1Fighters[1] = fg[1]
		bt.P1Fighters[2] = fg[0]
	}

	return bt
}

func (bt Battle) P1Action(action BattleAction, r *rand.Rand) Battle {
	if action == "" {
		return bt
	} else if action.IsMoveName() || action == BattleAction(STRUGGLE) {
		return bt.MoveUse(MoveName(action), r)
	} else {
		return bt.Switch(PokeName(action))
	}
}

func (bt Battle) P2Action(action BattleAction, r *rand.Rand) Battle {
	bt = bt.Reverse()
	if action == "" {
		return bt
	} else if action.IsMoveName() || action == BattleAction(STRUGGLE) {
		bt = bt.MoveUse(MoveName(action), r)
		return bt.Reverse()
	} else {
		bt = bt.Switch(PokeName(action))
		return bt.Reverse()
	}
}

func (bt *Battle) FinalSpeedWinner() Winner {
	p1 := NewFinalSpeed(bt)
	tb := bt.Reverse()
	p2 := NewFinalSpeed(&tb)

	if p1 > p2 {
		return WINNER_PLAYER1
	}

	if p1 < p2 {
		return WINNER_PLAYER2
	}
	return DRAW
}

func (bt *Battle) BattleActionPriorityWinner(p1Action, p2Action BattleAction) Winner {
	p1 := p1Action.Priority()
	p2 := p2Action.Priority()

	if p1 > p2 {
		return WINNER_PLAYER1
	}

	if p1 < p2 {
		return WINNER_PLAYER2
	}
	return DRAW
}

func (bt *Battle) ActionOrderWinner(p1, p2 BattleAction, r *rand.Rand) Winner {
	priorityWin := bt.BattleActionPriorityWinner(p1, p2)

	if priorityWin == WINNER_PLAYER1 {
		return WINNER_PLAYER1
	}

	if priorityWin == WINNER_PLAYER2 {
		return WINNER_PLAYER2
	}

	spWin := bt.FinalSpeedWinner()

	if spWin == WINNER_PLAYER1 {
		return WINNER_PLAYER1
	}

	if spWin == WINNER_PLAYER2 {
		return WINNER_PLAYER2
	}

	if omw.RandBool(r) {
		return WINNER_PLAYER1
	} else {
		return WINNER_PLAYER2
	}
}

func (bt Battle) TurnEnd(r *rand.Rand) Battle {
	//https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86
	p1First := func(bt Battle, f func(Battle) Battle) Battle {
		bt = f(bt)
		bt = bt.Reverse()
		bt = f(bt)
		bt = bt.Reverse()
		return bt
	}

	p2First := func(bt Battle, f func(Battle) Battle) Battle {
		bt = bt.Reverse()
		bt = f(bt)
		bt = bt.Reverse()
		bt = f(bt)
		return bt
	}

	run := func(fs []func(Battle) Battle) Battle {
		spWin := bt.FinalSpeedWinner()
		for _, f := range fs {
			if spWin == WINNER_PLAYER1 {
				bt = p1First(bt, f)
			} else if spWin == WINNER_PLAYER2 {
				bt = p2First(bt, f)
			} else {
				if omw.RandBool(r) {
					bt = p1First(bt, f)
				} else {
					bt = p2First(bt, f)
				}
			}
		}
		return bt
	}

	bt = run([]func(Battle) Battle{TurnEndLeftovers, TurnEndBlackSludge})
	bt = run([]func(Battle) Battle{TurnEndLeechSeed})
	bt = run([]func(Battle) Battle{TurnEndBadPoison})
	return bt
}

func (bt Battle) Push(p1Action, p2Action BattleAction, r *rand.Rand) Battle {
	isP1Faint := bt.P1Fighters[0].IsFaint()
	isP2Faint := bt.P2Fighters[0].IsFaint()

	if isP1Faint {
		bt = bt.P1Action(p1Action, r)
		if isP2Faint {
			bt = bt.P2Action(p2Action, r)
		}
		return bt
	}

	if isP2Faint {
		return bt.P2Action(p2Action, r)
	}

	orderWin := bt.ActionOrderWinner(p1Action, p2Action, r)

	var isP1s []bool

	if orderWin == WINNER_PLAYER1 {
		isP1s = []bool{true, false}
	} else {
		isP1s = []bool{false, true}
	}

	for _, isP1 := range isP1s {
		if isP1 {
			bt = bt.P1Action(p1Action, r)
		} else {
			bt = bt.P2Action(p2Action, r)
		}
	}

	if bt.IsGameEnd() {
		return bt
	}
	return bt.TurnEnd(r)
}

func (bt *Battle) IsGameEnd() bool {
	return bt.P1Fighters.IsAllFaint() || bt.P2Fighters.IsAllFaint()
}

func (bt *Battle) Winner() Winner {
	isP1AllFaint := bt.P1Fighters.IsAllFaint()
	isP2AllFaint := bt.P2Fighters.IsAllFaint()

	if isP1AllFaint && isP2AllFaint {
		return DRAW
	} else if isP1AllFaint {
		return WINNER_PLAYER2
	} else {
		return WINNER_PLAYER1
	}
}

func (bt *Battle) LegalActionss() []BattleActions {
	var p1 BattleActions
	var p2 BattleActions

	isP1Faint := bt.P1Fighters[0].IsFaint()
	isP2Faint := bt.P2Fighters[0].IsFaint()

	if !isP1Faint {
		p1 = bt.P1Fighters.LegalBattleActions()
	}

	if !isP2Faint {
		p2 = bt.P2Fighters.LegalBattleActions()
	}

	if isP1Faint && isP2Faint {
		p1 = bt.P1Fighters.LegalBattleActions()
		p2 = bt.P2Fighters.LegalBattleActions()
	}

	return []BattleActions{p1, p2}
}

func (bt *Battle) NewFinalDamage(moveName MoveName, isCrit bool, randDmgBonus RandomDamageBonus) FinalDamage {
	attaker := bt.P1Fighters[0]
	defender := bt.P2Fighters[0]
	return NewFinalDamage(&attaker, &defender, moveName, isCrit, randDmgBonus)
}

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

// 小数点以下がが0.5以上ならば、繰り上げ
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
	if poke.Item == "こだわりハチマキ" {
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
	if poke.Item == "こだわりメガネ" {
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
	return omw.Max(FinalAttack(y), 1)
}

type DefenseBonus int

const (
	INIT_DEFENSE_BONUS = DefenseBonus(4096)
)

func NewDefenseBonus(poke *Pokemon) DefenseBonus {
	y := INIT_DEFENSE_BONUS
	if poke.Item == "とつげきチョッキ" {
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
	return omw.Max(FinalDefense(y), 1)
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
var MAX_RANDOM_DAMAGE_BONUS = omw.Max(RANDOM_DAMAGE_BONUSES...)
var MEAN_RANDOM_DAMAGE_BONUS = omw.Mean(RANDOM_DAMAGE_BONUSES...)

type DamageBonus int

const (
	INIT_DAMAGE_BONUS = DamageBonus(4096)
)

func NewDamageBonus(poke *Pokemon) DamageBonus {
	y := INIT_DAMAGE_BONUS
	if poke.Item == "いのちのたま" {
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

const (
	PHYSICS = "物理"
	SPECIAL = "特殊"
	STATUS  = "変化"
)

func GetAllHPs(base BaseState) States {
	y := make(States, 0, len(ALL_INDIVIDUALS)*len(EFFECTIVE_EFFORTS))
	for _, iv := range ALL_INDIVIDUALS {
		for _, ev := range EFFECTIVE_EFFORTS {
			hp := GetHP(base, iv, ev)
			if !slices.Contains(y, hp) {
				y = append(y, hp)
			}
		}
	}
	return y
}

func GetAllStates(base BaseState) States {
	y := make(States, 0, len(ALL_INDIVIDUALS)*len(EFFECTIVE_EFFORTS)*len(ALL_NATURE_BONUSES))
	for _, iv := range ALL_INDIVIDUALS {
		for _, ev := range EFFECTIVE_EFFORTS {
			for _, bonus := range ALL_NATURE_BONUSES {
				state := GetState(base, iv, ev, bonus)
				if !slices.Contains(y, state) {
					y = append(y, state)
				}
			}
		}
	}
	return y
}