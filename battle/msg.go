package battle

import (
	"fmt"
	bp "github.com/sw965/bippa"
	//"strings"
	//"github.com/sw965/omw/fn"
)

func GetTrainerInfoMessageFunc(guestTrainerTitle, guestTrainerName string) func(bool) string {
	return func(isHost bool) string {
		if isHost {
			return ""
		}
		return guestTrainerTitle + "の " + guestTrainerName + "は "
	}
}

func GetTrainerNameMessageFunc(guestTrainerName string) func(bool) string {
	return func(isHost bool) string {
		if isHost {
			return ""
		}
		return guestTrainerName + "の "
	}
}

// https://wiki.xn--rckteqa2e.com/wiki/%E3%83%A9%E3%83%B3%E3%82%AF%E8%A3%9C%E6%AD%A3
func GetStandardRankFluctuationMessages(name bp.PokeName, v *bp.RankStat) []string {
	ms := make([]string, 0, 5)
	nameStr := name.ToString()

	add := func(v bp.Rank, s string) {
		if v == 1 {
			m := fmt.Sprintf("%sの %sが 上がった", nameStr, s)
			ms = append(ms, m)
		}
	
		if v == 2 {
			m := fmt.Sprintf("%sの %sが ぐーんと上がった", nameStr, s)
			ms = append(ms, m)
		}
	
		if v >= 3 {
			m := fmt.Sprintf("%sの %sが ぐぐーんと上がった", nameStr, s)
			ms = append(ms, m)
		}

		if v == -1 {
			m := fmt.Sprintf("%sの %sが 下がった", nameStr, s)
			ms = append(ms, m)
		}

		if v == -2 {
			m := fmt.Sprintf("%sの %sが がくんと下がった", nameStr, s)
			ms = append(ms, m)
		}

		if v <= -3 {
			m := fmt.Sprintf("%sの %sが がくーんと下がった", nameStr, s)
			ms = append(ms, m)
		}
	}

	add(v.Atk, "攻撃")
	add(v.Def, "防御")
	add(v.SpAtk, "特攻")
	add(v.SpDef, "特防")
	add(v.Speed, "素早さ")
	return ms
}

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
// 	HostTrainerTitle bp.TrainerTitle
// 	HostTrainerName bp.TrainerName
// 	GuestTrainerTitle bp.TrainerTitle
// 	GuestTrainerName bp.TrainerName
// 	IsHostView bool
// }

// func NewMessageMaker(m *Manager) MessageMaker {
// 	return MessageMaker{
// 		TrainerTitle:m.OpponentTrainerTitle,
// 		TrainerName:m.OpponentTrainerName,
// 		IsHostView:m.IsHostView,
// 	}
// }

// func (mm *MessageMaker) GetCurrentSelfTrainerName() TrainerName {
// 	if mm.IsHostView {
// 		return mm.HostTrainerName
// 	} else {
// 		return mm.GuestTrainerName
// 	}
// }

// func (mm *MessageMaker) SelfFullName() (Message, Message) {
// 	return Message(mm.OpponentTrainerName)
// }

// func (mm *MessageMaker) OpponentFullName() (Message, Message) {
// 	return Message(mm.OpponentTrainerTitle) + "の " + Message(mm.OpponentTrainerName)
// }

// func (mm *MessageMaker) OpponentTrainerNamePrefix() Message {
// 	if mm.IsHostView {
// 		return ""
// 	} else {
// 		return Message(mm.OpponentTrainerName) + "の "
// 	}
// }

// func (mm *MessageMaker) MoveUse(pokeName bp.PokeName, moveName bp.MoveName, isHost bool) Message {
// 	pokeNameStr := pokeName.ToString()
// 	moveNameStr := moveName.ToString()

// 	var hostTrainerName bp.TrainerName
// 	var guestTrainerName bp.TrainerName

// 	if isHost {
// 		hostPrefix = ""
// 		guestPrefix = mm.HostTrainerName + "の "
// 	} else {
// 		hostPrefix = mm.GuestTrainerName + "の "
// 		guestPrefix = ""
// 	}

