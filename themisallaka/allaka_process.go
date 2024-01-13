package themisallaka

import (
	"fmt"
	"strconv"
	"time"
)

type process uint64
type Pid string

var signatureThemisThread = "themis_pid#"

func (p process) toPid() Pid {
	if p > 0 {
		return Pid(signatureThemisThread + fmt.Sprintf("%v", p))
	}
	return ""
}

func (p Pid) toProcess() process {
	if len(p) > len(signatureThemisThread) {
		v := p[len(signatureThemisThread):]
		processUint, err := strconv.Atoi(string(v))
		if err != nil {
			return 0
		} else {
			return process(processUint)
		}
	}
	return 0
}

type ThreadBehaviour interface {
	Process(processFunc ThreadProcessFunc) ThreadBehaviour
	GetState() string
}

type ThreadProcessFunc interface {
	RunnableProcessFunc(behaviour *RunnableBehaviour) ThreadBehaviour
	RunningProcessFunc(behaviour *RunningBehaviour) ThreadBehaviour
	WaitCoodinateBehaviour(behaviour *WaitCoodinateBehaviour) ThreadBehaviour
}

type RunnableBehaviour struct {
	threadId uint32
}

func (r *RunnableBehaviour) GetState() string {
	return fmt.Sprintf("THREAD [ID=%v] is runnable", r.GetThreadId())
}

func (r *RunnableBehaviour) Process(processFunc ThreadProcessFunc) ThreadBehaviour {
	return processFunc.RunnableProcessFunc(r)
}

func (r *RunnableBehaviour) GetThreadId() uint32 {
	return r.threadId
}

type RunningBehaviour struct {
	threadId uint32
	Job      any
}

func (r *RunningBehaviour) SetThreadId(threadId uint32) {
	r.threadId = threadId
}

func (r *RunningBehaviour) GetState() string {
	return fmt.Sprintf("THREAD [ID=%v] is running with job %v\n", r.GetId(), r.Job)
}

func (r *RunningBehaviour) Process(processFunc ThreadProcessFunc) ThreadBehaviour {
	return processFunc.RunningProcessFunc(r)
}

func (r *RunningBehaviour) GetId() uint32 {
	return r.threadId
}

type WaitCoodinateBehaviour struct {
	Job             any
	threadId        uint32
	DurationTimeout *time.Duration
}

func (r *WaitCoodinateBehaviour) GetState() string {
	return fmt.Sprintf("THREAD [ID=%v] is waiting for coodinating from other", r.GetId())
}

func (r *WaitCoodinateBehaviour) Process(processFunc ThreadProcessFunc) ThreadBehaviour {
	return processFunc.WaitCoodinateBehaviour(r)
}

func (r *WaitCoodinateBehaviour) GetId() uint32 {
	return r.threadId
}
