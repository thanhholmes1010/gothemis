package themisallaka

import (
	"fmt"
)

// algebraic data types
type typeAction uint8

const (
	sendAction typeAction = iota
	executeAction
)

type ThreadBehaviour interface {
	Log(serviceName string, message string)
	Process() ThreadBehaviour
}

type RunnableBehaviour struct {
	threadId uint32
	pool     *pool
}

func (r *RunnableBehaviour) Log(serviceName string, message string) {
	fmt.Printf("[%v, ThreadId=%v]: state=[%v], message: %v\n", serviceName, r.threadId, "RUNNABLE", message)
}

func (r *RunnableBehaviour) Process() ThreadBehaviour {
	r.pool.cond_locker.Lock()
	for r.pool.q.Empty() {
		//r.Log(r.pool.serviceName, "No Job is Sleeping !!")
		r.pool.cond.Wait()
	}
	//r.Log(r.pool.serviceName, "Wake Up on having one job to process")
	popJob := r.pool.q.Pop().(*baseMailMessage)
	//r.Log(r.pool.serviceName, fmt.Sprintf("pop job %v", popJob))
	r.pool.cond_locker.Unlock()
	return &RunningBehaviour{
		jobPrepareToRun: popJob,
		threadId:        r.threadId,
		pool:            r.pool,
	}
}

type RunningBehaviour struct {
	threadId        uint32
	jobPrepareToRun *baseMailMessage
	pool            *pool
}

func (r *RunningBehaviour) Log(serviceName string, message string) {
	fmt.Printf("[%v, ThreadId=%v]: state=[%v], message: %v\n", serviceName, r.threadId, "RUNNABLE", message)
}

func (r *RunningBehaviour) Process() ThreadBehaviour {
	if r.pool.typeAction == sendAction {
		// when here mean dispatcher, send into executor pool one signal
		// have letter need be executed
		baseAllakaManager.managerLocker.Lock()
		receiverProcess := r.jobPrepareToRun.getReceiver()
		receiverContext := baseAllakaManager.container[receiverProcess]
		baseAllakaManager.managerLocker.Unlock()
		receiverContext.mutexMailBox.Lock()
		receiverContext.Mailbox.Push(false, r.jobPrepareToRun)
		receiverContext.mutexMailBox.Unlock()
		baseAllakaManager.executorManager.signalWakeUpPoolWithJob(r.jobPrepareToRun)
	}

	if r.pool.typeAction == executeAction {
		// access context with process id
		// process id where ?
		// on letter will have process id of destination
		//baseAllakaManager.managerLocker.Lock()
		receiverProcess := r.jobPrepareToRun.getReceiver()
		receiverContext := baseAllakaManager.container[receiverProcess]
		executeJobMail := receiverContext.Mailbox.Pop().(*baseMailMessage)
		//fmt.Println("executeJobMail: ", executeJobMail)
		receiverContext.newMailLatest = executeJobMail
		//baseAllakaManager.managerLocker.Unlock()

		receiverContext.actor.Receive(receiverContext)
	}

	return &WaitingBehaviour{
		pool:     r.pool,
		threadId: r.threadId,
	}
}

type WaitingBehaviour struct {
	threadId uint32
	pool     *pool
}

func (w *WaitingBehaviour) Log(serviceName string, message string) {
	fmt.Printf("[%v, ThreadId=%v]: state=[%v], message: %v\n", serviceName, w.threadId, "RUNNABLE", message)
}

func (w *WaitingBehaviour) Process() ThreadBehaviour {
	return &RunnableBehaviour{
		pool:     w.pool,
		threadId: w.threadId,
	}
}
