package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	gnode "github.com/dapperlabs/bamboo-node/pkg/network/gossip/v1"
	"github.com/gogo/protobuf/proto"
)

// Demo of for the gossip async node implementation
// How to run: just start three instances of this program. The nodes will
// communicate with each other and place gossip messages.

func main() {
	portPool := []string{"127.0.0.1:50000", "127.0.0.1:50001", "127.0.0.1:50002"}

	listener, err := pickPort(portPool)
	if err != nil {
		log.Fatal(err)
	}

	servePort := listener.Addr().String()

	fmt.Println(servePort)
	if err != nil {
		log.Fatal(err)
	}
	node := gnode.NewNode()

	Time := func(payloadBytes []byte) ([]byte, error) {
		newMsg := &Message{}
		if err := proto.Unmarshal(payloadBytes, newMsg); err != nil {
			return nil, fmt.Errorf("could not unmarshal payload: %v", err)
		}

		log.Printf("Payload: %v", string(newMsg.Text))
		time.Sleep(2 * time.Second)
		fmt.Printf("The time is: %v\n", time.Now().Unix())
		return []byte("Pong"), nil
	}

	node.RegisterFunc("Time", Time)

	go node.Serve(listener)

	peers := make([]string, 0)
	for _, port := range portPool {
		if port != servePort {
			peers = append(peers, port)
		}
	}

	t := time.Tick(5 * time.Second)

	for {
		select {
		case <-t:
			go func() {
				log.Println("Gossiping")
				payload := &Message{Text: []byte("Ping")}
				bytes, err := proto.Marshal(payload)
				if err != nil {
					log.Fatalf("could not marshal message: %v", err)
				}
				// You can try to change the line bellow to AsyncGossip(...), when you do
				// so you will notice that the responses returned to you will be empty
				// (that is because AyncGossip does not wait for the sent messages to be
				// processed)
				resps, err := node.SyncGossip(context.Background(), bytes, peers, "Time")
				if err != nil {
					log.Println(err)
				}
				for _, resp := range resps {
					if resp == nil {
						continue
					}
					log.Printf("Response: %v\n", string(resp.ResponseByte))
				}
			}()
		}
	}
}

func pickPort(portPool []string) (net.Listener, error) {
	for _, port := range portPool {
		ln, err := net.Listen("tcp4", port)
		if err == nil {
			return ln, nil
		}
	}

	return nil, fmt.Errorf("could not find an empty port in the given pool")
}
