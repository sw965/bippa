package battle

 import (
  	"fmt"
  	bp "github.com/sw965/bippa"
)

type DisplayUI struct {
	P1LeadPokeNames []string
	P1LeadLevels []bp.Level
	P1LeadMaxHPs []int
	P1LeadCurrentHPs []int

	P2LeadPokeNames []string
	P2LeadLevels []bp.Level
	P2LeadMaxHPs []int
	P2LeadCurrentHPs []int

	Message string
}

type ObserverUI struct {
	LastP1ViewManager Manager
	LastP2ViewManager Manager
}

func (ui *ObserverUI) LastManager(isPlayer1View bool) Manager {
	if isPlayer1View {
		return ui.LastP1ViewManager
	} else {
		return ui.LastP2ViewManager
	}
}

func (ui *ObserverUI) Observer(current *Manager) {
	//lastManager := ui.LastManager(current.CurrentSelfIsHost)
	fmt.Println(current.HostViewMessage)
}