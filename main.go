package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Prox struct {
	host   *url.URL
	target *url.URL
	proxy  *httputil.ReverseProxy
}

type Proxies []*Prox

func NewProxy(host string, target string) *Prox {
	hurl, _ := url.Parse(host)
	url, _ := url.Parse(target)

	return &Prox{host: hurl, target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (ps *Proxies) NewProxy(host string, target string) {
	hurl, _ := url.Parse(host)
	url, _ := url.Parse(target)

	*ps = append(*ps, &Prox{host: hurl, target: url, proxy: httputil.NewSingleHostReverseProxy(url)})
	return
}

//handle checks the source and see if it matches any of the
//proxies we have setup
func (ps *Proxies) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	url, err := url.Parse("http://" + r.Host)
	if err != nil {
		fmt.Println("Parse Error:" + err.Error())
		return
	} else {
		// fmt.Println("RequestHost:" + url.Hostname())
	}

	for i, _ := range *ps {
		if url.Hostname() == (*ps)[i].host.Hostname() {
			(*ps)[i].proxy.Transport = &myTransport{}
			(*ps)[i].proxy.ServeHTTP(w, r)
			return
		}
	} //EO go through proxies

	//if not found then we serve up the default
	UnknownProxyServer(w, r)
	return
}

var port *string
var redirecturl *string

var proxs Proxies

func main() {
	const (
		defaultPort        = "9090"
		defaultPortUsage   = "default server port, ':9090'"
		defaultTarget      = "http://127.0.0.1:80"
		defaultTargetUsage = "default redirect url, 'http://127.0.0.1:8080'"
	)

	// flags
	port = flag.String("port", defaultPort, defaultPortUsage)
	redirecturl = flag.String("url", defaultTarget, defaultTargetUsage)

	flag.Parse()

	fmt.Println("server will run on :", *port)
	fmt.Println("redirecting to :", *redirecturl)

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
	w.Write([]byte("Reverse proxy Server Running. Accepting at port:" + *port + " Redirecting to :" + *redirecturl))

}

func UnknownProxyServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Reverse proxy Server Running. Proxy not found"))
}
