/* The http-tunnel client listens on one or more (local) addresses, It will
accept incoming connection requests, read from the connection, encode
the read bytes in websafe base64 and wrap it a HTTP POST before sending
it off to the specified url.

The idea is something like this:

┌──────────┐   ┌──────────────────┐   ┌────────────────┐
│          ├──►│http-tunnel client├──►│Remote webserver│
│  Client  │   │                  │   │                │
│          │◄──┤    HTTP POST     │◄──┤  (e.g. nginx)  │
└──────────┘   └──────────────────┘   └──────┬─────────┘
                                             │ ▲
                                             │ │
                                             │ │
                                             ▼ │
                                   ┌───────────┴─────────┐
                                   │                     │
                    ... ◄────────► │  http-tunnel server │
                                   │                     │
                                   └─────────────────────┘
*/
package main

import (
	"net"
	"net/http"
	"log"
	"flag"
	"strings"
	"os"
	"fmt"
	//"sync"

	tunnel_config "github.com/s2ks/http-tunnel/internal/config"
	tunnel_encoding "github.com/s2ks/http-tunnel/internal/encoding"
	tunnel_util "github.com/s2ks/http-tunnel/internal/util"
)

var (
	/* TODO default search path for config file */
	config_file_option = flag.String("config", "", "Path to a configuration file")
	client_bufsize = 0xffff
)

type Client struct {
	Listener	net.Listener
	Url		string
}

func (c *Client) Handle(conn net.Conn) {
	/* TODO use a larger buffer, think about adding a configuration option
	for maximum body size. The largest the body can be depends on the webserver.
	For nginx the default is 1MB. */

	buf := make([]byte, client_bufsize)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			log.Print(err)
			return
		}

		/* Base64 encode */
		b64data := tunnel_encoding.Encode(buf[:n])

		/* Send a POST with base64 encoded data from the connection */
		resp, err := http.Post(c.Url, "text/plain", strings.NewReader(b64data))

		if err != nil {
			log.Print(err)
			return
		}

		tunnel_util.Forward(resp.Body, conn, nil)
		resp.Body.Close()
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
		accept	:= cfg_map[name]["accept"]
		url	:= cfg_map[name]["url"]

		if len(url) == 0 {
			log.Fatal(fmt.Errorf("No URL option specified:\n" +
			"\tSection [%s] does not appear to contain a URL" +
			"option (url=<dest>)", name))
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
