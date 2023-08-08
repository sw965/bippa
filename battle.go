package bippa

import (
	"math/rand"
	"golang.org/x/exp/slices"
	omwmath "github.com/sw965/omw/math"
	omwrand "github.com/sw965/omw/rand"
)

type Battle struct {
	P1Fighters Fighters
	P2Fighters Fighters
}

func (bt Battle) Reverse() Battle {
	return Battle{P1Fighters: bt.P2Fighters, P2Fighters: bt.P1Fighters}
}

func (bt1 *Battle) Equal(bt2 *Battle) bool {
	return bt1.P1Fighters.Equal(&bt2.P1Fighters) && bt1.P2Fighters.Equal(&bt2.P2Fighters)
}

func (bt *Battle) FinalAccuracy(moveName MoveName) int {
	if slices.Contains(NEVER_MISS_HIT_MOVE_NAMES, moveName) {
		return 100
	}

	a := float64(bt.P1Fighters[0].RankState.Accuracy - bt.P2Fighters[0].RankState.Evasion)
	var bonus float64
	if a <= -6 {
		bonus = 3.0/9.0
	} else if a >= 6 {
		bonus = 9.0/3.
	} else if a <= 0 {
		bonus = 3.0/(3.0-a)
	} else {
		bonus = (3.0+a)/3.0
	}
	b := float64(MOVEDEX[moveName].Accuracy) * bonus
	return omwmath.Min(int(b), 100)
}

func (bt *Battle) CriticalN(moveName MoveName) int {
	rank := MOVEDEX[moveName].CriticalRank
	return CRITICAL_N[rank]
}

func (bt *Battle) IsCritical(moveName MoveName, r *rand.Rand) bool {
	//https://wiki.xn--rckteqa2e.com/wiki/%E6%80%A5%E6%89%80
	n := bt.CriticalN(moveName)
	return r.Intn(n) == 0
}

func (bt *Battle) Damage(dmg int) {
	if bt.P1Fighters[0].IsFaint() {
		return
	}
	dmg = omwmath.Min(dmg, int(bt.P1Fighters[0].CurrentHP))
	bt.P1Fighters[0].CurrentHP -= State(dmg)
	bt.SitrusBerryHeal()
}

func (bt *Battle) Heal(heal int) {
	if bt.P1Fighters[0].IsFaint() {
		return
	}
	heal = omwmath.Max(heal, bt.P1Fighters[0].CurrentDamage())
	bt.P1Fighters[0].CurrentHP += State(heal)
}

func (bt *Battle) SitrusBerryHeal() {
	if bt.P1Fighters[0].Item != "オボンのみ" {
		return
	}

	if bt.P1Fighters[0].IsFaint() {
		return
	}

	max := bt.P1Fighters[0].MaxHP
	current := bt.P1Fighters[0].CurrentHP

	if int(current) <= int(float64(max)*1.0/2.0) {
		bt.P1Fighters[0].Item = EMPTY_ITEM
		heal := int(float64(max) * (1.0 / 4.0))
		bt.Heal(heal)
	}
}

func (bt *Battle) RankStateFluctuation(rs *RankState) {
	if bt.P1Fighters[0].IsFaint() {
		return
	}

	state := bt.P1Fighters[0].RankState.Add(rs)

	if bt.P1Fighters[0].Item == "しろいハーブ" && state.ContainsDown() {
		bt.P1Fighters[0].Item = EMPTY_ITEM
		state = state.ResetDown()
	}

	bt.P1Fighters[0].RankState = state.Regulate()
}

func (bt *Battle) AfterContact() {
	if bt.P2Fighters[0].Ability == "てつのトゲ" {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 8.0)
		bt.Damage(dmg)
	}

	if bt.P2Fighters[0].Item == "ゴツゴツメット" {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 6.0)
		bt.Damage(dmg)
	}
}

