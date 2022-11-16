package bippa

import (
	"encoding/json"
	"fmt"
	"github.com/sw965/omw"
	"io/ioutil"
	"math/rand"
)

func IsHit(n int, random *rand.Rand) bool {
	return random.Intn(100) < n
}

type CriticalRank int

const (
	FIGHTERS_LENGTH = 3
)

type Fighters [FIGHTERS_LENGTH]Pokemon

func (fighters1 *Fighters) Equal(fighters2 *Fighters) bool {
	for i, pokemon := range fighters1 {
		if !pokemon.Equal(&fighters2[i]) {
			return false
		}
	}
	return true
}

func (fighters *Fighters) Index(pokeName PokeName) int {
	for i, pokemon := range fighters {
		if pokemon.Name == pokeName {
			return i
		}
	}
	return -1
}

func (fighters *Fighters) PokeNames() PokeNames {
	result := make(PokeNames, FIGHTERS_LENGTH)
	for i, pokemon := range fighters {
		result[i] = pokemon.Name
	}
	return result
}

func (fighters *Fighters) IsAllFaint() bool {
	for _, pokemon := range fighters {
		if !pokemon.IsFaint() {
			return false
		}
	}
	return true
}

func (fighters *Fighters) FaintList() []bool {
	result := make([]bool, FIGHTERS_LENGTH)
	for i, pokemon := range fighters {
		result[i] = pokemon.IsFaint()
	}
	return result
}

func (fighters *Fighters) LegalActionCmdMoveNames() MoveNames {
	if fighters[0].IsFaint() {
		return MoveNames{}
	}

	tmpResult := MoveNames{}
	for moveName, powerPoint := range fighters[0].Moveset {
		if powerPoint.Current > 0 {
			tmpResult = append(tmpResult, moveName)
		}
	}

	var result MoveNames

	if fighters[0].ChoiceMoveName != "" {
		if tmpResult.In(fighters[0].ChoiceMoveName) {
			result = MoveNames{fighters[0].ChoiceMoveName}
		}
	} else if fighters[0].Item == "とつげきチョッキ" {
		for _, moveName := range tmpResult {
			if MOVEDEX[moveName].Category != STATUS {
				result = append(result, moveName)
			}
		}
	} else {
		result = tmpResult
	}

	if len(result) == 0 {
		return MoveNames{STRUGGLE}
	}
	return result
}

func (fighters *Fighters) LegalActionCmdPokeNames() []PokeName {
	result := make([]PokeName, 0)
	for _, pokemon := range fighters[1:] {
		if !pokemon.IsFaint() {
			result = append(result, pokemon.Name)
		}
	}
	return result
}

func (fighters *Fighters) LegalActionCmds() ActionCmds {
	legalActionCmdMoveNames := fighters.LegalActionCmdMoveNames()
	legalActionCmdPokeNames := fighters.LegalActionCmdPokeNames()
	result := make(ActionCmds, 0, len(legalActionCmdMoveNames)+len(legalActionCmdPokeNames))

	for _, moveName := range legalActionCmdMoveNames {
		result = append(result, ActionCmd(moveName))
	}

	for _, pokeName := range legalActionCmdPokeNames {
		result = append(result, ActionCmd(pokeName))
	}
	return result
}

func (fighters *Fighters) Save(filePath string) error {
	file, err := json.MarshalIndent(fighters, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, file, 0644)
}

func ReadFighters(filePath string) (Fighters, error) {
	result := Fighters{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Fighters{}, err
	}
	err = json.Unmarshal(file, &result)
	return result, err
}

//https://wiki.xn--rckteqa2e.com/wiki/%E3%81%99%E3%81%B0%E3%82%84%E3%81%95#.E8.A9.B3.E7.B4.B0.E3.81.AA.E4.BB.95.E6.A7.98
type SpeedBonus int

const (
	INIT_SPEED_BONUS = SpeedBonus(4096)
)

func NewSpeedBonus(battle *Battle) SpeedBonus {
	result := int(INIT_SPEED_BONUS)
	if battle.P1Fighters[0].Item == "こだわりスカーフ" {
		result = FiveOrMoreRounding(float64(result) * 6144.0 / 4096.0)
	}
	return SpeedBonus(result)
}

