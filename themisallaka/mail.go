package themisallaka

type baseMailMessage struct {
	messageId uint64
	receiver  Process
	sender    Process
	data      any
}

func (b *baseMailMessage) getReceiver() Process {
	return b.receiver
}

func (b *baseMailMessage) getSender() Process {
	return b.sender
}