// https://latest.pokewiki.net/%E3%83%90%E3%83%88%E3%83%AB%E4%B8%AD%E3%81%AE%E5%87%A6%E7%90%86%E3%81%AE%E9%A0%86%E7%95%AA
func (bt *Battle) MoveUse(moveName MoveName, r *rand.Rand) {
	if bt.P1Fighters[0].IsFaint() {
		return
	}

	if bt.P1Fighters[0].IsFlinch {
		return
	}

	if moveName == STRUGGLE {
		if bt.P2Fighters[0].IsFaint() {
			return
		}
		bt.P1Fighters[0].CurrentHP = 0
		return
	}

	moveData := MOVEDEX[moveName]
	bt.P1Fighters[0].Moveset[moveName].Current -= 1

	if bt.P2Fighters[0].IsFaint() {
		return
	}

	if bt.P1Fighters[0].Item.IsChoice() {
		bt.P1Fighters[0].ChoiceMoveName = moveName
	}

	accuracy := moveData.Accuracy
	if accuracy != -1 {
		if IsHit(accuracy, r) {
			return
		}
	}

	if moveData.Category == STATUS {
		move, ok := STATUS_MOVES[moveName]
		if !ok {
			return
		}
		move(bt, r)
		return
	}

	attackNum := omwrand.Int(moveData.MinAttackNum, moveData.MaxAttackNum+1, r)
	for i := 0; i < attackNum; i++ {
		isCrit := bt.IsCritical(moveName, r)
		dmg := bt.NewFinalDamage(moveName, isCrit, omwrand.Choice(RANDOM_DAMAGE_BONUSES, r))

		if dmg == 0 {
			return
		}

		bt.Reverse()
		bt.Damage(int(dmg))
		bt.Reverse()

		if !bt.P2Fighters[0].IsFlinch {
			bt.P2Fighters[0].IsFlinch = IsHit(moveData.FlinchPercent, r)
		}

		switch moveName {
			case "とんぼがえり":
				bt.P1Fighters[0].AfterUTurn = true
				return
		}

		if moveData.IsContact {
			bt.AfterContact()
		}

		if bt.P1Fighters[0].IsFaint() || bt.P2Fighters[0].IsFaint() {
			break
		}
	}

	if bt.P1Fighters[0].Item == "いのちのたま" {
		dmg := int(float64(bt.P1Fighters[0].MaxHP) * 1.0 / 10.0)
		bt.Damage(dmg)
	}
	return
}

func (bt *Battle) Switch(pokeName PokeName) {
	idx := slices.Index(bt.P1Fighters.PokeNames(), pokeName)

	bt.P1Fighters[0].RankState = RankState{}
	bt.P1Fighters[0].BadPoisonElapsedTurn = 0
	bt.P1Fighters[0].ChoiceMoveName = ""
	bt.P1Fighters[0].IsLeechSeed = false

	fg := bt.P1Fighters

	if idx == 1 {
		bt.P1Fighters[0] = fg[1]
		bt.P1Fighters[1] = fg[0]
		bt.P1Fighters[2] = fg[2]
	} else {
		bt.P1Fighters[0] = fg[2]
		bt.P1Fighters[1] = fg[1]
		bt.P1Fighters[2] = fg[0]
	}

	bt.P1Fighters[0].AfterUTurn = false
}

func (bt *Battle) P1Action(action Action, r *rand.Rand) {
	if action.IsMoveName() || action == Action(STRUGGLE) {
		bt.MoveUse(MoveName(action), r)
	} else {
		bt.Switch(PokeName(action))
	}
}

func (bt *Battle) P2Action(action Action, r *rand.Rand) {
	bt.Reverse()
	if action.IsMoveName() || action == Action(STRUGGLE) {
		bt.MoveUse(MoveName(action), r)
	} else {
		bt.Switch(PokeName(action))
	}
	bt.Reverse()
}

func (bt *Battle) FinalSpeedWinner() Winner {
	p1 := NewFinalSpeed(bt)
	tb := bt.Reverse()
	p2 := NewFinalSpeed(&tb)

	if p1 > p2 {
		return WINNER_PLAYER1
	}

	if p1 < p2 {
		return WINNER_PLAYER2
	}
	return DRAW
}

func (bt *Battle) ActionPriorityWinner(p1Action, p2Action Action) Winner {
	p1 := p1Action.Priority()
	p2 := p2Action.Priority()

	if p1 > p2 {
		return WINNER_PLAYER1
	}

	if p1 < p2 {
		return WINNER_PLAYER2
	}
	return DRAW
}

