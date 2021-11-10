package bippa

import (
	"fmt"
	"math/rand"
	"github.com/sw965/omw"
)

type SelfPointOfViewBattle struct {
	SelfFighters Fighters
	OpponentFighters Fighters

	SelfField OneSideField
	OpponentField OneSideField
}

func NewSelfPointOfView(selfFighters, opponentFIghters Fighters, needUI bool) SelfPointOfViewBattle {
	return SelfPointOfViewBattle{SelfFighters:selfFighters, OpponentFighters:opponentFIghters,
		SelfField:OneSideField{}, OpponentField:OneSideField{}}
}

func (spovb *SelfPointOfViewBattle) SwitchPointOfView() SelfPointOfViewBattle {
	return SelfPointOfViewBattle{SelfFighters:spovb.OpponentFighters, OpponentFighters:spovb.SelfFighters,
		SelfField:spovb.OpponentField, OpponentField:spovb.SelfField}
}

func (spovb *SelfPointOfViewBattle) Accuracy(moveName MoveName, moveData *MoveData) int {
	if moveName == "どくどく" && spovb.SelfFighters[0].Types.In(POISON) {
		return 100
	} else {
		return moveData.Accuracy
	}
}

func (spovb *SelfPointOfViewBattle) IsCritical(moveName MoveName, random *rand.Rand) bool {
	if ONE_HIT_KO_MOVE_NAMES.In(moveName) {
		return false
	} else {
		return IsCritical(random)
	}
}

func (spovb SelfPointOfViewBattle) ToDamage(damage int) SelfPointOfViewBattle {
	currentHP := spovb.SelfFighters[0].State.CurrentHP
	intCurrentHP := int(currentHP)

	if damage > intCurrentHP {
		damage = intCurrentHP
	}

	spovb.SelfFighters[0].State.CurrentHP -= State_(damage)
	return spovb.SitrusBerryHeal()
}

func (spovb SelfPointOfViewBattle) Heal(heal int) SelfPointOfViewBattle {
	currentDamage := spovb.SelfFighters[0].CurrentDamage()

	if currentDamage < heal {
		heal = currentDamage
	}

	spovb.SelfFighters[0].State.CurrentHP += State_(heal)
	return spovb
}

func (spovb SelfPointOfViewBattle) SitrusBerryHeal() SelfPointOfViewBattle {
	if spovb.SelfFighters[0].Item != "オボンのみ" {
		return spovb
	}

	if spovb.SelfFighters[0].IsFaint() {
		return spovb
	}

	maxHP := spovb.SelfFighters[0].State.MaxHP
	currentHP := spovb.SelfFighters[0].State.CurrentHP

	if int(currentHP) <= int(float64(maxHP)*1.0/2.0) {
		heal := int(float64(maxHP) * (1.0 / 4.0))
		spovb = spovb.Heal(heal)
	}

	return spovb
}

func (spovb SelfPointOfViewBattle) AfterContact() SelfPointOfViewBattle {
	if spovb.OpponentFighters[0].Item == "ゴツゴツメット" {
		damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 6.0)
		spovb = spovb.ToDamage(damage)
	}
	return spovb
}

