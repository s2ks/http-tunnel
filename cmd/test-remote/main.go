package main

import (
	"net"
	"flag"
	"log"
	"sync"

	tunnel_util "github.com/s2ks/http-tunnel/internal/util"
)

var (
	listen_opt = flag.String("listen", "", "Listen address")
)

func main() {
	flag.Parse()

	log.SetPrefix("http-tunnel test-remote: ")

	if *listen_opt == "" {
		log.Fatal("Please provide a valid -listen address (see -help)")
	}

	listener, err := net.Listen("tcp", *listen_opt)


	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s\n", *listen_opt)

	var wg sync.WaitGroup

	//tunnel_util.ForwarderBuffsize = 8
	//log.Print(tunnel_util.ForwarderBuffsize)

	for {
		log.Print("Poll...")
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		//wg.Add(1)
		go tunnel_util.Forward(conn, conn, nil)

		wg.Wait()
	}

}
