package battle

import (
	"fmt"
	"math"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	//"github.com/sw965/bippa/battle/dmgtools"
)

/*
	第4世代の技データ
	(あ行～た行) https://yakkun.com/dp/waza_list.htm#k
	(な行～わ行) https://yakkun.com/dp/waza_list2.htm
*/

type StatusEffect func(*Manager, *bp.Pokemon, *bp.Pokemon, *Context) error

// https://wiki.xn--rckteqa2e.com/wiki/%E8%BF%BD%E5%8A%A0%E5%8A%B9%E6%9E%9C
type AdditionalEffect func(*bp.Pokemon, *Context) error

type Move struct {
	StatusEffect StatusEffect
	SelfAdditionalEffect AdditionalEffect
	OpponentAdditionalEffect AdditionalEffect
}

func (m *Move) Run(battle *Manager, action *SoloAction, context *Context) error {
	src := &battle.SelfLeadPokemons[action.SrcIndex]
	if src.IsFainted() {
		m := fmt.Sprintf("%s は 瀕死状態なので、技を繰り出す事が出来ません。", src.Name.ToString())
		return fmt.Errorf(m)
	}

	if src.IsFlinchState {
		m := fmt.Sprintf("%s は 怯み状態なので、技を繰り出す事が出来ません。", src.Name.ToString())
		return fmt.Errorf(m)
	}

	switch action.MoveName {
		//あまごい
		case bp.RAIN_DANCE:
			battle.Weather = RAIN
			battle.RemainingTurnWeather = 5
			return nil
		//ねこだまし
		case bp.FAKE_OUT:
			if src.Turn != 1 {
				return nil
			}
		//トリックルーム
		case bp.TRICK_ROOM:
			if battle.IsTrickRoomState() {
				battle.RemainingTurnTrickRoom = 0
			} else {
				battle.RemainingTurnTrickRoom = 5
			}
			context.Observer(battle, TRICK_ROOM_EVENT)
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

	faintedCount := 0
	if action.MoveName == bp.SELF_DESTRUCT || action.MoveName == bp.EXPLOSION {
		src.Stat.CurrentHP = 0
	}

	var err error
	for _, target := range targetPokemons {
		var isHit bool
		if moveData.Accuracy == -1 {
			isHit = true
		} else {
			isHit, err = omwrand.IsPercentageMet(moveData.Accuracy, context.Rand)
			if err != nil {
				return err
			}
		}

		if !isHit {
			continue
		}

		//かんそうはだ
		if moveData.Type == bp.WATER && target.Ability == bp.DRY_SKIN {
			fmt.Println("かんそうはだが発動した！")
			heal := int(float64(target.Stat.CurrentHP) * 0.25)
			err := target.AddCurrentHP(heal)
			if err != nil {
				return err
			}
			continue
		}

		if moveData.Category == bp.STATUS {
			if target.IsSubstituteState() {
				if !moveData.CanSubstitute {
					m.StatusEffect(battle, src, target, context)
				}
			} else {
				m.StatusEffect(battle, src, target, context)
			}
		} else {
			dmg, isNoEffect, err := battle.CalculateDamage(action, target, isSingleDmg || (faintedCount-1) == targetN, context)
			if err != nil {
				return err
			}

			if isNoEffect {
				continue
			}

			var isBodyAttack bool
			var isFocusSash bool

			if target.IsSubstituteState() {
				if moveData.CanSubstitute {
					err = target.SubSubstituteHP(dmg)
				} else {
					isFocusSash, err = target.SubCurrentHP(dmg, true)
					isBodyAttack = true
				}
			} else {
				isFocusSash, err = target.SubCurrentHP(dmg, true)
				isBodyAttack = true
			}

			if err != nil {
				return err
			}

			if isFocusSash {
				context.Observer(battle, ITEM_USE_EVENT)
			}

			isTargetFainted := target.IsFainted()
			if isTargetFainted {
				faintedCount += 1	
			}

			if m.SelfAdditionalEffect != nil {
				m.SelfAdditionalEffect(src, context)
			}

			if isBodyAttack && !isTargetFainted && m.OpponentAdditionalEffect != nil {
				m.OpponentAdditionalEffect(target, context)
			}
		}
	}
	return nil
}

//10まんボルト
func NewThunderbolt() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			return target.SetStatusAilment(bp.PARALYSIS, 10, context.Rand)
		},
	}
}

//アームハンマー
func NewHammerArm() Move {
	return Move{
		SelfAdditionalEffect:func(src *bp.Pokemon, context *Context) error {
			isClearBodyValid := false
			return src.RankFluctuation(&bp.RankStat{Speed:-1}, 100, isClearBodyValid, context.Rand)
		},
	}
}

