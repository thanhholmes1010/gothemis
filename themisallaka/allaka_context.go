package themisallaka

import (
	"github.com/thaianhsoft/gothemis/themiscore/container"
	"sync"
)

type AllakaContext struct {
	Mailbox       container.IteratorContainer // queue
	pid           Pid
	actor         Actor
	newMailLatest *baseMailMessage
	mutexMailBox  *sync.Mutex
}

func newAllakaContext(pid Pid, actor Actor) *AllakaContext {
	return &AllakaContext{
		Mailbox:      container.NewQueue(),
		pid:          pid,
		actor:        actor,
		mutexMailBox: &sync.Mutex{},
	}
}

func (a *AllakaContext) GetProcess() Process {
	return a.pid.ToProcess()
}

func (a *AllakaContext) SpawnActor(actor Actor) *AllakaContext {
	child_pid := baseAllakaManager.globalProcessId + 1
	child_process := child_pid.ToProcess()
	child_context := newAllakaContext(child_pid, actor)
	baseAllakaManager.globalProcessId++
	baseAllakaManager.container[child_process] = child_context
	baseAllakaManager.registerSupervisionParentWithChilds(a.pid.ToProcess(), child_process)
	return child_context
}

func (a *AllakaContext) Send(receiverProcess Process, data any) {
	base_mail_message := &baseMailMessage{
		receiver: receiverProcess,
		sender:   a.pid.ToProcess(),
		data:     data,
	}
	baseAllakaManager.dispatcherManager.signalWakeUpPoolWithJob(base_mail_message)
}

func (a *AllakaContext) GetMailLatest() (senderProcess Process, data any) {
	return a.newMailLatest.sender, a.newMailLatest.data
}
