package mcts

import (
  "fmt"
  "math/rand"
  "github.com/sw965/crow"
  bp "github.com/sw965/bippa"
)

type UCBParam struct {
  TotalValue float64
  SimuNum int
}

func (ucbParam *UCBParam) Average() (float64, error) {
  if ucbParam.SimuNum == 0 {
    return 0.0, fmt.Errorf("0除算")
  }
  return float64(ucbParam.TotalValue) / float64(ucbParam.SimuNum), nil
}

type Node struct {
  Battle *bp.Battle
  AvailableBattleCommands bp.BattleCommands
  NextNodeParam NextNodeParam
  NextNodes Nodes
  IsAllExpansion bool
  IsP1Node bool
  SelectCount int
}

func NewNodePointer(battle *bp.Battle) *Node {
  isP1Phase := battle.IsP1Phase()
  var fighters bp.Fighters

  if isP1Phase {
    fighters = battle.P1Fighters
  } else {
    fighters = battle.P2Fighters
  }

  availableBattleCommands := fighters.AvailableBattleCommands()
  return &Node{Battle:battle, AvailableBattleCommands:availableBattleCommands, NextNodeParam:NextNodeParam{},
               IsAllExpansion:false, IsP1Node:isP1Phase, SelectCount:0}
}

func (node *Node) NewNoExpansionBattleCommands() bp.BattleCommands {
  availableBattleCommands := node.AvailableBattleCommands
  result := make(bp.BattleCommands, 0, len(availableBattleCommands) - len(node.NextNodeParam))

  for _, battleCommand := range availableBattleCommands {
    if _, ok := node.NextNodeParam[battleCommand]; !ok {
      result = append(result, battleCommand)
    }
  }
  return result
}

