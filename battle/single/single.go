package single

import (
	"fmt"
	"math/rand"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	omwmath "github.com/sw965/omw/math"
	"github.com/sw965/bippa/battle/dmgtools"
)

const (
	FIGHTER_NUM = 3
)

type Fighters [FIGHTER_NUM]bp.Pokemon

func (f *Fighters) Names() bp.PokeNames {
	ret := make(bp.PokeNames, FIGHTER_NUM)
	for i, pokemon := range f {
		ret[i] = pokemon.Name
	}
	return ret
}

func (f *Fighters) Equal(other *Fighters) bool {
	for i := range f {
		pokemon1 := f[i]
		pokemon2 := other[i]
		if !pokemon1.Equal(&pokemon2) {
			return false
		}
	}
	return true
}

func (f *Fighters) IndexByName(name bp.PokeName) int {
	for i := range f {
		if f[i].Name == name {
			return i
		}
	}
	return -1
}

func (f *Fighters) IsAllFaint() bool {
	for i := range f {
		if f[i].CurrentHP > 0 {
			return false
		}
	}
	return true
}

type Action struct {
	CmdMoveName bp.MoveName
	SwitchPokeName bp.PokeName
	IsPlayer1 bool
}

func (a *Action) ToString() string {
	p := map[bool]string{true:"player1", false:"player2"}[a.IsPlayer1]
	return bp.MOVE_NAME_TO_STRING[a.CmdMoveName] + bp.POKE_NAME_TO_STRING[a.SwitchPokeName] + " " + p
}

func (a *Action) IsEmpty() bool {
	return a.CmdMoveName == bp.EMPTY_MOVE_NAME && a.SwitchPokeName == bp.EMPTY_POKE_NAME
}

func (a *Action) IsCommandMove() bool {
	return a.CmdMoveName != bp.EMPTY_MOVE_NAME
}

func (a *Action) IsSwitch() bool {
	return a.SwitchPokeName != bp.EMPTY_POKE_NAME
}

type Actions []Action

func (as Actions) IsAllEmpty() bool {
	for i := range as {
		if !as[i].IsEmpty() {
			return false
		}
	}
	return true
}

func (as Actions) ToStrings() []string {
	ret := make([]string, len(as))
	for i, a := range as {
		ret[i] = a.ToString()
	}
	return ret
}

const (
	MAX_SIMULTANEOUS_ACTION_NUM = 2
)

type ActionSlices []Actions

type Step int

const (
	BEFORE_MOVE_USE_STEP Step = iota
	AFTER_MOVE_USE_STEP
	BEFORE_MOVE_DAMAGE_STEP
	AFTER_MOVE_DAMAGE_STEP
	BEFORE_SWITCH_STEP
	AFTER_SWITCH_STEP
	SELF_FAINT_STEP
	OPPONENT_FAINT_STEP
)

type Battle struct {
	SelfFighters Fighters
	OpponentFighters Fighters
	Turn int
	IsRealSelf bool

	RandDmgBonuses dmgtools.RandBonuses
	Observer func(*Battle, Step)
}

func (b Battle) SwapPlayers() Battle {
	b.SelfFighters, b.OpponentFighters = b.OpponentFighters, b.SelfFighters
	b.IsRealSelf = !b.IsRealSelf
	return b
}

func (b *Battle) CalcDamage(moveName bp.MoveName, r *rand.Rand) int {
	attacker := b.SelfFighters[0]
	defender := b.OpponentFighters[0]

	calculator := dmgtools.Calculator{
		dmgtools.Attacker{
			PokeName:attacker.Name,
			Level:attacker.Level,
			Atk:attacker.Atk,
			SpAtk:attacker.SpAtk,
		},
		dmgtools.Defender{
			PokeName:defender.Name,
			Level:defender.Level,
			Def:defender.Def,
			SpDef:defender.SpDef,
		},
	}
	randDmgBonus := omwrand.Choice(b.RandDmgBonuses, r)
	return calculator.Calculation(moveName, randDmgBonus)
}

func (b *Battle) p1CommandableMoveNames() bp.MoveNames {
	if b.SelfFighters[0].IsFaint() {
		return bp.MoveNames{}
	}
	
	if b.OpponentFighters[0].IsFaint() {
		return bp.MoveNames{}
	}

	names := make(bp.MoveNames, 0, bp.MAX_MOVESET)
	for moveName, pp := range b.SelfFighters[0].Moveset {
		if pp.Current > 0 {
			names = append(names, moveName)
		}
	}
	return names
}

