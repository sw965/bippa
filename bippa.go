package bippa

import (
	"fmt"
	omwmaps "github.com/sw965/omw/maps"
	"github.com/sw965/omw/fn"
	omwjson "github.com/sw965/omw/json"
	"golang.org/x/exp/slices"
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
)

var STRING_TO_MOVE_NAME = map[string]MoveName{
    "":              EMPTY_MOVE_NAME,
    "10まんボルト":   THUNDERBOLT,
    "アームハンマー": HAMMER_ARM,
    "ストーンエッジ": STONE_EDGE,
    "なみのり":      SURF,
    "れいとうビーム":  ICE_BEAM,
    "わるあがき":    STRUGGLE,
    "あまごい":      RAIN_DANCE,
    "いわなだれ":    ROCK_SLIDE,
    "おんがえし":    RETURN,
    "かみくだく":    CRUNCH,
    "がむしゃら":    ENDEAVOR,
    "こごえるかぜ":  ICY_WIND,
    "このゆびとまれ": FOLLOW_ME,
    "さいみんじゅつ": HYPNOSIS,
    "じこあんじ":    RECOVER,
    "じしん":        EARTHQUAKE,
    "じばく":        SELF_DESTRUCT,
    "たきのぼり":    WATERFALL,
    "だいばくはつ":  EXPLOSION,
    "ちょうはつ":    TAUNT,
    "でんじは":      THUNDER_WAVE,
    "ねこだまし":    FAKE_OUT,
    "ねっぷう":      HEAT_WAVE,
    "はらだいこ":    BELLY_DRUM,
    "ふいうち":      SUCKER_PUNCH,
    "ほのおのパンチ": FIRE_PUNCH,
    "まもる":        PROTECT,
    "みがわり":      SUBSTITUTE,
    "りゅうせいぐん": DRACO_METEOR,
    "クロスチョップ": CROSS_CHOP,
    "コメットパンチ": COMET_PUNCH,
    "サイコキネシス": PSYCHIC,
    "ジャイロボール": GYRO_BALL,
    "ダークホール":  DARK_VOID,
    "トリックルーム": TRICK_ROOM,
    "ハイドロポンプ": HYDRO_PUMP,
    "バレットパンチ": BULLET_PUNCH,
}

func StringToMoveName(s string) (MoveName, error) {
	if moveName, ok := STRING_TO_MOVE_NAME[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_MOVE_NAME に含まれていない為、MoveNameに変換出来ません。", s)
		return moveName, fmt.Errorf(msg)
	} else {
		return moveName, nil
	}
}

func StringsToMoveNames(ss []string) (MoveNames, error) {
	return fn.MapWithError[MoveNames](ss, StringToMoveName)
}

var STRING_TO_MOVE_CATEGORY = map[string]MoveCategory{
	"物理":PHYSICS,
	"特殊":SPECIAL,
	"変化":STATUS,
}

func StringToMoveCategory(s string) (MoveCategory, error) {
	if category, ok := STRING_TO_MOVE_CATEGORY[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_MOVE_CATEGORY に含まれていない為、MoveCategoryに変換出来ません。", s)
		return category, fmt.Errorf(msg)
	} else {
		return category, nil
	}
}

var STRING_TO_TARGET_RANGE = map[string]MoveTarget{
	"通常":NORMAL_TARGET,
	"相手2体":OPPONENT_TWO_TARGET,
	"自分":SELF_TARGET,
	"自分以外":OTHERS_TARGET,
	"全体":ALL_TARGET,
	"相手ランダム1体":OPPONENT_RANDOM_ONE_TARGET,
}

func StringToType(s string) (Type, error) {
	if t, ok := STRING_TO_TYPE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_TYPE に含まれていない為、Typeに変換出来ません。", s)
		return t, fmt.Errorf(msg)
	} else {
		return t, nil
	}
}

func StringsToTypes(ss []string) (Types, error) {
	return fn.MapWithError[Types](ss, StringToType)
}

func StringToTargetRange(s string) (MoveTarget, error) {
	if target, ok := STRING_TO_TARGET_RANGE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_TARGET_RANGE に含まれていない為、TargetRangeに変換出来ません。", s)
		return target, fmt.Errorf(msg)
	} else {
		return target, nil
	}
}

func StringToNature(s string) (Nature, error) {
	if nature, ok := STRING_TO_NATURE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_NATURE に含まれていない為、Natureに変換出来ません。", s)
		return nature, fmt.Errorf(msg)
	} else {
		return nature, nil
	}
}

func StringsToNatures(ss []string) (Natures, error) {
	return fn.MapWithError[Natures](ss, StringToNature)
}