func (bt *Battle) ActionOrderWinner(p1, p2 Action, r *rand.Rand) Winner {
	priorityWin := bt.ActionPriorityWinner(p1, p2)

	if priorityWin == WINNER_PLAYER1 {
		return WINNER_PLAYER1
	}

	if priorityWin == WINNER_PLAYER2 {
		return WINNER_PLAYER2
	}

	spWin := bt.FinalSpeedWinner()

	if spWin == WINNER_PLAYER1 {
		return WINNER_PLAYER1
	}

	if spWin == WINNER_PLAYER2 {
		return WINNER_PLAYER2
	}

	if omwrand.Bool(r) {
		return WINNER_PLAYER1
	} else {
		return WINNER_PLAYER2
	}
}

func (bt *Battle) TurnEnd(r *rand.Rand) {
	//https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86
	p1First := func(bt *Battle, f func(*Battle)) {
		f(bt)
		bt.Reverse()
		f(bt)
		bt.Reverse()
	}

	p2First := func(bt *Battle, f func(*Battle)) {
		bt.Reverse()
		f(bt)
		bt.Reverse()
		f(bt)
	}

	run := func(fs []func(*Battle) ) {
		spWin := bt.FinalSpeedWinner()
		for _, f := range fs {
			if spWin == WINNER_PLAYER1 {
				p1First(bt, f)
			} else if spWin == WINNER_PLAYER2 {
				p2First(bt, f)
			} else {
				if omwrand.Bool(r) {
					p1First(bt, f)
				} else {
					p2First(bt, f)
				}
			}
		}
	}

	bt.P1Fighters[0].IsFlinch = false
	bt.P2Fighters[0].IsFlinch = false
	run([]func(*Battle) {TurnEndLeftovers, TurnEndBlackSludge})
	run([]func(*Battle) {TurnEndLeechSeed})
	run([]func(*Battle) {TurnEndBadPoison})
}

func (bt Battle) Push(p1Action, p2Action Action, r *rand.Rand) Battle {
	isP1Faint := bt.P1Fighters[0].IsFaint()
	isP2Faint := bt.P2Fighters[0].IsFaint()

	if isP1Faint {
		bt.P1Action(p1Action, r)
		if isP2Faint {
			bt.P2Action(p2Action, r)
		}
		return bt
	}

	if isP2Faint {
		bt.P2Action(p2Action, r)
		return bt
	}

	orderWin := bt.ActionOrderWinner(p1Action, p2Action, r)

	var isP1s []bool

	if orderWin == WINNER_PLAYER1 {
		isP1s = []bool{true, false}
	} else {
		isP1s = []bool{false, true}
	}

	for _, isP1 := range isP1s {
		if isP1 {
			bt.P1Action(p1Action, r)
		} else {
			bt.P2Action(p2Action, r)
		}
	}

	if bt.IsGameEnd() {
		return bt
	}
	bt.TurnEnd(r)
	return bt
}

func (bt *Battle) IsGameEnd() bool {
	return bt.P1Fighters.IsAllFaint() || bt.P2Fighters.IsAllFaint()
}

func (bt *Battle) Winner() Winner {
	isP1AllFaint := bt.P1Fighters.IsAllFaint()
	isP2AllFaint := bt.P2Fighters.IsAllFaint()

	if isP1AllFaint && isP2AllFaint {
		return DRAW
	} else if isP1AllFaint {
		return WINNER_PLAYER2
	} else {
		return WINNER_PLAYER1
	}
}

func (bt *Battle) LegalActionss() []Actions {
	var p1 Actions
	var p2 Actions

	isP1Faint := bt.P1Fighters[0].IsFaint()
	isP2Faint := bt.P2Fighters[0].IsFaint()

	if !isP1Faint {
		p1 = bt.P1Fighters.LegalActions()
	}

	if !isP2Faint {
		p2 = bt.P2Fighters.LegalActions()
	}

	if isP1Faint && isP2Faint {
		p1 = bt.P1Fighters.LegalActions()
		p2 = bt.P2Fighters.LegalActions()
	}

	return []Actions{p1, p2}
}

func (bt *Battle) NewFinalDamage(moveName MoveName, isCrit bool, randDmgBonus RandomDamageBonus) FinalDamage {
	attaker := bt.P1Fighters[0]
	defender := bt.P2Fighters[0]
	return NewFinalDamage(&attaker, &defender, moveName, isCrit, randDmgBonus)
}