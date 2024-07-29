package single

import (
	"fmt"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	//"github.com/sw965/bippa/battle/dmgtools"
	omwmath "github.com/sw965/omw/math"
)

/*
	第4世代の技データ
	(あ行～た行) https://yakkun.com/dp/waza_list.htm#k
	(な行～わ行) https://yakkun.com/dp/waza_list2.htm
*/

type Status func(*Battle, *bp.Pokemon, *bp.Pokemon, *Context) error

func EmptyStatus(_ *Battle, _, _ *bp.Pokemon, _ *Context) error {
	return nil
}

// https://wiki.xn--rckteqa2e.com/wiki/%E8%BF%BD%E5%8A%A0%E5%8A%B9%E6%9E%9C
type AdditionalEffect func(*bp.Pokemon) error

func EmptyAdditionalEffect(_ *bp.Pokemon) error {
	return nil
}

func moveHelper(
	battle *Battle, action *SoloAction, context *Context,
	status Status, self, opponent AdditionalEffect,
	) error {
	src := &battle.SelfLeadPokemons[action.SrcIndex]

	if src.IsFainted() {
		return nil
	}

	if src.IsFlinchState {
		return nil
	}

	moveData := bp.MOVEDEX[action.MoveName]

	targetPokemons := battle.TargetPokemonPointers(action, context)
	targetPokemons.SortBySpeed()
	targetN := len(targetPokemons)

	var isSingleDmg bool
	if action.MoveName == bp.SELF_DESTRUCT || action.MoveName == bp.EXPLOSION {
		isSingleDmg = targetN <= 2
	} else {
		isSingleDmg = targetN == 1
	}

	var isBodyAttack bool
	var isTargetFainted bool
	var isContinue bool
	faintedCount := 0

	if action.MoveName == bp.SELF_DESTRUCT || action.MoveName == bp.EXPLOSION {
		src.CurrentHP = 0
	}

	fs := []func(*bp.Pokemon , *bp.Pokemon) error {
		//命中
		func(_, _ *bp.Pokemon) error {
			isHit, err := omwrand.IsPercentageMet(moveData.Accuracy, context.Rand)
			isContinue = isHit
			return err
		},

		//かんそうはだ
		func(_, target *bp.Pokemon) error {
			if moveData.Type == bp.WATER && target.Ability == bp.DRY_SKIN {
				err := target.AddCurrentHP(int(float64(target.CurrentHP) * 0.25))
				if err != nil {
					return err
				}
				isContinue = true
			}
			return nil
		},

		//変化技
		func(src, target *bp.Pokemon) error {
			isStatus := moveData.Category == bp.STATUS
			if moveData.Category == bp.STATUS {
				if target.IsSubstituteState() {
					if moveData.PiercingSubstitute != bp.NO_PIERCING_SUBSTITUTE {
						status(battle, src, target, context)
					}
				} else {
					status(battle, src, target, context)
				}
			}
			isContinue = isStatus
			return nil
		},

		//物理・特殊
		func(src, target *bp.Pokemon) error {
			dmg, isNoEffect, err := battle.CalculateDamage(action, target, isSingleDmg || (faintedCount-1) == targetN, context)
			if err != nil {
				return err
			}

			if isNoEffect {
				isContinue = isNoEffect
				return nil
			}

			if target.IsSubstituteState() {
				if moveData.PiercingSubstitute == bp.NO_PIERCING_SUBSTITUTE {
					err = target.SubSubstituteHP(dmg)
				} else {
					err = target.SubCurrentHP(dmg)
					isBodyAttack = true
				}
			} else {
				err = target.SubCurrentHP(dmg)
				isBodyAttack = true
			}
			return err
		},

		//攻撃された側の瀕死のカウント
		func(_, target *bp.Pokemon) error {
			isTargetFainted := target.IsFainted()
			if isTargetFainted {
				faintedCount += 1
			}
			return nil
		},

		//自分への追加効果
		func(src, target *bp.Pokemon) error {
			return self(src)
		},

		//相手への追加効果
		func(_, target *bp.Pokemon) error {
			if isBodyAttack && !isTargetFainted {
				return opponent(target)
			} else {
				return nil
			}
		},
	}

	for _, target := range targetPokemons {
		isBodyAttack = false
		isTargetFainted = false
		isContinue = false

		for _, f := range fs {
			err := f(src, target)
			if err != nil {
				return err
			}
			if isContinue {
				continue
			}
		}
	}
	return nil
}

