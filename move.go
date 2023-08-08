package bippa

import (
	"golang.org/x/exp/slices"
	"math/rand"
	"fmt"
)

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

type MoveName string

const (
	EMPTY_MOVE_NAME = MoveName("なし")
	STRUGGLE        = MoveName("わるあがき")
)

type MoveNames []MoveName

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
		if moveName == "" {
			return Moveset{}, fmt.Errorf("技名 が 空 に なっている")
		}

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

type StatusMove func(*Battle, *rand.Rand)

// あさのひざし
func MorningSun(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

// こうごうせい
func Synthesis(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

// じこさいせい
func Recover(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

//しっぽをふる
func TailWhip(bt *Battle, _ *rand.Rand) {
	bt.Reverse()
	bt.RankStateFluctuation(&RankState{Def:-1})
	bt.Reverse()
}

// すなあつめ
func ShoreUp(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

//タマゴうみ
func SoftBoiled(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

//つきのひかり
func Moonlight(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

//つめとぎ
func HoneClaws(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk:1, Accuracy:1})
}

// なまける
func SlackOff(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

// はねやすめ
func Roost(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

// ミルクのみ
func MilkDrink(bt *Battle, _ *rand.Rand) {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	bt.Heal(heal)
}

// どくどく
func Toxic(bt *Battle, _ *rand.Rand) {
	if bt.P2Fighters[0].StatusAilment != "" {
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

//なやみのタネ
func WorrySeed(bt *Battle, _ *rand.Rand) {
	bt.P2Fighters[0].Ability = "ふみん"
}

// やどりぎのタネ
func LeechSeed(bt *Battle, _ *rand.Rand) {
	if slices.Contains(bt.P2Fighters[0].Types, GRASS) {
		return
	}
	bt.P2Fighters[0].IsLeechSeed = true
}

// つるぎのまい
func SwordsDance(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk: 2})
}

// りゅうのまい
func DragonDance(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk: 1, Speed: 1})
}

// からをやぶる
func ShellSmash(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Atk: 2, Def: -1, SpAtk: 2, SpDef: -1, Speed: 2})
}

// てっぺき
func IronDefense(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{Def: 2})
}

// めいそう
func CalmMind(bt *Battle, _ *rand.Rand) {
	bt.RankStateFluctuation(&RankState{SpAtk: 1, SpDef: 1})
}

var STATUS_MOVES = map[MoveName]StatusMove{
	"あさのひざし":  Moonlight,
	"こうごうせい":  Synthesis,
	"じこさいせい":  Recover,
	"しっぽをふる": TailWhip,
	"すなあつめ":   ShoreUp,
	"タマゴうみ":   SoftBoiled,
	"つきのひかり":  Moonlight,
	"つめとぎ":HoneClaws,
	"なまける":    SlackOff,
	"なやみのタネ":WorrySeed,
	"はねやすめ":   Roost,
	"ミルクのみ":   MorningSun,
	"どくどく":    Toxic,
	"やどりぎのタネ": LeechSeed,
	"つるぎのまい":  SwordsDance,
	"りゅうのまい":  DragonDance,
	"からをやぶる":  ShellSmash,
	"てっぺき":    IronDefense,
}
