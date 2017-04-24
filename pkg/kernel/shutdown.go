package kernel

import "github.com/rai-project/inle/pkg/message"

// ShutdownReply encodes a boolean indication of shutdown/restart
type ShutdownReply struct {
	Restart bool `json:"restart"`
}

// HandleShutdownRequest sends a "shutdown" message
func (k *Kernel) HandleShutdownRequest(receipt message.Receipt) error {
	defer k.Shutdown()
	reply, err := message.New(message.ShutdownReplyType, receipt.Message)
	if err != nil {
		return err
	}
	content := receipt.Message.Content.(map[string]interface{})
	restart := content["restart"].(bool)
	reply.Content = ShutdownReply{restart}
	receipt.SendResponse(receipt.Connection.ShellSocket, reply)
	log.Debug("Shutting down in response to shutdown_request")

	DefaultSessionManager.Flush()

	return nil
}
