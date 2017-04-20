package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/rai-project/inle/pkg/protocol"
	"github.com/rai-project/logger"
	"github.com/rai-project/uuid"
)

type Message struct {
	Header Header `json:"header"`
	// In a chain of messages, the header from the parent is copied so that
	// clients can track where messages come from.
	ParentHeader Header `json:"parent_header"`
	// Any metadata associated with the message.
	Metadata map[string]interface{} `json:"metadata"`
	// The actual content of the message must be a dict, whose structure
	// depends on the message type.
	Content interface{} `json:"content"`
	// optional: buffers is a list of binary data buffers for implementations
	// that support binary extensions to the protocol.
	Buffer []string `json:"buffer"`
}

// ToWireMsg translates a ComposedMsg into a multipart ZMQ message ready to send, and
// signs it. This does not add the return identities or the delimiter.
func (msg Message) ToWire(signkey []byte) (msgparts [][]byte) {

	msgparts = make([][]byte, 5)
	header, err := json.Marshal(msg.Header)
	if err != nil {
		logger.Fatal(err)
	}
	msgparts[1] = header
	parentHeader, err := json.Marshal(msg.ParentHeader)
	if err != nil {
		logger.Fatal(err)
	}
	msgparts[2] = parentHeader
	if msg.Metadata == nil {
		msg.Metadata = make(map[string]interface{})
	}
	metadata, err := json.Marshal(msg.Metadata)
	if err != nil {
		logger.Fatal(err)
	}
	msgparts[3] = metadata
	content, err := json.Marshal(msg.Content)
	if err != nil {
		logger.Fatal(err)
	}
	msgparts[4] = content

	// Sign the message
	if len(signkey) != 0 {
		mac := hmac.New(sha256.New, signkey)
		for _, msgpart := range msgparts[1:] {
			mac.Write(msgpart)
		}
		msgparts[0] = make([]byte, hex.EncodedLen(mac.Size()))
		hex.Encode(msgparts[0], mac.Sum(nil))
	}

	return
}

// NewFromWire translates a multipart ZMQ messages received from a socket into
// a Message struct and a slice of return identities. This includes verifying the
// message signature.
func FromWire(msgparts [][]byte, signkey []byte) (*Message, [][]byte, error) {

	msg := &Message{}
	identities := [][]byte{}
	var err error

	i := 0
	for string(msgparts[i]) != "<IDS|MSG>" {
		i++
	}
	identities = msgparts[:i]

	// msgparts[i] is the delimiter

	// Validate signature
	if len(signkey) != 0 {
		mac := hmac.New(sha256.New, signkey)
		for _, msgpart := range msgparts[i+2 : i+6] {
			mac.Write(msgpart)
		}
		signature := make([]byte, hex.DecodedLen(len(msgparts[i+1])))
		_, err := hex.Decode(signature, msgparts[i+1])
		if err != nil {
			return nil, nil, err
		}
		if !hmac.Equal(mac.Sum(nil), signature) {
			return nil, nil, InvalidSignatureError
		}
	}
	json.Unmarshal(msgparts[i+2], &msg.Header)
	json.Unmarshal(msgparts[i+3], &msg.ParentHeader)
	json.Unmarshal(msgparts[i+4], &msg.Metadata)
	json.Unmarshal(msgparts[i+5], &msg.Content)
	return msg, identities, err
}

// New creates a new Message to respond to a parent message. This includes setting
// up its headers.
func New(msgType MessageType, parent Message) (*Message, error) {
	id := uuid.NewV4()
	return &Message{
		ParentHeader: parent.Header,
		Header: Header{
			ID:       id,
			Session:  parent.Header.Session,
			Username: parent.Header.Username,
			Type:     msgType,
			Version:  protocol.Version,
		},
	}, nil
}
