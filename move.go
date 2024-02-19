package bippa

import (
	"golang.org/x/exp/slices"
	"math/rand"
	"fmt"
	"math"
)

type MoveName int

const (
	NO_MOVE_NAME MoveName = iota
	A_KU_RO_BA_TTO
	A_SA_NO_HI_ZA_SI
	A_MA_E_RU
	U_SO_NA_KI
	E_NA_ZI_BO_RU
	KA_MI_TU_KU
	KA_RA_GE_N_KI
	GI_GA_DO_RE_I_N
	KI_RI_SA_KU
	KU_SA_NO_TI_KA_I
	KU_SA_MU_SU_BI
	KU_SA_WA_KE
	GU_RA_SU_FI_RU_DO
	KO_U_GO_U_SE_I
	KO_U_SO_KU_I_DO_U
	KO_NO_HA
	KO_RA_E_RU
	SA_I_DO_CHE_N_ZI
	ZI_KO_SA_I_SE_I
	SI_PPO_WO_HU_RU
	SYA_DO_KU_RO
	ZYA_RE_TU_KU
	SU_NA_A_TU_ME
	SU_PI_DO_SU_TA
	SO_RA_BI_MU
	TA_NE_BA_KU_DA_N
	TA_NE_MA_SI_N_GA_N
	TA_MA_GO_U_MI
	CHA_MU_BO_I_SU
	CHO_U_HA_TU
	TU_KI_NO_HI_KA_RI
	TU_ME_TO_GI
	TE_DA_SU_KE
	TE_RA_BA_SU_TO
	DE_N_KO_U_SE_KKA
	TO_SSI_N
	DO_RO_KA_KE
	TO_N_BO_GA_E_RI
	NA_MA_KE_RU
	NA_YA_MI_NO_TA_NE
	NE_GO_TO
	NE_MU_RU
	HA_NA_HU_BU_KI
	HA_NE_YA_SU_ME
	HI_KKA_KU
	HU_I_U_TI
	MA_ZI_KA_RU_RI_HU
	MA_NE_KKO
	MA_MO_RU
	MI_GA_WA_RI
	MI_RU_KU_NO_MI
	YA_DO_RI_GI_NO_TA_NE
	RI_HU_SU_TO_MU
	WA_RU_DA_KU_MI
	WA_RU_A_GA_KI
)

var STRING_TO_MOVE_NAME = map[string]MoveName{
	"アクロバット":A_KU_RO_BA_TTO,
	"あさのひざし":A_SA_NO_HI_ZA_SI,
	"あまえる":A_MA_E_RU,
	"うそなき":U_SO_NA_KI,
	"エナジーボール":E_NA_ZI_BO_RU,
	"かみつく":KA_MI_TU_KU,
	"からげんき":KA_RA_GE_N_KI,
	"ギガドレイン":GI_GA_DO_RE_I_N,
	"きりさく":KI_RI_SA_KU,
	"くさのちかい":KU_SA_NO_TI_KA_I,
	"くさむすび":KU_SA_MU_SU_BI,
	"くさわけ":KU_SA_WA_KE,
	"グラスフィールド":GU_RA_SU_FI_RU_DO,
	"こうごうせい":KO_U_GO_U_SE_I,
	"こうそくいどう":KO_U_SO_KU_I_DO_U,
	"このは":KO_NO_HA,
	"こらえる":KO_RA_E_RU,
	"サイドチェンジ":SA_I_DO_CHE_N_ZI,
	"じこさいせい":ZI_KO_SA_I_SE_I,
	"しっぽをふる":SI_PPO_WO_HU_RU,
	"シャドークロー":SYA_DO_KU_RO,
	"じゃれつく":ZYA_RE_TU_KU,
	"スピードスター":SU_PI_DO_SU_TA,
	"ソーラービーム":SO_RA_BI_MU,
	"タネばくだん":TA_NE_BA_KU_DA_N,
	"タネマシンガン":TA_NE_MA_SI_N_GA_N,
	"タマゴうみ":TA_MA_GO_U_MI,
	"チャームボイス":CHA_MU_BO_I_SU,
	"ちょうはつ":CHO_U_HA_TU,
	"つきのひかり":TU_KI_NO_HI_KA_RI,
	"つめとぎ":TU_ME_TO_GI,
	"てだすけ":TE_DA_SU_KE,
	"テラバースト":TE_RA_BA_SU_TO,
	"でんこうせっか":DE_N_KO_U_SE_KKA,
	"とっしん":TO_SSI_N,
	"どろかけ":DO_RO_KA_KE,
	"とんぼがえり":TO_N_BO_GA_E_RI,
	"なまける":NA_MA_KE_RU,
	"なやみのタネ":NA_YA_MI_NO_TA_NE,
	"ねごと":NE_GO_TO,
	"ねむる":NE_MU_RU,
	"はなふぶき":HA_NA_HU_BU_KI,
	"はねやすめ":HA_NE_YA_SU_ME,
	"ひっかく":HI_KKA_KU,
	"ふいうち":HU_I_U_TI,
	"マジカルリーフ":MA_ZI_KA_RU_RI_HU,
	"まねっこ":MA_NE_KKO,
	"まもる":MA_MO_RU,
	"みがわり":MI_GA_WA_RI,
	"ミルクのみ":MI_RU_KU_NO_MI,
	"やどりぎのタネ":YA_DO_RI_GI_NO_TA_NE,
	"リーフストーム":RI_HU_SU_TO_MU,
	"わるだくみ":WA_RU_DA_KU_MI,
}

