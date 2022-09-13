package main

import (
	"embed"
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/goserver"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/version"
	"log"
	"os"
)

const (
	defaultPort       = 8080
	defaultDBPort     = 5432
	defaultDBIp       = "127.0.0.1"
	defaultDBSslMode  = "prefer"
	defaultWebRootDir = "goCloudK8sUserGroupFront/dist/"
)

// content holds our static web server content.
//
//go:embed goCloudK8sUserGroupFront/dist/*
var content embed.FS

func main() {
	l := log.New(os.Stdout, fmt.Sprintf("%s ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("INFO: 'Starting %s v:%s  rev:%s  build: %s'", version.APP, version.VERSION, version.REVISION, version.BuildStamp)
	l.Printf("INFO: 'Repository url: https://%s'", version.REPOSITORY)
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		tools.ToSnakeCase(version.APP), tools.ToSnakeCase(version.APP), defaultDBSslMode)
	if err != nil {
		l.Fatalf("💥💥 error doing config.GetPgDbDsnUrlFromEnv. error: %v\n", err)
	}
	s, err := users.GetStorageInstance("postgres", dbDsn, l)
	if err != nil {
		l.Fatalf("💥💥 error doing users.GetStorageInstance(postgres, dbDsn  : %v\n", err)
	}
	defer s.Close()

	listenAddr, err := config.GetPortFromEnv(defaultPort)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'calling GetPortFromEnv got error: %v'\n", err)
	}
	l.Printf("INFO: 'Will start HTTP server listening on port %s'", listenAddr)
	server := goserver.NewGoHttpServer(listenAddr, l, s, defaultWebRootDir, content)
	err = server.StartServer()
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'calling echo.Start(%s) got error: %v'\n", listenAddr, err)
	}

}
