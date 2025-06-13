package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IvMaslov/ethernet"
	"github.com/IvMaslov/socket"
)

// simple ethernet echo server
func main() {
	sock, err := socket.New(socket.WithDevice("eth0"))
	if err != nil {
		log.Fatal(err)
	}

	etherSocker, err := ethernet.NewEtherSocket(sock, 0)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		etherEcho(ctx, etherSocker)
		wg.Done()
	}()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)

	<-sig

	sock.Close()

	cancel()
	wg.Wait()
}

func etherEcho(ctx context.Context, etherSock *ethernet.EtherSocket) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Exiting from ether echo")
			return
		default:
			f, err := etherSock.ReadFrame()
			if err != nil {
				log.Println("ERROR with reading -", err)
				continue
			}

			log.Printf("Packet %s -> %s\n", f.SrcHarwAddr, f.DestHarwAddr)

			f.DestHarwAddr, f.SrcHarwAddr = f.SrcHarwAddr, f.DestHarwAddr

			err = etherSock.Write(f.Payload)
			if err != nil {
				log.Println("ERROR with writing -", err)
			}
		}
	}
}
