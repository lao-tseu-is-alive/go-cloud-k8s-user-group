package goserver

import (
	"context"
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultProtocol        = "http"
	defaultPort            = 8080
	defaultServerIp        = ""
	defaultServerPath      = "/"
	defaultSecondsToSleep  = 3
	secondsShutDownTimeout = 5 * time.Second  // maximum number of second to wait before closing server
	defaultReadTimeout     = 10 * time.Second // max time to read request from the client
	defaultWriteTimeout    = 10 * time.Second // max time to write response to the client
	defaultIdleTimeout     = 2 * time.Minute  // max time for connections using TCP Keep-Alive
)

// GoHttpServer is a struct type to store information related to all handlers of web server
type GoHttpServer struct {
	listenAddress string
	logger        *log.Logger
	store         users.Storage
	e             *echo.Echo
	router        *http.ServeMux
	startTime     time.Time
	httpServer    http.Server
}

// WaitForHttpServer attempts to establish a TCP connection to listenAddress
// in a given amount of time. It returns upon a successful connection;
// otherwise exits with an error.
func WaitForHttpServer(listenAddress string, waitDuration time.Duration, numRetries int) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; i < numRetries; i++ {
		//conn, err := net.DialTimeout("tcp", listenAddress, dialTimeout)
		resp, err := httpClient.Get(listenAddress)
		if err != nil {
			fmt.Printf("\n[%d] Cannot make http get %s: %v\n", i, listenAddress, err)
			time.Sleep(waitDuration)
			continue
		}
		// All seems is good
		fmt.Printf("OK: Server responded after %d retries, with status code %d ", i, resp.StatusCode)
		return
	}
	log.Fatalf("Server %s not ready up after %d attempts", listenAddress, numRetries)
}

// waitForShutdownToExit will wait for interrupt signal SIGINT or SIGTERM and gracefully shutdown the server after secondsToWait seconds.
func waitForShutdownToExit(srv *http.Server, secondsToWait time.Duration) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	// wait for SIGINT (interrupt) 	: ctrl + C keypress, or in a shell : kill -SIGINT processId
	sig := <-interruptChan
	srv.ErrorLog.Printf("INFO: 'SIGINT %d interrupt signal received, about to shut down server after max %v seconds...'\n", sig, secondsToWait.Seconds())

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), secondsToWait)
	defer cancel()
	// gracefully shuts down the server without interrupting any active connections
	// as long as the actives connections last less than shutDownTimeout
	// https://pkg.go.dev/net/http#Server.Shutdown
	if err := srv.Shutdown(ctx); err != nil {
		srv.ErrorLog.Printf("💥💥 ERROR: 'Problem doing Shutdown %v'\n", err)
	}
	<-ctx.Done()
	srv.ErrorLog.Println("INFO: 'Server gracefully stopped, will exit'")
	os.Exit(0)
}

// NewGoHttpServer is a constructor that initializes the server,routes and all fields in GoHttpServer type
func NewGoHttpServer(listenAddress string, l *log.Logger, store users.Storage, webRootDir string, content embed.FS) *GoHttpServer {
	myServerMux := http.NewServeMux()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	myUsersApi := users.Service{
		Log:   l,
		Store: store,
	}
	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Printf("💥💥 ERROR: 'in customHTTPErrorHandler got error: %v'\n", err)
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		c.Logger().Error(err)
		errorPage := fmt.Sprintf("%s/%d.html", webRootDir, code)
		res, err := content.ReadFile(errorPage)
		if err != nil {
			log.Printf("💥💥 ERROR: 'in  content.ReadFile(%s) got error: %v'\n", errorPage, err)
		}
		if err := c.HTMLBlob(code, res); err != nil {
			log.Printf("💥💥 ERROR: 'in  c.HTMLBlob(%d, %s) got error: %v'\n", code, res, err)
			c.Logger().Error(err)
		}
	}
	var contentHandler = echo.WrapHandler(http.FileServer(http.FS(content)))

	// The embedded files will all be in the '/web' folder so need to rewrite the request (could also do this with fs.Sub)
	var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/web/$1"})
	//TODO  find a correct way to handle 404 in next handler, for now  is not used if we get /toto (only if method is not get)
	e.GET("/*", contentHandler, contentRewrite)
	// here the routes defined in OpenApi users.yaml are registered
	users.RegisterHandlers(e, &myUsersApi)
	// add another route for maxId
	e.GET("/users/maxid", myUsersApi.GetMaxId)

	myServer := GoHttpServer{
		listenAddress: listenAddress,
		logger:        l,
		store:         store,
		e:             e,
		router:        myServerMux,
		startTime:     time.Now(),
		httpServer: http.Server{
			Addr:         listenAddress,       // configure the bind address
			ErrorLog:     l,                   // set the logger for the server
			ReadTimeout:  defaultReadTimeout,  // max time to read request from the client
			WriteTimeout: defaultWriteTimeout, // max time to write response to the client
			IdleTimeout:  defaultIdleTimeout,  // max time for connections using TCP Keep-Alive
		},
	}
	//myServer.routes()

	return &myServer
}

// (*GoHttpServer) routes initializes all the handlers paths of this web server, it is called inside the NewGoHttpServer constructor
func (s *GoHttpServer) routes() {
	// s.router.Handle("/", s.getMyDefaultHandler())
	// s.e.GET("/readiness", s.getReadinessHandler())
	// s.e.GET("/health", s.getHealthHandler())

	//s.router.Handle("/hello", s.getHelloHandler())
}

// StartServer initializes all the handlers paths of this web server, it is called inside the NewGoHttpServer constructor
func (s *GoHttpServer) StartServer() error {

	// Starting the web server in his own goroutine
	go func() {
		s.logger.Printf("INFO: Starting http server listening at %s://%s/", defaultProtocol, s.listenAddress)
		err := s.e.StartServer(&s.httpServer)
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("💥💥 ERROR: 'Could not listen on %q: %s'\n", s.listenAddress, err)
		}
	}()
	s.logger.Printf("Server listening on : %s PID:[%d]", s.httpServer.Addr, os.Getpid())

	// Graceful Shutdown on SIGINT (interrupt)
	waitForShutdownToExit(&s.httpServer, secondsShutDownTimeout)
	return nil
}