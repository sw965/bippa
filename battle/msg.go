package msg

import (
	"fmt"
	bp "github.com/sw965/bippa"
	"strings"
	"github.com/sw965/omw/fn"
)

type Message string

func NewChallengeByTrainerMessage(trainerName string) Message {
	ret := fmt.Sprintf("%sが ", trainerName)
	ret += "勝負を しかけてきた！"
	return Message(ret)
}

func NewActionPromptMessage(pokeName bp.PokeName) Message {
    return Message(fmt.Sprintf("%sは どうする？", pokeName.ToString()))
}

func NewMoveUseMessage(pokeName bp.PokeName, moveName bp.MoveName, isSelf bool) Message {
	m := map[bool]string{
		true:"",
		false:"相手の ",
	}[isSelf]
	return Message(fmt.Sprintf(m + "%s の " + "%s！", pokeName.ToString(), moveName.ToString()))
}

func NewRecoilMessage(trainerName string, pokeName bp.PokeName) Message {
	return Message(fmt.Sprintf("%sの %sは 攻撃の 反動を 受けた", trainerName, pokeName.ToString()))
}

func NewGoMessage(trainerName string, pokeName bp.PokeName, isSelf bool) Message {
	if isSelf {
		return Message(fmt.Sprintf("行け！ %s！", pokeName.ToString()))
	} else {
		ret := fmt.Sprintf("%sは", trainerName)
		ret += fmt.Sprintf("%sを 繰り出した！", pokeName.ToString())
		return Message(ret)
	}
}

func NewBackMessage(trainerName string, pokeName bp.PokeName, isSelf bool) Message {
	if isSelf {
		return Message(fmt.Sprintf("戻れ！ %s", pokeName.ToString()))
	} else {
		return Message(fmt.Sprintf("%s は %s を 引っ込めた！", trainerName, pokeName.ToString()))
	}
}

func NewFaintMessage(trainerName string, pokeName bp.PokeName, isSelf bool) Message {
	m := map[bool]string{
		true:"",
		false:trainerName + "の ",
	}[isSelf]
	return Message(fmt.Sprintf("%s%s は 倒れた！", m, pokeName.ToString()))
}

func (m Message) ToSlice() []Message {
	slice := strings.Split(string(m), "")
	ret := make([]Message, len(slice))
	for i, s := range slice {
		ret[i] = Message(s)
	}
	return ret
}

func (m Message) Accumulate() []Message {
	return fn.Accumulate(m.ToSlice())
}