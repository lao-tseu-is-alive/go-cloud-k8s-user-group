package main

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/version"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort      = 8080
	defaultDBPort    = 5432
	defaultDBIp      = "127.0.0.1"
	defaultDBSslMode = "prefer"
	webRootDir       = "web"
)

// content holds our static web server content.
//go:embed web/*
var content embed.FS

var contentHandler = echo.WrapHandler(http.FileServer(http.FS(content)))

// The embedded files will all be in the '/web' folder so need to rewrite the request (could also do this with fs.Sub)
var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/web/$1"})

func customHTTPErrorHandler(err error, c echo.Context) {
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

// GetNewServer initialize a new Echo server and returns it
func GetNewServer(l *log.Logger, store users.Storage) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	myUsersApi := users.Service{
		Log:   l,
		Store: store,
	}
	e.HideBanner = true
	e.HTTPErrorHandler = customHTTPErrorHandler
	//TODO  find a correct way to handle 404 in next handler, for now  is not used if we get /toto (only if method is not get)
	e.GET("/*", contentHandler, contentRewrite)
	/*
		webRootDirPath, err := filepath.Abs(webRootDir)
		if err != nil {
			log.Fatalf("Problem getting absolute path of directory: %s\nError:\n%v\n", webRootDir, err)
		}
		if _, err := os.Stat(webRootDirPath); os.IsNotExist(err) {
			log.Fatalf("The webRootDir parameter is wrong, %s is not a valid directory\nError:\n%v\n", webRootDirPath, err)
		}
		l.Printf("Using live mode serving from %s", webRootDirPath)
		e.Static("/", webRootDirPath)
	*/
	// here the routes defined in OpenApi users.yaml are registered
	users.RegisterHandlers(e, &myUsersApi)
	// add another route for maxId
	e.GET("/users/maxid", myUsersApi.GetMaxId)
	return e
}

func main() {
	l := log.New(os.Stdout, fmt.Sprintf("%s ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("INFO: 'Starting %s v:%s  rev:%s  build: %s'", version.APP, version.VERSION, version.REVISION, version.BuildStamp)
	l.Printf("INFO: 'Repository url: https://%s'", version.REPOSITORY)
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		version.APP, version.APP, defaultDBSslMode)
	if err != nil {
		log.Fatalf("💥💥 error doing config.GetPgDbDsnUrlFromEnv. error: %v\n", err)
	}
	l.Printf("INFO: 'dbDsn: %s'", dbDsn)
	s, err := users.GetStorageInstance("postgres", dbDsn, l)
	if err != nil {
		l.Fatalf("💥💥 error doing users.GetStorageInstance(postgres, dbDsn  : %v\n", err)
	}
	defer s.Close()

	listenAddr, err := config.GetPortFromEnv(defaultPort)
	if err != nil {
		log.Fatalf("💥💥 ERROR: 'calling GetPortFromEnv got error: %v'\n", err)
	}
	l.Printf("INFO: 'Will start HTTP server listening on port %s'", listenAddr)
	e := GetNewServer(l, s)
	err = e.Start(listenAddr)
	if err != nil {
		log.Fatalf("💥💥 ERROR: 'calling echo.Start(%s) got error: %v'\n", listenAddr, err)
	}

}
