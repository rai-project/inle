package message

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/rai-project/inle/pkg/connection"
)

type Receipt struct {
	Message    Message
	Identities [][]byte
	Connection *connection.Connection
}

// SendResponse sends a message back to return identites of the received message.
func (receipt Receipt) SendResponse(socket *zmq.Socket, msg *Message) error {

	var err error

	for i := 0; i < len(receipt.Identities)-1; i++ {
		if _, err = socket.SendBytes(receipt.Identities[i], zmq.SNDMORE); err != nil {
			log.WithError(err).Error("unable to send bytes")
			return err
		}
	}
	_, err = socket.SendBytes(receipt.Identities[(len(receipt.Identities)-1)], zmq.SNDMORE)
	if err != nil {
		log.WithError(err).Error("unable to send bytes")
		return err
	}

	_, err = socket.Send("<IDS|MSG>", zmq.SNDMORE)
	if err != nil {
		log.WithError(err).Error("unable to send <IDS|MSG>")
		return err
	}

	newmsg := msg.ToWire(receipt.Connection.Key)

	for i := 0; i < len(newmsg)-1; i++ {
		if _, err = socket.SendBytes(newmsg[i], zmq.SNDMORE); err != nil {
			log.WithError(err).Error("unable to send bytes")
			return err
		}
	}
	_, err = socket.SendBytes(newmsg[(len(newmsg)-1)], 0)
	if err != nil {
		log.WithError(err).Error("unable to send bytes")
		return err
	}

	log.Debugf("<-- %s", msg.Header.Type)
	log.Debugf("%+v\n", msg.Content)
	return nil
}
