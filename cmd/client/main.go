package main

import (
	"net"
	"net/http"
	"log"
	"flag"
	"encoding/base64"
	"strings"

	//tunnel_config "github.com/s2ks/http-tunnel/internal/config"
	tunnel_encoding "github.com/s2ks/http-tunnel/internal/encoding"
)


/* TODO Use a scheme similar to stunnel with both the
client and the server having accept/connect options
to specify on what address to accept the connection and
to what address to send the encoded/decoded packet to. */

var (
	connect_option 	= flag.String("connect", "", "Address to connect to")
	protocol_option = flag.String("protocol", "tcp",
		"Network protocol to use: tcp/tcp4/tcp6/unix/unixpacket")
	accept_option = flag.String("accept", "localhost:12345",
		"Address to listen on: [host]:port")

	/* TODO default search path for config file */
	config_file_option = flag.String("config", "", "Path to a configuration file")
)

func Handle(conn net.Conn) {
	buf := make([]byte, 512)

	log.Print(conn.RemoteAddr().String())
	log.Print(conn.LocalAddr().String())

	for n, err := conn.Read(buf); n > 0; {
		if err != nil {
			log.Print(err)
			return
		}

		b64data := base64.RawURLEncoding.EncodeToString(buf)

		/* Send a POST with the base64 encoded original request
		to the host specified by 'connect_option' */
		resp, err := http.Post(*connect_option, "text/plain",
			strings.NewReader(b64data))

		if err != nil {
			log.Print(err)
			return
		}

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

	listener, err := net.Listen(*protocol_option, *accept_option)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go Handle(conn)
	}
}
