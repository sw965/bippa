package single

import (
	"fmt"
	omwmath "github.com/sw965/omw/math"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle/dmgtools"
)

func AttackMove(battle *Battle, action *SoloAction, context *Context) (bp.PokemonPointers, error) {
	targetPokemons, err := battle.TargetPokemonPointers(action, context.Rand)
	if err != nil {
		return bp.PokemonPointers{}, err
	}
	targetN := len(targetPokemons)
	if targetN == 0 {
		return bp.PokemonPointers{}, err
	}

	var isSingleDmg bool
	if action.MoveName == bp.SELF_DESTRUCT || action.MoveName == bp.EXPLOSION {
		isSingleDmg = targetN <= 2
	} else {
		isSingleDmg = targetN == 1
	}

	moveData := bp.MOVEDEX[action.MoveName]
	faintedCount := 0
	for i, defender := range targetPokemons {
		isHit, err := omwrand.IsPercentageMet(moveData.Accuracy, context.Rand)
		if err != nil {
			return bp.PokemonPointers{}, err
		}

		if !isHit {
			continue
		}

		if i == (targetN-1) {
			isSingleDmg = isSingleDmg || (faintedCount == i)
		}

		if moveData.Type == bp.WATER && defender.Ability == bp.DRY_SKIN {
			fmt.Println("ok")
			err := defender.AddCurrentHP(int(float64(defender.CurrentHP) * 0.25))
			if err != nil {
				return bp.PokemonPointers{}, err
			}
			continue
		}

		critRank := moveData.CriticalRank
		isCrit, err := dmgtools.IsCritical(critRank, context.Rand)
		if err != nil {
			return bp.PokemonPointers{} , err
		}
		dmg, err := battle.CalcDamage(action, defender, isCrit, isSingleDmg, context)
		if err != nil {
			return bp.PokemonPointers{}, err
		}

		dmg = omwmath.Min(dmg, defender.CurrentHP)
		err = defender.SubCurrentHP(dmg)
		if err != nil {
			return bp.PokemonPointers{}, err
		}

		if defender.IsFainted() {
			faintedCount += 1 
		}
	}
	return targetPokemons, nil
}

//10まんボルト
func Thunderbolt(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.THUNDERBOLT {
		return fmt.Errorf("HammerArmに渡されたAction.MoveNameがTHUNDERBOLTではない。")
	}
	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		return err
	}

	if len(targetPokemons) == 0 {
		return nil
	}

	if len(targetPokemons) != 1 {
		return fmt.Errorf("ライブラリのバグなので報告して。bp.Thunderbolt")
	}

	targetPokemon := targetPokemons[0]
	if targetPokemon.IsFainted() {
		return nil
	}

	ok, err := omwrand.IsPercentageMet(10, context.Rand)
	if ok {
		targetPokemon.StatusAilment = bp.PARALYSIS
	}
	return err
}

//アームハンマー
func HammerArm(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.HAMMER_ARM {
		return fmt.Errorf("HammerArmに渡されたAction.MoveNameがHAMMER_ARMではない。")
	}
	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		return err
	}

	if len(targetPokemons) == 0 {
		return nil
	}

	if len(targetPokemons) != 1 {
		return fmt.Errorf("ライブラリのバグなので報告して。bp.HammerArm")
	}

	targetPokemon := targetPokemons[0]
	if targetPokemon.IsFainted() {
		return nil
	}

	targetPokemon.Rank.Fluctuation(&bp.RankStat{Speed:-1})
	return nil
}

//ストーンエッジ
func StoneEdge(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.STONE_EDGE {
		return fmt.Errorf("StoneEdgeに渡されたAction.MoveNameがSTONE_EDGEではない。")
	}

	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		return err
	}

	if len(targetPokemons) == 0 {
		return nil
	}

	if len(targetPokemons) != 1 {
		return fmt.Errorf("ライブラリのバグなので報告して。bp.StoneEdge")
	}
	return nil
}

//なみのり
func Surf(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.SURF {
		return fmt.Errorf("Surfに渡されたAction.MoveNameがSURFではない。")
	}

	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		panic(err)
	}

	if len(targetPokemons) == 0 {
		return nil
	}
	return nil
}