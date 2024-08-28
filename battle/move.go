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

func NewSelfRankStatFluctuationAdditionalEffect(v *bp.RankStat, percentage int) AdditionalEffect {
	return func(m *Manager, p *bp.Pokemon) error {
		ok, err := omwrand.IsPercentageMet(percentage, GlobalContext.Rand)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		fluctuation := v.AdjustFluctuation(&p.RankStat)
		p.RankStat = p.RankStat.Add(&fluctuation)
		msgs := GetStandardRankFluctuationMessages(p.Name, &fluctuation)
		humanNameMsg := m.getHumanNameMessage(p.IsHost)
		for _, msg := range msgs {
			m.HostViewMessage = humanNameMsg + msg
			GlobalContext.Observer(m)
		}
		return nil
	}
}

func NewOpponentRankStatFluctuationAdditionalEffect(v *bp.RankStat, percentage int) AdditionalEffect {
	return func(m *Manager, p *bp.Pokemon) error {
		ok, err := omwrand.IsPercentageMet(percentage, GlobalContext.Rand)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		if p.Ability == bp.CLEAR_BODY {
			r := v.DownToZero()
			v = &r
		}
		fluctuation := v.AdjustFluctuation(&p.RankStat)
		p.RankStat = p.RankStat.Add(&fluctuation)

		msgs := GetStandardRankFluctuationMessages(p.Name, &fluctuation)
		humanNameMsg := m.getHumanNameMessage(p.IsHost)
		for _, msg := range msgs {
			m.HostViewMessage = humanNameMsg + msg
			GlobalContext.Observer(m)
		}
		return nil
	}
}

func NewPraalysisAdditionalEffect(percentage int) AdditionalEffect {
	return func(m *Manager, p *bp.Pokemon) error {
		if p.StatusAilment != bp.EMPTY_STATUS_AILMENT {
			return nil
		}

		ok, err := omwrand.IsPercentageMet(percentage, GlobalContext.Rand)
		if ok {
			p.StatusAilment = bp.PARALYSIS
			humanNameMsg := m.getHumanNameMessage(p.IsHost)
			m.HostViewMessage = humanNameMsg + fmt.Sprintf("%sは まひして わざが でにくくなった！", p.Name.ToString())
			GlobalContext.Observer(m)
		}
		return err
	}
}

func NewBurnAdditionalEffect(percentage int) AdditionalEffect {
	return func(m *Manager, p *bp.Pokemon) error {
		if p.StatusAilment != bp.EMPTY_STATUS_AILMENT {
			return nil
		}
	
		if slices.Contains(p.Types, bp.FIRE) {
			return nil
		}
	
		ok, err := omwrand.IsPercentageMet(percentage, GlobalContext.Rand)
		if ok {
			p.StatusAilment = bp.BURN
			humanNameMsg := m.getHumanNameMessage(p.IsHost)
			m.HostViewMessage = humanNameMsg + fmt.Sprintf("%sは やけどをおった！", p.Name.ToString())
			GlobalContext.Observer(m)
		}
		return err
	}
}

func NewFreezeAdditionalEffect(percentage int) AdditionalEffect {
	return func(m *Manager, p *bp.Pokemon) error {
		if p.StatusAilment != bp.EMPTY_STATUS_AILMENT {
			return nil
		}

		if slices.Contains(p.Types, bp.ICE) {
			return nil
		}

		ok, err := omwrand.IsPercentageMet(percentage, GlobalContext.Rand)
		if ok {
			p.StatusAilment = bp.FREEZE
			humanNameMsg := m.getHumanNameMessage(p.IsHost)
			m.HostViewMessage = humanNameMsg + fmt.Sprintf("%sは こおりついた！", p.Name.ToString())
			GlobalContext.Observer(m)
		}
		return err
	}
}

func NewFlinchAdditionalEffect(percentage int) AdditionalEffect {
	return func(m *Manager, p *bp.Pokemon) error {
		ok, err := omwrand.IsPercentageMet(percentage, GlobalContext.Rand)
		if ok {
			p.IsFlinchState = true
		}
		return err
	}
}

type Move struct {
	StatusEffect StatusEffect
	SelfAdditionalEffect AdditionalEffect
	OpponentAdditionalEffect AdditionalEffect
}

