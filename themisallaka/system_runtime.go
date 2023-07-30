package themisallaka

import (
	"github.com/thaianhsoft/gothemis/themiscore/container"
)

var baseAllaSystem *AllaSystem

type AllaSystem struct {
	*AllaKator
	containerAllakators *container.SplayTree
	globalProcessId     process
	pool                *Pool
}

func NewAllaSystem() *AllaSystem {
	if baseAllaSystem != nil {
		return baseAllaSystem
	}
	baseAllaSystem = &AllaSystem{
		containerAllakators: &container.SplayTree{},
		globalProcessId:     1,
	}
	baseAllakator := newAllaKator(baseAllaSystem.globalProcessId, nil)
	baseAllaSystem.AllaKator = baseAllakator
	baseAllaSystem.containerAllakators.Insert(uint64(baseAllaSystem.globalProcessId), baseAllakator)
	return baseAllaSystem
}

func (as *AllaSystem) nextProcessId() process {
	v := as.globalProcessId + 1
	as.globalProcessId++
	return v
}

func (as *AllaSystem) Run(port int) {

}
