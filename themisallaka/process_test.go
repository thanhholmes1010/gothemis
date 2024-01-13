package themisallaka

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type ProductActorRuntime struct {
	mutex     *sync.Mutex
	UpdateSum int
}

func (p *ProductActorRuntime) Receive(ctx *AllakaContext) {
	_, data := ctx.GetMailLatest()
	switch msg := data.(type) {
	case string:
		if msg == "update_sum" {
			p.UpdateSum += 1
			//fmt.Println("updating sum: ", p.UpdateSum)
		}

		if msg == "print_sum" {
			fmt.Println("update success all: ", p.UpdateSum)
		}
		if msg == "delay" {
			time.Sleep(20 * time.Second)
			fmt.Println("delaying...")
		}
	}
}

func TestPid(t *testing.T) {
	runtime.GOMAXPROCS(8)
	manager := DefaultAllakaManager()
	product_actor := manager.SpawnActor(&ProductActorRuntime{
		mutex: &sync.Mutex{},
	})
	for i := 0; i < 100000000; i++ {
		manager.Send(product_actor.GetProcess(), "update_sum")
	}

	manager.Send(product_actor.GetProcess(), "print_sum")
	select {}
}
