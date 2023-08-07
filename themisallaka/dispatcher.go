package themisallaka

type Dispatcher struct {
	pool *Pool
}

func (d *Dispatcher) RunnableProcessFunc(behaviour *RunnableBehaviour) ThreadBehaviour {
	d.pool.GetCondLocker().L.Lock()
	for d.pool.GetContainer().Empty() {
		//fmt.Println("[Dispatcher]: Thread is sleeping")
		d.pool.GetCondLocker().Wait()
	}
	tb := &RunningBehaviour{
		Job: d.pool.GetContainer().Pop(),
	}
	d.pool.GetCondLocker().L.Unlock()
	d.pool.GetCondLocker().Signal()
	return tb
}

func (d *Dispatcher) RunningProcessFunc(behaviour *RunningBehaviour) ThreadBehaviour {
	//TODO implement me
	panic("implement me")
}

func (d *Dispatcher) WaitCoodinateBehaviour(behaviour *WaitCoodinateBehaviour) ThreadBehaviour {
	//TODO implement me
	panic("implement me")
}

func NewDispatcher(threadSize int) *Dispatcher {
	d := &Dispatcher{}
	d.pool = NewPool(threadSize, d)
	defer d.pool.spawnThreads()
	return d
}
