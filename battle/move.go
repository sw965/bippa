package battle

import (
	"fmt"
	"math"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	"golang.org/x/exp/slices"
)

// /*
// 	第4世代の技データ
// 	(あ行～た行) https://yakkun.com/dp/waza_list.htm#k
// 	(な行～わ行) https://yakkun.com/dp/waza_list2.htm
// */

type StatusEffect func(*Manager, *bp.Pokemon, *bp.Pokemon) error

// https://wiki.xn--rckteqa2e.com/wiki/%E8%BF%BD%E5%8A%A0%E5%8A%B9%E6%9E%9C
type AdditionalEffect func(*Manager, *bp.Pokemon) error

type Move struct {
	StatusEffect StatusEffect
	SelfAdditionalEffect AdditionalEffect
	OpponentAdditionalEffect AdditionalEffect
}

func (m *Move) Run(manager *Manager, action *SoloAction) error {
	src := &manager.SelfLeadPokemons[action.SrcIndex]
	if src.IsFainted() {
		m := fmt.Sprintf("%s は 瀕死状態なので、技を繰り出す事が出来ません。", src.Name.ToString())
		return fmt.Errorf(m)
	}

	mm := MessageMaker{IsSelf:manager.IsHostView}
	manager.HostViewMessage = mm.MoveUse(src.Name, action.MoveName)
	GlobalContext.Observer(manager, MESSAGE_EVENT)

	switch action.MoveName {
		//あまごい
		case bp.RAIN_DANCE:
			manager.Weather = RAIN
			manager.RemainingTurn.Weather = 5
			return nil
		//ねこだまし
		case bp.FAKE_OUT:
			if src.TurnCount != 1 {
				return nil
			}
		//トリックルーム
		case bp.TRICK_ROOM:
			if manager.IsTrickRoomState() {
				manager.RemainingTurn.TrickRoom = 0
			} else {
				manager.RemainingTurn.TrickRoom = 5
			}
			return nil
		//じばく
		case bp.SELF_DESTRUCT:
			src.Stat.CurrentHP = 0
		//だいばくはつ
		case bp.EXPLOSION:
			src.Stat.CurrentHP = 0
	}

	moveData := bp.MOVEDEX[action.MoveName]
	targetPokemons := manager.TargetPokemonPointers(action)
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
			isHit, err = omwrand.IsPercentageMet(moveData.Accuracy, GlobalContext.Rand)
			if err != nil {
				return err
			}
		}

		if !isHit {
			continue
		}

		//かんそうはだ
		if moveData.Type == bp.WATER && target.Ability == bp.DRY_SKIN {
			heal := int(float64(target.Stat.CurrentHP) * 0.25)
			err := target.ApplyHealToBody(heal)
			if err != nil {
				return err
			}
			continue
		}

		if moveData.Category == bp.STATUS {
			if target.IsSubstituteState() {
				if !moveData.CanSubstitute {
					m.StatusEffect(manager, src, target)
				}
			} else {
				m.StatusEffect(manager, src, target)
			}
		} else {
			switch action.MoveName {
				//ふいうち
				case bp.SUCKER_PUNCH:
					if target.ThisTurnPlannedUseMoveName == bp.EMPTY_MOVE_NAME {
						break
					}
					md := bp.MOVEDEX[target.ThisTurnPlannedUseMoveName]
					if md.Category == bp.STATUS {
						break
					}
			}

			isCrit, err := IsCritical(moveData.CriticalRank, GlobalContext.Rand)
			if err != nil {
				return err
			}

			calc := DamageCalculator{
				Attacker:NewAttackerInfo(src),
				Defender:NewDefenderInfo(target),
				IsCritical:isCrit,
				RandBonus:GlobalContext.GetDamageRandBonus(),
				IsSingleDamage:isSingleDmg || (faintedCount-1) == targetN,
				IsDamageCappedByCurrentHP:true,
			}

			dmgResult := calc.Calculation(action.MoveName)
			if dmgResult.TypeEffective == bp.NO_EFFECTIVE {
				continue
			}
			dmg := dmgResult.Damage

			var isFocusSash bool
			var isBodyAttack bool

			bodyAttack := func() error {
				if target.Item == bp.FOCUS_SASH {
					if target.IsFullHP() && target.Stat.MaxHP == dmg {
						isFocusSash = true
						dmg -= 1
					}
				}
				target.ApplyDamageToBody(dmg)
				isBodyAttack = true
				return err
			}

			if target.IsSubstituteState() {
				if moveData.CanSubstitute {
					target.ApplyDamageToSubstitute(dmg)
				} else {
					err = bodyAttack()
				}
			} else {
				err = bodyAttack()
			}

			if err != nil {
				return err
			}

			if isFocusSash {
				target.Item = bp.EMPTY_ITEM
			}

			isTargetFainted := target.IsFainted()
			if isTargetFainted {
				faintedCount += 1
			}

			if m.SelfAdditionalEffect != nil && !src.IsFainted() {
				m.SelfAdditionalEffect(manager, src)
			}

			if isBodyAttack && !isTargetFainted && m.OpponentAdditionalEffect != nil {
				m.OpponentAdditionalEffect(manager, target)
			}
		}
	}
	return nil
}

