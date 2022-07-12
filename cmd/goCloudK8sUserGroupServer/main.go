package main

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/version"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort = 8080
	webRootDir  = "web/"
)

// content holds our static web server content.
//go:embed web/*
var content embed.FS

var contentHandler = echo.WrapHandler(http.FileServer(http.FS(content)))

// The embedded files will all be in the '/web' folder so need to rewrite the request (could also do this with fs.Sub)
var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/web/$1"})

func customHTTPErrorHandler(err error, c echo.Context) {
	errorPage := webRootDir + "index.html"
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func main() {
	l := log.New(os.Stdout, fmt.Sprintf("HTTP_SERVER_%s ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("INFO: 'Starting %s v:%s  rev:%s '", version.APP, version.VERSION, version.REVISION)
	listenAddr, err := config.GetPortFromEnv(defaultPort)
	if err != nil {
		log.Fatalf("💥💥 ERROR: 'calling GetPortFromEnv got error: %v'\n", err)
	}
	l.Printf("INFO: 'HTTP server listening on port %s'", listenAddr)
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.GET("/*", contentHandler, contentRewrite)
	err = e.Start(listenAddr)
	if err != nil {
		log.Fatalf("💥💥 ERROR: 'calling echo.Start(%s) got error: %v'\n", listenAddr, err)
	}

}
