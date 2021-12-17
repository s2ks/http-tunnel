package main

import (
	"flag"
	"log"
	"os"
	"net"
	"net/http"
	"fmt"
	"sync"

	tunnel_config "github.com/s2ks/http-tunnel/internal/config"
	tunnel_encoding "github.com/s2ks/http-tunnel/internal/encoding"
)

var (
	config_file_option = flag.String("config", "", "Path to a configuration file")
)

type TunnelHandler struct {
	Name            string
	Paths           []string
	Connect         []string

	conn 		net.Conn
	wg 		sync.WaitGroup
}

func (t *TunnelHandler) forward(src io.Reader, dest io.Writer) {
	t.wg.Add(1)
	buf := make([]byte, 0xffff)

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
	t.wg.Done()
}

func dialAny(string []addrv) (*net.TCPConn, error) {

	for _, addr := range addrv {
		tcpaddr, err := net.ResolveTCPAddr("tcp", addr)

		if err != nil {
			log.Print(err)
			continue
		}

		tcpconn, err := net.DialTCP("tcp", nil, tcpaddr)

		if err != nil {
			log.Print(err)
			continue
		} else {
			return tcpconn, nil
		}
	}

	return nil, fmt.Errorf("Unable to connect to any of the addresses given")
}

func (t *TunnelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	found := false
	for _, path := range t.Paths {
		found |= (path == r.URL.Path)
	}

	if found == false {
		http.NotFound(w, r)
		log.Print("Path " + r.URL.Path + " has no handler")
		return
	}

	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		log.Print("Method " + r.Method " will not be handled, we only" +
		"handle http POST")
		return
	}

	/* Should return a type that implements Read */
	body_decoder := tunnel_encoding.NewDecoderFromReader(r.Body) /* TODO implement */
	go t.forward(body_decoder, t.conn)
	go t.forward(t.conn, w)

	/* Wait for goroutines to finish */
	t.wg.Wait()
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

	for name, _ := range cfg_map {
		accept  := cfg_map[name]["accept"][0]
		connect := cfg_map[name]["connect"]
		paths   := cfg_map[name]["endpoint"]

		conn, err := dialAny(connect)

		if err != nil {
			log.Fatal(err)
		}

		for _, addr := range accept {
			tunnel_handler          := new(TunnelHandler)
			tunnel_handler.Connect  = connect
			tunnel_handler.Paths    = paths
			tunnel_handler.Name     = name
			tunnel_handler.conn 	= conn

			go func() {
				for _, path := range paths {
					http.Handle(path, tunnel_handler)
				}

				err := http.ListenAndServe(addr, nil)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}
}
