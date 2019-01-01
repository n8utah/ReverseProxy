package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:80", "http://localhost"}
	// config.AllowHeaders = []string{"Authorization"}
	// config.MaxAge = 12 * time.Hour
	// router.Use(cors.New(config))

	router.NoRoute(proxs.ginHandle)
	log.Fatal(router.Run(":" + *port))
}

func UnknownProxyServer(c *gin.Context) {
	w.Write([]byte("Reverse proxy Server Running. Proxy not found"))
}