type FinalSpeed float64

func NewFinalSpeed(battle *Battle) FinalSpeed {
	speed := battle.P1Fighters[0].Speed
	rankBonus := battle.P1Fighters[0].Rank.Speed.ToBonus()
	speedBonus := NewSpeedBonus(battle)

	result := FiveOrMoreRounding(float64(speed) * float64(rankBonus))
	result = FiveOverRounding(float64(result) * float64(speedBonus) / 4096.0)
	return FinalSpeed(result)
}

type ActionCmd string

func (actionCmd ActionCmd) IsMoveName() bool {
	_, ok := MOVEDEX[MoveName(actionCmd)]
	return ok
}

func (actionCmd ActionCmd) IsPokeName() bool {
	_, ok := POKEDEX[PokeName(actionCmd)]
	return ok
}

func (actionCmd ActionCmd) PriorityRank() int {
	if actionCmd == ActionCmd(STRUGGLE) {
		return 0
	} else if actionCmd.IsMoveName() {
		return MOVEDEX[MoveName(actionCmd)].PriorityRank
	} else {
		return 999
	}
}

type ActionCmds []ActionCmd

func (actionCmds ActionCmds) RandomChoice(random *rand.Rand) ActionCmd {
	index := random.Intn(len(actionCmds))
	return actionCmds[index]
}

type Battle struct {
	P1Fighters  Fighters
	P2Fighters  Fighters
	P1ActionCmd ActionCmd
}

func (battle *Battle) Reverse() Battle {
	return Battle{P1Fighters: battle.P2Fighters, P2Fighters: battle.P1Fighters, P1ActionCmd: battle.P1ActionCmd}
}

func (battle *Battle) Accuracy(moveName MoveName) int {
	return MOVEDEX[moveName].Accuracy
}

func (battle *Battle) CriticalN(moveName MoveName) int {
	criticalRank := MOVEDEX[moveName].CriticalRank
	return CRITICAL_N[criticalRank]
}

func (battle *Battle) IsCritical(moveName MoveName, random *rand.Rand) bool {
	//https://wiki.xn--rckteqa2e.com/wiki/%E6%80%A5%E6%89%80
	criticalN := battle.CriticalN(moveName)
	return random.Intn(criticalN) == 0
}

func (battle Battle) Damage(damage int) Battle {
	if battle.P1Fighters[0].IsFaint() {
		return battle
	}

	currentHP := battle.P1Fighters[0].CurrentHP

	if damage > currentHP {
		damage = currentHP
	}

	battle.P1Fighters[0].CurrentHP -= damage
	return battle.SitrusBerryHeal()
}

func (battle Battle) Heal(heal int) Battle {
	if battle.P1Fighters[0].IsFaint() {
		return battle
	}

	currentDamage := battle.P1Fighters[0].CurrentDamage()

	if currentDamage < heal {
		heal = currentDamage
	}

	battle.P1Fighters[0].CurrentHP += heal
	return battle
}

func (battle Battle) RankFluctuation(fluctuationRank *Rank) Battle {
	if battle.P1Fighters[0].IsFaint() {
		return battle
	}

	rank := battle.P1Fighters[0].Rank
	newRank := rank.Add(fluctuationRank)

	if battle.P1Fighters[0].Item == "しろいハーブ" && newRank.InDown() {
		battle.P1Fighters[0].Item = EMPTY_ITEM
		newRank = newRank.ResetDown()
	}

	battle.P1Fighters[0].Rank = newRank.Regulate()
	return battle
}

func (battle Battle) SitrusBerryHeal() Battle {
	if battle.P1Fighters[0].Item != "オボンのみ" {
		return battle
	}

	if battle.P1Fighters[0].IsFaint() {
		return battle
	}

	maxHP := battle.P1Fighters[0].MaxHP
	currentHP := battle.P1Fighters[0].CurrentHP

	if int(currentHP) <= int(float64(maxHP)*1.0/2.0) {
		battle.P1Fighters[0].Item = EMPTY_ITEM
		heal := int(float64(maxHP) * (1.0 / 4.0))
		battle = battle.Heal(heal)
	}
	return battle
}

