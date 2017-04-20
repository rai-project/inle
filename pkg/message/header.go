package message

import "time"

// The message header contains a pair of unique identifiers for the
// originating session and the actual message id, in addition to the
// username for the process that generated the message.  This is useful in
// collaborative settings where multiple users may be interacting with the
// same kernel simultaneously, so that frontends can label the various
// messages in a meaningful way.
type Header struct {
	ID       string      `json:"msg_id"`
	Username string      `json:"username"`
	Session  string      `json:"session"`
	Type     MessageType `json:"msg_type"`
	Date     time.Time   `json:"date"`
	Version  string      `json:"version"`
}