//https://latest.pokewiki.net/%E3%83%90%E3%83%88%E3%83%AB%E4%B8%AD%E3%81%AE%E5%87%A6%E7%90%86%E3%81%AE%E9%A0%86%E7%95%AA
func (spovb SelfPointOfViewBattle) MoveUse(moveName MoveName, random *rand.Rand) (SelfPointOfViewBattle, error) {
	var err error
	if spovb.SelfFighters[0].IsFaint() {
		return spovb, nil
	}

	if moveName == STRUGGLE {
		spovb.SelfFighters[0].State.CurrentHP = 0
		return spovb, nil
	}

	moveData := MOVEDEX[moveName]

	if _, ok := spovb.SelfFighters[0].Moveset[moveName]; !ok {
		errMsg := fmt.Sprintf("%v は %v を繰り出そうとしたが、覚えていない", spovb.SelfFighters[0].Name, moveName)
		return SelfPointOfViewBattle{}, fmt.Errorf(errMsg)
	}

	if spovb.SelfFighters[0].Moveset[moveName] <= 0 {
		errMsg := fmt.Sprintf("%v を繰り出そうとしたが、powerPointが0以下", moveName)
		return SelfPointOfViewBattle{}, fmt.Errorf(errMsg)
	}

	choiceMoveName := spovb.SelfFighters[0].ChoiceMoveName
	isChoiceState := choiceMoveName != ""

	if isChoiceState && choiceMoveName != moveName {
		errMsg := fmt.Sprintf("%vを選択したが、拘り状態な為、%vしか選択できない", moveName, choiceMoveName)
		return SelfPointOfViewBattle{}, fmt.Errorf(errMsg)
	}

	isSleep := spovb.SelfFighters[0].StatusAilmentDetail.StatusAilment == SLEEP

	if isSleep {
		spovb.SelfFighters[0].StatusAilmentDetail.SleepRemainingTurn -= 1
	}

	if spovb.SelfFighters[0].StatusAilmentDetail.SleepRemainingTurn == 0 && isSleep {
		spovb.SelfFighters[0].StatusAilmentDetail.StatusAilment = ""
	}

	if spovb.SelfFighters[0].StatusAilmentDetail.StatusAilment == SLEEP {
		return spovb, nil
	}

	copyMoveset := spovb.SelfFighters[0].Moveset.Copy()
	copyMoveset[moveName] -= 1
	spovb.SelfFighters[0].Moveset = copyMoveset

	if spovb.OpponentFighters[0].IsFaint() {
		return spovb, nil
	}

	if spovb.SelfFighters[0].Item.IsChoice() {
		spovb.SelfFighters[0].ChoiceMoveName = moveName
	}

	accuracy := spovb.Accuracy(moveName, moveData)
	if accuracy != -1 {
		if !IsHit(moveData.Accuracy, random) {
			return spovb, nil
		}
	}

	if moveData.Category == STATUS {
		return spovb, nil
	}

	attackNum, err := omw.RandomInt(moveData.MinAttackNum, moveData.MaxAttackNum+1, random)

	if err != nil {
		return SelfPointOfViewBattle{}, err
	}

	for i := 0; i < attackNum; i++ {
		isCritical := spovb.IsCritical(moveName, random)
		finalDamage, err := NewFinalDamage(&spovb, moveName, isCritical, NewDamageR(random))

		if err != nil {
			return SelfPointOfViewBattle{}, err
		}

		if finalDamage == 0 {
			return spovb, nil
		}

		realDamage := NewRealDamage(&spovb, finalDamage)
		if spovb.OpponentFighters[0].IsFocusSashOk(int(finalDamage)) {
			spovb.OpponentFighters[0].Item = ""
		}

		opovb := spovb.SwitchPointOfView()
		opovb = opovb.ToDamage(int(realDamage))
		spovb = opovb.SwitchPointOfView()

		if moveData.Contact == "接触" {
			spovb = spovb.AfterContact()
		}

		if spovb.SelfFighters[0].IsFaint() || spovb.OpponentFighters[0].IsFaint() {
			break
		}
	}

	if spovb.SelfFighters[0].Item == "いのちのたま" {
		damage := int(float64(spovb.SelfFighters[0].State.MaxHP) * 1.0 / 10.0)
		spovb = spovb.ToDamage(damage)
	}
	return spovb, nil
}

func (spovb SelfPointOfViewBattle) Switch(pokeName PokeName) (SelfPointOfViewBattle, error) {
	index := spovb.SelfFighters.Index(pokeName)

	if index == -1 {
		errMsg := fmt.Sprintf("%vに交代しようとしたが、fightersの中に存在しない", pokeName)
		return SelfPointOfViewBattle{}, fmt.Errorf(errMsg)
	}

	if index == 0 {
		errMsg := fmt.Sprintf("%vに交代しようとしたが、既に場に出ている", pokeName)
		return SelfPointOfViewBattle{}, fmt.Errorf(errMsg)
	}

	if spovb.SelfFighters[index].IsFaint() {
		errMsg := fmt.Sprintf("%vに交代しようとしたが、瀕死状態", pokeName)
		return SelfPointOfViewBattle{}, fmt.Errorf(errMsg)
	}

	spovb.SelfFighters[0].Rank = Rank{}
	spovb.SelfFighters[0].StatusAilmentDetail.BadPoisonElapsedTurn = 0
	spovb.SelfFighters[0].ChoiceMoveName = ""
	spovb.SelfFighters[0].IsLeechSeed = false

	if spovb.SelfFighters[0].Ability == "しぜんかいふく" {
		spovb.SelfFighters[0].StatusAilmentDetail.StatusAilment = ""
	}

	tmpFighters := spovb.SelfFighters

	if index == 1 {
		spovb.SelfFighters[0] = tmpFighters[1]
		spovb.SelfFighters[1] = tmpFighters[0]
		spovb.SelfFighters[2] = tmpFighters[2]
	} else {
		spovb.SelfFighters[0] = tmpFighters[2]
		spovb.SelfFighters[1] = tmpFighters[1]
		spovb.SelfFighters[2] = tmpFighters[0]
	}

	if spovb.SelfField.IsStealthRockActive {
		stealthRockDamage := spovb.SelfFighters[0].StealthRockDamage()
		spovb = spovb.ToDamage(stealthRockDamage)
	}
	return spovb, nil
}