func StringToMoveName(s string) MoveName {
	return STRING_TO_MOVE_NAME[s]
}

type MoveNames []MoveName

type Power int

func NewPower(bt *Battle, moveName MoveName) Power {
	moveData := MOVEDEX[moveName]
	if moveName == A_KU_RO_BA_TTO {
		return moveData.Power * 2
	} else {
		return moveData.Power
	}
}

type PowerPointUp int

const (
	MIN_POWER_POINT_UP = PowerPointUp(0)
	MAX_POWER_POINT_UP = PowerPointUp(3)
)

type PowerPointUps []PowerPointUp

func NewMaxPowerPointUps(n int) PowerPointUps {
	result := make(PowerPointUps, n)
	for i := 0; i < n; i++ {
		result[i] = MAX_POWER_POINT_UP
	}
	return result
}

var ALL_POWER_POINT_UPS = func() PowerPointUps {
	n := int(MAX_POWER_POINT_UP - MIN_POWER_POINT_UP) + 1
	result := make(PowerPointUps, n)
	for i := 0; i < n; i++ {
		result[i] = PowerPointUp(i)
	}
	return result
}

type PowerPoint struct {
	Max     int
	Current int
}

func NewPowerPoint(base int, up PowerPointUp) PowerPoint {
	bonus := (5.0 + float64(up)) / 5.0
	max := int(float64(base) * bonus)
	return PowerPoint{Max: max, Current: max}
}

func NewMaxPowerPoint(moveName MoveName) PowerPoint {
	base := MOVEDEX[moveName].BasePP
	return NewPowerPoint(base, MAX_POWER_POINT_UP)
}

type PowerPoints []PowerPoint

type Moveset map[MoveName]*PowerPoint

const (
	MIN_MOVESET_LENGTH = 1
	MAX_MOVESET_LENGTH = 4
)

func NewMoveset(pokeName PokeName, moveNames MoveNames, ppups PowerPointUps) (Moveset, error) {
	if len(moveNames) != len(ppups) {
		return Moveset{}, fmt.Errorf("len(moveNames) != len(ppups)")
	}

	y := Moveset{}
	learnset := POKEDEX[pokeName].Learnset
	for i, moveName := range moveNames {
		if !slices.Contains(learnset, moveName) {
			msg := fmt.Sprintf("%v は %v を 覚えない", pokeName, moveName)
			return Moveset{}, fmt.Errorf(msg)
		}

		base := MOVEDEX[moveName].BasePP
		up := ppups[i]
		pp := NewPowerPoint(base, up)
		y[moveName] = &pp
	}
	return y, nil
}

func (ms Moveset) Copy() Moveset {
	y := Moveset{}
	for k, v := range ms {
		pp := PowerPoint{Max: v.Max, Current: v.Current}
		y[k] = &pp
	}
	return y
}

func (ms1 Moveset) Equal(ms2 Moveset) bool {
	for k1, v1 := range ms1 {
		v2, ok := ms2[k1]
		if !ok {
			return false
		}

		if *v1 != *v2 {
			return false
		}
	}
	return true
}

type MoveCategory int

const (
	PHYSICS MoveCategory = iota
	SPECIAL
	STATUS
)

func NewMoveCategory(s string) (MoveCategory, error) {
	switch s {
		case "物理":
			return PHYSICS, nil
		case "特殊":
			return SPECIAL, nil
		case "変化":
			return STATUS, nil
		default:
			return -1, fmt.Errorf("不適なmoveCategory")
	}
}

type Target int

const (
	ONE_SELECT Target = iota
	OPPONENT_WHOLE
	SELF
	ALLY_ONE
	WHOLE
	OTHER_THAN_ONESELF
)

func NewTarget(s string) (Target, error) {
	switch s {
		case "１体選択":
			return ONE_SELECT, nil
		case "相手全体":
			return OPPONENT_WHOLE, nil
		case "自分":
			return SELF, nil
		case "味方１体":
			return ALLY_ONE, nil
		case "全体の場":
			return WHOLE, nil
		case "自分以外":
			return OTHER_THAN_ONESELF, nil
		default:
			return -1, fmt.Errorf("不適切なtarget")
	}
}

type StatusMove func(*Battle, *rand.Rand)

