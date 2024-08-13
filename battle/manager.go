package battle

 import (
 	"fmt"
 	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	"golang.org/x/exp/slices"
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

type Manager struct {
	SelfHumanTitle bp.HumanTitle
	SelfHumanName bp.HumanName
	OpponentHumanTitle bp.HumanTitle
	OpponentHumanName bp.HumanName

	SelfLeadPokemons bp.Pokemons
	SelfBenchPokemons  bp.Pokemons
	OpponentLeadPokemons bp.Pokemons
	OpponentBenchPokemons bp.Pokemons

	SelfFollowMePokemonPointers bp.PokemonPointers
	OpponentFollowMePokemonPointers bp.PokemonPointers

	Weather Weather
	RemainingTurn RemainingTurn
	Turn int
	
	IsSingle bool
	IsHostView bool
	HostViewMessage Message
}

func (m *Manager) Init() {
	id := 0
	for i := range m.SelfLeadPokemons {
		m.SelfLeadPokemons[i].Id = id
		id += 1
	}

	for i := range m.SelfBenchPokemons {
		m.SelfBenchPokemons[i].Id = id
		id += 1
	}

	for i := range m.OpponentLeadPokemons {
		m.OpponentLeadPokemons[i].Id = id
		id += 1
	}

	for i := range m.OpponentBenchPokemons {
		m.OpponentBenchPokemons[i].Id = id
		id += 1
	}
}

func (m *Manager) GetLeadPokemons(isSelf bool) bp.Pokemons {
	if isSelf {
		return m.SelfLeadPokemons
	} else {
		return m.OpponentLeadPokemons
	}
}

func (m Manager) Clone() Manager {
	m.SelfLeadPokemons = m.SelfLeadPokemons.Clone()
	m.SelfBenchPokemons = m.SelfBenchPokemons.Clone()
	m.OpponentLeadPokemons = m.OpponentLeadPokemons.Clone()
	m.OpponentBenchPokemons = m.OpponentBenchPokemons.Clone()
	m.SelfFollowMePokemonPointers = slices.Clone(m.SelfFollowMePokemonPointers)
	m.OpponentFollowMePokemonPointers = slices.Clone(m.OpponentFollowMePokemonPointers)
	return m
}

func (m *Manager) IsTrickRoomState() bool {
	return m.RemainingTurn.TrickRoom > 0
}

func (m *Manager) SwapView() {
	m.SelfLeadPokemons, m.SelfBenchPokemons, m.OpponentLeadPokemons, m.OpponentBenchPokemons =
		m.OpponentLeadPokemons, m.OpponentBenchPokemons, m.SelfLeadPokemons, m.SelfBenchPokemons
	m.SelfFollowMePokemonPointers, m.OpponentFollowMePokemonPointers = m.OpponentFollowMePokemonPointers, m.SelfFollowMePokemonPointers
}

//攻撃する側のポケモンは瀕死ではない事が前提で呼び出す関数
func (m *Manager) TargetPokemonPointers(action *SoloAction) bp.PokemonPointers {
	moveData := bp.MOVEDEX[action.MoveName]

	single := func() bp.PokemonPointers {
		switch moveData.Target {
			case bp.NORMAL_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.OPPONENT_TWO_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.SELF_TARGET:
				return m.SelfLeadPokemons.ToPointers().NotFainted()
			case bp.OTHERS_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.ALL_TARGET:
				return bp.PokemonPointers{}
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
		}
		return bp.PokemonPointers{}
	}

	doubleNormalTarget := func() bp.PokemonPointers {
		followMe := m.OpponentFollowMePokemonPointers.NotFainted()
		if len(followMe) != 0 {
			return bp.PokemonPointers{followMe[0]}
		}

		//味方への攻撃
		if action.IsSelfLeadTarget {
			ally := m.SelfLeadPokemons[action.TargetIndex]
			//攻撃対象の味方が瀕死ならば、ランダムに相手を攻撃する。
			if ally.IsFainted() {
				ps:= m.OpponentLeadPokemons.ToPointers().NotFainted()
				if len(ps) == 0 {
					return bp.PokemonPointers{}
				}
				return omwrand.Sample(ps, 1, GlobalContext.Rand)
			} else {
				return bp.PokemonPointers{&m.SelfLeadPokemons[action.TargetIndex]}
			}
		} else {
			target := m.OpponentLeadPokemons[action.TargetIndex]
			if target.IsFainted() {
				otherI := map[int]int{0:1, 1:0}[action.TargetIndex]
				other := m.OpponentLeadPokemons[otherI]
				if other.IsFainted() {
					return bp.PokemonPointers{}
				} else {
					return bp.PokemonPointers{&other}
				}						
			} else {
				return bp.PokemonPointers{&target}
			}
		}
	}

	doubleOpponentRandomOneTarget := func() bp.PokemonPointers {
		ps := m.OpponentLeadPokemons.ToPointers().NotFainted()
		if len(ps) == 0 {
			return bp.PokemonPointers{}
		} else {
			return omwrand.Sample(ps, 1, GlobalContext.Rand)
		}
	}

	double := func() bp.PokemonPointers {
		switch moveData.Target {
			case bp.NORMAL_TARGET:
				return doubleNormalTarget()
			case bp.OPPONENT_TWO_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.SELF_TARGET:
				return bp.PokemonPointers{&m.SelfLeadPokemons[action.TargetIndex]}
			case bp.OTHERS_TARGET:
				allyI := map[int]int{0:1, 1:0}[action.SrcIndex]
				ally := m.SelfLeadPokemons[allyI]
				return bp.PokemonPointers{&ally, &m.OpponentLeadPokemons[0], &m.OpponentLeadPokemons[1]}.NotFainted()
			case bp.ALL_TARGET:
				return bp.PokemonPointers{}
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return doubleOpponentRandomOneTarget()
		}
		return bp.PokemonPointers{}
	}

	if m.IsSingle {
		return single()
	} else {
		return double()
	}
}

// いかく
// https://wiki.xn--rckteqa2e.com/wiki/%E3%81%84%E3%81%8B%E3%81%8F
// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#%E3%82%BF%E3%83%BC%E3%83%B3%E3%81%AE%E8%A9%B3%E7%B4%B0

func (m *Manager) Switch(leadIdx, benchIdx int) error {
	if m.SelfBenchPokemons[benchIdx].IsFainted() {
		name := m.SelfBenchPokemons[benchIdx].Name
		msg := fmt.Sprintf("%d番目の %sに 交代しようとしたが、瀕死状態である為、交代出来ません。", benchIdx, name.ToString())
		return fmt.Errorf(msg)
	}

	mm := MessageMaker{IsSelf:m.IsHostView}
	m.HostViewMessage = mm.Back(m.SelfLeadPokemons[leadIdx].Name)
	GlobalContext.Observer(m, MESSAGE_EVENT)
	m.SelfLeadPokemons[leadIdx], m.SelfBenchPokemons[benchIdx] = m.SelfBenchPokemons[benchIdx], m.SelfLeadPokemons[leadIdx]
	m.HostViewMessage = mm.Go(m.SelfLeadPokemons[leadIdx].Name)
	GlobalContext.Observer(m, MESSAGE_EVENT)

	if m.SelfLeadPokemons[leadIdx].Ability == bp.INTIMIDATE {
		for i := range m.OpponentLeadPokemons {
			p := m.OpponentLeadPokemons[i]
			if p.Rank.Atk == bp.MIN_RANK {
				continue
			}

			if p.Ability == bp.CLEAR_BODY {
				m.HostViewMessage = Message(fmt.Sprintf("%sの クリアボディで %sの いかくは きかなかった！", p.Name.ToString(), m.SelfLeadPokemons[leadIdx].Name.ToString()))
				GlobalContext.Observer(m, MESSAGE_EVENT)
				continue
			}

			m.HostViewMessage = Message(fmt.Sprintf("%sの いかくで %sの こうげきが さがった！", m.SelfLeadPokemons[leadIdx].Name.ToString(), p.Name.ToString()))
			GlobalContext.Observer(m, MESSAGE_EVENT)

			if p.Ability != bp.CLEAR_BODY && p.Rank.Atk != bp.MIN_RANK {
				m.OpponentLeadPokemons[i].Rank.Atk -= 1
			}
		}
	}
	return nil
}

// https://latest.pokewiki.net/%E3%83%90%E3%83%88%E3%83%AB%E4%B8%AD%E3%81%AE%E5%87%A6%E7%90%86%E3%81%AE%E9%A0%86%E7%95%AA
// https://wiki.xn--rckteqa2e.com/wiki/%E3%82%BF%E3%83%BC%E3%83%B3#1.%E3%83%9D%E3%82%B1%E3%83%A2%E3%83%B3%E3%82%92%E7%B9%B0%E3%82%8A%E5%87%BA%E3%81%99

func (m *Manager) TurnEnd() error {
	return nil
}

func (m *Manager) ToEasyRead() EasyReadManager {
	return EasyReadManager{
		SelfLeadPokemons:m.SelfLeadPokemons.ToEasyRead(),
		SelfBenchPokemons:m.SelfBenchPokemons.ToEasyRead(),

		OpponentLeadPokemons:m.OpponentLeadPokemons.ToEasyRead(),
		OpponentBenchPokemons:m.OpponentBenchPokemons.ToEasyRead(),

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