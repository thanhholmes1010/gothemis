package main

import "github.com/thaianhsoft/gothemis/themisgraphorm/generator"

func main() {
	generator := generator.GraphEngineGeneratorImlp{}
	generator.Run()
	//Usage:
	//@Step1: type go run main.go
	//@Step2: command test: gen schema Student Id:UInt Name:Varchar(40) Age: UInt
}