//ストーンエッジ
func NewStoneEdge() Move {
	return Move{}
}

//なみのり
func NewSurf() Move {
	return Move{}
}

//れいとうビーム
func NewIceBeam() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			return target.SetStatusAilment(bp.FREEZE, 10, context.Rand)
		},
	}
}

//わるあがき
func NewStruggle() Move {
	return Move{
		SelfAdditionalEffect:func(src *bp.Pokemon, context *Context) error {
			dmg := int(float64(src.Stat.CurrentHP) / 4.0)
			_, err := src.SubCurrentHP(dmg, false)
			return err
		},
	}
}

//あまごい
func NewRainDance() Move {
	return Move{}
}

//いわなだれ
func NewRockSlide() Move {
	return Move{
		SelfAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			return target.SetIsFlinchState(30, context.Rand)
		},
	}
}

//おんがえし
func NewReturn() Move {
	return Move{}
}

//かみくだく
func NewCrunch() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			isClearBodyValid := true
			return target.RankFluctuation(&bp.RankStat{Def:-1}, 20, isClearBodyValid, context.Rand)
		},
	}
}

//がむしゃら
func NewEndeavor() Move {
	return Move{}
}

//こごえるかぜ
func NewIcyWind() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			isClearBodyValid := true
			return target.RankFluctuation(&bp.RankStat{Speed:-1}, 100, isClearBodyValid, context.Rand)
		},
	}
}

//このゆびとまれ
func NewFollowMe() Move {
	return Move{
		StatusEffect:func(b *Manager, src, target *bp.Pokemon, context *Context) error {
			if src != target {
				return fmt.Errorf("このゆびとまれ は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}
			b.SelfFollowMePokemonPointers = append(b.SelfFollowMePokemonPointers, src)
			context.Observer(b, FOLLOW_ME_EVENT)
			return nil
		},
	}
}

//さいみんじゅつ
func NewHypnosis() Move {
	return Move{
		StatusEffect:func(_ *Manager, _ *bp.Pokemon, target *bp.Pokemon, context *Context) error {
			return target.SetStatusAilment(bp.SLEEP, 100, context.Rand)
		},
	}
}

//じこあんじ
func NewRecover() Move {
	return Move{
		StatusEffect:func(_ *Manager, src *bp.Pokemon, target *bp.Pokemon, _ *Context) error {
			src.Rank = target.Rank.Clone()
			return nil
		},
	}
}

//じしん
func NewEarthquake() Move {
	return Move{}
}

//じばく
func NewSelfDestruct() Move {
	return Move{}
}

//たきのぼり
func NewWaterfall() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			return target.SetIsFlinchState(20, context.Rand)
		},
	}
}

//だいばくはつ
func NewExplosion() Move {
	return Move{}
}

//ちょうはつ
func NewTaunt() Move {
	return Move{
		StatusEffect:func(_ *Manager, _, target *bp.Pokemon, context *Context) error {
			target.SetRemainingTurnTauntState(context.Rand)
			return nil
		},
	}
}

//でんじは
func NewThunderWave() Move {
	return Move{
		StatusEffect:func(_ *Manager, _, target *bp.Pokemon, context *Context) error {
			return target.SetStatusAilment(bp.PARALYSIS, 100, context.Rand)
		},
	}
}

//ねこだまし
func NewFakeOut() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			return target.SetIsFlinchState(100, context.Rand)
		},
	}
}

//ねっぷう
func NewHeatWave() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			return target.SetStatusAilment(bp.BURN, 10, context.Rand)
		},
	}
}

//はらだいこ
func NewBellyDrum() Move {
	return Move{
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon, _ *Context) error {
			if src != target {
				return fmt.Errorf("はらだいこ は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}
			src.Rank.Atk = bp.MAX_RANK
			return nil
		},
	}
}

//ふいうち
func NewSuckerPunch() Move {
	return Move{}
}

//ほのおのパンチ
func NewFirePunch() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			return target.SetStatusAilment(bp.BURN, 10, context.Rand)
		},
	}
}

