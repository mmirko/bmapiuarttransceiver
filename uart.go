package bmapiuarttransceiver

import (
	"context"
	"log"

	"go.bug.st/serial"
)

func UartTransceiver(ctx context.Context, device string) (chan<- uint8, <-chan uint8) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal(err)
	}
	src := make(chan uint8)
	dst := make(chan uint8)
	buff := make([]byte, 100)
	go func() {
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
				break
			}
			if n == 0 {
				break
			}
			for i := 0; i < n; i++ {
				select {
				case <-ctx.Done():
					return
				case dst <- buff[i]:
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-src:
				_, err := port.Write([]byte{data})
				if err != nil {
					log.Fatal(err)
					break
				}
			}
		}
	}()

	return src, dst
}