func StringsToAbilities(ss []string) (Abilities, error) {
	return fn.MapWithError[Abilities](ss, StringToAbility)
}

func StringToItem(s string) (Item, error) {
	if item, ok := STRING_TO_ITEM[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_ITEM に含まれていない為、Itemに変換出来ません。", s)
		return item, fmt.Errorf(msg)
	} else {
		return item, nil
	}
}

func StringsToItems(ss []string) (Items, error) {
	return fn.MapWithError[Items](ss, StringToItem)
}

//ギャラドス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Gyarados() Pokemon {
	pokemon, err := NewPokemon(
		GYARADOS, STANDARD_LEVEL, JOLLY, INTIMIDATE, WACAN_BERRY,
		MoveNames{WATERFALL, STONE_EDGE, THUNDER_WAVE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP:164, Atk:92, Speed:MAX_EFFORT},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//カビゴン
func NewMoruhu2007Snorlax() Pokemon {
	ivStat := MAX_INDIVIDUAL_STAT.Clone()
	ivStat.Speed = MAX_INDIVIDUAL
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, FIRE_PUNCH, BELLY_DRUM, SUBSTITUTE},
		MAX_POINT_UPS,
		&ivStat, &HB252_D4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//カビゴン
func NewMoruhu2008Snorlax() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, RELAXED, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, FIRE_PUNCH, BELLY_DRUM, PROTECT},
		MAX_POINT_UPS,
		&ivStat, &HB252_D4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}


