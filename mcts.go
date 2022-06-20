package bippa

import (
	"fmt"
	"math/rand"
	"math"
)

type Eval struct {
	Func    func(*Battle) (float64, error)
	Reverse func(float64) float64
}

func NewPlayoutEval(trainer Trainer, random *rand.Rand) Eval {
	evalFunc := func(battle *Battle) (float64, error) {
		var winner Winner
		var err error

		if battle.IsGameEnd() {
			winner, err = battle.Winner()
		} else {
			battleV := *battle
			winner, err = trainer.Playout(trainer, battleV, random)
			if err != nil {
				return 0.0, err
			}
		}
		return WINNER_TO_REWARD[winner], err
	}

	reverse := func(evalY float64) float64 {
		//自分が勝ち(1)なら相手は負け(0)
		//自分が負け(0)なら相手は勝ち(1)
		//引き分けなら互いに(0.5)
		return 1 - evalY
	}
	return Eval{Func: evalFunc, Reverse: reverse}
}

type Policy func(*Battle) map[ActionCmd]float64

func NoPolicy(battle *Battle) map[ActionCmd]float64 {
	legalActionCmds := battle.P1Fighters.LegalActionCmds()
	result := map[ActionCmd]float64{}
	floatLegalActionNum := float64(len(legalActionCmds))
	for _, actionCmd := range legalActionCmds {
		result[actionCmd] = 1.0 / floatLegalActionNum
	}
	return result
}

type PolynomialUpperConfidenceBound struct {
	P float64
	AccumReward float64
	Trial       int
}

func (pucb *PolynomialUpperConfidenceBound) AverageReward() float64 {
	return float64(pucb.AccumReward) / float64(pucb.Trial + 1)
}

