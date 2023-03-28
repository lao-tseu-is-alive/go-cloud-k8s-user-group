package main

import (
	"embed"
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goserver"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/version"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"log"
	"os"
	"runtime"
)

const (
	defaultPort            = 8080
	defaultDBPort          = 5432
	defaultDBIp            = "127.0.0.1"
	defaultDBSslMode       = "prefer"
	defaultWebRootDir      = "goCloudK8sUserGroupFront/dist/"
	charsetUTF8            = "charset=UTF-8"
	MIMEAppJSON            = "application/json"
	MIMEHtml               = "text/html"
	MIMEAppJSONCharsetUTF8 = MIMEAppJSON + "; " + charsetUTF8
	MIMEHtmlCharsetUTF8    = MIMEHtml + "; " + charsetUTF8
)

var dbConn database.DB

// content holds our static web server content.
//
//go:embed goCloudK8sUserGroupFront/dist/*
var content embed.FS

func main() {
	l := log.New(os.Stdout, fmt.Sprintf("%s ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("INFO: 'Starting %s v:%s  rev:%s  build: %s'", version.APP, version.VERSION, version.REVISION, version.BuildStamp)
	l.Printf("INFO: 'Repository url: https://%s'", version.REPOSITORY)
	secret, err := config.GetJwtSecretFromEnv()
	if err != nil {
		l.Fatalf("💥💥 error doing config.GetJwtSecretFromEnv(). error: %v'\n", err)
	}
	tokenDuration, err := config.GetJwtDurationFromEnv(60)
	if err != nil {
		l.Fatalf("💥💥 error doing config.GetJwtDurationFromEnv(60). error: %v'\n", err)
	}
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	if err != nil {
		l.Fatalf("💥💥 error doing config.GetPgDbDsnUrlFromEnv(). error: %v\n", err)
	}
	dbConn, err = database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatalf("💥💥 error doing database.GetInstance(pgx, dbDsn:%v). error:%v\n", dbDsn, err)
	}
	defer dbConn.Close()

	userStore, err := users.GetStorageInstance("pgx", dbConn, l)
	if err != nil {
		l.Fatalf("💥💥 error doing users.GetStorageInstance(postgres, dbDsn  : %v\n", err)
	}

	userService := users.Service{
		Log:         l,
		Store:       userStore,
		JwtSecret:   []byte(secret),
		JwtDuration: tokenDuration,
	}

	listenAddr, err := config.GetPortFromEnv(defaultPort)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'calling GetPortFromEnv got error: %v'\n", err)
	}
	l.Printf("INFO: 'Will start HTTP server listening on port %s'", listenAddr)
	server := goserver.NewGoHttpServer(listenAddr, l, defaultWebRootDir, content, "/api")
	e := server.GetEcho()
	e.GET("/login", userService.GetLogin)
	e.POST("/login", userService.LoginUser)
	e.GET("/resetpassword", userService.GetResetPasswordEmail)
	e.POST("/resetpassword", userService.SendResetPassword)
	e.GET("/resetpassword/:resetPasswordToken", userService.GetResetPasswordToken)
	e.POST("/resetpassword/:resetPasswordToken", userService.ResetPassword)
	r := server.GetRestrictedGroup()
	users.RegisterHandlers(r, &userService) // register all openapi declared routes
	//some other restricted routes
	r.GET("/status", userService.GetStatus)
	r.GET("/users/maxid", userService.GetMaxId)

	err = server.StartServer()
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'calling echo.Start(%s) got error: %v'\n", listenAddr, err)
	}

}