func (spovb SelfPointOfViewBattle) Action(battleCommand BattleCommand, random *rand.Rand) (SelfPointOfViewBattle, error) {
	if battleCommand.IsMoveName() || battleCommand == BattleCommand(STRUGGLE) {
		return spovb.MoveUse(MoveName(battleCommand), random)
	}

	if battleCommand.IsPokeName() {
		return spovb.Switch(PokeName(battleCommand))
	}

	errMsg := fmt.Sprintf("「%v」は、battleCommandとして不適", battleCommand)
	return SelfPointOfViewBattle{}, fmt.Errorf(errMsg)
}

func (spovb *SelfPointOfViewBattle) NewBattle(p1BattleCommand BattleCommand) Battle {
	return Battle{P1Fighters:spovb.SelfFighters, P2Fighters:spovb.OpponentFighters,
		P1Field:spovb.SelfField, P2Field:spovb.OpponentField, P1Command:p1BattleCommand}
}

type Battle struct {
	P1Fighters Fighters
	P2Fighters Fighters

	P1Field OneSideField
	P2Field OneSideField

	P1Command BattleCommand
}

func (battle Battle) ReversePlayer() (Battle, error) {
	if battle.P1Command != "" {
		return Battle{}, fmt.Errorf("battle.ReversePlayerを呼び出す時のbattle.P1Commandは、ゼロ値でなければならない")
	}
	return Battle{P1Fighters:battle.P2Fighters, P2Fighters:battle.P1Fighters,
		P1Field:battle.P2Field, P2Field:battle.P1Field, P1Command:battle.P1Command}, nil
}

func (battle *Battle) NewP1PointOfViewBattle() SelfPointOfViewBattle {
	return SelfPointOfViewBattle{
		SelfFighters:battle.P1Fighters, OpponentFighters:battle.P2Fighters,
		SelfField:battle.P1Field, OpponentField:battle.P2Field,
	}
}

func (battle *Battle) NewP2PointOfViewBattle() SelfPointOfViewBattle {
	return SelfPointOfViewBattle{
		SelfFighters:battle.P2Fighters, OpponentFighters:battle.P1Fighters,
		SelfField:battle.P2Field, OpponentField:battle.P1Field,
	}
}

func (battle1 *Battle) Equal(battle2 *Battle) bool {
	return battle1.P1Fighters.Equal(&battle2.P1Fighters) &&
		battle1.P2Fighters.Equal(&battle2.P2Fighters) &&
		battle1.P1Command == battle2.P1Command
}

func (battle *Battle) FinalSpeedWinner() Winner {
	p1PointOfViewBattle := battle.NewP1PointOfViewBattle()
	p2PointOfViewBattle := battle.NewP2PointOfViewBattle()

	p1FinalSpeed := NewFinalSpeed(&p1PointOfViewBattle)
	p2FinalSpeed := NewFinalSpeed(&p2PointOfViewBattle)

	if p1FinalSpeed > p2FinalSpeed {
		return WINNER_PLAYER1
	}

	if p1FinalSpeed < p2FinalSpeed {
		return WINNER_PLAYER2
	}
	return DRAW
}

func (battle *Battle) PriorityWinner(p1BattleCommand, p2BattleCommand BattleCommand) Winner {
	p1PriorityRank := p1BattleCommand.PriorityRank()
	p2PriorityRank := p2BattleCommand.PriorityRank()

	if p1PriorityRank > p2PriorityRank {
		return WINNER_PLAYER1
	}

	if p1PriorityRank < p2PriorityRank {
		return WINNER_PLAYER2
	}
	return DRAW
}

func (battle *Battle) IsP1FirstAction(p1BattleCommand, p2BattleCommand BattleCommand, random *rand.Rand) bool {
	priorityWinner := battle.PriorityWinner(p1BattleCommand, p2BattleCommand)

	if priorityWinner == WINNER_PLAYER1 {
		return true
	}

	if priorityWinner == WINNER_PLAYER2 {
		return false
	}

	finalSpeedWinner := battle.FinalSpeedWinner()

	if finalSpeedWinner == WINNER_PLAYER1 {
		return true
	}

	if finalSpeedWinner == WINNER_PLAYER2 {
		return false
	}

	return omw.RandomBool(random)
}

func (battle Battle) P1Action(battleCommand BattleCommand, random *rand.Rand) (Battle, error) {
	var err error
	p1PointOfViewBattle := battle.NewP1PointOfViewBattle()
	p1PointOfViewBattle, err = p1PointOfViewBattle.Action(battleCommand, random)
	return p1PointOfViewBattle.NewBattle(battle.P1Command), err
}

func (battle Battle) P2Action(battleCommand BattleCommand, random *rand.Rand) (Battle, error) {
	var err error
	p2PointOfViewBattle := battle.NewP2PointOfViewBattle()
	p2PointOfViewBattle, err = p2PointOfViewBattle.Action(battleCommand, random)
	p1PointOfViewBattle := p2PointOfViewBattle.SwitchPointOfView()
	return p1PointOfViewBattle.NewBattle(battle.P1Command), err
}

