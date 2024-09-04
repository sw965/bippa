package battle

 import (
 	"fmt"
 	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	"golang.org/x/exp/slices"
	omwslices "github.com/sw965/omw/slices"
	omwmath "github.com/sw965/omw/math"
)

const (
	PLAYER_NUM = 2
	LEAD_NUM = 1
	BENCH_NUM = 2
	FIGHTERS_NUM = LEAD_NUM + BENCH_NUM
	DOUBLE_BATTLE_NUM = 2
)

type RemainingTurn struct {
	Weather int
	TrickRoom int
}

func (t *RemainingTurn) IsTrickRoomState() bool {
	return t.TrickRoom > 0
}

type Manager struct {
	CurrentSelfLeadPokemons bp.Pokemons
	CurrentSelfBenchPokemons  bp.Pokemons
	CurrentOpponentLeadPokemons bp.Pokemons
	CurrentOpponentBenchPokemons bp.Pokemons

	CurrentSelfFollowMePokemonPointers bp.PokemonPointers
	CurrentOpponentFollowMePokemonPointers bp.PokemonPointers

	Weather Weather
	RemainingTurn RemainingTurn
	Turn int

	CurrentSelfIsHost bool
	HostViewMessage string

	getHumanNameMessage func(bool) string
	getHumanInfoMessage func(bool) string
}

func (m *Manager) Init(guestHumanTitle, guestHumanName string) {
	for i := range m.CurrentSelfLeadPokemons {
		m.CurrentSelfLeadPokemons[i].IsHost = true
	}

	for i := range m.CurrentSelfBenchPokemons {
		m.CurrentSelfBenchPokemons[i].IsHost = true
	}
	m.CurrentSelfIsHost = true
	m.getHumanInfoMessage = GetHumanInfoMessageFunc(guestHumanTitle, guestHumanName)
	m.getHumanNameMessage = GetHumanNameMessageFunc(guestHumanName)
}

func (m *Manager) IsSingle() bool {
	return len(m.CurrentSelfLeadPokemons) == 1
}

func (m Manager) Clone() Manager {
	m.CurrentSelfLeadPokemons = m.CurrentSelfLeadPokemons.Clone()
	m.CurrentSelfBenchPokemons = m.CurrentSelfBenchPokemons.Clone()
	m.CurrentOpponentLeadPokemons = m.CurrentOpponentLeadPokemons.Clone()
	m.CurrentOpponentBenchPokemons = m.CurrentOpponentBenchPokemons.Clone()
	m.CurrentSelfFollowMePokemonPointers = slices.Clone(m.CurrentSelfFollowMePokemonPointers)
	m.CurrentOpponentFollowMePokemonPointers = slices.Clone(m.CurrentOpponentFollowMePokemonPointers)
	return m
}

func (m *Manager) SwapView() {
	m.CurrentSelfLeadPokemons, m.CurrentSelfBenchPokemons, m.CurrentOpponentLeadPokemons, m.CurrentOpponentBenchPokemons =
		m.CurrentOpponentLeadPokemons, m.CurrentOpponentBenchPokemons, m.CurrentSelfLeadPokemons, m.CurrentSelfBenchPokemons
	m.CurrentSelfFollowMePokemonPointers, m.CurrentOpponentFollowMePokemonPointers = m.CurrentOpponentFollowMePokemonPointers, m.CurrentSelfFollowMePokemonPointers
	m.CurrentSelfIsHost = !m.CurrentSelfIsHost
}

func (m *Manager) IsTrickRoomState() bool {
	return m.RemainingTurn.IsTrickRoomState()
}

func (m *Manager) ApplyDamageToBody(p *bp.Pokemon, dmg int) error {
	if dmg < 0 {
		return fmt.Errorf("ダメージは0以上でなければならない")
	}
	dmg = omwmath.Min(dmg, p.Stat.CurrentHP)
	p.Stat.CurrentHP -= dmg

	// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%AA%E3%83%9C%E3%83%B3%E3%81%AE%E3%81%BF
	halfMaxHP := int(float64(p.Stat.MaxHP) * 0.5)
	sitrusBerryOK := p.Item == bp.SITRUS_BERRY && halfMaxHP >= p.Stat.CurrentHP
	if sitrusBerryOK {
		heal := int(float64(p.Stat.MaxHP) * 0.5)
		err := p.ApplyHealToBody(heal)
		if err != nil {
			return err
		}
		m.HostViewMessage = m.getHumanNameMessage(p.IsHost) + p.Name.ToString() + "は オボンの実で 回復した！"
	}
	return nil
}

