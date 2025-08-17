package main

import (
	"fmt"
	"log"

	"github.com/pullrequest/Distributed-file-strorage-system/p2p"
)

func Onpeer(peer p2p.Peer) error {
	peer.Close()
	return nil

}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    "127.0.0.1:3000",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		OnPeer:        Onpeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		msg := <-tr.Consume()
		fmt.Printf("message: %+v\n", msg)
	}()

	if err := tr.ListenandAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
