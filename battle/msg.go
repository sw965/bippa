package battle

// import (
// 	"fmt"
// 	bp "github.com/sw965/bippa"
// 	"strings"
// 	"github.com/sw965/omw/fn"
// )

// type Message string

// func NewCriticalMessage() Message {
// 	return "きゅうしょに あたった！"
// }

// func NewSuperEffective() Message {
// 	return "こうかは ばつぐんだ！"
// }

// func NewNotVeryEffective() Message {
// 	return "こうかは いまひとつの ようだ..."
// }

// func NewParalysisCuredMessage() Message {
// 	return "まひが なおった！"
// }

// func NewSleepCuredMessage() Message {
// 	return "めを さました！"
// }

// func NewFreezeCuredMessage() Message {
// 	return "こおりが とけた！"
// }

// func NewBurnCuredMessage() Message {
// 	return "やけどが なおった！"
// }

// func NewStatusAilmentCuredMessage(s bp.StatusAilment) Message {
// 	switch s {
// 		case bp.PARALYSIS:
// 			return NewParalysisCuredMessage()
// 		case bp.SLEEP:
// 			return NewSleepCuredMessage()
// 		case bp.FREEZE:
// 			return NewFreezeCuredMessage()
// 		case bp.BURN:
// 			return NewBurnCuredMessage()
// 	}
// 	return ""
// }

// func (m Message) ToSlice() []Message {
// 	slice := strings.Split(string(m), "")
// 	ret := make([]Message, len(slice))
// 	for i, s := range slice {
// 		ret[i] = Message(s)
// 	}
// 	return ret
// }

// func (m Message) Accumulate() []Message {
// 	return fn.Accumulate(m.ToSlice())
// }

// type MessageMaker struct {
// 	HostHumanTitle bp.HumanTitle
// 	HostHumanName bp.HumanName
// 	GuestHumanTitle bp.HumanTitle
// 	GuestHumanName bp.HumanName
// 	IsHostView bool
// }

// func NewMessageMaker(m *Manager) MessageMaker {
// 	return MessageMaker{
// 		HumanTitle:m.OpponentHumanTitle,
// 		HumanName:m.OpponentHumanName,
// 		IsHostView:m.IsHostView,
// 	}
// }

// func (mm *MessageMaker) GetCurrentSelfHumanName() HumanName {
// 	if mm.IsHostView {
// 		return mm.HostHumanName
// 	} else {
// 		return mm.GuestHumanName
// 	}
// }

// func (mm *MessageMaker) SelfFullName() (Message, Message) {
// 	return Message(mm.OpponentHumanName)
// }

// func (mm *MessageMaker) OpponentFullName() (Message, Message) {
// 	return Message(mm.OpponentHumanTitle) + "の " + Message(mm.OpponentHumanName)
// }

// func (mm *MessageMaker) OpponentHumanNamePrefix() Message {
// 	if mm.IsHostView {
// 		return ""
// 	} else {
// 		return Message(mm.OpponentHumanName) + "の "
// 	}
// }

// func (mm *MessageMaker) MoveUse(pokeName bp.PokeName, moveName bp.MoveName, isHost bool) Message {
// 	pokeNameStr := pokeName.ToString()
// 	moveNameStr := moveName.ToString()

// 	var hostHumanName bp.HumanName
// 	var guestHumanName bp.HumanName

// 	if isHost {
// 		hostPrefix = ""
// 		guestPrefix = mm.HostHumanName + "の "
// 	} else {
// 		hostPrefix = mm.GuestHumanName + "の "
// 		guestPrefix = ""
// 	}

// 	host = fmt.Sprintf("%s%sの %s！", hostPrefix, pokeNameStr, moveNameStr)
// 	guest = fmt.Sprintf("%s%sの %s！", guestPrefix, pokeNameStr, moveNameStr)

// 	return Message(host), Message(guest)
// }

// func (mm *MessageMaker) NoEffective(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sには ", pokeName.ToString())
//     return mm.HumanNamePrefix() + Message(m) + "こうかが ない ようだ..."
// }

