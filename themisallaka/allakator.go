package themisallaka

type AllaKator struct {
	processId      process
	baseContext    Alla
	childProcesses []process
}

func newAllaKator(processId process, baseContext Alla) *AllaKator {
	return &AllaKator{
		processId:   processId,
		baseContext: baseContext,
	}
}

func (a *AllaKator) SelfPid() Pid {
	return a.processId.toPid()
}

func (a *AllaKator) SpawnChildProcess(initFuncState func() Alla) *AllaKator {
	processId := baseAllaSystem.nextProcessId()
	newChildAllaktor := newAllaKator(processId, initFuncState())
	a.childProcesses = append(a.childProcesses, processId)
	baseAllaSystem.containerAllakators.Insert(uint64(processId), newChildAllaktor)
	return newChildAllaktor
}

func (a *AllaKator) Send(pid Pid, message any) {

}

func (a *AllaKator) Message() (sender Pid, message any) {
	return "", nil
}