func (b *Battle) p2CommandableMoveNames() bp.MoveNames {
	bv := b.SwapPlayers()
	names := bv.p1CommandableMoveNames()
	return names
}

func (b *Battle) SeparateCommandableMoveNames() bp.MoveNamess {
	p1 := b.p1CommandableMoveNames()
	p2 := b.p2CommandableMoveNames()
	return bp.MoveNamess{p1, p2}
}

func (b Battle) CommandMove(moveName bp.MoveName, r *rand.Rand) (Battle, error) {
	if b.SelfFighters[0].IsFaint() {
		return b, nil
	}

	if b.OpponentFighters[0].IsFaint() {
		return b, nil
	}

	if _, ok := b.SelfFighters[0].Moveset[moveName]; !ok {
		msg := fmt.Sprintf("%s は %s を 繰り出そうとしたが、覚えていない", b.SelfFighters[0].Name.ToString(), moveName.ToString())
		return b, fmt.Errorf(msg)
	}

	moveset := b.SelfFighters[0].Moveset.Clone()
	moveset[moveName].Current -= 1

	b.Observer(&b, BEFORE_MOVE_USE_STEP)
	b.SelfFighters[0].Moveset = moveset
	b.Observer(&b, AFTER_MOVE_USE_STEP)

	dmg := b.CalcDamage(moveName, r)
	dmg = omwmath.Min(dmg, b.OpponentFighters[0].CurrentHP)

	b.Observer(&b, BEFORE_MOVE_DAMAGE_STEP)
	b.OpponentFighters[0].CurrentHP -= dmg
	b.Observer(&b, AFTER_MOVE_DAMAGE_STEP)

	return b, nil
}

func (b *Battle) p1SwitchablePokeNames() bp.PokeNames {
	//相手だけ瀕死状態ならば、自分は行動出来ない。
	if b.SelfFighters[0].CurrentHP > 0 && b.OpponentFighters[0].CurrentHP <= 0 {
		return bp.PokeNames{}
	}
	names := make(bp.PokeNames, 0, FIGHTER_NUM-1)
	for _, pokemon := range b.SelfFighters[1:] {
		if pokemon.CurrentHP > 0 {
			names = append(names, pokemon.Name)
		}
	}
	return names
}

func (b *Battle) p2SwitchablePokeNames() bp.PokeNames {
	bv := b.SwapPlayers()
	names := bv.p1SwitchablePokeNames()
	return names
}

func (b *Battle) SeparateSwitchablePokeNames() bp.PokeNamess {
	p1 := b.p1SwitchablePokeNames()
	p2 := b.p2SwitchablePokeNames()
	return bp.PokeNamess{p1, p2}
}

func (b Battle) Switch(pokeName bp.PokeName) (Battle, error) {
	idx := b.SelfFighters.IndexByName(pokeName)
	if idx == -1 {
		name := bp.POKE_NAME_TO_STRING[pokeName]
		msg := fmt.Sprintf("「%s]へ交代しようとしたが、Fightersの中に含まれていなかった。", name)
		return Battle{}, fmt.Errorf(msg)
	}

	if idx == 0 {
		name := bp.POKE_NAME_TO_STRING[pokeName]
		msg := fmt.Sprintf("「%s]へ交代しようとしたが、既に場に出ている。", name)
		return Battle{}, fmt.Errorf(msg)
	}
	
	b.Observer(&b, BEFORE_SWITCH_STEP)
	b.SelfFighters[0], b.SelfFighters[idx] = b.SelfFighters[idx], b.SelfFighters[0]
	b.Observer(&b, AFTER_SWITCH_STEP)
	return b, nil
}

func (b *Battle) Action(action Action, r *rand.Rand) (Battle, error) {
	if action.IsCommandMove() {
		return b.CommandMove(action.CmdMoveName, r)
	} else {
		return b.Switch(action.SwitchPokeName)
	}
}

