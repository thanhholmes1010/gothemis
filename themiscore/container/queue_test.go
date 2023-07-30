package container

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := NewQueue()
	data := []int{10, 30, 40, 6, 1, 2, 3}
	for i, one := range data {
		if i == 3 || i == 5 {
			fmt.Println("push front")
			queue.Push(true, one)
		} else {
			fmt.Println("push back")
			queue.Push(false, one)
		}
	}
}
