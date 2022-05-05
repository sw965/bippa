package bippa

import (
  "fmt"
  "math"
  "math/rand"
)

func UpperConfidenceBound1(averageReward float64, totalTrial, trial int, X float64) float64 {
  floatTotalTrial := float64(totalTrial)
  floatTrial := float64(trial)
  return averageReward + (X * math.Sqrt(2 * math.Log(floatTotalTrial) / floatTrial))
}

type UCT struct {
  AccumReward float64
  Trial int
}

func (uct UCT) AverageReward() float64 {
  return uct.AccumReward / float64(uct.Trial)
}

func (uct UCT) UCB(totalTrial int, X float64) float64 {
  averageReward := uct.AverageReward()
  return UpperConfidenceBound1(averageReward, totalTrial, uct.Trial, X)
}

type ActionCommandUCTs map[ActionCommand]*UCT

func (actionCommandUCTs ActionCommandUCTs) Keys() ActionCommand {
  result := make(ActionCommand, 0, len(actionCommandUCTs))
  for key, _ := range actionCommandUCTs {
    result = append(result, key)
  }
  return result
}

func (actionCommandUCTs ActionCommandUCTs) TotalTrial() int {
  result := 0
  for _, uct := range actionCommandUCTs {
    result += uct.Trial
  }
  return result
}

func (actionCommandUCTs ActionCommandUCTs) MaxUCB(X float64) float64 {
  actionCommands := actionCommandUCTs.Keys()
  totalTrial := actionCommandUCTs.TotalTrial()
  result := actionCommandUCTs[actionCommands[0]].UCB(totalTrial, X)

  for _, actionCommand := range actionCommands[1:] {
    ucb := actionCommandUCTs[actionCommand].UCB(totalTrial, X)
    if ucb > result {
      result = ucb
    }
  }
  return result
}

func (actionCommandUCTs ActionCommandUCTs) MaxUCBActionCommand(X float64) Actions {
  maxUCB := actionCommandUCTs.MaxUCB(X)
  totalTrial := actionCommandUCTs.TotalTrial()
  result := make(ActionCommands, 0)

  for actionCommand, uct := range actionCommandUCTs {
    ucb := uct.UCB(totalTrial, X)
    if ucb == maxUCB {
      result = append(result, actionCommand)
    }
  }
  return result
}

func (actionCommandUCTs ActionCommandUCTs) MaxTrial() int {
  actionCommands := actionCommandUCTs.Keys()
  result := actionCommandUCTs[actionCommands[0]].Trial

  for _, actionCommand := range actionCommands[1:] {
    trial := actionCommandUCTs[actionCommand].Trial

    if trial > result {
      result = trial
    }
  }
  return result
}

func (actionCommandUCTs ActionCommandUCTs) MaxTrialActionCommands() ActionCommands {
  maxTrial := actionCommandUCTs.MaxTrial()
  result := make(ActionCommands, 0)

  for actionCommand, uct := range actionCommandUCTs {
    trial := uct.Trial
    if trial == maxTrial {
      result = append(result, actionCommand)
    }
  }
  return result
}

type Node struct {
  Battle *Battle
  P1LegalActionCommands ActionCommands
  P2LegalActionCommands ActionCommands
  P1ActionCommandUCTs ActionCommandUCTs
  P2ActionCommandUCTs ActionCommandUCTs
  IsAllExpansion bool
  NextNodes Nodes
  SelectCount int
  AccumReward float64
}

func NewNodePointer(battle *Battle) *Node {
  var p1LegalActionCommands ActionCommands
  var p2LegalActionCommands ActionCommands

  if battle.IsP1OnlyAction() {
    p1LegalActionCommands = battle.P1Fighters.LegalActionCommands()
    p2LegalActionCommands = ActionCommands{}
  } else if battle.IsP2OnlyAction() {
    p1LegalActionCommands = ActionCommands{}
    p2LegalActionCommands = battle.P2Fighters.LegalActionCommands()
  } else {
    p1LegalActionCommands = battle.P1Fighters.LegalActionCommands()
    p2LegalActionCommands = battle.P2Fighters.LegalActionCommands()
  }

  return &Node{Battle:battle,
               P1LegalActionCommands:p1LegalActionCommands,
               P2LegalActionCommands:p2LegalActionCommands,
               P1ActionCommandUCTs:ActionCommandUCTs{},
               P2ActionCommandUCTs:ActionCommandUCTs{},
               IsAllExpansion:false, SelectCount:0}
}