func (m *Manager) targetPokemonPointersForSingle(a *SoloAction) (bp.PokemonPointers, error) {
	moveData := bp.MOVEDEX[a.MoveName]
	notFainted := m.CurrentOpponentLeadPokemons.ToPointers().FilterByNotFainted()
	switch moveData.Target {
		case bp.NORMAL_TARGET:
			return notFainted, nil
		case bp.OPPONENT_TWO_TARGET:
			return notFainted, nil
		// 攻撃するポケモンは、瀕死ではない事が前提なので、瀕死のチェックはしない。
		case bp.SELF_TARGET:
			if a.SrcIndex != a.TargetIndex {
				m := fmt.Sprintf("自分自身を対象とする技を使う場合、SoloAction.SrcIndex == SoloAction.TargetIndex でなければならない。")
				return bp.PokemonPointers{}, fmt.Errorf(m)
			}
			return bp.PokemonPointers{&m.CurrentSelfLeadPokemons[0]}, nil
		case bp.OTHERS_TARGET:
			return notFainted, nil
		case bp.OPPONENT_RANDOM_ONE_TARGET:
			return notFainted, nil
	}
	// ALL_TARGET の場合
	return bp.PokemonPointers{}, nil
}

func (m *Manager) targetPokemonPointersForDoubleNormalTarget(a *SoloAction) bp.PokemonPointers {
	followMe := m.CurrentOpponentFollowMePokemonPointers.FilterByNotFainted()
	if len(followMe) != 0 {
		// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%93%E3%81%AE%E3%82%86%E3%81%B3%E3%81%A8%E3%81%BE%E3%82%8C
		// 複数ポケモンがちゅうもくのまと状態である場合、最初のポケモンが優先される。
		return bp.PokemonPointers{followMe[0]}
	}

	//味方への攻撃
	if a.IsSelfLeadTarget {
		ally := m.CurrentSelfLeadPokemons[a.TargetIndex]
		// https://wiki.xn--rckteqa2e.com/wiki/%E3%83%80%E3%83%96%E3%83%AB%E3%83%90%E3%83%88%E3%83%AB
		//攻撃対象の味方が既に瀕死ならば、ランダムに相手を攻撃する。(第4世代)
		if ally.IsFainted() {
			ps := m.CurrentOpponentLeadPokemons.ToPointers().FilterByNotFainted()
			if len(ps) == 0 {
				return bp.PokemonPointers{}
			}
			return omwrand.Sample(ps, 1, GlobalContext.Rand)
		}
		return bp.PokemonPointers{&m.CurrentSelfLeadPokemons[a.TargetIndex]}
	}

	//相手への攻撃
	target := m.CurrentOpponentLeadPokemons[a.TargetIndex]
	if target.IsFainted() {
		otherI := map[int]int{0:1, 1:0}[a.TargetIndex]
		//攻撃対象の相手が既に瀕死ならば、もう片方のポケモンを攻撃対象にする。
		other := m.CurrentOpponentLeadPokemons[otherI]
		if other.IsFainted() {
			return bp.PokemonPointers{}
		}
		return bp.PokemonPointers{&other}						
	}
	return bp.PokemonPointers{&target}
}

func (m *Manager) targetPokemonPointersForDoubleRandomOneTarget() bp.PokemonPointers {
	ps := m.CurrentOpponentLeadPokemons.ToPointers().FilterByNotFainted()
	if len(ps) == 0 {
		return bp.PokemonPointers{}
	}
	return omwrand.Sample(ps, 1, GlobalContext.Rand)
}

