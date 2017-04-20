package connection

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/pkg/errors"
)

// ConnectionOptions stores the contents of the kernel connection file created by Jupyter.
type Options struct {
	SignatureScheme string `json:"signature_scheme"`
	Transport       string `json:"transport"`
	StdinPort       int    `json:"stdin_port"`
	ControlPort     int    `json:"control_port"`
	IOPubPort       int    `json:"iopub_port"`
	HeartBeatPort   int    `json:"hb_port"`
	ShellPort       int    `json:"shell_port"`
	Key             string `json:"key"`
	IP              string `json:"ip"`
}

// SocketGroup holds the sockets needed to communicate with the kernel, and
// the key for message signing.
type Connection struct {
	ShellSocket   *zmq.Socket
	ControlSocket *zmq.Socket
	StdinSocket   *zmq.Socket
	IOPubSocket   *zmq.Socket
	Key           []byte
	Options       *Options
}

func New(opts *Options) (*Connection, error) {

	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create zmq context")
	}
	shellSocket, err := ctx.NewSocket(zmq.ROUTER)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create shell socket")
	}
	controlSocket, err := ctx.NewSocket(zmq.ROUTER)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create control socket")
	}
	stdinSocket, err := ctx.NewSocket(zmq.ROUTER)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create stdin socket")
	}
	ioPubSocket, err := ctx.NewSocket(zmq.PUB)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create iopub socket")
	}

	address := fmt.Sprintf("%v://%v:%%v", opts.Transport, opts.IP)

	if err := shellSocket.Bind(fmt.Sprintf(address, opts.ShellPort)); err != nil {
		return nil, errors.Wrap(err, "cannot bind shell socket")
	}
	if err := controlSocket.Bind(fmt.Sprintf(address, opts.ControlPort)); err != nil {
		return nil, errors.Wrap(err, "cannot bind control socket")
	}
	if err := stdinSocket.Bind(fmt.Sprintf(address, opts.StdinPort)); err != nil {
		return nil, errors.Wrap(err, "cannot bind stdin socket")
	}
	if err := ioPubSocket.Bind(fmt.Sprintf(address, opts.IOPubPort)); err != nil {
		return nil, errors.Wrap(err, "cannot bind iopub socket")
	}

	// Message signing key
	key := opts.Key

	// Start the heartbeat device
	heartBeatSocket, err := ctx.NewSocket(zmq.REP)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get the Heartbeat device socket")
	}
	heartBeatSocket.Bind(fmt.Sprintf(address, opts.HeartBeatPort))

	go func() {
		err := zmq.Proxy(heartBeatSocket, heartBeatSocket, nil)
		if err != nil {
			log.WithError(err).Error("unable to create heatbeat proxy")
		}
	}()

	return &Connection{
		ShellSocket:   shellSocket,
		ControlSocket: controlSocket,
		StdinSocket:   stdinSocket,
		IOPubSocket:   ioPubSocket,
		Key:           []byte(key),
		Options:       opts,
	}, nil
}
