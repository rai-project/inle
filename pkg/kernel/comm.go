package kernel

import (
	"github.com/rai-project/inle/pkg/message"
)

func (k *Kernel) HandleCommOpen(receipt message.Receipt) error {
	// pp.Println("in HandleCommOpen", receipt)
	return nil
}

func (k *Kernel) HandleCommMessage(receipt message.Receipt) error {
	// pp.Println("in HandleCommMessage", receipt)
	return nil
}

func (k *Kernel) HandleCommClose(receipt message.Receipt) error {
	// pp.Println("in HandleCommClose", receipt)
	return nil
}

func (k *Kernel) HandleCommInfoRequest(receipt message.Receipt) error {
	// pp.Println("in HandleCommInfoRequest", receipt)
	return nil
}
