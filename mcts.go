package bippa

import (
	"fmt"
	"github.com/sw965/crow"
	"math/rand"
)

type PolynomialUpperConfidenceBound struct {
	P           float64
	AccumReward float64
	Trial       int
}

func (pucb *PolynomialUpperConfidenceBound) AverageReward() float64 {
	return float64(pucb.AccumReward) / float64(pucb.Trial+1)
}

func (pucb *PolynomialUpperConfidenceBound) Get(totalTrial int, X float64) float64 {
	averageReward := pucb.AverageReward()
	return crow.PolynomialUpperConfidenceBound(averageReward, pucb.P, totalTrial, pucb.Trial, X)
}

type PolynomialUpperConfidenceBounds []*PolynomialUpperConfidenceBound

func (pucbs PolynomialUpperConfidenceBounds) TotalTrial() int {
	result := 0
	for _, pucb := range pucbs {
		result += pucb.Trial
	}
	return result
}

func (pucbs PolynomialUpperConfidenceBounds) Max(X float64) float64 {
	totalTrial := pucbs.TotalTrial()
	result := pucbs[0].Get(totalTrial, X)
	for _, pucb := range pucbs[1:] {
		pucbY := pucb.Get(totalTrial, X)
		if pucbY > result {
			result = pucbY
		}
	}
	return result
}

func (pucbs PolynomialUpperConfidenceBounds) MaxPUCBsAndIndices(X float64) (PolynomialUpperConfidenceBounds, []int) {
	length := len(pucbs)
	totalTrial := pucbs.TotalTrial()
	max := pucbs.Max(X)
	maxPUCBs := make(PolynomialUpperConfidenceBounds, 0, length)
	indices := make([]int, 0, length)

	for i, pucb := range pucbs {
		pucbY := pucb.Get(totalTrial, X)
		if pucbY == max {
			maxPUCBs = append(maxPUCBs, pucb)
			indices = append(indices, i)
		}
	}
	return maxPUCBs, indices
}

type BattlePolicy func(*Battle) map[ActionCmd]float64

func NoBattlePolicy(battle *Battle) map[ActionCmd]float64 {
	legalActionCmds := battle.P1Fighters.LegalActionCmds()
	result := map[ActionCmd]float64{}
	floatLegalActionNum := float64(len(legalActionCmds))
	for _, actionCmd := range legalActionCmds {
		result[actionCmd] = 1.0 / floatLegalActionNum
	}
	return result
}

type ActionCmdPUCBs map[ActionCmd]*PolynomialUpperConfidenceBound

func (actionCmdPUCBs ActionCmdPUCBs) Keys() ActionCmds {
	result := make(ActionCmds, 0, len(actionCmdPUCBs))
	for actionCmd, _ := range actionCmdPUCBs {
		result = append(result, actionCmd)
	}
	return result
}

func (actionCmdPUCBs ActionCmdPUCBs) TotalTrial() int {
	result := 0
	for _, pucb := range actionCmdPUCBs {
		result += pucb.Trial
	}
	return result
}

func (actionCmdPUCBs ActionCmdPUCBs) Max(X float64) float64 {
	totalTrial := actionCmdPUCBs.TotalTrial()
	actionCmds := actionCmdPUCBs.Keys()
	result := actionCmdPUCBs[actionCmds[0]].Get(totalTrial, X)

	for _, actionCmd := range actionCmds[1:] {
		pucbY := actionCmdPUCBs[actionCmd].Get(totalTrial, X)
		if pucbY > result {
			result = pucbY
		}
	}
	return result
}

func (actionCmdPUCBs ActionCmdPUCBs) MaxActionCmds(X float64) ActionCmds {
	max := actionCmdPUCBs.Max(X)
	totalTrial := actionCmdPUCBs.TotalTrial()
	result := make(ActionCmds, 0)

	for actionCmd, pucb := range actionCmdPUCBs {
		pucbY := pucb.Get(totalTrial, X)
		if pucbY == max {
			result = append(result, actionCmd)
		}
	}
	return result
}

