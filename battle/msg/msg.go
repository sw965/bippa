package msg

import (
	"fmt"
	bp "github.com/sw965/bippa"
	"strings"
	"github.com/sw965/omw/fn"
)

type Message string

func Box() string {
	ret := "+-----------------------------------------------------------+"
	ret += "|                                                           |"
	ret += "|                                                           |"
	ret += "+-----------------------------------------------------------+"
	return ret
}

func NewChallengeByTrainer(trainerName string, s string) Message {
	ret := fmt.Sprintf("%sが", trainerName)
	ret += s
	ret += "勝負を しかけてきた！"
	return Message(ret)
}

func NewActionPrompt(pokeName bp.PokeName) Message {
    return Message(fmt.Sprintf("%sは どうする？", pokeName.ToString()))
}

func NewMoveUse(pokeName bp.PokeName, moveName bp.MoveName, isSelf bool) Message {
	m := map[bool]string{
		true:"",
		false:"相手の ",
	}[isSelf]
	return Message(fmt.Sprintf(m + "%s の " + "%s！", pokeName.ToString(), moveName.ToString()))
}

func NewGo(trainerName string, pokeName bp.PokeName, isSelf bool, s string) Message {
	if isSelf {
		return Message(fmt.Sprintf("行け！ %s！", pokeName.ToString()))
	} else {
		ret := fmt.Sprintf("%sは", trainerName)
		ret += s
		ret += fmt.Sprintf("%sを 繰り出した！", pokeName.ToString())
		return Message(ret)
	}
}

func NewBack(trainer string, pokeName bp.PokeName, isSelf bool) Message {
	if isSelf {
		return Message(fmt.Sprintf("戻れ！ %s", pokeName.ToString()))
	} else {
		return Message(fmt.Sprintf("%s は %s を 引っ込めた！", trainer, pokeName.ToString()))
	}
}

func NewFaint(pokeName bp.PokeName, isSelf bool) Message {
	m := map[bool]string{
		true:"",
		false:"相手の ",
	}[isSelf]
	return Message(fmt.Sprintf("%s%s は 倒れた！", m, pokeName.ToString()))
}

func (m Message) ToSlice() []string {
	return strings.Split(string(m), "")
}

func (m Message) Accumulate() []string {
	return fn.Accumulate(m.ToSlice())
}