//10まんボルト
func NewThunderbolt() Move {
	return Move{
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			if target.StatusAilment != bp.EMPTY_STATUS_AILMENT {
				return nil
			}

			ok, err := omwrand.IsPercentageMet(10, GlobalContext.Rand)
			if ok {
				target.StatusAilment = bp.PARALYSIS
			}
			return err
		},
	}
}

//アームハンマー
func NewHammerArm() Move {
	return Move{
		SelfAdditionalEffect:func(m *Manager, src *bp.Pokemon) error {
			if src.Rank.Speed != bp.MIN_RANK {
				src.Rank.Speed -= 1
			}
			return nil
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
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			if target.StatusAilment != bp.EMPTY_STATUS_AILMENT {
				return nil
			}

			if slices.Contains(target.Types, bp.ICE) {
				return nil
			}

			ok, err := omwrand.IsPercentageMet(10, GlobalContext.Rand)
			if ok {
				target.StatusAilment = bp.FREEZE
			}
			return err
		},
	}
}

//わるあがき
func NewStruggle() Move {
	return Move{
		SelfAdditionalEffect:func(m *Manager, src *bp.Pokemon) error {
			dmg := int(float64(src.Stat.CurrentHP) / 4.0)
			err := src.ApplyDamageToBody(dmg)
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
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			ok, err := omwrand.IsPercentageMet(30, GlobalContext.Rand)
			if ok {
				target.IsFlinchState = true
			}
			return err
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
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			if target.Ability == bp.CLEAR_BODY {
				return nil
			}

			if target.Rank.Def == bp.MIN_RANK {
				return nil
			}

			ok, err := omwrand.IsPercentageMet(20, GlobalContext.Rand)
			if ok {
				target.Rank.Def -= 1
			}
			return err
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
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			if target.Ability == bp.CLEAR_BODY {
				return nil
			}

			if target.Rank.Speed == bp.MIN_RANK {
				return nil
			}

			target.Rank.Speed -= 1
			return nil
		},
	}
}

