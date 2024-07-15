package single

import (
	bp "github.com/sw965/bippa"
)

const DOUBLE_BATTLE = 2

type SoloAction struct {
	MoveName bp.MoveName
	SrcIndex int
	TargetIndex int
	Speed int
	IsSelfLeadTarget bool
	IsSelfBenchTarget bool
	IsOpponentLeadTarget bool
	Target bp.TargetRange
}

type SoloActions []SoloAction

type Action [DOUBLE_BATTLE]SoloAction

type Actions []Action

type ActionsSlice []Actions