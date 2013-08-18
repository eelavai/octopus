package agent

import (
	"github.com/eelavai/octopus/state"
)

// List of agents that can be manipulated as a group
type List []*Agent

func NewList() List {
	list := make([]*Agent, 0)
	for i := 0; i < state.NodeCount(); i++ {
		agent := NewAgent(uint(i))
		list = append(list, agent)
	}
	return list
}

func (l List) Start() {
	for _, agent := range l {
		agent.Start()
	}
}

func (l List) Dryrun() {
	for _, agent := range l {
		agent.Dryrun()
	}
}

func (l List) Prepare() {
	for _, agent := range l {
		agent.Prepare()
	}
}
