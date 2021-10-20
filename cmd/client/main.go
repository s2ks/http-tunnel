package main

import (
	"net"
	"net/http"
	"log"
	"flag"
	"encoding/base64"
	"strings"

	tunnel_config "github.com/s2ks/http-tunnel/internal/config"
	tunnel_encoding "github.com/s2ks/http-tunnel/internal/encoding"
)

var (
	/* TODO default search path for config file */
	config_file_option = flag.String("config", "", "Path to a configuration file")
)

type Client struct {
	Listener net.Listener
	Connect []string
}

func Handle(conn net.Conn, connect []string) {
	log.Print(conn.RemoteAddr().String())
	log.Print(conn.LocalAddr().String())

	buf := make([]byte, 1024)

	for n, err := conn.Read(buf); n > 0; {
		if err != nil {
			log.Print(err)
			return
		}

		b64data := tunnel_encoding.Encode(buf[0:n])

		var resp http.Respone
		var err error

		/* Send a POST with the base64 encoded original request
		to the host(s) specified by connect.
		--
		Try all hosts specified in connect in order either
		until one succeeds or we run out of addresses. */
		for _, addr := range connect {
			resp, err := http.Post(addr, "text/plain",
				strings.NewReader(b64data))

			if err != nil {
				log.Print(err)
			} else {
				err = nil
				break
			}
		}

		/* All addresses were tried, and failed */
		if err != nil {
			return
		}

		/* The response body is empty */
		if resp.ContentLength == 0 {
			log.Print("ContentLength = 0; no data")
			return
		}

		respbuf := make([]byte, resp.ContentLength)
		rn, err := resp.Body.Read(respbuf)

		if err != nil {
			log.Print(err)
			return
		}

		if rn == 0 {
			log.Print("No data")
			return
		}

		decoded, err := tunnel_encoding.Decode(respbuf)

		log.Print(decoded)

		/* TODO send the decoded packet back to the remote */
		//conn.RemoteAddr().String()
	}
}

func main() {
	flag.Parse()

	if *config_file_option == "" {
		log.Fatal("Please provide a configuration file")
	}

	cfg_file, err := os.Open(*config_file_option)

	if err != nil {
		log.Fatal(err)
	}

	cfg_map, err := tunnel_config.Parse(cfg_file)

	if err != nil {
		log.Fatal(err)
	}

	clients := make([]Client, 0)

	for name, _ := range cfg_map {
		accept 	:= cfg_map[name]["accept"]
		connect := cfg_map[name]["connect"]

		 for _,  addr := range accept {
			client := new(Client)
			client.Listener, err = net.Listen("tcp", addr)

			if err != nil {
				log.Fatal(err)
			}

			client.Connect = connect

			clients = append(clients, client)
		 }
	}

	for {
		for _, client := range clients {
			conn, err := client.Listener.Accept()

			if err != nil {
				log.Fatal(err)
			}

			go Handle(conn, client.Connect)
		}
	}
}
