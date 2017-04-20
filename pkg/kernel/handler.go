package kernel

import "github.com/rai-project/inle/pkg/message"

type HandlerFunction func(message.Receipt) error

func (k *Kernel) HandleWithStatus(receipt message.Receipt, handler HandlerFunction) error {

	// Publish status: busy

	busy, err := message.New(message.StatusType, receipt.Message)
	if err != nil {
		return err
	}
	busy.Content = KernelStatus{"busy"}
	receipt.SendResponse(receipt.Connection.IOPubSocket, busy)

	// Call actual handler function

	handler(receipt)

	// Publish status: idle after processing

	idle, err := message.New(message.StatusType, receipt.Message)
	if err != nil {
		return err
	}
	idle.Content = KernelStatus{"idle"}
	receipt.SendResponse(receipt.Connection.IOPubSocket, idle)

	return nil
}
