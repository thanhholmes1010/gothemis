package themisallaka

import "fmt"

type Process string
type Pid uint32

var signature = "#Pid.("

func (p Pid) ToProcess() Process {
	return Process(fmt.Sprintf("%v%v)", signature, p))
}
