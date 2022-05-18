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

type ActionCmdUCB1s map[ActionCmd]*UpperConfidenceBound1

func (actionCmdUCB1s ActionCmdUCB1s) Keys() ActionCmds {
  result := make(ActionCmds, 0, len(actionCmdUCB1s))
  for actionCmd, _ := range actionCmdUCB1s {
    result = append(result, actionCmd)
  }
  return result
}

func (actionCmdUCB1s ActionCmdUCB1s) TotalTrial() int {
  result := 0
  for _, ucb1 := range actionCmdUCB1s {
    result += ucb1.Trial
  }
  return result
}

func (actionCmdUCB1s ActionCmdUCB1s) Max(X float64) float64 {
  totalTrial := actionCmdUCB1s.TotalTrial()
  actionCmds := actionCmdUCB1s.Keys()
  result := actionCmdUCB1s[actionCmds[0]].Get(totalTrial, X)

  for _, actionCmd := range actionCmds[1:] {
    ucb1v := actionCmdUCB1s[actionCmd].Get(totalTrial, X)
    if ucb1v > result {
      result = ucb1v
    }
  }
  return result
}

func (actionCmdUCB1s ActionCmdUCB1s) MaxActionCmds(X float64) ActionCmds {
  max := actionCmdUCB1s.Max(X)
  totalTrial := actionCmdUCB1s.TotalTrial()
  result := make(ActionCmds, 0)

  for actionCmd, ucb1 := range actionCmdUCB1s {
    ucb1v := ucb1.Get(totalTrial, X)
    if ucb1v == max {
      result = append(result, actionCmd)
    }
  }
  return result
}

func (actionCmdUCB1s ActionCmdUCB1s) MaxTrial() int {
  actionCmds := actionCmdUCB1s.Keys()
  result := actionCmdUCB1s[actionCmds[0]].Trial

  for _, actionCmd := range actionCmds[1:] {
    trial := actionCmdUCB1s[actionCmd].Trial
    if trial > result {
      result = trial
    }
  }
  return result
}

func (actionCmdUCB1s ActionCmdUCB1s) MaxTrialActionCmds() ActionCmds {
  maxTrial := actionCmdUCB1s.MaxTrial()
  result := make(ActionCmds, 0)
  for actionCmd, ucb1 := range actionCmdUCB1s {
    trial := ucb1.Trial
    if trial == maxTrial {
      result = append(result, actionCmd)
    }
  }
  return result
}

type Node struct {
  Battle *Battle
  LegalActionCmds ActionCmds
  ActionCmdUCB1s ActionCmdUCB1s
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

  legalActionCmds := fighters.LegalActionCmds()
  return &Node{Battle:battle, LegalActionCmds:legalActionCmds, ActionCmdUCB1s:ActionCmdUCB1s{},
               IsAllExpansion:false, IsP1Node:isP1Phase, SelectCount:0}
}

func (node *Node) NoExpansionActionCmds() ActionCmds {
  legalActionCmds := node.LegalActionCmds
  result := make(ActionCmds, 0, len(legalActionCmds) - len(node.ActionCmdUCB1s))

  for _, actionCmd := range legalActionCmds {
    if _, ok := node.ActionCmdUCB1s[actionCmd]; !ok {
      result = append(result, actionCmd)
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

    maxUCBActionCmds := node.ActionCmdUCB1s.MaxActionCmds(X)
    selectActionCmd := maxUCBActionCmds.RandomChoice(random)
    selects = append(selects, Select{Node:node, ActionCmd:selectActionCmd})
    node.SelectCount += 1

    battle, err = battle.Push(selectActionCmd, random)
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

  noExpansionActionCmd := node.NoExpansionActionCmds().RandomChoice(random)
  var err error

  //ランダムに選択した未展開ActionCmdでバトルを進める
  battle, err = battle.Push(noExpansionActionCmd, random)
  if err != nil {
    return 0.0, Selects{}, err
  }

  //ランダムに選択した未展開ActionCmdのUCBParamを新しく作り、selectsに追加する
  node.ActionCmdUCB1s[noExpansionActionCmd] = &UpperConfidenceBound1{}
  selects = append(selects, Select{Node:node, ActionCmd:noExpansionActionCmd})

  if len(node.ActionCmdUCB1s) == len(node.LegalActionCmds) {
    node.IsAllExpansion = true
  }

  //局面を評価する。
  evalY, err := eval.Func(&battle)
  return evalY, selects, err
}

func (node *Node) AverageReward() float64 {
  accumReward := 0.0
  for _, ucb1 := range node.ActionCmdUCB1s {
    accumReward += ucb1.AccumReward
  }
  return float64(accumReward) / float64(node.ActionCmdUCB1s.TotalTrial())
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
  ActionCmd ActionCmd
}

type Selects []Select

func (selects Selects) Backward(evalY float64, eval *Eval) {
  for _, select_ := range selects {
    node := select_.Node
    actionCmd := select_.ActionCmd

    if node.IsP1Node {
      node.ActionCmdUCB1s[actionCmd].AccumReward += evalY
    } else {
      node.ActionCmdUCB1s[actionCmd].AccumReward += eval.Reverse(evalY)
    }
    node.ActionCmdUCB1s[actionCmd].Trial += 1
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
  result := func(battle *Battle) (ActionCmd, error) {
    legalActionCmds := battle.P1Fighters.LegalActionCmds()

    if len(legalActionCmds) == 1 {
      return legalActionCmds[0], nil
    }

    allNodes, err := RunMCTS(*battle, simuNum, X, eval, random)

    if err != nil {
      return "", err
    }

    return allNodes[0].ActionCmdUCB1s.MaxTrialActionCmds().RandomChoice(random), nil
  }
  return result
}
