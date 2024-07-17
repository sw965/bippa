package single

import (
	"fmt"
	omwmath "github.com/sw965/omw/math"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle/dmgtools"
)

func Move(battle *Battle, action *SoloAction, context *Context) (bp.PokemonPointers, error) {
	attacker := &battle.SelfLeadPokemons[action.SrcIndex]
	if attacker.IsFlinch {
		return bp.PokemonPointers{}, nil
	}

	switch action.MoveName {
		//あまごい
		case RAIN_DANCE:
			battle.Weather = RAIN
			battle.WeatherRemainingTurn = 5
			return bp.PokemonPointers{}, nil
		//このゆびとまれ
		case FOLLOW_ME:
			attacker.IsFollowMe = true
			return bp.PokemonPointers{}, err
		//じばく
		case SELF_DESTRUCT:
			attacker.CurrentHP = 0
		//だいばくはつ
		case EXPLOSION:
			attacker.CurrentHP = 0
		//はらだいこ
		case BELLY_DRUM:
			attacker.Rank.Atk = MAX_RANK
	}

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
	moveHitHistory := make([]bool, len(targetPokemons)) 
	faintedCount := 0

	for i, defender := range targetPokemons {
		switch action.MoveName {
			// じこあんじ
			case RECOVER:
				attacker.Rank = defender.Rank
				break
		}

		if defender.IsProtect {
			continue
		}

		isHit, err := omwrand.IsPercentageMet(moveData.Accuracy, context.Rand)
		if err != nil {
			return bp.PokemonPointers{}, err
		}

		if !isHit {
			continue
		}

		switch action.MoveName {
			// がむしゃら
			case bp.ENDEAVOR:
				defender.CurrentHP = attacker.CurrentHP
				break
			// さいみんじゅつ
			case bp.HYPNOSIS:
				defender.StatusAilment = SLEEP
				break
			//ちょうはつ
			case TAUNT:
				defender.IsTaunt = true
				break
			//でんじは
			case THUNDER_WAVE:
				if !omwslices.Contains(defender.Types, GROUND) {
					defender.StatusAilment = PARALYSIS
				}
				break
		}

		if i == (targetN-1) {
			isSingleDmg = isSingleDmg || (faintedCount == i)
		}

		if moveData.Type == bp.WATER && defender.Ability == bp.DRY_SKIN {
			err := defender.AddCurrentHP(int(float64(defender.CurrentHP) * 0.25))
			if err != nil {
				return bp.PokemonPointers{}, err
			}
			continue
		}

		//ふいうち
		if action.MoveName == SUCKER_PUNCH {
			//相手が行動を終えていたら、失敗
			if defender.IsThisTurnActed {
				break
			}

			if data, ok := bp.MOVEDEX[defender.ThisTurnCommandMoveName]; !ok {
				msg := fmt.Sprintf(
					"%s の ThisTurnCommandMoveNameが %s になっているが、bp.MOVEDEXに含まれてない",
					defender.Name.ToString(), defender.ThisTurnCommandMoveName.ToString())
				return bp.PokemonPointers{}, fmt.Errorf("%v")
			} else {
				//相手が変化技を選択していたならば、失敗
				if data.Category.IsStatus() {
					break
				}
			}
		}

		//命中が確定
		hitHistory[i] = true

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

	switch action.MoveName {
		//アームハンマー
		case HAMMER_ARM:
			attacker.Rank = attacker.Rank.Fluctuation(&RankStat{Speed:-1})
	}

	for i, pokemon := range targetPokemons {
		if pokemon.IsFainted() {
			continue
		}

		if !moveHitHistory[i] {
			continue
		}

		switch action.MoveName {
			//10まんボルト
			case THUNDERBOLT:
				if omwrand.IsPercentageMet(10, context.Rand) {
					pokemon.StatusAilment = bp.PARALYSIS
				}
			//れいとうビーム
			case ICE_BEAM:
				if omwrand.IsPercentageMet(10, context.Rand) {
					pokemon.StatusAilment = bp.FREEZE
				}
			//かみくだく
			case CRUNCH:
				if ok, err := omwrand.IsPercentageMet(20, context.Rand); err {
					return bp.PokemonPointers{}, err
				} else {
					if ok {
						pokemon.Rank = pokemon.Rank.Fluctuation(&RankStat{Def:-1})
					}
				}
			//こごえるかぜ
			case ICY_WIND:
				pokemon.Rank = pokemon.Rank.Fluctuation(&RankStat{Speed:-1}) 
			//たきのぼり
			case WATERFALL:
				if ok, err := omwrand.IsPercentageMet(20, context.Rand); err {
					return bp.PokemonPointers{}, err
				} else {
					pokemon.IsFlinch = ok
				}
			//ねこだまし
			case FAKE_OUT:
				pokemon.IsFaint = true
			//ねっぷう
			case HEAT_WAVE:
				if ok, err := omwrand.IsPercentageMet(10, context.Rand); err {
					return bp.PokemonPointers{}, err
				} else {
					if ok {
						pokemon.StatusAilment = BURN
					}
				}
			//ほのおのパンチ
			case FIRE_PUNCH:
				if ok, err := omwrand.IsPercentageMet(10, context.Rand); err {
					return bp.PokemonPointers{}, err
				} else {
					if ok {
						pokemon.StatusAilment = BURN
					}
				}
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

	attacker := &battle.SelfLeadPokemons[action.SrcIndex]
	attacker.Rank = attacker.Rank.Fluctuation(&bp.RankStat{Speed:-1})
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

//れいとうビーム
func IceBeam(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.ICE_BEAM {
		return fmt.Errorf("IceBeamに渡されたAction.MoveNameがICE_BEAMではない")
	}

	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		return err
	}

	if len(targetPokemons) == 0 {
		return nil
	}

	if len(targetPokemons) != 1 {
		return fmt.Errorf("ライブラリのバグなので報告して。IceBeam")
	}

	if omwrand.IsPercentageMet(10, context.Rand) {
		targetPokemons[0].StatusAilment = FREEZE
	}
	return nil
}

//あまごい
func RainDance(battle *Battle, action *SoloAction, context *Context) error {
	battle.Weather = RAIN
	battle.WeatherRemainingTurn = 5
}

//いわなだれ
func RockSlide(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.ICE_BEAM {
		return fmt.Errorf("RockSlideに渡されたAction.MoveNameがbp.ICE_BEAMではない")
	}

	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		return err
	}

	if len(targetPokemons) == 0 {
		return nil
	}

	for _, pokemon := range targetPokemons {
		if !pokemon.IsFainted() && omwrand.IsPercentageMet(30, context.Rand) {
			pokemon.IsFlinch = true
		}
	}
	return nil
}

//おんがえし
func Return(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.RETURN {
		return fmt.Errorf("Returnに渡されたAction.MoveNameがbp.RETURNではない")
	}

	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		return err
	}

	if len(targetPokemons) == 0 {
		return nil
	}

	if len(targetPokemons) != 1 {
		return fmt.Errorf("ライブラリのバグなので報告して。Return")
	}
	return nil
}

//かみくだく
func Crunch(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.CRUNCH {
		return fmt.Errorf("Crunchに渡されたAction.MoveNameがbp.CRUNCHではない")
	}

	targetPokemons, err := AttackMove(battle, action, context)
	if err != nil {
		return err
	}

	if len(targetPokemons) == 0 {
		return nil
	}

	if len(targetPokemons) != 1 {
		return fmt.Errorf("ライブラリのバグなので報告して。Crunch")
	}

	targetPokemon := targetPokemons[0]
	if targetPokemon.IsFainted() && omwrand.IsPercentageMet(20, r) {
		targetPokemon.Rank.Def = targetPokemon.Rank.Fluctuation(&RankStat{Def:-1})
	}
	return nil
}