func (m *Manager) targetPokemonPointersForDoubleOthersTarget(a *SoloAction) bp.PokemonPointers {
	allyI := map[int]int{0:1, 1:0}[a.SrcIndex]
	ally := m.CurrentSelfLeadPokemons[allyI]
	left := m.CurrentOpponentLeadPokemons[0]
	right := m.CurrentOpponentLeadPokemons[1]
	return bp.PokemonPointers{&ally, &left, &right}.FilterByNotFainted()
}

func (m *Manager) targetPokemonPointersForDouble(a *SoloAction) (bp.PokemonPointers, error) {
	moveData := bp.MOVEDEX[a.MoveName]
	switch moveData.Target {
		case bp.NORMAL_TARGET:
			return m.targetPokemonPointersForDoubleNormalTarget(a), nil
		case bp.OPPONENT_TWO_TARGET:
			return m.CurrentOpponentLeadPokemons.ToPointers().FilterByNotFainted(), nil
		// 攻撃するポケモンは、瀕死ではない事が前提なので、瀕死のチェックはしない。
		case bp.SELF_TARGET:
			if a.SrcIndex != a.TargetIndex {
				m := fmt.Sprintf("自分自身を対象とする技を使う場合、SoloAction.SrcIndex == SoloAction.TargetIndex でなければならない。")
				return bp.PokemonPointers{}, fmt.Errorf(m)
			}
			return bp.PokemonPointers{&m.CurrentSelfLeadPokemons[a.TargetIndex]}, nil
		case bp.OTHERS_TARGET:
			return m.targetPokemonPointersForDoubleOthersTarget(a), nil
		case bp.OPPONENT_RANDOM_ONE_TARGET:
			return m.targetPokemonPointersForDoubleRandomOneTarget(), nil
	}
	//ALL_TARGETの場合
	return bp.PokemonPointers{}, nil
}

func (m *Manager) TargetPokemonPointers(a *SoloAction) (bp.PokemonPointers, error) {
	if m.IsSingle() {
		return m.targetPokemonPointersForSingle(a)
	}
	return m.targetPokemonPointersForDouble(a)
}

// いかく
// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%84%E3%81%8B%E3%81%8F
// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#%E3%82%BF%E3%83%BC%E3%83%B3%E3%81%AE%E8%A9%B3%E7%B4%B0

func (m *Manager) Switch(leadIdx, benchIdx int) error {
	beforePokeName := m.CurrentSelfLeadPokemons[benchIdx].Name
	beforePokeNameStr := beforePokeName.ToString()

	if m.CurrentSelfBenchPokemons[benchIdx].IsFainted() {
		msg := fmt.Sprintf("%d番目の %sに 交代しようとしたが、瀕死状態なので、交代出来ません。", benchIdx, beforePokeNameStr)
		return fmt.Errorf(msg)
	}

	if m.CurrentSelfIsHost {
		m.HostViewMessage = "戻れ！ " + beforePokeNameStr
	} else {
		m.HostViewMessage = m.getHumanInfoMessage(false) + beforePokeName.ToString() + " を 引っ込めた！"
	}

	GlobalContext.Observer(m)
	m.CurrentSelfLeadPokemons[leadIdx], m.CurrentSelfBenchPokemons[benchIdx] = m.CurrentSelfBenchPokemons[benchIdx], m.CurrentSelfLeadPokemons[leadIdx]

	afterPokemon := m.CurrentSelfLeadPokemons[leadIdx]
	afterPokeName := afterPokemon.Name
	afterPokeNameStr := afterPokeName.ToString()

	if m.CurrentSelfIsHost {
		m.HostViewMessage = "行け！ " + afterPokeNameStr
	} else {
		m.HostViewMessage =  m.getHumanInfoMessage(false) + afterPokeNameStr + "を 繰り出した！"
	}
	GlobalContext.Observer(m)

	if afterPokemon.Ability == bp.INTIMIDATE {
		intimidateStr := bp.INTIMIDATE.ToString()
		intimidateHumanNameMsg := m.getHumanNameMessage(m.CurrentSelfIsHost)
		targetHumanNameMsg := m.getHumanNameMessage(!m.CurrentSelfIsHost)

		for i := range m.CurrentOpponentLeadPokemons {
			target := &m.CurrentOpponentLeadPokemons[i]
			if target.RankStat.Atk == bp.MIN_RANK {
				continue
			}
			
			targetPokeNameStr := target.Name.ToString()

			if target.Ability == bp.CLEAR_BODY {
				clearbodyStr := bp.CLEAR_BODY.ToString()
				m.HostViewMessage = targetHumanNameMsg + targetPokeNameStr + "の " + clearbodyStr + "で " + intimidateHumanNameMsg + afterPokeNameStr + "の " + intimidateStr + "は きかなかった！"
				GlobalContext.Observer(m)
				continue
			}

			m.CurrentOpponentLeadPokemons[i].RankStat.Atk -= 1
			m.HostViewMessage = intimidateHumanNameMsg + afterPokeNameStr + "の " + intimidateStr + "で " + targetHumanNameMsg + targetPokeNameStr + "の こうげきが さがった！"
			GlobalContext.Observer(m)
		}
	}
	return nil
}