func (node *Node) Select(battle bp.Battle, allNodes Nodes, X float64,
                         random *rand.Rand) (*Node, bp.Battle, Nodes, Selects, bool, error) {
  selects := Selects{}
  isRoopSelect := false

  for {
    if !node.IsAllExpansion {
      break
    }

    maxUCBBattleCommands, err := node.NextNodeParam.MaxUCBBattleCommands(X)

    if err != nil {
      return &Node{}, bp.Battle{}, Nodes{}, Selects{}, false, err
    }

    selectBattleCommand := maxUCBBattleCommands.RandomChoice(random)
    selects = append(selects, Select{Node:node, BattleCommand:selectBattleCommand})
    node.SelectCount += 1

    battle, err = battle.Run(selectBattleCommand, random)
    if err != nil {
      return &Node{}, bp.Battle{}, Nodes{}, Selects{}, false, err
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

func (node *Node) ExpansionWithEvalY(battle bp.Battle, eval *Eval, selects Selects, random *rand.Rand) (float64, Selects, error) {
  if node.IsAllExpansion {
    return 0.0, Selects{}, fmt.Errorf("展開済みのNodeである")
  }

  noExpansionBattleCommand := node.NoExpansionBattleCommands().RandomChoice(random)
  var err error

  //ランダムに選択した未展開BattleCommandでバトルを進める
  battle, err = battle.Run(noExpansionBattleCommand, random)
  if err != nil {
    return 0.0, Selects{}, err
  }

  //ランダムに選択した未展開BattleCommandのUCBParamを新しく作り、selectsに追加する
  node.NextNodeParam[noExpansionBattleCommand] = &UCBParam{}
  selects = append(selects, Select{Node:node, BattleCommand:noExpansionBattleCommand})

  if len(node.NextNodeParam) == len(node.AvailableBattleCommands) {
    node.IsAllExpansion = true
  }

  //局面を評価する。
  evalY, err := eval.Func(&battle)
  return evalY, selects, err
}

func (node *Node) Average() float64 {
  totalValue := 0.0
  totalSimuNum := 0

  for _, ucbParam := range node.NextNodeParam {
    totalValue += ucbParam.TotalValue
    totalSimuNum += ucbParam.SimuNum
  }
  return float64(totalValue) / float64(totalSimuNum)
}

type Nodes []*Node

func (nodes Nodes) Find(battle *bp.Battle) (*Node, error) {
  for _, node := range nodes {
    if node.Battle.Equal(battle) {
      return node, nil
    }
  }
  return &Node{}, fmt.Errorf("battleが一致しているnodeが見つからなかった")
}

type NextNodeParam map[bp.BattleCommand]*UCBParam

func (nextNodeParam NextNodeParam) NewBattleCommands() bp.BattleCommands {
  result := make(bp.BattleCommands, 0, len(nextNodeParam))
  for battleCommand, _ := range nextNodeParam {
    result = append(result, battleCommand)
  }
  return result
}

func (nextNodeParam NextNodeParam) TotalSimuNum() int {
  result := 0
  for _, ucbParam := range nextNodeParam {
    result += ucbParam.SimuNum
  }
  return result
}

func (nextNodeParam NextNodeParam) UCBs(X float64) (map[bp.BattleCommand]float64, error) {
  totalSimuNum := nextNodeParam.TotalSimuNum()
  result := map[bp.BattleCommand]float64{}

  for battleCommand, ucbParam := range nextNodeParam {
    average, err := ucbParam.Average()
    if err != nil {
      return map[bp.BattleCommand]float64{}, err
    }

    ucb, err := crow.UpperConfidenceBound1(average, totalSimuNum, ucbParam.SimuNum, X)
    if err != nil {
      return map[bp.BattleCommand]float64{}, err
    }

    result[battleCommand] = ucb
  }

  return result, nil
}

func (nextNodeParam NextNodeParam) MaxUCB(X float64) (float64, error) {
  battleCommands := nextNodeParam.BattleCommands()
  ucbs, err := nextNodeParam.UCBs(X)

  if err != nil {
    return 0.0, err
  }

  maxUCB := ucbs[battleCommands[0]]
  for _, battleCommand := range battleCommands[1:] {
    ucb := ucbs[battleCommand]
    if ucb > maxUCB {
      maxUCB = ucb
    }
  }
  return maxUCB, nil
}

func (nextNodeParam NextNodeParam) NewMaxUCBBattleCommands(X float64) (bp.BattleCommands, error) {
  maxUCB, err := nextNodeParam.MaxUCB(X)
  if err != nil {
    return bp.BattleCommands{}, err
  }

  ucbs, err := nextNodeParam.UCBs(X)
  if err != nil {
    return bp.BattleCommands{}, err
  }

  result := make(bp.BattleCommands, 0)
  for battleCommand, ucb := range ucbs {
    if ucb == maxUCB {
      result = append(result, battleCommand)
    }
  }
  return result, nil
}

func (nextNodeParam NextNodeParam) MaxSimuNum() int {
  battleCommands := nextNodeParam.NewBattleCommands()
  result := nextNodeParam[battleCommands[0]].SimuNum

  for _, battleCommand := range battleCommands[1:] {
    simuNum := nextNodeParam[battleCommand].SimuNum
    if simuNum > result {
      result = simuNum
    }
  }
  return result
}

func (nextNodeParam NextNodeParam) NewMaxSimuNumBattleCommands() bp.BattleCommands {
  maxSimuNum := nextNodeParam.MaxSimuNum()
  result := make(bp.BattleCommands, 0)
  for battleCommand, ucbParam := range nextNodeParam {
    simuNum := ucbParam.SimuNum
    if simuNum == maxSimuNum {
      result = append(result, battleCommand)
    }
  }
  return result
}

type Select struct {
  Node *Node
  BattleCommand bp.BattleCommand
}

type Selects []Select

func (selects Selects) Backward(evalY float64, eval *Eval) {
  for _, select_ := range selects {
    node := select_.Node
    battleCommand := select_.BattleCommand
    if node.IsP1Node {
      node.NextNodeParam[battleCommand].TotalValue += evalY
    } else {
      node.NextNodeParam[battleCommand].TotalValue += eval.ReverseFunc(evalY)
    }
    node.NextNodeParam[battleCommand].SimuNum += 1
    node.SelectCount = 0
  }
}

func RunMCTS(rootBattle bp.Battle, simuNum int, X float64, eval *Eval, random *rand.Rand) (Nodes, error) {
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

func NewMCTSTrainer(simuNum int, X float64, eval *Eval, random *rand.Rand) bp.Trainer {
  result := func(battle *bp.Battle) (bp.BattleCommand, error) {
    allNodes, err := RunMCTS(*battle, simuNum, X, eval, random)

    if err != nil {
      return "", err
    }

    return allNodes[0].NextNodeParam.NewMaxSimuNumBattleCommands().RandomChoice(random), nil
  }
  return result
}
