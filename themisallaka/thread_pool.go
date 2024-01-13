package themisallaka

import (
	"fmt"
	"github.com/thaianhsoft/gothemis/themiscore/container"
	"sync"
)

type pool struct {
	threadSize  int
	q           container.IteratorContainer
	cond_locker *sync.Mutex
	cond        *sync.Cond
	serviceName string
	typeAction  typeAction
}

func newPool(threadSize int, serviceName string, typeAction typeAction) *pool {
	p := &pool{
		threadSize:  threadSize,
		q:           container.NewQueue(),
		cond_locker: &sync.Mutex{},
		serviceName: serviceName,
		typeAction:  typeAction,
	}

	p.cond = sync.NewCond(p.cond_locker)
	waitgroup := &sync.WaitGroup{}
	for i := 0; i < threadSize; i++ {
		waitgroup.Add(1)
		go p.spawnThread(uint32(i+1), waitgroup)
	}

	waitgroup.Wait()
	fmt.Println("all thread spawned")
	//fmt.Printf("pool: %p\n", p)
	return p
}

func (p *pool) spawnThread(threadId uint32, waitGroup *sync.WaitGroup) {
	//fmt.Printf("[Service=%v], spawn thread id=[%v] runnable!!!\n", p.serviceName, threadId)
	var initBehaviour ThreadBehaviour = &RunnableBehaviour{
		pool:     p,
		threadId: threadId,
	}
	firstBorn := true
	//fmt.Printf("from runnable pool copy address: %p\n", initBehaviour.(*RunnableBehaviour).pool)
	for {
		if firstBorn {
			waitGroup.Done()
			firstBorn = false
		}
		initBehaviour = initBehaviour.Process()
	}
}

func (p *pool) signalWakeUpPoolWithJob(job any) {
	// function used from goroutine outside module pool
	p.cond_locker.Lock()
	p.q.Push(false, job)
	p.cond_locker.Unlock()
	p.cond.Signal() // send signal to wake up have job on queue
}