//このゆびとまれ
func NewFollowMe() Move {
	return Move{
		StatusEffect:func(m *Manager, src, target *bp.Pokemon) error {
			if src != target {
				return fmt.Errorf("このゆびとまれ は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}
			m.SelfFollowMePokemonPointers = append(m.SelfFollowMePokemonPointers, src)
			return nil
		},
	}
}

//さいみんじゅつ
func NewHypnosis() Move {
	return Move{
		StatusEffect:func(m *Manager, src, target *bp.Pokemon) error {
			if target.StatusAilment != bp.EMPTY_STATUS_AILMENT {
				return nil
			}

			target.StatusAilment = bp.SLEEP
			if target.Item == bp.LUM_BERRY {
				target.Item = bp.EMPTY_ITEM
				target.StatusAilment = bp.EMPTY_STATUS_AILMENT
			}
			return nil
		},
	}
}

//じこあんじ
func NewRecover() Move {
	return Move{
		StatusEffect:func(_ *Manager, src *bp.Pokemon, target *bp.Pokemon) error {
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
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			ok, err := omwrand.IsPercentageMet(20, GlobalContext.Rand)
			if ok {
				target.IsFlinchState = true
			}
			return err
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
		StatusEffect:func(m *Manager, _, target *bp.Pokemon) error {
			target.RemainingTurnTauntState = omwrand.IntUniform(2, 5, GlobalContext.Rand)
			return nil
		},
	}
}

//でんじは
func NewThunderWave() Move {
	return Move{
		StatusEffect:func(m *Manager, src, target *bp.Pokemon) error {
			if slices.Contains(target.Types, bp.GROUND) {
				return nil
			}
			if target.StatusAilment != bp.EMPTY_STATUS_AILMENT {
				return nil
			}
			target.StatusAilment = bp.PARALYSIS
			return nil
		},
	}
}

//ねこだまし
func NewFakeOut() Move {
	return Move{
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			target.IsFlinchState = true
			return nil
		},
	}
}

//ねっぷう
func NewHeatWave() Move {
	return Move{
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			if target.StatusAilment != bp.EMPTY_STATUS_AILMENT {
				return nil
			}

			if slices.Contains(target.Types, bp.FIRE) {
				return nil
			}

			ok, err := omwrand.IsPercentageMet(10, GlobalContext.Rand)
			if ok {
				target.StatusAilment = bp.BURN
			}
			return err
		},
	}
}

//はらだいこ
func NewBellyDrum() Move {
	return Move{
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon) error {
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
		OpponentAdditionalEffect:func(m *Manager, target *bp.Pokemon) error {
			if target.StatusAilment != bp.EMPTY_STATUS_AILMENT {
				return nil
			}

			if slices.Contains(target.Types, bp.FIRE) {
				return nil
			}

			ok, err := omwrand.IsPercentageMet(10, GlobalContext.Rand)
			if ok {
				target.StatusAilment = bp.BURN
			}
			return err
		},
	}
}

//まもる
func NewProtect() Move {
	return Move{
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon) error {
			if src != target {
				return fmt.Errorf("まもる は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}

			// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%BE%E3%82%82%E3%82%8B#%E6%88%90%E5%8A%9F%E7%8E%87
			isSuccess := math.Pow(0.5, float64(src.ProtectConsecutiveSuccess)) > GlobalContext.Rand.Float64()
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
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon) error {
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
		SelfAdditionalEffect:func(_ *Manager, src *bp.Pokemon) error {
			if src.Rank.SpAtk >= bp.MIN_RANK - 2 {
				src.Rank.Speed -= 2
			} else if src.Rank.SpAtk != bp.MIN_RANK {
				src.Rank.Speed -= 1
			}
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
		SelfAdditionalEffect:func(_ *Manager, src *bp.Pokemon) error {
			if src.Rank.Atk == bp.MAX_RANK {
				return nil
			}
			ok, err := omwrand.IsPercentageMet(20, GlobalContext.Rand)
			if ok {
				src.Rank.Atk += 1
			}
			return err
		},
	}
}

//サイコキネシス
func NewPsychic() Move {
	return Move{
		OpponentAdditionalEffect:func(_ *Manager, target *bp.Pokemon) error {
			if target.Rank.SpDef == bp.MIN_RANK {
				return nil
			}

			if target.Ability == bp.CLEAR_BODY {
				return nil
			}

			ok, err := omwrand.IsPercentageMet(10, GlobalContext.Rand)
			if ok {
				target.Rank.SpDef -= 1
			}
			return err
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
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon) error {
			if target.StatusAilment != bp.EMPTY_STATUS_AILMENT {
				return nil
			}
			target.StatusAilment = bp.SLEEP
			return nil
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