func A_sa_no_hi_za_si(bt *Battle, _ *rand.Rand) {
	badWeathers := Weathers{RAIN, SANDSTORM, SNOW}
	var p float64
	if bt.Weather == SUNNY_DAY {
		p = 2.0/3.0
	} else if slices.Contains(badWeathers, bt.Weather) {
		p = 1.0/4.0
	} else {
		p = 1.0/2.0
	}
	heal := int(float64(bt.P1Fighters[0].MaxHP) * p)
	bt.Heal(heal)
}

func A_ma_e_ru(bt *Battle, _ *rand.Rand) {
	bt.Reverse()
	bt.RankStateFluctuation(&RankState{Atk:-2})
	bt.Reverse()
}

func U_so_na_ki(bt *Battle, _ *rand.Rand) {
	bt.Reverse()
	bt.RankStateFluctuation(&RankState{SpDef:-2})
	bt.Reverse()
}

func Gu_ra_su_fi_ru_do(bt *Battle, _ *rand.Rand) {
	if bt.Field.Type != GRASS_FIELD {
		bt.Field = Field{Type:GRASS_FIELD, RemainingTurn:5}
	}
}

func Ko_u_go_u_se_i(bt *Battle, _ *rand.Rand) {
	badWeathers := Weathers{RAIN, SANDSTORM, SNOW}
	var p float64
	if bt.Weather == SUNNY_DAY {
		p = 2.0/3.0
	} else if slices.Contains(badWeathers, bt.Weather) {
		p = 1.0/4.0
	} else {
		p = 1.0/2.0
	}
	heal := int(float64(bt.P1Fighters[0].MaxHP) * p)
	bt.Heal(heal)
}

func Ko_u_so_ku_i_do_u(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Speed:2})
}

func Ko_ra_e_ru(bt *Battle, r *rand.Rand) {
	count := bt.P1Fighters[0].EndureConsecutiveSuccessCount
	p := math.Pow(1.0/3.0, float64(count))
	isSuccess := p > r.Float64()
	bt.P1Fighters[0].IsEndure = isSuccess
	if isSuccess {
		bt.P1Fighters[0].EndureConsecutiveSuccessCount += 1
	} else {
		bt.P1Fighters[0].EndureConsecutiveSuccessCount = 0
	}
}

func Zi_ko_sa_i_se_i(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

func Si_ppo_wo_hu_ru(bt *Battle, _ *rand.Rand) {
	bt.Reverse()
	bt.RankStateFluctuation(&RankState{Def:-1})
	bt.Reverse()
}

func Su_na_a_tu_me(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

func Ta_ma_go_u_mi(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

func Cho_u_ha_tu(bt *Battle, _ *rand.Rand) {
	var turn int
	if bt.IsP2Acted {
		turn = 4
	} else {
		turn = 3
	}
	bt.P2Fighters[0].TauntTurn = turn
}
 
func Tu_ki_no_hi_ka_ri(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

func Tu_me_to_gi(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk:1, Accuracy:1})
}

func Na_ma_ke_ru(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

func Ha_ne_ya_su_me(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

func Mi_ru_ku_no_mi(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

func Do_ku_do_ku(bt *Battle, _ *rand.Rand) {
	if bt.P2Fighters[0].StatusAilment != NO_AILMENT {
		return
	}

	p2PokeTypes := bt.P2Fighters[0].Types

	if slices.Contains(p2PokeTypes, POISON) {
		return
	}

	if slices.Contains(p2PokeTypes, STEEL) {
		return
	}

	bt.P2Fighters[0].StatusAilment = BAD_POISON
}

func Na_ya_mi_no_ta_ne(bt *Battle, _ *rand.Rand) {
	bt.P2Fighters[0].Ability = "ふみん"
}

func Mi_ga_w_ri(bt *Battle, _ *rand.Rand) {
	a := int(float64(bt.P1Fighters[0].MaxHP) * 1.0/4.0)
	if int(bt.P1Fighters[0].CurrentHP) > a && bt.P1Fighters[0].SubstituteDollHP == 0 {
		bt.P1Fighters[0].SubstituteDollHP = a
	}
}

func Ya_do_ri_gi_no_ta_ne(bt *Battle, _ *rand.Rand) {
	if slices.Contains(bt.P2Fighters[0].Types, GRASS) {
		return
	}
	bt.P2Fighters[0].IsLeechSeed = true
}

func Tu_ru_gi_no_ma_i(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk: 2})
}

func Ryu_u_no_ma_i(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk: 1, Speed: 1})
}

func Ka_ra_wo_ya_bu_ru(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk: 2, Def: -1, SpAtk: 2, SpDef: -1, Speed: 2})
}

func Te_ppe_ki(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Def: 2})
}

func Me_i_so_u(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{SpAtk: 1, SpDef: 1})
}

func Wa_ru_da_ku_mi(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{SpAtk:2})
}

var STATUS_MOVES = map[MoveName]StatusMove{
	A_SA_NO_HI_ZA_SI:A_sa_no_hi_za_si,
}