func (battle Battle) TurnEnd(random *rand.Rand) Battle {
	//https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86
	p1First := func(battle Battle, f func(SelfPointOfViewBattle) SelfPointOfViewBattle) Battle {
		p1PointOfViewBattle := battle.NewP1PointOfViewBattle()
		p1PointOfViewBattle = f(p1PointOfViewBattle)
		p2PointOfViewBattle := p1PointOfViewBattle.SwitchPointOfView()
		p2PointOfViewBattle = f(p2PointOfViewBattle)
		p1PointOfViewBattle = p2PointOfViewBattle.SwitchPointOfView()
		return p1PointOfViewBattle.NewBattle(battle.P1Command)
	}

	p2First := func(battle Battle, f func(SelfPointOfViewBattle) SelfPointOfViewBattle) Battle {
		p2PointOfViewBattle := battle.NewP2PointOfViewBattle()
		p2PointOfViewBattle = f(p2PointOfViewBattle)
		p1PointOfViewBattle := p2PointOfViewBattle.SwitchPointOfView()
		p1PointOfViewBattle = f(p1PointOfViewBattle)
		return p1PointOfViewBattle.NewBattle(battle.P1Command)
	}

	run := func(fs []func(SelfPointOfViewBattle) SelfPointOfViewBattle) Battle {
		finalSpeedWinner := battle.FinalSpeedWinner()

		for _, f := range fs {
			if finalSpeedWinner == WINNER_PLAYER1 {
				battle = p1First(battle, f)
			} else if finalSpeedWinner == WINNER_PLAYER2 {
				battle = p2First(battle, f)
			} else {
				if omw.RandomBool(random) {
					battle = p1First(battle, f)
				} else {
					battle = p2First(battle, f)
				}
			}
		}
		return battle
	}

	battle = run([]func(SelfPointOfViewBattle)SelfPointOfViewBattle{TurnEndLeftovers, TurnEndBlackSludge})
	battle = run([]func(SelfPointOfViewBattle)SelfPointOfViewBattle{TurnEndLeechSeed})
	battle = run([]func(SelfPointOfViewBattle)SelfPointOfViewBattle{TurnEndBadPoison})
	battle = run([]func(SelfPointOfViewBattle)SelfPointOfViewBattle{TurnEndBurn})
	return battle
}

func (battle Battle) Run(battleCommand BattleCommand, random *rand.Rand) (Battle, error) {
	var err error

	if battle.P1Fighters[0].IsFaint() {
		return battle.P1Action(battleCommand, random)
	}

	if battle.P2Fighters[0].IsFaint() {
		return battle.P2Action(battleCommand, random)
	}

	if battle.P1Command == BattleCommand("") {
		battle.P1Command = battleCommand
		return battle, nil
	}

	isP1FirstAction := battle.IsP1FirstAction(battle.P1Command, battleCommand, random)

	var isP1Actions []bool
	var actionOrderBattleCommands []BattleCommand

	if isP1FirstAction {
		isP1Actions = []bool{true, false}
		actionOrderBattleCommands = []BattleCommand{battle.P1Command, battleCommand}
	} else {
		isP1Actions = []bool{false, true}
		actionOrderBattleCommands = []BattleCommand{battleCommand, battle.P1Command}
	}

	for i, iBattleCommand := range actionOrderBattleCommands {
		if isP1Actions[i] {
			battle, err = battle.P1Action(iBattleCommand, random)
		} else {
			battle, err = battle.P2Action(iBattleCommand, random)
		}

		if err != nil {
			return Battle{}, err
		}
	}

	battle = battle.TurnEnd(random)
	battle.P1Command = BattleCommand("")
	return battle, err
}

func (battle *Battle) IsP1Phase() bool {
	if battle.P1Fighters[0].IsFaint() {
		return true
	}

	if battle.P2Fighters[0].IsFaint() {
		return false
	}

	return battle.P1Command == ""
}

func (battle *Battle) IsGameEnd() bool {
	return battle.P1Fighters.IsAllFaint() || battle.P2Fighters.IsAllFaint()
}

func (battle *Battle) Winner() (Winner, error) {
	if !battle.IsGameEnd() {
		return Winner{}, fmt.Errorf("ゲームが終了していない状態でWinnerは求める事は出来ない")
	}

	isP1AllFaint := battle.P1Fighters.IsAllFaint()
	isP2AllFaint := battle.P2Fighters.IsAllFaint()

	if isP1AllFaint && isP2AllFaint {
		return DRAW, nil
	}

	if isP1AllFaint {
		return WINNER_PLAYER1, nil
	}

	return WINNER_PLAYER2, nil
}