// 	host = fmt.Sprintf("%s%sの %s！", hostPrefix, pokeNameStr, moveNameStr)
// 	guest = fmt.Sprintf("%s%sの %s！", guestPrefix, pokeNameStr, moveNameStr)

// 	return Message(host), Message(guest)
// }

// func (mm *MessageMaker) NoEffective(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sには ", pokeName.ToString())
//     return mm.TrainerNamePrefix() + Message(m) + "こうかが ない ようだ..."
// }

// func (mm *MessageMaker) TypeEffective(pokeName bp.PokeName, effective bp.TypeEffective) Message {
// 	switch effective {
// 		case bp.NEUTRAL_EFFECTIVE:
// 			return ""
// 		case bp.SUPER_EFFECTIVE:
// 			m := fmt.Sprintf("%sに 効果は バツグンだ！", pokeName.ToString())
// 			return mm.TrainerNamePrefix() + Message(m)
// 		case bp.NOT_VERY_EFFECTIVE:
// 			m := fmt.Sprintf("%sに 効果は いまひとつだ", pokeName.ToString())
// 			return mm.TrainerNamePrefix() + Message(m)
// 		case bp.NO_EFFECTIVE:
// 			m := fmt.Sprintf("%sには こうかが ない ようだ", pokeName.ToString())
// 			return mm.TrainerNamePrefix() + Message(m)
// 	}
// 	return ""
// }

// func (mm *MessageMaker) Recoil(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは 攻撃の 反動を 受けた", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
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
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Paralysis(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは まひして わざが でにくくなった！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) ParalysisLock(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは しびれてうごけない！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Burn(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは やけどをおった！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) BurnDamage(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは やけどの ダメージを うけている！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Sleep(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ねむってしまった！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) SleepLock(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ぐうぐう ねむっている", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) SleepCured(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m) + NewSleepCuredMessage()
// }

// func (mm *MessageMaker) Freeze(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは こおりついた！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) FreezeLock(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは こおって しまって うごかない！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) FreezeCured(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sの ", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m) + NewFreezeCuredMessage()
// }

// func (mm *MessageMaker) IntimidateAndClearBody(iName, cName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sの クリアボディで %sの いかくは きかなかった！", cName.ToString(), iName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) ChestoBerry(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは カゴのみで ", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m) + NewSleepCuredMessage()
// }

// func (mm *MessageMaker) LumBerry(pokeName bp.PokeName, status bp.StatusAilment) Message {
// 	m := fmt.Sprintf("%sは ラムのみで", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m) + NewStatusAilmentCuredMessage(status)
// }

// func (mm *MessageMaker) TrickRoom(pokeName bp.PokeName, distort bool) Message {
// 	var m string
// 	if distort {
// 		m = fmt.Sprintf("%sは じくうを ゆがめた！", pokeName.ToString())
// 	} else {
// 		m = fmt.Sprintf("%sは じくうを もどした！", pokeName.ToString())
// 	}
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) FollowMe(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ちゅうもくの まとになった！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) BellyDrum(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは たいりょくを けずって パワーぜんかいに なった！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Taunt(pokeName bp.PokeName) Message {
// 	m := fmt.Sprintf("%sは ちょうはつに のってしまった！", pokeName.ToString())
// 	return mm.TrainerNamePrefix() + Message(m)
// }

// func (mm *MessageMaker) Intimidate(src, target bp.PokeName) Message {
// 	if mm.IsSelf {
// 		return Message(fmt.Sprintf("%sの いかくで %sの %sの こうげきが さがった！", src.ToString(), mm.TrainerName, target.ToString()))
// 	} else {
// 		return Message(fmt.Sprintf("%sの %sの いかくで %sの こうげきがさがった！", mm.TrainerName, src.ToString(), target.ToString()))
// 	}
// }

// func (mm *MessageMaker) Init() {

// }