//カビゴン
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Snorlax() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&ivStat, &EffortStat{HP:204, Atk:52, Def:156, SpDef:96},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//カビゴン
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Snorlax() Pokemon {
	pokemon, err := NewPokemon(
		SNORLAX, STANDARD_LEVEL, BRAVE, THICK_FAT, SITRUS_BERRY,
		MoveNames{RETURN, CRUNCH, SELF_DESTRUCT, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT,
		&EffortStat{HP:204, Atk:52, Def:156, SpDef:60, Speed:36},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドーブル
func NewMoruhu2007Smeargle() Pokemon {
	pokemon, err := NewPokemon(
		SMEARGLE, MIN_LEVEL, BRAVE, OWN_TEMPO, FOCUS_SASH,
		MoveNames{FAKE_OUT, FOLLOW_ME, DARK_VOID, ENDEAVOR},
		MAX_POINT_UPS,
		&MIN_INDIVIDUAL_STAT, &EffortStat{},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドーブル
func NewMoruhu2008Smeargle() Pokemon {
	return NewMoruhu2007Smeargle()
}

//ボーマンダ
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagiSalamence2009() Pokemon {
	pokemon, err := NewPokemon(
		SALAMENCE, STANDARD_LEVEL, MODEST, INTIMIDATE, SITRUS_BERRY,
		MoveNames{DRACO_METEOR, HEAT_WAVE, RAIN_DANCE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:20, SpAtk:236, Speed:252},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//メタグロス
func NewMoruhu2007Metagross() Pokemon {
	ivStat := MAX_INDIVIDUAL_STAT.Clone()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, BRAVE, CLEAR_BODY, LUM_BERRY,
		MoveNames{EARTHQUAKE, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&ivStat, &EffortStat{HP:MAX_EFFORT, Def:128, SpDef:128},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//メタグロス
func NewMoruhu2008Metagross() Pokemon {
	ivStat := MAX_INDIVIDUAL_STAT.Clone()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, BRAVE, CLEAR_BODY, LUM_BERRY,
		MoveNames{HAMMER_ARM, BULLET_PUNCH, ROCK_SLIDE, RECOVER},
		MAX_POINT_UPS,
		&ivStat, &EffortStat{HP:MAX_EFFORT, Def:128, SpDef:128},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//メタグロス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Metagross() Pokemon {
	pokemon, err := NewPokemon(
		METAGROSS, STANDARD_LEVEL, ADAMANT, CLEAR_BODY, LUM_BERRY,
		MoveNames{COMET_PUNCH, BULLET_PUNCH, EARTHQUAKE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:236, Atk:36, Def:4, SpDef:172, Speed:60},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ラティオス
// https://matsu-1129.hatenadiary.org/entry/20090308/1236586122
func NewRomanStan2009Latios() Pokemon {
	pokemon, err := NewPokemon(
		LATIOS, STANDARD_LEVEL, TIMID, LEVITATE, FOCUS_SASH,
		MoveNames{DRACO_METEOR, THUNDERBOLT, RAIN_DANCE, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &CS252_H4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//エンペルト
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Empoleon() Pokemon {
	pokemon, err := NewPokemon(
		EMPOLEON, STANDARD_LEVEL, MODEST, TORRENT, WACAN_BERRY,
		MoveNames{HYDRO_PUMP, SURF, ICY_WIND, PROTECT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &EffortStat{HP:68, Def:12, SpAtk:252, SpDef:4, Speed:172},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドータクン
func NewMoruhu2007Bronzong() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		BRONZONG, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
		MoveNames{GYRO_BALL, EXPLOSION, TRICK_ROOM, HYPNOSIS},
		MAX_POINT_UPS,
		&ivStat, &HD252_B4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドータクン
func NewMoruhu2008Bronzong() Pokemon {
	ivStat := NewMaxIndividualStat()
	ivStat.Speed = MIN_INDIVIDUAL
	pokemon, err := NewPokemon(
		BRONZONG, STANDARD_LEVEL, SASSY, HEATPROOF, CHESTO_BERRY,
		MoveNames{PSYCHIC, EXPLOSION, TRICK_ROOM, HYPNOSIS},
		MAX_POINT_UPS,
		&ivStat, &HD252_B4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//ドクロッグ
// https://detail.chiebukuro.yahoo.co.jp/qa/question_detail/q1267938327
func NewKusanagi2009Toxicroak() Pokemon {
	pokemon, err := NewPokemon(
		TOXICROAK, STANDARD_LEVEL, ADAMANT, DRY_SKIN, FOCUS_SASH,
		MoveNames{CROSS_CHOP, SUCKER_PUNCH, FAKE_OUT, TAUNT},
		MAX_POINT_UPS,
		&MAX_INDIVIDUAL_STAT, &AS252_B4,
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

var ABILITY_TO_STRING = omwmaps.Invert[map[Ability]string](STRING_TO_ABILITY)
var MOVE_NAME_TO_STRING = omwmaps.Invert[map[MoveName]string](STRING_TO_MOVE_NAME)
var MOVE_CATEGORY_TO_STRING = omwmaps.Invert[map[MoveCategory]string](STRING_TO_MOVE_CATEGORY)

type TrainerName string

const (
	LEAF = "リーフ"
	LYRA = "コトネ"
	MAY = "ハルカ"
	DAWN = "ヒカリ"
	HILDA = "トウコ"
	ROSA = "メイ"
	SERENA = "セレナ"
	SELENE = "ミズキ"
	GLORIA = "ユウリ"
)

type Type int

const (
	NORMAL Type  = iota
	FIRE
	WATER
	GRASS
	ELECTRIC
	ICE
	FIGHTING
	POISON
	GROUND
	FLYING
	PSYCHIC_TYPE
	BUG
	ROCK
	GHOST
	DRAGON
	DARK
	STEEL
	FAIRY
)

var STRING_TO_TYPE = map[string]Type{
	"ノーマル":NORMAL,
	"ほのお":FIRE,
	"みず":WATER,
	"くさ":GRASS,
	"でんき":ELECTRIC,
	"こおり":ICE,
	"かくとう":FIGHTING,
	"どく":POISON,
	"じめん":GROUND,
	"ひこう":FLYING,
	"エスパー":PSYCHIC_TYPE,
	"むし":BUG,
	"いわ":ROCK,
	"ゴースト":GHOST,
	"ドラゴン":DRAGON,
	"あく":DARK,
	"はがね":STEEL,
	"フェアリー":FAIRY,
}

var TYPE_TO_STRING = omwmaps.Invert[map[Type]string](STRING_TO_TYPE)

func (t Type) ToString() string {
	return TYPE_TO_STRING[t]
}

type Types []Type

var ALL_TYPES = func() Types {
	buff, err := omwjson.Load[[]string](ALL_TYPES_PATH)
	if err != nil {
		panic(err)
	}
	ret := make(Types, len(buff))
	for i, s := range buff {
		ret[i] = STRING_TO_TYPE[s]
	}
	return ret
}()

func (ts Types) ToStrings() []string {
	ret := make([]string, len(ts))
	for i, t := range ts {
		ret[i] = t.ToString()
	}
	return ret
}

func (ts Types) Sort() Types {
	ret := slices.Clone(ts)
	slices.SortFunc(ret, func(t1, t2 Type) bool { return slices.Index(ALL_TYPES, t1) < slices.Index(ALL_TYPES, t2) } )
	return ret
}

type TypesSlice []Types

var ALL_TWO_TYPESS = func() TypesSlice {
	return omwslices.Concat(
		fn.Map[TypesSlice](ALL_TYPES, func(t Type) Types { return Types{t} }),
		omwslices.Combination[TypesSlice, Types](ALL_TYPES, 2),
	)
}()