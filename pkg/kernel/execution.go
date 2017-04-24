package kernel

import (
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

	content := make(map[string]interface{})

	// Actual execution handling

	reply, err := message.New(message.ExecuteReplyType, receipt.Message)
	if err != nil {
		return err
	}

	reqcontent := receipt.Message.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	silent := reqcontent["silent"].(bool)
	if !silent {
		k.ExecCounter++
	}

	sessionName := receipt.Message.Header.Session

	if !DefaultSessionManager.Has(sessionName) {
		NewSession(sessionName)
	}

	sess, err := DefaultSessionManager.Get(sessionName)
	if err != nil {
		content["status"] = "error"
		content["ename"] = "ERROR"
		content["evalue"] = err.Error()
		content["traceback"] = code
		errormsg, merr := message.New(message.ErrorReplyType, receipt.Message)
		if merr != nil {
			log.WithError(err).Error("unable to create ErrorReplyType")
			return err
		}
		errormsg.Content = ErrMsg{"Error", err.Error(), []string{sess.buf.String()}}
		defer sess.buf.Reset()
		receipt.SendResponse(receipt.Connection.IOPubSocket, errormsg)

		reply.Content = content
		return receipt.SendResponse(receipt.Connection.ShellSocket, reply)
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

		// lng := linguist.Detect(code)
		// pp.Println("lng = ", string(lng))
		// _ = lng

		if err := sess.Exec(code); err != nil {
			log.WithError(err).Error("unable to create exec code")
		}

		outContent.Data["text/plain"] = sess.buf.String()
		defer sess.buf.Reset()
		outContent.Metadata = make(map[string]interface{})
		out.Content = outContent
		receipt.SendResponse(receipt.Connection.IOPubSocket, out)
	}

	// send the output back to the notebook
	reply.Content = content
	return receipt.SendResponse(receipt.Connection.ShellSocket, reply)
}
