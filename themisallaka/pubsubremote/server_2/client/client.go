package main

import "github.com/thaianhsoft/gothemis/themisallaka"

type ProductClient struct {
	*themisallaka.PersistEntity
}

func main() {
	clientSystem := themisallaka.NewAllaSystem()
	clientSystem.ConnectRemoteServer("localhost:8000")
	for {

	}
}
