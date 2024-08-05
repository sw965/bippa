package battle

import (
	"fmt"
	bp "github.com/sw965/bippa"
	"strings"
	"github.com/sw965/omw/fn"
)

type Message string

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

type MessageMaker struct {
	HumanTitle bp.HumanTitle
	HumanName bp.HumanName
	IsSelf bool
}

func (mm *MessageMaker) FullName() Message {
	if mm.IsSelf {
		return Message(string(mm.HumanName))
	} else {
		return Message(string(mm.HumanTitle) + "の " + string(mm.HumanName))
	}
}

func (mm *MessageMaker) ChallengeByTrainer() Message {
	m := fmt.Sprintf("%sが 勝負をしかけてきた！", mm.HumanName)
	return Message(m)
}

func (mm *MessageMaker) ActionPrompt(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは どうする？", pokeName.ToString())
    return Message(m)
}

func (mm *MessageMaker) MoveUse(pokeName bp.PokeName, moveName bp.MoveName) Message {
	var h string
	if mm.IsSelf {
		h = ""
	} else {
		h = string(mm.HumanName) + "の "
	}
	m := fmt.Sprintf(h + "%sの " + "%s！", pokeName.ToString(), moveName.ToString())
	return Message(m)
}

func (mm *MessageMaker) Recoil(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sの %sは 攻撃の 反動を 受けた", mm.HumanName, pokeName.ToString())
	return Message(m)
}

func (mm *MessageMaker) Go(pokeName bp.PokeName) Message {
	if mm.IsSelf {
		return Message(fmt.Sprintf("行け！ %s！", pokeName.ToString()))
	} else {
		m := fmt.Sprintf("%sは ", mm.FullName())
		m += fmt.Sprintf("%sを 繰り出した！", pokeName.ToString())
		return Message(m)
	}
}

func (mm *MessageMaker) Back(pokeName bp.PokeName) Message {
	if mm.IsSelf {
		return Message(fmt.Sprintf("戻れ！ %s", pokeName.ToString()))
	} else {
		return Message(fmt.Sprintf("%sは %sを 引っ込めた！", mm.FullName(), pokeName.ToString()))
	}
}

func (mm *MessageMaker) Faint(pokeName bp.PokeName) Message {
	var h string
	if mm.IsSelf {
		h = ""
	} else {
		h = string(mm.HumanName) + "の "
	}
	m := fmt.Sprintf(h + "%s は 倒れた！", h, pokeName.ToString())
	return Message(m)
}

func (mm *MessageMaker) TrickRoom(pokeName bp.PokeName, distort bool) Message {
	var h string
	if mm.IsSelf {
		h = ""
	} else {
		h = string(mm.HumanName) + "の "
	}

	var m string
	if distort {
		m = fmt.Sprintf(h + "%sは じくうを ゆがめた！", pokeName)
	} else {
		m = fmt.Sprintf(h + "%sは じくうを もどした！", pokeName)
	}
	return Message(m)
}