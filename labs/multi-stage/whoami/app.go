package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Units.
const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
)

var (
	cert string
	key  string
	port string
	name string
	mode string
)

func init() {
	flag.StringVar(&cert, "cert", "", "give me a certificate")
	flag.StringVar(&key, "key", "", "give me a key")
	flag.StringVar(&port, "port", "80", "give me a port number")
	flag.StringVar(&name, "name", os.Getenv("WHOAMI_NAME"), "give me a name")
	flag.StringVar(&mode, "mode", os.Getenv("WHOAMI_MODE"), "set mode: q=quiet")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()

	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/bench", benchHandler)
	http.HandleFunc("/", whoamiHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Starting up on port " + port)

	if len(cert) > 0 && len(key) > 0 {
		log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, nil))
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func benchHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, "1")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		printBinary(p)
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}

func printBinary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	queryParams := u.Query()

	size, err := strconv.ParseInt(queryParams.Get("size"), 10, 64)
	if err != nil {
		size = 1
	}
	if size < 0 {
		size = 0
	}

	unit := queryParams.Get("unit")
	switch strings.ToLower(unit) {
	case "kb":
		size *= KB
	case "mb":
		size *= MB
	case "gb":
		size *= GB
	case "tb":
		size *= TB
	}

	attachment, err := strconv.ParseBool(queryParams.Get("attachment"))
	if err != nil {
		attachment = false
	}

	content := fillContent(size)

	if attachment {
		w.Header().Set("Content-Disposition", "Attachment")
		http.ServeContent(w, r, "data.txt", time.Now(), content)
		return
	}

	if _, err := io.Copy(w, content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func whoamiHandler(w http.ResponseWriter, req *http.Request) {
	mutexHealthState.RLock()
	defer mutexHealthState.RUnlock()
	if currentHealthState.StatusCode != 200 {
		w.WriteHeader(currentHealthState.StatusCode)
		return
	}

	u, _ := url.Parse(req.URL.String())
	wait := u.Query().Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}

	if name != "" {
		_, _ = fmt.Fprintln(w, "Name:", name)
	}

	hostname, _ := os.Hostname()
	_, _ = fmt.Fprintln(w, "Hostname:", hostname)

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !ip.IsLoopback() {
				_, _ = fmt.Fprintln(w, "IP:", ip)
			}
		}
	}

	if (mode == "q") {
		return
	}

	_, _ = fmt.Fprintln(w, "RemoteAddr:", req.RemoteAddr)
	if err := req.Write(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()

	data := struct {
		Hostname string      `json:"hostname,omitempty"`
		IP       []string    `json:"ip,omitempty"`
		Headers  http.Header `json:"headers,omitempty"`
		URL      string      `json:"url,omitempty"`
		Host     string      `json:"host,omitempty"`
		Method   string      `json:"method,omitempty"`
		Name     string      `json:"name,omitempty"`
	}{
		Hostname: hostname,
		IP:       []string{},
		Headers:  req.Header,
		URL:      req.URL.RequestURI(),
		Host:     req.Host,
		Method:   req.Method,
		Name:     name,
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				data.IP = append(data.IP, ip.String())
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type healthState struct {
	StatusCode int
}

var (
	currentHealthState = healthState{http.StatusOK}
	mutexHealthState   = &sync.RWMutex{}
)

func healthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var statusCode int

		if err := json.NewDecoder(req.Body).Decode(&statusCode); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Update health check status code [%d]\n", statusCode)

		mutexHealthState.Lock()
		defer mutexHealthState.Unlock()
		currentHealthState.StatusCode = statusCode
	} else {
		mutexHealthState.RLock()
		defer mutexHealthState.RUnlock()
		w.WriteHeader(currentHealthState.StatusCode)
	}
}

func fillContent(length int64) io.ReadSeeker {
	charset := "-ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[i%len(charset)]
	}

	if length > 0 {
		b[0] = '|'
		b[length-1] = '|'
	}

	return bytes.NewReader(b)
}
