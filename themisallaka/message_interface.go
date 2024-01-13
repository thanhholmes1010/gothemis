package themisallaka

type IMessage interface {
	Id() uint64
	Data() any
}
