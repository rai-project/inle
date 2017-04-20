package kernel

import "github.com/rai-project/inle/pkg/message"

// ConnectReply encodes the ports necessary for connecting to the kernel
type ConnectReply struct {
	ShellPort int `json:"shell_port"`
	IOPubPort int `json:"iopub_port"`
	StdinPort int `json:"stdin_port"`
	HBPort    int `json:"hb_port"`
}

func (k *Kernel) HandleConnectRequest(receipt message.Receipt) error {
	reply, err := message.New(message.ConnectReplyType, receipt.Message)
	if err != nil {
		return err
	}

	connectionInfo := k.connectionOptions

	reply.Content = ConnectReply{
		ShellPort: connectionInfo.ShellPort,
		IOPubPort: connectionInfo.IOPubPort,
		StdinPort: connectionInfo.StdinPort,
		HBPort:    connectionInfo.HeartBeatPort,
	}

	return receipt.SendResponse(receipt.Connection.ShellSocket, reply)
}