// https://latest.pokewiki.net/%E3%83%90%E3%83%88%E3%83%AB%E4%B8%AD%E3%81%AE%E5%87%A6%E7%90%86%E3%81%AE%E9%A0%86%E7%95%AA
// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#1.%E3%83%9D%E3%82%B1%E3%83%A2%E3%83%B3%E3%82%92%E7%B9%B0%E3%82%8A%E5%87%BA%E3%81%99

func (m *Manager) TurnEnd() error {
	for i := range m.CurrentSelfLeadPokemons {
		m.CurrentSelfBenchPokemons[i].ThisTurnPlannedUseMoveName = bp.EMPTY_MOVE_NAME
	}

	for i := range m.CurrentOpponentLeadPokemons {
		m.CurrentOpponentLeadPokemons[i].ThisTurnPlannedUseMoveName = bp.EMPTY_MOVE_NAME
	}

	leadPokemons := omwslices.Concat(m.CurrentSelfLeadPokemons.ToPointers(), m.CurrentOpponentLeadPokemons.ToPointers())
	slices.SortFunc(leadPokemons, func(p1, p2 *bp.Pokemon) bool {
		if p1.Stat.Speed > p2.Stat.Speed {
			return !m.IsTrickRoomState()
		} else if p1.Stat.Speed < p2.Stat.Speed {
			return m.IsTrickRoomState()
		} else {
			return omwrand.Bool(GlobalContext.Rand)
		}
	})

	if m.RemainingTurn.Weather > 0 {
		m.RemainingTurn.Weather -= 1
	}

	if m.RemainingTurn.Weather == 0 {
		if m.Weather == RAIN {
			m.HostViewMessage = "雨が 降り止んだ！"
		}
		m.Weather = NORMAL_WEATHER
	}

	for _, p := range leadPokemons {
		// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%84%E3%81%91%E3%81%A9
		if p.StatusAilment == bp.BURN {
			dmg := int(float64(p.Stat.MaxHP) * 0.125)
			dmg = omwmath.Max(dmg, 1)
			err := m.ApplyDamageToBody(p, dmg)
			if err != nil {
				return err
			}
			m.HostViewMessage = m.getHumanNameMessage(p.IsHost) + p.Name.ToString() + "は " + bp.BURN.ToString() + " の ダメージを 受けている！"
		}
	}
	return nil
}

func (m *Manager) ToEasyRead() EasyReadManager {
	return EasyReadManager{
		SelfLeadPokemons:m.CurrentSelfLeadPokemons.ToEasyRead(),
		SelfBenchPokemons:m.CurrentSelfBenchPokemons.ToEasyRead(),

		OpponentLeadPokemons:m.CurrentOpponentLeadPokemons.ToEasyRead(),
		OpponentBenchPokemons:m.CurrentOpponentBenchPokemons.ToEasyRead(),

		Turn:m.Turn,
		Weather:m.Weather.ToString(),
		RemainingTurn:m.RemainingTurn,
	}
}

type EasyReadManager struct {
	SelfLeadPokemons bp.EasyReadPokemons
	SelfBenchPokemons bp.EasyReadPokemons

	OpponentLeadPokemons bp.EasyReadPokemons
	OpponentBenchPokemons bp.EasyReadPokemons

	Turn int
	IsHostView bool

	Weather string
	RemainingTurn RemainingTurn
}