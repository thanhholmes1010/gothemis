package themisallaka

import (
	"github.com/thaianhsoft/gothemis/themiscore/container"
	"sync"
)

type Pool struct {
	container       container.IteratorContainer
	condLocker      *sync.Cond
	operationLocker *sync.Mutex
	chansThread     []*chan struct{}
	maxSize         int
	tpFunc          ThreadProcessFunc
}

func NewPool(maxSize int, tpFunc ThreadProcessFunc) *Pool {
	p := &Pool{
		container:       container.NewQueue(),
		condLocker:      sync.NewCond(&sync.Mutex{}),
		operationLocker: &sync.Mutex{},
		maxSize:         maxSize,
		tpFunc:          tpFunc,
	}
	p.spawnThreads()
	return p
}

func (p *Pool) PushJob(data any, highestPriority bool) {
	p.operationLocker.Lock()
	defer func() {
		p.operationLocker.Unlock()
		p.condLocker.Signal()
	}()
	if !highestPriority {
		p.container.Push(false, data)
	} else {
		p.container.Push(true, data)
	}
}

func (p *Pool) PopJob() any {
	p.operationLocker.Lock()
	defer func() {
		p.operationLocker.Unlock()
		p.condLocker.Signal()
	}()
	return p.container.Pop()
}

func (p *Pool) spawnThreads() {
	for i := 0; i < p.maxSize; i++ {
		id := uint32(i)
		go p.runThread(id)
	}
}

func (p *Pool) runThread(threadId uint32) {
	var initBehaviour ThreadBehaviour = &RunnableBehaviour{}
	for {
		initBehaviour = initBehaviour.Process(p.tpFunc)
	}
}

func (p *Pool) GetContainer() container.IteratorContainer {
	return p.container
}

func (p *Pool) GetCondLocker() *sync.Cond {
	return p.condLocker
}

func (p *Pool) GetOperationLocker() *sync.Mutex {
	return p.operationLocker
}
