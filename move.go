package bippa

import (
	"fmt"
	"golang.org/x/exp/slices"
	omwjson "github.com/sw965/omw/json"
	omwmaps "github.com/sw965/omw/maps"
)

type MoveData struct {
	Type Type
	Category MoveCategory
	Power int
	Accuracy int
	BasePP int
	IsContact bool
    PriorityRank int
    CriticalRank CriticalRank
    Target MoveTarget
}

func LoadMoveData(moveName MoveName) (MoveData, error) {
	path := MOVE_DATA_PATH + MOVE_NAME_TO_STRING[moveName] + omwjson.EXTENSION
	buff, err := omwjson.Load[EasyReadMoveData](path)
	if err != nil {
		return MoveData{}, err
	}
	return buff.From()
}

func (m *MoveData) ToEasyRead() EasyReadMoveData {
	return EasyReadMoveData{
		Type:m.Type.ToString(),
		Category:m.Category.ToString(),
		Power:m.Power,
		Accuracy:m.Accuracy,
		BasePP:m.BasePP,
		IsContact:m.IsContact,
		PriorityRank:m.PriorityRank,
		CriticalRank:m.CriticalRank,
		Target:m.Target.ToString(),
	}
}

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

func (name MoveName) ToString() string {
	return MOVE_NAME_TO_STRING[name]
}

type MoveNames []MoveName

var ALL_MOVE_NAMES = func() MoveNames {
	buff, err := omwjson.Load[[]string](ALL_MOVE_NAMES_PATH)
	if err != nil {
		panic(err)
	}

	ret, err := StringsToMoveNames(buff)
	if err != nil {
		panic(err)
	}
	return ret
}()

func (ns MoveNames) ToStrings() []string {
	ss := make([]string, len(ns))
	for i, n := range ns {
		ss[i] = n.ToString()
	}
	return ss
}

type MoveNamesSlice []MoveNames

type MoveCategory int

const (
	PHYSICS MoveCategory = iota
	SPECIAL
	STATUS
)

func (m MoveCategory) ToString() string {
	return MOVE_CATEGORY_TO_STRING[m]
}

type MoveTarget int

const (
    NORMAL_TARGET MoveTarget = iota // 通常
    OPPONENT_TWO_TARGET              // 相手2体
    SELF_TARGET                      // 自分
    OTHERS_TARGET                    // 自分以外
    ALL_TARGET                       // 全体
    OPPONENT_RANDOM_ONE_TARGET       // 相手ランダム1体
)

var TARGET_RANGE_TO_STRING = omwmaps.Invert[map[MoveTarget]string](STRING_TO_TARGET_RANGE)

func (t MoveTarget) ToString() string {
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

type Movedex map[MoveName]*MoveData

var MOVEDEX = func() Movedex {
	ret := Movedex{}
	for _, name := range ALL_MOVE_NAMES {
		data, err := LoadMoveData(name)
		if err != nil {
			panic(err)
		}
		ret[name] = &data
	}
	return ret
}()

func (m Movedex) ToEasyRead() EasyReadMovedex {
	ret := EasyReadMovedex{}
	for name, data := range m {
		ret[name.ToString()] = data.ToEasyRead()
	}
	return ret
}

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