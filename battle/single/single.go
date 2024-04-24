package single

import (
	"fmt"
	"math/rand"
	bp "github.com/sw965/bippa"
	"github.com/sw965/omw"
	"github.com/sw965/bippa/dmgtools"
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/crow/mcts/dpuct"
)

const (
	FIGHTER_NUM = 3
)

type Fighters [FIGHTER_NUM]bp.Pokemon

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

func (a *Action) IsCommandMove() bool {
	return a.CmdMoveName != bp.EMPTY_MOVE_NAME
}

func (a *Action) IsSwitch() bool {
	return a.SwitchPokeName != bp.EMPTY_POKE_NAME
}

type Actions []Action

const (
	SINGLE_BATTLE_MAX_SIMULTANEOUS_ACTION_NUM = 2
)

type Actionss []Actions

type Battle struct {
	P1Fighters Fighters
	P2Fighters Fighters
}

func (b *Battle) SwapPlayers() Battle {
    return Battle{P1Fighters: b.P2Fighters, P2Fighters: b.P1Fighters}
}

func (b *Battle) CalcDamage(moveName bp.MoveName, randDmgBonus float64) int {
	attacker := b.P1Fighters[0]
	defender := b.P2Fighters[0]

	calculator := dmgtools.Calculator{
		dmgtools.Attacker{
			PokeName:attacker.Name,
			Level:attacker.Level,
			Atk:attacker.Atk,
			SpAtk:attacker.SpAtk,
			MoveName:moveName,
		},
		dmgtools.Defender{
			PokeName:attacker.Name,
			Level:defender.Level,
			Def:defender.Def,
			SpDef:defender.SpDef,
		},
	}
	return calculator.Execute(randDmgBonus)
}

func (b *Battle) CommandableMoveNamess() bp.MoveNamess {
	namess := make(bp.MoveNamess, SINGLE_BATTLE_MAX_SIMULTANEOUS_ACTION_NUM)
	for i := range namess {
		namess[i] = make(bp.MoveNames, 0, bp.MAX_MOVESET_NUM)
	}

	appendCommandableMoveNames := func(sellFighters, opponentFighters *Fighters, idx int) {
		if opponentFighters[0].CurrentHP > 0 {
			for moveName, pp := range sellFighters[0].Moveset {
				if pp.Current > 0 {
					namess[idx] = append(namess[idx], moveName)
				}
			}
		}
	}

	appendCommandableMoveNames(&b.P1Fighters, &b.P2Fighters, 0)
	appendCommandableMoveNames(&b.P2Fighters, &b.P1Fighters, 1)
	return namess
}

func (b Battle) CommandMove(moveName bp.MoveName, r *rand.Rand) Battle {
	randDmgBonus := omw.RandChoice(dmgtools.RANDOM_BONUS, r)
	dmg := b.CalcDamage(moveName, randDmgBonus)
	dmg = omw.Min(dmg, b.P2Fighters[0].CurrentHP)
	moveset := b.P1Fighters[0].Moveset.Clone()
	fmt.Println(moveset, moveset[moveName], b.P1Fighters[0].Name)
	moveset[moveName].Current -= 1
	b.P1Fighters[0].Moveset = moveset
	b.P2Fighters[0].CurrentHP -= dmg
	return b
}

func (b *Battle) SwitchablePokeNamess() bp.PokeNamess {
	namess := make(bp.PokeNamess, SINGLE_BATTLE_MAX_SIMULTANEOUS_ACTION_NUM)
	appendSwitchablePokeNames := func(fighters *Fighters, idx int) {
		for _, pokemon := range fighters {
			if pokemon.CurrentHP > 0 {
				namess[idx] = append(namess[idx], pokemon.Name)
			}
		}
	}
	appendSwitchablePokeNames(&b.P1Fighters, 0)
	appendSwitchablePokeNames(&b.P2Fighters, 1)
	return namess
}

