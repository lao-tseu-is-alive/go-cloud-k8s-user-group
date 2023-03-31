package main

import (
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goserver"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/metadata"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/version"
	"runtime"
	"strings"
)

const (
	defaultPort                = 8080
	defaultDBPort              = 5432
	defaultDBIp                = "127.0.0.1"
	defaultDBSslMode           = "prefer"
	defaultWebRootDir          = "goCloudK8sUserGroupFront/dist/"
	defaultSqlDbMigrationsPath = "db/migrations"
	charsetUTF8                = "charset=UTF-8"
	MIMEAppJSON                = "application/json"
	MIMEHtml                   = "text/html"
	MIMEAppJSONCharsetUTF8     = MIMEAppJSON + "; " + charsetUTF8
	MIMEHtmlCharsetUTF8        = MIMEHtml + "; " + charsetUTF8
)

var db database.DB

// content holds our static web server content.
//
//go:embed goCloudK8sUserGroupFront/dist/*
var content embed.FS

// sqlMigrations holds our db migrations sql files using https://github.com/golang-migrate/migrate
// in the line above you SHOULD have the same path  as const defaultSqlDbMigrationsPath
//
//go:embed db/migrations/*.sql
var sqlMigrations embed.FS

func main() {
	// l := log.New(os.Stdout, fmt.Sprintf("%s ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	prefix := fmt.Sprintf("%s ", version.APP)
	l, _ := golog.NewLogger("dev", golog.DebugLevel, prefix)
	l.Info("Starting %s v:%s  rev:%s  build: %s'", version.APP, version.VERSION, version.REVISION, version.BuildStamp)
	l.Info("Repository url: https://%s'", version.REPOSITORY)
	secret, err := config.GetJwtSecretFromEnv()
	if err != nil {
		l.Fatal("💥💥 error doing config.GetJwtSecretFromEnv(). error: %v'\n", err)
	}
	tokenDuration, err := config.GetJwtDurationFromEnv(60)
	if err != nil {
		l.Fatal("💥💥 error doing config.GetJwtDurationFromEnv(60). error: %v'\n", err)
	}
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	if err != nil {
		l.Fatal("💥💥 error doing config.GetPgDbDsnUrlFromEnv(). error: %v\n", err)
	}
	db, err = database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatal("💥💥 error doing database.GetInstance(pgx, dbDsn:%v). error:%v\n", dbDsn, err)
	}
	defer db.Close()

	metadataService := metadata.Service{
		Log: l,
		Db:  db,
	}

	err = metadataService.CreateMetadataTableIfNeeded()
	if err != nil {
		l.Fatal("💥💥 error doing metadataService.CreateMetadataTableIfNeeded  error: %v", err)
	}

	found, ver, err := metadataService.GetServiceVersion(version.APP)
	if err != nil {
		l.Fatal("💥💥 error doing metadataService.CreateMetadataTableIfNeeded  error: %v\n", err)
	}
	if found {
		l.Info("service %s was found in metadata with version: %s", version.APP, ver)
	} else {
		l.Info("service %s was not found in metadata", version.APP)
	}
	err = metadataService.SetServiceVersion(version.APP, version.VERSION)
	if err != nil {
		l.Fatal("💥💥 error doing metadataService.SetServiceVersion  error: %v\n", err)
	}
	// example of go-migrate db migration with embed files in go program
	// https://github.com/golang-migrate/migrate
	// https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
	d, err := iofs.New(sqlMigrations, defaultSqlDbMigrationsPath)
	if err != nil {
		l.Fatal("💥💥 error doing iofs.New for db migrations  error: %v\n", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, strings.Replace(dbDsn, "postgres", "pgx", 1))
	if err != nil {
		l.Fatal("💥💥 error doing migrate.NewWithSourceInstance(iofs, dbURL:%s)  error: %v\n", dbDsn, err)
	}

	err = m.Up()
	if err != nil {
		//if err == m.
		if err != migrate.ErrNoChange {
			l.Fatal("💥💥 error doing migrate.Up error: %v\n", err)
		}
	}

	userStore, err := users.GetStorageInstance("pgx", db, l)
	if err != nil {
		l.Fatal("💥💥 error doing users.GetStorageInstance(postgres, dbDsn  : %v\n", err)
	}

	userService := users.Service{
		Log:         l,
		Store:       userStore,
		JwtSecret:   []byte(secret),
		JwtDuration: tokenDuration,
	}

	listenAddr, err := config.GetPortFromEnv(defaultPort)
	if err != nil {
		l.Fatal("💥💥 ERROR: 'calling GetPortFromEnv got error: %v'\n", err)
	}
	l.Info("Will start HTTP server listening on port %s'", listenAddr)
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
		l.Fatal("💥💥 ERROR: 'calling echo.Start(%s) got error: %v'\n", listenAddr, err)
	}

}
