package themisallaka

import (
	"sync/atomic"
)

var transactionId uint64 = 0

type TxOperation uint8

const (
	WriteOperation TxOperation = iota + 1
	ReadOperation
)

type TransactorMessage struct {
	Id               uint64
	MessageIds       []uint64
	PersistedSignals []bool
	event            any
	txAllakator      *AllaKator
}

func (m *TransactorMessage) AddReadOperation(allakator *AllaKator) {
	m.txAllakator = allakator
}

func (m *TransactorMessage) AddWriteOperation(changeWriteFunc func()) {
	if _, ok := defaulTransactManager.casFreeLockWriter[m.txAllakator]; !ok {
		defaulTransactManager.casFreeLockWriter[m.txAllakator] = true
	}
}

func NewTransactMessage(event any) *TransactorMessage {
	atomic.AddUint64(&transactionId, 1)
	return &TransactorMessage{
		Id:    atomic.LoadUint64(&transactionId),
		event: event,
	}
}

var defaulTransactManager *transactManager = &transactManager{
	casFreeLockWriter: make(map[*AllaKator]bool),
}

type transactManager struct {
	casFreeLockWriter map[*AllaKator]bool
}

func (m *transactManager) Commit(txMessage *TransactorMessage) {

}