// func (mm *MessageMaker) TypeEffective(pokeName bp.PokeName, effective bp.TypeEffective) Message {
// 	switch effective {
// 		case bp.NEUTRAL_EFFECTIVE:
// 			return ""
// 		case bp.SUPER_EFFECTIVE:
// 			m := fmt.Sprintf("%sに 効果は バツグンだ！", pokeName.ToString())
// 			return mm.HumanNamePrefix() + Message(m)
// 		case bp.NOT_VERY_EFFECTIVE:
// 			m := fmt.Sprintf("%sに 効果は いまひとつだ", pokeName.ToString())
// 			return mm.HumanNamePrefix() + Message(m)
// 		case bp.NO_EFFECTIVE:
// 			m := fmt.Sprintf("%sには こうかが ない ようだ", pokeName.ToString())
// 			return mm.HumanNamePrefix() + Message(m)
// 	}
// 	return ""
// }

// func (mm *MessageMaker) Recoil(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは 攻撃の 反動を 受けた", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Go(pokeName bp.PokeName) Message {
// 	if mm.IsSelf {
// 		return Message(fmt.Sprintf("行け！ %s！", pokeName.ToString()))
// 	} else {
// 		m := fmt.Sprintf("%sは ", mm.FullName())
// 		m += fmt.Sprintf("%sを 繰り出した！", pokeName.ToString())
// 		return Message(m)
// 	}
// }

// func (mm *MessageMaker) Back(pokeName bp.PokeName) Message {
// 	if mm.IsSelf {
// 		return Message(fmt.Sprintf("戻れ！ %s", pokeName.ToString()))
// 	} else {
// 		return Message(fmt.Sprintf("%sは %sを 引っ込めた！", mm.FullName(), pokeName.ToString()))
// 	}
// }

// func (mm *MessageMaker) Faint(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%s は 倒れた！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Paralysis(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは まひして わざが でにくくなった！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) ParalysisLock(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは しびれてうごけない！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Burn(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは やけどをおった！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) BurnDamage(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは やけどの ダメージを うけている！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Sleep(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ねむってしまった！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) SleepLock(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ぐうぐう ねむっている", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) SleepCured(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m) + NewSleepCuredMessage()
// }

// func (mm *MessageMaker) Freeze(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは こおりついた！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) FreezeLock(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは こおって しまって うごかない！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) FreezeCured(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sの ", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m) + NewFreezeCuredMessage()
// }

// func (mm *MessageMaker) IntimidateAndClearBody(iName, cName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sの クリアボディで %sの いかくは きかなかった！", cName.ToString(), iName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) ChestoBerry(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは カゴのみで ", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m) + NewSleepCuredMessage()
// }

// func (mm *MessageMaker) LumBerry(pokeName bp.PokeName, status bp.StatusAilment) Message {
// 	m := fmt.Sprintf("%sは ラムのみで", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m) + NewStatusAilmentCuredMessage(status)
// }

// func (mm *MessageMaker) TrickRoom(pokeName bp.PokeName, distort bool) Message {
// 	var m string
// 	if distort {
// 		m = fmt.Sprintf("%sは じくうを ゆがめた！", pokeName.ToString())
// 	} else {
// 		m = fmt.Sprintf("%sは じくうを もどした！", pokeName.ToString())
// 	}
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) FollowMe(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ちゅうもくの まとになった！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) BellyDrum(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは たいりょくを けずって パワーぜんかいに なった！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Taunt(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ちょうはつに のってしまった！", pokeName.ToString())
// 	return mm.HumanNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Intimidate(src, target bp.PokeName) Message {
// 	if mm.IsSelf {
// 		return Message(fmt.Sprintf("%sの いかくで %sの %sの こうげきが さがった！", src.ToString(), mm.HumanName, target.ToString()))
// 	} else {
// 		return Message(fmt.Sprintf("%sの %sの いかくで %sの こうげきがさがった！", mm.HumanName, src.ToString(), target.ToString()))
// 	}
// }

// func (mm *MessageMaker) Init() {

// }