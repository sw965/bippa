package bippa

import (
	"os"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/maps"
	"github.com/sw965/omw"
	"github.com/sw965/crow"
	"math/rand"
	"strings"
)

var (
	SW965_PATH = os.Getenv("GOPATH") + "sw965/"

	DATA_PATH      = SW965_PATH + "arbok/data/"
	POKEDEX_PATH   = DATA_PATH + "pokedex/"
	MOVEDEX_PATH   = DATA_PATH + "movedex/"
	NATUREDEX_PATH = DATA_PATH + "naturedex.json"
	TYPEDEX_PATH   = DATA_PATH + "typedex.json"

	ALL_POKE_NAMES_PATH = DATA_PATH + "all_poke_names.json"
	ALL_NATURES_PATH    = DATA_PATH + "all_natures.json"
	ALL_MOVE_NAMES_PATH = DATA_PATH + "all_move_names.json"
	ALL_ITEMS_PATH      = DATA_PATH + "all_items.json"

	RATTA_PATH               = SW965_PATH + "ratta/"
)

type BaseState int

type PokeData struct {
	NormalAbilities Abilities
	HiddenAbility   Ability
	AllAbilities    Abilities

	Gender string
	Types  Types

	BaseHP    BaseState
	BaseAtk   BaseState
	BaseDef   BaseState
	BaseSpAtk BaseState
	BaseSpDef BaseState
	BaseSpeed BaseState

	Height    float64
	Weight    float64
	EggGroups []string
	Category  string

	Learnset MoveNames
}

func LoadPokeData(path string) PokeData {
	y, err := omw.LoadJson[PokeData](path)
	if err != nil {
		panic(err)
	}
	return y
}

type Pokedex map[PokeName]*PokeData

var POKEDEX = func() Pokedex {
	names, err := omw.DirNames(POKEDEX_PATH)
	if err != nil {
		panic(err)
	}
	y := Pokedex{}
	for _, name := range names {
		full := POKEDEX_PATH + name
		pokeName := strings.TrimRight(name, ".json")
		pokeData := LoadPokeData(full)
		y[PokeName(pokeName)] = &pokeData
	}
	return y
}()

var ALL_POKE_NAMES = func() PokeNames {
	y, err := omw.LoadJson[PokeNames](ALL_POKE_NAMES_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_ABILITIES = func() Abilities {
	y := make(Abilities, 0)
	for _, pokeData := range POKEDEX {
		for _, ability := range pokeData.AllAbilities {
			if !slices.Contains(y, ability) {
				y = append(y, ability)
			}
		}
	}
	return y
}()

type MoveData struct {
	Type     Type
	Category string
	Power    int
	Accuracy int
	BasePP   int
	Target   string

	Contact    string
	Protect    string
	MagicCoat  string
	Snatch     string
	MirrorMove string
	Substitute string

	GigantamaxMove  string
	GigantamaxPower int

	PriorityRank int
	CriticalRank CriticalRank

	MinAttackNum int
	MaxAttackNum int
}

func LoadMoveData(path string) MoveData {
	y, err := omw.LoadJson[MoveData](path)
	if err != nil {
		panic(err)
	}
	return y
}

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	y := Movedex{}
	names, err := omw.DirNames(MOVEDEX_PATH)
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		moveName := strings.TrimRight(name, ".json")
		moveData := LoadMoveData(MOVEDEX_PATH + name)
		y[MoveName(moveName)] = &moveData
	}
	return y
}()

var ALL_MOVE_NAMES = func() MoveNames {
	y, err := omw.LoadJson[MoveNames](ALL_MOVE_NAMES_PATH)	
	if err != nil {
		panic(err)
	}
	return y
}()

type NatureData struct {
	AtkBonus   NatureBonus
	DefBonus   NatureBonus
	SpAtkBonus NatureBonus
	SpDefBonus NatureBonus
	SpeedBonus NatureBonus
}

type Naturedex map[Nature]*NatureData

