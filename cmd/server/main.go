package main

import (
	"flag"
	"log"
	"os"
	"net"
	"net/http"
	"fmt"
	"sync"
	//"io"
	"time"

	tunnel_config "github.com/s2ks/http-tunnel/internal/config"
	tunnel_encoding "github.com/s2ks/http-tunnel/internal/encoding"
	//tunnel_util "github.com/s2ks/http-tunnel/internal/util"
)

var (
	config_file_option = flag.String("config", "", "Path to a configuration file")
)

type TunnelHandler struct {
	Name            string
	Connect         []string

	conn 		net.Conn
	wg 		sync.WaitGroup
	recvchan 	chan []byte
}

/* Dial any one of the addresses specified in addrv */
func dialAny(addrv []string) (*net.TCPConn, error) {

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


func (t *TunnelHandler) connReader() {
	buf := make([]byte, 512)
	for {
		n, err := t.conn.Read(buf)

		if err != nil {
			log.Print(err)
		}

		if n == 0 {
			break
		}

		go func(b []byte) {
			t.recvchan <- b
		}(buf[:n])
	}
}
func (t *TunnelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		log.Print("Method " + r.Method + " will not be handled, we only" +
		"handle http POST")
		return
	}

	body_decoder := tunnel_encoding.NewDecoderFromReader(r.Body)

	/* TODO: get buffer size from a constant */
	buf := make([]byte, 512)
	for {
		n, err := body_decoder.Read(buf)

		if n == 0 {
			log.Print(err)
			break
		}

		t.conn.Write(buf[:n])

		timeout := false
		select {
			case recv := <-t.recvchan:
				log.Printf("Received %d bytes\n", len(recv))
				w.Write(recv)
			case <-time.After(5 * time.Second):
				timeout = true

		}

		/* Read additional bytes if there are any */
		if timeout == false {
			for stop := false; stop == false; {
				select {
				case recv := <-t.recvchan:
					log.Printf("Received %d bytes\n", len(recv))
					w.Write(recv)
				default:
					stop = true
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	log.SetPrefix("http-tunnel server: ")

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

	var wg sync.WaitGroup

	for name, _ := range cfg_map {
		accept  := cfg_map[name]["accept"]
		connect := cfg_map[name]["connect"]
		//paths   := cfg_map[name]["endpoint"]

		conn, err := dialAny(connect)

		if err != nil {
			log.Fatal(err)
		}

		tunnel_handler          := new(TunnelHandler)
		tunnel_handler.Connect  = connect
		//tunnel_handler.Paths    = paths
		tunnel_handler.Name     = name
		tunnel_handler.conn 	= conn
		tunnel_handler.recvchan = make(chan []byte)

		go tunnel_handler.connReader()

		http.Handle("/", tunnel_handler)

		for _, addr := range accept {
			wg.Add(1)
			go func() {
				defer wg.Done()
				//err := http.ListenAndServe(addr, servemux)
				log.Printf("listening on %s\n", addr)
				err := http.ListenAndServe(addr, nil)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}

	wg.Wait()
}