func (battle Battle) AfterContact() Battle {
	if battle.P2Fighters[0].Ability == "てつのトゲ" {
		damage := int(float64(battle.P1Fighters[0].MaxHP) * 1.0 / 8.0)
		battle = battle.Damage(damage)
	}

	if battle.P2Fighters[0].Item == "ゴツゴツメット" {
		damage := int(float64(battle.P1Fighters[0].MaxHP) * 1.0 / 6.0)
		battle = battle.Damage(damage)
	}
	return battle
}

//https://latest.pokewiki.net/%E3%83%90%E3%83%88%E3%83%AB%E4%B8%AD%E3%81%AE%E5%87%A6%E7%90%86%E3%81%AE%E9%A0%86%E7%95%AA
func (battle Battle) MoveUse(moveName MoveName, random *rand.Rand) (Battle, error) {
	if battle.P1Fighters[0].IsFaint() {
		return battle, nil
	}

	if moveName == STRUGGLE {
		if battle.P2Fighters[0].IsFaint() {
			return battle, nil
		}
		battle.P1Fighters[0].CurrentHP = 0
		return battle, nil
	}

	moveData := MOVEDEX[moveName]

	if _, ok := battle.P1Fighters[0].Moveset[moveName]; !ok {
		errMsg := fmt.Sprintf("%v は %v を繰り出そうとしたが、覚えていない", battle.P1Fighters[0].Name, moveName)
		return Battle{}, fmt.Errorf(errMsg)
	}

	if battle.P1Fighters[0].Moveset[moveName].Current <= 0 {
		errMsg := fmt.Sprintf("%v は %v を繰り出そうとしたが、powerPointが0以下", battle.P1Fighters[0].Name, moveName)
		return Battle{}, fmt.Errorf(errMsg)
	}

	choiceMoveName := battle.P1Fighters[0].ChoiceMoveName
	isChoiceState := choiceMoveName != ""

	if isChoiceState && choiceMoveName != moveName {
		errMsg := fmt.Sprintf("%vを選択したが、拘り状態な為、%vしか選択できない", moveName, choiceMoveName)
		return Battle{}, fmt.Errorf(errMsg)
	}

	copyMoveset := battle.P1Fighters[0].Moveset.Copy()
	copyMoveset[moveName].Current -= 1
	battle.P1Fighters[0].Moveset = copyMoveset

	if battle.P2Fighters[0].IsFaint() {
		return battle, nil
	}

	if battle.P1Fighters[0].Item.IsChoice() {
		battle.P1Fighters[0].ChoiceMoveName = moveName
	}

	accuracy := moveData.Accuracy
	if accuracy != -1 {
		if !IsHit(moveData.Accuracy, random) {
			return battle, nil
		}
	}

	if moveData.Category == STATUS {
		statusMove, ok := STATUS_MOVES[moveName]
		if !ok {
			return battle, nil
		}
		return statusMove(battle, random), nil
	}

	attackNum, err := omw.RandomInt(moveData.MinAttackNum, moveData.MaxAttackNum+1, random)

	if err != nil {
		return Battle{}, err
	}

	for i := 0; i < attackNum; i++ {
		isCritical := battle.IsCritical(moveName, random)
		finalDamage, err := battle.NewFinalDamage(moveName, isCritical, NewRandomDamageBonus(random))

		if err != nil {
			return Battle{}, err
		}

		if finalDamage == 0 {
			return battle, nil
		}

		battle = battle.Reverse()
		battle = battle.Damage(int(finalDamage))
		battle = battle.Reverse()

		if moveData.Contact == "接触" {
			battle = battle.AfterContact()
		}

		if battle.P1Fighters[0].IsFaint() || battle.P2Fighters[0].IsFaint() {
			break
		}
	}

	if battle.P1Fighters[0].Item == "いのちのたま" {
		damage := int(float64(battle.P1Fighters[0].MaxHP) * 1.0 / 10.0)
		battle = battle.Damage(damage)
	}
	return battle, nil
}

