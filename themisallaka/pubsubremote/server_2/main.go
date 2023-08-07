package main

import "github.com/thaianhsoft/gothemis/themisallaka"

func main() {
	server := themisallaka.NewAllaSystem()
	server.WaitSignal()
}