func (b Battle) Switch(pokeName bp.PokeName) (Battle, error) {
	idx := b.P1Fighters.IndexByName(pokeName)
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

	b.P1Fighters[0], b.P1Fighters[idx] = b.P1Fighters[idx], b.P1Fighters[0]
	return b, nil
}

func (b *Battle) Action(action Action, r *rand.Rand) (Battle, error) {
	if action.IsCommandMove() {
		return b.CommandMove(action.CmdMoveName, r), nil
	} else {
		return b.Switch(action.SwitchPokeName)
	}
}

func (b *Battle) IsActionFirst(p1Action, p2Action *Action, r *rand.Rand) bool {
	if p1Action.IsSwitch() && p2Action.IsCommandMove() {
		return true
	}

	if p1Action.IsCommandMove() && p2Action.IsSwitch() {
		return false
	}

	attacker := b.P1Fighters[0]
	defender := b.P2Fighters[0]

	if attacker.Speed > defender.Speed {
		return true
	} else if attacker.Speed < defender.Speed {
		return false
	} else {
		return omw.RandBool(r)
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
	return b1.P1Fighters.Equal(&b2.P1Fighters) && b1.P2Fighters.Equal(&b2.P2Fighters) 
}

func IsEnd(b *Battle) bool {
	return b.P1Fighters.IsAllFaint() || b.P2Fighters.IsAllFaint()
}

func LegalActionss(b *Battle) Actionss {
	moveNamess := b.CommandableMoveNamess()
	pokeNamess := b.SwitchablePokeNamess()
	actionss := make(Actionss, SINGLE_BATTLE_MAX_SIMULTANEOUS_ACTION_NUM)
	isPlayer1s := []bool{true, false}
	for playerI := range actionss {
		isPlayer1 := isPlayer1s[playerI]
		moveNames := moveNamess[playerI]
		pokeNames := pokeNamess[playerI]
		actions := make(Actions, 0, len(moveNames) + len(pokeNames))
		for _, name := range moveNames {
			actions = append(actions, Action{CmdMoveName:name, IsPlayer1:isPlayer1})
		}
		for _, name := range pokeNames {
			actions = append(actions, Action{SwitchPokeName:name, IsPlayer1:isPlayer1})
		}
		actionss[playerI] = actions
	}
	return actionss
}

func Push(r *rand.Rand) func(Battle, Actions) (Battle, error) {
	return func(battle Battle, actions Actions) (Battle, error) {
		var err error
		sorted := battle.SortActionsByOrder(&actions[0], &actions[1], r)
		for i := range sorted {
			action := sorted[i]
			if action.IsPlayer1 {
				battle, err = battle.Action(action, r)
			} else {
				battle = battle.SwapPlayers()
				battle, err = battle.Action(action, r)
				battle = battle.SwapPlayers()
			}
			if err != nil {
				return Battle{}, err
			}
		}
		return battle, nil
	}
}

func NewMCTS(r *rand.Rand) dpuct.MCTS[Battle, Actionss, Actions, Action] {
	game := simultaneous.Game[Battle, Actionss, Actions, Action]{
		Equal:Equal,
		IsEnd:IsEnd,
		LegalActionss:LegalActionss,
		Push:Push(r),
	}
	game.SetRandomActionPlayer(r)

	leafNodeEvalsFunc := func(battle *Battle) dpuct.LeafNodeEvalYs {
		b1 := battle.P1Fighters.IsAllFaint()
		b2 := battle.P2Fighters.IsAllFaint()
		if b1 && b2 {
			return dpuct.LeafNodeEvalYs{0.5, 0.5}
		} else if b1 {
			return dpuct.LeafNodeEvalYs{1.0, 0.0}
		} else {
			return dpuct.LeafNodeEvalYs{0.0, 1.0}
		}
	}

	mcts := dpuct.MCTS[Battle, Actionss, Actions, Action]{
		Game:game,
		LeafNodeEvalsFunc:leafNodeEvalsFunc,
	}
	mcts.SetUniformActionPoliciesFunc()
	return mcts
}