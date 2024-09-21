package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IvMaslov/ethernet"
	"github.com/IvMaslov/socket"
)

func main() {
	ifce, err := socket.New(socket.WithDevice("virbr0"))
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		etherEcho(ctx, ifce)
		wg.Done()
	}()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)

	<-sig

	cancel()
	wg.Wait()

	ifce.Close()
}

func etherEcho(ctx context.Context, ifce *socket.Interface) {
	buf := make([]byte, 1500)

	for {
		select {
		case <-ctx.Done():
			log.Println("Exiting from ether echo")
			return
		default:
			n, err := ifce.Read(buf)
			if err != nil {
				log.Println("ERROR with reading -", err)
			}

			f := &ethernet.Frame{}
			err = f.Unmarshal(buf[:n])
			if err != nil {
				log.Println("ERROR with unmarshaling -", err)
				continue
			}

			fmt.Printf("Packet %s -> %s | with len %d\n", f.SrcHarwAddr, f.DestHarwAddr, n)

			f.DestHarwAddr, f.SrcHarwAddr = f.SrcHarwAddr, f.DestHarwAddr

			data, err := f.Marshal()
			if err != nil {
				log.Println("ERROR with marshaling -", err)
				continue
			}

			_, err = ifce.Write(data)
			if err != nil {
				log.Println("ERROR with writing -", err)
			}
		}
	}
}
