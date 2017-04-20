package kernel

import (
	"runtime"

	"github.com/rai-project/config"
	"github.com/rai-project/inle/pkg/message"
	"github.com/rai-project/inle/pkg/protocol"
)

// KernelInfo holds information about the igo kernel, for kernel_info_reply messages.
type KernelInfo struct {
	// Version of messaging protocol.
	// The first integer indicates major version.  It is incremented when
	// there is any backward incompatible change.
	// The second integer indicates minor version.  It is incremented when
	// there is any backward compatible change.
	ProtocolVersion string `json:"protocol_version"`
	// The kernel implementation name
	// (e.g. 'ipython' for the IPython kernel)
	Implementation string `json:"implementation"`
	// Implementation version number.
	// The version number of the kernel's implementation
	//(e.g. IPython.__version__ for the IPython kernel)
	ImplementationVersion string `json:"implementation_version"`
	// Information about the language of code for the kernel
	LanguageInfo KernelLanguageInfo `json:"language_info"`
	// A banner of information about the kernel,
	// which may be desplayed in console environments.
	Banner string `json:"banner"`
}

type KernelLanguageInfo struct {
	// Name of the programming language that the kernel implements.
	// Kernel included in IPython returns 'python'.
	Name string `json:"name"`
	// Language version number.
	// It is Python version number (e.g., '2.7.3') for the kernel
	// included in IPython.
	Version string `json:"version"`
	// mimetype for script files in this language
	Mimetype string `json:"mimetype"`
	// Extension including the dot, e.g. '.py'
	FileExtension string `json:"file_extension"`
	//Pygments lexer, for highlighting
	// Only needed if it differs from the 'name' field.
	PpygmentsLexer string `json:"pygments_lexer,omitempty"`
	// Codemirror mode, for for highlighting in the notebook.
	// Only needed if it differs from the 'name' field.
	CodeMirrorMode string `json:"codemirror_mode,omitempty"`
}

// SendKernelInfo sends a kernel_info_reply message.
func (k *Kernel) HandleKernelInfoRequest(receipt message.Receipt) error {
	reply, err := message.New(message.KernelInfoReplyType, receipt.Message)
	if err != nil {
		return err
	}

	reply.Content = KernelInfo{
		ProtocolVersion:       protocol.Version,
		Implementation:        config.App.Name,
		ImplementationVersion: config.App.Version.Version,
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
