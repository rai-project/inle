package kernel

import (
	"runtime"

	"github.com/rai-project/inle/pkg/message"
	"github.com/rai-project/inle/pkg/protocol"
)

// KernelInfo holds information about the igo kernel, for kernel_info_reply messages.
type KernelInfo struct {
	ProtocolVersion       string             `json:"protocol_version"`
	Implementation        string             `json:"implementation"`
	ImplementationVersion string             `json:"implementation_version"`
	LanguageInfo          KernelLanguageInfo `json:"language_info"`
	Banner                string             `json:"banner"`
}

type KernelLanguageInfo struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	Mimetype      string `json:"mimetype"`
	FileExtension string `json:"file_extension"`
}

// SendKernelInfo sends a kernel_info_reply message.
func (k *Kernel) HandleKernelInfoRequest(receipt message.Receipt) error {
	reply, err := message.New(message.KernelInfoReplyType, receipt.Message)
	if err != nil {
		return err
	}

	reply.Content = KernelInfo{
		ProtocolVersion:       protocol.Version,
		Implementation:        "inle",
		ImplementationVersion: "0.1",
		LanguageInfo: KernelLanguageInfo{
			Name:          "rai",
			Version:       runtime.Version(),
			Mimetype:      "application/x-rai", // text/plain would be possible, too
			FileExtension: ".rai",
		},
		Banner: "RAI-Inle - https://github.com/rai-project/inle",
	}

	return receipt.SendResponse(receipt.Connection.ShellSocket, reply)
}