func (move *Move) Run(m *Manager, action *SoloAction) error {
	src := &m.CurrentSelfLeadPokemons[action.SrcIndex]
	if src.IsFainted() {
		s := fmt.Sprintf("%sは 瀕死状態なので、技を繰り出す事が出来ません。", src.Name.ToString())
		return fmt.Errorf(s)
	}

	srcPokeNameStr := src.Name.ToString()
	srcHumanNameMsg := m.getHumanNameMessage(src.IsHost)

	// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%BE%E3%81%B2
	if src.StatusAilment == bp.PARALYSIS {
		ok, err := omwrand.IsPercentageMet(25, GlobalContext.Rand)
		if err != nil {
			return err
		}

		if ok {
			m.HostViewMessage = srcHumanNameMsg + srcPokeNameStr + "は からだが しびれて うごけない！"
			GlobalContext.Observer(m)
			return nil
		}
	}

	// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%93%E3%81%8A%E3%82%8A_(%E7%8A%B6%E6%85%8B%E7%95%B0%E5%B8%B8)
	if src.StatusAilment == bp.FREEZE {
		ok, err := omwrand.IsPercentageMet(20, GlobalContext.Rand)
		if err != nil {
			return err
		}

		if !ok {
			m.HostViewMessage = srcHumanNameMsg + srcPokeNameStr + "こおって しまって うごかない！"
			GlobalContext.Observer(m)
			return nil
		}

		src.StatusAilment = bp.EMPTY_STATUS_AILMENT
		m.HostViewMessage = srcHumanNameMsg + srcPokeNameStr + "の こおりが 溶けた！" 
		GlobalContext.Observer(m)
		return nil
	}

	if src.IsFlinchState {
		m.HostViewMessage = srcHumanNameMsg + fmt.Sprintf("%sは ひるんで わざが だせない！", srcPokeNameStr)
		GlobalContext.Observer(m)
		return nil
	}

	moveNameStr := action.MoveName.ToString()
	m.HostViewMessage = srcHumanNameMsg + fmt.Sprintf("%sの %s!", srcPokeNameStr, moveNameStr)
	GlobalContext.Observer(m)
	src.Moveset[action.MoveName].Current -= 1

	switch action.MoveName {
		//あまごい
		case bp.RAIN_DANCE:
			m.Weather = RAIN
			m.RemainingTurn.Weather = 5
			m.HostViewMessage = "雨が 降り始めた！"
			GlobalContext.Observer(m)
			return nil
		//ねこだまし
		case bp.FAKE_OUT:
			if src.TurnCount != 1 {
				return nil
			}
		//トリックルーム
		case bp.TRICK_ROOM:
			if m.IsTrickRoomState() {
				m.RemainingTurn.TrickRoom = 0
				m.HostViewMessage = srcHumanNameMsg + fmt.Sprintf("%sは じくうを もどした！", srcPokeNameStr)
				GlobalContext.Observer(m)
			} else {
				m.RemainingTurn.TrickRoom = 5
				m.HostViewMessage = srcHumanNameMsg + fmt.Sprintf("%sは じくうを ゆがめた！", srcPokeNameStr)
				GlobalContext.Observer(m)
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
	targetPokemons, err := m.TargetPokemonPointers(action)
	if err != nil {
		return err
	}
	// https://wiki.xn--rckteqa2e.com/wiki/%E3%83%80%E3%83%96%E3%83%AB%E3%83%90%E3%83%88%E3%83%AB
	// 複数を対象とする技は、第四世代ではすばやさが高いポケモンから処理される。
	targetPokemons.SortBySpeed()
	targetNum := len(targetPokemons)

	var isSingleDmg bool
	if action.MoveName == bp.SELF_DESTRUCT || action.MoveName == bp.EXPLOSION {
		isSingleDmg = targetNum <= 2
	} else {
		isSingleDmg = targetNum == 1
	}

	faintedCount := 0
	if action.MoveName == bp.SELF_DESTRUCT || action.MoveName == bp.EXPLOSION {
		src.Stat.CurrentHP = 0
	}

	for _, target := range targetPokemons {
		targetPokeNameStr := target.Name.ToString()

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

		targetHumanNameMsg := m.getHumanNameMessage(target.IsHost)

		//かんそうはだ
		if moveData.Type == bp.WATER && target.Ability == bp.DRY_SKIN {
			heal := int(float64(target.Stat.CurrentHP) * 0.25)
			isFullHP := target.IsFullHP()
			err := target.ApplyHealToBody(heal)
			if err != nil {
				return err
			}

			drySkinStr := bp.DRY_SKIN.ToString()
			if isFullHP {
				m.HostViewMessage = targetHumanNameMsg + fmt.Sprintf("%sの %sで %sは こうかが なかった！", targetPokeNameStr, drySkinStr, moveNameStr)
			} else {
				m.HostViewMessage = targetHumanNameMsg + fmt.Sprintf("%sは %sで かいふくした！", targetPokeNameStr, drySkinStr)
			}
			GlobalContext.Observer(m)
			continue
		}

		if moveData.Category == bp.STATUS {
			if target.IsSubstituteState() {
				if !moveData.CanSubstitute {
					move.StatusEffect(m, src, target)
				}
			} else {
				move.StatusEffect(m, src, target)
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
				IsSingleDamage:isSingleDmg || (faintedCount-1) == targetNum,
				IsDamageCappedByCurrentHP:true,
			}

			noEffectiveMsg := targetHumanNameMsg + fmt.Sprintf("%sには こうかは ないようだ...", targetPokeNameStr)

			dmgDetailResult := calc.Calculation(action.MoveName)
			if dmgDetailResult.TypeEffective == bp.NO_EFFECTIVE {
				m.HostViewMessage = noEffectiveMsg
				GlobalContext.Observer(m)
				continue
			} else if dmgDetailResult.IsEndeavorFailure {
				m.HostViewMessage = noEffectiveMsg
				GlobalContext.Observer(m)
				continue
			}
			dmg := dmgDetailResult.Damage

			var isFocusSash bool
			var isBodyAttack bool

			bodyAttack := func() error {
				if target.Item == bp.FOCUS_SASH {
					if target.IsFullHP() && target.Stat.MaxHP == dmg {
						isFocusSash = true
						dmg -= 1
					}
				}
				m.ApplyDamageToBody(target, dmg)
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

			if isCrit {
				m.HostViewMessage = targetHumanNameMsg + fmt.Sprintf("%sの きゅうしょに あたった！", targetPokeNameStr)
				GlobalContext.Observer(m)
			}

			switch dmgDetailResult.TypeEffective {
				case bp.SUPER_EFFECTIVE:
					m.HostViewMessage = targetHumanNameMsg + fmt.Sprintf("%sに こうかは ばつぐんだ！", targetPokeNameStr)
					GlobalContext.Observer(m)
				case bp.NOT_VERY_EFFECTIVE:
					m.HostViewMessage = targetHumanNameMsg + fmt.Sprintf("%sに こうかは いまひとつのようだ...", targetPokeNameStr)
					GlobalContext.Observer(m)
			}

			if isFocusSash {
				target.Item = bp.EMPTY_ITEM
				focusSashStr := bp.FOCUS_SASH.ToString()
				m.HostViewMessage = targetHumanNameMsg + fmt.Sprintf("%sは %sで もちこたえた！", targetPokeNameStr, focusSashStr)
				GlobalContext.Observer(m)
			}

			isTargetFainted := target.IsFainted()
			if isTargetFainted {
				faintedCount += 1
			}

			if move.SelfAdditionalEffect != nil && !src.IsFainted() {
				move.SelfAdditionalEffect(m, src)
			}

			if isBodyAttack && !isTargetFainted && move.OpponentAdditionalEffect != nil {
				move.OpponentAdditionalEffect(m, target)
			}
		}
	}
	return nil
}

//10まんボルト
func NewThunderbolt() Move {
	return Move{
		OpponentAdditionalEffect:NewPraalysisAdditionalEffect(10),
	}
}

//アームハンマー
func NewHammerArm() Move {
	return Move{
		SelfAdditionalEffect:func(m *Manager, src *bp.Pokemon) error {
			if src.RankStat.Speed != bp.MIN_RANK {
				src.RankStat.Speed -= 1
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
		OpponentAdditionalEffect:NewFreezeAdditionalEffect(10),
	}
}

//わるあがき
func NewStruggle() Move {
	return Move{
		SelfAdditionalEffect:func(m *Manager, src *bp.Pokemon) error {
			dmg := int(float64(src.Stat.CurrentHP) / 4.0)
			err := m.ApplyDamageToBody(src, dmg)
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
		OpponentAdditionalEffect:NewFlinchAdditionalEffect(30),
	}
}

//おんがえし
func NewReturn() Move {
	return Move{}
}

//かみくだく
func NewCrunch() Move {
	return Move{
		OpponentAdditionalEffect:NewOpponentRankStatFluctuationAdditionalEffect(&bp.RankStat{Def:-1}, 20),
	}
}

//がむしゃら
func NewEndeavor() Move {
	return Move{}
}

//こごえるかぜ
func NewIcyWind() Move {
	return Move{
		OpponentAdditionalEffect:NewOpponentRankStatFluctuationAdditionalEffect(&bp.RankStat{Speed:-1}, 100),
	}
}

//このゆびとまれ
func NewFollowMe() Move {
	return Move{
		StatusEffect:func(m *Manager, src, target *bp.Pokemon) error {
			if src != target {
				return fmt.Errorf("このゆびとまれ は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}
			m.CurrentSelfFollowMePokemonPointers = append(m.CurrentSelfFollowMePokemonPointers, src)
			m.HostViewMessage = m.getHumanNameMessage(src.IsHost) + src.Name.ToString() + " は ちゅうもくのまとになった！"
			GlobalContext.Observer(m)
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
		StatusEffect:func(m *Manager, src *bp.Pokemon, target *bp.Pokemon) error {
			src.RankStat = target.RankStat.Clone()
			m.HostViewMessage = m.getHumanNameMessage(src.IsHost) + fmt.Sprintf("%sは %sの のうりょうへんかを コピーした！", src.Name.ToString(), target.Name.ToString())
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
		OpponentAdditionalEffect:NewFlinchAdditionalEffect(20),
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
		OpponentAdditionalEffect:NewBurnAdditionalEffect(10),
	}
}

//はらだいこ
func NewBellyDrum() Move {
	return Move{
		StatusEffect:func(_ *Manager, src, target *bp.Pokemon) error {
			if src != target {
				return fmt.Errorf("はらだいこ は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}
			src.RankStat.Atk = bp.MAX_RANK
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
		OpponentAdditionalEffect:NewBurnAdditionalEffect(10),
	}
}

//まもる
func NewProtect() Move {
	return Move{
		StatusEffect:func(m *Manager, src, target *bp.Pokemon) error {
			if src != target {
				return fmt.Errorf("まもる は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}

			// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%BE%E3%82%82%E3%82%8B#%E6%88%90%E5%8A%9F%E7%8E%87
			isSuccess := math.Pow(0.5, float64(src.ProtectConsecutiveSuccess)) > GlobalContext.Rand.Float64()
			src.IsProtectState = isSuccess
			if isSuccess {
				src.ProtectConsecutiveSuccess += 1
				m.HostViewMessage = m.getHumanNameMessage(src.IsHost) + fmt.Sprintf("%sは まもりの たいせいに はいった！", src.Name.ToString())
			} else {
				src.ProtectConsecutiveSuccess = 0
				m.HostViewMessage = fmt.Sprintf("しかし うまく きまらなかった")
			}
			return nil
		},
	}
}

//みがわり
func NewSubstitute() Move {
	return Move{
		StatusEffect:func(m *Manager, src, target *bp.Pokemon) error {
			if src != target {
				return fmt.Errorf("みがわり は 技を繰り出したポケモン と 対象になるポケモン の アドレスが 一致していなければならない。")
			}

			// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%BF%E3%81%8C%E3%82%8F%E3%82%8A

			if src.IsSubstituteState() {
				m.HostViewMessage = "しかし" + m.getHumanNameMessage(src.IsHost) + src.Name.ToString() +"の みがわりは すでに でていた"
				GlobalContext.Observer(m)
				return nil
			}

			cost := int(float64(src.Stat.MaxHP) * 0.25)
			if src.Stat.CurrentHP > cost {
				src.Stat.CurrentHP -= cost
				src.SubstituteHP = cost
				m.HostViewMessage = m.getHumanNameMessage(src.IsHost) + fmt.Sprintf("%sの ぶんしんが あらわれた！", src.Name.ToString())
			} else {
				m.HostViewMessage = fmt.Sprintf("しかし みがわりを だすには たいりょくが たりなかった！")
			}
			GlobalContext.Observer(m)
			return nil
		},
	}
}

//りゅうせいぐん
func NewDracoMeteor() Move {
	return Move{
		SelfAdditionalEffect:func(_ *Manager, src *bp.Pokemon) error {
			if src.RankStat.SpAtk >= bp.MIN_RANK - 2 {
				src.RankStat.Speed -= 2
			} else if src.RankStat.SpAtk != bp.MIN_RANK {
				src.RankStat.Speed -= 1
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
			if src.RankStat.Atk == bp.MAX_RANK {
				return nil
			}
			ok, err := omwrand.IsPercentageMet(20, GlobalContext.Rand)
			if ok {
				src.RankStat.Atk += 1
			}
			return err
		},
	}
}

//サイコキネシス
func NewPsychic() Move {
	return Move{
		OpponentAdditionalEffect:NewOpponentRankStatFluctuationAdditionalEffect(&bp.RankStat{SpDef:-1}, 10),
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