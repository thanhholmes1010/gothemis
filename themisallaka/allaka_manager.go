package themisallaka

import "sync"

type AllakaManager struct {
	*AllakaContext        // itself also be context [root-context]
	globalProcessId   Pid // 0
	container         map[Process]*AllakaContext
	supervisorTree    map[Process][]Process
	dispatcherManager *pool
	executorManager   *pool
	managerLocker     *sync.Mutex
}

var baseAllakaManager *AllakaManager

func DefaultAllakaManager() *AllakaManager {
	if baseAllakaManager == nil {
		baseAllakaManager = &AllakaManager{
			container:         make(map[Process]*AllakaContext),
			supervisorTree:    make(map[Process][]Process),
			managerLocker:     &sync.Mutex{},
			dispatcherManager: newPool(100000, "Dispatcher", sendAction),
			executorManager:   newPool(100000, "Executor", executeAction),
		}
	}

	new_pid := baseAllakaManager.globalProcessId + 1
	new_process := new_pid.ToProcess()
	baseAllakaManager.AllakaContext = newAllakaContext(baseAllakaManager.globalProcessId+1, nil)
	baseAllakaManager.globalProcessId++
	baseAllakaManager.container[new_process] = baseAllakaManager.AllakaContext
	return baseAllakaManager
}

func (a *AllakaManager) registerSupervisionParentWithChilds(parent_process Process, child_process Process) {
	if _, ok := a.supervisorTree[parent_process]; !ok {
		a.supervisorTree[parent_process] = make([]Process, 0)
	}
	a.supervisorTree[parent_process] = append(a.supervisorTree[parent_process], child_process)
}