func (battle Battle) Switch(pokeName PokeName) (Battle, error) {
	index := battle.P1Fighters.Index(pokeName)

	if index == -1 {
		errMsg := fmt.Sprintf("%vに交代しようとしたが、fightersの中に存在しない", pokeName)
		return Battle{}, fmt.Errorf(errMsg)
	}

	if index == 0 {
		errMsg := fmt.Sprintf("%vに交代しようとしたが、既に場に出ている", pokeName)
		return Battle{}, fmt.Errorf(errMsg)
	}

	if battle.P1Fighters[index].IsFaint() {
		errMsg := fmt.Sprintf("%vに交代しようとしたが、瀕死状態", pokeName)
		return Battle{}, fmt.Errorf(errMsg)
	}

	battle.P1Fighters[0].Rank = Rank{}
	battle.P1Fighters[0].BadPoisonElapsedTurn = 0
	battle.P1Fighters[0].ChoiceMoveName = ""
	battle.P1Fighters[0].IsLeechSeed = false

	tmpFighters := battle.P1Fighters

	if index == 1 {
		battle.P1Fighters[0] = tmpFighters[1]
		battle.P1Fighters[1] = tmpFighters[0]
		battle.P1Fighters[2] = tmpFighters[2]
	} else {
		battle.P1Fighters[0] = tmpFighters[2]
		battle.P1Fighters[1] = tmpFighters[1]
		battle.P1Fighters[2] = tmpFighters[0]
	}

	return battle, nil
}

func (battle Battle) P1Action(actionCmd ActionCmd, random *rand.Rand) (Battle, error) {
	if actionCmd.IsMoveName() || actionCmd == ActionCmd(STRUGGLE) {
		return battle.MoveUse(MoveName(actionCmd), random)
	}

	if actionCmd.IsPokeName() {
		return battle.Switch(PokeName(actionCmd))
	}

	errMsg := fmt.Sprintf("「%v」は、actionCmdとして不適", actionCmd)
	return Battle{}, fmt.Errorf(errMsg)
}

func (battle Battle) P2Action(actionCmd ActionCmd, random *rand.Rand) (Battle, error) {
	var err error
	battle = battle.Reverse()

	if actionCmd.IsMoveName() || actionCmd == ActionCmd(STRUGGLE) {
		battle, err = battle.MoveUse(MoveName(actionCmd), random)
		return battle.Reverse(), err
	}

	if actionCmd.IsPokeName() {
		battle, err = battle.Switch(PokeName(actionCmd))
		return battle.Reverse(), err
	}

	errMsg := fmt.Sprintf("「%v」は、actionCmdとして不適", actionCmd)
	return Battle{}, fmt.Errorf(errMsg)
}

func (battle1 *Battle) Equal(battle2 *Battle) bool {
	return battle1.P1Fighters.Equal(&battle2.P1Fighters) && battle1.P2Fighters.Equal(&battle2.P2Fighters) && battle1.P1ActionCmd == battle2.P1ActionCmd
}

func (battle *Battle) FinalSpeedWinner() Winner {
	p1FinalSpeed := NewFinalSpeed(battle)
	reverseBattle := battle.Reverse()
	p2FinalSpeed := NewFinalSpeed(&reverseBattle)

	if p1FinalSpeed > p2FinalSpeed {
		return WINNER_P1
	}

	if p1FinalSpeed < p2FinalSpeed {
		return WINNER_P2
	}
	return DRAW
}

func (battle *Battle) ActionCmdPriorityRankWinner(p1ActionCmd, p2ActionCmd ActionCmd) Winner {
	p1PriorityRank := p1ActionCmd.PriorityRank()
	p2PriorityRank := p2ActionCmd.PriorityRank()

	if p1PriorityRank > p2PriorityRank {
		return WINNER_P1
	}

	if p1PriorityRank < p2PriorityRank {
		return WINNER_P2
	}
	return DRAW
}

func (battle *Battle) ActionPriorityWinner(p1ActionCmd, p2ActionCmd ActionCmd, random *rand.Rand) Winner {
	actionCmdPriorityRankWinner := battle.ActionCmdPriorityRankWinner(p1ActionCmd, p2ActionCmd)

	if actionCmdPriorityRankWinner == WINNER_P1 {
		return WINNER_P1
	}

	if actionCmdPriorityRankWinner == WINNER_P2 {
		return WINNER_P2
	}

	finalSpeedWinner := battle.FinalSpeedWinner()

	if finalSpeedWinner == WINNER_P1 {
		return WINNER_P1
	}

	if finalSpeedWinner == WINNER_P2 {
		return WINNER_P2
	}

	if omw.RandomBool(random) {
		return WINNER_P1
	} else {
		return WINNER_P2
	}
}

