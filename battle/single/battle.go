package single

import (
	"fmt"
	"math"
	"math/rand"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	omwmath "github.com/sw965/omw/math"
	"github.com/sw965/bippa/battle/dmgtools"
)

const (
	PLAYER_NUM = 2
	LEAD_NUM = 1
	BENCH_NUM = 2
	FIGHTERS_NUM = LEAD_NUM + BENCH_NUM
)

type Battle struct {
	SelfFighters bp.Pokemons
	OpponentFighters bp.Pokemons
	Turn int
	IsRealSelf bool
}

func (b Battle) Clone() Battle {
	return Battle{
		SelfFighters:b.SelfFighters.Clone(),
		OpponentFighters:b.OpponentFighters.Clone(),
		Turn:b.Turn,
		IsRealSelf:b.IsRealSelf,
	}
}

func (b Battle) SwapView() Battle {
	b.SelfFighters, b.OpponentFighters = b.OpponentFighters, b.SelfFighters
	b.IsRealSelf = !b.IsRealSelf
	return b
}

func (b *Battle) CalcDamage(moveName bp.MoveName, context *Context) int {
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
	return calculator.Calculation(moveName, context.DamageRandBonus())
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

	if len(names) == 0 {
		return bp.MoveNames{bp.STRUGGLE}
	} else {
		return names
	}
}

func (b *Battle) p2CommandableMoveNames() bp.MoveNames {
	bv := b.SwapView()
	names := bv.p1CommandableMoveNames()
	return names
}

func (b *Battle) SeparateCommandableMoveNames() bp.MoveNamess {
	p1 := b.p1CommandableMoveNames()
	p2 := b.p2CommandableMoveNames()
	return bp.MoveNamess{p1, p2}
}

func (b Battle) CommandMove(moveName bp.MoveName, context *Context) (Battle, error) {
	if b.SelfFighters[0].IsFaint() {
		return b, nil
	}

	if b.OpponentFighters[0].IsFaint() {
		return b, nil
	}

	pp, ok := b.SelfFighters[0].Moveset[moveName]
	if !ok {
		msg := fmt.Sprintf("%s は %s を 繰り出そうとしたが、覚えていない", b.SelfFighters[0].Name.ToString(), moveName.ToString())
		return b, fmt.Errorf(msg)
	}

	if pp.Current <= 0 {
		msg := fmt.Sprintf("%s は %s を 繰り出そうとしたが、PPが0", b.SelfFighters[0].Name.ToString(), moveName.ToString())
		return b, fmt.Errorf(msg)
	}

	b.SelfFighters = b.SelfFighters.Clone()
	b.OpponentFighters = b.OpponentFighters.Clone()
	
	b.SelfFighters[0].Moveset[moveName].Current -= 1
	context.Observer(&b, MOVE_USE_EVENT)

	dmg := b.CalcDamage(moveName, context)
	dmg = omwmath.Min(dmg, b.OpponentFighters[0].CurrentHP)

	b.OpponentFighters[0].CurrentHP -= dmg
	context.Observer(&b, OPPONENT_DAMAGE_EVENT)

	if moveName == bp.STRUGGLE {
		dmg := int(math.Round(float64(b.SelfFighters[0].CurrentHP) * 0.25))
		b.SelfFighters[0].CurrentHP -= dmg
		context.Observer(&b, RECOIL_EVENT)
	}
	return b, nil
}

func (b *Battle) p1SwitchablePokeNames() bp.PokeNames {
	//相手だけ瀕死状態ならば、自分は行動出来ない。
	if b.SelfFighters[0].CurrentHP > 0 && b.OpponentFighters[0].CurrentHP <= 0 {
		return bp.PokeNames{}
	}
	names := make(bp.PokeNames, 0, FIGHTERS_NUM)
	for _, pokemon := range b.SelfFighters[1:] {
		if pokemon.CurrentHP > 0 {
			names = append(names, pokemon.Name)
		}
	}
	return names
}

func (b *Battle) p2SwitchablePokeNames() bp.PokeNames {
	bv := b.SwapView()
	names := bv.p1SwitchablePokeNames()
	return names
}

func (b *Battle) SeparateSwitchablePokeNames() bp.PokeNamess {
	p1 := b.p1SwitchablePokeNames()
	p2 := b.p2SwitchablePokeNames()
	return bp.PokeNamess{p1, p2}
}

func (b Battle) Switch(pokeName bp.PokeName, context *Context) (Battle, error) {
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

	selfFighters := b.SelfFighters.Clone()
	selfFighters[0], selfFighters[idx] = selfFighters[idx], selfFighters[0]
	b.SelfFighters = selfFighters
	context.Observer(&b, SWITCH_EVENT)
	return b, nil
}

func (b *Battle) Action(action Action, context *Context) (Battle, error) {
	if action.IsCommandMove() {
		return b.CommandMove(action.CmdMoveName, context)
	} else {
		return b.Switch(action.SwitchPokeName, context)
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

func (b *Battle) ToEasyRead() EasyReadBattle {
	return EasyReadBattle{
		SelfFighters:b.SelfFighters.ToEasyRead(),
		OpponentFighters:b.OpponentFighters.ToEasyRead(),
		Turn:b.Turn,
		IsRealSelf:b.IsRealSelf,
	}
}