func (b *Battle) IsActionFirst(p1Action, p2Action *Action, r *rand.Rand) bool {
	if p1Action.IsEmpty() {
		return false
	}

	if p2Action.IsEmpty() {
		return true
	}

	if p1Action.IsSwitch() && p2Action.IsCommandMove() {
		return true
	}

	if p1Action.IsCommandMove() && p2Action.IsSwitch() {
		return false
	}

	attacker := b.SelfFighters[0]
	defender := b.OpponentFighters[0]

	if attacker.Speed > defender.Speed {
		return true
	} else if attacker.Speed < defender.Speed {
		return false
	} else {
		return omwrand.Bool(r)
	}
}

func (b *Battle) SortActionsByOrder(p1Action, p2Action *Action, r *rand.Rand) Actions {
	if b.IsActionFirst(p1Action, p2Action, r) {
		return Actions{*p1Action, *p2Action}
	} else {
		return Actions{*p2Action, *p1Action}
	}
}

func Equal(b1, b2 *Battle) bool {
	return b1.SelfFighters.Equal(&b2.SelfFighters) && b1.OpponentFighters.Equal(&b2.OpponentFighters) && b1.Turn == b2.Turn
}

func IsEnd(b *Battle) (bool, []float64) {
	isP1AllFaint := b.SelfFighters.IsAllFaint()
	isP2AllFaint := b.OpponentFighters.IsAllFaint()

	if isP1AllFaint && isP2AllFaint {
		return true, []float64{0.5, 0.5}
	} else if isP1AllFaint {
		return true, []float64{0.0, 1.0}
	} else if isP2AllFaint {
		return true, []float64{1.0, 0.0}
	} else {
		return false, []float64{}
	}
}

func LegalSeparateActions(b *Battle) ActionSlices {
	separateMoveNames := b.SeparateCommandableMoveNames()
	separatePokeNames := b.SeparateSwitchablePokeNames()
	ret := make(ActionSlices, MAX_SIMULTANEOUS_ACTION_NUM)
	isPlayer1s := []bool{true, false}
	for playerI := range ret {
		isPlayer1 := isPlayer1s[playerI]
		moveNames := separateMoveNames[playerI]
		pokeNames := separatePokeNames[playerI]
		actions := make(Actions, 0, len(moveNames) + len(pokeNames))
		for _, name := range moveNames {
			actions = append(actions, Action{CmdMoveName:name, IsPlayer1:isPlayer1})
		}
		for _, name := range pokeNames {
			actions = append(actions, Action{SwitchPokeName:name, IsPlayer1:isPlayer1})
		}
		if len(actions) == 0 {
			actions = append(actions, Action{IsPlayer1:isPlayer1})
		}
		ret[playerI] = actions
	}
	return ret
}

func NewPushFunc(r *rand.Rand) func(Battle, Actions) (Battle, error) {
	return func(battle Battle, actions Actions) (Battle, error) {
		if len(actions) != 2 {
			return Battle{}, fmt.Errorf("len(actions) != 2 (NewPushFunc)")
		}

		for actions[0].IsPlayer1 == actions[1].IsPlayer1 {
			return Battle{}, fmt.Errorf("プレイヤー1もしくはプレイヤー2が連続で行動しようとした。(actions[0].IsPlayer1 == actions[1].IsPlayer1)")
		}
		if actions.IsAllEmpty() {
			return Battle{}, fmt.Errorf("両プレイヤーのActionがEmptyになっているため、Pushできません。Emptyじゃないようにするには、Action.CmdMoveNameかAction.SwitchPokeNameのいずれかは、ゼロ値以外の値である必要があります。")
		}

		var err error
		sorted := battle.SortActionsByOrder(&actions[0], &actions[1], r)
		for i := range sorted {
			action := sorted[i]
			if action.CmdMoveName == bp.EMPTY_MOVE_NAME && action.SwitchPokeName == bp.EMPTY_POKE_NAME {
				continue
			}
			if action.IsPlayer1 {
				battle, err = battle.Action(action, r)
			} else {
				battle = battle.SwapPlayers()
				battle, err = battle.Action(action, r)
				battle = battle.SwapPlayers()
			}

			if battle.SelfFighters[0].IsFaint() {
				battle.Observer(&battle, SELF_FAINT_STEP)
			}

			if battle.OpponentFighters[0].IsFaint() {
				battle.Observer(&battle, OPPONENT_FAINT_STEP)
			}

			if err != nil {
				return Battle{}, err
			}
		}
		battle.Turn += 1
		return battle, nil
	}
}