func (battle Battle) TurnEnd(random *rand.Rand) Battle {
	//https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#5..E3.82.BF.E3.83.BC.E3.83.B3.E7.B5.82.E4.BA.86.E6.99.82.E3.81.AE.E5.87.A6.E7.90.86
	p1First := func(battle Battle, f func(Battle) Battle) Battle {
		battle = f(battle)
		battle = battle.Reverse()
		battle = f(battle)
		battle = battle.Reverse()
		return battle
	}

	p2First := func(battle Battle, f func(Battle) Battle) Battle {
		battle = battle.Reverse()
		battle = f(battle)
		battle = battle.Reverse()
		battle = f(battle)
		return battle
	}

	run := func(fs []func(Battle) Battle) Battle {
		finalSpeedWinner := battle.FinalSpeedWinner()

		for _, f := range fs {
			if finalSpeedWinner == WINNER_P1 {
				battle = p1First(battle, f)
			} else if finalSpeedWinner == WINNER_P2 {
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

	battle = run([]func(Battle) Battle{TurnEndLeftovers, TurnEndBlackSludge})
	battle = run([]func(Battle) Battle{TurnEndLeechSeed})
	battle = run([]func(Battle) Battle{TurnEndBadPoison})
	return battle
}

func (battle *Battle) IsP1Phase() bool {
	if battle.P1Fighters[0].IsFaint() {
		return true
	}

	if battle.P2Fighters[0].IsFaint() {
		return false
	}

	return battle.P1ActionCmd == ""
}

func (battle Battle) Push(actionCmd ActionCmd, random *rand.Rand) (Battle, error) {
	if battle.P1Fighters[0].IsFaint() {
		return battle.P1Action(actionCmd, random)
	}

	if battle.P2Fighters[0].IsFaint() {
		return battle.P2Action(actionCmd, random)
	}

	if battle.P1ActionCmd == "" {
		battle.P1ActionCmd = actionCmd
		return battle, nil
	}

	p2ActionCmd := actionCmd
	actionPriorityWinner := battle.ActionPriorityWinner(battle.P1ActionCmd, p2ActionCmd, random)

	var isP1Action []bool

	if actionPriorityWinner == WINNER_P1 {
		isP1Action = []bool{true, false}
	} else {
		isP1Action = []bool{false, true}
	}

	var err error
	for _, isP1 := range isP1Action {
		if isP1 {
			battle, err = battle.P1Action(battle.P1ActionCmd, random)
		} else {
			battle, err = battle.P2Action(p2ActionCmd, random)
		}

		if err != nil {
			return Battle{}, err
		}
	}

	if battle.IsGameEnd() {
		return battle, nil
	}

	battle = battle.TurnEnd(random)
	battle.P1ActionCmd = ""
	return battle, err
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
	} else if isP1AllFaint {
		return WINNER_P2, nil
	} else {
		return WINNER_P1, nil
	}
}

func (battle *Battle) NewFinalDamage(moveName MoveName, isCritical bool, randomDamageBonus RandomDamageBonus) (FinalDamage, error) {
	attackPokemon := battle.P1Fighters[0]
	defensePokemon := battle.P2Fighters[0]
	return NewFinalDamage(&attackPokemon, &defensePokemon, moveName, isCritical, randomDamageBonus)
}

func (battle *Battle) AttackDamageProbabilityDistribution(moveName MoveName) (DamageProbabilityDistribution, error) {
	accuracy := battle.Accuracy(moveName)
	criticalN := battle.CriticalN(moveName)
	return NewAttackDamageProbabilityDistribution(&battle.P1Fighters[0], &battle.P2Fighters[0], moveName, accuracy, criticalN)
}

type Winner struct {
	IsP1 bool
	IsP2 bool
}

var (
	WINNER_P1 = Winner{IsP1: true, IsP2: false}
	WINNER_P2 = Winner{IsP1: false, IsP2: true}
	DRAW      = Winner{IsP1: false, IsP2: false}
)

var WINNER_TO_REWARD = map[Winner]float64{WINNER_P1: 1.0, WINNER_P2: 0.0, DRAW: 0.5}
