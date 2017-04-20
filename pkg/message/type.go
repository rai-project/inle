package message

type MessageType string

const (
	CommCloseType         MessageType = "comm_close"
	CommInfoRequestType   MessageType = "comm_info_request"
	CommMessageType       MessageType = "comm_msg"
	CommOpenType          MessageType = "comm_open"
	CompleteRequestType   MessageType = "complete_request"
	ConnectReplyType      MessageType = "connect_reply"
	ConnectRequestType    MessageType = "connect_request"
	ExecuteReplyType      MessageType = "execute_reply"
	ExecuteRequestType    MessageType = "execute_request"
	ExecuteResultType     MessageType = "execute_result"
	HistoryRequestType    MessageType = "history_request"
	InspectRequestType    MessageType = "inspect_request"
	IsCompleteRequestType MessageType = "is_complete_request"
	KernelInfoReplyType   MessageType = "kernel_info_reply"
	KernelInfoRequestType MessageType = "kernel_info_request"
	ShutdownReplyType     MessageType = "shutdown_reply"
	ShutdownRequestType   MessageType = "shutdown_request"
	StatusType            MessageType = "status"
)