var NATUREDEX = func() Naturedex {
	y, err := omw.LoadJson[Naturedex](NATUREDEX_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_NATURES = func() Natures {
	y, err := omw.LoadJson[Natures](ALL_NATURES_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

type TypeData map[Type]float64
type Typedex map[Type]TypeData

var TYPEDEX = func() Typedex {
	y, err := omw.LoadJson[Typedex](TYPEDEX_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

var ALL_ITEMS = func() Items {
	y, err := omw.LoadJson[Items](ALL_ITEMS_PATH)
	if err != nil {
		panic(err)
	}
	return y
}()

type PokeName string
type PokeNames []PokeName

func (pns PokeNames) Sort() {
	isSwap := func(name1, name2 PokeName) bool { return slices.Index(ALL_POKE_NAMES, name1) > slices.Index(ALL_POKE_NAMES, name2)}
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

type MoveName string

const (
	EMPTY_MOVE_NAME = MoveName("なし")
	STRUGGLE        = MoveName("わるあがき")
)

type MoveNames []MoveName

type PowerPointUp int

var (
	MIN_POWER_POINT_UP = PowerPointUp(0)
	MAX_POWER_POINT_UP = PowerPointUp(3)
)

type PowerPointUps []PowerPointUp

var ALL_POWER_POINT_UPS = omw.MakeIntegerRange[PowerPointUps](MIN_POWER_POINT_UP, MAX_POWER_POINT_UP+1, 1)

type PowerPoint struct {
	Max     int
	Current int
}

func NewPowerPoint(base int, up PowerPointUp) PowerPoint {
	bonus := (5.0 + float64(up)) / 5.0
	max := int(float64(base) * bonus)
	return PowerPoint{Max: max, Current: max}
}

func NewMaxPowerPoint(moveName MoveName) PowerPoint {
	base := MOVEDEX[moveName].BasePP
	return NewPowerPoint(base, MAX_POWER_POINT_UP)
}

type PowerPoints []PowerPoint

type Moveset map[MoveName]*PowerPoint

const (
	MIN_MOVESET_LENGTH = 1
	MAX_MOVESET_LENGTH = 4
)

func (ms Moveset) Copy() Moveset {
	y := Moveset{}
	for k, v := range ms {
		pp := PowerPoint{Max:v.Max, Current:v.Current}
		y[k] = &pp
	}
	return y
}

func (ms1 Moveset) Equal(ms2 Moveset) bool {
	for k1, v1 := range ms1 {
		v2, ok := ms2[k1]
		if !ok {
			return false
		}

		if *v1 != *v2 {
			return false
		}
	}
	return true
}

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

type StatusAilment string

const (
	NORMAL_POISON = StatusAilment("どく")
	BAD_POISON    = StatusAilment("もうどく")
	SLEEP         = StatusAilment("ねむり")
	BURN          = StatusAilment("やけど")
	PARALYSIS     = StatusAilment("まひ")
	FREEZE        = StatusAilment("こおり")
)

type Individual int

const (
	EMPTY_INDIVIDUAL = Individual(-1)
	MIN_INDIVIDUAL   = Individual(0)
	MAX_INDIVIDUAL   = Individual(31)
)

type Individuals []Individual

var ALL_INDIVIDUALS Individuals = omw.MakeIntegerRange[Individuals](MIN_INDIVIDUAL, MAX_INDIVIDUAL+1, 1)

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

var INIT_INDIVIDUAL_STATE = IndividualState{HP:-1, Atk:-1, Def:-1, SpAtk:-1, SpDef:-1, Speed:-1}

type Effort int

var (
	EMPTY_EFFORT   = Effort(-1)
	MIN_EFFORT     = Effort(0)
	MAX_EFFORT     = Effort(252)
	MAX_SUM_EFFORT = Effort(510)
	EFFECTIVE_EFFORT = Effort(4)
)

type Efforts []Effort

var ALL_EFFORTS Efforts = omw.MakeIntegerRange[Efforts](MIN_EFFORT, MAX_EFFORT+1, 1)
var EFFECTIVE_EFFORTS Efforts = omw.Filter(ALL_EFFORTS, omw.IsRemainderZero(EFFECTIVE_EFFORT))

type EffortState struct {
	HP    Effort
	Atk   Effort
	Def   Effort
	SpAtk Effort
	SpDef Effort
	Speed Effort
}

var INIT_EFFORT_STATE = EffortState{HP:-1, Atk:-1, Def:-1, SpAtk:-1, SpDef:-1, Speed:-1}

func (es *EffortState) Sum() Effort {
	hp := es.HP
	atk := es.Atk
	def := es.Def
	spAtk := es.SpAtk
	spDef := es.SpDef
	speed := es.Speed
	return hp + atk + def + spAtk + spDef + speed
}

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

const (
	MIN_TEAM_LENGTH = 3
	MAX_TEAM_LENGTH = 6
)

type Team []Pokemon

func (team Team) IsValidLength() bool {
	n := len(team)
	return MIN_TEAM_LENGTH <= n && n <= MAX_TEAM_LENGTH
}

func (team Team) PokeNames() PokeNames {
	y := make(PokeNames, len(team))
	for i, poke := range team {
		y[i] = poke.Name
	}
	return y
}

func (team Team) Items() Items {
	y := make(Items, len(team))
	for i, poke := range team {
		y[i] = poke.Item
	}
	return y
}

func (team Team) Find(name PokeName) Pokemon {
	idx := slices.Index(team.PokeNames(), name)
	return team[idx]
}

func (team Team) Sort() Team {
	names := team.PokeNames()
	names.Sort()
	y := make(Team, 0, len(team))
	for _, name := range names {
		poke := team.Find(name)
		y = append(y, poke)
	}
	return y
}

func (team1 Team) Equal(team2 Team) bool {
	if len(team1) != len(team2) {
		return false
	}
	team1 = team1.Sort()
	team2 = team2.Sort()
	for i, poke := range team1 {
		if !poke.Equal(&team2[i]) {
			return false
		}
	}
	return true
}

func (team Team) LegalBuildCmds() TeamBuildCmds {
	pokeNames := team.PokeNames()
	if len(team.PokeNames()) < MAX_TEAM_LENGTH {
		f := func(pokeName PokeName) bool { return !slices.Contains(pokeNames, pokeName) }
		legalPokeNames := omw.Filter(ALL_POKE_NAMES, f)
		y := make(TeamBuildCmds, len(legalPokeNames))
		for i, pokeName := range legalPokeNames {
			y[i] = TeamBuildCmd{PokeName:pokeName}
		}
		return y
	}

	for i, pokemon := range team {
		if pokemon.Gender == "" {
			legalGenders := NewValidGenders(pokemon.Name)
			y := make(TeamBuildCmds, len(legalGenders))
			for j, gender := range legalGenders {
				y[j] = TeamBuildCmd{Gender:gender, Index:i}
			}
			return y
		}
	}

	for i, pokemon := range team {
		if pokemon.Nature == "" {
			y := make(TeamBuildCmds, len(ALL_NATURES))
			for j, nature := range ALL_NATURES {
				y[j] = TeamBuildCmd{Nature:nature, Index:i}
			}
			return y
		}
	}

	for i, pokemon := range team {
		if pokemon.Ability == "" {
			legalAbilities := POKEDEX[pokemon.Name].AllAbilities
			y := make(TeamBuildCmds, len(legalAbilities))
			for j, ability := range legalAbilities {
				y[j] = TeamBuildCmd{Ability:ability, Index:i}
			}
			return y
		}
	}

	for i, pokemon := range team {
		if pokemon.Item == "" {
			y := make(TeamBuildCmds, 0, len(ALL_ITEMS))
			items := team.Items()
			for _, item := range ALL_ITEMS {
				if slices.Contains(items, item) {
					continue
				}
				y = append(y, TeamBuildCmd{Item:item, Index:i})
			}
			return y
		}
	}

	for i, pokemon := range team {
		learnset := POKEDEX[pokemon.Name].Learnset
		n := omw.Min(len(learnset), MAX_MOVESET_LENGTH)
		if len(pokemon.Moveset) < n {
			legalMoveNames := POKEDEX[pokemon.Name].Learnset
			y := make(TeamBuildCmds, len(legalMoveNames))
			for j, moveName := range legalMoveNames {
				y[j] = TeamBuildCmd{MoveName:moveName, Index:i}
			}
			return y
		}
	}
	return TeamBuildCmds{}
}

func (team Team) Push(cmd *TeamBuildCmd) Team {
	y := make(Team, 0, MAX_TEAM_LENGTH)
	for _, pokemon := range team {
		y = append(y, pokemon)
	}

	if cmd.PokeName != "" {
		y = append(y, Pokemon{Name:cmd.PokeName, Moveset:Moveset{}, IndividualState:INIT_INDIVIDUAL_STATE, EffortState:INIT_EFFORT_STATE})
		return y
	}

	idx := cmd.Index
	if cmd.Gender != "" {
		y[idx].Gender = cmd.Gender
		return y
	}

	if cmd.Nature != "" {
		y[idx].Nature = cmd.Nature
		return y
	}

	if cmd.Ability != "" {
		y[idx].Ability = cmd.Ability
		return y
	}

	if cmd.Item != "" {
		y[idx].Item = cmd.Item
		return y
	}

	if cmd.MoveName != "" {
		pp := NewMaxPowerPoint(cmd.MoveName)
		y[idx].Moveset[cmd.MoveName] = &pp
		return y
	}
	return Team{}
}

type TeamBuildCmd struct {
	PokeName PokeName
	Gender Gender
	Nature Nature
	Ability Ability
	Item Item
	MoveName MoveName
	Index int
}

type TeamBuildCmds []TeamBuildCmd

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
	P1Fighters  Fighters
	P2Fighters  Fighters
}

func (bt *Battle) Reverse() Battle {
	return Battle{P1Fighters: bt.P2Fighters, P2Fighters: bt.P1Fighters}
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

func (bt *Battle) Winner() Winner{
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
	DRAW      = Winner{IsPlayer1: false, IsPlayer2: false}
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

type StatusMove func(Battle, *rand.Rand) Battle

// あさのひざし
func MorningSun(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// こうごうせい
func Synthesis(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// じこさいせい
func Recover(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// すなあつめ
func ShoreUp(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// タマゴうみ
func SoftBoiled(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// つきのひかり
func Moonlight(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// なまける
func SlackOff(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// はねやすめ
func Roost(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// ミルクのみ
func MilkDrink(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// どくどく
func Toxic(bt Battle, _ *rand.Rand) Battle {
	if bt.P2Fighters[0].StatusAilment != "" {
		return bt
	}

	p2PokeTypes := bt.P2Fighters[0].Types

	if slices.Contains(p2PokeTypes, POISON) {
		return bt
	}

	if slices.Contains(p2PokeTypes, STEEL) {
		return bt
	}

	bt.P2Fighters[0].StatusAilment = BAD_POISON
	return bt
}

// やどりぎのタネ
func LeechSeed(bt Battle, _ *rand.Rand) Battle {
	if slices.Contains(bt.P2Fighters[0].Types, GRASS) {
		return bt
	}
	bt.P2Fighters[0].IsLeechSeed = true
	return bt
}

// つるぎのまい
func SwordsDance(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Atk: 2})
}

// りゅうのまい
func DragonDance(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Atk: 1, Speed: 1})
}

// からをやぶる
func ShellSmash(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Atk: 2, Def: -1, SpAtk: 2, SpDef: -1, Speed: 2})
}

// てっぺき
func IronDefense(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Def: 2})
}

// めいそう
func CalmMind(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{SpAtk: 1, SpDef: 1})
}

var STATUS_MOVES = map[MoveName]StatusMove{
	"あさのひざし":  Moonlight,
	"こうごうせい":  Synthesis,
	"じこさいせい":  Recover,
	"すなあつめ":   ShoreUp,
	"タマゴうみ":   SoftBoiled,
	"つきのひかり":  Moonlight,
	"なまける":    SlackOff,
	"はねやすめ":   Roost,
	"ミルクのみ":   MorningSun,
	"どくどく":    Toxic,
	"やどりぎのタネ": LeechSeed,
	"つるぎのまい":  SwordsDance,
	"りゅうのまい":  DragonDance,
	"からをやぶる":  ShellSmash,
	"てっぺき":    IronDefense,
}

// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86
func TurnEndLeftovers(bt Battle) Battle {
	if bt.P1Fighters[0].Item != "たべのこし" {
		return bt
	}

	if bt.P1Fighters[0].IsFullHP() {
		return bt
	}

	heal := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 16.0)
	bt = bt.Heal(heal)
	return bt
}

func TurnEndBlackSludge(bt Battle) Battle {
	if bt.P1Fighters[0].Item != "くろいヘドロ" {
		return bt
	}

	if slices.Contains(bt.P1Fighters[0].Types, POISON) {
		heal := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 16.0)
		heal = omw.Max(heal, 1)
		bt = bt.Heal(heal)
	} else {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 8.0)
		dmg = omw.Max(dmg, 1)
		bt = bt.Damage(dmg)
	}
	return bt
}

func TurnEndLeechSeed(bt Battle) Battle {
	if bt.P1Fighters[0].IsFaint() {
		return bt
	}

	if bt.P2Fighters[0].IsFaint() {
		return bt
	}

	if !bt.P2Fighters[0].IsLeechSeed {
		return bt
	}

	dmg := int(float64(bt.P2Fighters[0].MaxHP) * 1.0 / 8.0)
	heal := dmg

	bt = bt.Reverse()
	bt = bt.Damage(dmg)
	bt = bt.Reverse()
	bt = bt.Heal(heal)
	return bt
}

func TurnEndBadPoison(bt Battle) Battle {
	if bt.P1Fighters[0].StatusAilment != BAD_POISON {
		return bt
	}

	if bt.P1Fighters[0].BadPoisonElapsedTurn < 16 {
		bt.P1Fighters[0].BadPoisonElapsedTurn += 1
	}

	dmg := bt.P1Fighters[0].BadPoisonDamage()
	return bt.Damage(dmg)
}

func GetAllHPs(base BaseState) States {
	y := make(States, 0, len(ALL_INDIVIDUALS) * len(EFFECTIVE_EFFORTS))
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
	y := make(States, 0, len(ALL_INDIVIDUALS) * len(EFFECTIVE_EFFORTS) * len(ALL_NATURE_BONUSES))
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

type BaseModel[K comparable] map[K]float64
type TwoRelationshipModel[K1, K2 comparable] map[K1]BaseModel[K2]

func NewInitSameFeatureValueTwoRelationshipModel[FS ~[]K, K comparable](features FS, r *rand.Rand) TwoRelationshipModel[K, K] {
	y := TwoRelationshipModel[K, K]{}
	permutation2 := omw.Permutation[[]FS](features, 2)
	for _, vs := range permutation2 {
		k1 := vs[0]
		if _, ok := y[k1]; !ok {
			y[k1] = BaseModel[K]{}
		}
		k2 := vs[1]
		y[k1][k2] = omw.RandFloat64(0.0, 16.0, r)
	}
	return y
}

func NewInitDifferentFeatureValueTwoRelationshipModel[FS1 ~[]K1, FS2 ~[]K2, K1, K2 comparable](features1 FS1, features2 FS2, r *rand.Rand) TwoRelationshipModel[K1, K2] {
	y := TwoRelationshipModel[K1, K2]{}
	for _, k1 := range features1 {
		if _, ok := y[k1]; !ok {
			y[k1] = BaseModel[K2]{}
		}
		for _, k2 := range features2 {
			y[k1][k2] = omw.RandFloat64(0.0, 16.0, r)
		}
	}
	return y
}

type PokemonModel struct {
	MoveNameAndMoveName TwoRelationshipModel[MoveName, MoveName]
	MoveNameAndGender TwoRelationshipModel[MoveName, Gender]
	MoveNameAndAbility TwoRelationshipModel[MoveName, Ability]
	MoveNameAndItem TwoRelationshipModel[MoveName, Item]
	MoveNameAndHP TwoRelationshipModel[MoveName, State]
	MoveNameAndAtk TwoRelationshipModel[MoveName, State]
	MoveNameAndDef TwoRelationshipModel[MoveName, State]
	MoveNameAndSpAtk TwoRelationshipModel[MoveName, State]
	MoveNameAndSpDef TwoRelationshipModel[MoveName, State]
	MoveNameAndSpeed TwoRelationshipModel[MoveName, State]
}

func NewInitPokemonModel(pokeName PokeName, r *rand.Rand) PokemonModel {
	pokeData := POKEDEX[pokeName]
	genders := NewValidGenders(pokeName)
	hps := GetAllHPs(pokeData.BaseHP)
	atks := GetAllStates(pokeData.BaseAtk)
	defs := GetAllStates(pokeData.BaseDef)
	spAtks := GetAllStates(pokeData.BaseSpAtk)
	spDefs := GetAllStates(pokeData.BaseSpDef)
	speeds := GetAllStates(pokeData.BaseSpeed)

	moveNameAndMoveName := NewInitSameFeatureValueTwoRelationshipModel[MoveNames](pokeData.Learnset, r)
	moveNameAndGender := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, Genders](pokeData.Learnset, genders, r)
	moveNameAndAbility := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, Abilities](pokeData.Learnset, pokeData.AllAbilities, r)
	moveNameAndItem := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, Items](pokeData.Learnset, ALL_ITEMS, r)
	moveNameAndHP := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, States](pokeData.Learnset, hps, r)
	moveNameAndAtk := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, States](pokeData.Learnset, atks, r)
	moveNameAndDef := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, States](pokeData.Learnset, defs, r)
	moveNameAndSpAtk := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, States](pokeData.Learnset, spAtks, r)
	moveNameAndSpDef := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, States](pokeData.Learnset, spDefs, r)
	moveNameAndSpeed := NewInitDifferentFeatureValueTwoRelationshipModel[MoveNames, States](pokeData.Learnset, speeds, r)

	return PokemonModel{
		MoveNameAndMoveName:moveNameAndMoveName,
		MoveNameAndGender:moveNameAndGender,
		MoveNameAndAbility:moveNameAndAbility,
		MoveNameAndItem:moveNameAndItem,
		MoveNameAndHP:moveNameAndHP,
		MoveNameAndAtk:moveNameAndAtk,
		MoveNameAndDef:moveNameAndDef,
		MoveNameAndSpAtk:moveNameAndSpAtk,
		MoveNameAndSpDef:moveNameAndSpDef,
		MoveNameAndSpeed:moveNameAndSpeed,
	}
}