type Move func(*Battle, *SoloAction, *Context) error

func EmptyMove(_ *Battle, _ *SoloAction, _ *Context) error {
	return nil
}

//10まんボルト
func Thunderbolt(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.THUNDERBOLT {
		return fmt.Errorf("Thunderboltに渡されたAction.MoveNameがTHUNDERBOLTではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		func(target *bp.Pokemon) error {
			return target.SetStatusAilment(bp.PARALYSIS, 10, context.Rand)
		},
	)
}

//アームハンマー
func HammerArm(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.HAMMER_ARM {
		return fmt.Errorf("HammerArmに渡されたAction.MoveNameがbp.HammerArmではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		func(src *bp.Pokemon) error {
			isClearBodyValid := false
			return src.RankFluctuation(&bp.RankStat{Speed:-1}, 100, isClearBodyValid, context.Rand)
		},
		EmptyAdditionalEffect,
	)
}

//ストーンエッジ
func StoneEdge(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.STONE_EDGE {
		return fmt.Errorf("StoneEdgeに渡されたAction.MoveNameがbp.STONE_EDGEではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
}

//なみのり
func Surf(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.SURF {
		return fmt.Errorf("Surfに渡されたAction.MoveNameがbp.SURFではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
}

//れいとうビーム
func IceBeam(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.ICE_BEAM {
		return fmt.Errorf("IceBeamに渡されたAction.MoveNameがbp.ICE_BEAMではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		func(target *bp.Pokemon) error {
			return target.SetStatusAilment(bp.FREEZE, 10, context.Rand)
		},
	)
}

//わるあがき
func Struggle(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.STRUGGLE {
		return fmt.Errorf("Struggleに渡されたAction.MoveNameがbp.STRUGGLEではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		func(src *bp.Pokemon) error {
			dmg := int(float64(src.CurrentHP) / 4.0)
			return src.SubCurrentHP(dmg)
		},
		EmptyAdditionalEffect,
	)
}

//あまごい
func RainDance(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.RAIN_DANCE {
		return fmt.Errorf("RainDanceに渡されたAction.MoveNameがbp.RAIN_DANCEではない。")
	}
	battle.Weather = RAIN
	battle.RemainingTurnWeather = 5
	return nil
}

//いわなだれ
func RockSlide(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.ROCK_SLIDE {
		return fmt.Errorf("RockSlideに渡されたAction.MoveNameがbp.ROCK_SLIDEではない。")
	}

	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		func(target *bp.Pokemon) error {
			return target.SetIsFlinchState(30, context.Rand)
		},
	)
}

//おんがえし
func Return(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.RETURN {
		return fmt.Errorf("Returnに渡されたAction.MoveNameがbp.RETURNではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
}

/*
	第4世代のデータ
	かみくだく	あく	ぶつり	80	100	15	○	○	×
	通常	20%の確率で相手の『ぼうぎょ』ランクを1段階下げる。
*/
func Crunch(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.CRUNCH {
		return fmt.Errorf("Crunchに渡されたAction.MoveNameがbp.CRUNCHではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		func(target *bp.Pokemon) error {
			isClearBodyValid := true
			return target.RankFluctuation(&bp.RankStat{Def:-1}, 20, isClearBodyValid, context.Rand)
		},
	)
}

//がむしゃら
func Endeavor(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.ENDEAVOR {
		return fmt.Errorf("Endeavorに渡されたAction.MoveNameがbp.ENDEAVORではない。")
	}
	err := moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
	return err
}

//こごえるかぜ
func IcyWind(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.ICY_WIND {
		return fmt.Errorf("IcyWindに渡されたAction.MoveNameがbp.ICY_WINDではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		func(targetPokemon *bp.Pokemon) error {
			isClearBodyValid := true
			return targetPokemon.RankFluctuation(&bp.RankStat{Speed:-1}, 100, isClearBodyValid, context.Rand)
		},
	)
}

//このゆびとまれ
func FollowMe(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.FOLLOW_ME {
		return fmt.Errorf("FollowMeに渡されたAction.MoveNameがbp.FOLLOW_MEではない。")
	}

	maxPriority := omwmath.Max(
		omwmath.Max(battle.SelfLeadPokemons.FollowMePriorities()...),
		omwmath.Max(battle.OpponentLeadPokemons.FollowMePriorities()...),
	)

	if maxPriority == 0 {
		battle.SelfLeadPokemons[action.SrcIndex].FollowMePriority = bp.MAX_FOLLOW_ME_PRIORITY
	} else {
		battle.SelfLeadPokemons[action.SrcIndex].FollowMePriority = maxPriority - 1
	}
	return nil
}

//さいみんじゅつ
func Hypnosis(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.HYPNOSIS {
		return fmt.Errorf("Hypnosisに渡されたAction.MoveNameがbp.HYPNOSISではない。")
	}
	return moveHelper(
		battle, action, context,
		func(_ *Battle, _, target *bp.Pokemon, context *Context) error {
			//命中が確定している
			target.SetStatusAilment(bp.SLEEP, 100, context.Rand)
			return nil
		},
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
}

//じこあんじ
func Recover(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.RECOVER {
		return fmt.Errorf("Recoverに渡されたAction.MoveNameがbp.RECOVERではない。")
	}
	return moveHelper(
		battle, action, context,
		func(_ *Battle, src, target *bp.Pokemon, _ *Context) error {
			src.Rank = target.Rank.Clone()
			return nil
		},
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
}

//じしん
func Earthquake(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.EARTHQUAKE {
		return fmt.Errorf("Earthquakeに渡されたAction.MoveNameがbp.EARTHQUAKEではない。")
	}
	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
}

//じばく
func SelfDestruct(battle *Battle, action *SoloAction, context *Context) error {
	if action.MoveName != bp.SELF_DESTRUCT {
		return fmt.Errorf("SelfDestructに渡されたAction.MoveNameがbp.SELF_DESTRUCTではない。")
	}

	return moveHelper(
		battle, action, context,
		EmptyStatus,
		EmptyAdditionalEffect,
		EmptyAdditionalEffect,
	)
}

func GetMoveFunc(moveName bp.MoveName) Move {
	switch moveName {
		case bp.THUNDERBOLT:
			return Thunderbolt
		case bp.HAMMER_ARM:
			return HammerArm
		case bp.STONE_EDGE:
			return StoneEdge
		case bp.SURF:
			return Surf
		case bp.ICE_BEAM:
			return IceBeam
		case bp.STRUGGLE:
			return Struggle
		case bp.RAIN_DANCE:
			return RainDance
		case bp.ROCK_SLIDE:
			return RockSlide
		case bp.RETURN:
			return Return
		case bp.CRUNCH:
			return Crunch
		case bp.ENDEAVOR:
			return Endeavor
		case bp.ICY_WIND:
			return IcyWind
		case bp.FOLLOW_ME:
			return FollowMe
		case bp.HYPNOSIS:
			return Hypnosis
		case bp.RECOVER:
			return Recover
		case bp.EARTHQUAKE:
			return Earthquake
		case bp.SELF_DESTRUCT:
			return SelfDestruct
	}
	return EmptyMove
}