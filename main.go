package main

import (
	"flag"
	"log"
	"net/http"
)

var port *string

// var redirecturl *string

func main() {
	const (
		defaultPort      = "9090"
		defaultPortUsage = "default server port, ':9090'"
	// 	defaultTarget      = "http://127.0.0.1:80"
	// 	defaultTargetUsage = "default redirect url, 'http://127.0.0.1:8080'"
	)

	// // flags
	port = flag.String("port", defaultPort, defaultPortUsage)
	// redirecturl = flag.String("url", defaultTarget, defaultTargetUsage)

	// flag.Parse()

	// fmt.Println("server will run on :", *port)
	// fmt.Println("redirecting to :", *redirecturl)

	proxs = make(Proxies, 0)
	// proxy
	proxs.NewProxy("http://localhost", "http://127.0.0.1:80")
	proxs.NewProxy("http://127.0.0.2:80", "http://127.0.0.1:81")

	http.HandleFunc("/proxyServer", ProxyServer)

	// server redirection
	http.HandleFunc("/", proxs.handle)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func ProxyServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Reverse proxy Server Running. Accepting at port:" + *port))
}

func UnknownProxyServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Reverse proxy Server Running. Proxy not found"))
}
