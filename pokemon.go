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

var STRING_TO_POKE_NAME = map[string]PokeName{
	"ニャオハ":NYA_O_HA,
}

func StringToPokeName(s string) PokeName {
	return STRING_TO_POKE_NAME[s]
}

type PokeNames []PokeName

func (pns PokeNames) Sort() {
	isSwap := func(name1, name2 PokeName) bool {
		return slices.Index(ALL_POKE_NAMES, name1) > slices.Index(ALL_POKE_NAMES, name2)
	}
	slices.SortFunc(pns, isSwap)
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