func (actionCmdPUCBs ActionCmdPUCBs) MaxTrial() int {
	actionCmds := actionCmdPUCBs.Keys()
	result := actionCmdPUCBs[actionCmds[0]].Trial

	for _, actionCmd := range actionCmds[1:] {
		trial := actionCmdPUCBs[actionCmd].Trial
		if trial > result {
			result = trial
		}
	}
	return result
}

func (actionCmdPUCBs ActionCmdPUCBs) MaxTrialActionCmds() ActionCmds {
	maxTrial := actionCmdPUCBs.MaxTrial()
	result := make(ActionCmds, 0)
	for actionCmd, pucb := range actionCmdPUCBs {
		trial := pucb.Trial
		if trial == maxTrial {
			result = append(result, actionCmd)
		}
	}
	return result
}

type BattleNode struct {
	Battle          *Battle
	LegalActionCmds ActionCmds
	ActionCmdPUCBs  ActionCmdPUCBs
	NextBattleNodes       BattleNodes
	IsP1            bool
	SelectCount     int
}

func NewBattleNodePointer(battle *Battle, battlePolicy BattlePolicy) *BattleNode {
	isP1Phase := battle.IsP1Phase()
	var fighters Fighters

	if isP1Phase {
		fighters = battle.P1Fighters
	} else {
		fighters = battle.P2Fighters
	}

	legalActionCmds := fighters.LegalActionCmds()
	actionCmdPUCBs := ActionCmdPUCBs{}
	var battlePolicyY map[ActionCmd]float64

	if isP1Phase {
		battlePolicyY = battlePolicy(battle)
	} else {
		reverseBattle := battle.Reverse()
		battlePolicyY = battlePolicy(&reverseBattle)
	}

	for _, actionCmd := range legalActionCmds {
		actionCmdPUCBs[actionCmd] = &PolynomialUpperConfidenceBound{P: battlePolicyY[actionCmd]}
	}

	return &BattleNode{Battle: battle, LegalActionCmds: legalActionCmds, ActionCmdPUCBs: actionCmdPUCBs,
		IsP1: isP1Phase, SelectCount: 0}
}

func (battleNode *BattleNode) SelectAndExpansion(battle Battle, allBattleNodes BattleNodes, battlePolicy BattlePolicy, X float64, random *rand.Rand) (Battle, BattleNodes, Selects, error) {
	selects := make(Selects, 0, 64)
	var err error

	for {
		maxPUCBActionCmds := battleNode.ActionCmdPUCBs.MaxActionCmds(X)
		selectActionCmd := maxPUCBActionCmds.RandomChoice(random)
		selects = append(selects, Select{BattleNode: battleNode, ActionCmd: selectActionCmd})
		battleNode.SelectCount += 1

		battle, err = battle.Push(selectActionCmd, random)
		if err != nil {
			return Battle{}, BattleNodes{}, Selects{}, err
		}

		if battle.IsGameEnd() {
			break
		}

		//NextBattleNodesの中に、同じ局面のbattleが存在するならば、それを次のbattleNodeとする
		//NextBattleNodesの中に、同じ局面のbattleが存在しないなら、allBattleNodesの中から同じ局面のbattleが存在しないかを調べる。
		//allBattleNodesの中に、同じ局面のbattleが存在するならば、次回から高速に探索出来るように、NextBattleNodesに追加して、次のbattleNodeとする。
		//NextBattleNodesにもallBattleNodesにも同じ局面のbattleが存在しないなら、新しくbattleNodeを作り、
		//NextBattleNodesと、allBattleNodesに追加し、新しく作ったbattleNodeを次のbattleNodeとし、select処理を終了する。

		nextBattleNode, err := battleNode.NextBattleNodes.Find(&battle)
		if err != nil {
			nextBattleNode, err = allBattleNodes.Find(&battle)
			if err == nil {
				battleNode.NextBattleNodes = append(battleNode.NextBattleNodes, nextBattleNode)
			} else {
				nextBattleNode = NewBattleNodePointer(&battle, battlePolicy)
				allBattleNodes = append(allBattleNodes, nextBattleNode)
				battleNode.NextBattleNodes = append(battleNode.NextBattleNodes, nextBattleNode)
				break
			}
		}

		if nextBattleNode.SelectCount == 1 {
			break
		}
		battleNode = nextBattleNode
	}
	return battle, allBattleNodes, selects, nil
}

