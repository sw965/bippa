package bippa

import (
	"golang.org/x/exp/slices"
	"math/rand"
	"github.com/sw965/omw"
	"fmt"
)

type MoveName string

const (
	EMPTY_MOVE_NAME = MoveName("なし")
	STRUGGLE        = MoveName("わるあがき")
)

type MoveNames []MoveName

func (mns MoveNames) Sort() {
	isSwap := func(moveName1, moveName2 MoveName) bool {
		return slices.Index(ALL_MOVE_NAMES, moveName1) > slices.Index(ALL_MOVE_NAMES, moveName2)
	}
	slices.SortFunc(mns, isSwap)
}

type PowerPointUp int

var (
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

var ALL_POWER_POINT_UPS = omw.MakeIntegerRange[PowerPointUps](MIN_POWER_POINT_UP, MAX_POWER_POINT_UP+1, 1)

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

type StatusMove func(Battle, *rand.Rand) Battle

// あさのひざし
func MorningSun(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// こうごうせい
func Synthesis(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// じこさいせい
func Recover(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// すなあつめ
func ShoreUp(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// タマゴうみ
func SoftBoiled(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// つきのひかり
func Moonlight(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// なまける
func SlackOff(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// はねやすめ
func Roost(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// ミルクのみ
func MilkDrink(bt Battle, _ *rand.Rand) Battle {
	heal := int(float64(bt.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return bt.Heal(heal)
}

// どくどく
func Toxic(bt Battle, _ *rand.Rand) Battle {
	if bt.P2Fighters[0].StatusAilment != "" {
		return bt
	}

	p2PokeTypes := bt.P2Fighters[0].Types

	if slices.Contains(p2PokeTypes, POISON) {
		return bt
	}

	if slices.Contains(p2PokeTypes, STEEL) {
		return bt
	}

	bt.P2Fighters[0].StatusAilment = BAD_POISON
	return bt
}

// やどりぎのタネ
func LeechSeed(bt Battle, _ *rand.Rand) Battle {
	if slices.Contains(bt.P2Fighters[0].Types, GRASS) {
		return bt
	}
	bt.P2Fighters[0].IsLeechSeed = true
	return bt
}

// つるぎのまい
func SwordsDance(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Atk: 2})
}

// りゅうのまい
func DragonDance(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Atk: 1, Speed: 1})
}

// からをやぶる
func ShellSmash(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Atk: 2, Def: -1, SpAtk: 2, SpDef: -1, Speed: 2})
}

// てっぺき
func IronDefense(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{Def: 2})
}

// めいそう
func CalmMind(bt Battle, _ *rand.Rand) Battle {
	return bt.RankStateFluctuation(&RankState{SpAtk: 1, SpDef: 1})
}

var STATUS_MOVES = map[MoveName]StatusMove{
	"あさのひざし":  Moonlight,
	"こうごうせい":  Synthesis,
	"じこさいせい":  Recover,
	"すなあつめ":   ShoreUp,
	"タマゴうみ":   SoftBoiled,
	"つきのひかり":  Moonlight,
	"なまける":    SlackOff,
	"はねやすめ":   Roost,
	"ミルクのみ":   MorningSun,
	"どくどく":    Toxic,
	"やどりぎのタネ": LeechSeed,
	"つるぎのまい":  SwordsDance,
	"りゅうのまい":  DragonDance,
	"からをやぶる":  ShellSmash,
	"てっぺき":    IronDefense,
}