func (model PokemonModel) ParameterSize() int {
	y := 0
	y += omw.NestMapSize(model.MoveNameAndMoveName)
	y += omw.NestMapSize(model.MoveNameAndGender)
	y += omw.NestMapSize(model.MoveNameAndAbility)
	y += omw.NestMapSize(model.MoveNameAndItem)
	y += omw.NestMapSize(model.MoveNameAndHP)
	y += omw.NestMapSize(model.MoveNameAndAtk)
	y += omw.NestMapSize(model.MoveNameAndDef)
	y += omw.NestMapSize(model.MoveNameAndSpAtk)
	y += omw.NestMapSize(model.MoveNameAndSpDef)
	y += omw.NestMapSize(model.MoveNameAndSpeed)
	return y
}

func (model *PokemonModel) Output(pokemon *Pokemon) float64 {
	y := 0.0
	moveNames := maps.Keys(pokemon.Moveset)

	permutation2 := omw.Permutation[[]MoveNames, MoveNames](moveNames, 2)
	for _, mns := range permutation2 {
		y += model.MoveNameAndMoveName[mns[0]][mns[1]]
	}

	for _, moveName := range moveNames {
		y += model.MoveNameAndGender[moveName][pokemon.Gender]
		y += model.MoveNameAndAbility[moveName][pokemon.Ability]
		y += model.MoveNameAndItem[moveName][pokemon.Item]
		y += model.MoveNameAndHP[moveName][pokemon.MaxHP]
		y += model.MoveNameAndAtk[moveName][pokemon.Atk]
		y += model.MoveNameAndDef[moveName][pokemon.Def]
		y += model.MoveNameAndSpAtk[moveName][pokemon.SpAtk]
		y += model.MoveNameAndSpDef[moveName][pokemon.SpDef]
		y += model.MoveNameAndSpeed[moveName][pokemon.Speed]
	}
	return y
}