func (node *Node) NoExpansionP1ActionCommands()  {
  p1LegalActionCommands := node.P1LegalActionCommands
  result := make(Actions, 0, len(p1LegalActionCommands) - len(node.ActionCommandUCTs))

  for _, actionCommand := range p1LegalActionCommands {
    if _, ok := node.P1ActionCommandUCTs[actionCommand]; !ok {
      result = append(result, actionCommand)
    }
  }
  return result
}

func (node *Node) NoExpansionP2ActionCommands()

func (node *Node) Select(battle Battle, allNodes Nodes, X float64, random *rand.Rand) (*Node, Battle, Nodes, Selects, bool, error) {
  selects := Selects{}
  isRoopSelect := false
  var err error

  for {
    if !node.IsAllExpansion {
      break
    }

    maxUCBActions := node.ActionCommandUCTs.MaxUCBActions(X)
    selectAction := maxUCBActions.RandomChoice(random)
    selects = append(selects, Select{Node:node, P1ActionCommand:selectAction.P1, P2ActionCommand:selectAction.P2})
    node.SelectCount += 1

    battle, err = battle.Push(&selectAction, random)

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
  battle, err = battle.Push(&noExpansionAction, random)
  if err != nil {
    return 0.0, Selects{}, err
  }

  //ランダムに選択した未展開ActionのDecoupledUCTを新しく作り、selectsに追加する
  node.ActionCommandUCTs[noExpansionAction] = &DecoupledUCT{}
  selects = append(selects, Select{Node:node, P1ActionCommand:noExpansionAction.P1, P2ActionCommand:noExpansionAction.P2})

  if len(node.ActionCommandUCTs) == len(node.LegalActions) {
    node.IsAllExpansion = true
  }

  //局面を評価する。
  evalY, err := eval.Func(&battle)
  return evalY, selects, err
}

func (node *Node) P1AndP2AverageReward() (float64, float64) {
  p1AccumReward := 0.0
  p2AccumReward := 0.0

  for _, decoupledUCT := range node.ActionCommandUCTs {
    p1AccumReward += decoupledUCT.P1.AccumReward
    p2AccumReward += decoupledUCT.P2.AccumReward
  }

  p1TotalTrial, p2TotalTrial := node.ActionCommandUCTs.P1AndP2TotalTrial()
  p1Result := float64(p1AccumReward) / float64(p1TotalTrial)
  p2Result := float64(p2AccumReward) / float64(p2TotalTrial)
  return p1Result, p2Result
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
  P1ActionCommand ActionCommand
  P2ActionCommand ActionCommand
}

type Selects []Select

func (selects Selects) Backward(evalY float64, eval *Eval) {
  //selectはキーワード構文なので後ろにアンダーバーを付けている
  for _, select_ := range selects {
    node := select_.Node
    action := Action{P1:select_.P1ActionCommand, P2:select_.P2ActionCommand}

    node.ActionCommandUCTs[action].P1.AccumReward += evalY
    node.ActionCommandUCTs[action].P2.AccumReward += eval.ReverseFunc(evalY)
    node.ActionCommandUCTs[action].P1.Trial += 1
    node.ActionCommandUCTs[action].P2.Trial += 1
    node.SelectCount = 0
  }
}

func RunDMCTS(rootBattle Battle, simuNum int, X float64, eval *Eval, random *rand.Rand) (Nodes, error) {
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
    allNodes, err := RunDMCTS(*battle, simuNum, X, eval, random)
    if err != nil {
      return Action{}, err
    }
    return allNodes[0].ActionCommandUCTs.MaxTrialActions().RandomChoice(random), nil
  }
  return result
}