func (battleNode *BattleNode) AverageReward() float64 {
	accumReward := 0.0
	for _, pucb := range battleNode.ActionCmdPUCBs {
		accumReward += pucb.AverageReward()
	}
	return float64(accumReward) / float64(len(battleNode.ActionCmdPUCBs))
}

type BattleNodes []*BattleNode

func (battleNodes BattleNodes) Find(battle *Battle) (*BattleNode, error) {
	for _, battleNode := range battleNodes {
		if battleNode.Battle.Equal(battle) {
			return battleNode, nil
		}
	}
	return &BattleNode{}, fmt.Errorf("battleが一致しているbattleNodeが見つからなかった")
}

type Select struct {
	BattleNode      *BattleNode
	ActionCmd ActionCmd
}

type Selects []Select

func (selects Selects) Backward(battleEvalY float64, battleEval *BattleEval) {
	for _, select_ := range selects {
		battleNode := select_.BattleNode
		actionCmd := select_.ActionCmd

		if battleNode.IsP1 {
			battleNode.ActionCmdPUCBs[actionCmd].AccumReward += battleEvalY
		} else {
			battleNode.ActionCmdPUCBs[actionCmd].AccumReward += battleEval.Reverse(battleEvalY)
		}
		battleNode.ActionCmdPUCBs[actionCmd].Trial += 1
		battleNode.SelectCount = 0
	}
}

func RunMCTS(rootBattle Battle, simuNum int, X float64, battlePolicy BattlePolicy, battleEval *BattleEval, random *rand.Rand) (BattleNodes, error) {
	rootBattleNode := NewBattleNodePointer(&rootBattle, battlePolicy)
	allBattleNodes := BattleNodes{rootBattleNode}
	battle := rootBattle

	var selects Selects
	var battleEvalY float64
	var err error

	for i := 0; i < simuNum; i++ {
		battle, allBattleNodes, selects, err = rootBattleNode.SelectAndExpansion(battle, allBattleNodes, battlePolicy, X, random)
		if err != nil {
			return BattleNodes{}, err
		}

		battleEvalY, err = battleEval.Func(&battle)

		if err != nil {
			return BattleNodes{}, err
		}

		selects.Backward(battleEvalY, battleEval)
		battle = rootBattle
	}
	return allBattleNodes, nil
}

func NewMCTSTrainer(simuNum int, X float64, battlePolicy BattlePolicy, battleEval *BattleEval, random *rand.Rand) Trainer {
	result := func(battle *Battle) (ActionCmd, error) {
		legalActionCmds := battle.P1Fighters.LegalActionCmds()

		if len(legalActionCmds) == 1 {
			return legalActionCmds[0], nil
		}

		allBattleNodes, err := RunMCTS(*battle, simuNum, X, battlePolicy, battleEval, random)

		if err != nil {
			return "", err
		}

		return allBattleNodes[0].ActionCmdPUCBs.MaxTrialActionCmds().RandomChoice(random), nil
	}
	return result
}

type TeamNode struct {
	Team Team
	LegalPokeNames PokeNames
	LegalAbilities Abilities
	LegalItems Items
	LegalMoveNames MoveNames
	LegalNatures Natures

	Policies []float64
}