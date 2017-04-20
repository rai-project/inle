package message

type MessageType string

const (
	StatusType            MessageType = "status"
	KernelInfoReplyType   MessageType = "kernel_info_reply"
	KernelInfoRequestType MessageType = "kernel_info_request"
	ConnectReplyType      MessageType = "connect_reply"
	ConnectRequestType    MessageType = "connect_request"
	ExecuteReplyType      MessageType = "execute_reply"
	ExecuteResultType     MessageType = "execute_result"
	ExecuteRequestType    MessageType = "execute_request"
	ShutdownReplyType     MessageType = "shutdown_reply"
	ShutdownRequestType   MessageType = "shutdown_request"
)
