package kernel

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/k0kubun/pp"
	zmq "github.com/pebbe/zmq4"
	"github.com/pkg/errors"
	"github.com/rai-project/inle/pkg/connection"
	"github.com/rai-project/inle/pkg/message"
)

type Kernel struct {
	id                string
	done              chan struct{}
	ExecCounter       int
	connectionOptions *connection.Options
}

func New(id string, connectionConfig io.Reader) (*Kernel, error) {
	buf, err := ioutil.ReadAll(connectionConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read connectionConfig")
	}
	connectionOptions := &connection.Options{}

	if err := json.Unmarshal(buf, connectionOptions); err != nil {
		return nil, errors.Wrap(err, "unable to parse connectionOptions")
	}

	return &Kernel{
		id:                id,
		done:              make(chan struct{}, 1),
		ExecCounter:       0,
		connectionOptions: connectionOptions,
	}, nil
}

func (k *Kernel) Start() error {
	conn, err := connection.New(k.connectionOptions)
	if err != nil {
		return err
	}

	pi := zmq.NewPoller()

	pi.Add(conn.ShellSocket, zmq.POLLIN)
	pi.Add(conn.StdinSocket, zmq.POLLIN)
	pi.Add(conn.ControlSocket, zmq.POLLIN)

	var msgparts [][]byte
	var polled []zmq.Polled
	// Message receiving loop:
	go func() {
		for {
			polled, err = pi.Poll(-1)
			if err != nil {
				log.WithError(err).Error("failed to poll item")
				continue
			}
			switch {
			case polled[0].Events&zmq.POLLIN != 0: // shell socket
				msgparts, _ = polled[0].Socket.RecvMessageBytes(0)
				msg, ids, err := message.FromWire(msgparts, conn.Key)
				if err != nil {
					log.WithError(err).Error("cannot read message from wire")
					continue
				}
				if false {
					log.Println("received shell message: ", msg)
				}
				k.HandleShellMsg(message.Receipt{
					Message:    *msg,
					Identities: ids,
					Connection: conn,
				})
			case polled[1].Events&zmq.POLLIN != 0: // stdin socket - not implemented.
				polled[1].Socket.RecvMessageBytes(0)
			case polled[2].Events&zmq.POLLIN != 0: // control socket - treat like shell socket.
				msgparts, _ = polled[2].Socket.RecvMessageBytes(0)
				msg, ids, err := message.FromWire(msgparts, conn.Key)
				if err != nil {
					log.WithError(err).Error("cannot read message from wire")
					continue
				}
				log.Debug("received control message: ", msg)
				k.HandleShellMsg(message.Receipt{
					Message:    *msg,
					Identities: ids,
					Connection: conn,
				})
			}
		}
	}()

	return nil
}

func (k *Kernel) Wait() error {
	<-k.done
	return nil
}

func (k *Kernel) Stop() error {
	return nil
}

func (k *Kernel) Shutdown() error {
	k.done <- struct{}{}
	return nil
}

func (k *Kernel) Run() error {
	if err := k.Start(); err != nil {
		return err
	}
	if err := k.Wait(); err != nil {
		return err
	}
	return k.Stop()
}

// HandleShellMsg responds to a message on the shell ROUTER socket.
func (k *Kernel) HandleShellMsg(receipt message.Receipt) error {
	switch receipt.Message.Header.Type {
	case message.KernelInfoRequestType:
		k.HandleWithStatus(receipt, k.HandleKernelInfoRequest)
	case message.ConnectRequestType:
		k.HandleWithStatus(receipt, k.HandleConnectRequest)
	case message.CommOpenType:
		k.HandleWithStatus(receipt, k.HandleCommOpen)
	case message.CommCloseType:
		k.HandleWithStatus(receipt, k.HandleCommClose)
	case message.CommInfoRequestType:
		k.HandleWithStatus(receipt, k.HandleCommInfoRequest)
	case message.CommMessageType:
		k.HandleWithStatus(receipt, k.HandleCommMessage)
	case message.ExecuteRequestType:
		k.HandleWithStatus(receipt, k.HandleExecuteRequest)
	case message.ShutdownRequestType:
		k.HandleWithStatus(receipt, k.HandleShutdownRequest)
	default:
		err := errors.Errorf("Unhandled shell message: %v", pp.Sprint(receipt))
		log.WithError(err).Error()
		return err
	}
	return nil
}
