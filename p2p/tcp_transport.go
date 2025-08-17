package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over TCP established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer.
	conn net.Conn
	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	rpcChan chan RPC

	// mu    sync.RWMutex
	// peers map[net.Addr]Peer
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implements the peer interface.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		rpcChan:          make(chan RPC),
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenandAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil
}

// Consume implements the transport interface, which will return
// read-only channel
// for reading the incoming msgs received from another peer in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChan
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handleConn(conn)
	}

}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	// Readloop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			fmt.Printf("TCP read error: %s\n", err)
			return
		}
		rpc.From = conn.RemoteAddr()

		t.rpcChan <- rpc
	}

}