//まもる
func NewProtect() Move {
	return Move{
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon, context *Context) error {
			if src != target {
				return fmt.Errorf("まもる は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}

			// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%BE%E3%82%82%E3%82%8B#%E6%88%90%E5%8A%9F%E7%8E%87
			isSuccess := math.Pow(0.5, float64(src.ProtectConsecutiveSuccess)) > context.Rand.Float64()
			src.IsProtectState = isSuccess
			if isSuccess {
				src.ProtectConsecutiveSuccess += 1
			} else {
				src.ProtectConsecutiveSuccess = 0
			}
			return nil
		},
	}
}

//みがわり
func NewSubstitute() Move {
	return Move{
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon, _ *Context) error {
			if src != target {
				return fmt.Errorf("みがわり は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}

			// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%BF%E3%81%8C%E3%82%8F%E3%82%8A

			if src.IsSubstituteState() {
				return nil
			}

			cost := int(float64(src.Stat.MaxHP) * 0.25)

			if src.Stat.CurrentHP > cost {
				src.Stat.CurrentHP -= cost
				src.SubstituteHP = cost
			}

			return nil
		},
	}
}

//りゅうせいぐん
func NewDracoMeteor() Move {
	return Move{
		SelfAdditionalEffect:func(src *bp.Pokemon, context *Context) error {
			isClearBodyValid := false
			src.RankFluctuation(&bp.RankStat{SpAtk:-2}, 100, isClearBodyValid, context.Rand)
			return nil
		},
	}
}

//クロスチョップ
func NewCrossChop() Move {
	return Move{}
}


//コメットパンチ
func NewCometPunch() Move {
	return Move{
		SelfAdditionalEffect:func(src *bp.Pokemon, context *Context) error {
			isClearBodyValid := false
			src.RankFluctuation(&bp.RankStat{Atk:1}, 20, isClearBodyValid, context.Rand)
			return nil
		},
	}
}

//サイコキネシス
func NewPsychic() Move {
	return Move{
		OpponentAdditionalEffect:func(target *bp.Pokemon, context *Context) error {
			isClearBodyValid := true
			target.RankFluctuation(&bp.RankStat{SpDef:-1}, 10, isClearBodyValid, context.Rand)
			return nil
		},
	}
}

//ジャイロボール
func NewGyroBall() Move {
	return Move{}
}

//ダークホール
func NewDarkVoid() Move {
	return Move{
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon, context *Context) error {
			return target.SetStatusAilment(bp.SLEEP, 100, context.Rand)
		},
	}
}

//トリックルーム
func NewTrickRoom() Move {
	return Move{}
}

//ハイドロポンプ
func NewHydroPump() Move {
	return Move{}
}


//バレットパンチ
func NewBulletPunch() Move {
	return Move{}
}

func GetMove(moveName bp.MoveName) Move {
	switch moveName {
		case bp.THUNDERBOLT:
			return NewThunderbolt()
		case bp.HAMMER_ARM:
			return NewHammerArm()
		case bp.STONE_EDGE:
			return NewStoneEdge()
		case bp.SURF:
			return NewSurf()
		case bp.ICE_BEAM:
			return NewIceBeam()
		case bp.STRUGGLE:
			return NewStruggle()
		case bp.RAIN_DANCE:
			return NewRainDance()
		case bp.ROCK_SLIDE:
			return NewRockSlide()
		case bp.RETURN:
			return NewReturn()
		case bp.CRUNCH:
			return NewCrunch()
		case bp.ENDEAVOR:
			return NewEndeavor()
		case bp.ICY_WIND:
			return NewIcyWind()
		case bp.FOLLOW_ME:
			return NewFollowMe()
		case bp.HYPNOSIS:
			return NewHypnosis()
		case bp.RECOVER:
			return NewRecover()
		case bp.EARTHQUAKE:
			return NewEarthquake()
		case bp.SELF_DESTRUCT:
			return NewSelfDestruct()
		case bp.WATERFALL:
			return NewWaterfall()
		case bp.EXPLOSION:
			return NewExplosion()
		case bp.TAUNT:
			return NewTaunt()
		case bp.THUNDER_WAVE:
			return NewThunderWave()
		case bp.FAKE_OUT:
			return NewFakeOut()
		case bp.HEAT_WAVE:
			return NewHeatWave()
		case bp.BELLY_DRUM:
			return NewBellyDrum()
		case bp.SUCKER_PUNCH:
			return NewSuckerPunch()
		case bp.FIRE_PUNCH:
			return NewFirePunch()
		case bp.PROTECT:
			return NewProtect()
		case bp.SUBSTITUTE:
			return NewSubstitute()
		case bp.DRACO_METEOR:
			return NewDracoMeteor()
		case bp.CROSS_CHOP:
			return NewCrossChop()
		case bp.COMET_PUNCH:
			return NewCometPunch()
		case bp.PSYCHIC:
			return NewHypnosis()
		case bp.GYRO_BALL:
			return NewGyroBall()
		case bp.DARK_VOID:
			return NewDarkVoid()
		case bp.TRICK_ROOM:
			return NewTrickRoom()
		case bp.HYDRO_PUMP:
			return NewHydroPump()
		case bp.BULLET_PUNCH:
			return NewBulletPunch()
	}
	return Move{}
}