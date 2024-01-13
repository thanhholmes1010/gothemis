package themisallaka

import "testing"

type OrderAllakator struct {
	Id     uint32
	Amount float32
}

type CheckoutOrderEvent struct {
	OrderId   uint32
	Amount    float32
	buyerPid  Pid
	sellerPid Pid
}

func (o *OrderAllakator) Receive(allaProc *AllaKator) {

}

type UserAllakator struct {
	Id      uint32
	Balance float32
}

type UpdateBalanceEvent struct {
	UserId uint32
	Id     uint32
}

func (u *UserAllakator) Receive(allaProc *AllaKator) {
	sender, msg := allaProc.Message()
	switch msg.(type) {
	case *TransactorMessage:
		txMessage := msg.(*TransactorMessage)
		switch v := txMessage.event.(type) {
		case *CheckoutOrderEvent:
			txMessage.AddReadOperation(u)
			txMessage.AddWriteOperation(func(data any) any {

			})
		}
	}
}

func TestAllakator(t *testing.T) {
	manager := NewAllaSystem()
	buyer_1 := manager.SpawnChildProcess(func() Alla {
		return &UserAllakator{
			Id:      1,
			Balance: 50.04,
		}
	})
	seller_1 := manager.SpawnChildProcess(func() Alla {
		return &UserAllakator{
			Id:      2,
			Balance: 10.4,
		}
	})
	checkout := NewTransactMessage(&CheckoutOrderEvent{
		OrderId:   34,
		Amount:    15.04,
		buyerPid:  buyer_1.SelfPid(),
		sellerPid: seller_1.SelfPid(),
	})

}
