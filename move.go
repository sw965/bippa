package bippa

import (
	"fmt"
	"golang.org/x/exp/slices"
	osliecs "github.com/sw965/omw/slices"
	omwmaps "github.com/sw965/omw/maps"
)

type MoveName int

const (
    EMPTY_MOVE_NAME MoveName = iota
    THUNDERBOLT    // 10まんボルト
    HAMMER_ARM     // アームハンマー
    STONE_EDGE     // ストーンエッジ
    SURF           // なみのり
    ICE_BEAM       // れいとうビーム
    STRUGGLE       // わるあがき
    RAIN_DANCE     // あまごい
    ROCK_SLIDE     // いわなだれ
    RETURN         // おんがえし
    CRUNCH         // かみくだく
    ENDEAVOR       // がむしゃら
    ICY_WIND       // こごえるかぜ
    FOLLOW_ME      // このゆびとまれ
    HYPNOSIS       // さいみんじゅつ
    RECOVER        // じこあんじ
    EARTHQUAKE     // じしん
    SELF_DESTRUCT  // じばく
    WATERFALL      // たきのぼり
    EXPLOSION      // だいばくはつ
    TAUNT          // ちょうはつ
    THUNDER_WAVE   // でんじは
    FAKE_OUT       // ねこだまし
    HEAT_WAVE      // ねっぷう
    BELLY_DRUM     // はらだいこ
    SUCKER_PUNCH   // ふいうち
    FIRE_PUNCH     // ほのおのパンチ
    PROTECT        // まもる
    SUBSTITUTE     // みがわり
    DRACO_METEOR   // りゅうせいぐん
    CROSS_CHOP     // クロスチョップ
    COMET_PUNCH    // コメットパンチ
    PSYCHIC        // サイコキネシス
    GYRO_BALL      // ジャイロボール
    DARK_VOID      // ダークホール
    TRICK_ROOM     // トリックルーム
    HYDRO_PUMP     // ハイドロポンプ
    BULLET_PUNCH   // バレットパンチ
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

var MOVE_NAME_TO_STRING = omwmaps.Invert[map[MoveName]string](STRING_TO_MOVE_NAME)

func (name MoveName) ToString() string {
	return MOVE_NAME_TO_STRING[name]
}

type MoveNames []MoveName

func (names MoveNames) ToStrings() []string {
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = name.ToString()
	}
	return ret
}

func (names MoveNames) Sort() MoveNames {
	ret := make(MoveNames, len(names))
	for i := 0; i < osliecs.Count(names, EMPTY_MOVE_NAME); i++ {
		ret = append(ret, EMPTY_MOVE_NAME)
	}

	for _, name := range ALL_MOVE_NAMES {
		if slices.Contains(names, name) {
			ret = append(ret, name)
		}
	}
	return ret
}

type MoveNamess []MoveNames

type MoveCategory int

const (
	PHYSICS MoveCategory = iota
	SPECIAL
	STATUS
)

var STRING_TO_MOVE_CATEGORY = map[string]MoveCategory{
	"物理":PHYSICS,
	"特殊":SPECIAL,
	"変化":STATUS,
}

var MOVE_CATEGORY_TO_STRING = omwmaps.Invert[map[MoveCategory]string](STRING_TO_MOVE_CATEGORY)

func (m MoveCategory) ToString() string {
	return MOVE_CATEGORY_TO_STRING[m]
}

type TargetRange int

const (
    NORMAL_TARGET TargetRange = iota // 通常
    OPPONENT_TWO_TARGET              // 相手2体
    SELF_TARGET                      // 自分
    OTHERS_TARGET                    // 自分以外
    ALL_TARGET                       // 全体
    OPPONENT_RANDOM_ONE_TARGET       // 相手ランダム1体
)

var STRING_TO_TARGET_RANGE = map[string]TargetRange{
	"通常":NORMAL_TARGET,
	"相手2体":OPPONENT_TWO_TARGET,
	"自分":SELF_TARGET,
	"自分以外":OTHERS_TARGET,
	"全体":ALL_TARGET,
	"相手ランダム1体":OPPONENT_RANDOM_ONE_TARGET,
}

var TARGET_RANGE_TO_STRING = omwmaps.Invert[map[TargetRange]string](STRING_TO_TARGET_RANGE)

func (t TargetRange) ToString() string {
	return TARGET_RANGE_TO_STRING[t]
}

type PointUp int

const (
	MAX_POINT_UP = 3
)

type PointUps []PointUp

var (
	MAX_POINT_UPS = PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP}
)

type PowerPoint struct {
	Max int
	Current int
}

const (
	MIN_MOVESET = 1
	MAX_MOVESET = 4
)

func NewPowerPoint(base int, up PointUp) PowerPoint {
    increment := int(float64(base) / 5.0)
    max := base + (increment * int(up))
	return PowerPoint{Max:max, Current:max}
}

type Moveset map[MoveName]*PowerPoint

func NewMoveset(pokeName PokeName, moveNames MoveNames) (Moveset, error) {
	if len(moveNames) == 0 {
		msg := fmt.Sprintf("覚えさせる技が指定されていません。ポケモンには、少なくとも%dつ以上の技を覚えさせる必要があります。", MIN_MOVESET)
		return Moveset{}, fmt.Errorf(msg)
	}

	if len(moveNames) > MAX_MOVESET {
		msg := fmt.Sprintf("覚えさせる技の数が、限度を超えています。技は最大で%dつまで覚えさせることが出来ます。", MAX_MOVESET)
		return Moveset{}, fmt.Errorf(msg)
	}

	learnset := POKEDEX[pokeName].Learnset
	moveset := Moveset{}
	for i := range moveNames {
		moveName := moveNames[i]
		if !slices.Contains(learnset, moveNames[i]) {
			pokeNameStr := POKE_NAME_TO_STRING[pokeName]
			moveNameStr := MOVE_NAME_TO_STRING[moveName]
			msg := fmt.Sprintf("「%s」 は 「%s」 を覚えることができません。覚えられる技を選択してください。", pokeNameStr, moveNameStr)
			return Moveset{}, fmt.Errorf(msg)
		}
		basePP := MOVEDEX[moveName].BasePP
		moveset[moveName] = &PowerPoint{Max:basePP, Current:basePP}
	}
	return moveset, nil
}

func (m Moveset) Equal(other Moveset) bool {
	for moveName, pp := range m {
		otherPP, ok := other[moveName]
		if !ok {
			return false
		}
		if *pp != *otherPP {
			return false
		}
	}
	return true
}

func (m Moveset) Clone() Moveset {
	clone := Moveset{}
	for moveName, pp := range m {
		clone[moveName] = &PowerPoint{Max:pp.Max, Current:pp.Current}
	}
	return clone
}

func (m Moveset) ToEasyRead() EasyReadMoveset {
	ret := EasyReadMoveset{}
	for moveName, pp := range m {
		ret[moveName.ToString()] = *pp
	}
	return ret
}