package battle

import (
	"fmt"
	bp "github.com/sw965/bippa"
	"strings"
	"github.com/sw965/omw/fn"
)

type Message string

func NewCriticalMessage() Message {
	return "きゅうしょに あたった！"
}

func NewSuperEffective() Message {
	return "こうかは ばつぐんだ！"
}

func NewNotVeryEffective() Message {
	return "こうかは いまひとつの ようだ..."
}

func NewParalysisCuredMessage() Message {
	return "まひが なおった！"
}

func NewSleepCuredMessage() Message {
	return "めを さました！"
}

func NewFreezeCuredMessage() Message {
	return "こおりが とけた！"
}

func NewBurnCuredMessage() Message {
	return "やけどが なおった！"
}

func NewStatusAilmentCuredMessage(s bp.StatusAilment) Message {
	switch s {
		case bp.PARALYSIS:
			return NewParalysisCuredMessage()
		case bp.SLEEP:
			return NewSleepCuredMessage()
		case bp.FREEZE:
			return NewFreezeCuredMessage()
		case bp.BURN:
			return NewBurnCuredMessage()
	}
	return ""
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

type MessageMaker struct {
	HumanTitle bp.HumanTitle
	HumanName bp.HumanName
	IsSelf bool
}

func (mm *MessageMaker) FullName() Message {
	return Message(mm.HumanTitle) + "の " + Message(mm.HumanName)
}

func (mm *MessageMaker) HumanNamePrefix() Message {
	if mm.IsSelf {
		return ""
	} else {
		return Message(mm.HumanName) + "の "
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
	m := fmt.Sprintf("%sの " + "%s！", pokeName.ToString(), moveName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) NoEffective(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sには ", pokeName.ToString())
    return mm.HumanNamePrefix() + Message(m) + "こうかが ない ようだ..."
}

func (mm *MessageMaker) TypeEffective(pokeName bp.PokeName, effective bp.TypeEffective) Message {
	switch effective {
		case bp.NEUTRAL_EFFECTIVE:
			return ""
		case bp.SUPER_EFFECTIVE:
			m := fmt.Sprintf("%sに 効果は バツグンだ！", pokeName.ToString())
			return mm.HumanNamePrefix() + Message(m)
		case bp.NOT_VERY_EFFECTIVE:
			m := fmt.Sprintf("%sに 効果は いまひとつだ", pokeName.ToString())
			return mm.HumanNamePrefix() + Message(m)
		case bp.NO_EFFECTIVE:
			m := fmt.Sprintf("%sには こうかが ない ようだ", pokeName.ToString())
			return mm.HumanNamePrefix() + Message(m)
	}
	return ""
}

func (mm *MessageMaker) Recoil(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは 攻撃の 反動を 受けた", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
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
	m := fmt.Sprintf("%s は 倒れた！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) Paralysis(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは まひして わざが でにくくなった！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) ParalysisLock(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは しびれてうごけない！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) Burn(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは やけどをおった！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) BurnDamage(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは やけどの ダメージを うけている！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) Sleep(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは ねむってしまった！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) SleepLock(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは ぐうぐう ねむっている", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) SleepCured(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは ", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m) + NewSleepCuredMessage()
}

func (mm *MessageMaker) Freeze(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは こおりついた！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) FreezeLock(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは こおって しまって うごかない！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) FreezeCured(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sの ", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m) + NewFreezeCuredMessage()
}

func (mm *MessageMaker) IntimidateAndClearBody(iName, cName bp.PokeName) Message {
	m := fmt.Sprintf("%sの クリアボディで %sの いかくは きかなかった！", cName.ToString(), iName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) ChestoBerry(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは カゴのみで ", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m) + NewSleepCuredMessage()
}

func (mm *MessageMaker) LumBerry(pokeName bp.PokeName, status bp.StatusAilment) Message {
	m := fmt.Sprintf("%sは ラムのみで", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m) + NewStatusAilmentCuredMessage(status)
}

func (mm *MessageMaker) TrickRoom(pokeName bp.PokeName, distort bool) Message {
	var m string
	if distort {
		m = fmt.Sprintf("%sは じくうを ゆがめた！", pokeName.ToString())
	} else {
		m = fmt.Sprintf("%sは じくうを もどした！", pokeName.ToString())
	}
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) FollowMe(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは ちゅうもくの まとになった！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) BellyDrum(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは たいりょくを けずって パワーぜんかいに なった！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}

func (mm *MessageMaker) Taunt(pokeName bp.PokeName) Message {
	m := fmt.Sprintf("%sは ちょうはつに のってしまった！", pokeName.ToString())
	return mm.HumanNamePrefix() + Message(m)
}