func NewTeamBuildPUCTFunCaller(model *PokemonModel, r *rand.Rand) crow.PUCT_FunCaller[Team, TeamBuildCmd] {
	legalActions := func(team *Team) []TeamBuildCmd {
		return team.LegalBuildCmds()
	}

	push := func(team Team, action TeamBuildCmd) Team {
		return team.Push(&action)
	}

	isEnd := func(team *Team) bool {
		return len(team.LegalBuildCmds()) == 0
	}

	equal := func(team1, team2 *Team) bool {
		return team1.Equal(*team2)
	}

	game := crow.AlternatelyMoveGameFunCaller[Team, TeamBuildCmd]{
		LegalActions:legalActions,
		Push:push,
		EqualState:equal,
		IsEnd:isEnd,
	}

	game.SetRandomActionPlayer(r)
	
	leaf := func(team *Team) crow.PUCT_LeafEvalY {
		y := 0.0
		for _, pokemon := range *team {
			y += model.Output(&pokemon)
		}
		return crow.PUCT_LeafEvalY(y)
	}

	backward := func(y crow.PUCT_LeafEvalY, team *Team) crow.PUCT_BackwardEvalY {
		return crow.PUCT_BackwardEvalY(y)
	}

	eval := crow.PUCT_EvalFunCaller[Team]{
		Leaf:leaf,
		Backward:backward,
	}

	fnCaller := crow.PUCT_FunCaller[Team, TeamBuildCmd]{
		Game:game,
		Eval:eval,
	}

	fnCaller.SetNoPolicy()
	return fnCaller
}