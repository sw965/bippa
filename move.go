package bippa

import (
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
	CanSubstitute bool
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

type EasyReadMoveData struct {
    Type         string
    Category     string
    Power        int
    Accuracy     int
    BasePP       int
	IsContact    bool
	PriorityRank int
	CriticalRank CriticalRank
	Target       string
	CanSubstitute bool
}

func (m *EasyReadMoveData) From() (MoveData, error) {
	t, err := StringToType(m.Type)
	if err != nil {
		return MoveData{}, err
	}

	category, err := StringToMoveCategory(m.Category)
	if err != nil {
		return MoveData{}, err
	}

	target, err := StringToMoveTarget(m.Target)
	if err != nil {
		return MoveData{}, err
	}

	return MoveData{
		Type:t,
		Category:category,
		Power:m.Power,
		Accuracy:m.Accuracy,
		BasePP:m.BasePP,
		IsContact:m.IsContact,
		PriorityRank:m.PriorityRank,
		CriticalRank:m.CriticalRank,
		Target:target,
		CanSubstitute:m.CanSubstitute,
	}, nil
}

type EasyReadMovedex map[string]EasyReadMoveData

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

func (mn MoveName) ToString() string {
	return MOVE_NAME_TO_STRING[mn]
}

type MoveNames []MoveName

var ALL_MOVE_NAMES = func() MoveNames {
	buff, err := omwjson.Load[[]string](ALL_MOVE_NAMES_PATH)
	if err != nil {
		panic(err)
	}

	mns, err := StringsToMoveNames(buff)
	if err != nil {
		panic(err)
	}
	return mns
}()

func (mns MoveNames) ToStrings() []string {
	ss := make([]string, len(mns))
	for i, mn := range mns {
		ss[i] = mn.ToString()
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

func (c MoveCategory) ToString() string {
	return MOVE_CATEGORY_TO_STRING[c]
}

type MoveTarget int

const (
    NORMAL_TARGET MoveTarget = iota  // 通常
    OPPONENT_TWO_TARGET              // 相手2体
    SELF_TARGET                      // 自分
    OTHERS_TARGET                    // 自分以外
    ALL_TARGET                       // 全体
    OPPONENT_RANDOM_ONE_TARGET       // 相手ランダム1体
)

var MOVE_TARGET_TO_STRING = omwmaps.Invert[map[MoveTarget]string](STRING_TO_MOVE_TARGET)

func (t MoveTarget) ToString() string {
	return MOVE_TARGET_TO_STRING[t]
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
	d := Movedex{}
	for _, name := range ALL_MOVE_NAMES {
		data, err := LoadMoveData(name)
		if err != nil {
			panic(err)
		}
		d[name] = &data
	}
	return d
}()

func (m Movedex) ToEasyRead() EasyReadMovedex {
	e := EasyReadMovedex{}
	for name, data := range m {
		e[name.ToString()] = data.ToEasyRead()
	}
	return e
}

type PowerPoint struct {
	Max int
	Current int
}

const (
	MIN_MOVESET_LENGTH = 1
	MAX_MOVESET_LENGTH = 4
)

func NewPowerPoint(base int, up PointUp) PowerPoint {
    v := int(float64(base) / 5.0)
    max := base + (v * int(up))
	return PowerPoint{Max:max, Current:max}
}

type Moveset map[MoveName]*PowerPoint

func (m Moveset) Equal(other Moveset) bool {
	for k, v1 := range m {
		v2, ok := other[k]
		if !ok {
			return false
		}
		if *v1 != *v2 {
			return false
		}
	}
	return true
}

func (m Moveset) Clone() Moveset {
	c := Moveset{}
	for k, v := range m {
		c[k] = &PowerPoint{Max:v.Max, Current:v.Current}
	}
	return c
}

func (m Moveset) ToEasyRead() EasyReadMoveset {
	e := EasyReadMoveset{}
	for k, v := range m {
		e[k.ToString()] = *v
	}
	return e
}

type EasyReadMoveset map[string]PowerPoint

func (e EasyReadMoveset) From() (Moveset, error) {
	m := Moveset{}
	for k, v := range e {
		moveName, err := StringToMoveName(k)
		if err != nil {
			return Moveset{}, err
		}
		pp := PowerPoint{Max:v.Max, Current:v.Current}
		m[moveName] = &pp
	}
	return m, nil
}