package bippa

import (
  "fmt"
  "math"
  "math/rand"
)

type UpperConfidenceBound1 struct {
  AccumReward float64
  Trial int
}

func (ucb1 *UpperConfidenceBound1) AverageReward() float64 {
  return float64(ucb1.AccumReward) / float64(ucb1.Trial)
}

func (ucb1 *UpperConfidenceBound1) Get(totalTrial int, X float64) float64 {
  floatTotalTrial := float64(totalTrial)
  floatTrial := float64(ucb1.Trial)
  averageReward := ucb1.AverageReward()
  return averageReward + (X * math.Sqrt(2 * math.Log(floatTotalTrial) / floatTrial))
}

type ActionUCB1s map[Action]*UpperConfidenceBound1

func (actionUCB1s ActionUCB1s) Keys() Actions {
  result := make(Actions, 0, len(actionUCB1s))
  for action, _ := range actionUCB1s {
    result = append(result, action)
  }
  return result
}

func (actionUCB1s ActionUCB1s) TotalTrial() int {
  result := 0
  for _, ucb1 := range actionUCB1s {
    result += ucb1.Trial
  }
  return result
}

func (actionUCB1s ActionUCB1s) Max(X float64) float64 {
  totalTrial := actionUCB1s.TotalTrial()
  actions := actionUCB1s.Keys()
  result := actionUCB1s[actions[0]].Get(totalTrial, X)

  for _, action := range actions[1:] {
    ucb1v := actionUCB1s[action].Get(totalTrial, X)
    if ucb1v > result {
      result = ucb1v
    }
  }
  return result
}

func (actionUCB1s ActionUCB1s) MaxActions(X float64) Actions {
  max := actionUCB1s.Max(X)
  totalTrial := actionUCB1s.TotalTrial()
  result := make(Actions, 0)

  for action, ucb1 := range actionUCB1s {
    ucb1v := ucb1.Get(totalTrial, X)
    if ucb1v == max {
      result = append(result, action)
    }
  }
  return result
}

func (actionUCB1s ActionUCB1s) MaxTrial() int {
  actions := actionUCB1s.Keys()
  result := actionUCB1s[actions[0]].Trial

  for _, action := range actions[1:] {
    trial := actionUCB1s[action].Trial
    if trial > result {
      result = trial
    }
  }
  return result
}

func (actionUCB1s ActionUCB1s) MaxTrialActions() Actions {
  maxTrial := actionUCB1s.MaxTrial()
  result := make(Actions, 0)
  for action, ucb1 := range actionUCB1s {
    trial := ucb1.Trial
    if trial == maxTrial {
      result = append(result, action)
    }
  }
  return result
}

type Node struct {
  Battle *Battle
  LegalActions Actions
  ActionUCB1s ActionUCB1s
  IsAllExpansion bool
  NextNodes Nodes
  IsP1Node bool
  SelectCount int
}

func NewNodePointer(battle *Battle) *Node {
  isP1Phase := battle.IsP1Phase()
  var fighters Fighters

  if isP1Phase {
    fighters = battle.P1Fighters
  } else {
    fighters = battle.P2Fighters
  }

  legalActions := fighters.LegalActions()
  return &Node{Battle:battle, LegalActions:legalActions, ActionUCB1s:ActionUCB1s{},
               IsAllExpansion:false, IsP1Node:isP1Phase, SelectCount:0}
}

func (node *Node) NoExpansionActions() Actions {
  legalActions := node.LegalActions
  result := make(Actions, 0, len(legalActions) - len(node.ActionUCB1s))

  for _, action := range legalActions {
    if _, ok := node.ActionUCB1s[action]; !ok {
      result = append(result, action)
    }
  }
  return result
}

