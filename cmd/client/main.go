package main

import (
	"net"
	"net/http"
	"log"
	"flag"
	"strings"
	"os"
	"fmt"

	tunnel_config "github.com/s2ks/http-tunnel/internal/config"
	tunnel_encoding "github.com/s2ks/http-tunnel/internal/encoding"
)

var (
	/* TODO default search path for config file */
	config_file_option = flag.String("config", "", "Path to a configuration file")
)

type Client struct {
	Listener 	net.Listener
	Url 		string
}

func (c *Client) Handle(conn net.Conn) {
	buf := make([]byte, 1024)

	for n, err := conn.Read(buf); n > 0; {
		if err != nil {
			log.Fatal(err)
		}

		/* Base64 encode */
		b64data := tunnel_encoding.Encode(buf)

		/* Send a POST with base64 encoded data from the connection */
		resp, err := http.Post(c.Url, "text/plain", strings.NewReader(b64data))

		if err != nil {
			log.Fatal(err)
		}

		/* The response body is empty */
		if resp.ContentLength > 0 {
			respbuf := make([]byte, resp.ContentLength)
			_, err := resp.Body.Read(respbuf)

			if err != nil {
				log.Fatal(err)
			}

			decoded, err := tunnel_encoding.Decode(respbuf)

			if err != nil {
				log.Fatal(err)
			}

			conn.Write(decoded)
		}
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

	clients := make([]*Client, 0)

	for name, _ := range cfg_map {
		accept 	:= cfg_map[name]["accept"]
		url 	:= cfg_map[name]["url"]

		if len(url) == 0 {
			log.Fatal(fmt.Errorf("No url option specified in configuration" +
				"file"))
		}

		for _,  addr := range accept {
			client := new(Client)
			client.Listener, err = net.Listen("tcp", addr)

			if err != nil {
				log.Fatal(err)
			}

			client.Url = url[0]

			clients = append(clients, client)
		}
	}

	for {
		for _, c := range clients {
		 	conn, err := c.Listener.Accept()

			if err != nil {
				log.Fatal(err)
			}

			go c.Handle(conn)
		}
	}
}
