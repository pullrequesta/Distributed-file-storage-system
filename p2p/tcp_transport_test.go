package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	tcpOpts := TCPTransportOpts{
		ListenAddr:    "127.0.0.1:3000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOpts)

	assert.Equal(t, tr.ListenAddr, "127.0.0.1:3000")

	//Server
	assert.Nil(t, tr.ListenandAccept())
}
