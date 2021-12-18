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
		wg.Add(1)
		defer wg.Done()
	}
	for {
		n, err := src.Read(buf)

		if err != nil {
			log.Print(err)
			return
		}

		_, err = dest.Write(buf[:n])

		if err != nil {
			log.Print(err)
			return
		}
	}
}
