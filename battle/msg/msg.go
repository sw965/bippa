package msg

import (
	"fmt"
	bp "github.com/sw965/bippa"
	"strings"
	"github.com/sw965/omw/fn"
)

type Message string

func NewChallengeByTrainer(trainerName string) Message {
	return Message(fmt.Sprintf("%s が 勝負を しかけてきた！", trainerName))
}

func NewMoveUse(pokeName bp.PokeName, moveName bp.MoveName, isSelf bool) Message {
	m := map[bool]string{
		true:"",
		false:"相手の ",
	}[isSelf]
	return Message(fmt.Sprintf(m + "%s の " + "%s！", pokeName.ToString(), moveName.ToString()))
}

func NewBack(trainer string, pokeName bp.PokeName, isSelf bool) Message {
	if isSelf {
		return Message(fmt.Sprintf("戻れ！ %s", pokeName.ToString()))
	} else {
		return Message(fmt.Sprintf("%s は %s を 引っ込めた！", trainer, pokeName.ToString()))
	}
}

func NewGo(trainer string, pokeName bp.PokeName, isSelf bool) Message {
	if isSelf {
		return Message(fmt.Sprintf("行け！ %s", pokeName.ToString()))
	} else {
		return Message(fmt.Sprintf("%s は %s を 繰り出した！", trainer, pokeName.ToString()))
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