func (pucb *PolynomialUpperConfidenceBound) Get(totalTrial int, X float64) float64 {
	floatTotalTrial := float64(totalTrial)
	floatTrial := float64(pucb.Trial + 1)
	averageReward := pucb.AverageReward()
	return averageReward + (X * pucb.P * math.Sqrt(floatTotalTrial) / floatTrial)
	//return averageReward + (X * pucb.P * math.Sqrt(2*math.Log(floatTotalTrial)/floatTrial))
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

type Node struct {
	Battle          *Battle
	LegalActionCmds ActionCmds
	ActionCmdPUCBs  ActionCmdPUCBs
	NextNodes       Nodes
	IsP1        bool
	SelectCount     int
}

func NewNodePointer(battle *Battle, policy Policy) *Node {
	isP1Phase := battle.IsP1Phase()
	var fighters Fighters

	if isP1Phase {
		fighters = battle.P1Fighters
	} else {
		fighters = battle.P2Fighters
	}

	legalActionCmds := fighters.LegalActionCmds()
	actionCmdPUCBs := ActionCmdPUCBs{}
	var policyY map[ActionCmd]float64

	if isP1Phase {
		policyY = policy(battle)
	} else {
		reverseBattle := battle.Reverse()
		policyY = policy(&reverseBattle)
	}

	for _, actionCmd := range legalActionCmds {
		actionCmdPUCBs[actionCmd] = &PolynomialUpperConfidenceBound{P:policyY[actionCmd]}
	}

	return &Node{Battle: battle, LegalActionCmds: legalActionCmds, ActionCmdPUCBs: actionCmdPUCBs,
		IsP1: isP1Phase, SelectCount: 0}
}

func (node *Node) SelectAndExpansion(battle Battle, allNodes Nodes, policy Policy, X float64, random *rand.Rand) (Battle, Nodes, Selects, error) {
	selects := Selects{}
	var err error

	for {
		maxPUCBActionCmds := node.ActionCmdPUCBs.MaxActionCmds(X)
		selectActionCmd := maxPUCBActionCmds.RandomChoice(random)
		selects = append(selects, Select{Node: node, ActionCmd: selectActionCmd})
		node.SelectCount += 1

		battle, err = battle.Push(selectActionCmd, random)
		if err != nil {
			return Battle{}, Nodes{}, Selects{}, err
		}

		if battle.IsGameEnd() {
			break
		}

		//NextNodesの中に、同じ局面のbattleが存在するならば、それを次のnodeとする
		//NextNodesの中に、同じ局面のbattleが存在しないなら、allNodesの中から同じ局面のbattleが存在しないかを調べる。
		//allNodesの中に、同じ局面のbattleが存在するならば、次回から高速に探索出来るように、NextNodesに追加して、次のnodeとする。
		//NextNodesにもallNodesにも同じ局面のbattleが存在しないなら、新しくnodeを作り、
		//NextNodesと、allNodesに追加し、新しく作ったnodeを次のnodeとし、select処理を終了する。

		nextNode, err := node.NextNodes.Find(&battle)
		if err != nil {
			nextNode, err = allNodes.Find(&battle)
			if err == nil {
				node.NextNodes = append(node.NextNodes, nextNode)
			} else {
				nextNode = NewNodePointer(&battle, policy)
				allNodes = append(allNodes, nextNode)
				node.NextNodes = append(node.NextNodes, nextNode)
				break
			}
		}

		if nextNode.SelectCount == 1 {
			break
		}
		node = nextNode
	}
	return battle, allNodes, selects, nil
}

func (node *Node) AverageReward() float64 {
	accumReward := 0.0
	for _, pucb := range node.ActionCmdPUCBs {
		accumReward += pucb.AverageReward()
	}
	return float64(accumReward) / float64(len(node.ActionCmdPUCBs))
}

type Nodes []*Node

func (nodes Nodes) Find(battle *Battle) (*Node, error) {
	for _, node := range nodes {
		if node.Battle.Equal(battle) {
			return node, nil
		}
	}
	return &Node{}, fmt.Errorf("battleが一致しているnodeが見つからなかった")
}

type Select struct {
	Node      *Node
	ActionCmd ActionCmd
}

type Selects []Select

func (selects Selects) Backward(evalY float64, eval *Eval) {
	for _, select_ := range selects {
		node := select_.Node
		actionCmd := select_.ActionCmd

		if node.IsP1 {
			node.ActionCmdPUCBs[actionCmd].AccumReward += evalY
		} else {
			node.ActionCmdPUCBs[actionCmd].AccumReward += eval.Reverse(evalY)
		}
		node.ActionCmdPUCBs[actionCmd].Trial += 1
		node.SelectCount = 0
	}
}

func RunMCTS(rootBattle Battle, simuNum int, X float64, policy Policy, eval *Eval, random *rand.Rand) (Nodes, error) {
	rootNode := NewNodePointer(&rootBattle, policy)
	allNodes := Nodes{rootNode}
	battle := rootBattle

	var selects Selects
	var evalY float64
	var err error

	for i := 0; i < simuNum; i++ {
		battle, allNodes, selects, err = rootNode.SelectAndExpansion(battle, allNodes, policy, X, random)
		if err != nil {
			return Nodes{}, err
		}

		evalY, err = eval.Func(&battle)

		if err != nil {
			return Nodes{}, err
		}

		selects.Backward(evalY, eval)
		battle = rootBattle
	}
	return allNodes, nil
}

func NewMCTSTrainer(simuNum int, X float64, policy Policy, eval *Eval, random *rand.Rand) Trainer {
	result := func(battle *Battle) (ActionCmd, error) {
		legalActionCmds := battle.P1Fighters.LegalActionCmds()

		if len(legalActionCmds) == 1 {
			return legalActionCmds[0], nil
		}

		allNodes, err := RunMCTS(*battle, simuNum, X, policy, eval, random)

		if err != nil {
			return "", err
		}

		return allNodes[0].ActionCmdPUCBs.MaxTrialActionCmds().RandomChoice(random), nil
	}
	return result
}
