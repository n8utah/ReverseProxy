package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
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
func (ps *Proxies) ginHandle(c *gin.Context) {
	c.Writer.Header().Set("X-GoProxy", "GoProxy")
	url, err := url.Parse("http://" + c.Request.Host)
	if err != nil {
		fmt.Println("Parse Error:" + err.Error())
		return
	} else {
		// fmt.Println("RequestHost:" + url.Hostname())
	}

	for i, _ := range *ps {
		if url.Hostname() == (*ps)[i].host.Hostname() {
			(*ps)[i].proxy.Transport = &myTransport{}
			(*ps)[i].proxy.ServeHTTP(c.Writer, c.Request)
			return
		}
	} //EO go through proxies

	//if not found then we serve up the default
	UnknownProxyServer(c)
	return
}

var proxs Proxies

func AddProxyHandle(c *gin.Context) {
	c.JSON(http.StatusOK, "Reverse proxy Server Running. Accepting at port:"+*port)
}
