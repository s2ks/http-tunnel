package util

import (
	"io"
	"sync"
	"log"
)

var (
	ForwarderBuffsize = 0xffff
)

func Forward(src io.Reader, dest io.Writer, wg *sync.WaitGroup) {
	buf := make([]byte, ForwarderBuffsize)

	if wg != nil {
		defer wg.Done()
		defer log.Print("Forwarder done...")
	}

	for {
		n, err := src.Read(buf)

		if err != nil && n == 0 {
			log.Print(err)
			return
		}

		_, err = dest.Write(buf[:n])

		if err != nil {
			log.Print(err)
		}
	}
}
