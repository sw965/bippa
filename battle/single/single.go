package single

import (
	"fmt"
	"math/rand"
	"github.com/sw965/crow/ucb"
	bp "github.com/sw965/bippa"
	orand "github.com/sw965/omw/rand"
	omath "github.com/sw965/omw/math"
	"github.com/sw965/bippa/dmgtools"
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/crow/mcts/duct"
	"golang.org/x/exp/slices"
	"github.com/sw965/crow/tensor"
	"github.com/sw965/crow/model"
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

type Actionss []Actions

type Battle struct {
	P1Fighters Fighters
	P2Fighters Fighters
	Turn int
	IsSelf bool
}

func (b Battle) SwapPlayers() Battle {
	b.P1Fighters, b.P2Fighters = b.P2Fighters, b.P1Fighters
	b.IsSelf = !b.IsSelf
	return b
}

func (b *Battle) CalcDamage(moveName bp.MoveName, randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) int {
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
			PokeName:defender.Name,
			Level:defender.Level,
			Def:defender.Def,
			SpDef:defender.SpDef,
		},
	}
	randDmgBonus := orand.Choice(randDmgBonuses, r)
	return calculator.Calculation(randDmgBonus)
}

func (b *Battle) p1CommandableMoveNames() bp.MoveNames {
	if b.P1Fighters[0].IsFaint() {
		return bp.MoveNames{}
	}
	
	if b.P2Fighters[0].IsFaint() {
		return bp.MoveNames{}
	}

	names := make(bp.MoveNames, 0, bp.MAX_MOVESET_NUM)
	for moveName, pp := range b.P1Fighters[0].Moveset {
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

func (b *Battle) CommandableMoveNamess() bp.MoveNamess {
	p1 := b.p1CommandableMoveNames()
	p2 := b.p2CommandableMoveNames()
	return bp.MoveNamess{p1, p2}
}

func (b Battle) CommandMove(moveName bp.MoveName, randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) Battle {
	if b.P1Fighters[0].IsFaint() {
		return b
	}
	if b.P2Fighters[0].IsFaint() {
		return b
	}
	moveset := b.P1Fighters[0].Moveset.Clone()
	//moveset[moveName].Current -= 1

	b.P1Fighters[0].Moveset = moveset
	dmg := b.CalcDamage(moveName, randDmgBonuses, r)
	dmg = omath.Min(dmg, b.P2Fighters[0].CurrentHP)
	b.P2Fighters[0].CurrentHP -= dmg
	return b
}

func (b *Battle) p1SwitchablePokeNames() bp.PokeNames {
	//相手だけ瀕死状態ならば、自分は行動出来ない。
	if b.P1Fighters[0].CurrentHP > 0 && b.P2Fighters[0].CurrentHP <= 0 {
		return bp.PokeNames{}
	}
	names := make(bp.PokeNames, 0, FIGHTER_NUM-1)
	for _, pokemon := range b.P1Fighters[1:] {
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

func (b *Battle) SwitchablePokeNamess() bp.PokeNamess {
	p1 := b.p1SwitchablePokeNames()
	p2 := b.p2SwitchablePokeNames()
	return bp.PokeNamess{p1, p2}
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

func (b *Battle) Action(action Action, randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) (Battle, error) {
	if action.IsCommandMove() {
		return b.CommandMove(action.CmdMoveName, randDmgBonuses, r), nil
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

	attacker := b.P1Fighters[0]
	defender := b.P2Fighters[0]

	if attacker.Speed > defender.Speed {
		return true
	} else if attacker.Speed < defender.Speed {
		return false
	} else {
		return orand.Bool(r)
	}
}

func (b *Battle) SortActionsByOrder(p1Action, p2Action *Action, r *rand.Rand) Actions {
	if b.IsActionFirst(p1Action, p2Action, r) {
		return Actions{*p1Action, *p2Action}
	} else {
		return Actions{*p2Action, *p1Action}
	}
}

func (b *Battle) EndLeafNodeEvalYs() (duct.LeafNodeEvalYs, error) {
	isP1AllFaint := b.P1Fighters.IsAllFaint()
	isP2AllFaint := b.P2Fighters.IsAllFaint()
	if isP1AllFaint && isP2AllFaint {
		return duct.LeafNodeEvalYs{0.5, 0.5}, nil
	} else if isP1AllFaint {
		return duct.LeafNodeEvalYs{0.0, 1.0}, nil
	} else if isP2AllFaint {
		return duct.LeafNodeEvalYs{1.0, 0.0}, nil
	} else {
		return duct.LeafNodeEvalYs{}, fmt.Errorf("ゲームが終わっていない")
	}
}

func Equal(b1, b2 *Battle) bool {
	return b1.P1Fighters.Equal(&b2.P1Fighters) && b1.P2Fighters.Equal(&b2.P2Fighters) && b1.Turn == b2.Turn
}

func IsEnd(b *Battle) bool {
	return b.P1Fighters.IsAllFaint() || b.P2Fighters.IsAllFaint()
}

func LegalActionss(b *Battle) Actionss {
	moveNamess := b.CommandableMoveNamess()
	pokeNamess := b.SwitchablePokeNamess()
	actionss := make(Actionss, MAX_SIMULTANEOUS_ACTION_NUM)
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
		if len(actions) == 0 {
			actions = append(actions, Action{IsPlayer1:isPlayer1})
		}
		actionss[playerI] = actions
	}
	return actionss
}

func Push(randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) func(Battle, Actions) (Battle, error) {
	return func(battle Battle, actions Actions) (Battle, error) {
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
				battle, err = battle.Action(action, randDmgBonuses, r)
			} else {
				battle = battle.SwapPlayers()
				battle, err = battle.Action(action, randDmgBonuses, r)
				battle = battle.SwapPlayers()
			}
			if err != nil {
				return Battle{}, err
			}
		}
		battle.Turn += 1
		return battle, nil
	}
}

func NewGame(randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) simultaneous.Game[Battle, Actionss, Actions, Action] {
	return simultaneous.Game[Battle, Actionss, Actions, Action]{
		Equal:Equal,
		IsEnd:IsEnd,
		LegalActionss:LegalActionss,
		Push:Push(randDmgBonuses, r),
	}
}

func NewMCTSRandPlayout(ucbFunc ucb.Func, randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) duct.MCTS[Battle, Actionss, Actions, Action] {
	game := NewGame(randDmgBonuses, r)
	game.SetRandActionPlayer(r)

	leafNodeEvalsFunc := func(battle *Battle) duct.LeafNodeEvalYs {
		endBattle, err := game.Playout(*battle)
		if err != nil {
			panic(err)
		}
		ys, err := endBattle.EndLeafNodeEvalYs()
		if err != nil {
			panic(err)
		}
		return ys
	}

	mcts := duct.MCTS[Battle, Actionss, Actions, Action]{
		UCBFunc:ucbFunc,
		Game:game,
		LeafNodeEvalsFunc:leafNodeEvalsFunc,
		NextNodesCap:64,
	}
	mcts.SetUniformActionPoliciesFunc()
	return mcts
}

func NewMCTSNNEval(ucbFunc ucb.Func, nn *model.SequentialInputOutput1D, randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) duct.MCTS[Battle, Actionss, Actions, Action] {
	game := simultaneous.Game[Battle, Actionss, Actions, Action]{
		Equal:Equal,
		IsEnd:IsEnd,
		LegalActionss:LegalActionss,
		Push:Push(randDmgBonuses, r),
	}
	game.SetRandActionPlayer(r)

	leafNodeEvalsFunc := func(battle *Battle) duct.LeafNodeEvalYs {
		isP1AllFaint := battle.P1Fighters.IsAllFaint()
		isP2AllFaint := battle.P2Fighters.IsAllFaint()
		if isP1AllFaint && isP2AllFaint {
			return duct.LeafNodeEvalYs{0.5, 0.5}
		} else if isP1AllFaint {
			return duct.LeafNodeEvalYs{0.0, 1.0}
		} else if isP2AllFaint {
			return duct.LeafNodeEvalYs{1.0, 0.0}
		} else {
			x := battle.ToIndexFeature(15000)
			y, err := nn.Predict(x)
			if err != nil {
				panic(err)
			}
			v := duct.LeafNodeEvalY(y[0])
			return duct.LeafNodeEvalYs{v, 1.0 - v}
		}
	}

	mcts := duct.MCTS[Battle, Actionss, Actions, Action]{
		UCBFunc:ucbFunc,
		Game:game,
		LeafNodeEvalsFunc:leafNodeEvalsFunc,
		NextNodesCap:64,
	}
	mcts.SetUniformActionPoliciesFunc()
	return mcts
}

func (b *Battle) ToIndexFeature(baseIndex int) tensor.D1 {
	makePowerIndexFeature := func(pokemon *bp.Pokemon, category bp.MoveCategory) tensor.D1 {
		var stat int
		if category == bp.PHYSICS {
			stat = pokemon.Atk
		} else {
			stat = pokemon.SpAtk
		}
		ret := make(tensor.D1, len(bp.ALL_TYPES))
		for i, t := range bp.ALL_TYPES {
			max := 0.0
			for moveName, pp := range pokemon.Moveset {
				moveData := bp.MOVEDEX[moveName]
				if pp.Current <= 0 || moveData.Type != t || moveData.Category != category {
					continue
				}
				index := float64(moveData.Power) * float64(stat) * float64(moveData.Accuracy / 100.0)
				if index > max {
					max = index
				}
			}
			ret[i] = max / float64(baseIndex)
		}
		return ret
	}

	makeDurabilityIndexFeature := func(pokemon *bp.Pokemon, category bp.MoveCategory) tensor.D1 {
		var stat int
		if category == bp.PHYSICS {
			stat = pokemon.Def
		} else {
			stat = pokemon.SpDef
		}
		pokeData := bp.POKEDEX[pokemon.Name]
		ret := make(tensor.D1, len(bp.ALL_TYPES))
		for i, t := range bp.ALL_TYPES {
			if !slices.Contains(pokeData.Types, t) {
				continue
			}
			ret[i] = float64(pokemon.CurrentHP) * float64(stat) / float64(baseIndex)
		}
		return ret
	}

	ret := make(tensor.D1, 0, len(bp.ALL_TYPES) * 4 * FIGHTER_NUM * 2)
	add := func(fighters Fighters) {
		for _, pokemon := range fighters {
			ret = append(ret, makePowerIndexFeature(&pokemon, bp.PHYSICS)...)
			ret = append(ret, makePowerIndexFeature(&pokemon, bp.SPECIAL)...)
			ret = append(ret, makeDurabilityIndexFeature(&pokemon, bp.PHYSICS)...)
			ret = append(ret, makeDurabilityIndexFeature(&pokemon, bp.SPECIAL)...)
		}
	}
	add(b.P1Fighters)
	add(b.P2Fighters)
	return ret
}