func (node *Node) Select(battle Battle, allNodes Nodes, X float64, random *rand.Rand) (*Node, Battle, Nodes, Selects, bool, error) {
  selects := Selects{}
  isRoopSelect := false
  var err error

  for {
    if !node.IsAllExpansion {
      break
    }

    maxUCBActions := node.ActionUCB1s.MaxActions(X)
    selectAction := maxUCBActions.RandomChoice(random)
    selects = append(selects, Select{Node:node, Action:selectAction})
    node.SelectCount += 1

    battle, err = battle.Push(selectAction, random)
    if err != nil {
      return &Node{}, Battle{}, Nodes{}, Selects{}, false, err
    }

    if battle.IsGameEnd() {
      break
    }

    //NextNodesの中に、同じ局面のbattleが存在するならば、それを次のnodeとする
    //NextNodesの中に、同じ局面のbattleが存在しないなら、allNodesの中から同じ局面のbattleが存在しないかを調べる。
    //allNodesの中に、同じ局面のbattleが存在するならば、次回から高速に探索出来るように、NextNodesに追加して、次のnodeとする。
    //NextNodesにもallNodesにも同じ局面のbattleが存在しないなら、新しくnodeを作り、
    //NextNodesと、allNodesに追加し、新しく作ったnodeを次のnodeとする。
    //またnodeを新しく作った場合は、一番上のbreak条件に必ず引っかかる。

    nextNode, err := node.NextNodes.Find(&battle)
    if err != nil {
      nextNode, err = allNodes.Find(&battle)
      if err == nil {
        node.NextNodes = append(node.NextNodes, nextNode)
      } else {
        nextNode = NewNodePointer(&battle)
        allNodes = append(allNodes, nextNode)
        node.NextNodes = append(node.NextNodes, nextNode)
      }
    }

    if nextNode.SelectCount == 1 {
      isRoopSelect = true
      break
    }

    node = nextNode
  }

  return node, battle, allNodes, selects, isRoopSelect, nil
}

func (node *Node) ExpansionWithEvalY(battle Battle, eval *Eval, selects Selects, random *rand.Rand) (float64, Selects, error) {
  if node.IsAllExpansion {
    return 0.0, Selects{}, fmt.Errorf("展開済みのNodeである")
  }

  noExpansionAction := node.NoExpansionActions().RandomChoice(random)
  var err error

  //ランダムに選択した未展開Actionでバトルを進める
  battle, err = battle.Push(noExpansionAction, random)
  if err != nil {
    return 0.0, Selects{}, err
  }

  //ランダムに選択した未展開ActionのUCBParamを新しく作り、selectsに追加する
  node.ActionUCB1s[noExpansionAction] = &UpperConfidenceBound1{}
  selects = append(selects, Select{Node:node, Action:noExpansionAction})

  if len(node.ActionUCB1s) == len(node.LegalActions) {
    node.IsAllExpansion = true
  }

  //局面を評価する。
  evalY, err := eval.Func(&battle)
  return evalY, selects, err
}

func (node *Node) AverageReward() float64 {
  accumReward := 0.0
  for _, ucb1 := range node.ActionUCB1s {
    accumReward += ucb1.AccumReward
  }
  return float64(accumReward) / float64(node.ActionUCB1s.TotalTrial())
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
  Node *Node
  Action Action
}

type Selects []Select

func (selects Selects) Backward(evalY float64, eval *Eval) {
  for _, select_ := range selects {
    node := select_.Node
    action := select_.Action

    if node.IsP1Node {
      node.ActionUCB1s[action].AccumReward += evalY
    } else {
      node.ActionUCB1s[action].AccumReward += eval.Reverse(evalY)
    }
    node.ActionUCB1s[action].Trial += 1
    node.SelectCount = 0
  }
}

func RunMCTS(rootBattle Battle, simuNum int, X float64, eval *Eval, random *rand.Rand) (Nodes, error) {
  rootNode := NewNodePointer(&rootBattle)
  allNodes := Nodes{rootNode}

  node := rootNode
  battle := rootBattle

  var selects Selects
  var isRoopSelect bool
  var evalY float64
  var err error

  for i := 0; i < simuNum; i++ {
    node, battle, allNodes, selects, isRoopSelect, err = node.Select(battle, allNodes, X, random)
    if err != nil {
      return Nodes{}, err
    }

    //Select処理を行い、ゲームが終了したならば、その状態で評価関数を呼び出し、評価する。
    //RoopSelect状態ならば、必ず全てを展開しているNodeになっている(=展開出来ない)ので、そのまま評価を得る。
    //ゲームが終了せず、RoopSelect状態でなければ、必ず展開していないNodeがある状態なので、展開処理をして、評価を得る。

    if battle.IsGameEnd() || isRoopSelect {
      evalY, err = eval.Func(&battle)
    } else {
      evalY, selects, err = node.ExpansionWithEvalY(battle, eval, selects, random)
    }

    if err != nil {
      return Nodes{}, err
    }

    selects.Backward(evalY, eval)
    node = rootNode
    battle = rootBattle
  }
  return allNodes, nil
}

func NewMCTSTrainer(simuNum int, X float64, eval *Eval, random *rand.Rand) Trainer {
  result := func(battle *Battle) (Action, error) {
    allNodes, err := RunMCTS(*battle, simuNum, X, eval, random)

    if err != nil {
      return "", err
    }

    return allNodes[0].ActionUCB1s.MaxTrialActions().RandomChoice(random), nil
  }
  return result
}
