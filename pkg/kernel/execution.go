package kernel

import (
	"fmt"

	"github.com/rai-project/inle/pkg/message"
)

// OutputMsg holds the data for a pyout message.
type OutputMsg struct {
	Execcount int                    `json:"execution_count"`
	Data      map[string]string      `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ErrMsg encodes the traceback of errors output to the notebook
type ErrMsg struct {
	EName     string   `json:"ename"`
	EValue    string   `json:"evalue"`
	Traceback []string `json:"traceback"`
}

// HandleExecuteRequest runs code from an execute_request method, and sends the various
// reply messages.
func (k *Kernel) HandleExecuteRequest(receipt message.Receipt) error {

	// Actual execution handling

	reply, err := message.New(message.ExecuteReplyType, receipt.Message)
	if err != nil {
		return err
	}

	content := make(map[string]interface{})
	reqcontent := receipt.Message.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	silent := reqcontent["silent"].(bool)
	if !silent {
		k.ExecCounter++
	}
	content["execution_count"] = k.ExecCounter

	// TODO:: Hookup

	content["status"] = "ok"
	content["payload"] = make([]map[string]interface{}, 0)
	content["user_variables"] = make(map[string]string)
	content["user_expressions"] = make(map[string]string)

	out, err := message.New(message.ExecuteResultType, receipt.Message)
	if err != nil {
		log.WithError(err).Error("unable to create ExecuteResultType")
	} else {
		var outContent OutputMsg
		outContent.Execcount = k.ExecCounter
		outContent.Data = make(map[string]string)

		// lng, _ := lang.Detect(code)
		// pp.Println("lng = ", string(lng))
		// _ = lng
		outContent.Data["text/plain"] = fmt.Sprint("Hello.... got " + code + " ... detected the langauge to be ")
		outContent.Metadata = make(map[string]interface{})
		out.Content = outContent
		receipt.SendResponse(receipt.Connection.IOPubSocket, out)
	}

	// send the output back to the notebook
	reply.Content = content
	return receipt.SendResponse(receipt.Connection.ShellSocket, reply)
}
